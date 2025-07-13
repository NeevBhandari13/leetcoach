package models

// interview response to front end
type startInterviewResponse struct {
	SessionID    string `json:"sessionId"`
	ResponseText string `json:"responseText"`
}

func PackageStartInterviewResponse(sessionID string, responseText string) startInterviewResponse {
	return startInterviewResponse{
		SessionID:    sessionID,
		ResponseText: responseText,
	}
}

// request from frontend with session id and user input
type ContinueInterviewRequest struct {
	SessionID string `json:"sessionId"`
	Input     string `json:"input"`
}

// response from backend with session id and response text
type ContinueInterviewResponse struct {
	ResponseText string `json:"responseText"`
}

func PackageContinueInterviewResponse(responseText string) ContinueInterviewResponse {
	return ContinueInterviewResponse{
		ResponseText: responseText,
	}
}
