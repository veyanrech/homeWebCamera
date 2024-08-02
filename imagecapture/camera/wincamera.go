package camera

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/veyanrech/homeWebCamera/config"
	"github.com/veyanrech/homeWebCamera/utils"
)

type winCamera struct {
	DevicesNames   []string
	PicturesFolder string
	conf           config.Config
}

func NewWinCamera(picdir string, c config.Config) Camera {

	dn := c.GetSliceOfStrings("devices")
	if dn == nil {
		dn = askForMacOsDevicesNames()
	}

	if dn == nil {
		fmt.Println("No devices names entered")
		return nil
	}

	return &winCamera{
		DevicesNames:   dn,
		PicturesFolder: picdir,
		conf:           c,
	}
}

func (c *winCamera) TakePicture() error {
	ffmpegCommandFormat := "ffmpeg -f dshow -i video=\"%s\" -vframes 1 %s/%s"

	for _, v := range c.DevicesNames {
		finalCommand := fmt.Sprintf(ffmpegCommandFormat, v, c.PicturesFolder, utils.GenerateFilename("output.jpg"))
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

func askForWinDevicesNames() []string {
	fmt.Println("Here is the list of devices:")
	//list the devices
	cmdToList := "ffmpeg -list_devices true -f dshow -i dummy"
	err := exec.Command("sh", "-c", cmdToList).Run()
	if err != nil {
		fmt.Println("Error listing the devices")
		return nil
	}
	fmt.Println("Enter the device indexes names separated by comma")
	var devicesNames string
	fmt.Scanln(&devicesNames)

	if devicesNames == "" {
		fmt.Println("No devices names entered")
		return nil
	}

	return strings.Split(devicesNames, ",")
}
