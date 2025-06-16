package prompts

import (
	"fmt"

	"github.com/neevbhandari13/leetcoach/internal/models"
)

// map of the instructions for each state
var StateInstructions = map[models.State]string{
	models.IntroState: `
Start the interview by greeting the candidate warmly and professionally.
Let them know you’ll be working together on a coding problem.
Do not present the problem yet.
`,

	models.PresentProblemState: `
Present a vague, open-ended version of the coding problem to encourage the candidate to ask clarifying questions.
Avoid using specific terms like 'graph', 'palindrome', etc.
`,

	models.ClarifyState: `
Answer the candidate's clarifying questions clearly and concisely.
Confirm or deny their assumptions without revealing more than necessary.
`,

	models.InitialSolutionState: `
Collaboratively discuss the candidate’s initial approach to the problem.
Encourage them to explain their reasoning, including time and space complexity.
Avoid providing direct answers or writing code.
`,

	models.OptimisationState: `
Support the candidate in refining their solution.
Ask guiding questions about trade-offs, edge cases, or improvements.
Let the candidate lead the thinking and exploration.
`,

	models.WrapUpState: `
Conclude the interview with positive and constructive feedback.
Acknowledge any thoughtful questions, optimizations, or communication strengths.
End the session with encouragement and thanks.
No further transitions are needed after this.
`,
}

func getStateInstructions(state models.State) string {
	return StateInstructions[state]
}

var DeveloperPrompts = map[models.State]string{
	models.IntroState: fmt.Sprintf(`
You are in the 'intro' state. Once the candidate acknowledges and is ready to begin, transition to 'present_problem'.

intro:
%s

present_problem:
%s

`, getStateInstructions(models.IntroState), getStateInstructions(models.PresentProblemState)),

	models.PresentProblemState: fmt.Sprintf(`
You are in the 'present_problem' state, if the candidate starts asking clarifying questions, transition to 'clarify', if they jump into a solution, transtion to 'initial_solution'.

present_problem:
%s

clarify:
%s

initial_solution:
%s

`, getStateInstructions(models.PresentProblemState), getStateInstructions(models.ClarifyState), getStateInstructions(models.InitialSolutionState)),

	models.ClarifyState: fmt.Sprintf(`
You are in the 'clarify' state. When the candidate begins discussing a solution approach, transition to 'initial_solution'.

clarify:
%s

initial_solution:
%s

`, getStateInstructions(models.ClarifyState), getStateInstructions(models.InitialSolutionState)),

	models.InitialSolutionState: fmt.Sprintf(`
You are in the 'initial_solution' state. Once the candidate has a working baseline solution, transition to 'optimisation'.

initial_solution:
%s

optimisation:
%s

`, getStateInstructions(models.InitialSolutionState), getStateInstructions(models.OptimisationState)),

	models.OptimisationState: fmt.Sprintf(`
You are in the 'optimisation' state. When the candidate has finalized their approach or there’s nothing more to improve, transition to 'wrap_up'.

optimisation:
%s	

wrap_up:
%s

`, getStateInstructions(models.OptimisationState), getStateInstructions(models.WrapUpState)),

	models.WrapUpState: fmt.Sprintf(`
You are in the 'wrap_up' state. No further transitions are needed after this.

wrap_up:
%s

`, getStateInstructions(models.WrapUpState)),
}

func GetDeveloperPrompt(state models.State) string {
	return DeveloperPrompts[state]
}
