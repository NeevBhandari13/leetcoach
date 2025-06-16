package interview

import (
	"github.com/neevbhandari13/leetcoach/internal/models"
)

// turns role and content into a message
func PackageMessage(role string, content string) models.Message {
	return models.Message{
		Role:    role,
		Content: content,
	}
}

// Gets chat history from sessionID
func (s *SessionStore) GetChatHistory(sessionID string) []models.Message {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.sessions[sessionID].ChatHistory
}

func (s *SessionStore) UpdateChatHistory(sessionID string, message models.Message) {
	// lock the sessions map
	s.mu.Lock()
	defer s.mu.Unlock()
	// get chat history from session and append the new message to it
	// need to update directly for it to work
	s.sessions[sessionID].ChatHistory = append(s.sessions[sessionID].ChatHistory, message)
}

// useful for when we get user message and need to send to api
func (s *SessionStore) AppendAndReadChatHistory(sessionID string, message models.Message) []models.Message {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.sessions[sessionID].ChatHistory = append(s.sessions[sessionID].ChatHistory, message)
	return s.sessions[sessionID].ChatHistory
}
