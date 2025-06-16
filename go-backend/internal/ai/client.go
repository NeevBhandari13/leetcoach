package ai

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/neevbhandari13/leetcoach/internal/models"
)

type GPTClient struct {
	BaseURL string
}

// init runs before anything else in the package good for initialising variables etc
func NewGPTClient() *GPTClient {
	// need to specifically load in .env file
	err := godotenv.Load("../.env")
	if err != nil {
		log.Println("No .env file found or could not load .env file")
	}
	// get base URL to assign to gpt client
	baseURL := os.Getenv("PYTHON_SERVICE_URL")
	if baseURL == "" {
		log.Fatal("PYTHON_SERVICE_URL is not set")
	}

	return &GPTClient{
		BaseURL: baseURL,
	}

}

func (client *GPTClient) CallGPT(gptRequest models.GPTRequest) (string, models.State, error) {

	// get url for chat endpoint in python microservice
	chatURL := fmt.Sprintf("%s/chat", client.BaseURL)

	// convert gptRequest to json
	// is a slice of bytes under the hood
	bodyBytes, err := json.Marshal(gptRequest)
	if err != nil {
		return "", models.NilState, err
	}

	// make HTTP POST request to python microservice
	// http.Post requires a body of type io.Reader
	// bytes.NewBuffer wraps our []bytes into a *bytes.Buffer
	// "application/json" is the content type header
	response, err := http.Post(chatURL, "application/json", bytes.NewBuffer(bodyBytes))
	if err != nil {
		return "", models.NilState, err
	}

	// parse response body
	var gptResponse models.GPTResponse

	// decode json into gptResponse
	err = json.NewDecoder(response.Body).Decode(&gptResponse)
	if err != nil {
		return "", models.NilState, err
	}
	// close response body
	defer response.Body.Close()

	// handle error when no state is sent back
	if gptResponse.Reply == "" {
		return "", models.NilState, fmt.Errorf("empty response from GPT")
	}

	return gptResponse.Reply, gptResponse.CurrentState, nil

}
