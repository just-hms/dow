package termx

import (
	"fmt"
	"os"

	"golang.org/x/term"
)

func Read() (rune, error) {
	// switch stdin into 'raw' mode
	oldState, err := term.MakeRaw(int(os.Stdin.Fd()))
	if err != nil {
		return 0, err
	}
	defer term.Restore(int(os.Stdin.Fd()), oldState)

	b := make([]byte, 1)
	_, err = os.Stdin.Read(b)
	if err != nil {
		return 0, err
	}
	return rune(b[0]), nil
}

type spinner struct {
	count int
}

func NewSpin() *spinner {
	return &spinner{}
}

var frames = []string{"⣷", "⣯", "⣟", "⡿", "⢿", "⣻", "⣽", "⣾"}

// Spin displays a spinner animation until the context is canceled.
func (s *spinner) Spin(text string) {
	fmt.Printf("\r%s %s  ", text, frames[s.count])
	s.count = (s.count + 1) % len(frames)
}

func (s *spinner) Flush() {
	fmt.Print("\r")
}
