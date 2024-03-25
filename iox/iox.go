package iox

import (
	"io"
)

func ReadChar(s io.Reader) rune {
	b := make([]byte, 1)
	s.Read(b)

	if rune(string(b)[0]) == 13 {
		return '\n'
	}

	return rune(string(b)[0])
}
