package main

import (
	"io/fs"
	"path/filepath"
)

func getPyenvDirs() []FileLocator {
	// TODO test Windows `getUserHome` now works on Windows
	homeDir, found := getUserHome()
	dirs := []FileLocator{}

	if found {
		pyenv := filepath.Join(homeDir, ".pyenv", "versions")
		filepath.WalkDir(pyenv, func(path string, info fs.DirEntry, err error) error {
			if err != nil {
				return err
			}

			if path != pyenv && info.IsDir() {
				dirs = append(dirs, FileLocator{filepath.Join(path, "bin"), true})
				return filepath.SkipDir
			}

			return nil
		})
	}

	return dirs

}
