package internal

import (
	"context"
	"sort"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/gdamore/tcell/v2/views"
	"github.com/mattn/go-runewidth"
	"github.com/rivo/uniseg"

	"github.com/env25/mpdlrc/internal/timerpool"
)

var _ Widget = &WidgetLyrics{}

// LyricsWidget is a Widget implementation.
type WidgetLyrics struct {
	widgetCommon

	cellView *views.CellView

	//*WidgetLyricsData /* not needed */
}

type WidgetLyricsData struct {
	Playing bool
	Times   []time.Duration
	Lines   []string
	Elapsed time.Duration
	index   int
	total   int
}

// NewWidgetLyrics allocates new LyricsWidget.
func NewWidgetLyrics(events chan<- tcell.Event) *WidgetLyrics {
	w := &WidgetLyrics{}
	w.events = events
	w.cellView = views.NewCellView()
	return w
}

func (w *WidgetLyrics) Update(ctx context.Context) {
	w.mu.Lock()
	defer w.mu.Unlock()

	// panic if not exist
	d := ctx.Value((*WidgetLyricsData)(nil)).(*WidgetLyricsData)
	_ = *d

	t := ctx.Value((*time.Time)(nil)).(time.Time)
	d.Elapsed += time.Since(t)

	// w.WidgetLyricsData = d /* not needed */

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

func (w *WidgetLyrics) update(ctx context.Context, d *WidgetLyricsData) {
	w.updateModel(d)

	go func() {
		select {
		case <-ctx.Done():
			return
		case w.events <- NewEventFunction(w.Draw):
		}
	}()

	if !d.Playing || d.index+1 >= d.total {
		return
	}

	timer := timerpool.Get(d.Times[d.index+1] - d.Elapsed)
	go func() {
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

func (w *WidgetLyrics) updateModel(d *WidgetLyricsData) {
	m := &lyricsModel{}

	x, y := w.view.Size()
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

			// StringWidth uses uniseg.Graphemes under the hood.
			// We copy its code to avoid running uniseg again.
			// wd := runewidth.StringWidth(grphms.Str())

			wd := 1
			for _, r := range runes {
				wd = runewidth.RuneWidth(r)
				if wd > 0 {
					break // Our best guess at this point is to use the width of the first non-zero-width rune.
				}
			}

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

	w.cellView.SetModel(m)
}

var _ views.CellModel = &lyricsModel{}

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
func (m *lyricsModel) GetBounds() (int, int)             { return m.width, m.height }
func (m *lyricsModel) GetCursor() (int, int, bool, bool) { return 0, 0, false, false }
func (m *lyricsModel) MoveCursor(int, int)               {}
func (m *lyricsModel) SetCursor(int, int)                {}

func (w *WidgetLyrics) SetView(view views.View) {
	w.mu.Lock()
	defer w.mu.Unlock()
	w.view = view
	w.cellView.SetView(view)
}

func (w *WidgetLyrics) Draw() {
	w.mu.Lock()
	defer w.mu.Unlock()
	w.cellView.Draw()
}

func (w *WidgetLyrics) Resize() {
	w.mu.Lock()
	defer w.mu.Unlock()
	w.cellView.Resize()
}

func (w *WidgetLyrics) Size() (int, int) {
	w.mu.Lock()
	defer w.mu.Unlock()
	return w.cellView.Size()
}
