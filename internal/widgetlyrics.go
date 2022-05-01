package internal

import (
	"sort"
	"sync"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/gdamore/tcell/v2/views"
	"github.com/neeharvi/textwidth"
	"github.com/rivo/uniseg"
	"golang.org/x/text/unicode/norm"
)

var _ Widget = &WidgetLyrics{}

// LyricsWidget is a Widget implementation.
type WidgetLyrics struct {
	sync.RWMutex
	postFunc func(fn func())

	view     views.View
	cellView *views.CellView

	toCall struct {
		sync.Mutex
		*time.Timer
	}
	id   string
	quit <-chan struct{}
}

// NewLyricsWidget allocates new LyricsWidget.
func NewLyricsWidget(postFunc func(fn func()), quit <-chan struct{}) *WidgetLyrics {
	w := &WidgetLyrics{
		postFunc: postFunc,
		quit:     quit,
	}
	w.cellView = views.NewCellView()
	w.cellView.Init()
	return w
}

func (w *WidgetLyrics) Cancel() {
	w.toCall.Lock()
	defer w.toCall.Unlock()
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
	w.id = id
	w.Unlock()

	w.RLock()
	defer w.RUnlock()

	for i := range lines {
		lines[i] = norm.NFC.String(lines[i])
	}

	total := len(lines)

	// This is index is the first line after the one to be displayed.
	index := sort.Search(total, func(i int) bool { return times[i] >= elapsed })

	if index < 0 || index > total {
		// This path is chosen when index is out of bounds for whatever reason.
		// Will display nothing. Will not start AfterFunc chain.

		index = 0
		total = 1
		lines = make([]string, 1)
		playing = false
	} else {
		// select previous line
		index--
	}

	select {
	case <-w.quit:
		return
	default:
		if w.id != id {
			return
		}
	}

	if playing {
		go w.update(id, times, lines, elapsed, index, total)
	} else {
		go func() {
			w.updateModel(lines, index)
			w.postFunc(w.Draw)
		}()
	}
}

func (w *WidgetLyrics) update(
	id string,
	times []time.Duration,
	lines []string,
	elapsed time.Duration,
	index int,
	total int,
) {
	w.RLock()
	defer w.RUnlock()

	select {
	case <-w.quit:
		return
	default:
		if w.id != id || index >= (total-1) {
			return
		}
	}

	w.updateModel(lines, index)
	w.postFunc(w.Draw)

	w.toCall.Lock()
	defer w.toCall.Unlock()
	w.toCall.Timer = time.AfterFunc((times[index+1] - elapsed), func() {
		index += 1
		elapsed = times[index]
		w.update(id, times, lines, elapsed, index, total)
	})
}

func (w *WidgetLyrics) updateModel(lines []string, index int) {
	m := &lyricsModel{}

	x, y := w.view.Size()
	mid := y / 2
	midoff := 0

	m.width = 0
	m.height = y + 1

	// nothing is highlighted when index is -1 like it should
	i1 := index - mid
	i2 := index + mid + 1

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

	for i := i1; i < i2 && i < len(lines); i++ {
		width := textwidth.WidthString(lines[i])
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

		graphemes := uniseg.NewGraphemes(lines[i])

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

	w.Lock()
	w.cellView.SetModel(m)
	w.Unlock()
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
func (m *lyricsModel) MoveCursor(int, int)               { return }
func (m *lyricsModel) SetCursor(int, int)                { return }

func (w *WidgetLyrics) SetView(view views.View) {
	w.Lock()
	w.view = view
	w.cellView.SetView(view)
	w.Unlock()
}

func (w *WidgetLyrics) Draw() {
	w.RLock()
	defer w.RUnlock()
	w.cellView.Draw()
}

func (w *WidgetLyrics) Resize() {
	w.Lock()
	w.cellView.Resize()
	w.Unlock()
}

func (w *WidgetLyrics) HandleEvent(ev tcell.Event) bool {
	return false
}

func (w *WidgetLyrics) Size() (int, int) {
	w.RLock()
	defer w.RUnlock()
	return w.cellView.Size()
}