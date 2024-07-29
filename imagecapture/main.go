package main

import (
	"fmt"
	"strings"
)

func main() {
	// Code
	//ask for console user input to enter devices names
	dn := askForDevicesNames()

	if dn == nil {
		fmt.Println("Exitig program")
		return
	}

}

func askForDevicesNames() []string {
	fmt.Println("Enter the devices names separated by comma")
	var devicesNames string
	fmt.Scanln(&devicesNames)

	if devicesNames == "" {
		fmt.Println("No devices names entered")
		return nil
	}

	return strings.Split(devicesNames, ",")
}
