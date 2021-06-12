package config

import (
	"errors"
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/pelletier/go-toml/v2"
)

type Config struct {
	MusicDir  string
	LyricsDir string
	Debug     bool

	MPD struct {
		Protocol string
		Address  string
	}
}

func DefaultConfig() (cfg *Config) {
	cfg = new(Config)
	cfg.MusicDir = "~/Music"
	cfg.LyricsDir = ""
	cfg.MPD.Protocol = "tcp"
	cfg.MPD.Address = "localhost:6600"
	cfg.Debug = false
	return cfg
}

// MergeTOMLFile merges TOML file with Config.
func (cfg *Config) MergeTOMLFile(fpath string) (err error) {
	var f *os.File

	if f, err = os.Open(fpath); err != nil {
		return fmt.Errorf("open config file: %w", err)
	}
	defer f.Close()

	if err = toml.NewDecoder(f).Decode(cfg); err != nil {
		return fmt.Errorf("decode config file: %w", err)
	}

	return nil
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
		return HomeDir() + str[1:]
	} else if str[:1] == "~" {
		sp := strings.Split(str[1:], "/")
		return path.Join(HomeDirUser(sp[0]), path.Join(sp[1:]...))
	} else {
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
