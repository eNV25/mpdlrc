package client

import "github.com/gdamore/tcell/v2"

type Watcher interface {
	Start() error
	Stop() error
	PostEvents(ch chan<- tcell.Event, quit <-chan struct{})
}
