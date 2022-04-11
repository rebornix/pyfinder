package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
)

func getCandidatesFromEnvironmentsTxt(home string, foundHome bool, suffix string) []string {
	folders := []string{}

	if foundHome {
		env := filepath.Join(home, ".conda", "environments.txt")

		f, err := os.Open(env)
		if err != nil {
			fmt.Println(err)
			return folders
		}
		defer f.Close()
		scanner := bufio.NewScanner(f)

		for scanner.Scan() {
			// do something with a line
			line := scanner.Text()
			line = strings.TrimSpace(line)
			if !strings.HasPrefix(line, "#") {
				folders = append(folders, filepath.Join(scanner.Text(), suffix))
			}
		}

		if err := scanner.Err(); err != nil {
			fmt.Println(err)
			return folders
		}
	}

	return folders
}

func getCandidatesFromKnownPaths(home string, foundHome bool, suffix string) []string {
	prefixes := []string{}

	if runtime.GOOS == "windows" {
		programData, found := os.LookupEnv("PROGRAMDATA")
		if !found {
			programData = "C:\\ProgramData"
		}

		prefixes = append(prefixes, programData)
		if foundHome {
			localAppData, appDataFound := os.LookupEnv("localAppData")

			if !appDataFound {
				localAppData = filepath.Join(home, "AppData", "Local")
			}

			prefixes = append(prefixes, filepath.Join(localAppData, "Continuum"))
		}
	} else {
		prefixes = append(prefixes, "/usr/share", "/usr/local/share", "/opt")

		if foundHome {
			prefixes = append(prefixes, home, filepath.Join(home, "opt"))
		}
	}

	folders := []string{}
	for _, prefix := range prefixes {
		filepath.WalkDir(prefix, func(path string, info fs.DirEntry, err error) error {
			if err != nil {
				fmt.Println(err)
				return nil
			}

			if path != prefix && info.IsDir() {
				return filepath.SkipDir
			}

			if strings.Contains(strings.ToLower(info.Name()), "conda") {
				folders = append(folders, filepath.Join(path, suffix))
			}

			return nil
		})
	}

	return folders
}

type CondaInfo struct {
	Envs []string `json:"envs"`
}

func getConda(wg *sync.WaitGroup) {
	defer wg.Done()

	dirs := []string{}
	home, found := getUserHome()
	// TODO customCondaPath
	if found {
		suffix := ""

		if runtime.GOOS == "windows" {
			suffix = "Scripts\\conda.exe"
		} else {
			suffix = "bin/conda"
		}

		dirs = append(dirs, getCandidatesFromKnownPaths(home, found, suffix)...)
		dirs = append(dirs, getCandidatesFromEnvironmentsTxt(home, found, suffix)...)
	}

	for _, dir := range dirs {
		_, err := os.Stat(dir)
		if !os.IsNotExist(err) {
			// it exist
			cmd := exec.Command(dir, "info", "--json")
			out, err := cmd.Output()

			if err != nil {
				continue
			}

			var condaInfo CondaInfo
			json.Unmarshal([]byte(string(out)), &condaInfo)

			for _, env := range condaInfo.Envs {
				var pyPath string
				if runtime.GOOS == "windows" {
					pyPath = filepath.Join(env, "python.exe")
				} else {
					pyPath = filepath.Join(env, "bin/python")
				}

				_, err := os.Stat(pyPath)
				if !os.IsNotExist(err) {
					files = append(files, pyPath)
				}
			}
			break
		}
	}
}
