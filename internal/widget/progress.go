package widget

import (
	"context"
	"sort"
	"time"

	"github.com/env25/mpdlrc/internal/client"
	"github.com/env25/mpdlrc/internal/event"
	"github.com/env25/mpdlrc/internal/events"
	"github.com/env25/mpdlrc/internal/panics"
	"github.com/env25/mpdlrc/internal/styles"
	"github.com/env25/mpdlrc/internal/timerpool"
)

var _ Widget = &Progress{}

type Progress struct {
	common
}

type progressData struct {
	Playing  bool
	Elapsed  time.Duration
	Duration time.Duration
	elapsedX int
	totalX   int
}

func NewProgress() *Progress {
	ret := &Progress{}
	return ret
}

func (w *Progress) Update(ctx context.Context) {
	defer panics.Handle(ctx)

	w.mu.Lock()
	defer w.mu.Unlock()

	vx, _ := w.Size()
	status := client.StatusFromContext(ctx)

	d := &progressData{
		Playing:  status.State() == "play",
		Elapsed:  status.Elapsed(),
		Duration: status.Duration(),
		totalX:   vx,
	}

	d.Elapsed += time.Since(event.FromContext(ctx).When())

	d.Duration = d.Duration / time.Duration(vx)
	d.elapsedX = sort.Search(d.totalX, func(i int) bool { return (time.Duration(i) * d.Duration) >= d.Elapsed })

	w.update(ctx, d)
}

func (w *Progress) update(ctx context.Context, d *progressData) {
	go events.PostFunc(ctx, func() { w.draw(d) })

	if !d.Playing || d.elapsedX+1 >= d.totalX {
		return
	}

	timer := timerpool.Get(d.Duration)
	go func() {
		defer panics.Handle(ctx)

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

	w.Fill(' ', styles.Default())

	{
		r := styles.BorderU
		s := styles.Border()
		for x := 0; x < d.totalX; x++ {
			w.SetContent(x, 0, r, nil, s)
		}
	}

	for x := 0; x < d.elapsedX; x++ {
		w.SetContent(x, 1, '=', nil, styles.Default().Bold(true))
	}
	w.SetContent(d.elapsedX, 1, '>', nil, styles.Default().Bold(true))
	for x := d.elapsedX + 1; x < d.totalX; x++ {
		w.SetContent(x, 1, '-', nil, styles.Default().Dim(true))
	}

	{
		r := styles.BorderD
		s := styles.Border()
		for x := 0; x < d.totalX; x++ {
			w.SetContent(x, 2, r, nil, s)
		}
	}
}
