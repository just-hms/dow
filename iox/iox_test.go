package iox_test

import (
	"bytes"
	"testing"

	"github.com/just-hms/dow/iox"
	"github.com/stretchr/testify/require"
)

func TestReadChar(t *testing.T) {
	t.Parallel()

	req := require.New(t)

	testcases := []struct {
		name  string
		input []byte
		exp   rune
	}{
		{
			input: []byte("y"),
			exp:   'y',
		},
		{
			input: []byte("a"),
			exp:   'a',
		},
		{
			input: []byte{13},
			exp:   '\n',
		},
	}

	for _, tt := range testcases {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			buf := bytes.NewBuffer([]byte(tt.input))
			res := iox.ReadChar(buf)
			req.Equal(res, tt.exp)
		})
	}
}
