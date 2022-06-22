package internal

import (
	"context"
	"log"
	"os"
	"path/filepath"
	"reflect"
	"runtime"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/gdamore/tcell/v2/views"
	"go.uber.org/multierr"

	"github.com/env25/mpdlrc/internal/config"
	"github.com/env25/mpdlrc/lrc"
)

// Application struct. Call (*Application).Run to run.
type Application struct {
	tcell.Screen

	bctx   context.Context
	cancel func()
	quit   func()
	events chan tcell.Event

	cfg     *config.Config
	client  ClientType
	watcher WatcherType

	wlyrics   *WidgetLyrics
	wprogress *WidgetProgress

	id    string
	times []time.Duration
	lines []string
}

// NewApplication allocates new Application from cfg.
func NewApplication(cfg *config.Config) *Application {
	app := &Application{
		cfg:    cfg,
		events: make(chan tcell.Event),
	}

	app.client = NewMPDClient(cfg.MPD.Connection, cfg.MPD.Address, cfg.MPD.Password)
	app.watcher = NewMPDWatcher(cfg.MPD.Connection, cfg.MPD.Address, cfg.MPD.Password)

	app.wlyrics = NewWidgetLyrics(app.events)
	app.wprogress = NewWidgetProgress(app.events)

	app.bctx, app.quit = context.WithCancel(context.Background())
	_, app.cancel = context.WithCancel(app.bctx)
	return app
}

// Update subwidgets after querying information from client.
func (app *Application) Update(ctx context.Context) {
	song, status := app.client.NowPlaying(), app.client.Status()
	if song == nil || status == nil {
		return
	}

	var (
		playing  bool
		elapsed  = status.Elapsed()
		duration = status.Duration()
	)

	// cancel previous context
	app.cancel()

	ctx, app.cancel = context.WithCancel(ctx)

	switch status.State() {
	case StatePlay:
		playing = true
	default:
		playing = false
	}

	if id := song.ID(); id != app.id {
		app.id = id
		app.times, app.lines = app.lyrics(song)
	}

	go app.wprogress.Update(context.WithValue(ctx, (*WidgetProgressData)(nil), &WidgetProgressData{
		Playing:  playing,
		Elapsed:  elapsed,
		Duration: duration,
	}))
	go app.wlyrics.Update(context.WithValue(ctx, (*WidgetLyricsData)(nil), &WidgetLyricsData{
		Playing: playing,
		Elapsed: elapsed,
		Times:   app.times,
		Lines:   app.lines,
	}))
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
		app.Resize()
		return true
	case *EventPlayer:
		app.Update(app.bctx)
		return true
	case *EventPing:
		app.client.Ping()
		return true
	case *EventFunction:
		if config.Debug {
			log.Println(
				"event: *event.Function: ev.Func:",
				runtime.FuncForPC(reflect.ValueOf(ev.Func).Pointer()).Name(),
			)
		}
		ev.Func()
		return true
	}
	return false
}

// Resize is run after a resize event.
func (app *Application) Resize() {
	app.Screen.Fill(' ', tcell.Style{})
	app.Screen.Sync()
	app.wprogress.View().Resize(0, 0, -1, 1)
	app.wlyrics.View().Resize(0, 1, -1, -1)
	app.wprogress.Resize()
	app.wlyrics.Resize()
	app.Update(app.bctx)
}

// lyrics fetches lyrics using information from song.
func (app *Application) lyrics(song SongType) ([]time.Duration, []string) {
	if r, err := os.Open(
		filepath.Join(app.cfg.LyricsDir, song.LRCFile()),
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
	app.quit()
}

// Run the application.
func (app *Application) Run() (err error) {
	defer func() {
	}()

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
	defer func() { err = multierr.Append(err, app.client.Stop()) }()

	err = app.watcher.Start()
	if err != nil {
		return
	}
	defer func() { err = multierr.Append(err, app.watcher.Stop()) }()

	defer app.Quit()

	// Screen.ChannelEvents closes events
	go app.Screen.ChannelEvents(app.events, app.bctx.Done())
	go app.watcher.PostEvents(app.bctx, app.events)
	go sendNewEventEvery(app.bctx, app.events, NewEventPing, 5*time.Second)

	app.wlyrics.SetView(views.NewViewPort(app.Screen, 0, 0, 0, 0))
	app.wprogress.SetView(views.NewViewPort(app.Screen, 0, 0, 0, 0))

	for ev := range app.events {
		app.HandleEvent(ev)
		app.Show()
	}
	return
}
