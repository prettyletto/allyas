package imports

import (
	"fmt"
	"regexp"
	"strings"
)

type ParsedEntry struct {
	Name    string
	Command string
}

var (
	aliasLine = regexp.MustCompile(`^\s*alias\s+([A-Za-z_][A-Za-z0-9_]*)=(.*)$`)
	fnStart   = regexp.MustCompile(`^\s*([A-Za-z_][A-Za-z0-9_]*)\s*\(\)\s*\{\s*$`)
	fnOneLine = regexp.MustCompile(`^\s*([A-Za-z_][A-Za-z0-9_]*)\s*\(\)\s*\{\s*(.*?)\s*;\s*\}\s*$`)
)

func Parse(content string) ([]ParsedEntry, error) {
	lines := strings.Split(content, "\n")
	var entries []ParsedEntry

	for i := 0; i < len(lines); i++ {
		line := strings.TrimSpace(lines[i])
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		if m := aliasLine.FindStringSubmatch(lines[i]); m != nil {
			cmd := strings.TrimSpace(m[2])
			cmd = strings.Trim(cmd, `"'`)
			entries = append(entries, ParsedEntry{Name: m[1], Command: cmd})
			continue
		}

		if m := fnOneLine.FindStringSubmatch(lines[i]); m != nil {
			command := strings.TrimSpace(m[2])
			if command == "" {
				command = ":"
			}
			entries = append(entries, ParsedEntry{Name: m[1], Command: command})
			continue
		}

		if m := fnStart.FindStringSubmatch(lines[i]); m != nil {
			name := m[1]
			var body []string

			i++
			for ; i < len(lines); i++ {
				if strings.TrimSpace(lines[i]) == "}" {
					break
				}
				body = append(body, strings.TrimSpace(lines[i]))
			}

			if i >= len(lines) {
				return nil, fmt.Errorf("function %q is missing closing }", name)
			}

			command := strings.Join(body, "\n")
			if strings.TrimSpace(command) == "" {
				command = ":"
			}

			entries = append(entries, ParsedEntry{Name: name, Command: command})
		}

	}
	return entries, nil
}
