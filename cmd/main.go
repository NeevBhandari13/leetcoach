package main

import (
	"fmt"
	"os"

	"github.com/NeevBhandari13/leetcoach/internal/api"
	"github.com/NeevBhandari13/leetcoach/internal/chat"
	"github.com/NeevBhandari13/leetcoach/internal/db"
	"github.com/NeevBhandari13/leetcoach/internal/llm"
	"github.com/NeevBhandari13/leetcoach/internal/session"
	"github.com/anthropics/anthropic-sdk-go"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(".env"); err != nil {
		fmt.Println("godotenv error:", err)
	}
	sqlDb, err := db.Open(os.Getenv("DSN"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "db error: %v\n", err)
		os.Exit(1)
	}
	if err := db.Migrate(sqlDb); err != nil {
		fmt.Fprintf(os.Stderr, "migrate error: %v\n", err)
		os.Exit(1)
	}
	if err := db.Seed(sqlDb); err != nil {
		fmt.Fprintf(os.Stderr, "seed error: %v\n", err)
		os.Exit(1)
	}
	client := NewLLMClient()                               // llm client interface
	chatService := chat.NewChatService(client)             // chat service with the reply function
	sessionStore := session.NewSessionStore(sqlDb, client) // session store backed by postgres
	router := api.NewRouter(chatService, sessionStore)     // creates new router
	router.Run(":8080")                                    // router on port 8080

}

func NewLLMClient() llm.Client {
	provider := os.Getenv("LLM_PROVIDER")
	switch provider {
	case "anthropic":
		return llm.NewAnthropicClient(os.Getenv("ANTHROPIC_API_KEY"), anthropic.Model(os.Getenv("LLM_MODEL")))
	default:
		panic("Invalid LLM Provider")

	}
}
