package client

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"sync"
	"sync/atomic"
	"time"

	"github.com/gdamore/tcell/v2"
	"go.uber.org/multierr"

	"github.com/env25/mpdlrc/internal/config"
	"github.com/env25/mpdlrc/internal/events"
	"github.com/env25/mpdlrc/internal/lyrics"
	"github.com/env25/mpdlrc/internal/mpd"
	"github.com/env25/mpdlrc/internal/panics"
)

// MPDClient implents a [Client] for the Music Player Daemon (MPD).
type MPDClient struct {
	mu     sync.Mutex
	closed atomic.Bool
	idling atomic.Bool
	locked atomic.Bool
	cond   sync.Cond

	c   *mpd.Client
	cfg *config.Config

	id  atomic.Value // string
	lrc *lyrics.Lyrics
}

var _ Client = &MPDClient{}

// NewMPDClient returns a pointer to an instance of [MPDClient].
// A password of "" can be used if there is no password.
func NewMPDClient(cfg *config.Config) (*MPDClient, error) {
	for _, cs := range &[...]struct{ net, addr string }{
		{cfg.MPD.Connection, cfg.MPD.Address},
		{"unix", filepath.Join(config.GetEnv("XDG_RUNTIME_DIR"), "mpd", "socket")},
		{"unix", filepath.Join(config.RootDir(), "run", "mpd", "socket")},
		{"tcp", ":6600"},
	} {
		if cs.net == "" || cs.addr == "" {
			continue
		}
		c, err := mpd.DialAuthenticated(cs.net, cs.addr, cfg.MPD.Password)
		if err != nil {
			continue
		}
		cfg.MPD.Connection = cs.net
		cfg.MPD.Address = cs.addr
		return newMPDClient(c, cfg), nil
	}
	return nil, fmt.Errorf("NewMPDClient: client not found: %w", os.ErrNotExist)
}

func newMPDClient(c *mpd.Client, cfg *config.Config) *MPDClient {
	ret := &MPDClient{
		c:   c,
		cfg: cfg,
	}
	ret.cond.L = &ret.mu
	return ret
}

// Close finalilly closes the connection.
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
	multierr.AppendInto(&err, c.c.Close())
	if err != nil {
		goto err
	}
	return nil
}

func (c *MPDClient) idle() ([]string, error) {
	c.idling.Store(true)
	mpdevs, err := c.c.Idle()
	c.idling.Store(false)
	return mpdevs, err
}

func (c *MPDClient) noIdle() error {
	if !c.idling.Load() {
		return nil
	}
	return c.c.NoIdle()
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
	c.mu.Unlock()
	c.cond.Broadcast()
}

// Data returns [Data] for the currectly playing song.
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

// MusicDir return the music directory, if available.
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

func (c *MPDClient) lyrics(song Song) *lyrics.Lyrics {
	id := song.ID()
	old := c.id.Swap(id)
	if id != old {
		c.lrc = lyrics.ForFile(filepath.Join(c.cfg.LyricsDir, filepath.FromSlash(song.File())))
	}
	return c.lrc
}

// Ping sends no-op message.
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

// TogglePause toggles the pause state.
func (c *MPDClient) TogglePause() bool {
	c.lock()
	defer c.unlock()
	status, _ := c.c.Status()
	pause := MPDStatus(status).State() != "pause"
	_ = c.c.Pause(pause)
	return pause
}

// PostEvents listens for events, and forwards them to [events.PostEvent].
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

// MPDSong implements [Song].
type MPDSong mpd.Attrs

var _ Song = MPDSong{}

// ID returns the song id.
func (s MPDSong) ID() string { return s["Id"] }

// Title returns the song title.
func (s MPDSong) Title() string { return s["Title"] }

// Artist returns the song artist.
func (s MPDSong) Artist() string { return s["Artist"] }

// Album returns the song Album.
func (s MPDSong) Album() string { return s["Album"] }

// Date returns the song Date.
func (s MPDSong) Date() string { return s["Date"] }

// File returns the song File.
func (s MPDSong) File() string { return s["file"] }

// MPDStatus implements [Status].
type MPDStatus map[string]string

var _ Status = MPDStatus{}

// State returns the player state. "pause", "play"
func (s MPDStatus) State() string { return s["state"] }

// Duration returns the player song duration.
func (s MPDStatus) Duration() time.Duration { return s.timeDuration("duration", time.Second) }

// Elapsed returns the pleyer elapsed duration.
func (s MPDStatus) Elapsed() time.Duration { return s.timeDuration("elapsed", time.Second) }

// Repeat returns repeat option
func (s MPDStatus) Repeat() bool { return s["repeat"] != "0" }

// Random returns random option.
func (s MPDStatus) Random() bool { return s["random"] != "0" }

// Single returns single option.
func (s MPDStatus) Single() bool { return s["single"] != "0" }

// Consume returns consume option.
func (s MPDStatus) Consume() bool { return s["consume"] != "0" }

func (s MPDStatus) timeDuration(key string, unit time.Duration) time.Duration {
	parsed, _ := strconv.ParseFloat(s[key], 64)
	return time.Duration(parsed * float64(unit))
}
