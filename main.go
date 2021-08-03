package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"github.com/pelletier/go-toml/v2"
	"github.com/spf13/pflag"

	"github.com/env25/mpdlrc/internal"
	"github.com/env25/mpdlrc/internal/config"
)

var exitCode = 0

func exit() { os.Exit(exitCode) }

func main() {
	defer exit()

	var (
		usage   = false
		dumpcfg = false
		cfg     = config.DefaultConfig()
	)

	pflag.StringVar(&cfg.MusicDir, `musicdir`, cfg.MusicDir, `override MusicDir`)
	pflag.StringVar(&cfg.LyricsDir, `lyricsdir`, cfg.LyricsDir, `override LyricsDir`)
	pflag.StringVar(&cfg.MPD.Connection, `mpd-connection`, cfg.MPD.Connection, `override MPD.Connection (possible "unix", "tcp")`)
	pflag.StringVar(&cfg.MPD.Address, `mpd-address`, cfg.MPD.Address, `override MPD.Address (use unix socket path or "host:port")`)
	pflag.StringVar(&cfg.MPD.Password, `mpd-password`, cfg.MPD.Password, `override MPD.Password`)
	pflag.BoolVar(&cfg.Debug, `debug`, cfg.Debug, `enable debug`)
	pflag.BoolVar(&dumpcfg, `dump-config`, dumpcfg, `dump config`)
	pflag.BoolVarP(&usage, `help`, `h`, usage, `show this help message`)

	for _, fpath := range config.ConfigFiles {
		f, err := os.Open(fpath)
		if err != nil {
			switch err.(type) {
			case *os.PathError:
				// don't do anything if file doesn't exist //
			default:
				fmt.Fprintln(os.Stderr, fmt.Errorf("open config file: %w", err))
			}
		}
		err = toml.NewDecoder(f).Decode(cfg)
		if err != nil {
			fmt.Fprintln(os.Stderr, fmt.Errorf("decode config file: %w", err))
		}
	}

	pflag.Parse()

	if usage {
		pflag.Usage()
		return
	}

	if dumpcfg {
		cfg.Expand()
		var b strings.Builder
		toml.NewEncoder(&b).Encode(cfg)
		fmt.Fprint(os.Stdout, b.String()[:b.Len()-1])
		return
	}

	log.SetFlags(0)

	if cfg.Debug {
		var logBuilder strings.Builder
		log.SetOutput(&logBuilder)
		defer fmt.Fprint(os.Stderr, logBuilder.String())
	} else {
		log.SetOutput(io.Discard)
	}

	cfg.Expand()

	if err := cfg.Assert(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		exitCode = 1
		return
	}

	if err := internal.NewApplication(cfg).Run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		exitCode = 1
		return
	}
}
