package client

import (
	"context"
	"os"
	"runtime"
	"strconv"
	"time"

	"github.com/gdamore/tcell/v2"
	"go.uber.org/atomic"
	"go.uber.org/multierr"

	"github.com/env25/mpdlrc/internal/events"
	"github.com/env25/mpdlrc/internal/panics"
)

type MPDClient struct {
	closed atomic.Bool

	c *mpd_Client
	w *mpd_Watcher

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
	c.c, err = mpd_DialAuthenticated(c.net, c.addr, c.password)
	if err != nil {
		return err
	}
	c.w, err = mpd_NewWatcher(c.net, c.addr, c.password)
	if err != nil {
		return err
	}
	runtime.SetFinalizer(c, func(c *MPDClient) { _ = c.Stop() })
	return
}

func (c *MPDClient) Stop() (err error) {
	if !c.closed.CAS(false, true) {
		return os.ErrClosed
	}
	multierr.AppendInvoke(&err, multierr.Invoke(c.c.Close))
	multierr.AppendInvoke(&err, multierr.Invoke(c.w.Close))
	return
}

func (c *MPDClient) Pause() error {
	if c.closed.Load() {
		return os.ErrClosed
	}
	return c.c.Pause(true)
}

func (c *MPDClient) Play() error {
	if c.closed.Load() {
		return os.ErrClosed
	}
	return c.c.Pause(false)
}

func (c *MPDClient) Ping() error {
	if c.closed.Load() {
		return os.ErrClosed
	}
	return c.c.Ping()
}

func (c *MPDClient) NowPlaying() (Song, error) {
	if c.closed.Load() {
		return nil, os.ErrClosed
	}
	attrs, err := c.c.CurrentSong()
	return MPDSong(attrs), err
}

func (c *MPDClient) Status() (Status, error) {
	if c.closed.Load() {
		return nil, os.ErrClosed
	}
	attrs, err := c.c.Status()
	return MPDStatus(attrs), err
}

func (c *MPDClient) MusicDir() (string, error) {
	if c.closed.Load() {
		return "", os.ErrClosed
	}
	attrs, err := c.c.Command("config").Attrs()
	return attrs["music_directory"], err
}

func (c *MPDClient) PostEvents(ctx context.Context) {
	defer panics.Handle(ctx)
	var newEvent func() tcell.Event
	for {
		select {
		case <-ctx.Done():
			return
		case <-c.w.Error:
			// no-op
		case mpdev := <-c.w.Event:
			switch mpdev {
			case "player":
				newEvent = newPlayerEvent
			case "options":
				newEvent = newOptionsEvent
			}
			if newEvent != nil {
				events.PostEvent(ctx, newEvent())
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
