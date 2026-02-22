package handler

import "strings"

// SplitCommaSeparated converts a single comma-separated value into a slice.
// Empty segments are ignored; if all segments are empty, it returns the original value as a single element.
func SplitCommaSeparated(value string) []string {
	parts := strings.Split(value, ",")
	res := make([]string, 0, len(parts))
	for _, p := range parts {
		v := strings.TrimSpace(p)
		if v == "" {
			continue
		}
		res = append(res, v)
	}

	if len(res) == 0 {
		return []string{value}
	}

	return res
}
