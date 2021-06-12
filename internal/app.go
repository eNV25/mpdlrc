package mpdlrc

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

	view   views.View
	widget views.Widget

	client  client.Client
	watcher client.Watcher
	song    song.Song
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
	app.widget = app.lyricsw
}

// Draw implements the root Widget.
func (app *Application) Draw() {
	if app.widget == app.lyricsw {
		app.lyricsw.SetLyrics(app.lyrics, -1)
	}
	app.widget.Draw()
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
				app.client.TogglePlay()
				return true
			}
		}
	case *events.PlayerEvent:
		switch app.client.State() {
		case state.PauseState:
			app.lyricsw.SetPaused(true)
		case state.PlayState:
			app.lyricsw.SetPaused(false)
		}
		app.SongChange(app.client.NowPlaying())
		app.Draw()
		return true
	case *events.TickerEvent:
		// no-op
		return true
	}
	return app.widget.HandleEvent(ev)
}

// SetScreen overrides views.Application.SetScreen.
func (app *Application) SetScreen(screen tcell.Screen) {
	app.Screen = screen
	app.Application.SetScreen(screen)
}

// SetView implements the root Widget.
func (app *Application) SetView(view views.View) {
	app.view = view
	app.widget.SetView(view)
}

func (app *Application) SongChange(song song.Song) {
	app.song = song
	if r, err := os.Open(
		path.Join(app.cfg.LyricsDir, app.song.LRCFile()),
	); err != nil {
		app.lyrics = lrc.NewLyrics(make([]time.Duration, 1), make([]string, 1)) // blank screen
	} else {
		if l, err := lrc.NewParser(r).Parse(); err != nil {
			panic(err)
		} else {
			app.lyrics = l
		}
	}
}

// Start overrides views.Application.Start.
func (app *Application) Start() {
	app.client.Start()
	app.SongChange(app.client.NowPlaying())
	app.Application.Start()
	go events.PostTickerEvents(app.PostEvent, 1*time.Second, app.quitch) // ticker events
	go app.watcher.PostEvents(app.PostEvent, app.quitch)                 // mpd player events
}

// Resize implements the root Widget.
func (app *Application) Resize() {
	app.widget.Resize()
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
	var screen tcell.Screen

	screen, err = tcell.NewScreen()
	if err != nil {
		return fmt.Errorf("allocate screen: %w", err)
	}

	app.SetScreen(screen)
	app.SetRootWidget(app)

	defer func() {
		app.Quit()
		err = app.Wait()
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
