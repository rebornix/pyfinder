package main

import (
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"sync"
)

func isStorePythonInstalled(storeRoot string) bool {
	idlePath := filepath.Join(storeRoot, "idle.exe")
	pipPath := filepath.Join(storeRoot, "pip.exe")

	_, err := os.Stat(idlePath)
	if !os.IsNotExist(err) {
		return true
	}

	_, err = os.Stat(pipPath)
	return !os.IsNotExist(err)
}

func getWindowsStoreDirs(wg *sync.WaitGroup) {
	defer wg.Done()

	localAppData, appDataFound := os.LookupEnv("localAppData")
	if appDataFound {
		storeRoot := filepath.Join(localAppData, "Microsoft", "WindowsApps")
		installed := isStorePythonInstalled(storeRoot)

		exes := []string{}

		if installed {
			r, _ := regexp.Compile(`python3\.([0-9]|[0-9][0-9])\.exe`)
			filepath.WalkDir(storeRoot, func(path string, info fs.DirEntry, err error) error {
				if err != nil {
					return err
				}

				if path != storeRoot && info.IsDir() {
					return filepath.SkipDir
				}

				matches := r.FindAllString(info.Name(), -1)
				if len(matches) > 0 {
					exes = append(exes, path)
				}

				return nil
			})
		}

		files = append(files, exes...)
	}
}
