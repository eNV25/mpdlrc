package client

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"sync"
	"time"

	_ "github.com/env25/gompd/v2/mpd" // for mpd_bundle.go
	"github.com/gdamore/tcell/v2"
	"go.uber.org/atomic"
	"go.uber.org/multierr"

	"github.com/env25/mpdlrc/internal/config"
	"github.com/env25/mpdlrc/internal/events"
	"github.com/env25/mpdlrc/internal/lyrics"
	"github.com/env25/mpdlrc/internal/panics"
	"github.com/env25/mpdlrc/internal/ufilepath"
)

type MPDClient struct {
	mu     sync.Mutex
	closed atomic.Bool
	idling atomic.Bool
	locked atomic.Bool
	cond   sync.Cond

	c   *mpd_Client
	cfg *config.Config

	id  atomic.Value // string
	lrc *lyrics.Lyrics
}

var _ Client = &MPDClient{}

// NewMPDClient returns a pointer to an instance of MPDClient.
// A password of "" can be used if there is no password.
func NewMPDClient(cfg *config.Config) (*MPDClient, error) {
	for _, cs := range &[...]struct{ net, addr string }{
		{"unix", filepath.Join(os.Getenv("XDG_RUNTIME_DIR"), "mpd", "socket")},
		{"unix", filepath.Join(string(filepath.Separator), "run", "mpd", "socket")},
		{cfg.MPD.Connection, cfg.MPD.Address},
	} {
		c, err := mpd_DialAuthenticated(cs.net, cs.addr, cfg.MPD.Password)
		if err != nil {
			continue
		}
		cfg.MPD.Connection = cs.net
		cfg.MPD.Address = cs.addr
		return newMPDClient(c, cfg), nil
	}
	return nil, fmt.Errorf("NewMPDClient: %w", os.ErrNotExist)
}

func newMPDClient(c *mpd_Client, cfg *config.Config) *MPDClient {
	ret := &MPDClient{
		c:   c,
		cfg: cfg,
	}
	ret.cond.L = &ret.mu
	return ret
}

func (c *MPDClient) Ping() (err error) {
	c.lock()
	defer c.unlock()
	goto normal
err:
	return fmt.Errorf("MPDClient: Ping: %w", err)
normal:
	if c.closed.Load() {
		err = os.ErrClosed
		goto err
	}
	err = c.c.Ping()
	if err != nil {
		goto err
	}
	return nil
}

func (c *MPDClient) MusicDir() (_ string, err error) {
	c.lock()
	defer c.unlock()
	goto normal
err:
	err = fmt.Errorf("MPDClient: MusicDir: %w", err)
	return
normal:
	if c.closed.Load() {
		err = os.ErrClosed
		goto err
	}
	attrs, err := c.c.Command("config").Attrs()
	if err != nil {
		goto err
	}
	return attrs["music_directory"], nil
}

func (c *MPDClient) Data() (data Data, err error) {
	c.lock()
	defer c.unlock()
	goto start
err:
	err = fmt.Errorf("MPDClient: Data: %w", err)
	return
start:
	if c.closed.Load() {
		err = os.ErrClosed
		goto err
	}

	cmdlist := c.c.BeginCommandList()
	songFuture := cmdlist.CurrentSong()
	statusFuture := cmdlist.Status()

	err = cmdlist.End()
	if err != nil {
		goto err
	}

	song, songErr := songFuture.Value()
	status, statusErr := statusFuture.Value()
	err = multierr.Append(songErr, statusErr)
	if err != nil {
		goto err
	}

	{
		song := MPDSong(song)
		status := MPDStatus(status)
		return Data{
			Song:   song,
			Status: status,
			Lyrics: c.lyrics(song),
		}, nil
	}
}

func (c *MPDClient) lyrics(song Song) *lyrics.Lyrics {
	id := song.ID()
	old := c.id.Swap(id)
	if id != old {
		c.lrc = lyrics.ForFile(filepath.Join(c.cfg.LyricsDir, ufilepath.FromSlash(song.File())))
	}
	return c.lrc
}

func (c *MPDClient) idle() ([]string, error) {
	c.idling.Store(true)
	mpdevs, err := c.c.idle()
	c.idling.Store(false)
	return mpdevs, err
}

func (c *MPDClient) noIdle() error {
	if !c.idling.Load() {
		return nil
	}
	return c.c.noIdle()
}

func (c *MPDClient) lock() {
	if !c.locked.CompareAndSwap(false, true) {
		return
	}
	_ = c.noIdle()
	c.mu.Lock()
}

func (c *MPDClient) unlock() {
	if !c.locked.CompareAndSwap(true, false) {
		return
	}
	c.cond.Broadcast()
	c.mu.Unlock()
}

func (c *MPDClient) PostEvents(ctx context.Context) {
	defer panics.Handle(ctx)

	for {
		c.mu.Lock()

		mpdevs, err := c.idle()

		for c.locked.Load() && !c.closed.Load() {
			c.cond.Wait()
		}

		c.mu.Unlock()

		if c.closed.Load() {
			return
		}
		if err != nil {
			continue
		}
		select {
		case <-ctx.Done():
			return
		default:
		}

		for _, mpdev := range mpdevs {
			var newEvent func(Data) tcell.Event

			switch mpdev {
			case "player":
				newEvent = newPlayerEvent
			case "options":
				newEvent = newOptionsEvent
			default:
				continue
			}

			data, err := c.Data()
			if err != nil {
				if errors.Is(err, os.ErrClosed) {
					return
				}
				continue
			}

			if !events.PostEvent(ctx, newEvent(data)) {
				return
			}
		}
	}
}

func (c *MPDClient) Close() (err error) {
	if !c.closed.CompareAndSwap(false, true) {
		err = os.ErrClosed
		return
	}
	err = c.noIdle()
	c.mu.Lock()
	c.cond.Broadcast()
	defer c.mu.Unlock()
	goto normal
err:
	return fmt.Errorf("MPDClient: Close: %w", err)
normal:
	err = multierr.Append(err, c.c.Close())
	if err != nil {
		goto err
	}
	return nil
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
func (s MPDStatus) Duration() time.Duration { return s.timeDuration("duration", time.Second) }
func (s MPDStatus) Elapsed() time.Duration  { return s.timeDuration("elapsed", time.Second) }
func (s MPDStatus) Repeat() bool            { return s["repeat"] != "0" }
func (s MPDStatus) Random() bool            { return s["random"] != "0" }
func (s MPDStatus) Single() bool            { return s["single"] != "0" }
func (s MPDStatus) Consume() bool           { return s["consume"] != "0" }

func (s MPDStatus) timeDuration(key string, unit time.Duration) time.Duration {
	parsed, _ := strconv.ParseFloat(s[key], 64)
	return time.Duration(parsed * float64(unit))
}
