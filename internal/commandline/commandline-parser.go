// Package commandline retrieves the command line parameters and verify them. Not using flag so this way any parameter order works
package commandline

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

// GetCommandLineParams returns path, normalization factor and if help needs to be displayed
func GetCommandLineParams() (string, float64, bool, error) {
	folder := getFolder()
	if !folderExists(folder) {
		return "", 0, false, fmt.Errorf("folder " + folder + " does not exists or not readable, or it is a file.")
	}

	factor := getFactor()
	if factor <= 0 || factor > 1 {
		return "", 0, false, fmt.Errorf("factor must be larger then 0 and smaller or equal then 1, example 0.8")
	}

	return folder, factor, getHelp(), nil
}

func getFolder() string {
	pars := getNonFlagParams()
	if len(pars) > 0 {
		return pars[0]
	}

	return "."
}

func getFactor() float64 {
	flagPars := getFlagParams()

	if v, ok := flagPars["-factor"]; ok {
		floatValue, err := strconv.ParseFloat(v, 64)
		if err != nil {
			return -1
		}

		return floatValue
	}

	return 1
}

func getHelp() bool {
	flagPars := getFlagParams()
	_, ok := flagPars["-help"]
	return ok
}

func folderExists(folderName string) bool {
	fileInfo, err := os.Stat(folderName)
	if err != nil {
		return false
	}

	if fileInfo.IsDir() {
		return true
	}

	return false
}

func getNonFlagParams() []string {
	res := make([]string, 0)
	for i := 1; i < len(os.Args); i++ {
		if !strings.HasPrefix(os.Args[i], "-") {
			res = append(res, os.Args[i])
		}
	}

	return res
}

func getFlagParams() map[string]string {
	res := make(map[string]string, 0)
	for i := 1; i < len(os.Args); i++ {
		if strings.HasPrefix(os.Args[i], "-") {
			parts := strings.Split(os.Args[i], "=")
			if len(parts) == 2 {
				res[parts[0]] = parts[1]
			} else {
				res[parts[0]] = ""
			}
		}
	}

	return res
}
