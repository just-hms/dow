package osx

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestIsLocked(t *testing.T) {
	req := require.New(t)

	path := "testdata/test.bin"

	locked := IsLocked(path)
	req.False(locked)

	file, err := os.Open(path)
	req.NoError(err)
	defer file.Close()

	locked = IsLocked(path)
	req.True(locked)
}
