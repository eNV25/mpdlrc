package internal

import (
	"sort"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/gdamore/tcell/v2/views"
	"github.com/neeharvi/textwidth"
	"github.com/rivo/uniseg"
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

	toCall *time.Timer
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

	total := len(lines)
	elapsed := status.Elapsed()

	// This is index is the first line after the one to be displayed.
	index := sort.Search(total, func(i int) bool { return times[i] >= elapsed })

	if index < 0 || index > total {
		// This path is chosen when index is out of bounds for whatever reason.
		// Will display nothing. Will not start AfterFunc chain.

		index = 0
		total = 1
		lines = make([]string, 1)
		playing = false
	} else if index == 0 || index == total {
		// This path is chosen when:
		//
		//  - index == total (no match from sort.Search)
		//    Will highlight last line by setting index to total-1,
		//    because all needed lines have already been shown.
		//
		//  - index == 0 (first march from sort.Search)
		//    Will highlight nothing by setting index to -1,
		//    because no line needs to be displayed yet.
		//
		// Will not start AfterFunc chain.

		// select previous line
		index--
		playing = false
	} else {
		// 0 < index < total
		// This path is chosen normally.
		// Will display previous line, which is the line to be shown.

		if index >= (total - 1) {
			return
		}

		// select previous line
		index--
	}

	if playing {
		w.update(times, lines, elapsed, index, total)
	} else {
		w.updateModel(lines, index)
	}
	w.postFunc(w.Draw)
}

func (w *LyricsWidget) update(times []time.Duration, lines []string, elapsed time.Duration, index int, total int) {
	w.updateModel(lines, index)
	w.postFunc(w.Draw)

	if index >= (total - 1) {
		return
	}

	w.toCall = time.AfterFunc((times[index+1] - elapsed), func() {
		index += 1
		elapsed = times[index]
		w.update(times, lines, elapsed, index, total)
	})
}

func (w *LyricsWidget) updateModel(lines []string, index int) {
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
		return ' ', tcell.StyleDefault, nil, 1
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
