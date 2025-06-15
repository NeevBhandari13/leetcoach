package ai

import (
	"github.com/neevbhandari13/leetcoach/internal/interview"
	"github.com/neevbhandari13/leetcoach/internal/models"
)

// gets instructions, developer prompt and chat history and packages into a GPTRequest
func PackageGPTRequest(instructions string, developerPrompt string, chatHistory []models.Message) models.GPTRequest {
	Input := append(chatHistory, interview.PackageMessage("developer", developerPrompt))
	return models.GPTRequest{
		Instructions: instructions,
		Input:        Input,
	}

}
