package mpd

import (
	"strconv"
	"time"

	"local/mpdlrc/types"

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
	case types.PlayState:
		c.Pause()
	case types.PauseState:
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

func (c *MPDClient) NowPlaying() types.Song {
	if c.closed {
		return nil
	}
	if attrs, err := c.client.CurrentSong(); err != nil {
		panic(err)
	} else {
		return Song(attrs)
	}
}

func (c *MPDClient) State() types.State {
	if c.closed {
		return 0
	}
	if status, err := c.client.Status(); err != nil || status == nil {
		return 0
	} else {
		switch status["state"] {
		case "play":
			return types.PlayState
		case "stop":
			return types.StopState
		case "pause":
			return types.PauseState
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
