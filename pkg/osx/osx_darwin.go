package osx

import (
	"errors"
	"io/fs"
	"syscall"
	"time"
)

// getCreationTime returns the last inode change time (ctime), which on
// macOS is bumped when a download completes and the browser renames the
// temp file into place, making it a good proxy for "just downloaded".
func getCreationTime(info fs.FileInfo) (time.Time, error) {
	statT, ok := info.Sys().(*syscall.Stat_t)
	if !ok {
		return time.Time{}, errors.New("failed to get raw syscall.Stat_t")
	}
	return time.Unix(statT.Ctimespec.Sec, statT.Ctimespec.Nsec), nil
}
