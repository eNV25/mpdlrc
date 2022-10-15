package widget

import (
	"context"
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/mattn/go-runewidth"
	"github.com/rivo/uniseg"

	"github.com/env25/mpdlrc/internal/client"
	"github.com/env25/mpdlrc/internal/events"
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

type statusData struct {
	client.Data
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

func (w *Status) Update(ctx context.Context, ev tcell.Event) {
	defer panics.Handle(ctx)

	switch ev.(type) {
	case *tcell.EventResize:
		w.resize()
	case *client.OptionsEvent:
		// no-op
	case *client.PlayerEvent:
		// no-op
	default:
		return
	}

	w.mu.Lock()
	defer w.mu.Unlock()

	data := client.DataFromContext(ctx)

	d := &statusData{data}

	go events.PostFunc(ctx, func() { w.draw(d) })
}

func (w *Status) draw(d *statusData) {
	w.mu.Lock()
	defer w.mu.Unlock()

	vx, _ := w.Size()

	w.Fill(' ', styles.Default())

	{
		r := styles.RuneBorderUpper
		s := styles.Border()
		for x := 0; x < vx; x++ {
			w.SetContent(x, 0, r, nil, s)
		}
	}

	{
		r := styles.RuneBorderLower
		s := styles.Border()
		for x := 0; x < vx; x++ {
			w.SetContent(x, 2, r, nil, s)
		}
	}

	{
		state := d.State()
		var status []byte
		if state == "play" {
			status = append(status, "[playing] "...)
		} else if state == "pause" {
			status = append(status, "[paused] "...)
		} else if state == "stop" {
			status = append(status, "[stopped] "...)
		}
		if d.Lyrics == nil || len(d.Lines) == 0 {
			status = append(status, "no lyrics "...)
		}
		for x, c := range status {
			w.SetContent(x, 1, rune(c), nil, styles.Default())
		}
	}

	{
		status := append([]byte{}, "[------]"...)
		if d.Repeat() {
			status[1] = 'r'
		}
		if d.Random() {
			status[2] = 'x'
		}
		if d.Single() {
			status[3] = 's'
		}
		if d.Consume() {
			status[4] = 'c'
		}
		for o, c := range status {
			w.SetContent(vx-len(status)+o, 1, rune(c), nil, styles.Default())
		}
	}

	{
		var pre strings.Builder
		var title strings.Builder
		var suf strings.Builder
		pre.WriteString("  ")
		if s := d.Artist(); s != "" {
			pre.WriteString(s)
			pre.WriteString(" - ")
		}
		if s := d.Title(); s != "" {
			title.WriteString(s)
		} else {
			title.WriteString(d.Song.File())
		}
		if s := d.Album(); s != "" {
			suf.WriteString(" - ")
			suf.WriteString(s)
		}
		if s := d.Date(); s != "" {
			suf.WriteString(" (")
			suf.WriteString(s)
			suf.WriteString(")")
		}
		suf.WriteString("  ")
		x := (vx-runewidth.StringWidth(title.String()))/2 - runewidth.StringWidth(pre.String())
		for _, cs := range &[...]*struct {
			c string
			s tcell.Style
		}{
			{pre.String(), styles.Default()},
			{title.String(), styles.Default().Bold(true)},
			{suf.String(), styles.Default()},
		} {
			for clstr, st := "", -1; len(cs.c) > 0; {
				clstr, cs.c, _, st = uniseg.FirstGraphemeClusterInString(cs.c, st)
				rs := []rune(clstr)
				w.SetContent(x, 1, rs[0], rs[1:], cs.s)
				x += urunewidth.GraphemeWidth(rs)
			}
		}
	}
}
