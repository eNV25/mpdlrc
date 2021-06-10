package main

import (
	"fmt"
	"os"

	"local/mpdlrc/lrc"
)

var fpath string

func init() {
	fpath = os.Args[1]
}

func main() {
	f, err := os.Open(fpath)
	if err != nil {
		panic(err)
	} else {
		defer f.Close()
	}

	l, err := lrc.NewParser(f).Parse()
	if err != nil {
		panic(err)
	}

	lines := l.Lines()
	times := l.Times()

	for i := range lines {
		fmt.Printf("%10v     %v\n", times[i], lines[i])
	}
}
