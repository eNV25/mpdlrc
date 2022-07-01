package widget

import (
	"context"

	"github.com/gdamore/tcell/v2"
	"github.com/mattn/go-runewidth"
	"github.com/rivo/uniseg"

	"github.com/env25/mpdlrc/internal/client"
	"github.com/env25/mpdlrc/internal/events"
	"github.com/env25/mpdlrc/internal/lyrics"
	"github.com/env25/mpdlrc/internal/panics"
	"github.com/env25/mpdlrc/internal/styles"
	"github.com/env25/mpdlrc/internal/urunewidth"
)

// [playing]     Logic/Russ - *Therapy Music* - Vinyl Days (2022)    [rzscxu]
//                 repeat random single consume crossfade update

var _ Widget = &Status{}

type Status struct {
	common
}

func NewStatus() *Status {
	ret := &Status{}
	return ret
}

type statusData struct {
	Song   client.Song
	Status client.Status
	Lyrics *lyrics.Lyrics
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
	defer panics.Handle(ctx)

	w.mu.Lock()
	defer w.mu.Unlock()

	song := client.SongFromContext(ctx)
	status := client.StatusFromContext(ctx)
	lyrics := lyrics.FromContext(ctx)

	d := &statusData{
		Song:   song,
		Status: status,
		Lyrics: lyrics,
	}

	go events.PostFunc(ctx, func() { w.draw(d) })
}

func (w *Status) draw(d *statusData) {
	w.mu.Lock()
	defer w.mu.Unlock()

	var (
		styleDefault = tcell.Style{}
		styleBold    = styleDefault.Bold(true)
	)

	vx, _ := w.Size()

	w.Fill(' ', styleDefault)

	{
		r := styles.BorderU
		s := styles.BorderStyle()
		for x := 0; x < vx; x++ {
			w.SetContent(x, 0, r, nil, s)
		}
	}

	{
		r := styles.BorderD
		s := styles.BorderStyle()
		for x := 0; x < vx; x++ {
			w.SetContent(x, 2, r, nil, s)
		}
	}

	{
		state := d.Status.State()
		var status []byte
		if state == "play" {
			status = append(status, "[playing] "...)
		} else if state == "pause" {
			status = append(status, "[paused] "...)
		} else if state == "stop" {
			status = append(status, "[stopped] "...)
		}
		if len(d.Lyrics.Lines) == 0 {
			status = append(status, "no lyrics "...)
		}
		for x, c := range status {
			w.SetContent(x, 1, rune(c), nil, styleDefault)
		}
	}

	{
		status := append([]byte{}, "[------]"...)
		if d.Status.Repeat() {
			status[1] = 'r'
		}
		if d.Status.Random() {
			status[2] = 'x'
		}
		if d.Status.Single() {
			status[3] = 's'
		}
		if d.Status.Consume() {
			status[4] = 'c'
		}
		for o, c := range status {
			w.SetContent(vx-len(status)+o, 1, rune(c), nil, styleDefault)
		}
	}

	{
		const margin = "  "
		pre := margin + d.Song.Artist() + " - "
		title := d.Song.Title()
		suf := " - " + d.Song.Album() + " " + "(" + d.Song.Date() + ")" + margin
		cs := [...]*struct {
			c string
			s tcell.Style
		}{
			{pre, styleDefault},
			{title, styleBold},
			{suf, styleDefault},
		}
		x := ((vx - runewidth.StringWidth(title)) / 2) - runewidth.StringWidth(pre)
		for _, c := range cs {
			gr := uniseg.NewGraphemes(c.c)
			for gr.Next() {
				rs := gr.Runes()
				x += urunewidth.GraphemeWidth(rs)
				w.SetContent(x, 1, rs[0], rs[1:], c.s)
			}
		}
	}
}
