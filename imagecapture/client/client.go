package client

import (
	"os"

	"github.com/veyanrech/homeWebCamera/config"
)

type Client struct {
	certs []string
	conf  config.Config
	q     *RoundBufferQueue
}

func NewClient(c config.Config) *Client {
	return &Client{
		conf: c,
		q:    NewRoundBufferQueue(5),
	}
}

func UploadCerts(c *Client, conf config.Config) error {
	return nil
}

func (c *Client) Run() {
	// Run the client
	go func() {
		for {
			// Load file to queue
			c.loadFileToQueue()

			// Send file
			c.sendFile()

			// Remove file
			c.removeFile()
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

func (c *Client) sendFile() {}

func (c *Client) removeFile() {}
