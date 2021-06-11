package main

import (
	"flag"
	"fmt"
	"os"

	"local/mpdlrc/lrc"
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

	l, err := lrc.NewParser(f).Parse()
	if err != nil {
		panic(err)
	}

	lines := l.Lines()
	times := l.Times()

	for i := range lines {
		fmt.Printf("%10v  %v\n", times[i], lines[i])
	}
}
