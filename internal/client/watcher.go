package client

import "github.com/gdamore/tcell/v2"

type Watcher interface {
	Start() error
	Stop() error
	PostEvents(postEvent func(tcell.Event) error, quit <-chan struct{})
}
