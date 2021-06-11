package mpdlrc

import (
	"strings"
	"time"

	"local/mpdlrc/config"
	"local/mpdlrc/types"

	"github.com/gdamore/tcell/v2"
	"github.com/gdamore/tcell/v2/views"
)

// LyricsWidget is a Widget implentaion.
type LyricsWidget struct {
	*views.TextArea

	app    *Application
	cfg    *config.Config
	toCall *time.Timer // from AfterFunc
	scroll bool
	paused bool
}

// NewLyricsWidget allocates new LyricsWidget.
func NewLyricsWidget(app *Application) (ret *LyricsWidget) {
	ret = &LyricsWidget{
		TextArea: new(views.TextArea),
		app:      app,
		cfg:      app.cfg,
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

func (w *LyricsWidget) Scroll(d ScrollDirection) {
	ev := tcell.NewEventKey(tcell.Key(d), rune(0), tcell.ModMask(0))
	w.TextArea.HandleEvent(ev)
}

func (w *LyricsWidget) SetContent(text string) {
	w.SetLines(strings.Split(strings.ReplaceAll(text, "\r\n", "\n"), "\n"))
}

func (w *LyricsWidget) SetLines(lines []string) {
	x, y := w.app.Size()

	for i := range lines {
		offset := (x - len(lines[i])) / 2
		if offset < 0 {
			msg := "screen is too small"
			offset = (x - len(msg)) / 2
			lines = append(make([]string, (y/2)-1), (strings.Repeat(" ", offset) + msg))
			break
		}
		lines[i] = strings.Repeat(" ", offset+1) + lines[i]
	}

	w.TextArea.SetLines(lines)
}

func (w *LyricsWidget) SetLyrics(lyrics types.Lyrics, i int) {
	if w.paused || lyrics == nil {
		return
	}

	times := lyrics.Times()
	lines := lyrics.Lines()

	if i < 0 {
		if w.toCall != nil {
			w.toCall.Stop() // cancel
		}
		i = lyrics.Search(w.app.client.Elapsed()) - 1
		if i < 0 {
			i = 0
			lines = []string{""}
		}
	}

	_, y := w.app.Size()
	lines = append(make([]string, (y/2)-1), lines[i])
	w.SetLines(lines)

	if i >= (lyrics.N())-1 {
		return
	}

	w.toCall = time.AfterFunc(times[i+1]-times[i], func() {
		w.SetLyrics(lyrics, i+1)
	})
}

func (w *LyricsWidget) SetScroll(v bool) {
	w.scroll = v
}

func (w *LyricsWidget) SetPaused(v bool) {
	w.paused = v
}
