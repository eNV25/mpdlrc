// Package main
package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/env25/mpdlrc/lrc"
)

func main() {
	var fpath string

	flag.StringVar(&fpath, "file", "", "select file")

	flag.Parse()

	var f *os.File

	if fpath == "" && flag.Arg(0) != "" {
		fpath = flag.Arg(0)
	}

	if err := error(nil); fpath == "" {
		f = os.Stdin
	} else {
		f, err = os.Open(fpath)
		if err != nil {
			panic(err)
		} else {
			defer f.Close()
		}
	}

	times, lines, err := lrc.ParseReader(f)
	if err != nil {
		panic(err)
	}

	for i := range lines {
		fmt.Printf("%10v  %v\n", times[i], lines[i])
	}
}
