package models

import (
	interview "github.com/neevbhandari13/leetcoach/internal/interview"
)

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type Session struct {
	ID          string
	State       interview.State
	ChatHistory []Message
	ProblemText string
}
