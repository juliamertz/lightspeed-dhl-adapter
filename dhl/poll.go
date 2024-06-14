package dhl

import (
	"fmt"
	// "jorismertz/lightspeed-dhl/database"
	"time"
)

const (
	every = 5 // minutes
)

func StartPolling() {
	go func() {
		for {
			// orders := database.
			// err := GetDrafts()
			// if err != nil {
			// 	fmt.Println(err)
			// }

			fmt.Println("Polling...")
			time.Sleep(every * time.Minute)
		}
	}()
}
