//go:build linux || darwin

package osx

import "os/exec"

func IsLocked(path string) bool {
	cmd := exec.Command("lsof", path)
	out, err := cmd.Output()
	return err == nil && len(out) != 0
}
