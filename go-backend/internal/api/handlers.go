package api

import (
	"github.com/gin-gonic/gin"
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

}
