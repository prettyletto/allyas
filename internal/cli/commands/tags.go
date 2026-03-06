package commands

import "strings"

func splitTags(raw string) []string {
	parts := strings.Split(raw, ",")
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		t := strings.TrimSpace(p)
		if t == "" {
			continue
		}
		out = append(out, t)
	}
	return out
}

func appendUniqueTags(dst []string, values []string) []string {
	seen := make(map[string]bool, len(dst))
	for _, t := range dst {
		seen[strings.ToLower(strings.TrimSpace(t))] = true
	}
	for _, t := range values {
		key := strings.ToLower(strings.TrimSpace(t))
		if key == "" || seen[key] {
			continue
		}
		dst = append(dst, t)
		seen[key] = true
	}
	return dst
}
