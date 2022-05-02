package internal

import (
	"sort"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/gdamore/tcell/v2/views"

	"github.com/env25/mpdlrc/internal/util"
)

type WidgetProgress struct {
	postFunc func(fn func())

	view views.View

	duration time.Duration
	elapsedX int
	totalX   int
	offsetY  int
	runes    [3]rune
	styles   [3]tcell.Style

	toCall struct {
		once util.Once
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
		go w.update()
	} else {
		go func() {
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

	if !w.toCall.once.Do(func() {
		w.toCall.Timer = time.AfterFunc(w.duration, func() {
			w.elapsedX += 1
			w.update()
		})
	}) {
		w.toCall.Reset(w.duration)
	}
}

func (w *WidgetProgress) Draw() {
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
	w.view = view
}

func (w *WidgetProgress) Resize() {
	w.totalX, _ = w.view.Size()
}

func (w *WidgetProgress) Size() (int, int) {
	x, _ := w.view.Size()
	return x, 1
}

func (*WidgetProgress) HandleEvent(tcell.Event) bool { return false }
