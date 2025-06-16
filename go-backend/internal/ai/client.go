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

// declare this as a package level variable so all the functions can access it
// only variable declarations can be done outside a function like this
var (
	PythonServiceURL string
)

// init runs before anything else in the package good for initialising variables etc
func init() {
	// need to specifically load in .env file
	err := godotenv.Load("../.env")
	if err != nil {
		log.Println("No .env file found or could not load .env file")
	}
	// using = instead of := assigns the value to the package variable
	// rather than declaring a local variable
	PythonServiceURL = os.Getenv("PYTHON_SERVICE_URL")
	if PythonServiceURL == "" {
		log.Fatal("PYTHON_SERVICE_URL is not set")
	}

}

func CallGPT(gptRequest models.GPTRequest) (string, models.State, error) {

	// get url for chat endpoint in python microservice
	chatURL := fmt.Sprintf("%s/chat", PythonServiceURL)

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

	return gptResponse.Reply, gptResponse.CurrentState, nil

}
