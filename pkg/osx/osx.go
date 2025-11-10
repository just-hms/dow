package osx

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"syscall"

	"github.com/just-hms/dow/pkg/bytex"
)

func DownloadFolderPath() (string, error) {
	s := os.Getenv("DOW_DOWNLOAD_PATH")
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

func getCreationTime(info fs.FileInfo) (t syscall.Timespec, err error) {
	statT, ok := info.Sys().(*syscall.Stat_t)
	if !ok {
		return t, errors.New("failed to get raw syscall.Stat_t")
	}
	return statT.Ctimespec, nil // Birth time is stored in Ctimespec on Darwin
}

// LatestFile returns the latest created file (ignoring dot-files)
func LatestFile(files []fs.DirEntry) (os.FileInfo, error) {
	var latestFile os.FileInfo
	var latestTime syscall.Timespec

	for _, file := range files {
		if file.IsDir() || file.Name()[0] == '.' {
			continue
		}

		fInfo, err := file.Info()
		if os.IsNotExist(err) {
			continue
		}
		if err != nil {
			return nil, err
		}

		ctime, err := getCreationTime(fInfo)
		if err != nil {
			continue
		}

		if latestFile == nil || ctime.Sec > latestTime.Sec || (ctime.Sec == latestTime.Sec && ctime.Nsec > latestTime.Nsec) {
			latestFile = fInfo
			latestTime = ctime
		}
	}

	if latestFile == nil {
		return nil, errors.New("no valid files found")
	}

	return latestFile, nil
}

func Size(f fs.FileInfo) string {
	s := float64(f.Size())
	if s > bytex.GB {
		return fmt.Sprintf("%.2f GB", s/bytex.GB)
	}
	if s > bytex.MB {
		return fmt.Sprintf("%.2f MB", s/bytex.MB)
	}
	if s > bytex.KB {
		return fmt.Sprintf("%.2f KB", s/bytex.KB)
	}
	return fmt.Sprintf("%.2f B", s)
}
