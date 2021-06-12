package mpd

import (
	"strconv"
	"time"

	"github.com/env25/mpdlrc/internal/song"
	"github.com/env25/mpdlrc/internal/state"

	"github.com/fhs/gompd/v2/mpd"
)

type MPDClient struct {
	client *mpd.Client
	net    string
	addr   string
	closed bool
}

func NewMPDClient(net string, addr string) *MPDClient {
	return &MPDClient{
		net:  net,
		addr: addr,
	}
}

func (c *MPDClient) Pause() {
	if c.closed {
		return
	}
	_ = c.client.Pause(true)
}

func (c *MPDClient) Play() {
	if c.closed {
		return
	}
	_ = c.client.Pause(false)
}

func (c *MPDClient) TogglePlay() {
	if c.closed {
		return
	}
	switch c.State() {
	case state.PlayState:
		c.Pause()
	case state.PauseState:
		c.Play()
	}
}

func (c *MPDClient) Start() {
	if client, err := mpd.Dial(c.net, c.addr); err != nil {
		panic(err)
	} else {
		c.client = client
	}
}

func (c *MPDClient) Stop() {
	c.closed = true
	c.client.Close()
}

func (c *MPDClient) NowPlaying() song.Song {
	if c.closed {
		return nil
	}
	if attrs, err := c.client.CurrentSong(); err != nil {
		panic(err)
	} else {
		return Song(attrs)
	}
}

func (c *MPDClient) State() state.State {
	if c.closed {
		return 0
	}
	if status, err := c.client.Status(); err != nil || status == nil {
		return 0
	} else {
		switch status["state"] {
		case "play":
			return state.PlayState
		case "stop":
			return state.StopState
		case "pause":
			return state.PauseState
		}
	}
	return 0
}

func (c *MPDClient) Elapsed() time.Duration {
	if c.closed {
		return 0
	}
	if status, err := c.client.Status(); err != nil || status == nil {
		return 0
	} else {
		elapsed, _ := strconv.ParseFloat(status["elapsed"], 64)
		return time.Duration(elapsed * float64(time.Second))
	}
}
