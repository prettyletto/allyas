package printer

import (
	"fmt"
	"io"
	"strings"

	"github.com/Prettyletto/Allyas/internal/app/aliases/list"
)

type TerminalInfo struct {
	Width int
	IsTTY bool
}

func PrintCompact(w io.Writer, items []list.ListOutput, term TerminalInfo) error {
	nameWidth := longestName(items)
	nameWidth = min(nameWidth, 24)
	nameWidth = max(nameWidth, 12)

	cmdWidth := term.Width - nameWidth - 2
	cmdWidth = max(cmdWidth, 20)

	for _, item := range items {
		name := truncate(item.Name, nameWidth)
		cmd := truncate(item.Command, cmdWidth)
		if _, err := fmt.Fprintf(w, "%-*s  %s\n", nameWidth, name, cmd); err != nil {
			return err
		}
	}

	return nil
}

func PrintDetailed(w io.Writer, items []list.ListOutput) error {
	for i, item := range items {
		if i > 0 {
			if _, err := fmt.Fprintln(w); err != nil {
				return err
			}
		}

		if _, err := fmt.Fprintln(w, item.Name); err != nil {
			return err
		}
		if _, err := fmt.Fprintf(w, "  command: %s\n", item.Command); err != nil {
			return err
		}
		if item.Description != "" {
			if _, err := fmt.Fprintf(w, "  desc:    %s\n", item.Description); err != nil {
				return err
			}
		}
		if item.Group != "" {
			if _, err := fmt.Fprintf(w, "  group:   %s\n", item.Group); err != nil {
				return err
			}
		}
		if len(item.Tags) > 0 {
			if _, err := fmt.Fprintf(w, "  tags:    %s\n", strings.Join(item.Tags, ", ")); err != nil {
				return err
			}
		}
		if item.CreatedAt != "" {
			if _, err := fmt.Fprintf(w, "  created: %s\n", item.CreatedAt); err != nil {
				return err
			}
		}
		if item.UpdatedAt != "" {
			if _, err := fmt.Fprintf(w, "  updated: %s\n", item.UpdatedAt); err != nil {
				return err
			}
		}
		if item.ShowStats {
			if _, err := fmt.Fprintf(w, "  usage:   %d\n", item.UsageCount); err != nil {
				return err
			}

			if item.LastUsedAt != "" {
				if _, err := fmt.Fprintf(w, "  last:    %s\n", item.LastUsedAt); err != nil {
					return err
				}
			} else {
				if _, err := fmt.Fprintln(w, "  last:    never"); err != nil {
					return err
				}
			}

			if item.LastExitCode != nil {
				if _, err := fmt.Fprintf(w, "  exit:    %d\n", *item.LastExitCode); err != nil {
					return err
				}
			}
		}
	}

	return nil
}

func longestName(items []list.ListOutput) int {
	maxLen := 0
	for _, item := range items {
		n := len([]rune(item.Name))
		maxLen = max(n, maxLen)
	}
	return maxLen
}

func truncate(s string, maxLen int) string {
	if maxLen <= 0 {
		return ""
	}
	r := []rune(strings.TrimSpace(s))
	if len(r) <= maxLen {
		return string(r)
	}
	if maxLen == 1 {
		return "..."
	}

	return string(r[:maxLen-1]) + "..."
}
