package models

type State int

const (
	IntroState State = iota
	PresentProblemState
	ClarifyState
	InitialSolutionState
	OptimisationState
	WrapUpState
)
