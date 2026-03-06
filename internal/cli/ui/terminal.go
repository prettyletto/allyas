package ui

import (
	"os"

	"golang.org/x/term"
)

func GetTerminalWith() int {
	fd := int(os.Stdout.Fd())
	if !term.IsTerminal(fd){
		return 80
	}
	w,_, err := term.GetSize(fd)
	if err != nil || w <= 0 {
		return 80
	}
	return 80
}
