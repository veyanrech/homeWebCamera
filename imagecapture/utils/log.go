package utils

import (
	"fmt"
	"os"
	"time"
)

type Logger interface {
	Info(string)
	Error(string)
	Disable() //don't want to log on google cloud
}

type consoleLogger struct {
	disable bool
}

func NewConsoleLogger() Logger {
	return &consoleLogger{
		disable: false,
	}
}

func (cl *consoleLogger) Disable() {
	cl.disable = true
}

func (cl *consoleLogger) Info(msg string) {
	if !cl.disable {
		fmt.Println("INFO:", msg)
	}
}

func (cl *consoleLogger) Error(msg string) {
	if !cl.disable {
		fmt.Println("ERROR:", msg)
	}
}

type fileLogger struct {
	filepathInfo  *os.File
	filepathError *os.File
	disable       bool
}

func NewFileLogger(infof, errorf string) Logger {
	res := &fileLogger{}

	path, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	path += string(os.PathSeparator) + "logs"

	err = os.MkdirAll(path, 0755)
	if err != nil {
		panic(err)
	}

	f, err := os.OpenFile(path+string(os.PathSeparator)+infof, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		panic(err)
	}
	res.filepathInfo = f

	f, err = os.OpenFile(path+string(os.PathSeparator)+errorf, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		panic(err)
	}
	res.filepathError = f

	return res
}

func (fl *fileLogger) Disable() {
	fl.disable = true
}

func (fl *fileLogger) Info(msg string) {
	if fl.disable {
		return
	}
	logtime := time.Now().Format("01.02.2006 15:04:05.000")
	fmt.Fprintf(fl.filepathInfo, "%s: %s\n", logtime, msg)
}

func (fl *fileLogger) Error(msg string) {
	if fl.disable {
		return
	}
	logtime := time.Now().Format("01.02.2006 15:04:05.000")
	fmt.Fprintf(fl.filepathError, "%s: %s\n", logtime, msg)
}
