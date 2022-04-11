package main

import (
	"os"
	"path/filepath"
)

func getGlobalVirtualEnvDirs() []string {
	homeDir, found := getUserHome()
	dirs := []string{}

	if found {
		for _, s := range []string{"Envs", ".direnv", ".venvs", ".virtualenvs", ".local/share/virtualenvs"} {
			venvDir := filepath.Join(homeDir, s)
			_, err := os.Stat(venvDir)
			if !os.IsNotExist(err) {
				dirs = append(dirs, venvDir)
			}
		}
	}

	return dirs
}
