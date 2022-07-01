package widget

import (
	"context"
	"sort"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/mattn/go-runewidth"
	"github.com/rivo/uniseg"

	"github.com/env25/mpdlrc/internal/client"
	"github.com/env25/mpdlrc/internal/event"
	"github.com/env25/mpdlrc/internal/events"
	"github.com/env25/mpdlrc/internal/lyrics"
	"github.com/env25/mpdlrc/internal/panics"
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

		var t time.Time
		select {
		case <-ctx.Done():
			timerpool.Put(timer, false)
			return
		case t = <-timer.C:
			timerpool.Put(timer, true)
		}

		w.mu.Lock()
		defer w.mu.Unlock()
		d.index += 1
		d.Elapsed = d.Times[d.index] + time.Since(t)
		w.update(ctx, d)
	}()
}

func (w *Lyrics) model(d *lyricsData) *lyricsModel {
	m := &lyricsModel{}

	x, y := w.Size()
	mid := y / 2
	midoff := 0

	m.width = 0
	m.height = y + 1

	// nothing is highlighted when index is -1 like it should
	i1 := d.index - mid
	i2 := d.index + mid + 1

	var (
		styleDefault = tcell.Style{}
		styleHl      = styleDefault.Bold(true).Reverse(true)
	)

	m.maincs = make([][]rune, m.height)
	m.combcs = make([][][]rune, m.height)
	m.widths = make([][]int, m.height)
	m.styles = make([][]tcell.Style, m.height)
	m.cells = make([]int, m.height)

	row := 0

	for ; i1 < 0; i1++ {
		row++
	}

	for i := i1; i < i2 && i < len(d.Lines); i++ {
		// TODO: StringWidth uses uniseg.Graphemes under the hood.
		//       Can we do this without running uniseg twice?
		width := runewidth.StringWidth(d.Lines[i])
		off := (x - width) / 2
		if off < 0 {
			off = 0
		}
		width += off

		m.maincs[row] = make([]rune, width)
		m.combcs[row] = make([][]rune, width)
		m.widths[row] = make([]int, width)
		m.styles[row] = make([]tcell.Style, width)

		cell := 0

		if row == mid {
			midoff = off
		}

		for ; off > 0; off-- {
			m.maincs[row][cell] = ' '
			m.widths[row][cell] = 1
			cell += 1
		}

		grphms := uniseg.NewGraphemes(d.Lines[i])

		for grphms.Next() {
			runes := grphms.Runes()

			wd := urunewidth.GraphemeWidth(runes)

			m.maincs[row][cell] = runes[0]
			m.combcs[row][cell] = runes[1:]
			m.widths[row][cell] = wd

			cell += wd
		}

		m.cells[row] = cell

		if cell > m.width {
			m.width = cell
		}

		row++
	}

	for cell := midoff; cell < len(m.maincs[mid]); cell++ {
		m.styles[mid][cell] = styleHl
	}

	return m
}

var _ cellModel = &lyricsModel{}

type lyricsModel struct {
	maincs [][]rune
	combcs [][][]rune
	widths [][]int
	styles [][]tcell.Style
	cells  []int
	width  int
	height int
}

func (m *lyricsModel) GetCell(x, y int) (rune, tcell.Style, []rune, int) {
	var styleDefault tcell.Style
	if x < 0 || y < 0 || y >= m.height || x >= m.cells[y] {
		return ' ', styleDefault, nil, 1
	}
	return m.maincs[y][x], m.styles[y][x], m.combcs[y][x], m.widths[y][x]
}
func (m *lyricsModel) GetBounds() (int, int) { return m.width, m.height }

func (w *Lyrics) draw(model *lyricsModel) {
	w.mu.Lock()
	defer w.mu.Unlock()

	styleDefault := tcell.Style{}

	w.Fill(' ', styleDefault)

	ex, ey := model.GetBounds()
	vx, vy := w.Size()
	if ex < vx {
		ex = vx
	}
	if ey < vy {
		ey = vy
	}

	for y := 0; y < ey; y++ {
		for x := 0; x < ex; {
			ch, style, comb, wid := model.GetCell(x, y)
			w.SetContent(x, y, ch, comb, style)
			x += wid
		}
	}
}
