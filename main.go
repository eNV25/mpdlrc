// Package main
package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"os"
	"strings"

	"github.com/docopt/docopt-go"

	"github.com/env25/mpdlrc/internal"
	"github.com/env25/mpdlrc/internal/client"
	"github.com/env25/mpdlrc/internal/config"
	"github.com/env25/mpdlrc/internal/dirs"
)

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

func maine() (_ int, err error) {
	opts, err := docopt.ParseDoc(usage)
	if err != nil {
		return 1, err
	}

	cfg := config.DefaultConfig()

	err = cfg.FromFiles(opts["--config"].([]string))
	if err != nil {
		return 1, err
	}

	cfg.FromEnv(dirs.GetEnv)
	cfg.FromOpts(opts)

	conn, err := client.NewMPDClient(&cfg.MPD.Connection, &cfg.MPD.Address, &cfg.MPD.Password, &cfg.LyricsDir)
	if err != nil {
		return 1, err
	}
	defer conn.Close()

	cfg.FromClient(conn)

	cfg.Expand()

	if opts["--dump-config"].(bool) {
		fmt.Print(cfg)
		return
	}

	err = cfg.Assert()
	if err != nil {
		return 1, err
	}

	if config.Debug {
		fmt.Fprint(os.Stderr, "\n", cfg, "\n")
	}

	err = internal.NewApplication(cfg, conn).Run(context.Background())
	if err != nil {
		return 1, err
	}

	return
}

func main() {
	var logBuilder strings.Builder
	log.SetFlags(0)
	if config.Debug {
		slog.SetDefault(slog.New(slog.NewTextHandler(&logBuilder, &slog.HandlerOptions{Level: slog.LevelDebug})))
	} else {
		log.SetOutput(&logBuilder)
	}

	code, err := maine()
	if err != nil {
		slog.Error("main", err)
	}

	fmt.Fprint(os.Stderr, &logBuilder)
	os.Exit(code)
}
