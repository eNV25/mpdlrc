package widget

import (
	"sync"

	"github.com/gdamore/tcell/v2/views"
)

type common struct {
	mu sync.Mutex

	views.ViewPort
}

func (w *common) View() views.View {
	return &w.ViewPort
}

func (w *common) Resize() {
	w.mu.Lock()
	defer w.mu.Unlock()
	w.ViewPort.ValidateView()
}
