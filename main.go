package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/just-hms/dow/osx"
	"github.com/just-hms/dow/termx"
)

const maxElapsedBeforeAsking = time.Minute

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

	lastFile, err := osx.LatestFile(files)
	if err != nil {
		log.Fatalf("Failed to find latest file: %v", err)
	}

	if lastFile == nil {
		log.Fatalf("The download directory is empty")
	}

	if lastFile.IsDir() {
		log.Fatalf("The most recent file %q is a directory", lastFile.Name())
	}

	sourcePath := filepath.Join(downloadPath, lastFile.Name())

	ctx, cancel := context.WithCancel(context.Background())
	err = termx.Spin(ctx)
	if err != nil {
		log.Fatal(err)
	}

	for osx.IsLocked(sourcePath) {
		time.Sleep(200 * time.Millisecond)
	}

	cancel()

	if time.Since(lastFile.ModTime()) > maxElapsedBeforeAsking {
		fmt.Printf(
			"%q is older than %v. Proceed? (y/N) ",
			lastFile.Name(),
			maxElapsedBeforeAsking,
		)

		resp, _ := termx.Read()
		if resp != 'y' {
			fmt.Println("No")
			return
		}
		fmt.Println()
	}

	err = osx.Move(sourcePath, destPath)
	if err != nil {
		log.Fatalf("Failed to move %v", err)
	}
}
