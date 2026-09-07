package osx

import (
	"errors"
	"io/fs"
	"syscall"
	"time"
)

// errSharingViolation is ERROR_SHARING_VIOLATION, returned by CreateFile
// when the file is already open elsewhere without a compatible share mode.
const errSharingViolation = syscall.Errno(32)

// IsLocked reports whether path is currently open by another process (e.g. a
// browser still writing to it), by trying to open it with no sharing allowed.
func IsLocked(path string) bool {
	pathPtr, err := syscall.UTF16PtrFromString(path)
	if err != nil {
		return false
	}

	handle, err := syscall.CreateFile(
		pathPtr,
		syscall.GENERIC_READ,
		0, // no sharing: fails if another process has the file open
		nil,
		syscall.OPEN_EXISTING,
		syscall.FILE_ATTRIBUTE_NORMAL,
		0,
	)
	if err != nil {
		return err == errSharingViolation
	}
	syscall.CloseHandle(handle)
	return false
}

// getCreationTime returns the last modification time, matching the Unix
// implementations' use of a timestamp that's bumped when a download
// completes and the browser renames the temp file into place.
func getCreationTime(info fs.FileInfo) (time.Time, error) {
	attrs, ok := info.Sys().(*syscall.Win32FileAttributeData)
	if !ok {
		return time.Time{}, errors.New("failed to get raw syscall.Win32FileAttributeData")
	}
	return time.Unix(0, attrs.LastWriteTime.Nanoseconds()), nil
}
