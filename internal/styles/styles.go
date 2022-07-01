package styles

import "github.com/gdamore/tcell/v2"

func BorderStyle() tcell.Style {
	return tcell.Style{}.Foreground(tcell.ColorGrey)
}

const (
	BorderD rune = '🭶'
	BorderU rune = '🭻'
)