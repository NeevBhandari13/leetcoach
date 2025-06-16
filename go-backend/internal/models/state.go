package models

type State string

const (
	IntroState           State = "intro"
	PresentProblemState  State = "present_problem"
	ClarifyState         State = "clarify"
	InitialSolutionState State = "initial_solution"
	OptimisationState    State = "optimisation"
	WrapUpState          State = "wrap_up"
	NilState             State = ""
)

func ParseState(s string) State {
	switch s {
	case "intro":
		return IntroState
	case "present_problem":
		return PresentProblemState
	case "clarify":
		return ClarifyState
	case "initial_solution":
		return InitialSolutionState
	case "optimisation":
		return OptimisationState
	case "wrap_up":
		return WrapUpState
	default:
		return NilState // or panic/log error depending on your needs
	}
}
