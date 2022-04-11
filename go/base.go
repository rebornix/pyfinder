package main

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
)

var files []string

func VisitFile(path string, info fs.DirEntry, err error) error {
	if err != nil {
		fmt.Println(err)
		return nil
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
}

func getUserHome() (string, bool) {
	localOS := runtime.GOOS

	if localOS == "windows" {
		return os.LookupEnv("USERPROFILE")
	} else {
		p, found := os.LookupEnv("HOME")
		if found {
			return p, found
		} else {
			return os.LookupEnv("HOMEPATH")
		}
	}
}
