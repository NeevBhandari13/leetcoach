package interview

import (
	"log"
	"runtime/debug"

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

	session, ok := s.sessions[sessionID]
	if !ok {
		log.Printf("❌ sessionID not found in session store: %s", sessionID)
		debug.PrintStack() // helpful to trace where this was called
		return
	} else {
		log.Printf("✅ sessionID found in session store: %s", sessionID)
	}

	session.ChatHistory = append(session.ChatHistory, message)
	// get chat history from session and append the new message to it
	// need to update directly for it to work
}

// useful for when we get user message and need to send to api
func (s *SessionStore) AppendAndReadChatHistory(sessionID string, message models.Message) []models.Message {
	s.mu.Lock()
	defer s.mu.Unlock()

	session, ok := s.sessions[sessionID]
	if !ok {
		log.Printf("❌ sessionID not found in session store: %s", sessionID)
		debug.PrintStack() // helpful to trace where this was called
		return []models.Message{}
	} else {
		log.Printf("✅ sessionID found in session store: %s", sessionID)
	}

	session.ChatHistory = append(session.ChatHistory, message)
	return session.ChatHistory
}
