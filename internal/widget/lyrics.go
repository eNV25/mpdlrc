package widget

import (
	"context"
	"sort"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/uniseg"

	"github.com/env25/mpdlrc/internal/client"
	"github.com/env25/mpdlrc/internal/event"
	"github.com/env25/mpdlrc/internal/events"
	"github.com/env25/mpdlrc/internal/lyrics"
	"github.com/env25/mpdlrc/internal/panics"
	"github.com/env25/mpdlrc/internal/styles"
	"github.com/env25/mpdlrc/internal/timerpool"
	"github.com/env25/mpdlrc/internal/urunewidth"
)

var _ Widget = &Lyrics{}

// LyricsWidget is a Widget implementation.
type Lyrics struct {
	common

	//*WidgetLyricsData /* not needed */
}

type lyricsData struct {
	Playing bool
	Times   []time.Duration
	Lines   []string
	Elapsed time.Duration
	index   int
	total   int
}

// NewLyrics allocates new LyricsWidget.
func NewLyrics() *Lyrics {
	w := &Lyrics{}
	return w
}

func (w *Lyrics) Update(ctx context.Context) {
	defer panics.Handle(ctx)

	w.mu.Lock()
	defer w.mu.Unlock()

	status := client.StatusFromContext(ctx)
	lyrics := lyrics.FromContext(ctx)

	d := &lyricsData{
		Playing: status.State() == "play",
		Elapsed: status.Elapsed(),
		Times:   lyrics.Times,
		Lines:   lyrics.Lines,
	}

	d.Elapsed += time.Since(event.FromContext(ctx).When())

	d.total = len(d.Lines)

	// This index is the first line after the one to be displayed.
	d.index = sort.Search(d.total, func(i int) bool { return d.Times[i] > d.Elapsed })

	if d.index < 0 || d.index > d.total {
		// This path is chosen when index is out of bounds for whatever reason.
		// Will display nothing. Will not start AfterFunc chain.

		d.index = 0
		d.total = 1
		d.Lines = []string{}
		d.Playing = false
	} else {
		// select previous line
		d.index--
	}

	w.update(ctx, d)
}

func (w *Lyrics) update(ctx context.Context, d *lyricsData) {
	m := w.model(d)

	go events.PostFunc(ctx, func() { w.draw(m) })

	if !d.Playing || d.index+1 >= d.total {
		return
	}

	timer := timerpool.Get(d.Times[d.index+1] - d.Elapsed)
	go func() {
		defer panics.Handle(ctx)

		select {
		case <-ctx.Done():
			timerpool.Put(timer, false)
			return
		case t := <-timer.C:
			timerpool.Put(timer, true)

			w.mu.Lock()
			defer w.mu.Unlock()
			d.index += 1
			d.Elapsed = d.Times[d.index] + time.Since(t)
			w.update(ctx, d)
		}
	}()
}

func (w *Lyrics) model(d *lyricsData) *lyricsModel {
	m := &lyricsModel{}

	vx, vy := w.Size()
	ymid := vy / 2

	m.width = 0
	m.height = vy + 1

	// nothing is highlighted when index is -1 like it should
	i1 := d.index - ymid
	i2 := d.index + ymid + 1

	m.vx = vx
	m.xwidth = make([]int, m.height)
	m.combcs = make([][][]rune, m.height)
	m.widths = make([][]int, m.height)
	m.styles = make([]tcell.Style, m.height)

	y := 0

	for ; i1 < 0; i1++ {
		y++
	}

	for i := i1; i < i2 && i < len(d.Lines); i++ {
		x := 0
		max := len(d.Lines[i]) * 2
		m.combcs[y] = make([][]rune, max)
		m.widths[y] = make([]int, max)

		gr := uniseg.NewGraphemes(d.Lines[i])

		for wd := 0; gr.Next(); x += wd {
			rs := gr.Runes()
			wd = urunewidth.GraphemeWidth(rs)
			m.combcs[y][x] = rs
			m.widths[y][x] = wd
		}

		m.xwidth[y] = x

		if x > m.width {
			m.width = x
		}

		y++
	}

	m.styles[ymid] = styles.Default().Bold(true).Reverse(true)

	return m
}

var _ cellModel = &lyricsModel{}

type lyricsModel struct {
	vx     int
	width  int
	height int
	xwidth []int
	combcs [][][]rune
	styles []tcell.Style
	widths [][]int
}

func (m *lyricsModel) GetCell(x, y int) (rune, []rune, tcell.Style, int) {
	x = x - (m.vx-m.xwidth[y])/2 // centre
	if y < 0 || x < 0 || y >= m.height || x >= m.xwidth[y] {
		return ' ', nil, styles.Default(), 1
	}
	return m.combcs[y][x][0], m.combcs[y][x][1:], m.styles[y], m.widths[y][x]
}
func (m *lyricsModel) GetBounds() (int, int) { return m.width, m.height }

func (w *Lyrics) draw(m *lyricsModel) {
	w.mu.Lock()
	defer w.mu.Unlock()

	w.Fill(' ', styles.Default())

	ex, ey := m.GetBounds()
	vx, vy := w.Size()
	if ex < vx {
		ex = vx
	}
	if ey < vy {
		ey = vy
	}

	for y := 0; y < ey; y++ {
		for x := 0; x < ex; {
			ch, comb, style, wid := m.GetCell(x, y)
			w.SetContent(x, y, ch, comb, style)
			x += wid
		}
	}
}
