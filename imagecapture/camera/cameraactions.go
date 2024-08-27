package camera

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"runtime"
	"strings"
	"syscall"
	"time"

	"github.com/veyanrech/homeWebCamera/imagecapture/config"
	"github.com/veyanrech/homeWebCamera/imagecapture/utils"
)

type Camera interface {
	TakePicture() error
}

func NewCameraByOS(c config.Config, l utils.Logger) Camera {

	picdir := createFolder()

	//clean content of the folder, but not the folder itself
	folderContents, err := os.ReadDir(picdir)
	if err != nil {
		panic(err)
	}

	for _, file := range folderContents {
		err = os.Remove(picdir + string(os.PathSeparator) + file.Name())
		if err != nil {
			panic(err)
		}
	}

	c.Set("pictures_folder", picdir)

	switch os := utils.GetOS(); os {
	case "darwin":
		return NewMacOSCamera(picdir, c, l)
	case "windows":
		return NewWinCamera(picdir, c, l)
	}

	return nil
}

type CameraService struct {
	cam      Camera
	interval int
	l        utils.Logger
}

func NewCameraService(cam Camera, c config.Config, l utils.Logger) *CameraService {

	//test if the camera is working
	err := cam.TakePicture()
	if err != nil {
		//one attempt to kill ffmpeg
		err = killFFMPEG()
		if err != nil {
			l.Error("Error killing ffmpeg process")
			panic(err)
		} else {
			//second attempt to take picture
			err = cam.TakePicture()
			if err != nil {
				l.Error("Error taking picture")
				panic(err)
			}
		}
	}

	return &CameraService{
		cam:      cam,
		interval: c.GetInt("take_picture_interval_min"),
		l:        l,
	}
}

func (cs *CameraService) TakePictureEvery() {
	go func() {
		signalChannel := make(chan os.Signal, 1)
		signal.Notify(signalChannel, syscall.SIGINT, syscall.SIGTERM)

		ticker := time.NewTicker(time.Duration(cs.interval) * time.Minute)

		for {
			select {
			case <-ticker.C:
				err := cs.TakePictureWithFail()
				if err != nil {
					panic(err)
				}
			case <-signalChannel:
				os.Exit(0)
			}
		}
	}()
}

func (cs *CameraService) TakePictureWithFail() error {
	err := cs.cam.TakePicture()
	if err != nil {
		//one attempt to kill ffmpeg
		err = killFFMPEG()
		if err != nil {
			cs.l.Error("Error killing ffmpeg process")
			return err
		}
		//second attempt to take picture
		err = cs.cam.TakePicture()
		if err != nil {
			cs.l.Error("Error taking picture")
			return err
		}
	}
	return nil
}

func killFFMPEG() error {
	switch os := utils.GetOS(); os {
	case "darwin", "linux":
		return killNixFFMPEG()
	case "windows":
		return killWinFFMPEG()
	}

	return errors.New("OS not supported")
}

func killNixFFMPEG() error {
	var psCommand, grepCommand *exec.Cmd
	psCommand = exec.Command("ps", "aux")
	grepCommand = exec.Command("grep", "ffmpeg")
	return KillFFMPEGBYID(psCommand, grepCommand)
}

func killWinFFMPEG() error {
	var psCommand, grepCommand *exec.Cmd
	psCommand = exec.Command("tasklist")
	grepCommand = exec.Command("findstr", "ffmpeg")
	return KillFFMPEGBYID(psCommand, grepCommand)
}

func KillFFMPEGBYID(psc, grepc *exec.Cmd) error {

	// Get the list of processes
	psOutput, err := psc.Output()
	if err != nil {
		return err
	}

	// Filter for ffmpeg process
	grepc.Stdin = strings.NewReader(string(psOutput))
	grepOutput, err := grepc.Output()
	if err != nil {
		return err
	}

	// Parse and kill the process
	for _, line := range strings.Split(string(grepOutput), "\n") {
		if strings.Contains(line, "ffmpeg") {
			fields := strings.Fields(line)
			if len(fields) > 1 {
				pid := fields[1]
				return killProcess(pid)
			}
		}
	}
	return nil

}

func killProcess(pid string) error {
	var killCommand *exec.Cmd
	switch os := runtime.GOOS; os {
	case "windows":
		killCommand = exec.Command("taskkill", "/F", "/PID", pid)
	case "darwin", "linux":
		killCommand = exec.Command("kill", "-9", pid)
	default:
		return fmt.Errorf("unsupported operating system: %s", os)
	}

	if err := killCommand.Run(); err != nil {
		return err
	}
	return nil
}

func createFolder() string {
	//get the location where the program is running
	dir, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	folder := dir + string(os.PathSeparator) + "pictures"

	//create folder
	err = os.MkdirAll(folder, os.ModePerm)
	if err != nil {
		panic(err)
	}

	return folder
}
