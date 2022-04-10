package main

import (
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

func getHomeBrewEnvDirs() []string {
	// home brew package path
	cellar := "/usr/local/Cellar"
	dirs := []string{}

	_, err := os.Stat(cellar)
	if !os.IsNotExist(err) {
		// pyenvsDir := []string{}
		err = filepath.WalkDir(cellar, func(path string, info fs.DirEntry, err error) error {
			if err != nil {
				// fmt.Println(err)
				return err
			}

			if info.IsDir() && path != cellar {
				if strings.HasPrefix(info.Name(), "python") {
					dirs = append(dirs, path)
				}

				return filepath.SkipDir
			}

			return nil
		})
	}

	return dirs
}
