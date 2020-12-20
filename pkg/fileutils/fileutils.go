// Package fileutils contains utilities for file operations.
// Based on https://github.com/facebook/openbmc/tools/flashy/lib/fileutils/fileutils.go
// Original author: lhl2617, Facebook, (GPLv2)
package fileutils

import (
	"io/ioutil"
	"log"
	"os"
	"path"
	"runtime"
)

// GetCurrFilePath gets the current file path.
func GetCurrFilePath() string {
	return getCurrFilePath(1)
}

func getCurrFilePath(offset int) string {
	_, filename, _, ok := runtime.Caller(1 + offset)
	if !ok {
		panic("No caller information!")
	}
	return filename
}

// GetCurrFileDir gets the current file's directory.
func GetCurrFileDir() string {
	return getCurrFileDir(1)
}

func getCurrFileDir(offset int) string {
	return path.Dir(getCurrFilePath(1 + offset))
}

// GetAllFilesInCurrDir gets all files in the folder of the caller's file.
func GetAllFilesInCurrDir() ([]os.FileInfo, error) {
	return ioutil.ReadDir(getCurrFileDir(1))
}

// PathExists returns true when the path exists (can be file/directory).
// It defaults to `false` if os.Stat returns any other non-nil error.
func PathExists(path string) bool {
	return pathExists(path, os.Stat)
}

func pathExists(path string, stat func(name string) (os.FileInfo, error)) bool {
	_, err := stat(path)
	if err == nil {
		// exists for sure
		return true
	} else if os.IsNotExist(err) {
		// does not exist for sure
		return false
	} else {
		// may or may not exist
		log.Printf("Existence check of path '%v' returned error '%v', defaulting to false",
			path, err)
		return false
	}
}

// RemovePathIfExists removes the path pointed by path if it exists.
func RemovePathIfExists(path string) error {
	return removePathIfExists(path, os.RemoveAll, PathExists)
}

func removePathIfExists(
	path string,
	remove func(name string) error,
	pathExists func(path string) bool,
) error {
	if pathExists(path) {
		return remove(path)
	}
	return nil
}
