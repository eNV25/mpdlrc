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

	times, lines, err := lrc.ParseReader(f)
	if err != nil {
		panic(err)
	}

	for i := range lines {
		fmt.Printf("%10v  %v\n", times[i], lines[i])
	}
}
