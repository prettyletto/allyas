package ui

import (
	"bufio"
	"fmt"
	"io"
	"strings"
)

func AskYesNo(in io.Reader, out io.Writer, question string, defaultYes bool) (bool, error) {
	suffix := " [y/N]: "
	if defaultYes {
		suffix = " [Y/n]: "
	}
	fmt.Fprint(out, question, suffix)

	r := bufio.NewReader(in)
	line, err := r.ReadString('\n')

	if err != nil && err != io.EOF {
		return false, err
	}
	s := strings.ToLower(strings.TrimSpace(line))

	if s == "" {
		return defaultYes, nil
	}

	return s == "y" || s == "yes", nil
}

func AskInput(in io.Reader, out io.Writer, question, suffix, defaultValue string) (string, error) {
	fmt.Fprint(out, question, suffix)

	r := bufio.NewReader(in)
	line, err := r.ReadString('\n')
	if err != nil && err != io.EOF {
		return "", err
	}
	s := strings.ToLower(strings.TrimSpace(line))

	if s == "" {
		return "", nil
	}

	return s, nil
}
