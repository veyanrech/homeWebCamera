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
	l              utils.Logger
}

func NewWinCamera(picdir string, c config.Config, l utils.Logger) Camera {

	dn := c.GetSliceOfStrings("devices")
	if dn == nil {
		dn = askForWinDevicesNames()
	}

	if dn == nil {
		fmt.Println("No devices names entered")
		return nil
	}

	return &winCamera{
		DevicesNames:   dn,
		PicturesFolder: picdir,
		conf:           c,
		l:              l,
	}
}

func (c *winCamera) TakePicture() error {
	ffmpegCommandFormat := "ffmpeg -f dshow -i video=\"%s\" -vframes 1 %s\\%s"

	cmdSPlit := strings.Split(ffmpegCommandFormat, " ")

	for _, v := range c.DevicesNames {

		// finalCommand := fmt.Sprintf(ffmpegCommandFormat, v, c.PicturesFolder, utils.GenerateFilename("output.jpg"))
		cmdSPlit[4] = fmt.Sprintf("video=\"%s\"", v)
		cmdSPlit[6] = fmt.Sprintf("%s\\%s", c.PicturesFolder, utils.GenerateFilename("output.jpg"))

		cmdSPlitcopy := make([]string, len(cmdSPlit))
		copy(cmdSPlitcopy, cmdSPlit)

		err := c.runWinCommandsplitted(cmdSPlitcopy)
		if err != nil {
			c.l.Error(err.Error())
			return err
		}
	}

	return nil

}

func (c *winCamera) runWinCommand(command string) error {
	cmd := exec.Command("cmd", c.PicturesFolder, command)

	b, err := cmd.CombinedOutput()

	c.l.Info(string(b))

	if err != nil {
		return err
	}
	return nil
}

func (c *winCamera) runWinCommandsplitted(command []string) error {
	cmd := exec.Command(command[0], command[1:]...)

	b, err := cmd.CombinedOutput()

	c.l.Info(string(b))

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
