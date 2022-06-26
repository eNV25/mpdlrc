package internal

import (
	"context"
	"log"
	"path/filepath"
	"reflect"
	"runtime"
	"time"

	"github.com/gdamore/tcell/v2"
	"go.uber.org/multierr"

	"github.com/env25/mpdlrc/internal/client"
	"github.com/env25/mpdlrc/internal/config"
	"github.com/env25/mpdlrc/internal/event"
	"github.com/env25/mpdlrc/internal/events"
	"github.com/env25/mpdlrc/internal/lyrics"
	"github.com/env25/mpdlrc/internal/ufilepath"
	"github.com/env25/mpdlrc/internal/upath"
	"github.com/env25/mpdlrc/internal/widget"
)

// Application struct. Call (*Application).Run to run.
type Application struct {
	tcell.Screen

	bctx   context.Context
	cancel func()
	quit   func()
	events chan tcell.Event

	cfg     *config.Config
	client  *client.MPDClient
	watcher *client.MPDWatcher

	wprogress *widget.Progress
	wlyrics   *widget.Lyrics

	id     string
	lyrics *lyrics.Lyrics
}

// NewApplication allocates new Application from cfg.
func NewApplication(cfg *config.Config) *Application {
	app := &Application{
		cfg:       cfg,
		events:    make(chan tcell.Event),
		client:    client.NewMPDClient(cfg.MPD.Connection, cfg.MPD.Address, cfg.MPD.Password),
		watcher:   client.NewMPDWatcher(cfg.MPD.Connection, cfg.MPD.Address, cfg.MPD.Password),
		wprogress: widget.NewProgress(),
		wlyrics:   widget.NewLyrics(),
	}

	app.bctx, app.quit = context.WithCancel(context.Background())
	app.bctx = events.ContextWith(app.bctx, app.events)

	_, app.cancel = context.WithCancel(app.bctx)
	return app
}

// update subwidgets after querying information from client.
func (app *Application) update(ev tcell.Event) {
	app.cancel()

	song, _ := app.client.NowPlaying() // TODO
	status, _ := app.client.Status()   // TODO
	if song == nil || status == nil {
		return
	}

	if id := song.ID(); id != app.id {
		file := filepath.Join(app.cfg.LyricsDir, ufilepath.FromSlash(upath.ReplaceExt(song.File(), ".lrc")))
		app.lyrics = lyrics.New(file)
	}

	ctx := app.bctx
	ctx, app.cancel = context.WithCancel(ctx)
	ctx = event.ContextWith(ctx, ev)
	ctx = client.ContextWithSong(ctx, song)
	ctx = client.ContextWithStatus(ctx, status)
	ctx = lyrics.ContextWith(ctx, app.lyrics)

	go app.wprogress.Update(ctx)
	go app.wlyrics.Update(ctx)
}

// handleEvent handles dem events.
func (app *Application) handleEvent(ev tcell.Event) bool {
	if config.Debug {
		log.Printf("event: %T", ev)
	}
	var x, y int
	switch ev := ev.(type) {
	case *tcell.EventKey:
		switch ev.Key() {
		case tcell.KeyCtrlL:
			x, y = app.Screen.Size()
			goto resize
		case tcell.KeyCtrlC, tcell.KeyEscape:
			goto quit
		case tcell.KeyRune:
			switch ev.Rune() {
			case 'q':
				goto quit
			case ' ':
				return true
			}
		}
	case *tcell.EventResize:
		// guaranteed to run at program start
		x, y = ev.Size()
		goto resize
	case *event.Player:
		goto update
	case *event.Ping:
		_ = app.client.Ping()
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
	default:
	}
	return false
resize:
	app.resize(x, y)
	goto update
update:
	app.update(ev)
	return true
quit:
	app.Quit()
	return true
}

// resize is run after a resize event.
func (app *Application) resize(x, y int) {
	app.cancel()
	app.Screen.Fill(' ', tcell.Style{})
	app.Screen.Sync()
	app.wprogress.View().Resize(0, 0, x, 1)
	app.wlyrics.View().Resize(0, 1, x, y-1)
	app.wprogress.Resize()
	app.wlyrics.Resize()
}

// Quit the application.
func (app *Application) Quit() {
	app.quit()
}

// Run the application.
func (app *Application) Run() (err error) {
	app.Screen, err = tcell.NewScreen()
	if err != nil {
		return
	}

	err = app.Screen.Init()
	if err != nil {
		return
	}
	defer app.Screen.Fini()

	err = app.client.Start()
	if err != nil {
		return
	}
	defer multierr.AppendInvoke(&err, multierr.Invoke(app.client.Stop))

	err = app.watcher.Start()
	if err != nil {
		return
	}
	defer multierr.AppendInvoke(&err, multierr.Invoke(app.watcher.Stop))

	defer app.Quit()

	app.cfg.FromClient(app.client)
	if config.Debug {
		log.Print("\n", app.cfg)
	}

	// Screen.ChannelEvents closes events
	go app.Screen.ChannelEvents(app.events, app.bctx.Done())
	go app.watcher.PostEvents(app.bctx)
	go events.PostEveryTick(app.bctx, event.NewPing, 5*time.Second)

	app.wprogress.SetView(app.Screen)
	app.wlyrics.SetView(app.Screen)

	for ev := range app.events {
		app.handleEvent(ev)
		app.Show()
	}
	return
}
