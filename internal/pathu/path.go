package pathu

import (
	"path"
	"strings"
)

// TrimExt returns the path without its file name extension.
// The extension is the suffix beginning at the final dot
// in the final slash-separated element of path;
// it is empty if there is no dot.
func TrimExt(p string) string {
	return strings.TrimSuffix(p, path.Ext(p))
}
