package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/env25/mpdlrc/lrc"
)

func main() {
	var fpath string

	flag.StringVar(&fpath, "file", os.Stdin.Name(), "select file")

	flag.Parse()

	var f *os.File
	var err error

	if fpath == "" && flag.Arg(0) != "" {
		fpath = flag.Arg(0)
	}

	if fpath == os.Stdin.Name() {
		f = os.Stdin
	} else {
		f, err = os.Open(fpath)
		if err != nil {
			panic(err)
		}
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
