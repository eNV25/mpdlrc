package state

type State uint

const (
	none State = iota
	Play
	Stop
	Pause
)
