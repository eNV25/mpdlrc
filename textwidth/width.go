package textwidth

import (
	"unicode"

	"golang.org/x/text/width"
)

// RuneWidth returns fixed-width width of rune.
// https://en.wikipedia.org/wiki/Halfwidth_and_fullwidth_forms#In_Unicode
func RuneWidth(r rune) int {
	//      non-printing,       combing character,       null character
	if !unicode.IsPrint(r) || unicode.Is(unicode.Mn, r) || r == '\x00' {
		return 0
	}
	switch width.LookupRune(r).Kind() {
	case width.EastAsianWide, width.EastAsianFullwidth:
		return 2
	case width.EastAsianNarrow, width.EastAsianHalfwidth, width.EastAsianAmbiguous, width.Neutral:
		return 1
	default:
		return 0
	}
}

// StringWidth returns fixed-width width of string.
// https://en.wikipedia.org/wiki/Halfwidth_and_fullwidth_forms#In_Unicode
func StringWidth(s string) (n int) {
	for _, r := range s {
		n += RuneWidth(r)
	}
	return n
}
