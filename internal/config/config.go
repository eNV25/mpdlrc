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
	cfg.MusicDir = expandTilde(os.ExpandEnv(cfg.MusicDir))
	cfg.LyricsDir = expandTilde(os.ExpandEnv(cfg.LyricsDir))
	if strings.ContainsRune(cfg.MPD.Address, os.PathSeparator) {
		cfg.MPD.Address = expandTilde(os.ExpandEnv(cfg.MPD.Address))
	}

	if cfg.LyricsDir == "" && cfg.MusicDir != "" {
		cfg.LyricsDir = cfg.MusicDir
	}
}

func expandTilde(str string) string {
	switch {
	case strings.HasPrefix(str, "~"):
		// ~ or ~/path or ~user/path
		u, p, _ := stringsCut(str[1:], string(os.PathSeparator))
		return filepath.Join(HomeDirUser(u), p) // calls filepath.Clean
	default:
		// path or /path
		return str
	}
}

// stringsCut slices s around the first instance of sep,
// returning the text before and after sep.
// The found result reports whether sep appears in s.
// If sep does not appear in s, cut returns s, "", false.
//
// Copied from standard library to allow build with go pre-1.18
func stringsCut(s, sep string) (before, after string, found bool) {
	if i := strings.Index(s, sep); i >= 0 {
		return s[:i], s[i+len(sep):], true
	}
	return s, "", false
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
