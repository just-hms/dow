package osx

import (
	"fmt"
	"log"
	"syscall"
	"unsafe"
)

var (
	modntdll                   = syscall.NewLazyDLL("ntdll.dll")
	modkernel32                = syscall.NewLazyDLL("kernel32.dll")
	procNtQueryInformationFile = modntdll.NewProc("NtQueryInformationFile")
	procCreateFileW            = modkernel32.NewProc("CreateFileW")
)

const (
	FILE_READ_ATTRIBUTES               = 0x80
	FILE_SHARE_READ                    = 0x1
	OPEN_EXISTING                      = 0x3
	FILE_FLAG_BACKUP_SEMANTICS         = 0x02000000
	INVALID_HANDLE_VALUE               = uintptr(^syscall.Handle(0))
	FileProcessIdsUsingFileInformation = 47
)

type IO_STATUS_BLOCK struct {
	Status      uintptr
	Information uintptr
}

type FILE_PROCESS_IDS_USING_FILE_INFORMATION struct {
	NumberOfProcessIdsInList int64
	ProcessIdList            [64]int64
}

func fileInfo(path string) (*FILE_PROCESS_IDS_USING_FILE_INFORMATION, error) {
	ptr, err := syscall.UTF16PtrFromString(path)
	if err != nil {
		return nil, err
	}
	// Create a file handle
	handle, _, err := procCreateFileW.Call(
		uintptr(unsafe.Pointer(ptr)), // lpFileName
		FILE_READ_ATTRIBUTES,         // dwDesiredAccess
		FILE_SHARE_READ,              // dwShareMode
		0,                            // lpSecurityAttributes
		OPEN_EXISTING,                // dwCreationDisposition
		FILE_FLAG_BACKUP_SEMANTICS,   // dwFlagsAndAttributes
		0)                            // hTemplateFile
	if handle == INVALID_HANDLE_VALUE {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer syscall.CloseHandle(syscall.Handle(handle))

	// Prepare the IO_STATUS_BLOCK and FILE_PROCESS_IDS_USING_FILE_INFORMATION structures
	var iosb IO_STATUS_BLOCK
	var info FILE_PROCESS_IDS_USING_FILE_INFORMATION

	// Call NtQueryInformationFile to get the list of PIDs
	ret, _, _ := procNtQueryInformationFile.Call(
		handle,
		uintptr(unsafe.Pointer(&iosb)),
		uintptr(unsafe.Pointer(&info)),
		unsafe.Sizeof(info),
		FileProcessIdsUsingFileInformation,
	)

	// Check the return value
	if int(ret) < 0 {
		return nil, syscall.Errno(ret)
	}

	return &info, nil
}
func IsLocked(path string) bool {
	info, err := fileInfo(path)
	if err != nil {
		return false
	}

	log.Println(info.ProcessIdList)

	return len(info.ProcessIdList) > 1
}
