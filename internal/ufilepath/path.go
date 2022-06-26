package ufilepath

import (
	"path/filepath"
	"runtime"
)

// FromSlash returns the result of replacing each slash ('/') character
// in path with a separator character. Multiple slashes are replaced
// by multiple separators.
func FromSlash(path string) string {
	if runtime.GOOS == "windows" {
		// Avoid strings.ReplaceAll since windows accepts forward slash.
		return path
	}
	return filepath.FromSlash(path)
}
