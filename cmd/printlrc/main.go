package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/env25/mpdlrc/lrc"
)

var fpath string

func init() {
	flag.StringVar(&fpath, "file", os.Args[1], "select file")
}

func main() {
	flag.Parse()

	f, err := os.Open(fpath)
	if err != nil {
		panic(err)
	} else {
		defer f.Close()
	}

	lrcs, err := lrc.NewParser(f).Parse()
	if err != nil {
		panic(err)
	}

	lines := lrcs.Lines()
	times := lrcs.Times()

	for i := range lines {
		fmt.Printf("%10v  %v\n", times[i], lines[i])
	}
}
