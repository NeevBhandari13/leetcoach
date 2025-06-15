package models

type GPTRequest struct {
	Instructions string    `json:"instructions"`
	Input        []Message `json:"input"`
}

type GPTResponse struct {
	Reply string `json:"reply"`
}
