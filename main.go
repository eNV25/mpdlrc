// Package main
package main

import (
	"context"
	"fmt"
	"io"
	stdlog "log"
	"os"
	"strings"

	"github.com/docopt/docopt-go"

	"github.com/env25/mpdlrc/internal"
	"github.com/env25/mpdlrc/internal/client"
	"github.com/env25/mpdlrc/internal/config"
	"github.com/env25/mpdlrc/internal/dirs"
	"github.com/env25/mpdlrc/internal/zerolog"
	"github.com/env25/mpdlrc/internal/zerolog/log"
)

func main() {
	os.Exit(maine())
}

const usage = `
Display synchronized lyrics for track playing in MPD.

Usage:
	mpdlrc -h|--help
	mpdlrc [options]
	mpdlrc [options] --config=FILE...

Options:
	-h, --help              Show this help and exit
	--config=FILE           Use config file
	--dump-config           Print final config

Configuration Options:
	--lyricsdir=DIR         override cfg.LyricsDir
	--musicdir=DIR          override cfg.MusicDir
	--mpd-connection=CONN   override cfg.MPD.Connection
	--mpd-address=ADDR      override cfg.MPD.Address
	--mpd-password=PASSWD   override cfg.MPD.Password
`

func init() {
	stdlog.SetFlags(0)

	if config.Debug {
		var logBuilder strings.Builder
		defer fmt.Fprint(os.Stderr, &logBuilder)
		log.Logger = zerolog.New(&zerolog.ConsoleWriter{Out: &logBuilder}).With().Timestamp().Logger()

		zerolog.SetGlobalLevel(zerolog.DebugLevel)
		stdlog.SetOutput(&log.Logger)
	} else {
		zerolog.SetGlobalLevel(zerolog.Disabled)
		stdlog.SetOutput(io.Discard)
	}
}

func maine() int {
	ctx := context.Background()

	opts, err := docopt.ParseDoc(usage)
	if err != nil {
		log.Err(err).Send()
		return 1
	}

	cfg := config.DefaultConfig()

	err = cfg.FromFiles(opts["--config"].([]string))
	if err != nil {
		log.Err(err).Send()
		return 1
	}

	cfg.FromEnv(dirs.GetEnv)
	cfg.FromOpts(opts)

	conn, err := client.NewMPDClient(&cfg.MPD.Connection, &cfg.MPD.Address, &cfg.MPD.Password, &cfg.LyricsDir)
	if err != nil {
		log.Err(err).Send()
		return 1
	}
	defer conn.Close()

	cfg.FromClient(conn)

	cfg.Expand()

	if opts["--dump-config"].(bool) {
		fmt.Print(cfg)
		return 0
	}

	err = cfg.Assert()
	if err != nil {
		log.Err(err).Send()
		return 1
	}

	if config.Debug {
		fmt.Fprint(os.Stderr, "\n", cfg, "\n")
	}

	err = internal.NewApplication(cfg, conn).Run(ctx)
	if err != nil {
		log.Err(err).Send()
		return 1
	}

	return 0
}
