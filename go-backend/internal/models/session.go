package models

type Session struct {
	SessionID   string    `json:"sessionId"`
	State       State     `json:"state"`
	ChatHistory []Message `json:"chatHistory"`
	ProblemText string    `json:"problemText"`
}
