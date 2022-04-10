package main

import (
	"io/fs"
	"os"
	"path/filepath"
)

func getGlobalVirtualEnvDirs() []string {
	homeDir, found := os.LookupEnv("HOME")
	dirs := []string{"/usr/local/bin", "/usr/bin"}

	if found {
		for _, s := range []string{"Envs", ".direnv", ".venvs", ".virtualenvs", ".local/share/virtualenvs"} {
			venvDir := filepath.Join(homeDir, s)
			_, err := os.Stat(venvDir)
			if !os.IsNotExist(err) {
				dirs = append(dirs, venvDir)
			}
		}

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
