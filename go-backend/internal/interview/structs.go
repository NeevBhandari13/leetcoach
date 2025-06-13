package interview

import ()

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type Session struct {
	SessionID   string    `json:"session_id"`
	State       State     `json:"state"`
	ChatHistory []Message `json:"chat_history"`
	ProblemText string    `json:"problem_text"`
}
