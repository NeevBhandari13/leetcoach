package models

// interview response to front end
type InterviewResponse struct {
	SessionID    string `json:"sessionID"`
	ResponseText string `json:"responseText"`
}

func PackageInterviewResponse(sessionID string, responseText string) InterviewResponse {
	return InterviewResponse{
		SessionID:    sessionID,
		ResponseText: responseText,
	}
}

// request from frontend with session id and user input
type ContinueInterviewRequest struct {
	SessionID string `json:"sessionID"`
	Input     string `json:"input"`
}
