package internal

import (
	"fmt"
	"os"
	"path"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/gdamore/tcell/v2/views"

	"github.com/env25/mpdlrc/internal/client"
	"github.com/env25/mpdlrc/internal/config"
	"github.com/env25/mpdlrc/internal/event"
	"github.com/env25/mpdlrc/internal/mpd"
	"github.com/env25/mpdlrc/internal/song"
	"github.com/env25/mpdlrc/internal/state"
	"github.com/env25/mpdlrc/internal/status"
	"github.com/env25/mpdlrc/internal/widget"
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
	id      interface{}
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

	app.lyricsw = NewLyricsWidget(app.PostFunc)
	app.progressw = NewProgressWidget(app.PostFunc)
	app.focused = app.lyricsw

	return app
}

// Update subwidgets after querying information from client.
func (app *Application) Update() {
	app.song = app.client.NowPlaying()
	app.status = app.client.Status()

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
		app.times, app.lines = app.Lyrics(app.song)
	}

	app.progressw.Update(app.playing, app.status)
	app.lyricsw.Update(app.playing, app.status, app.times, app.lines)
}

// Resize is run after a resize event.
func (app *Application) Resize() {
	app.SetView(app.Screen)
	app.progressw.Resize()
	app.lyricsw.Resize()
}

// HandleEvent handles dem events.
func (app *Application) HandleEvent(ev tcell.Event) bool {
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
		ev.Run()
		return true
	}
	return app.focused.HandleEvent(ev)
}

// PostFunc runs function fn in the event loop.
func (app *Application) PostFunc(fn func()) error {
	ev := event.NewFunctionEvent(fn)
	return app.PostEvent(ev)
}

// SetView updates the views of subwidgets.
func (app *Application) SetView(view views.View) {
	if app.lyricsv == nil || app.progressv == nil {
		app.progressv = views.NewViewPort(view, 0, 0, 0, 0)
		app.progressw.SetView(app.progressv)
		app.lyricsv = views.NewViewPort(view, 0, 0, 0, 0)
		app.lyricsw.SetView(app.lyricsv)
	}
	app.progressv.Resize(0, 0, -1, 1)
	app.lyricsv.Resize(0, 1, -1, -1)
}

// Lyrics fetches lyrics using information from song.
func (app *Application) Lyrics(song song.Song) ([]time.Duration, []string) {
	if r, err := os.Open(
		path.Join(app.cfg.LyricsDir, app.song.LRCFile()),
	); err != nil {
		// TODO: better error messages
		return make([]time.Duration, 1), make([]string, 1) // blank screen
	} else {
		if times, lines, err := lrc.NewParser(r).Parse(); err != nil {
			// TODO: better error messages
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
func (app *Application) Run() error {
	var err error

	app.Screen, err = tcell.NewScreen()
	if err != nil {
		err = fmt.Errorf("new screen: %w", err)
		goto quit
	}

	err = app.Screen.Init()
	if err != nil {
		err = fmt.Errorf("init screen: %w", err)
		goto quit
	}

	err = app.client.Start()
	if err != nil {
		err = fmt.Errorf("starting client: %w", err)
		goto quit
	}

	err = app.watcher.Start()
	if err != nil {
		err = fmt.Errorf("starting watcher: %w", err)
		goto quit
	}

	defer app.watcher.Stop()
	defer app.client.Stop()
	defer app.Screen.Fini()

	app.PostFunc(app.Update)

	go app.Screen.ChannelEvents(app.events, app.quit)
	go app.watcher.PostEvents(
		app.PostEvent, app.quit)
	go event.PostTickerEvents(
		app.PostEvent, 5*time.Second,
		event.NewPingEvent, app.quit)

	for {
		app.Show()

		select {
		case <-app.quit:
			goto rtrn
		case ev := <-app.events:
			app.HandleEvent(ev)
		}
	}

quit:
	close(app.quit)
rtrn:
	return err
}
