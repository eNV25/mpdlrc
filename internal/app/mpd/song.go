package mpd

import (
	"path"
)

type Song map[string]string

func (s Song) ID() string {
	return s["Id"]
}

func (s Song) Title() string {
	return s["Title"]
}

func (s Song) Artist() string {
	return s["Artist"]
}

func (s Song) Album() string {
	return s["Album"]
}

func (s Song) File() string {
	return s["file"]
}

func (s Song) LRCFile() string {
	file := s.File()
	return file[:(len(file)-len(path.Ext(file)))] + ".lrc"
}
