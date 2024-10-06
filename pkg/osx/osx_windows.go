package osx

import (
	"os"
	"syscall"
)

const ERROR_SHARING_VIOLATION = 32

func IsLocked(path string) bool {
	file, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		// If the error is a sharing violation, it means the file is in use
		if pathErr, ok := err.(*os.PathError); ok {
			if errno, ok := pathErr.Err.(syscall.Errno); ok && errno == ERROR_SHARING_VIOLATION {
				return true
			}
		}
		return false
	}
	defer file.Close() // Ensure we close the file if we successfully opened it
	return false
}
