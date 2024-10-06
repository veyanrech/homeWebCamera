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

	var conf config.Config
	var filename string

	switch opsys := utils.GetOS(); opsys {
	case "darwin":
		filename = "." + string(os.PathSeparator) + "macos.config.json"
	case "windows":
		filename = "." + string(os.PathSeparator) + "win.config.json"
	}

	conf = config.NewConfig(filename)

	cam := camera.NewCameraByOS(conf, lggr)

	if cam == nil {
		panic("Camera not found") //no need to recover
	}

	cs := camera.NewCameraService(cam, conf, lggr)

	cs.TakePictureEvery()

	camclient := client.NewClient(conf, lggr)

	camclient.Run()

	signalChannel := make(chan os.Signal, 1)

	signal.Notify(signalChannel, os.Interrupt)

	<-signalChannel

}
