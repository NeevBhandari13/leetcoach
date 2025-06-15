package models

type GPTRequest struct {
	Instructions string    `json:"instructions"`
	Input        []Message `json:"input"`
}

type GPTResponse struct {
	ResponseText string `json:"responseText"`
}
