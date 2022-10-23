package widget

import (
	"context"

	"github.com/gdamore/tcell/v2"
	"github.com/gdamore/tcell/v2/views"
)

// Widget is an iterface for widgets.
type Widget interface {
	View() views.View
	SetView(view views.View)
	Size() (x int, y int)
	Update(ctx context.Context, ev tcell.Event)
}

type cellModel interface {
	GetCell(x, y int) (mainc rune, combc []rune, style tcell.Style, cwidth int)
	GetBounds() (width int, height int)
}
