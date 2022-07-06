package widget

import (
	"context"

	"github.com/gdamore/tcell/v2"
	"github.com/gdamore/tcell/v2/views"
)

type Widget interface {
	View() views.View
	SetView(view views.View)
	Size() (x int, y int)
	Resize()
	Update(ctx context.Context)
}

type cellModel interface {
	GetCell(x, y int) (mainc rune, combc []rune, style tcell.Style, cwidth int)
	GetBounds() (width int, height int)
}
