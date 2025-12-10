package utils

import "strings"

func NormalizeName(s string) string {

	return strings.ToLower(strings.TrimSpace(s))
}

func NormalizeSlice(slice []string) []string {
	for i, s := range slice {
		slice[i] = NormalizeName(s)
	}
	return slice
}
