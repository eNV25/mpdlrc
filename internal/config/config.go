package config

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
)

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
		cfg.MPD.Connection = "tcp"
		cfg.MPD.Address = "localhost:6600"
		cfg.MPD.Password = ""
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

// Expand expands tilde ("~") and variables ("$VAR" or "${VAR}") in paths in Config.
// Sets LyricsDir to MusicDir if empty.
func (cfg *Config) Expand() {
	cfg.MusicDir = ExpandTilde(os.ExpandEnv(cfg.MusicDir))
	cfg.LyricsDir = ExpandTilde(os.ExpandEnv(cfg.LyricsDir))
	if strings.Contains(cfg.MPD.Address, string(os.PathSeparator)) {
		cfg.MPD.Address = ExpandTilde(os.ExpandEnv(cfg.MPD.Address))
	}

	if cfg.LyricsDir == "" && cfg.MusicDir != "" {
		cfg.LyricsDir = cfg.MusicDir
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
