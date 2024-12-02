// Package main, Audio volume normalizer main entry point
package main

import (
	"fmt"

	"github.com/olbrichattila/audionorm/internal/commandline"
	"github.com/olbrichattila/audionorm/internal/normalizer"
)

func main() {
	folder, factor, help, err := commandline.GetCommandLineParams()
	if err != nil {
		fmt.Println(err.Error())
		displayHelp()
		return
	}

	if help {
		displayHelp()
		return
	}

	normalizer := normalizer.New(func(message string) {
		fmt.Println(message)
	})

	normalizer.Normalize(folder, factor)
}

func displayHelp() {
	fmt.Println(
		`Usage:
  audionorm <path> -factor=0.8 -help
    Where factor is a number between 0 and 1 and help displays this message
    Non of the parameters are mandatory.
	Examples:
	audionorm (uses current directory and factor is 1)
	audionorm -factor=0.8 (uses current directory and factor is 0.8)
	audionorm ./myfolder (uses "myfolder" directory and factor is 1)
	audionorm ./myfolder -factor=0.8 (uses "myfolder" directory and factor is 1)
	audionorm ./myfolder -factor=0.8 -help (Display help only, and quit)
	audionorm -help (display help only)
	`)
}
