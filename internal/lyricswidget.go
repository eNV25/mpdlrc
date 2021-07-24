package internal

import (
	"sort"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/gdamore/tcell/v2/views"
	"github.com/neeharvi/textwidth"
	"golang.org/x/text/unicode/norm"

	"github.com/env25/mpdlrc/internal/status"
	"github.com/env25/mpdlrc/internal/widget"
)

var _ widget.Widget = &LyricsWidget{}

// LyricsWidget is a Widget implementation.
type LyricsWidget struct {
	postFunc func(fn func()) error

	view     views.View
	cellView *views.CellView

	toCall  *time.Timer
	elapsed time.Duration
	times   []time.Duration
	lines   []string
	total   int
	index   int
}

// NewLyricsWidget allocates new LyricsWidget.
func NewLyricsWidget(postFunc func(fn func()) error) *LyricsWidget {
	w := &LyricsWidget{
		postFunc: postFunc,
		cellView: views.NewCellView(),
	}
	w.cellView.Init()
	return w
}

func (w *LyricsWidget) Cancel() {
	if w.toCall != nil {
		w.toCall.Stop()
	}
}

func (w *LyricsWidget) Update(playing bool, status status.Status, times []time.Duration, lines []string) {
	if status == nil || times == nil || lines == nil {
		return
	}

	for i := range lines {
		lines[i] = norm.NFC.String(lines[i])
	}

	w.lines = lines
	w.times = times
	w.total = len(lines)

	w.elapsed = status.Elapsed()
	w.index = sort.Search(w.total, func(i int) bool { return w.times[i] >= w.elapsed })

	if w.index < 0 || w.index >= w.total {
		w.index = 0
		w.total = 1
		w.lines = make([]string, 1)
	} else {
		w.index -= 1
	}

	if w.index >= (w.total - 1) {
		return
	}

	if playing {
		w.update()
	} else {
		w.updateModel(w.lines, w.index)
	}
	w.postFunc(w.Draw)
}

func (w *LyricsWidget) update() {
	w.updateModel(w.lines, w.index)

	if w.index >= (w.total - 1) {
		return
	}

	w.toCall = time.AfterFunc((w.times[w.index+1] - w.elapsed), func() {
		w.index += 1
		w.elapsed = w.times[w.index]
		w.update()
		w.postFunc(w.Draw)
	})
}

func (w *LyricsWidget) updateModel(lines []string, index int) {
	m := &lyricsModel{}

	x, y := w.view.Size()

	m.width = 0
	m.height = y + 1

	m.maincs = make([][]rune, m.height)
	m.combcs = make([][][]rune, m.height)
	m.widths = make([][]int, m.height)
	m.styles = make([][]tcell.Style, m.height)
	m.cells = make([]int, m.height)

	i1 := (index - (y / 2))
	i2 := (index + (y / 2)) + 1

	row := 0

	for ; i1 < 0; i1++ {
		row++
	}

	for i := i1; i < i2 && i < len(lines); i++ {
		n := textwidth.WidthString(lines[i])
		off := ((x - n) / 2) + 1
		if off < 0 {
			off = 1
		}
		n += off

		cell := 0

		m.maincs[row] = make([]rune, n)
		m.combcs[row] = make([][]rune, n)
		m.widths[row] = make([]int, n)
		m.styles[row] = make([]tcell.Style, n)

		for ; off > 0; off-- {
			m.maincs[row][cell] = ' '
			m.widths[row][cell] = 1
			cell += 1
		}

		for _, r := range lines[i] {
			w := textwidth.WidthRune(r)

			m.maincs[row][cell] = r
			m.widths[row][cell] = w

			if row == (y / 2) {
				m.styles[row][cell] = tcell.StyleDefault.Bold(true).Reverse(true)
			}

			cell += w
		}

		m.cells[row] = cell

		if cell > m.width {
			m.width = cell
		}

		row++
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
	if x < 0 || y < 0 || y >= m.height || x >= m.cells[y] {
		return 0, tcell.StyleDefault, nil, 1
	}
	return m.maincs[y][x], m.styles[y][x], m.combcs[y][x], m.widths[y][x]
}
func (m *lyricsModel) GetBounds() (int, int)             { return m.width, m.height }
func (m *lyricsModel) GetCursor() (int, int, bool, bool) { return 0, 0, false, false }
func (m *lyricsModel) MoveCursor(int, int)               { return }
func (m *lyricsModel) SetCursor(int, int)                { return }

func (w *LyricsWidget) SetView(view views.View) {
	w.view = view
	w.cellView.SetView(view)
}

func (w *LyricsWidget) Draw() {
	w.cellView.Draw()
}

func (w *LyricsWidget) Resize() {
	w.cellView.Resize()
}

func (w *LyricsWidget) HandleEvent(ev tcell.Event) bool {
	return false
}

func (w *LyricsWidget) Size() (int, int) {
	return w.cellView.Size()
}
