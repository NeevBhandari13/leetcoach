package models

type GPTRequest struct {
	Instructions string    `json:"instructions"`
	ChatHistory  []Message `json:"chatHistory"`
}

type GPTResponse struct {
	Reply string `json:"reply"`
}
