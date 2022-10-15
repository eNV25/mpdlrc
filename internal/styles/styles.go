package styles

import "github.com/gdamore/tcell/v2"

func Default() tcell.Style {
	return tcell.Style{}
}

func Border() tcell.Style {
	return tcell.Style{}.Foreground(tcell.ColorGrey)
}

const (
	RuneBorderUpper rune = tcell.RuneS9
	RuneBorderLower rune = tcell.RuneS1
)
