package internal

import (
	"fmt"
	"os"
	"path"
	"time"

	"github.com/env25/mpdlrc/internal/client"
	"github.com/env25/mpdlrc/internal/config"
	"github.com/env25/mpdlrc/internal/events"
	"github.com/env25/mpdlrc/internal/lyrics"
	"github.com/env25/mpdlrc/internal/mpd"
	"github.com/env25/mpdlrc/internal/song"
	"github.com/env25/mpdlrc/internal/state"
	"github.com/env25/mpdlrc/internal/status"
	"github.com/env25/mpdlrc/lrc"

	"github.com/gdamore/tcell/v2"
	"github.com/gdamore/tcell/v2/views"
)

// Application struct. It embeds and overrides views.Application. It also implemets views.Widget
// so that it can be used as the root Widget. Call (*Application).Run to run.
type Application struct {
	tcell.Screen
	*views.Application

	WW *views.WidgetWatchers

	focused views.Widget

	client  client.Client
	watcher client.Watcher
	song    song.Song
	status  status.Status
	lyrics  lyrics.Lyrics
	cfg     *config.Config
	lyricsw *LyricsWidget
	quitch  chan struct{}
}

// NewApplication allocates new Application from cfg.
func NewApplication(cfg *config.Config) (app *Application) {
	app = new(Application)
	app.cfg = cfg
	app.Init()
	return app
}

// Init initiatialises Application. Can be called multiple times.
func (app *Application) Init() {
	app.Application = new(views.Application)
	app.WW = new(views.WidgetWatchers)
	app.quitch = make(chan struct{})
	app.client = mpd.NewMPDClient(app.cfg.MPD.Protocol, app.cfg.MPD.Address)
	app.watcher = mpd.NewMPDWatcher(app.cfg.MPD.Protocol, app.cfg.MPD.Address)
	app.lyricsw = NewLyricsWidget(app)
	app.focused = app.lyricsw
}

// Draw implements the root Widget.
func (app *Application) Draw() {
	if app.focused == app.lyricsw {
		app.lyricsw.Draw()
	}
}

// HandleEvent implements the root Widget.
func (app *Application) HandleEvent(ev tcell.Event) bool {
	switch ev := ev.(type) {
	case *tcell.EventKey:
		switch ev.Key() {
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
	case *events.PlayerEvent:
		app.Update()
		app.Draw()
		return true
	case *events.TickerEvent:
		app.client.Ping()
		return true
	}
	return app.focused.HandleEvent(ev)
}

// SetView implements the root Widget.
func (app *Application) SetView(view views.View) {
	app.lyricsw.SetView(view)
}

func (app *Application) Lyrics(song song.Song) lyrics.Lyrics {
	if r, err := os.Open(
		path.Join(app.cfg.LyricsDir, app.song.LRCFile()),
	); err != nil {
		return lrc.NewLyrics(make([]time.Duration, 1), make([]string, 1)) // blank screen
	} else {
		if l, err := lrc.NewParser(r).Parse(); err != nil {
			return lrc.NewLyrics(make([]time.Duration, 1), make([]string, 1)) // blank screen
		} else {
			return l
		}
	}
}

func (app *Application) Update() {
	app.song = app.client.NowPlaying()
	app.lyrics = app.Lyrics(app.song)
	app.status = app.client.Status()
	switch app.status.State() {
	case state.PauseState:
		app.lyricsw.SetPaused(true)
	case state.PlayState:
		app.lyricsw.SetPaused(false)
	}
	app.lyricsw.Update(app.status, app.lyrics)
}

// Start overrides views.Application.Start.
func (app *Application) Start() {
	app.client.Start()
	app.Application.Start()
	go events.PostTickerEvents(app.PostEvent, 1*time.Second, app.quitch) // ticker events
	go app.watcher.PostEvents(app.PostEvent, app.quitch)                 // mpd player events
}

// Resize implements the root Widget.
func (app *Application) Resize() {
	app.Update()
	app.lyricsw.Resize()
}

// Quit performs shotdown steps, overrides views.Application.Quit.
func (app *Application) Quit() {
	app.Application.Quit()
	{
		// if already closed; no-op
		// else; close
		select {
		case _, ok := <-app.quitch:
			if ok {
				close(app.quitch)
			}
		default:
			close(app.quitch)
		}
	}
	app.client.Stop()
}

// Run runs the application, overrides views.Application.Run.
func (app *Application) Run() (err error) {
	app.Screen, err = tcell.NewScreen()
	if err != nil {
		return fmt.Errorf("allocate screen: %w", err)
	}

	app.Application.SetScreen(app.Screen)
	app.Application.SetRootWidget(app)

	defer func() {
		if err == nil {
			app.Quit()
			err = app.Wait()
		}
	}()

	app.Start()
	return app.Wait()
}

// Unwatch implements the root Widget.
func (app *Application) Unwatch(handler tcell.EventHandler) {
	app.WW.Unwatch(handler)
}

// Watch implements the root Widget.
func (app *Application) Watch(handler tcell.EventHandler) {
	app.WW.Watch(handler)
}
