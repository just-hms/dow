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
	if len(os.Args) != 1 && len(os.Args) != 2 {
		log.Fatal("Usage: dow <destination path>")
	}

	destPath := "."
	if len(os.Args) == 2 {
		destPath = os.Args[1]
	}

	downloadPath, err := osx.DownloadFolderPath()
	if err != nil {
		log.Fatalf("Error getting the dowload path %v", err)
	}

	files, err := os.ReadDir(downloadPath)
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
			"%q is older than %v. Proceed? (y/N) ",
			lastFile.Name(),
			maxElapsedBeforeAsking,
		)

		resp, _ := iox.Read()
		if resp != 'y' {
			fmt.Println("No")
			return
		}
		fmt.Println()
	}

	sourcePath := filepath.Join(downloadPath, lastFile.Name())

	err = osx.Move(sourcePath, destPath)
	if err != nil {
		log.Fatalf("Failed to move %v", err)
	}
}
