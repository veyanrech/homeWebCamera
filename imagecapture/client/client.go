package client

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/veyanrech/homeWebCamera/config"
	"github.com/veyanrech/homeWebCamera/utils"
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

		for {
			select {
			case <-ticker.C:

				// Load file to queue
				c.loadFileToQueue()

				// Send file
				c.sendFile(deleteQueue)

				// Remove file
				c.removeFile(deleteQueue)

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
			c.q.Add(f.Name())
		}
	}
}

func (c *Client) sendFile(delqu *RoundBufferQueue) {
	formdataBody := bytes.Buffer{}
	writer := multipart.NewWriter(&formdataBody)
	counter := 0
	v, ok := c.q.Get()
	for ok {

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
		if counter == c.conf.GetInt("devices_count") {

			writer.Close()

			req, err := http.NewRequest("POST", c.conf.GetString("bot_url"), &formdataBody)
			if err != nil {
				c.l.Error(fmt.Sprint("Error creating request: ", err))
				fopen.Close()
				continue
			}

			req.Header.Set("Content-Type", writer.FormDataContentType())
			req.Header.Set("X-Chat-Registration-Token", c.conf.GetString("registered_chat_token"))

			//send files
			client := http.Client{}
			resp, err := client.Do(req)
			if err != nil {
				c.l.Error(fmt.Sprint("Error sending file: ", err))
				fopen.Close()
				continue
			}

			if resp.StatusCode != http.StatusOK {
				c.l.Error(fmt.Sprint("Error sending file: ", resp.Status))
				fopen.Close()
				continue
			}

			fopen.Close()
		}

		fopen.Close()

		delqu.Add(v)

		v, ok = c.q.Get()
	}

}

func (c *Client) removeFile(delqu *RoundBufferQueue) {
	v, ok := delqu.Get()
	for ok {
		err := os.Remove(c.conf.GetString("pictures_folder") + string(os.PathSeparator) + v)
		if err != nil {
			c.l.Error(fmt.Sprint("Error removing file: ", err))
		}
		v, ok = delqu.Get()
	}
}
