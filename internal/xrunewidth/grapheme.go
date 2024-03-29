package xrunewidth

import "github.com/mattn/go-runewidth"

// GraphemeWidth returns the number of cells in rs.
func GraphemeWidth(rs []rune) (wd int) {
	// copied from [runewidth.StringWidth]
	for _, r := range rs {
		wd = runewidth.RuneWidth(r)
		if wd > 0 {
			break // Our best guess at this point is to use the width of the first non-zero-width rune.
		}
	}
	return
}
