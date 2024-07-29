package camera

import (
	"fmt"
	"os/exec"
)

type winCamera struct {
	DevicesNames   []string
	PicturesFolder string
}

func NewWinCamera(dn []string, picdir string) Camera {
	return &winCamera{
		DevicesNames:   dn,
		PicturesFolder: picdir,
	}
}

func (c *winCamera) TakePicture() error {
	ffmpegCommandFormat := "ffmpeg -f dshow -i video=\"%s\" -vframes 1 %s/%s"

	for _, v := range c.DevicesNames {
		finalCommand := fmt.Sprintf(ffmpegCommandFormat, v, c.PicturesFolder, generateFilename("output.jpg"))
		err := c.runWinCommand(finalCommand)
		if err != nil {
			return err
		}
	}

	return nil

}

func (c *winCamera) runWinCommand(command string) error {
	err := exec.Command("cmd", c.PicturesFolder, command).Run()
	if err != nil {
		return err
	}
	return nil
}
