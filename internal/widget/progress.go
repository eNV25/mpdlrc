package widget

import (
	"context"
	"sort"
	"time"

	"github.com/gdamore/tcell/v2"

	"github.com/env25/mpdlrc/internal/client"
	"github.com/env25/mpdlrc/internal/event"
	"github.com/env25/mpdlrc/internal/events"
	"github.com/env25/mpdlrc/internal/timerpool"
)

var _ Widget = &Progress{}

type Progress struct {
	common

	totalX int
}

type progressData struct {
	Playing  bool
	Elapsed  time.Duration
	Duration time.Duration
	elapsedX int
}

func NewProgress() *Progress {
	ret := &Progress{}
	return ret
}

func (w *Progress) Update(ctx context.Context) {
	w.mu.Lock()
	defer w.mu.Unlock()

	status := client.StatusFromContext(ctx)

	d := &progressData{
		Playing:  status.State() == "play",
		Elapsed:  status.Elapsed(),
		Duration: status.Duration(),
	}

	d.Elapsed += time.Since(event.FromContext(ctx).When())

	d.Duration = d.Duration / time.Duration(w.totalX)
	d.elapsedX = sort.Search(w.totalX, func(i int) bool { return (time.Duration(i) * d.Duration) >= d.Elapsed })

	w.update(ctx, d)
}

func (w *Progress) update(ctx context.Context, d *progressData) {
	go events.PostFunc(ctx, func() { w.draw(d) })

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

func (w *Progress) draw(d *progressData) {
	w.mu.Lock()
	defer w.mu.Unlock()

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

	w.Fill(' ', styleDefault)
	for x := 0; x < d.elapsedX; x++ {
		w.SetContent(x, 0, rune0, nil, style0)
	}
	w.SetContent(d.elapsedX, 0, rune1, nil, style1)
	for x := d.elapsedX + 1; x < w.totalX; x++ {
		w.SetContent(x, 0, rune2, nil, style2)
	}
}

func (w *Progress) Resize() {
	w.common.Resize()
	w.mu.Lock()
	defer w.mu.Unlock()
	w.totalX, _ = w.ViewPort.Size()
}

func (w *Progress) Size() (int, int) {
	w.mu.Lock()
	defer w.mu.Unlock()
	return w.totalX, 1
}
