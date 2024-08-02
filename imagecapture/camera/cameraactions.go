package camera

import (
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/veyanrech/homeWebCamera/config"
	"github.com/veyanrech/homeWebCamera/utils"
)

type Camera interface {
	TakePicture() error
}

func NewCameraByOS(c config.Config) Camera {

	picdir := createFolder()

	c.Set("pictures_folder", picdir)

	switch os := utils.GetOS(); os {
	case "darwin":
		return NewMacOSCamera(picdir, c)
	case "windows":
		return NewWinCamera(picdir, c)
	}

	return nil
}

type CameraService struct {
	cam      Camera
	interval int
}

func NewCameraService(cam Camera, c config.Config) *CameraService {
	return &CameraService{
		cam:      cam,
		interval: c.GetInt("take_picture_interval_min"),
	}
}

func (cs *CameraService) TakePictureEvery() {
	go func() {
		go func() {
			signalChannel := make(chan os.Signal, 1)
			signal.Notify(signalChannel, syscall.SIGINT, syscall.SIGTERM)

			ticker := time.NewTicker(time.Duration(cs.interval) * time.Minute)

			for {
				select {
				case <-ticker.C:
					err := cs.cam.TakePicture()
					if err != nil {
						panic(err)
					}
				case <-signalChannel:
					os.Exit(0)
				}
			}
		}()

	}()
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
