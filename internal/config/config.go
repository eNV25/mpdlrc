package config

import (
	"errors"
	"os"
	"path"
	"strings"
)

type Config struct {
	MusicDir  string
	LyricsDir string
	Debug     bool

	MPD struct {
		Connection string
		Address    string
		Password   string
	}
}

func DefaultConfig() (cfg *Config) {
	cfg = new(Config)
	cfg.MusicDir = "~/Music"
	cfg.LyricsDir = ""
	cfg.MPD.Connection = "tcp"
	cfg.MPD.Address = "localhost:6600"
	cfg.MPD.Password = ""
	cfg.Debug = false
	return cfg
}

// Expand expands tilde ("~") and variables ("$VAR" or "${VAR}") in paths in Config.
// Sets LyricsDir to MusicDir if empty.
func (cfg *Config) Expand() {
	cfg.MusicDir = expandTilde(cfg.MusicDir)
	cfg.MusicDir = os.ExpandEnv(cfg.MusicDir)
	cfg.LyricsDir = expandTilde(cfg.LyricsDir)
	cfg.LyricsDir = os.ExpandEnv(cfg.LyricsDir)
	cfg.MPD.Address = expandTilde(cfg.MPD.Address)
	cfg.MPD.Address = os.ExpandEnv(cfg.MPD.Address)

	if cfg.LyricsDir == "" && cfg.MusicDir != "" {
		cfg.LyricsDir = cfg.MusicDir
	}
}

func expandTilde(str string) string {
	if str != "" && (str == "~" || str[:2] == "~/") {
		// ~/path/
		return HomeDir() + str[1:]
	} else if str[:1] == "~" {
		// ~user/path/
		sp := strings.Split(str[1:], "/")
		return path.Join(HomeDirUser(sp[0]), path.Join(sp[1:]...))
	} else {
		// path/ or /path/
		return str
	}
}

// Assert return error if Config is invalid.
func (cfg *Config) Assert() error {
	if cfg.MusicDir == "" || cfg.MusicDir[:1] != "/" {
		return errors.New("Invalid path in MusicDir")
	}
	if cfg.LyricsDir == "" || cfg.LyricsDir[:1] != "/" {
		return errors.New("Invalid path in LyricsDir")
	}
	return nil
}
