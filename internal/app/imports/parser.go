package imports

import (
	"fmt"
	"regexp"
	"strings"
)

type ParsedEntry struct {
	Name    string
	Command string
	Line    int
	Renamed bool
}

type ParseWarning struct {
	Line    int
	Content string
	Reason  string
}

type ParseResult struct {
	Entries  []ParsedEntry
	Warnings []ParseWarning
}

var (
	aliasName = regexp.MustCompile(`^[A-Za-z_][A-Za-z0-9_]*$`)
	aliasLine = regexp.MustCompile(`^\s*alias\s+([A-Za-z_][A-Za-z0-9_]*)=(.*)$`)
	fnHeader  = regexp.MustCompile(`^\s*(?:function\s+)?([A-Za-z_][A-Za-z0-9_]*)(?:\s*\(\))?\s*\{\s*(.*)$`)
)

func Parse(content string) ([]ParsedEntry, error) {
	result, err := ParseDetailed(content)
	if err != nil {
		return nil, err
	}
	return result.Entries, nil
}

func ParseDetailed(content string) (ParseResult, error) {
	lines := strings.Split(content, "\n")
	var result ParseResult

	for i := 0; i < len(lines); i++ {
		line := strings.TrimSpace(lines[i])
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		fullLine, end := collectContinuation(lines, i)
		if m := aliasLine.FindStringSubmatch(fullLine); m != nil {
			cmd := normalizeAliasValue(m[2])
			if cmd == "" {
				result.Warnings = append(result.Warnings, ParseWarning{
					Line:    i + 1,
					Content: strings.TrimSpace(lines[i]),
					Reason:  "alias command is empty",
				})
				i = end
				continue
			}
			result.Entries = append(result.Entries, ParsedEntry{Name: m[1], Command: cmd, Line: i + 1})
			i = end
			continue
		}

		if strings.HasPrefix(line, "alias ") {
			result.Warnings = append(result.Warnings, ParseWarning{
				Line:    i + 1,
				Content: line,
				Reason:  "unsupported alias syntax",
			})
			continue
		}

		if m := fnHeader.FindStringSubmatch(lines[i]); m != nil {
			name := m[1]
			body, next, err := readFunctionBody(lines, i, m[2])
			if err != nil {
				result.Warnings = append(result.Warnings, ParseWarning{
					Line:    i + 1,
					Content: line,
					Reason:  err.Error(),
				})
				break
			}

			command := strings.TrimSpace(body)
			if command == "" {
				command = ":"
			}

			result.Entries = append(result.Entries, ParsedEntry{Name: name, Command: command, Line: i + 1})
			i = next
			continue
		}

		if looksLikeFunctionStart(line) {
			result.Warnings = append(result.Warnings, ParseWarning{
				Line:    i + 1,
				Content: line,
				Reason:  "unsupported function syntax",
			})
			continue
		}

		result.Warnings = append(result.Warnings, ParseWarning{
			Line:    i + 1,
			Content: line,
			Reason:  "unsupported line",
		})
	}
	return result, nil
}

func collectContinuation(lines []string, start int) (string, int) {
	out := strings.TrimRight(lines[start], " \t")
	end := start
	for strings.HasSuffix(out, `\`) && end+1 < len(lines) {
		out = strings.TrimSuffix(out, `\`) + strings.TrimLeft(lines[end+1], " \t")
		end++
	}
	return out, end
}

func normalizeAliasValue(value string) string {
	value = stripInlineComment(strings.TrimSpace(value))
	value = strings.TrimSpace(value)
	if len(value) < 2 {
		return value
	}

	quote := value[0]
	if quote != '\'' && quote != '"' {
		return value
	}
	if value[len(value)-1] != quote {
		return value
	}

	inner := value[1 : len(value)-1]
	if quote == '\'' {
		return strings.ReplaceAll(inner, `'\''`, `'`)
	}

	replacer := strings.NewReplacer(`\"`, `"`, `\\`, `\`, `\$`, `$`, "\\`", "`")
	return replacer.Replace(inner)
}

func stripInlineComment(s string) string {
	var quote rune
	escaped := false
	for i, r := range s {
		if escaped {
			escaped = false
			continue
		}
		if r == '\\' {
			escaped = true
			continue
		}
		if quote != 0 {
			if r == quote {
				quote = 0
			}
			continue
		}
		if r == '\'' || r == '"' {
			quote = r
			continue
		}
		if r == '#' && (i == 0 || s[i-1] == ' ' || s[i-1] == '\t') {
			return strings.TrimSpace(s[:i])
		}
	}
	return s
}

func readFunctionBody(lines []string, start int, firstRest string) (string, int, error) {
	body := []string{}
	rest := strings.TrimSpace(firstRest)
	if rest != "" {
		if before, ok := splitFunctionClose(rest); ok {
			return trimFunctionStatement(before), start, nil
		}
		body = append(body, rest)
	}

	for i := start + 1; i < len(lines); i++ {
		line := strings.TrimSpace(lines[i])
		if before, ok := splitFunctionClose(line); ok {
			before = trimFunctionStatement(before)
			if before != "" {
				body = append(body, before)
			}
			return strings.Join(body, "\n"), i, nil
		}
		body = append(body, line)
	}

	return "", len(lines) - 1, fmt.Errorf("function starting on line %d is missing closing }", start+1)
}

func trimFunctionStatement(s string) string {
	return strings.TrimSpace(strings.TrimSuffix(strings.TrimSpace(s), ";"))
}

func splitFunctionClose(line string) (string, bool) {
	var quote rune
	escaped := false
	for i, r := range line {
		if escaped {
			escaped = false
			continue
		}
		if r == '\\' {
			escaped = true
			continue
		}
		if quote != 0 {
			if r == quote {
				quote = 0
			}
			continue
		}
		if r == '\'' || r == '"' {
			quote = r
			continue
		}
		if r == '}' {
			return line[:i], true
		}
	}
	return "", false
}

func looksLikeFunctionStart(line string) bool {
	return strings.HasPrefix(line, "function ") || strings.Contains(line, "()")
}
