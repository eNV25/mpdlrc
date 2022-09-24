package main

import (
	"context"
	"log"
	"os"
	"strings"

	"github.com/docopt/docopt-go"

	"github.com/env25/mpdlrc/internal"
	"github.com/env25/mpdlrc/internal/client"
	"github.com/env25/mpdlrc/internal/config"
)

func init() {
	log.SetFlags(0)
	log.SetPrefix("mpdlrc: ")
}

func main() {
	var err error
	var exitCode int

	defer func() {
		if err != nil {
			log.Println(err)
		}
		os.Exit(exitCode)
	}()

	const usage = `
Display synchronized lyrics for track playing in MPD.

Usage:
	mpdlrc [options] [--config=FILE]...

Options:
	--config=FILE           Use config file
	--dump-config           Print final config
	-h, --help              Show this help and exit

Configuration Options:
	--lyricsdir=DIR         override cfg.LyricsDir
	--musicdir=DIR          override cfg.MusicDir
	--mpd-address=ADDR      override cfg.MPD.Address
	--mpd-connection=CONN   override cfg.MPD.Connection
	--mpd-password=PASSWD   override cfg.MPD.Password
`

	opts, err := docopt.ParseDoc(usage)
	if err != nil {
		log.Println("docopt parse:", err)
		exitCode = 1
		return
	}

	cfg := config.DefaultConfig()

	err = cfg.FromFiles(opts["--config"].([]string))
	if err != nil {
		log.Println(err)
		return
	}

	cfg.FromEnv(config.GetEnv)
	cfg.FromOpts(opts)

	conn, err := client.NewMPDClient(cfg)
	if err != nil {
		return
	}
	defer conn.Close()

	cfg.FromClient(conn.MusicDir())
	cfg.Expand()

	if opts["--dump-config"].(bool) {
		log.Print(cfg)
		return
	}

	if config.Debug {
		log.Print("\n\n", cfg, "\n")
	}

	err = cfg.Assert()
	if err != nil {
		log.Println(err)
		exitCode = 1
		return
	}

	var logBuilder strings.Builder
	defer log.Println(&logBuilder)
	defer log.SetOutput(log.Writer())
	log.SetOutput(&logBuilder)

	err = internal.NewApplication(cfg, conn).Run(context.Background())
}
