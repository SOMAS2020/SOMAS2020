// Package fileutils contains utilities for file operations.
package fileutils

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
)

// GetAllFilesInCurrDir gets all files in the folder of the caller's file.
func GetAllFilesInCurrDir() ([]os.FileInfo, error) {
	_, filename, _, ok := runtime.Caller(1)
	if !ok {
		panic("No caller information")
	}
	return ioutil.ReadDir(filepath.Dir(filename))
}
