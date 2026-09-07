package osx

import (
	"path/filepath"
	"slices"
	"strings"
)

// inProgressExts are extensions browsers use for downloads that aren't done
// yet, e.g. Chrome's "Unconfirmed 12345.crdownload" while a SafeBrowsing
// keep/discard prompt is pending: the file may already be unlocked at that
// point, but it isn't the final file, so it must be skipped like a dot-file.
var inProgressExts = []string{".crdownload", ".part", ".download"}

func IsInProgressDownload(name string) bool {
	ext := strings.ToLower(filepath.Ext(name))
	return slices.Contains(inProgressExts, ext)
}
