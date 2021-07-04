package internal

import (
	"fmt"
	"os"
	"path"
	"time"

	"github.com/env25/mpdlrc/internal/client"
	"github.com/env25/mpdlrc/internal/config"
	"github.com/env25/mpdlrc/internal/events"
	"github.com/env25/mpdlrc/internal/mpd"
	"github.com/env25/mpdlrc/internal/song"
	"github.com/env25/mpdlrc/internal/state"
	"github.com/env25/mpdlrc/internal/status"
	"github.com/env25/mpdlrc/internal/widget"
	"github.com/env25/mpdlrc/lrc"

	"github.com/gdamore/tcell/v2"
	"github.com/gdamore/tcell/v2/views"
)

// Application struct. It embeds and overrides views.Application. It also implemets views.Widget
// so that it can be used as the root Widget. Call (*Application).Run to run.
type Application struct {
	tcell.Screen

	client  client.Client
	watcher client.Watcher
	song    song.Song
	status  status.Status
	cfg     *config.Config

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

	app.client = mpd.NewMPDClient(cfg.MPD.Protocol, cfg.MPD.Address, cfg.MPD.Password)
	app.watcher = mpd.NewMPDWatcher(cfg.MPD.Protocol, cfg.MPD.Address, cfg.MPD.Password)

	app.lyricsw = NewLyricsWidget(app, app.quit)
	app.progressw = NewProgressWidget(app, app.quit)
	app.focused = app.lyricsw

	return app
}

// Draw implements the root Widget.
func (app *Application) Draw() {
	app.progressw.Draw()
	app.lyricsw.Draw()
}

// Update subwidgets after querying information from client.
func (app *Application) Update() {
	app.song = app.client.NowPlaying()
	app.status = app.client.Status()

	app.progressw.Cancel()
	app.lyricsw.Cancel()
	switch app.status.State() {
	case state.PlayState:
		app.progressw.Update(app.status)
		times, lines := app.Lyrics(app.song)
		app.lyricsw.Update(app.status, times, lines)
	}
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
		app.Resize()
		app.Update()
		app.Screen.Fill(' ', tcell.StyleDefault)
		app.Screen.Sync()
		app.Draw()
		app.Screen.Sync()
		return true
	case *events.PlayerEvent:
		app.Update()
		app.Draw()
		return true
	case *events.PingEvent:
		go app.client.Ping()
		return true
	case *events.FunctionEvent:
		ev.Run()
		return true
	}
	return app.focused.HandleEvent(ev)
}

// PostFunc runs function fn in the event loop.
func (app *Application) PostFunc(fn func()) error {
	ev := events.NewFunctionEvent(fn)
	return app.PostEvent(ev)
}

// SetView updates the views of subwidgets.
func (app *Application) SetView(view views.View) {
	if app.lyricsv == nil {
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
	// NOTE: put all shutdown actions under the select case
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

	defer app.client.Stop()
	defer app.Screen.Fini()

	app.PostFunc(app.Update)
	app.PostFunc(app.Draw)

	go app.Screen.ChannelEvents(app.events, app.quit)
	go app.watcher.PostEvents(
		app.PostEvent, app.quit)
	go events.PostTickerEvents(
		app.PostEvent, 5*time.Second,
		events.NewPingEvent, app.quit)

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
