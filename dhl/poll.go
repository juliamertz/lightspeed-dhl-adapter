package dhl

import (
	"fmt"
	"time"
)

const (
	every = 5 // minutes
)

func StartPolling() {
	go func() {
		for {
			fmt.Println("Polling...")
			time.Sleep(every * time.Minute)
		}
	}()
}
