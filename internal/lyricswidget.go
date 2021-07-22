package internal

import (
	"sort"
	"strings"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/gdamore/tcell/v2/views"
	"github.com/neeharvi/textwidth"
	"golang.org/x/text/width"

	"github.com/env25/mpdlrc/internal/status"
)

// LyricsWidget is a Widget implementation.
type LyricsWidget struct {
	postFunc func(fn func()) error

	view     views.View
	textArea *views.TextArea

	toCall  *time.Timer
	elapsed time.Duration
	times   []time.Duration
	lines   []string
	total   int
	index   int
	scroll  bool
}

// NewLyricsWidget allocates new LyricsWidget.
func NewLyricsWidget(postFunc func(fn func()) error) *LyricsWidget {
	w := &LyricsWidget{
		postFunc: postFunc,
		textArea: new(views.TextArea),
		scroll:   false,
	}
	w.textArea.Init()
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
		if w.index < 0 {
			w.SetLine("")
		} else {
			w.SetLine(w.lines[w.index])
		}
	}
	w.postFunc(w.Draw)
}

func (w *LyricsWidget) update() {
	if w.index < 0 {
		w.SetLine("")
	} else {
		w.SetLine(w.lines[w.index])
	}

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

func (w *LyricsWidget) SetLine(line string) {
	line = width.Fold.String(line)
	x, y := w.view.Size()
	offset := ((x - textwidth.StringWidth(line)) / 2) + 1
	if offset < 0 {
		offset = 1
	}
	lines := append(make([]string, ((y/2)-1), (y/2)), (strings.Repeat(" ", offset) + line))
	w.textArea.SetLines(lines)
}

// ScrollDirection represents scroll direction for Scroll methods.
type ScrollDirection tcell.Key

// Constants of type ScrollDirection.
const (
	ScrollUp    = ScrollDirection(tcell.KeyUp)
	ScrollDown  = ScrollDirection(tcell.KeyDown)
	ScrollRight = ScrollDirection(tcell.KeyRight)
	ScrollLeft  = ScrollDirection(tcell.KeyLeft)
)

func (w *LyricsWidget) Draw() {
	w.textArea.Draw()
}

func (w *LyricsWidget) Resize() {
	w.textArea.Resize()
}

func (w *LyricsWidget) HandleEvent(ev tcell.Event) bool {
	switch ev := ev.(type) {
	case *tcell.EventKey:
		if w.scroll {
			// scroll
			switch ev.Rune() {
			case 'k':
				w.Scroll(ScrollUp)
				return true
			case 'j':
				w.Scroll(ScrollDown)
				return true
			case 'l':
				w.Scroll(ScrollRight)
				return true
			case 'h':
				w.Scroll(ScrollLeft)
				return true
			}
		} else {
			// no scroll
			switch ev.Key() {
			case tcell.KeyUp, tcell.KeyDown, tcell.KeyRight, tcell.KeyLeft:
				return true
			}
		}
	}
	return w.textArea.HandleEvent(ev)
}

// Scroll in the direction represented by d.
func (w *LyricsWidget) Scroll(d ScrollDirection) {
	ev := tcell.NewEventKey(tcell.Key(d), rune(0), tcell.ModMask(0))
	w.textArea.HandleEvent(ev)
}

func (w *LyricsWidget) SetScroll(v bool) {
	w.scroll = v
}

func (w *LyricsWidget) SetView(view views.View) {
	w.view = view
	w.textArea.SetView(view)
}

func (w *LyricsWidget) Size() (int, int) {
	return w.textArea.Size()
}
