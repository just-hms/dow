package osx

import (
	"os"
)

func IsLocked(path string) bool {
	file, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		// If the error is a sharing violation, it means the file is in use
		if _, ok := err.(*os.PathError); ok {
			return true
		}
		return false
	}
	defer file.Close() // Ensure we close the file if we successfully opened it
	return false
}
