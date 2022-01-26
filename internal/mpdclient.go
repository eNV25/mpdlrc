package internal

import (
	"errors"
	"path"
	"strconv"
	"sync/atomic"
	"time"

	"github.com/fhs/gompd/v2/mpd"
	"github.com/gdamore/tcell/v2"
)

var (
	ErrAlreadyClosed = errors.New("MPDClient: already closed")
)

type MPDClient struct {
	client              *mpd.Client
	net, addr, password string

	_closed uintptr // atomic
}

var _ ClientInterface = &MPDClient{}

// NewMPDClient returns a pointer to an instance of MPDClient.
// A password of "" can be used if there is no password.
func NewMPDClient(net, addr, password string) *MPDClient {
	return &MPDClient{
		net:      net,
		addr:     addr,
		password: password,
	}
}

func (c *MPDClient) closed() bool { return atomic.LoadUintptr(&c._closed) != 0 }

func (c *MPDClient) setClosed() bool { return atomic.CompareAndSwapUintptr(&c._closed, 0, 1) }

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

func (c *MPDClient) NowPlaying() Song {
	if c.closed() {
		return nil
	}
	if attrs, err := c.client.CurrentSong(); err != nil || attrs == nil {
		return nil
	} else {
		return MPDSong(attrs)
	}
}

func (c *MPDClient) Status() Status {
	if c.closed() {
		return nil
	}
	if attrs, err := c.client.Status(); err != nil || attrs == nil {
		return nil
	} else {
		return MPDStatus(attrs)
	}
}

type MPDWatcher struct {
	watcher             *mpd.Watcher
	net, addr, password string
}

var _ WatcherInterface = &MPDWatcher{}

func NewMPDWatcher(net, addr, password string) *MPDWatcher {
	return &MPDWatcher{net: net, addr: addr, password: password}
}

func (w *MPDWatcher) Start() (err error) {
	w.watcher, err = mpd.NewWatcher(w.net, w.addr, w.password, "player")
	return
}

func (w *MPDWatcher) Stop() error { return w.watcher.Close() }

func (w *MPDWatcher) PostEvents(ch chan<- tcell.Event, quit <-chan struct{}) {
	var newEvent (func() tcell.Event)
	for {
		select {
		case <-quit:
			return
		case mpdev := <-w.watcher.Event:
			switch mpdev {
			case "player":
				newEvent = NewEventPlayer
			}
			if newEvent != nil {
				ch <- newEvent()
				newEvent = nil
			}
		}
	}
}

type MPDSong map[string]string

var _ SongInterface = MPDSong{}

func (s MPDSong) ID() string {
	return s["Id"]
}

func (s MPDSong) Title() string {
	return s["Title"]
}

func (s MPDSong) Artist() string {
	return s["Artist"]
}

func (s MPDSong) Album() string {
	return s["Album"]
}

func (s MPDSong) File() string {
	return s["file"]
}

func (s MPDSong) LRCFile() string {
	file := s.File()
	return file[:(len(file)-len(path.Ext(file)))] + ".lrc"
}

type MPDStatus map[string]string

var _ StatusInterface = MPDStatus{}

func (s MPDStatus) Duration() time.Duration {
	return secondStringToDuration(s["duration"])
}

func (s MPDStatus) Elapsed() time.Duration {
	return secondStringToDuration(s["elapsed"])
}

func secondStringToDuration(str string) time.Duration {
	parsed, _ := strconv.ParseFloat(str, 64)
	return time.Duration(parsed * float64(time.Second))
}

func (s MPDStatus) State() State {
	switch s["state"] {
	case "play":
		return StatePlay
	case "stop":
		return StateStop
	case "pause":
		return StatePause
	}
	return 0
}
