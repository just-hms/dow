package termx

import (
	"context"
	"fmt"
	"os"
	"sync"
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

func NewSpinner(text string) *spinner {
	return &spinner{text: text}
}

type spinner struct {
	cancel context.CancelFunc
	text   string
	wg     sync.WaitGroup
}

// Spin displays a spinner animation until the context is canceled.
func (s *spinner) Spin() error {
	ctx, cancel := context.WithCancel(context.Background())
	s.cancel = cancel

	oldState, err := term.MakeRaw(int(os.Stdin.Fd()))
	if err != nil {
		return err
	}
	defer term.Restore(int(os.Stdin.Fd()), oldState)

	fmt.Print(s.text)

	frames := []rune{'|', '/', '-', '\\'}
	delay := 100 * time.Millisecond

	s.wg.Add(1)
	go func() {
		defer s.wg.Done()
		time.Sleep(delay * 5)

		for i := 0; ; i = (i + 1) % len(frames) {
			select {
			case <-ctx.Done():
				for range s.text {
					fmt.Print("\r")
				}
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

func (s *spinner) Stop() {
	s.cancel()
	s.wg.Wait()
}
