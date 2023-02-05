// Package mpd provides the client side interface to MPD (Music Player Daemon).
package mpd

import (
	_ "github.com/fhs/gompd/v2/mpd" // needed by bundle
)

// Idle sends waits for changes in subsystems and returns the ones that changed.
func (c *Client) Idle(subsystems ...string) ([]string, error) {
	return c.idle(subsystems...)
}

// NoIdle cancels the current call to [Client.Idle].
func (c *Client) NoIdle() error {
	return c.noIdle()
}
