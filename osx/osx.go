package osx

import (
	"errors"
	"io/fs"
	"os"
	"path/filepath"
)

func DownloadFolderPath() (string, error) {
	s := os.Getenv("DOWNLOAD_FOLDER_PATH")
	if s != "" {
		return s, nil
	}

	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	return filepath.Join(homeDir, "Downloads"), nil
}

func Move(sourcePath, destPath string) error {
	fInfo, err := os.Stat(destPath)
	if err == nil && fInfo.IsDir() {
		destPath = filepath.Join(destPath, filepath.Base(sourcePath))
	}

	return os.Rename(sourcePath, destPath)
}

func LatestFile(files []fs.DirEntry) (os.FileInfo, error) {
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
