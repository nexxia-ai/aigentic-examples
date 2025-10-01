package main

import (
	"fmt"
	"log"
	"os"

	"github.com/nexxia-ai/aigentic"
	openai "github.com/nexxia-ai/aigentic-openai"
	"github.com/nexxia-ai/aigentic/utils"
)

var (
	openaiModel = openai.NewModel("gpt-4o-mini", getAPIKey())

	// change this if you like to use ollama instead of openai
	// ollamaModel = ollama.NewModel("qwen3:1.7b", "")
	// model = openaiModel
)

func getAPIKey() string {
	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		fmt.Println("Warning: OPENAI_API_KEY environment variable not set")
		fmt.Println("Please set your OpenAI API key: export OPENAI_API_KEY=your_api_key_here")
		os.Exit(1)
	}
	return apiKey
}

// simpleAgent demonstrates a simple agent that takes an user input and returns a response
var simpleAgent = aigentic.Agent{
	Model:        openaiModel, // show use of ollama model
	Name:         "SimpleAgent",
	Description:  "A simple agent that responds to user messages",
	Instructions: "Respond to user questions in a friendly and informative way.",
}

func main() {
	utils.LoadEnvFile("../.env")

	fmt.Println("=== Running Simple Agent ===")
	response, err := simpleAgent.Execute("Hello! Can you tell me a fun fact about space?")
	if err != nil {
		log.Fatalf("Error running simple agent: %v", err)
	}
	fmt.Printf("Simple Agent Response: %s\n\n", response)
}
