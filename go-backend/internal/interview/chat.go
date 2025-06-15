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
func GetChatHistory(sessionID string) []models.Message {
	sessionsMutex.Lock()
	defer sessionsMutex.Unlock()
	return sessions[sessionID].ChatHistory
}

func UpdateChatHistory(sessionID string, message models.Message) {
	// lock the sessions map
	sessionsMutex.Lock()
	defer sessionsMutex.Unlock()
	// get chat history from session and append the new message to it
	// need to update directly for it to work
	sessions[sessionID].ChatHistory = append(sessions[sessionID].ChatHistory, message)
}

// useful for when we get user message and need to send to api
func AppendAndReadChatHistory(sessionID string, message models.Message) []models.Message {
	sessionsMutex.Lock()
	defer sessionsMutex.Unlock()
	sessions[sessionID].ChatHistory = append(sessions[sessionID].ChatHistory, message)
	return sessions[sessionID].ChatHistory
}
