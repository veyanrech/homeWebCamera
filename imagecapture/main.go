package main

import (
	"os"
	"os/signal"

	"github.com/veyanrech/homeWebCamera/camera"
	"github.com/veyanrech/homeWebCamera/config"
)

func main() {
	//run Camera Capturing

	conf := config.NewConfig()

	cam := camera.NewCameraByOS(conf)

	if cam == nil {
		panic("Camera not found")
	}

	cs := camera.NewCameraService(cam, conf)

	cs.TakePictureEvery()

	signalChannel := make(chan os.Signal, 1)

	signal.Notify(signalChannel, os.Interrupt)

	<-signalChannel

}
