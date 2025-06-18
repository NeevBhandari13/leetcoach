package models

// interview response to front end
type InterviewResponse struct {
	SessionID    string `json:"session_id"`
	ResponseText string `json:"response_text"`
}

func PackageInterviewResponse(sessionID string, responseText string) InterviewResponse {
	return InterviewResponse{
		SessionID:    sessionID,
		ResponseText: responseText,
	}
}

// request from frontend with session id and user input
type ContinueInterviewRequest struct {
	SessionID string `json:"session_id"`
	Input     string `json:"input"`
}
