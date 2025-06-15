package api

import (
	"github.com/gin-gonic/gin"
	"github.com/neevbhandari13/leetcoach/internal/ai"
	"github.com/neevbhandari13/leetcoach/internal/interview"
	"github.com/neevbhandari13/leetcoach/internal/prompts"
	"net/http"
)

// gin.Context is the context of the HTTP request,
// it is a wrapper around the incoming HTTP request and the outgoing HTTP response
// tool for reading what the client sent and deciding what to send back
func testHandler(c *gin.Context) {
	// c.JSON sends a JSON response back to the client
	// gin.H is type H map[string]interface{}
	// defines a map that has keys of string and values of interface{} (any type)
	c.JSON(http.StatusOK, gin.H{
		"message": "Test",
	})
}

func startInterviewHandler(c *gin.Context) {
	session := interview.CreateSession()

	instructions := prompts.GetInstructions()
	developerPrompt := prompts.GetDeveloperPrompt(session.State)
	chatHistory := interview.GetChatHistory(session.SessionID)

	gptRequest := ai.PackageGPTRequest(instructions, developerPrompt, chatHistory)

	response, err := ai.CallGPT(gptRequest)
	// handle error
	if err != nil {
		// send back response with StatusInternalServerError code and error message in gin.H
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	interview.UpdateChatHistory(session.SessionID, interview.PackageMessage("assistant", response))

	c.JSON(http.StatusOK, gin.H{
		"session_id": session.SessionID,
		"problem":    session.ProblemText,
		"state":      session.State,
		"chat":       session.ChatHistory,
	})

}
