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
	cfg.MPD.Connection = "tcp"
	cfg.MPD.Address = "localhost:6600"
	cfg.MPD.Password = ""
	return cfg
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
