package config

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/docopt/docopt-go"
	"github.com/pelletier/go-toml/v2"
	"go.uber.org/multierr"
)

var _ fmt.Stringer = (*Config)(nil)

type Config struct {
	LyricsDir string
	MusicDir  string

	MPD struct {
		Address    string
		Connection string
		Password   string
	}
}

func DefaultConfig() (cfg *Config) {
	cfg = &Config{}
	cfg.MusicDir = "~/Music"
	cfg.LyricsDir = ""
	host := GetEnv("MPD_HOST")
	port := GetEnv("MPD_PORT")
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

func (cfg *Config) FromFiles(files []string) (err error) {
	for _, fpath := range files {
		errr := cfg.FromFile(fpath)
		if errr != nil && !errors.Is(errr, os.ErrNotExist) {
			multierr.AppendInto(&err, errr)
		}
	}
	return
}

func (cfg *Config) FromFile(fpath string) error {
	f, err := os.Open(ExpandEnv(ExpandTilde(fpath)))
	if err != nil {
		return err
	}
	defer f.Close()
	return toml.NewDecoder(f).Decode(cfg)
}

func (cfg *Config) FromClient(musicDir string, err error) {
	if err != nil {
		return
	}
	cfg.MusicDir = musicDir
	cfg.LyricsDir = musicDir
}

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

func (cfg *Config) FromEnv(getEnv func(string) string) {
	if getEnv == nil {
		getEnv = GetEnv
	}
	cfgLyricsDir := getEnv("MPDLRC_LYRICSDIR")
	cfgMusicDir := getEnv("MPDLRC_MUSICDIR")
	cfgMPDAddress := getEnv("MPDLRC_MPD_ADDRESS")
	cfgMPDConnection := getEnv("MPDLRC_MPD_CONNECTION")
	cfgMPDPassword := getEnv("MPDLRC_MPD_PASSWORD")
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
	cfg.MusicDir = ExpandEnv(ExpandTilde(cfg.MusicDir))
	cfg.LyricsDir = ExpandEnv(ExpandTilde(cfg.LyricsDir))
	if strings.Contains(cfg.MPD.Address, string(os.PathSeparator)) {
		cfg.MPD.Address = ExpandTilde(cfg.MPD.Address)
	}
	cfg.MPD.Connection = ExpandEnv(cfg.MPD.Connection)
	cfg.MPD.Address = ExpandEnv(cfg.MPD.Address)
	cfg.MPD.Password = ExpandEnv(cfg.MPD.Password)
	if cfg.LyricsDir == "" && cfg.MusicDir != "" {
		cfg.LyricsDir = cfg.MusicDir
	}
}

// Assert return error if cfg is invalid.
func (cfg *Config) Assert() error {
	var err error
	if !filepath.IsAbs(cfg.MusicDir) {
		multierr.AppendInto(&err, errors.New("Invalid path in MusicDir"))
	}
	if !filepath.IsAbs(cfg.LyricsDir) {
		multierr.AppendInto(&err, errors.New("Invalid path in LyricsDir"))
	}
	return err
}
