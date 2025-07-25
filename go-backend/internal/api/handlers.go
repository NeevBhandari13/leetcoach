package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/neevbhandari13/leetcoach/internal/ai"
	"github.com/neevbhandari13/leetcoach/internal/interview"
	"github.com/neevbhandari13/leetcoach/internal/models"
	"github.com/neevbhandari13/leetcoach/internal/prompts"
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

func startInterviewHandler(client *ai.GPTClient, sessionStore *interview.SessionStore) gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessionStore.CreateSession("")

		reply := "Hello! Welcome to LeetCoach! Today we are going to be running a technical interview going over a coding problem together. Remember I'm here to help you! Are you ready to get started?"

		// add the llm response to chat history
		sessionStore.UpdateChatHistory(session.SessionID, interview.PackageMessage("assistant", reply))

		// package response to front end as models.InterviewResponse
		response := models.PackageStartInterviewResponse(session.SessionID, reply)

		c.JSON(http.StatusOK, response)
	}

}

func continueInterviewHandler(client *ai.GPTClient, sessionStore *interview.SessionStore) gin.HandlerFunc {
	return func(c *gin.Context) {
		// request is a variable to hold the request body
		var request models.ContinueInterviewRequest

		// BindJSON takes the request body and binds it to the request variable
		if err := c.BindJSON(&request); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		// unpack request body
		sessionID := request.SessionID
		userInput := request.Input

		// convert user input into a message
		userMessage := interview.PackageMessage("user", userInput)

		// append user message to chat history and retrieve it
		chatHistory := sessionStore.AppendAndReadChatHistory(sessionID, userMessage)

		// get instructions and developer prompt
		instructions := prompts.GetInstructions()
		state, err := sessionStore.GetState(sessionID)
		// handle error
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}
		developerPrompt := prompts.GetDeveloperPrompt(state)

		// package gpt request
		gptRequest := ai.PackageGPTRequest(instructions, developerPrompt, chatHistory)

		reply, nextState, err := client.CallGPT(gptRequest)
		// handle error
		if err != nil {
			// send back response with StatusInternalServerError code and error message in gin.H
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}

		// add ai response to chat history
		sessionStore.UpdateChatHistory(sessionID, interview.PackageMessage("assistant", reply))

		// set next state
		err = sessionStore.SetState(sessionID, nextState)

		// handle error
		if err != nil {
			// send back response with StatusInternalServerError code and error message in gin.H
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}

		// package response to front end as models.InterviewResponse
		response := models.PackageContinueInterviewResponse(reply)

		// send back response
		c.JSON(http.StatusOK, response)
	}

}
