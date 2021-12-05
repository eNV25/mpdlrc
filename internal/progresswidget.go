package internal

import (
	"sort"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/gdamore/tcell/v2/views"

	"github.com/env25/mpdlrc/internal/status"
)

type ProgressWidget struct {
	postFunc func(fn func())

	view views.View

	toCall *time.Timer
	quit   chan struct{}

	elapsedX int
	totalX   int
	offsetY  int
	runes    [3]rune
	styles   [3]tcell.Style
}

func NewProgressWidget(postFunc func(fn func())) *ProgressWidget {
	return &ProgressWidget{
		postFunc: postFunc,
		runes:    [3]rune{'=', '>', '-'},
		styles: [3]tcell.Style{
			tcell.StyleDefault.Attributes(tcell.AttrBold),
			tcell.StyleDefault.Attributes(tcell.AttrBold),
			tcell.StyleDefault.Attributes(tcell.AttrDim),
		},
	}
}

func (w *ProgressWidget) Cancel() {
	if w.toCall != nil {
		w.toCall.Stop()
	}
}

func (w *ProgressWidget) Update(playing bool, status status.Status) {
	elapsed := status.Elapsed()
	duration := status.Duration() / time.Duration(w.totalX)
	w.elapsedX = sort.Search(w.totalX, func(i int) bool { return (time.Duration(i) * duration) >= elapsed })

	if w.elapsedX >= w.totalX {
		return
	}

	if playing {
		go w.update(duration)
	} else {
		go func() {
			w.postFunc(w.Draw)
		}()
	}
}

func (w *ProgressWidget) update(duration time.Duration) {
	if w.elapsedX >= w.totalX {
		return
	}

	w.toCall = time.AfterFunc(duration, func() {
		w.elapsedX += 1
		w.update(duration)
		w.postFunc(w.Draw)
	})
}

func (w *ProgressWidget) Draw() {
	w.view.Fill(' ', tcell.StyleDefault)
	for x := 0; x < w.elapsedX; x++ {
		w.view.SetContent(x, w.offsetY, w.runes[0], nil, w.styles[0])
	}
	w.view.SetContent(w.elapsedX, w.offsetY, w.runes[1], nil, w.styles[1])
	for x := w.elapsedX + 1; x < w.totalX; x++ {
		w.view.SetContent(x, w.offsetY, w.runes[2], nil, w.styles[2])
	}
}

func (w *ProgressWidget) SetView(view views.View) {
	w.view = view
}

func (w *ProgressWidget) Resize() {
	w.totalX, _ = w.view.Size()
}

func (w *ProgressWidget) Size() (int, int) {
	x, _ := w.view.Size()
	return x, 1
}

func (*ProgressWidget) HandleEvent(tcell.Event) bool { return false }
