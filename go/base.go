package main

import (
	"os"
	"runtime"
)

type FileLocator struct {
	path    string
	skipDir bool
}

var files []string

// export function matchPythonBinFilename(filename: string): boolean {
//     /**
//      * This Reg-ex matches following file names:
//      * python
//      * python3
//      * python38
//      * python3.8
//      */
//     const posixPythonBinPattern = /^python(\d+(\.\d+)?)?$/;

//     return posixPythonBinPattern.test(path.basename(filename));
// }

// export function matchPythonBinFilename(filename: string): boolean {
//     /**
//      * This Reg-ex matches following file names:
//      * python.exe
//      * python3.exe
//      * python38.exe
//      * python3.8.exe
//      */
//     const windowsPythonExes = /^python(\d+(.\d+)?)?\.exe$/;

//     return windowsPythonExes.test(path.basename(filename));
// }

// func VisitFile(path string, info fs.DirEntry, err error) error {
// 	if err != nil {
// 		fmt.Println(err)
// 		return nil
// 	}

// 	if !info.IsDir() {
// 		if strings.HasPrefix(info.Name(), "python") {
// 			matched := r.MatchString(info.Name())
// 			if matched {
// 				if info.Type()&fs.ModeSymlink == fs.ModeSymlink {
// 					realPath, err := filepath.EvalSymlinks(path)
// 					if err == nil {
// 						files = append(files, realPath)
// 					}
// 				} else {
// 					files = append(files, path)
// 				}
// 			}
// 		}
// 	}

// 	return nil
// }

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
