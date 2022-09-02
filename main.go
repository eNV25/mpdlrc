package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/docopt/docopt-go"
	"github.com/pelletier/go-toml/v2"
	"go.uber.org/multierr"

	"github.com/env25/mpdlrc/internal"
	"github.com/env25/mpdlrc/internal/config"
)

func init() {
	log.SetFlags(0)
}

func main() {
	exitCode := 0

	defer func() {
		os.Exit(exitCode)
	}()

	const usage = `
Display MPD synchronized lyrics.

Usage:
    mpdlrc [options] [--config=FILE]...

Options:
    --config=FILE           Use config file
    --dump-config           Print final config

Configuration Options:
    --lyricsdir=DIR         override cfg.LyricsDir
    --musicdir=DIR          override cfg.MusicDir
    --mpd-address=ADDR      override cfg.MPD.Address
    --mpd-connection=CONN   override cfg.MPD.Connection
    --mpd-password=PASSWD   override cfg.MPD.Password
`

	opts, err := docopt.ParseDoc(usage)
	if err != nil {
		fmt.Println("docopt parse:", err)
		exitCode = 1
		return
	}

	cfg := config.DefaultConfig()

	for _, fpath := range opts["--config"].([]string) {
		var err error
		func() {
			var f *os.File
			f, err = os.Open(fpath)
			if err != nil {
				return
			}
			defer multierr.AppendInvoke(&err, multierr.Invoke(f.Close))
			err = toml.NewDecoder(f).Decode(cfg)
			if err != nil {
				return
			}
		}()
		if err != nil && !errors.Is(err, os.ErrNotExist) {
			log.Println("config file:", err)
		}
	}

	cfg.FromOpts(opts)
	cfg.Expand()

	var logBuilder strings.Builder
	defer fmt.Fprint(os.Stderr, &logBuilder)

	if opts["--dump-config"].(bool) {
		_ = toml.NewEncoder(&logBuilder).Encode(cfg)
		return
	}

	log.SetOutput(&logBuilder)

	err = internal.NewApplication(cfg).Run()
	if err != nil {
		log.Println(err)
		exitCode = 1
		return
	}
}
