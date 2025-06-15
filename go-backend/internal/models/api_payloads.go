package models

type InterviewResponse struct {
	SessionID    string `json:"sessionID"`
	ResponseText string `json:"responseText"`
}

type ContinueInterviewRequest struct {
	SessionID string `json:"sessionID"`
	Input     string `json:"input"`
}
