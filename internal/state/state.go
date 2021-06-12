package state

type State uint32

const (
	_ State = iota
	PlayState
	StopState
	PauseState
)
