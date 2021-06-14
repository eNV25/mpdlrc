package main

import (
	"fmt"

	"github.com/env25/mpdlrc/internal/client"
	"github.com/env25/mpdlrc/internal/mpd"
)

func main() {
	var c client.Client = mpd.NewMPDClient("unix", "/run/user/1000/mpd/socket")
	defer c.Stop()
	c.Start()
	fmt.Println(c.NowPlaying())
	fmt.Println(c.Elapsed())
}
