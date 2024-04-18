//go:build windows

package osx

import (
	"syscall"

	"golang.org/x/sys/windows"
)

func IsLocked(filepath string) bool {
	utf16, _ := syscall.UTF16PtrFromString(filepath)

	// Try to open the file with a flag that fails if the file is in use.
	// GENERIC_READ is chosen here; the specific flags may need to be adjusted based on the use case.
	handle, err := windows.CreateFile(
		utf16,
		windows.GENERIC_READ,
		0, // This disallows sharing, meaning if the file is open by another process, it should fail.
		nil,
		windows.OPEN_EXISTING,
		windows.FILE_ATTRIBUTE_NORMAL,
		0,
	)
	if err != nil {
		return true
	}
	windows.CloseHandle(handle)
	return false
}
