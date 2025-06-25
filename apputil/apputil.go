package apputil

import (
	"os"
	"path/filepath"

	"manifold.mleku.dev/chk"
)

// EnsureDir checks a file could be written to a path, creates the directories as needed
func EnsureDir(fileName string) {
	dirName := filepath.Dir(fileName)
	if _, err := os.Stat(dirName); chk.E(err) {
		merr := os.MkdirAll(dirName, os.ModePerm)
		if chk.E(merr) {
			panic(merr)
		}
	}
}

// FileExists reports whether the named file or directory exists.
func FileExists(filePath string) bool {
	_, e := os.Stat(filePath)
	return e == nil
}
