package termx

import (
	"context"
	"fmt"
	"os"
	"time"

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

// Spin displays a spinner animation until the context is canceled.
func Spin(ctx context.Context) error {
	oldState, err := term.MakeRaw(int(os.Stdin.Fd()))
	if err != nil {
		return err
	}
	defer term.Restore(int(os.Stdin.Fd()), oldState)

	frames := []rune{'|', '/', '-', '\\'}
	delay := 100 * time.Millisecond

	go func() {
		time.Sleep(delay * 5)

		for i := 0; ; i = (i + 1) % len(frames) {
			select {
			case <-ctx.Done():
				return
			default:
				fmt.Printf("\r%c", frames[i])
				time.Sleep(delay)
				fmt.Print("\r")
			}
		}
	}()

	return nil
}
