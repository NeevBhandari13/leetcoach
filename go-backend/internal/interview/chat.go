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

// adds message to chat history
func AddMessage(chatHistory []models.Message, message models.Message) []models.Message {
	chatHistory = append(chatHistory, message)
	return chatHistory
}
