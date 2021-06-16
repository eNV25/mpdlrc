package internal

import (
	"strings"
	"time"

	"github.com/env25/mpdlrc/internal/lyrics"
	"github.com/env25/mpdlrc/internal/status"

	"github.com/gdamore/tcell/v2"
	"github.com/gdamore/tcell/v2/views"
)

// LyricsWidget is a Widget implementation.
type LyricsWidget struct {
	*views.TextArea

	view views.View

	lyrics  lyrics.Lyrics
	app     *Application
	toCall  *time.Timer // from AfterFunc
	elapsed time.Duration
	scroll  bool
	paused  bool
	index   int
}

// NewLyricsWidget allocates new LyricsWidget.
func NewLyricsWidget(app *Application) (ret *LyricsWidget) {
	ret = &LyricsWidget{
		TextArea: new(views.TextArea),
		app:      app,
		scroll:   false,
		paused:   false,
	}
	ret.Init()
	return ret
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
	return w.TextArea.HandleEvent(ev)
}

// Scroll in the direction represented by d.
func (w *LyricsWidget) Scroll(d ScrollDirection) {
	ev := tcell.NewEventKey(tcell.Key(d), rune(0), tcell.ModMask(0))
	w.TextArea.HandleEvent(ev)
}

func (w *LyricsWidget) SetContent(text string) {
	w.SetLines(strings.Split(strings.ReplaceAll(text, "\r\n", "\n"), "\n"))
}

func (w *LyricsWidget) SetLines(lines []string) {
	x, _ := w.view.Size()

	for i := range lines {
		offset := (x - len(lines[i])) / 2
		if offset < 0 {
			lines = make([]string, 1) // empty
			break
		}
		lines[i] = strings.Repeat(" ", offset+1) + lines[i] // centre line
	}

	w.TextArea.SetLines(lines)
}

func (w *LyricsWidget) Update(status status.Status, lyrics lyrics.Lyrics) {
	if w.paused {
		return
	}

	if status != nil && lyrics != nil {
		if w.toCall != nil {
			w.toCall.Stop() // cancel
		}
		w.lyrics = lyrics
		w.elapsed = status.Elapsed()
		w.index = lyrics.Search(w.elapsed) - 1
	}

	times := w.lyrics.Times()
	lines := w.lyrics.Lines()

	if w.index < 0 {
		// blank screen
		w.index = 0
		lines = []string{""}
	}

	_, y := w.view.Size()
	lines = append(make([]string, y/2), lines[w.index]) // centre line
	w.SetLines(lines)

	if w.index >= len(times)-1 {
		return
	}

	w.toCall = time.AfterFunc(times[w.index+1]-w.elapsed, func() {
		w.index += 1
		w.elapsed = times[w.index]
		w.Update(nil, nil)
	})
}

func (w *LyricsWidget) SetScroll(v bool) {
	w.scroll = v
}

func (w *LyricsWidget) SetPaused(v bool) {
	w.paused = v
}

func (w *LyricsWidget) SetView(view views.View) {
	w.view = view
	w.TextArea.SetView(view)
}
