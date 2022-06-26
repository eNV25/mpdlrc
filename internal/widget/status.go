package widget

import (
	"context"

	"github.com/gdamore/tcell/v2"
	"github.com/mattn/go-runewidth"
	"github.com/rivo/uniseg"

	"github.com/env25/mpdlrc/internal/client"
	"github.com/env25/mpdlrc/internal/events"
	"github.com/env25/mpdlrc/internal/urunewidth"
)

// [playing]     Logic/Russ - *Therapy Music* - Vinyl Days (2022)    [rzscxu]
//                 repeat random single consume crossfade update

var _ Widget = &Status{}

type Status struct {
	common

	totalX int
}

func NewStatus() *Status {
	ret := &Status{}
	return ret
}

type statusData struct {
	Song   client.Song
	Status client.Status
	// Album string
	// Artist string
	// Title string
	// Date string
	// Filename string
	// State string
	// Repeat string
	// Random string
	// Single string
	// Consume string
}

func (w *Status) Update(ctx context.Context) {
	w.mu.Lock()
	defer w.mu.Unlock()

	song := client.SongFromContext(ctx)
	status := client.StatusFromContext(ctx)

	d := &statusData{
		Song:   song,
		Status: status,
	}

	go events.PostFunc(ctx, func() { w.draw(d) })
}

func (w *Status) draw(d *statusData) {
	w.mu.Lock()
	defer w.mu.Unlock()

	styleDefault := tcell.Style{} // styleBold    = styleDefault.Bold(true)

	title := d.Song.Artist() + " - " + d.Song.Title() + " - " + d.Song.Album() + "(" + d.Song.Date() + ")"
	width := runewidth.StringWidth(title)
	off := (w.totalX - width) / 2

	x := off

	gr := uniseg.NewGraphemes(title)
	for gr.Next() {
		rs := gr.Runes()
		wd := urunewidth.GraphemeWidth(rs)
		w.SetContent(x, 0, rs[0], rs[1:], styleDefault)
		x += wd
	}
}

func (w *Status) Resize() {
	w.common.Resize()
	w.mu.Lock()
	defer w.mu.Unlock()
	w.totalX, _ = w.ViewPort.Size()
}

func (w *Status) Size() (int, int) {
	w.mu.Lock()
	defer w.mu.Unlock()
	return w.totalX, 1
}
