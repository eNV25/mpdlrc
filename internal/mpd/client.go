package mpd

import (
	"github.com/env25/mpdlrc/internal/song"
	"github.com/env25/mpdlrc/internal/status"

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

func (c *MPDClient) Ping() {
	_ = c.client.Ping()
}

func (c *MPDClient) Start() {
	var err error
	for {
		c.client, err = mpd.Dial(c.net, c.addr)
		if err == nil {
			break
		}
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
		return nil
	} else {
		return Song(attrs)
	}
}

func (c *MPDClient) Status() status.Status {
	if c.closed {
		return nil
	}
	if status, err := c.client.Status(); err != nil || status == nil {
		return nil
	} else {
		return Status(status)
	}
}
