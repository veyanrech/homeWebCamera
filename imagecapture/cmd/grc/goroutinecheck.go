package main

import (
	"fmt"
	"time"
)

func main() {
	fmt.Println("start")

	go func() {
		fmt.Println("goroutine")

		defer func() {
			if r := recover(); r != nil {
				fmt.Println("recovered")
			}
		}()

		counter := 0

		for {
			counter++
			time.Sleep(1 * time.Second)

			for i := 0; i < 10; i++ {
				fmt.Println(counter, i)
			}

			// panic("panic")
		}
	}()

	time.Sleep(10 * time.Second)
	fmt.Println("end")
}
