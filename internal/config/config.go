// Package config implements the [Config] structure, and holds related functions.
package config

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/docopt/docopt-go"
	"github.com/pelletier/go-toml/v2"

	"github.com/env25/mpdlrc/internal/client"
	"github.com/env25/mpdlrc/internal/dirs"
	"github.com/env25/mpdlrc/internal/mpdconf"
)

var _ fmt.Stringer = (*Config)(nil)

// Config holds the configuration.
type Config struct {
	LyricsDir string
	MusicDir  string

	MPD struct {
		Address    string
		Connection string
		Password   string
	}
}

// DefaultConfig returns the default configuration.
func DefaultConfig() (cfg *Config) {
	cfg = &Config{}
	cfg.MusicDir = "~/Music"
	cfg.LyricsDir = ""
	host := dirs.GetEnv("MPD_HOST")
	port := dirs.GetEnv("MPD_PORT")
	if host == "" {
		if port == "" {
			// Not enough information.
			// [client.NewMPDClient] can try guess a few times and set these fields.
			cfg.MPD.Connection = ""
			cfg.MPD.Address = ""
			cfg.MPD.Password = ""
		} else {
			cfg.MPD.Connection = "tcp"
			cfg.MPD.Address = ":" + port
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
		if port == "" || strings.ContainsRune(host, os.PathSeparator) {
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

// FromFiles merges the configuration fore files.
func (cfg *Config) FromFiles(files []string) error {
	defaultFiles := [...]string{filepath.Join(dirs.GetEnv("XDG_CONFIG_HOME"), "mpdlrc", "config.toml")}
	if len(files) == 0 {
		files = defaultFiles[:]
	}
	var errs []error
	for _, fpath := range files {
		err := cfg.FromFile(fpath)
		if err != nil && !errors.Is(err, os.ErrNotExist) {
			errs = append(errs, err)
		}
	}
	return errors.Join(errs...)
}

// FromFile merges the configuration from fpath.
func (cfg *Config) FromFile(fpath string) error {
	f, err := os.Open(dirs.ExpandEnv(dirs.ExpandTilde(fpath)))
	if err != nil {
		return err
	}
	defer f.Close()
	return toml.NewDecoder(f).Decode(cfg)
}

// FromClient merges the configuration from client.
func (cfg *Config) FromClient(c client.Client) {
	musicDir, err := c.MusicDir()
	if err != nil {
		if _, ok := c.(*client.MPDClient); ok {
			cfg.fromMPDConfig()
		}
		return
	}
	cfg.MusicDir = musicDir
	cfg.fixLyricsDir()
}

func (cfg *Config) fromMPDConfig() {
	if cfg.MusicDir != "" {
		return
	}
	for _, fpath := range &[...]string{
		filepath.Join(dirs.GetEnv("XDG_CONFIG_HOME"), "mpd", "mpd.conf"),
		filepath.Join(dirs.HomeDir(""), ".mpdconf"),
		filepath.Join(dirs.HomeDir(""), ".mpd", "mpd.conf"),
		filepath.Join(dirs.RootDir(), "etc", "mpd.conf"),
	} {
		func(fpath string) {
			f, err := os.Open(fpath)
			if err != nil {
				return
			}
			defer f.Close()
			var s mpdconf.Scanner
			s.Init(f)
			for s.Next() {
				if v, ok := s.Str("music_directory"); ok {
					cfg.MusicDir = v
				}
			}
			cfg.fixLyricsDir()
		}(fpath)
	}
}

// FromOpts merges the configuration from [docopt.Opts] opts.
func (cfg *Config) FromOpts(opts docopt.Opts) {
	cfgLyricsDir, _ := opts["--lyricsdir"].(string)
	cfgMusicDir, _ := opts["--musicdir"].(string)
	cfgMPDAddress, _ := opts["--mpd-address"].(string)
	cfgMPDConnection, _ := opts["--mpd-connection"].(string)
	cfgMPDPassword, _ := opts["--mpd-connection"].(string)
	for _, x := range &[...]*struct{ from, to *string }{
		{&cfgLyricsDir, &cfg.LyricsDir},
		{&cfgMusicDir, &cfg.MusicDir},
		{&cfgMPDAddress, &cfg.MPD.Address},
		{&cfgMPDConnection, &cfg.MPD.Connection},
		{&cfgMPDPassword, &cfg.MPD.Password},
	} {
		if *x.from != "" {
			*x.to = *x.from
		}
	}
}

// FromEnv merges the configuration from env.
func (cfg *Config) FromEnv(env func(string) string) {
	if env == nil {
		env = dirs.GetEnv
	}
	cfgLyricsDir := env("MPDLRC_LYRICSDIR")
	cfgMusicDir := env("MPDLRC_MUSICDIR")
	cfgMPDAddress := env("MPDLRC_MPD_ADDRESS")
	cfgMPDConnection := env("MPDLRC_MPD_CONNECTION")
	cfgMPDPassword := env("MPDLRC_MPD_PASSWORD")
	for _, x := range &[...]*struct{ from, to *string }{
		{&cfgLyricsDir, &cfg.LyricsDir},
		{&cfgMusicDir, &cfg.MusicDir},
		{&cfgMPDAddress, &cfg.MPD.Address},
		{&cfgMPDConnection, &cfg.MPD.Connection},
		{&cfgMPDPassword, &cfg.MPD.Password},
	} {
		if *x.from != "" {
			*x.to = *x.from
		}
	}
}

// Expand expands tilde ("~") and variables ("$VAR" or "${VAR}") in paths in cfg.
// Sets LyricsDir to MusicDir if empty.
func (cfg *Config) Expand() {
	cfg.MusicDir = dirs.ExpandEnv(dirs.ExpandTilde(cfg.MusicDir))
	cfg.LyricsDir = dirs.ExpandEnv(dirs.ExpandTilde(cfg.LyricsDir))
	if strings.ContainsRune(cfg.MPD.Address, os.PathSeparator) {
		cfg.MPD.Address = dirs.ExpandTilde(cfg.MPD.Address)
	}
	cfg.MPD.Connection = dirs.ExpandEnv(cfg.MPD.Connection)
	cfg.MPD.Address = dirs.ExpandEnv(cfg.MPD.Address)
	cfg.MPD.Password = dirs.ExpandEnv(cfg.MPD.Password)
	cfg.fixLyricsDir()
}

func (cfg *Config) fixLyricsDir() {
	if cfg.LyricsDir == "" {
		cfg.LyricsDir = cfg.MusicDir
	}
}

// Assert return error if cfg is invalid.
func (cfg *Config) Assert() error {
	var errs []error
	if !filepath.IsAbs(cfg.MusicDir) {
		errs = append(errs, fmt.Errorf("Config: invalid path: MusicDir must be absolute: %q", cfg.MusicDir))
	}
	if !filepath.IsAbs(cfg.LyricsDir) {
		errs = append(errs, fmt.Errorf("Config: invalid path: LyricsDir must be absolute: %q", cfg.LyricsDir))
	}
	if strings.ContainsRune(cfg.MPD.Address, os.PathSeparator) && !filepath.IsAbs(cfg.MPD.Address) {
		errs = append(errs, fmt.Errorf("Config: invalid path: MPD.Address (unix socket) must be absolute: %q", cfg.MPD.Address))
	}
	return errors.Join(errs...)
}
