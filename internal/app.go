package internal

import (
	"log"
	"os"
	"path"
	"reflect"
	"runtime"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/gdamore/tcell/v2/views"

	"github.com/env25/mpdlrc/internal/app/client"
	"github.com/env25/mpdlrc/internal/app/event"
	"github.com/env25/mpdlrc/internal/app/mpd"
	"github.com/env25/mpdlrc/internal/app/song"
	"github.com/env25/mpdlrc/internal/app/state"
	"github.com/env25/mpdlrc/internal/app/status"
	"github.com/env25/mpdlrc/internal/app/widget"
	"github.com/env25/mpdlrc/internal/config"
	"github.com/env25/mpdlrc/lrc"
)

// Application struct. It embeds and overrides views.Application. It also implemets views.Widget
// so that it can be used as the root Widget. Call (*Application).Run to run.
type Application struct {
	tcell.Screen

	cfg *config.Config

	client  client.Client
	watcher client.Watcher
	song    song.Song
	status  status.Status
	times   []time.Duration
	lines   []string
	id      string
	playing bool

	focused   widget.Widget
	lyricsv   *views.ViewPort
	lyricsw   *LyricsWidget
	progressv *views.ViewPort
	progressw *ProgressWidget

	quit   chan struct{}
	events chan tcell.Event
}

// NewApplication allocates new Application from cfg.
func NewApplication(cfg *config.Config) *Application {
	app := &Application{
		cfg:    cfg,
		quit:   make(chan struct{}),
		events: make(chan tcell.Event),
	}

	app.client = mpd.NewMPDClient(cfg.MPD.Connection, cfg.MPD.Address, cfg.MPD.Password)
	app.watcher = mpd.NewMPDWatcher(cfg.MPD.Connection, cfg.MPD.Address, cfg.MPD.Password)

	app.lyricsw = NewLyricsWidget(app.postFunc, app.quit)
	app.progressw = NewProgressWidget(app.postFunc, app.quit)
	app.focused = app.lyricsw

	return app
}

// Update subwidgets after querying information from client.
func (app *Application) Update() {
	song, status := app.client.NowPlaying(), app.client.Status()
	if song != nil && status != nil {
		app.song, app.status = song, status
	} else {
		return
	}

	switch app.status.State() {
	case state.Play:
		app.playing = true
	case state.Pause:
		app.playing = false
	default:
		// nothing to do
		return
	}

	app.progressw.Cancel()
	app.lyricsw.Cancel()

	if id := app.song.ID(); id != app.id {
		app.id = id
		app.times, app.lines = app.lyrics(app.song)
	}

	app.progressw.Update(app.playing, app.id, app.status.Elapsed(), app.status.Duration())
	app.lyricsw.Update(app.playing, app.id, app.status.Elapsed(), app.times, app.lines)
}

// Resize is run after a resize event.
func (app *Application) Resize() {
	app.SetView(app.Screen)
	app.progressw.Resize()
	app.lyricsw.Resize()
}

// HandleEvent handles dem events.
func (app *Application) HandleEvent(ev tcell.Event) bool {
	if config.Debug {
		log.Printf("event: %T", ev)
	}
	switch ev := ev.(type) {
	case *tcell.EventKey:
		switch ev.Key() {
		case tcell.KeyCtrlL:
			// fake resize event
			return app.HandleEvent(tcell.NewEventResize(app.Size()))
		case tcell.KeyCtrlC, tcell.KeyEscape:
			app.Quit()
			return true
		case tcell.KeyRune:
			switch ev.Rune() {
			case 'q':
				app.Quit()
				return true
			case ' ':
				return true
			}
		}
	case *tcell.EventResize:
		// guaranteed to run at program start
		app.Screen.Fill(' ', tcell.StyleDefault)
		app.Screen.Sync()
		app.Resize()
		app.Update()
		return true
	case *event.Player:
		app.Update()
		return true
	case *event.Ping:
		app.client.Ping()
		return true
	case *event.Function:
		if config.Debug {
			log.Println(
				"event: *event.Function: ev.Func:",
				runtime.FuncForPC(reflect.ValueOf(ev.Func).Pointer()).Name(),
			)
		}
		ev.Func()
		return true
	}
	return app.focused.HandleEvent(ev)
}

// postFunc runs function fn in the event loop. uses an unbuffered channel.
func (app *Application) postFunc(fn func()) {
	app.events <- event.NewFunctionEvent(fn)
}

// SetView updates the views of subwidgets.
func (app *Application) SetView(view views.View) {
	if app.lyricsv == nil || app.progressv == nil {
		// init
		app.progressv = views.NewViewPort(view, 0, 0, 0, 0)
		app.progressw.SetView(app.progressv)
		app.lyricsv = views.NewViewPort(view, 0, 0, 0, 0)
		app.lyricsw.SetView(app.lyricsv)
	}
	app.progressv.Resize(0, 0, -1, 1)
	app.lyricsv.Resize(0, 1, -1, -1)
}

// lyrics fetches lyrics using information from song.
func (app *Application) lyrics(song song.Song) ([]time.Duration, []string) {
	if r, err := os.Open(
		path.Join(app.cfg.LyricsDir, app.song.LRCFile()),
	); err != nil {
		return make([]time.Duration, 1), make([]string, 1) // blank screen
	} else {
		if times, lines, err := lrc.ParseReader(r); err != nil {
			return make([]time.Duration, 1), make([]string, 1) // blank screen
		} else {
			return times, lines
		}
	}
}

// Quit the application.
func (app *Application) Quit() {
	close(app.quit)
}

// Run the application.
func (app *Application) Run() (err error) {

	app.Screen, err = tcell.NewScreen()
	if err != nil {
		goto quit
	}

	err = app.Screen.Init()
	if err != nil {
		goto quit
	}
	defer app.Screen.Fini()

	err = app.client.Start()
	if err != nil {
		goto quit
	}
	defer app.client.Stop()

	err = app.watcher.Start()
	if err != nil {
		goto quit
	}
	defer app.watcher.Stop()

	go app.Screen.ChannelEvents(app.events, app.quit)
	go app.watcher.PostEvents(app.events, app.quit)
	go event.PostTickerEvents(app.events, 5*time.Second, event.NewPingEvent, app.quit)

	for ev := range app.events {
		app.HandleEvent(ev)
		app.Show()
	}
	return

quit:
	app.Quit()
	return
}
