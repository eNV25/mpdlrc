package main

import (
	"log"
	"os"

	mpdlrc "github.com/env25/mpdlrc/internal"
	"github.com/env25/mpdlrc/internal/config"

	"github.com/spf13/pflag"
)

var (
	usage = false
	cfg   = config.DefaultConfig()
)

func init() {
	pflag.StringVar(&cfg.MusicDir, `musicdir`, cfg.MusicDir, `override MusicDir`)
	pflag.StringVar(&cfg.LyricsDir, `lyricsdir`, cfg.LyricsDir, `override LyricsDir`)
	pflag.StringVar(&cfg.MPD.Protocol, `mpd.protocol`, cfg.MPD.Protocol, `override MPD.Protocol (possible "unix", "tcp")`)
	pflag.StringVar(&cfg.MPD.Address, `mpd.address`, cfg.MPD.Address, `override MPD.Address (use unix socket path or "host:port")`)
	pflag.BoolVar(&cfg.Debug, `debug`, cfg.Debug, `enable debug`)
	pflag.BoolVarP(&usage, `help`, `h`, usage, `show this help message`)
}

func main() {
	for _, f := range config.ConfigFiles {
		err := cfg.MergeTOMLFile(f)
		if err != nil {
			switch err.(type) {
			case *os.PathError:
				// no-op
			default:
				log.Fatalln(err)
			}
		}
	}

	pflag.Parse()

	if usage {
		pflag.Usage()
		os.Exit(0)
	}

	cfg.Expand()

	if err := cfg.Assert(); err != nil {
		log.Fatalln(err)
	}

	if err := mpdlrc.NewApplication(cfg).Run(); err != nil {
		log.Fatalln(err)
	}

	os.Exit(0)
}
