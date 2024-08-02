package camera

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/veyanrech/homeWebCamera/config"
	"github.com/veyanrech/homeWebCamera/utils"
)

const macosCommand = "ffmpeg -f avfoundation -list_devices true -i \"\""

type macOSCamera struct {
	DevicesNames   []string
	PicturesFolder string
	conf           config.Config
}

func NewMacOSCamera(folder string, c config.Config) Camera {

	dn := c.GetSliceOfStrings("devices")
	if dn == nil {
		dn = askForMacOsDevicesNames()
	}

	if dn == nil {
		fmt.Println("No devices names entered")
		return nil
	}

	return &macOSCamera{
		DevicesNames:   dn,
		PicturesFolder: folder,
		conf:           c,
	}
}

func (c *macOSCamera) TakePicture() error {

	for _, v := range c.DevicesNames {
		err := c.runTakePictureMacosCommand(v, c.PicturesFolder+"/"+utils.GenerateFilename("output.jpg"))
		if err != nil {
			return err
		}
	}

	return nil
}

func (c *macOSCamera) runTakePictureMacosCommand(devicename string, picturesPlace string) error {
	cmd := exec.Command("sh", "-c", "echo $(ffmpeg -f avfoundation -framerate 30 -video_size 1280x720 -i "+devicename+" -frames:v 1 "+picturesPlace+")")

	o, err := cmd.CombinedOutput()

	if err != nil {
		return err
	}

	fmt.Println(string(o))

	return nil
}

func runBrowseDevicesMacosCommand() error {
	// cmd := exec.Command("ffmpeg", "-f avfoundation", "-list_devices true", "-i \"\"")
	// cmd := exec.Command("zsh", "-c", "$(\"ffmpeg -f avfoundation\")")
	// cmd := exec.Command("ls", "-la")
	cmd := exec.Command("sh", "-c", "echo $(ffmpeg -f avfoundation -list_devices true -i \"\")")

	output, err := cmd.CombinedOutput()

	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	fmt.Println(string(output))

	return nil
}

func askForMacOsDevicesNames() []string {
	fmt.Println("Here is the list of devices:")
	//list the devices

	err := runBrowseDevicesMacosCommand()
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
