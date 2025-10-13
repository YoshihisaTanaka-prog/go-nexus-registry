package customError

import (
	"fmt"
	"os"
)

func Exit1(args ...any) {
	fmt.Fprintln(os.Stderr, args...)
	os.Exit(1)
}