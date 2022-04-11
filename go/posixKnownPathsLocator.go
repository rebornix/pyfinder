package main

import "os"

func getSearchPathEntries() []string {
	return []string{}
}

func getCommonPosixBinPaths() []FileLocator {
	knownPaths := []string{
		"/bin",
		"/etc",
		"/lib",
		"/lib/x86_64-linux-gnu",
		"/lib64",
		"/sbin",
		"/snap/bin",
		"/usr/bin",
		"/usr/games",
		"/usr/include",
		"/usr/lib",
		"/usr/lib/x86_64-linux-gnu",
		"/usr/lib64",
		"/usr/libexec",
		"/usr/local",
		"/usr/local/bin",
		"/usr/local/etc",
		"/usr/local/games",
		"/usr/local/lib",
		"/usr/local/sbin",
		"/usr/sbin",
		"/usr/share",
		"~/.local/bin",
	}

	knownPaths = append(knownPaths, getSearchPathEntries()...)

	dirs := []FileLocator{}

	for _, p := range knownPaths {
		_, err := os.Stat(p)
		if !os.IsNotExist(err) {
			dirs = append(dirs, FileLocator{p, true})
		}
	}

	return dirs
}

func getPosixBinPaths() []FileLocator {
	return getCommonPosixBinPaths()
}
