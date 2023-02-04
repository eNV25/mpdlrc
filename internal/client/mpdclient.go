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

	"github.com/env25/mpdlrc/internal/dirs"
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
	wait   atomic.Bool
	nowait sync.Cond

	c *mpd.Client

	id        atomic.Value // string
	lrc       *lyrics.Lyrics
	lyricsDir *string
}

var _ Client = &MPDClient{}

// NewMPDClient returns a pointer to an instance of [MPDClient].
// A password of "" can be used if there is no password.
func NewMPDClient(conntype, addr, passwd, lyricsdir *string) (*MPDClient, error) {
	for _, cs := range &[...]struct{ net, addr string }{
		{*conntype, *addr},
		{"unix", filepath.Join(dirs.GetEnv("XDG_RUNTIME_DIR"), "mpd", "socket")},
		{"unix", filepath.Join(dirs.RootDir(), "run", "mpd", "socket")},
		{"tcp", ":6600"},
	} {
		if cs.net == "" || cs.addr == "" {
			continue
		}
		c, err := mpd.DialAuthenticated(cs.net, cs.addr, *passwd)
		if err != nil {
			continue
		}
		*conntype = cs.net
		*addr = cs.addr
		return newMPDClient(c, lyricsdir), nil
	}
	return nil, fmt.Errorf("MPD client not found: %w", os.ErrNotExist)
}

func newMPDClient(c *mpd.Client, lyricsdir *string) *MPDClient {
	ret := &MPDClient{
		c:         c,
		lyricsDir: lyricsdir,
	}
	ret.nowait.L = &ret.mu
	return ret
}

// Close finalilly closes the connection.
func (c *MPDClient) Close() error {
	if !c.closed.CompareAndSwap(false, true) {
		return os.ErrClosed
	}

	c.wait.Store(false)
	c.noIdle()
	c.nowait.Broadcast()
	c.mu.Lock()
	defer c.mu.Unlock()

	err := c.c.Close()
	if err != nil {
		return err
	}
	return nil
}

func (c *MPDClient) idle() ([]string, error) {
	c.idling.Store(true)
	mpdevs, err := c.c.Idle()
	c.idling.Store(false)
	return mpdevs, err
}

func (c *MPDClient) noIdle() {
	for c.idling.Load() && c.c.NoIdle() != nil {
	}
}

func (c *MPDClient) lock() {
	if !c.wait.CompareAndSwap(false, true) {
		return
	}
	c.noIdle()
	c.mu.Lock()
}

func (c *MPDClient) unlock() {
	if !c.wait.CompareAndSwap(true, false) {
		return
	}
	c.mu.Unlock()
	c.nowait.Broadcast()
}

// Data returns [Data] for the currectly playing song.
func (c *MPDClient) Data() (Data, error) {
	if c.closed.Load() {
		return Data{}, os.ErrClosed
	}
	c.lock()
	defer c.unlock()

	cmdlist := c.c.BeginCommandList()
	songFuture := cmdlist.CurrentSong()
	statusFuture := cmdlist.Status()

	err := cmdlist.End()
	if err != nil {
		return Data{}, err
	}

	song, songErr := songFuture.Value()
	status, statusErr := statusFuture.Value()
	if err = errors.Join(songErr, statusErr); err != nil {
		return Data{}, err
	}

	return Data{
		Song:   MPDSong(song),
		Status: MPDStatus(status),
		Lyrics: c.lyrics(MPDSong(song)),
	}, nil
}

// MusicDir return the music directory, if available.
func (c *MPDClient) MusicDir() (string, error) {
	if c.closed.Load() {
		return "", os.ErrClosed
	}
	c.lock()
	defer c.unlock()
	attrs, err := c.c.Command("config").Attrs()
	if err != nil {
		return "", err
	}
	return attrs["music_directory"], nil
}

func (c *MPDClient) lyrics(song Song) *lyrics.Lyrics {
	id := song.ID()
	old := c.id.Swap(id)
	if id != old {
		c.lrc = lyrics.ForFile(filepath.Join(*c.lyricsDir, filepath.FromSlash(song.File())))
	}
	return c.lrc
}

// Ping sends no-op message.
func (c *MPDClient) Ping() error {
	if c.closed.Load() {
		return os.ErrClosed
	}
	c.lock()
	defer c.unlock()
	err := c.c.Ping()
	if err != nil {
		return err
	}
	return nil
}

// TogglePause toggles the pause state.
func (c *MPDClient) TogglePause() bool {
	if c.closed.Load() {
		return false
	}
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

		for c.wait.Load() {
			c.nowait.Wait() // blocks
		}

		if c.closed.Load() {
			return
		}

		mpdevs, err := c.idle() // blocks

		c.mu.Unlock()

		if c.closed.Load() {
			return
		}

		if err != nil {
			continue
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
				continue
			}

			if !events.PostEvent(ctx, newEvent(data)) { // blocks
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

// State returns the player state. "pause", "play".
func (s MPDStatus) State() string { return s["state"] }

// Duration returns the player song duration.
func (s MPDStatus) Duration() time.Duration { return s.timeDuration("duration", time.Second) }

// Elapsed returns the pleyer elapsed duration.
func (s MPDStatus) Elapsed() time.Duration { return s.timeDuration("elapsed", time.Second) }

// Repeat returns repeat option.
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
