package internal

import (
	"sort"
	"sync"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/gdamore/tcell/v2/views"

	"github.com/env25/mpdlrc/internal/util"
)

type WidgetProgress struct {
	postFunc func(fn func())

	sync.Mutex
	view views.View

	duration time.Duration
	elapsedX int
	totalX   int
	offsetY  int
	runes    [3]rune
	styles   [3]tcell.Style

	toCall struct {
		util.Once
		*time.Timer
	}
	id   string
	quit <-chan struct{}
}

func NewWidgetProgress(postFunc func(fn func()), quit <-chan struct{}) *WidgetProgress {
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
	w.Lock()
	defer w.Unlock()
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
	defer w.Unlock()

	w.id = id
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
			w.Lock()
			defer w.Unlock()
			w.update()
		}()
	} else {
		go func() {
			w.Lock()
			defer w.Unlock()
			w.postFunc(w.Draw)
		}()
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
			w.Lock()
			defer w.Unlock()
			w.elapsedX += 1
			w.update()
		})
	}) {
		w.toCall.Reset(w.duration)
	}
}

func (w *WidgetProgress) Draw() {
	w.Lock()
	defer w.Unlock()
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
	defer w.Unlock()
	w.view = view
}

func (w *WidgetProgress) Resize() {
	w.Lock()
	defer w.Unlock()
	w.totalX, _ = w.view.Size()
}

func (w *WidgetProgress) Size() (int, int) {
	w.Lock()
	defer w.Unlock()
	x, _ := w.view.Size()
	return x, 1
}

func (*WidgetProgress) HandleEvent(tcell.Event) bool { return false }
