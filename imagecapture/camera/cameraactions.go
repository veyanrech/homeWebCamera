package camera

import (
	"os"
	"runtime"
)

type Camera interface {
	TakePicture() error
}

func getOS() string {
	//get the OS type
	return runtime.GOOS
}

func NewCameraByOS(dn []string) Camera {

	picdir := createFolder()

	switch os := getOS(); os {
	case "darwin":
		return NewMacOSCamera(dn, picdir)
	case "windows":
		return NewWinCamera(dn, picdir)
	}

	return nil
}

func createFolder() string {
	//get the location where the program is running
	dir, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	folder := dir + "/pictures"

	//create folder
	err = os.MkdirAll(folder, os.ModePerm)
	if err != nil {
		panic(err)
	}

	return folder
}
