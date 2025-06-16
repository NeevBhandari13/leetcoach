package models

type GPTRequest struct {
	Instructions string    `json:"instructions"`
	Input        []Message `json:"input"`
}

type AiServiceResponse struct {
	Response string `json:"response"`
}

type GPTResponse struct {
	Reply        string `json:"reply"`
	CurrentState State  `json:"current_state"`
}
