package models

// interview response to front end
type startInterviewResponse struct {
	SessionID string `json:"sessionId"`
	Reply     string `json:"reply"`
}

func PackageStartInterviewResponse(sessionID string, reply string) startInterviewResponse {
	return startInterviewResponse{
		SessionID: sessionID,
		Reply:     reply,
	}
}

// request from frontend with session id and user input
type ContinueInterviewRequest struct {
	SessionID string `json:"sessionId"`
	Input     string `json:"input"`
}

// response from backend with session id and response text
type ContinueInterviewResponse struct {
	Reply string `json:"reply"`
}

func PackageContinueInterviewResponse(reply string) ContinueInterviewResponse {
	return ContinueInterviewResponse{
		Reply: reply,
	}
}
