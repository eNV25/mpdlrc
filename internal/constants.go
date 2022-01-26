package internal

type State uint

const (
	stateNone State = iota
	StatePlay
	StateStop
	StatePause
)
