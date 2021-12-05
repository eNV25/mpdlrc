package mpd

import (
	"errors"
	"sync/atomic"

	"github.com/fhs/gompd/v2/mpd"

	"github.com/env25/mpdlrc/internal/client"
	"github.com/env25/mpdlrc/internal/song"
	"github.com/env25/mpdlrc/internal/status"
)

var (
	ErrAlreadyClosed = errors.New("MPDClient: already closed")
)

type MPDClient struct {
	client              *mpd.Client
	net, addr, password string

	closedval uint32 // atomic
}

var _ client.Client = &MPDClient{}

// NewMPDClient returns a pointer to an instance of MPDClient.
// A password of "" can be used if there is no password.
func NewMPDClient(net, addr, password string) *MPDClient {
	return &MPDClient{
		net:      net,
		addr:     addr,
		password: password,
	}
}

func (c *MPDClient) closed() bool { return atomic.LoadUint32(&c.closedval) != 0 }

func (c *MPDClient) setClosed() bool { return atomic.CompareAndSwapUint32(&c.closedval, 0, 1) }

func (c *MPDClient) Start() (err error) {
	c.client, err = mpd.DialAuthenticated(c.net, c.addr, c.password)
	return
}

func (c *MPDClient) Pause() {
	if c.closed() {
		return
	}
	_ = c.client.Pause(true)
}

func (c *MPDClient) Play() {
	if c.closed() {
		return
	}
	_ = c.client.Pause(false)
}

func (c *MPDClient) Ping() {
	if c.closed() {
		return
	}
	_ = c.client.Ping()
}

func (c *MPDClient) Stop() error {
	if !c.setClosed() {
		return ErrAlreadyClosed
	}
	return c.client.Close()
}

func (c *MPDClient) NowPlaying() song.Song {
	if c.closed() {
		return nil
	}
	if attrs, err := c.client.CurrentSong(); err != nil || attrs == nil {
		return nil
	} else {
		return Song(attrs)
	}
}

func (c *MPDClient) Status() status.Status {
	if c.closed() {
		return nil
	}
	if attrs, err := c.client.Status(); err != nil || attrs == nil {
		return nil
	} else {
		return Status(attrs)
	}
}
