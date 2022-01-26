package internal

import (
	"sort"
	"sync"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/gdamore/tcell/v2/views"
)

type WidgetProgress struct {
	sync.RWMutex

	postFunc func(fn func())

	view views.View

	elapsedX int
	totalX   int
	offsetY  int
	runes    [3]rune
	styles   [3]tcell.Style

	toCall struct {
		sync.Mutex
		*time.Timer
	}
	id   string
	quit <-chan struct{}
}

func NewProgressWidget(postFunc func(fn func()), quit <-chan struct{}) *WidgetProgress {
	return &WidgetProgress{
		postFunc: postFunc,
		quit:     quit,
		runes:    [3]rune{'=', '>', '-'},
		styles: [3]tcell.Style{
			tcell.StyleDefault.Attributes(tcell.AttrBold),
			tcell.StyleDefault.Attributes(tcell.AttrBold),
			tcell.StyleDefault.Attributes(tcell.AttrDim),
		},
	}
}

func (w *WidgetProgress) Cancel() {
	w.toCall.Lock()
	defer w.toCall.Unlock()
	if w.toCall.Timer != nil {
		w.toCall.Stop()
	}
}

func (w *WidgetProgress) Update(
	playing bool,
	id string,
	elapsed time.Duration,
	duration time.Duration,
) {
	w.Lock()
	w.id = id
	d := duration / time.Duration(w.totalX)
	w.elapsedX = sort.Search(w.totalX, func(i int) bool { return (time.Duration(i) * d) >= elapsed })
	w.Unlock()

	w.RLock()
	defer w.RUnlock()

	select {
	case <-w.quit:
		return
	default:
		if w.id != id || w.elapsedX >= w.totalX {
			return
		}
	}

	if playing {
		go w.update(id, d)
	} else {
		go func() {
			w.postFunc(w.Draw)
		}()
	}
}

func (w *WidgetProgress) update(id string, d time.Duration) {
	w.RLock()
	defer w.RUnlock()

	select {
	case <-w.quit:
		return
	default:
		if w.id != id || w.elapsedX >= w.totalX {
			return
		}
	}

	w.postFunc(w.Draw)

	w.toCall.Lock()
	defer w.toCall.Unlock()
	w.toCall.Timer = time.AfterFunc(d, func() {
		w.Lock()
		w.elapsedX += 1
		w.Unlock()
		w.update(id, d)
	})
}

func (w *WidgetProgress) Draw() {
	w.RLock()
	defer w.RUnlock()
	w.view.Fill(' ', tcell.StyleDefault)
	for x := 0; x < w.elapsedX; x++ {
		w.view.SetContent(x, w.offsetY, w.runes[0], nil, w.styles[0])
	}
	w.view.SetContent(w.elapsedX, w.offsetY, w.runes[1], nil, w.styles[1])
	for x := w.elapsedX + 1; x < w.totalX; x++ {
		w.view.SetContent(x, w.offsetY, w.runes[2], nil, w.styles[2])
	}
}

func (w *WidgetProgress) SetView(view views.View) {
	w.Lock()
	w.view = view
	w.Unlock()
}

func (w *WidgetProgress) Resize() {
	w.Lock()
	w.totalX, _ = w.view.Size()
	w.Unlock()
}

func (w *WidgetProgress) Size() (int, int) {
	w.RLock()
	defer w.RUnlock()
	x, _ := w.view.Size()
	return x, 1
}

func (*WidgetProgress) HandleEvent(tcell.Event) bool { return false }
