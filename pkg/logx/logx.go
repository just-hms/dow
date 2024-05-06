package logx

import (
	"fmt"
	"os"
)

type Logger struct{}

func (l Logger) Printf(format string, v ...any) {
	fmt.Fprintf(os.Stderr, format, v...)
}
func (l Logger) Println(a ...any) {
	fmt.Fprintln(os.Stderr, a...)
}
