//go:build windows

// enable ANSI escape sequences in the console
// https://learn.microsoft.com/en-us/windows/console/console-virtual-terminal-sequences#output-sequences

package prettyslog

import (
	"os"

	"golang.org/x/sys/windows"
)

func init() {
	var mode uint32
	console := windows.Handle(os.Stdout.Fd())
	if err := windows.GetConsoleMode(console, &mode); err != nil {
		return
	}
	mode |= windows.ENABLE_PROCESSED_OUTPUT | windows.ENABLE_VIRTUAL_TERMINAL_PROCESSING
	_ = windows.SetConsoleMode(console, mode)
}
