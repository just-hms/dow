//go:build linux || darwin

package osx_test

import (
	"os"
	"testing"

	"github.com/just-hms/dow/pkg/osx"
	"github.com/stretchr/testify/require"
)

func TestIsLocked(t *testing.T) {
	req := require.New(t)

	path := "testdata/test.bin"

	locked := osx.IsLocked(path)
	req.False(locked)

	file, err := os.Open(path)
	req.NoError(err)
	defer file.Close()

	locked = osx.IsLocked(path)
	req.True(locked)
}
