package models

type Session struct {
	SessionID   string    `json:"session_id"`
	State       State     `json:"state"`
	ChatHistory []Message `json:"chat_history"`
	ProblemText string    `json:"problem_text"`
}
