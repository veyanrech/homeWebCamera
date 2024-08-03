package utils

import (
	"runtime"
	"time"
)

func GenerateFilename(additional string) string {

	timeNow := time.Now()
	firstpart := timeNow.Format("2006_01_02_15_04_05")
	secondpart := timeNow.Format(".000")

	firstandsecond := firstpart + "_" + secondpart[1:]

	return firstandsecond + "_" + additional

}

func GetOS() string {
	//get the OS type
	return runtime.GOOS
}
