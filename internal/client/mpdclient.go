package client

import (
	"context"
	"os"
	"runtime"
	"strconv"
	"time"

	"github.com/fhs/gompd/v2/mpd"
	"github.com/gdamore/tcell/v2"
	"go.uber.org/atomic"

	"github.com/env25/mpdlrc/internal/event"
	"github.com/env25/mpdlrc/internal/events"
)

type MPDClient struct {
	closed atomic.Bool

	client              *mpd.Client
	net, addr, password string
}

var _ Client = &MPDClient{}

// NewMPDClient returns a pointer to an instance of MPDClient.
// A password of "" can be used if there is no password.
func NewMPDClient(net, addr, password string) *MPDClient {
	return &MPDClient{
		net:      net,
		addr:     addr,
		password: password,
	}
}

func (c *MPDClient) Start() (err error) {
	c.client, err = mpd.DialAuthenticated(c.net, c.addr, c.password)
	runtime.SetFinalizer(c, func(c *MPDClient) { _ = c.Stop() })
	return
}

func (c *MPDClient) Pause() error {
	if c.closed.Load() {
		return os.ErrClosed
	}
	return c.client.Pause(true)
}

func (c *MPDClient) Play() error {
	if c.closed.Load() {
		return os.ErrClosed
	}
	return c.client.Pause(false)
}

func (c *MPDClient) Ping() error {
	if c.closed.Load() {
		return os.ErrClosed
	}
	return c.client.Ping()
}

func (c *MPDClient) Stop() error {
	if !c.closed.CAS(false, true) {
		return os.ErrClosed
	}
	return c.client.Close()
}

func (c *MPDClient) NowPlaying() (Song, error) {
	if c.closed.Load() {
		return nil, os.ErrClosed
	}
	attrs, err := c.client.CurrentSong()
	return MPDSong(attrs), err
}

func (c *MPDClient) Status() (Status, error) {
	if c.closed.Load() {
		return nil, os.ErrClosed
	}
	attrs, err := c.client.Status()
	return MPDStatus(attrs), err
}

func (c *MPDClient) MusicDir() (string, error) {
	if c.closed.Load() {
		return "", os.ErrClosed
	}
	attrs, err := c.client.Command("config").Attrs()
	return attrs["music_directory"], err
}

type MPDWatcher struct {
	watcher             *mpd.Watcher
	net, addr, password string
}

var _ Watcher = &MPDWatcher{}

func NewMPDWatcher(net, addr, password string) *MPDWatcher {
	return &MPDWatcher{net: net, addr: addr, password: password}
}

func (w *MPDWatcher) Start() (err error) {
	w.watcher, err = mpd.NewWatcher(w.net, w.addr, w.password, "player")
	runtime.SetFinalizer(w, func(w *MPDWatcher) { _ = w.Stop() })
	return
}

func (w *MPDWatcher) Stop() error { return w.watcher.Close() }

func (w *MPDWatcher) PostEvents(ctx context.Context) {
	ch := events.FromContext(ctx)
	var newEvent (func() tcell.Event)
	for {
		select {
		case <-ctx.Done():
			return
		case mpdev := <-w.watcher.Event:
			switch mpdev {
			case "player":
				newEvent = event.NewPlayer
			}
			if newEvent != nil {
				select {
				case <-ctx.Done():
					return
				case ch <- newEvent():
				}
				newEvent = nil
			}
		}
	}
}

type MPDSong map[string]string

var _ Song = MPDSong{}

func (s MPDSong) ID() string     { return s["Id"] }
func (s MPDSong) Title() string  { return s["Title"] }
func (s MPDSong) Artist() string { return s["Artist"] }
func (s MPDSong) Album() string  { return s["Album"] }
func (s MPDSong) Date() string   { return s["Date"] }
func (s MPDSong) File() string   { return s["file"] }

type MPDStatus map[string]string

var _ Status = MPDStatus{}

func (s MPDStatus) State() string           { return s["state"] }
func (s MPDStatus) Duration() time.Duration { return s.timeDuration("duration") }
func (s MPDStatus) Elapsed() time.Duration  { return s.timeDuration("elapsed") }
func (s MPDStatus) Repeat() bool            { return s["repeat"] != "0" }
func (s MPDStatus) Random() bool            { return s["random"] != "0" }
func (s MPDStatus) Single() bool            { return s["single"] != "0" }
func (s MPDStatus) Consume() bool           { return s["consume"] != "0" }

func (s MPDStatus) timeDuration(key string) time.Duration {
	parsed, _ := strconv.ParseFloat(s[key], 64)
	return time.Duration(parsed * float64(time.Second))
}
