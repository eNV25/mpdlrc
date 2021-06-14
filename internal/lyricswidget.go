package internal

import (
	"strings"
	"time"

	"github.com/env25/mpdlrc/internal/client"
	"github.com/env25/mpdlrc/internal/config"
	"github.com/env25/mpdlrc/internal/lyrics"

	"github.com/gdamore/tcell/v2"
	"github.com/gdamore/tcell/v2/views"
)

// LyricsWidget is a Widget implementation.
type LyricsWidget struct {
	*views.TextArea

	client client.Client

	app     *Application
	cfg     *config.Config
	toCall  *time.Timer // from AfterFunc
	elapsed time.Duration
	scroll  bool
	paused  bool
}

// NewLyricsWidget allocates new LyricsWidget.
func NewLyricsWidget(app *Application) (ret *LyricsWidget) {
	ret = &LyricsWidget{
		TextArea: new(views.TextArea),
		app:      app,
		cfg:      app.cfg,
		client:   app.client,
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
	x, _ := w.app.Size()

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

// SetLyrics sets a particular line i of lyrics to be displayed.  Each call sets
// an AfterFunc for the next line that needs to be displayed, so this
// method only needs to be called when the lyrics change. If i is -1 it cancels
// the AfterFunc and queries the current time from the client.
func (w *LyricsWidget) SetLyrics(lyrics lyrics.Lyrics, i int) {
	if w.paused || lyrics == nil {
		return
	}

	times := lyrics.Times()
	lines := lyrics.Lines()

	if i < 0 {
		if w.toCall != nil {
			w.toCall.Stop() // cancel
		}
		w.elapsed = w.client.Elapsed()
		i = lyrics.Search(w.elapsed) - 1
		if i < 0 {
			// blank screen
			i = 0
			lines = []string{""}
		}
	}

	_, y := w.app.Size()
	lines = append(make([]string, (y/2)-1), lines[i]) // centre line
	w.SetLines(lines)

	if i >= (lyrics.N())-1 {
		return
	}

	w.toCall = time.AfterFunc(times[i+1]-w.elapsed, func() {
		w.SetLyrics(lyrics, i+1)
	})

	w.elapsed = w.elapsed + (times[i+1] - w.elapsed)
}

func (w *LyricsWidget) SetScroll(v bool) {
	w.scroll = v
}

func (w *LyricsWidget) SetPaused(v bool) {
	w.paused = v
}
