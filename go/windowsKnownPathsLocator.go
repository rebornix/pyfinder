package main

import (
	"os"
	"runtime"
	"strings"
)

func getSearchPathEnvVarNames() []string {
	os := runtime.GOOS
	if os == "windows" {
		return []string{"PATH", "Path"}
	}

	return []string{}
}

func getWindowsKnownDirs() []FileLocator {
	envVars := getSearchPathEnvVarNames()
	knownDirs := []FileLocator{}

	for _, envVar := range envVars {
		val, found := os.LookupEnv(envVar)
		if found {
			parts := strings.Split(val, string(os.PathListSeparator))

			for _, part := range parts {
				// fmt.Println(part)
				_, err := os.Stat(part)
				if !os.IsNotExist(err) {
					knownDirs = append(knownDirs, FileLocator{part, true})
				}
			}
		}
	}

	return knownDirs
}
