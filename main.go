package main

import (
	"errors"
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/just-hms/dow/iox"
	"github.com/just-hms/dow/osx"
)

const maxElapsedBeforeAsking = time.Minute

func getDownloadsPath() string {
	return filepath.Join(os.Getenv("HOME"), "Downloads")
}

func findLatestFile(files []fs.DirEntry) (os.FileInfo, error) {
	var latestFile os.FileInfo

	for _, file := range files {
		fInfo, err := file.Info()

		// Skip files that were deleted since listing
		if os.IsNotExist(err) {
			continue
		}

		// Return on unexpected error
		if err != nil {
			return nil, err
		}

		if latestFile == nil || fInfo.ModTime().After(latestFile.ModTime()) {
			latestFile = fInfo
		}
	}

	if latestFile == nil {
		return nil, errors.New("no file found")
	}

	return latestFile, nil
}

func main() {
	if len(os.Args) != 2 {
		log.Fatal("Usage: dow <destination path>")
	}

	downloadsPath := getDownloadsPath()

	files, err := os.ReadDir(downloadsPath)
	if err != nil {
		log.Fatalf("Error reading the download folder %v", err)
	}

	lastFile, err := findLatestFile(files)
	if err != nil {
		log.Fatalf("Failed to find latest file: %v", err)
	}

	if lastFile == nil {
		log.Fatalf("The download directory is empty")
	}

	if time.Since(lastFile.ModTime()) > maxElapsedBeforeAsking {
		fmt.Printf(
			"%q is older than %v. Proceed? [y]/N ",
			lastFile.Name(),
			maxElapsedBeforeAsking,
		)

		resp := iox.ReadChar(os.Stdin)
		if resp != 'y' && resp != '\n' {
			return
		}
	}

	sourcePath := filepath.Join(downloadsPath, lastFile.Name())
	destPath := os.Args[1]

	err = osx.Move(sourcePath, destPath)
	if err != nil {
		log.Fatalf("Failed to move %v", err)
	}
}
