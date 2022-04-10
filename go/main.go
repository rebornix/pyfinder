package main

// https://github.com/microsoft/vscode-python/wiki/Interpreter-and-Environment-Discovery

import (
	"fmt"
	"path/filepath"
	"sync"
	"time"
)

func search(path string, wg *sync.WaitGroup) {
	defer wg.Done()

	filepath.WalkDir(path, VisitFile)
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

	dirs = append(dirs, getGlobalVirtualEnvDirs()...)
	// dirs = append(dirs, getHomeBrewEnvDirs()...)

	var wg sync.WaitGroup
	wg.Add(len(dirs))

	for _, dir := range dirs {
		go search(dir, &wg)
	}
	wg.Wait()
	return removeDuplicateStr(dirs)
}

func main() {
	start := time.Now()
	dirs := []string{}

	dirs = append(dirs, getGlobalVirtualEnvDirs()...)
	// dirs = append(dirs, getHomeBrewEnvDirs()...)
	var wg sync.WaitGroup
	wg.Add(len(dirs))

	for _, dir := range dirs {
		go search(dir, &wg)
	}
	wg.Wait()
	duration := time.Since(start)
	fmt.Println(duration)

	for _, file := range removeDuplicateStr(files) {
		fmt.Println(file)
	}
}
