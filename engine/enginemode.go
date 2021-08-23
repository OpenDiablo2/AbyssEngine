package engine

type EngineMode int

const (
	EngineModeBoot EngineMode = iota
	EngineModeGame
	EngineModeError
)
