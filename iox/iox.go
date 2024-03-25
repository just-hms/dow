package iox

import (
	"io"
)

func ReadChar(s io.Reader) rune {
	b := make([]byte, 1)
	s.Read(b)
	return rune(string(b)[0])
}
