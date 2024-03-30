package main

import (
	"errors"
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/just-hms/dow/osx"
	"github.com/just-hms/dow/termx"
)

const maxElapsedBeforeAsking = time.Minute

func getLastFile(path string) (fs.FileInfo, error) {
	files, err := os.ReadDir(path)
	if err != nil {
		return nil, err
	}
	lastFile, err := osx.LatestFile(files)
	if err != nil {
		return nil, err
	}
	if lastFile.IsDir() {
		return nil, errors.New("cannot move a directory")
	}
	return lastFile, nil
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

	var (
		lastFile   fs.FileInfo
		sourcePath string
	)

	// TODO: this is basically a wait for folder
	// https://superuser.com/questions/860064/how-can-i-find-all-files-open-within-a-given-directory#:~:text=lsof%20has%20switches,open%20files%20recursively)

	for {
		lastFile, err = getLastFile(downloadPath)
		if err != nil {
			log.Fatalf("Error getting the last file %v", err)
		}

		sourcePath = filepath.Join(downloadPath, lastFile.Name())

		s := termx.NewSpinner("Downloading")
		s.Spin()
		waited := false
		for osx.IsLocked(sourcePath) {
			waited = true
			time.Sleep(200 * time.Millisecond)
		}
		s.Stop()
		if !waited {
			break
		}
	}

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
		log.Fatalf("Failed to %v", err)
	}
}
