package internal

import (
	"sort"
	"sync"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/gdamore/tcell/v2/views"

	"github.com/env25/mpdlrc/internal/syncu"
)

type WidgetProgress struct {
	mu sync.Mutex

	toCall struct {
		syncu.Once
		*time.Timer
	}

	view views.View

	widgetProgressData

	quit     <-chan struct{}
	postFunc func(fn func())
}

type widgetProgressData struct {
	duration time.Duration
	elapsedX int
	totalX   int
	offsetY  int
}

func NewWidgetProgress(postFunc func(fn func()), quit <-chan struct{}) *WidgetProgress {
	return &WidgetProgress{
		postFunc: postFunc,
		quit:     quit,
	}
}

func (w *WidgetProgress) Cancel() {
	w.mu.Lock()
	defer w.mu.Unlock()
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
	w.mu.Lock()
	defer w.mu.Unlock()

	w.duration = duration / time.Duration(w.totalX)
	w.elapsedX = sort.Search(w.totalX, func(i int) bool { return (time.Duration(i) * w.duration) >= elapsed })

	select {
	case <-w.quit:
		return
	default:
		if w.elapsedX >= w.totalX {
			return
		}
	}

	if playing {
		go func() {
			w.mu.Lock()
			defer w.mu.Unlock()
			w.update()
		}()
	} else {
		go w.postFunc(w.Draw)
	}
}

func (w *WidgetProgress) update() {
	select {
	case <-w.quit:
		return
	default:
		if w.elapsedX >= w.totalX {
			return
		}
	}

	go w.postFunc(w.Draw)

	if !w.toCall.Once.Do(func() {
		w.toCall.Timer = time.AfterFunc(w.duration, func() {
			w.mu.Lock()
			defer w.mu.Unlock()
			w.elapsedX += 1
			w.update()
		})
	}) {
		w.toCall.Reset(w.duration)
	}
}

func (w *WidgetProgress) Draw() {
	w.mu.Lock()
	defer w.mu.Unlock()

	const (
		rune0 rune = '='
		rune1 rune = '>'
		rune2 rune = '-'
	)

	var (
		styleDefault tcell.Style
		style0       = styleDefault.Bold(true)
		style1       = style0
		style2       = styleDefault.Dim(true)
	)

	w.view.Fill(' ', styleDefault)
	for x := 0; x < w.elapsedX; x++ {
		w.view.SetContent(x, w.offsetY, rune0, nil, style0)
	}
	w.view.SetContent(w.elapsedX, w.offsetY, rune1, nil, style1)
	for x := w.elapsedX + 1; x < w.totalX; x++ {
		w.view.SetContent(x, w.offsetY, rune2, nil, style2)
	}
}

func (w *WidgetProgress) SetView(view views.View) {
	w.mu.Lock()
	defer w.mu.Unlock()
	w.view = view
}

func (w *WidgetProgress) Resize() {
	w.mu.Lock()
	defer w.mu.Unlock()
	w.totalX, _ = w.view.Size()
}

func (w *WidgetProgress) Size() (int, int) {
	w.mu.Lock()
	defer w.mu.Unlock()
	x, _ := w.view.Size()
	return x, 1
}

func (*WidgetProgress) HandleEvent(tcell.Event) bool { return false }
