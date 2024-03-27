//go:build linux

package osx

import "os/exec"

func IsLocked(filePath string) bool {
	cmd := exec.Command("lsof", filePath)
	out, err := cmd.Output()
	return err == nil && len(out) != 0
}
