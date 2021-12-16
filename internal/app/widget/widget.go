package widget

import (
	"github.com/gdamore/tcell/v2"
	"github.com/gdamore/tcell/v2/views"
)

type Widget interface {
	Draw()
	Resize()
	HandleEvent(ev tcell.Event) bool
	SetView(view views.View)
	Size() (x int, y int)
}
