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
	"github.com/env25/mpdlrc/internal/panics"
	"github.com/env25/mpdlrc/internal/ufilepath"
	"github.com/env25/mpdlrc/internal/upath"
	"github.com/env25/mpdlrc/internal/widget"
)

// Application struct. Call (*Application).Run to run.
type Application struct {
	tcell.Screen

	events chan tcell.Event
	quit   func()
	cancel func()

	cfg     *config.Config
	client  *client.MPDClient
	watcher *client.MPDWatcher

	wprogress *widget.Progress
	wlyrics   *widget.Lyrics
	wstatus   *widget.Status

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
		wstatus:   widget.NewStatus(),
		lyrics:    &lyrics.Lyrics{},
	}

	return app
}

// update subwidgets after querying information from client.
func (app *Application) update(ctx context.Context) {
	app.cancel()

	song, err := app.client.NowPlaying() // TODO
	status, errr := app.client.Status()  // TODO
	if err != nil || errr != nil {
		log.Printf("%+v\n", multierr.Append(err, errr))
		return
	}

	if id := song.ID(); id != app.id {
		file := filepath.Join(app.cfg.LyricsDir, ufilepath.FromSlash(upath.ReplaceExt(song.File(), ".lrc")))
		app.lyrics = lyrics.New(file)
	}

	ctx = client.ContextWithSong(ctx, song)
	ctx = client.ContextWithStatus(ctx, status)
	ctx = lyrics.ContextWith(ctx, app.lyrics)
	ctx, app.cancel = context.WithCancel(ctx)

	go app.wprogress.Update(ctx)
	go app.wlyrics.Update(ctx)
	go app.wstatus.Update(ctx)
}

// handleEvent handles dem events.
func (app *Application) handleEvent(ctx context.Context, ev tcell.Event) bool {
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
	ctx = event.ContextWith(ctx, ev)
	app.update(ctx)
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
	app.wprogress.View().Resize(0, 0, x, 3)
	app.wlyrics.View().Resize(0, 3, x, y-6)
	app.wstatus.View().Resize(0, y-3, x, 3)
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

	// Update config with data from MPD
	app.cfg.FromClient(app.client)
	if config.Debug {
		log.Print("\n", app.cfg)
	}

	ctx := context.Background()
	ctx = panics.ContextWithHook(ctx, app.Quit)
	ctx = events.ContextWith(ctx, app.events)
	ctx, app.quit = context.WithCancel(ctx)

	// We make sure this function value is never nil
	app.cancel = func() {}

	go app.Screen.ChannelEvents(app.events, ctx.Done())
	go app.watcher.PostEvents(ctx)
	go events.PostEveryTick(ctx, event.NewPing, 5*time.Second)

	app.wprogress.SetView(app.Screen)
	app.wlyrics.SetView(app.Screen)
	app.wstatus.SetView(app.Screen)

	for ev := range app.events {
		app.handleEvent(ctx, ev)
		app.Show()
	}
	return
}
