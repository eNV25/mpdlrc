package internal

import (
	"sync"

	"github.com/gdamore/tcell/v2"
	"github.com/gdamore/tcell/v2/views"
)

type widgetCommon struct {
	mu sync.Mutex

	view   views.View
	events chan<- tcell.Event
}

func (w *widgetCommon) SetView(view views.View) {
	w.mu.Lock()
	defer w.mu.Unlock()
	w.view = view
}

func (w *widgetCommon) View() views.View {
	w.mu.Lock()
	defer w.mu.Unlock()
	return w.view
}
