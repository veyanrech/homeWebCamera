package utils

import (
	"fmt"
	"os"
	"time"
)

type Logger interface {
	Info(string)
	Error(string)
}

type consoleLogger struct {
}

func NewConsoleLogger() Logger {
	return &consoleLogger{}
}

func (cl *consoleLogger) Info(msg string) {
	fmt.Println("INFO:", msg)
}

func (cl *consoleLogger) Error(msg string) {
	fmt.Println("ERROR:", msg)
}

type fileLogger struct {
	filepathInfo  *os.File
	filepathError *os.File
}

func NewFileLogger(infof, errorf string) Logger {
	res := &fileLogger{}

	path, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	path += "/logs"

	err = os.MkdirAll(path, 0755)
	if err != nil {
		panic(err)
	}

	f, err := os.OpenFile(path+"/"+infof, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		panic(err)
	}
	res.filepathInfo = f

	f, err = os.OpenFile(path+"/"+errorf, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		panic(err)
	}
	res.filepathError = f

	return res
}

func (fl *fileLogger) Info(msg string) {
	logtime := time.Now().Format("01.02.2006 15:04:05.000")
	fmt.Fprintf(fl.filepathInfo, "%s: %s\n", logtime, msg)
}

func (fl *fileLogger) Error(msg string) {
	logtime := time.Now().Format("01.02.2006 15:04:05.000")
	fmt.Fprintf(fl.filepathError, "%s: %s\n", logtime, msg)
}
