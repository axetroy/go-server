package pio

import (
	"io"
	"runtime"

	"github.com/kataras/pio/terminal"
)

func isTerminal(output io.Writer) bool {
	isTerminal := !IsNop(output) || terminal.IsTerminal(output)
	if isTerminal && runtime.GOOS == "windows" {
		// if on windows then return true only when it does support 256-bit colors,
		// this is why we initialy do that terminal check for the "output" writer.
		return terminal.SupportColors
	}

	return isTerminal
}
