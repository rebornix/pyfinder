package main

import (
	"io/fs"
	"path/filepath"
)

func getPyenvDirs() []string {
	// TODO support Windows
	homeDir, found := getUserHome()
	dirs := []string{}

	if found {
		pyenv := filepath.Join(homeDir, ".pyenv", "versions")
		filepath.WalkDir(pyenv, func(path string, info fs.DirEntry, err error) error {
			if err != nil {
				return err
			}

			if path != pyenv && info.IsDir() {
				dirs = append(dirs, filepath.Join(path, "bin"))
				return filepath.SkipDir
			}

			return nil
		})
	}

	return dirs

}
