package utils

import (
	"fmt"
	"os"
)

func Assert(statement bool, msg string) {
	if statement {
		fmt.Printf("Assertion failed: %s", msg)
		os.Exit(1)
	}
}
