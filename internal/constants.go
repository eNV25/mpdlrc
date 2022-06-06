package internal

type State uint

const (
	_ State = iota
	StatePlay
	StateStop
	StatePause
)
