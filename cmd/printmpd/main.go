package main

import (
	"fmt"
	"reflect"
	"runtime"

	"local/mpdlrc/client"
	"local/mpdlrc/mpd"
)

func FunctionName(i interface{}) string {
	return runtime.FuncForPC(reflect.ValueOf(i).Pointer()).Name()
}

func handle(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	var c client.Client = mpd.NewMPDClient("unix", "/run/user/1000/mpd/socket")
	defer c.Stop()
	c.Start()
	fmt.Println(c.NowPlaying())
	fmt.Println(c.Elapsed())
}
