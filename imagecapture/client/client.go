package client

import (
	"bytes"
	"context"
	"crypto/tls"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/veyanrech/homeWebCamera/imagecapture/config"
	"github.com/veyanrech/homeWebCamera/imagecapture/utils"
)

type Client struct {
	certs []string
	conf  config.Config
	q     *RoundBufferQueue
	l     utils.Logger
}

func NewClient(c config.Config, l utils.Logger) *Client {
	return &Client{
		conf: c,
		q:    NewRoundBufferQueue(5),
		l:    l,
	}
}

func UploadCerts(c *Client, conf config.Config) error {
	return nil
}

func (c *Client) Run() {
	// Run the client
	go func() {
		ticker := time.NewTicker(time.Duration(c.conf.GetInt("send_file_interval_min")) * time.Minute)
		signalChannel := make(chan os.Signal, 1)
		signal.Notify(signalChannel, os.Interrupt)

		deleteQueue := NewRoundBufferQueue(c.conf.GetInt("devices_count"))

		allowtodeleteChannel := make(chan bool, 1)

		for {
			select {
			case <-ticker.C:

				// Load file to queue
				c.loadFileToQueue()

				// Send file
				c.sendFile(deleteQueue, allowtodeleteChannel)

				// Remove file
				go func() {
					<-allowtodeleteChannel
					c.removeFile(deleteQueue)
				}()

			case <-signalChannel:
				os.Exit(0)
			}
		}
	}()
}

func (c *Client) loadFileToQueue() {
	fs, err := os.ReadDir(c.conf.GetString("pictures_folder"))
	if err != nil {
		panic(err)
	}

	for _, f := range fs {
		if !f.IsDir() {
			// Add file to queue
			c.l.Info(fmt.Sprint("Adding file to queue: ", f.Name()))
			c.q.Add(f.Name())
		}
	}
}

func (c *Client) sendFile(delqu *RoundBufferQueue, allowToDeleteCh chan bool) {

	if c.q.UnprocessedLen() < 2 {
		return
	}

	formdataBody := bytes.Buffer{}
	writer := multipart.NewWriter(&formdataBody)
	tempqueue := NewRoundBufferQueue(c.conf.GetInt("devices_count"))
	counter := 0
	v, ok := c.q.Get()
	for ok {

		c.l.Info(fmt.Sprint("Sending file: ", v))

		tempqueue.Add(v)

		counter++

		//check if file exists
		fopen, err := os.Open(c.conf.GetString("pictures_folder") + "/" + v)
		if err != nil {
			c.l.Error(fmt.Sprint("Error reading file: ", err))
			continue
		}

		//add content to formdata
		filepart, err := writer.CreateFormFile(fmt.Sprintf("file%d", counter), v)
		if err != nil {
			c.l.Error(fmt.Sprint("Error creating form file: ", err))
			fopen.Close()
			continue
		}

		_, err = io.Copy(filepart, fopen)
		if err != nil {
			c.l.Error(fmt.Sprint("Error copying file content: ", err))
			fopen.Close()
			continue
		}

		//send file
		if counter != c.conf.GetInt("devices_count") {
			fopen.Close()
			v, ok = c.q.Get()
			continue
		}

		writer.Close()

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		// req, err := http.NewRequest("POST", c.conf.GetString("photos_receiver_url"), &formdataBody)
		req, err := http.NewRequestWithContext(ctx, "POST", c.conf.GetString("photos_receiver_url"), &formdataBody)
		if err != nil {
			c.l.Error(fmt.Sprint("Error creating request: ", err))
			fopen.Close()
			continue
		}

		req.Header.Set("Content-Type", writer.FormDataContentType())
		req.Header.Set("X-Chat-Registration-Token", c.conf.GetString("registered_chat_token"))

		//send files
		client := http.Client{
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{
					InsecureSkipVerify: true, // Disable certificate verification
				},
			},
		}
		resp, err := client.Do(req)
		if err != nil {
			c.l.Error(fmt.Sprint("Error sending file: ", err))
			fopen.Close()
			continue
		}

		if resp.StatusCode != http.StatusOK {
			c.l.Error(fmt.Sprint("Error sending file: ", resp.Status))
			fopen.Close()

			//return files to queue
			v2, ok2 := tempqueue.Get()
			for ok2 {
				c.q.Add(v2)
				v2, ok2 = tempqueue.Get()
			}

			continue
		} else {

			//add files to delete queue
			v2, ok2 := tempqueue.Get()
			for ok2 {
				c.l.Info(fmt.Sprint("File sent and ready to be deleted: ", v2))
				delqu.Add(v2)
				v2, ok2 = tempqueue.Get()
			}

			go func() {
				allowToDeleteCh <- true
			}()
		}

		fopen.Close()

		v, ok = c.q.Get()
	}

}

func (c *Client) removeFile(delqu *RoundBufferQueue) {
	v, ok := delqu.Get()
	for ok {
		c.l.Info(fmt.Sprint("Removing file: ", v))
		err := os.Remove(c.conf.GetString("pictures_folder") + string(os.PathSeparator) + v)
		if err != nil {
			c.l.Error(fmt.Sprint("Error removing file: ", err))
		}
		v, ok = delqu.Get()
	}
}
