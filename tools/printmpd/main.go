package main

import (
	"encoding/json"
	"flag"
	"fmt"

	"github.com/fhs/gompd/v2/mpd"
)

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	flag.Parse()

	cnxn := "unix"
	addr := "/run/user/1000/mpd/socket"

	m, err := mpd.DialAuthenticated(cnxn, addr, "")
	check(err)
	defer m.Close()

	a, err := m.CurrentSong()
	ret, err := json.MarshalIndent(a, "", "  ")
	fmt.Printf("%s // %v\n", ret, err)

	a, err = m.Status()
	ret, err = json.MarshalIndent(a, "", "  ")
	fmt.Printf("%s // %v\n", ret, err)

	as, err := m.Command("listmounts").AttrsList("mount")
	ret, err = json.MarshalIndent(as, "", "  ")
	fmt.Printf("%s // %v\n", ret, err)

	a, err = m.Command("config").Attrs()
	ret, err = json.MarshalIndent(a, "", "  ")
	fmt.Printf("%s // %v\n", ret, err)
}
