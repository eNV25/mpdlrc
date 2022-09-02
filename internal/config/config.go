package config

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/docopt/docopt-go"
	"github.com/pelletier/go-toml/v2"
)

var _ fmt.Stringer = (*Config)(nil)

type Config struct {
	MusicDir  string
	LyricsDir string

	MPD struct {
		Connection string
		Address    string
		Password   string
	}
}

func DefaultConfig() (cfg *Config) {
	cfg = &Config{}
	cfg.MusicDir = "~/Music"
	cfg.LyricsDir = ""
	host := os.Getenv("MPD_HOST")
	port := os.Getenv("MPD_PORT")
	if host == "" {
		if port == "" {
			cfg.MPD.Connection = "tcp"
			cfg.MPD.Address = "localhost:6600"
			cfg.MPD.Password = ""
		} else {
			cfg.MPD.Connection = "tcp"
			cfg.MPD.Address = "localhost:" + port
			cfg.MPD.Password = ""
		}
	} else {
		// If @ is found in host, it has two possible meanings.
		// @path          : At the beginning, this is a Linux abstract socket path.
		// password@host  : Elsewhere, it's the password and host.
		// password@@path : Using @@ probably means both password and abstract socket,
		//                  but we don't need to handle this explicitly.
		// Summarised from docs of --host option in man:mpc(1).
		password := ""
		at := strings.Index(host, "@")
		if at > 0 {
			password, host = host[:at], host[at+1:]
		}
		if port == "" || strings.Contains(host, string(os.PathSeparator)) {
			cfg.MPD.Connection = "unix"
			cfg.MPD.Address = host
			cfg.MPD.Password = password
		} else {
			cfg.MPD.Connection = "tcp"
			cfg.MPD.Address = host + ":" + port
			cfg.MPD.Password = password
		}
	}
	return
}

func (cfg *Config) String() string {
	var b strings.Builder
	_ = toml.NewEncoder(&b).Encode(cfg)
	return b.String()
}

func (cfg *Config) FromClient(musicDir string, err error) {
	// Don't use client here to avoid import cycle
	if cfg.LyricsDir != "" {
		return
	}
	if err != nil {
		cfg.LyricsDir = cfg.MusicDir
		return
	}
	cfg.MusicDir = musicDir
	cfg.LyricsDir = musicDir
}

func (cfg *Config) FromOpts(opts docopt.Opts) {
	cfgMusicDir, _ := opts["--musicdir"].(string)
	cfgLyricsDir, _ := opts["--lyricsdir"].(string)
	cfgMPDConnection, _ := opts["--mpd-connection"].(string)
	cfgMPDAddress, _ := opts["--mpd-address"].(string)
	cfgMPDPassword, _ := opts["--mpd-password"].(string)
	for _, x := range &[...]*struct{ to, from *string }{
		{&cfg.MusicDir, &cfgMusicDir},
		{&cfg.LyricsDir, &cfgLyricsDir},
		{&cfg.MPD.Connection, &cfgMPDConnection},
		{&cfg.MPD.Address, &cfgMPDAddress},
		{&cfg.MPD.Password, &cfgMPDPassword},
	} {
		if *x.from != "" {
			*x.to = *x.from
		}
	}
}

// Expand expands tilde ("~") and variables ("$VAR" or "${VAR}") in paths in Config.
// Sets LyricsDir to MusicDir if empty.
func (cfg *Config) Expand() {
	cfg.MusicDir = ExpandTilde(os.ExpandEnv(cfg.MusicDir))
	cfg.LyricsDir = ExpandTilde(os.ExpandEnv(cfg.LyricsDir))
	cfg.MPD.Connection = os.ExpandEnv(cfg.MPD.Connection)
	cfg.MPD.Address = os.ExpandEnv(cfg.MPD.Address)
	cfg.MPD.Password = os.ExpandEnv(cfg.MPD.Password)
	if strings.Contains(cfg.MPD.Address, string(os.PathSeparator)) {
		cfg.MPD.Address = ExpandTilde(cfg.MPD.Address)
	}
}

// Assert return error if Config is invalid.
func (cfg *Config) Assert() error {
	if !filepath.IsAbs(cfg.MusicDir) {
		return errors.New("Invalid path in MusicDir")
	}
	if !filepath.IsAbs(cfg.LyricsDir) {
		return errors.New("Invalid path in LyricsDir")
	}
	return nil
}
