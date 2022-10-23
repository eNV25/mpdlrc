// Package internal contains the [Application] struct and other packages.
package internal

import (
	"context"
	"log"
	"reflect"
	"runtime"

	"github.com/gdamore/tcell/v2"

	"github.com/env25/mpdlrc/internal/client"
	"github.com/env25/mpdlrc/internal/config"
	"github.com/env25/mpdlrc/internal/event"
	"github.com/env25/mpdlrc/internal/events"
	"github.com/env25/mpdlrc/internal/panics"
	"github.com/env25/mpdlrc/internal/widget"
)

// Application struct. Call (*Application).Run to run.
type Application struct {
	tcell.Screen

	wprogress widget.Progress
	wlyrics   widget.Lyrics
	wstatus   widget.Status

	cfg    *config.Config
	client *client.MPDClient
	events chan tcell.Event
	quit   func()
	cancel func()
}

// NewApplication allocates new Application from cfg.
func NewApplication(cfg *config.Config, client *client.MPDClient) *Application {
	app := &Application{
		cfg:    cfg,
		client: client,
		events: make(chan tcell.Event),
		quit:   noop,
		cancel: noop,
	}
	return app
}

func noop() {}

func (app *Application) update(ctx context.Context, ev tcell.Event) {
	if config.Debug {
		switch ev := ev.(type) {
		case *event.Func:
			log.Printf("update: %T: %s", ev, runtime.FuncForPC(reflect.ValueOf(ev.Func).Pointer()).Name())
		case *tcell.EventKey:
			switch ev.Key() {
			case tcell.KeyRune:
				log.Printf("update: %T: %q", ev, string(ev.Rune()))
			default:
				log.Printf("update: %T: %s", ev, tcell.KeyNames[ev.Key()])
			}
		default:
			log.Printf("update: %T", ev)
		}
	}
	switch ev := ev.(type) {
	case *client.PlayerEvent:
		app.updateData(ctx, ev, ev.Data)
	case *client.OptionsEvent:
		app.updateData(ctx, ev, ev.Data)
	case *event.Func:
		ev.Func()
	case *tcell.EventKey:
		switch ev.Key() {
		case tcell.KeyCtrlL:
			x, y := app.Screen.Size()
			app.updateResize(ctx, ev, x, y)
		case tcell.KeyCtrlC, tcell.KeyEscape:
			app.Quit()
		case tcell.KeyRune:
			switch ev.Rune() {
			case 'q':
				app.Quit()
			case 'p':
				app.client.TogglePause()
			case ' ':
				app.client.TogglePause()
			}
		}
	case *tcell.EventResize:
		// guaranteed to run at program start
		x, y := ev.Size()
		app.updateResize(ctx, ev, x, y)
	}
}

func (app *Application) updateData(ctx context.Context, ev tcell.Event, data client.Data) {
	app.cancel()

	ctx = client.ContextWithData(ctx, data)
	ctx, app.cancel = context.WithCancel(ctx)

	go app.wprogress.Update(ctx, ev)
	go app.wlyrics.Update(ctx, ev)
	go app.wstatus.Update(ctx, ev)
}

func (app *Application) updateResize(ctx context.Context, ev tcell.Event, x, y int) {
	app.cancel()
	app.Screen.Fill(' ', tcell.Style{})
	app.Screen.Sync()
	app.wprogress.View().Resize(0, 0, x, 3)
	app.wlyrics.View().Resize(0, 3, x, y-6)
	app.wstatus.View().Resize(0, y-3, x, 3)
	data, err := app.client.Data()
	if err != nil {
		log.Println("updateResize:", err)
		return
	}
	app.updateData(ctx, ev, data)
}

// Quit the application.
func (app *Application) Quit() {
	app.quit()
}

// Run the application.
func (app *Application) Run(ctx context.Context) (err error) {
	app.Screen, err = tcell.NewScreen()
	if err != nil {
		return
	}

	err = app.Screen.Init()
	if err != nil {
		return
	}
	defer app.Screen.Fini()
	defer app.Quit()

	ctx = panics.ContextWithHook(ctx, app.Quit)
	ctx = events.ContextWith(ctx, app.events)
	ctx, app.quit = context.WithCancel(ctx)

	go app.Screen.ChannelEvents(app.events, ctx.Done())
	go app.client.PostEvents(ctx)

	app.wprogress.SetView(app.Screen)
	app.wlyrics.SetView(app.Screen)
	app.wstatus.SetView(app.Screen)

	for ev := range app.events {
		app.update(ctx, ev)
		app.Show()
	}
	return
}
