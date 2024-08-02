package utils

import (
	"runtime"
	"time"
)

func GenerateFilename(additional string) string {

	timeNow := time.Now().Format("2006-01-02-15:04:05.000")

	return timeNow + "-" + additional

}

func GetOS() string {
	//get the OS type
	return runtime.GOOS
}
