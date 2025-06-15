package interview

import (
	"github.com/neevbhandari13/leetcoach/internal/models"
)

var DeveloperPrompts = map[models.State]string{
	models.IntroState:           "Start the interview by greeting the candidate warmly and professionally. Let them know you’ll be working together on a coding problem. Do not present the problem yet.",
	models.PresentProblemState:  "Present a vague, open-ended version of the problem to encourage the candidate to ask clarifying questions. Do not mention specific terms like 'graph', 'palindrome', etc.",
	models.ClarifyState:         "Answer the candidate's clarifying questions clearly and concisely. Confirm or deny what they ask without revealing more than necessary.",
	models.InitialSolutionState: "Discuss the candidate’s initial approach collaboratively. Encourage them to explain their reasoning and analyze time and space complexity. Do not provide answers or write code.",
	models.OptimisationState:    "Support the candidate in refining their solution. Ask guiding questions about potential tradeoffs or improvements, but let them lead the thinking. Be collaborative and exploratory.",
	models.WrapUpState:          "Conclude the interview with positive feedback. Acknowledge good questions, improvements made, or strong communication. Keep it constructive and affirming.",
}

func GetDeveloperPrompt(state models.State) string {
	return DeveloperPrompts[state]
}
