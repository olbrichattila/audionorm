// Package main, Audio volume normalizer main entry point
package main

import (
	"fmt"

	"github.com/olbrichattila/audionorm/internal/commandline"
	"github.com/olbrichattila/audionorm/internal/normalizer"
)

func main() {
	folder, factor, tolerance, help, err := commandline.GetCommandLineParams()
	if err != nil {
		fmt.Println("Error: \033[31m" + err.Error() + "\033[0m")
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

	normalizer.Normalize(folder, factor, tolerance)
}

func displayHelp() {
	fmt.Println(
		`Usage:
audionorm <path> -factor=<value> -tolerance=<value> -help
Description:
<path>: Specifies the folder containing audio files to process. Defaults to the current directory if not provided.
-factor=<value>: A number between 0 and 1 that defines the normalization factor. Defaults to 1 if not specified.
-tolerance=<value>: Specifies the over-amplification tolerance, a number between 0 and 20. If set to 0 or omitted, over-amplification is disabled.
-help: Displays this help message and exits.
Notes:

All parameters are optional.
If -help is specified, the program will display this message and exit without performing any processing.
Examples:

1. audionorm
   Uses the current directory, normalization factor of 1, and disables over-amplification.

2. audionorm -factor=0.8
   Uses the current directory, with a normalization factor of 0.8.

3. audionorm ./myfolder
   Processes the "myfolder" directory with a normalization factor of 1.

4. audionorm ./myfolder -factor=0.8
   Processes the "myfolder" directory with a normalization factor of 0.8.

5. audionorm ./myfolder -factor=0.8 -help
   Displays this help message and exits without processing any files.

6. audionorm -tolerance=2
   Uses the current directory, with over-amplification tolerance set to 2.

7. audionorm -help
   Displays this help message and exits.`)
}
