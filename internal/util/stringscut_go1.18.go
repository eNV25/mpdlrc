//go:build go1.18

package util

import "strings"

// StringsCut slices s around the first instance of sep,
// returning the text before and after sep.
// The found result reports whether sep appears in s.
// If sep does not appear in s, cut returns s, "", false.
func StringsCut(s, sep string) (before, after string, found bool) {
	return strings.Cut(s, sep)
}
