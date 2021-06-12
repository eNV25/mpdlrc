package client

import "github.com/gdamore/tcell/v2"

type Watcher interface {
	PostEvents(
		postEvent func(tcell.Event) error,
		quit <-chan struct{},
	)
}
