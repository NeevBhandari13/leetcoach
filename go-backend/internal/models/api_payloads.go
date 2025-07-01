package models

// interview response to front end
type startInterviewResponse struct {
	SessionID    string `json:"session_id"`
	ResponseText string `json:"response_text"`
}

func PackageStartInterviewResponse(sessionID string, responseText string) startInterviewResponse {
	return startInterviewResponse{
		SessionID:    sessionID,
		ResponseText: responseText,
	}
}

// request from frontend with session id and user input
type ContinueInterviewRequest struct {
	SessionID string `json:"session_id"`
	Input     string `json:"input"`
}

// response from backend with session id and response text
type ContinueInterviewResponse struct {
	ResponseText string `json:"response_text"`
}

func PackageContinueInterviewResponse(responseText string) ContinueInterviewResponse {
	return ContinueInterviewResponse{
		ResponseText: responseText,
	}
}
