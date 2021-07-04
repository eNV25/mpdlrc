package internal

import (
	"sort"
	"time"

	"github.com/env25/mpdlrc/internal/status"

	"github.com/gdamore/tcell/v2"
	"github.com/gdamore/tcell/v2/views"
)

type ProgressWidget struct {
	app *Application

	view views.View

	toCall *time.Timer
	quit   chan struct{}

	elapsed  time.Duration
	duration time.Duration
	elapsedX int
	totalX   int
	offsetY  int
	runes    [3]rune
	styles   [3]tcell.Style
}

func NewProgressWidget(app *Application, quit chan struct{}) *ProgressWidget {
	w := &ProgressWidget{
		app:   app,
		quit:  quit,
		runes: [3]rune{'=', '>', ' '},
		styles: [3]tcell.Style{
			tcell.StyleDefault.Bold(true),
			tcell.StyleDefault.Bold(true),
			tcell.StyleDefault,
		},
	}
	return w
}

func (w *ProgressWidget) Cancel() {
	if w.toCall != nil {
		w.toCall.Stop()
	}
}

func (w *ProgressWidget) Update(status status.Status) {
	w.elapsed = status.Elapsed()
	w.duration = status.Duration() / time.Duration(w.totalX)
	w.elapsedX = sort.Search(w.totalX, func(i int) bool { return (time.Duration(i) * w.duration) >= w.elapsed }) - 1

	w.update()
}

func (w *ProgressWidget) update() {
	if w.elapsedX >= (w.totalX - 1) {
		return
	}

	w.elapsedX += 1

	w.toCall = time.AfterFunc(w.duration, func() {
		w.update()
		w.app.PostFunc(func() {
			w.app.Draw()
		})
	})
}

func (w *ProgressWidget) Draw() {
	w.view.Fill(w.runes[2], w.styles[2])
	for x := 0; x < w.elapsedX; x++ {
		w.view.SetContent(x, w.offsetY, w.runes[0], nil, w.styles[0])
	}
	w.view.SetContent(w.elapsedX, w.offsetY, w.runes[1], nil, w.styles[1])
}

func (w *ProgressWidget) SetView(view views.View) {
	w.view = view
}

func (w *ProgressWidget) Resize() {
	w.totalX, _ = w.view.Size()
}

func (w *ProgressWidget) Size() (int, int) { return w.view.Size() }

func (*ProgressWidget) HandleEvent(tcell.Event) bool { return false }
