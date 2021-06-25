// textwidth provides functions for getting the fixed-width width of unicode
// byte slices, runes and strings.
//
// https://en.wikipedia.org/wiki/Halfwidth_and_fullwidth_forms#In_Unicode
//
package textwidth

import (
	"unicode"
	"unicode/utf8"

	"golang.org/x/text/width"
)

// WidthOfRune returns fixed-width width of byte slice.
func Width(b []byte) (n int) {
	for len(b) > 0 {
		r, size := utf8.DecodeRune(b)

		if r == utf8.RuneError {
			return -1
		}

		n += WidthOfRune(r)

		b = b[size:]
	}
	return n
}

// WidthOfRune returns fixed-width width of rune.
func WidthOfRune(r rune) int {
	switch {
	case unicode.Is(unicode.Mn, r), !unicode.IsGraphic(r):
		return 0
	default:
		switch width.LookupRune(r).Kind() {
		case width.EastAsianWide, width.EastAsianFullwidth:
			return 2
		case width.EastAsianNarrow, width.EastAsianHalfwidth, width.EastAsianAmbiguous, width.Neutral:
			return 1
		default:
			return 0
		}
	}
}

// WidthOfString returns fixed-width width of string.
func WidthOfString(s string) (n int) {
	for _, r := range s {
		n += WidthOfRune(r)
	}
	return n
}
