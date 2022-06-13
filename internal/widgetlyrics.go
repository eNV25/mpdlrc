package internal

import (
	"sort"
	"sync"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/gdamore/tcell/v2/views"
	"github.com/neeharvi/textwidth"
	"github.com/rivo/uniseg"

	"github.com/env25/mpdlrc/internal/util"
)

var _ Widget = &WidgetLyrics{}

// LyricsWidget is a Widget implementation.
type WidgetLyrics struct {
	postFunc func(fn func())

	sync.Mutex
	view     views.View
	cellView *views.CellView

	toCall struct {
		util.Once
		*time.Timer
	}

	times   []time.Duration
	lines   []string
	elapsed time.Duration
	index   int
	total   int

	quit <-chan struct{}
}

// NewWidgetLyrics allocates new LyricsWidget.
func NewWidgetLyrics(postFunc func(fn func()), quit <-chan struct{}) *WidgetLyrics {
	w := &WidgetLyrics{
		postFunc: postFunc,
		quit:     quit,
	}
	w.cellView = views.NewCellView()
	w.cellView.Init()
	return w
}

func (w *WidgetLyrics) Cancel() {
	w.Lock()
	defer w.Unlock()
	if w.toCall.Timer != nil {
		w.toCall.Stop()
	}
}

func (w *WidgetLyrics) Update(
	playing bool,
	id string,
	elapsed time.Duration,
	times []time.Duration,
	lines []string,
) {
	w.Lock()
	defer w.Unlock()

	w.elapsed = elapsed
	w.times = times
	w.lines = lines

	w.total = len(w.lines)

	// This index is the first line after the one to be displayed.
	w.index = sort.Search(w.total, func(i int) bool { return w.times[i] > w.elapsed })

	if w.index < 0 || w.index > w.total {
		// This path is chosen when index is out of bounds for whatever reason.
		// Will display nothing. Will not start AfterFunc chain.

		w.index = 0
		w.total = 1
		w.lines = make([]string, 1)
		playing = false
	} else {
		// select previous line
		w.index--
	}

	select {
	case <-w.quit:
		return
	default:
	}

	if playing {
		go func() {
			w.Lock()
			defer w.Unlock()
			w.update()
		}()
	} else {
		go func() {
			w.Lock()
			defer w.Unlock()
			w.updateModel()
		}()
	}
}

func (w *WidgetLyrics) update() {
	select {
	case <-w.quit:
		return
	default:
		if w.index >= (w.total - 1) {
			return
		}
	}

	w.updateModel()

	if !w.toCall.Once.Do(func() {
		w.toCall.Timer = time.AfterFunc((w.times[w.index+1] - w.elapsed), func() {
			w.Lock()
			defer w.Unlock()
			w.index += 1
			w.elapsed = w.times[w.index]
			w.update()
		})
	}) {
		w.toCall.Reset((w.times[w.index+1] - w.elapsed))
	}
}

func (w *WidgetLyrics) updateModel() {
	m := &lyricsModel{}

	x, y := w.view.Size()
	mid := y / 2
	midoff := 0

	m.width = 0
	m.height = y + 1

	// nothing is highlighted when index is -1 like it should
	i1 := w.index - mid
	i2 := w.index + mid + 1

	hlStyle := tcell.StyleDefault.Attributes(tcell.AttrBold | tcell.AttrReverse)

	m.maincs = make([][]rune, m.height)
	m.combcs = make([][][]rune, m.height)
	m.widths = make([][]int, m.height)
	m.styles = make([][]tcell.Style, m.height)
	m.cells = make([]int, m.height)

	row := 0

	for ; i1 < 0; i1++ {
		row++
	}

	for i := i1; i < i2 && i < len(w.lines); i++ {
		width := textwidth.WidthString(w.lines[i])
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

		graphemes := uniseg.NewGraphemes(w.lines[i])

		for graphemes.Next() {
			runes := graphemes.Runes()
			wd := textwidth.WidthRunes(runes)

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
		m.styles[mid][cell] = hlStyle
	}

	w.cellView.SetModel(m)
	go w.postFunc(w.Draw)
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
	if x < 0 || y < 0 || y >= m.height || x >= m.cells[y] {
		return ' ', tcell.StyleDefault, nil, 1
	}
	return m.maincs[y][x], m.styles[y][x], m.combcs[y][x], m.widths[y][x]
}
func (m *lyricsModel) GetBounds() (int, int)             { return m.width, m.height }
func (m *lyricsModel) GetCursor() (int, int, bool, bool) { return 0, 0, false, false }
func (m *lyricsModel) MoveCursor(int, int)               {}
func (m *lyricsModel) SetCursor(int, int)                {}

func (w *WidgetLyrics) SetView(view views.View) {
	w.Lock()
	defer w.Unlock()
	w.view = view
	w.cellView.SetView(view)
}

func (w *WidgetLyrics) Draw() {
	w.Lock()
	defer w.Unlock()
	w.cellView.Draw()
}

func (w *WidgetLyrics) Resize() {
	w.Lock()
	defer w.Unlock()
	w.cellView.Resize()
}

func (w *WidgetLyrics) HandleEvent(ev tcell.Event) bool {
	return false
}

func (w *WidgetLyrics) Size() (int, int) {
	w.Lock()
	defer w.Unlock()
	return w.cellView.Size()
}
