package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path"
	"strconv"

	"github.com/fhs/gompd/v2/mpd"
	"github.com/spf13/pflag"
)

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	pflag.Parse()

	c, err := mpd.DialAuthenticated("unix", path.Join("/run/user", strconv.Itoa(os.Geteuid()), "mpd/socket"), "")
	check(err)
	defer func() {
		err = c.Close()
		check(err)
	}()

	attrs, err := c.CurrentSong()
	ret, errr := json.MarshalIndent(attrs, "", "  ")
	fmt.Printf("%s %v %v\n", ret, err, errr)

	attrs, err = c.Status()
	ret, errr = json.MarshalIndent(attrs, "", "  ")
	fmt.Printf("%s %v %v\n", ret, err, errr)

	attrss, err := c.Command("listmounts").AttrsList("mount")
	ret, errr = json.MarshalIndent(attrss, "", "  ")
	fmt.Printf("%s %v %v\n", ret, err, errr)

	attrs, err = c.Command("config").Attrs()
	ret, errr = json.MarshalIndent(attrs, "", "  ")
	fmt.Printf("%s %v %v\n", ret, err, errr)
}
