package upath

import (
	"path"
)

// TrimExt returns the path without its file name extension.
// The extension is the suffix beginning at the final dot
// in the final slash-separated element of path;
// it is empty if there is no dot.
func TrimExt(p string) string {
	return p[:len(p)-len(path.Ext(p))]
}

// ReplaceExt returns the path with its file name extension replaced
// by the provided one. The extension is the suffix beginning at
// the final dot in the final slash-separated element of path;
// it is empty if there is no dot.
func ReplaceExt(p string, ext string) string {
	return TrimExt(p) + ext
}
