package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/nexxia-ai/aigentic"
	openai "github.com/nexxia-ai/aigentic-openai"
	"github.com/nexxia-ai/aigentic/memory"
	"github.com/nexxia-ai/aigentic/utils"
)

func getAPIKey() string {
	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		fmt.Println("Error: OPENAI_API_KEY environment variable not set")
		fmt.Println("Please set your OpenAI API key: export OPENAI_API_KEY=your_api_key_here")
		os.Exit(1)
	}
	return apiKey
}

func main() {
	utils.LoadEnvFile("../.env")

	fmt.Println("💾 Aigentic Memory System Example")
	fmt.Println("==================================")
	fmt.Println()

	model := openai.NewModel("gpt-4o-mini", getAPIKey())

	session := aigentic.NewSession(context.Background())

	agent := aigentic.Agent{
		Model:        model,
		Name:         "PersonalAssistant",
		Description:  "A personal assistant that remembers user preferences and context",
		Instructions: "You are a personal assistant. Remember user preferences, past interactions, and important information using session memory. Use save_memory tool with compartment 'session' for information that should persist across conversations.",
		Session:      session,
		Memory:       memory.NewMemory(),
	}

	fmt.Println("First conversation:")
	response, err := agent.Execute("My name is Alice and I prefer morning meetings. I'm working on a project about renewable energy.")
	if err != nil {
		log.Fatalf("Error: %v", err)
	}
	fmt.Printf("Response: %s\n\n", response)

	fmt.Println("Second conversation (new agent run, same session):")
	response, err = agent.Execute("What's my name and what project am I working on?")
	if err != nil {
		log.Fatalf("Error: %v", err)
	}
	fmt.Printf("Response: %s\n\n", response)

	fmt.Println("Third conversation:")
	response, err = agent.Execute("When do I prefer to have meetings?")
	if err != nil {
		log.Fatalf("Error: %v", err)
	}
	fmt.Printf("Response: %s\n\n", response)

	fmt.Println("✅ Example completed successfully!")
}
