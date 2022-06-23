package internal

import (
	"context"
	"sort"
	"time"

	"github.com/gdamore/tcell/v2"

	"github.com/env25/mpdlrc/internal/timerpool"
)

var _ Widget = &WidgetProgress{}

type WidgetProgress struct {
	widgetCommon

	totalX int
	*WidgetProgressData
}

type WidgetProgressData struct {
	Playing  bool
	Elapsed  time.Duration
	Duration time.Duration
	elapsedX int
	offsetY  int
}

func NewWidgetProgress(events chan<- tcell.Event) *WidgetProgress {
	ret := &WidgetProgress{}
	ret.events = events
	return ret
}

func (w *WidgetProgress) Update(ctx context.Context) {
	w.mu.Lock()
	defer w.mu.Unlock()

	// panic if not exist
	d := ctx.Value((*WidgetProgressData)(nil)).(*WidgetProgressData)
	_ = *d
	w.WidgetProgressData = d

	d.Duration = d.Duration / time.Duration(w.totalX)
	d.elapsedX = sort.Search(w.totalX, func(i int) bool { return (time.Duration(i) * d.Duration) >= d.Elapsed })

	w.update(ctx, d)
}

func (w *WidgetProgress) update(ctx context.Context, d *WidgetProgressData) {
	go func() {
		select {
		case <-ctx.Done():
			return
		case w.events <- NewEventFunction(w.Draw):
		}
	}()

	if !d.Playing || d.elapsedX+1 >= w.totalX {
		return
	}

	timer := timerpool.Get(d.Duration)
	go func() {
		select {
		case <-ctx.Done():
			timerpool.Put(timer, false)
			return
		case <-timer.C:
			timerpool.Put(timer, true)
		}

		w.mu.Lock()
		defer w.mu.Unlock()
		d.elapsedX += 1
		w.update(ctx, d)
	}()
}

func (w *WidgetProgress) Draw() {
	w.mu.Lock()
	defer w.mu.Unlock()

	if w.WidgetProgressData == nil {
		return
	}

	const (
		rune0 rune = '='
		rune1 rune = '>'
		rune2 rune = '-'
	)

	var (
		styleDefault = tcell.Style{}
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
