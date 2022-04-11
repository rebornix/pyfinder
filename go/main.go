package main

// https://github.com/microsoft/vscode-python/wiki/Interpreter-and-Environment-Discovery

import (
	"fmt"
	"io/fs"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
	"sync"
	"time"
)

func search(rootPath string, skipFolder bool, wg *sync.WaitGroup) {
	defer wg.Done()

	filepath.WalkDir(rootPath, func(path string, info fs.DirEntry, err error) error {
		if err != nil {
			fmt.Println(err)
			return nil
		}

		if skipFolder && info.IsDir() && path != rootPath {
			return filepath.SkipDir
		}

		if !info.IsDir() {
			if strings.HasPrefix(info.Name(), "python") {
				matched, _ := regexp.MatchString(`^python[0-9\.]*$`, info.Name())
				if matched {
					if info.Type()&fs.ModeSymlink == fs.ModeSymlink {
						realPath, err := filepath.EvalSymlinks(path)
						if err == nil {
							files = append(files, realPath)
						}
					} else {
						files = append(files, path)
					}
				}
			}
		}

		return nil
	})
}

func removeDuplicateStr(strSlice []string) []string {
	allKeys := make(map[string]bool)
	list := []string{}
	for _, item := range strSlice {
		if _, value := allKeys[item]; !value {
			allKeys[item] = true
			list = append(list, item)
		}
	}
	return list
}

func getpython() []string {
	dirs := []string{}
	os := runtime.GOOS

	dirs = append(dirs, getGlobalVirtualEnvDirs()...)
	dirs = append(dirs, getPyenvDirs()...)

	if os != "windows" {
		dirs = append(dirs, getPosixBinPaths()...)
	}
	// dirs = append(dirs, getHomeBrewEnvDirs()...)

	var wg sync.WaitGroup
	wg.Add(len(dirs))

	for _, dir := range dirs {
		go search(dir, true, &wg)
	}
	wg.Wait()
	dedups := removeDuplicateStr(dirs)

	//TODO posix filter pyenv shims

	isMacPython2Deprecated := false
	if os == "darwin" {
		isMacPython2Deprecated = true
	}

	if isMacPython2Deprecated {
		macfiltered := []string{}
		for _, value := range dedups {
			if strings.HasPrefix(value, "python") || strings.HasPrefix(value, "/usr/bin/python") || strings.HasPrefix(value, "/usr/bin/python2") {
			} else {
				macfiltered = append(macfiltered, value)
			}
		}
		return macfiltered
	}

	return dedups
}

func main() {
	start := time.Now()
	getpython()
	duration := time.Since(start)
	fmt.Println(duration)

	for _, file := range removeDuplicateStr(files) {
		fmt.Println(file)
	}
}
