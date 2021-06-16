package main

import (
	"encoding/json"
	"fmt"

	"github.com/env25/mpdlrc/internal/client"
	"github.com/env25/mpdlrc/internal/mpd"
)

func main() {
	var c client.Client = mpd.NewMPDClient("unix", "/run/user/1000/mpd/socket")
	defer c.Stop()
	c.Start()
	var ret []byte
	ret, _ = json.MarshalIndent(c.NowPlaying(), "", "  ")
	fmt.Printf("%s\n", ret)
	ret, _ = json.MarshalIndent(c.Status(), "", "  ")
	fmt.Printf("%s\n", ret)
}
