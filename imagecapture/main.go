package main

import (
	"os"
	"os/signal"

	"github.com/veyanrech/homeWebCamera/imagecapture/camera"
	"github.com/veyanrech/homeWebCamera/imagecapture/client"
	"github.com/veyanrech/homeWebCamera/imagecapture/config"
	"github.com/veyanrech/homeWebCamera/imagecapture/utils"
)

func main() {
	//run Camera Capturing
	lggr := utils.NewFileLogger("info.log", "error.log")

	conf := config.NewConfig()

	cam := camera.NewCameraByOS(conf, lggr)

	if cam == nil {
		panic("Camera not found")
	}

	cs := camera.NewCameraService(cam, conf, lggr)

	cs.TakePictureEvery()

	camclient := client.NewClient(conf, lggr)

	camclient.Run()

	signalChannel := make(chan os.Signal, 1)

	signal.Notify(signalChannel, os.Interrupt)

	<-signalChannel

}
