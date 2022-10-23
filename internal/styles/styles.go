// Package styles implements some elements in common used by widgets.
package styles

import "github.com/gdamore/tcell/v2"

// Default returns the default [tcell.Style].
func Default() tcell.Style {
	return tcell.Style{}
}

// Border returns the style for borders.
func Border() tcell.Style {
	return tcell.Style{}.Foreground(tcell.ColorGrey)
}

const (
	// RuneBorderUpper is the rune used as the upper border of the status bar.
	RuneBorderUpper rune = tcell.RuneS9

	// RuneBorderLower is the rune used as the lower border of the status bar.
	RuneBorderLower rune = tcell.RuneS1
)
