package ui

import (
	"os"

	"golang.org/x/term"
)

// IsTerminal проверяет является ли stdout терминалом
func IsTerminal() bool {
	return term.IsTerminal(int(os.Stdout.Fd()))
}
