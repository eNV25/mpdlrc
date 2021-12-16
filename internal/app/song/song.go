package song

type Song interface {
	ID() string
	Title() string
	Artist() string
	Album() string
	File() string
	LRCFile() string
}
