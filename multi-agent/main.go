package main

import (
	"fmt"
	"log"
	"log/slog"
	"os"

	"github.com/nexxia-ai/aigentic"
	openai "github.com/nexxia-ai/aigentic-openai"
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

	fmt.Println("👥 Aigentic Multi-Agent System Example")
	fmt.Println("======================================")
	fmt.Println()

	model := openai.NewModel("gpt-4o-mini", getAPIKey())

	researchAgent := aigentic.Agent{
		Model:        model,
		Name:         "Researcher",
		Description:  "Expert at gathering and analyzing information on any topic",
		Instructions: "You are a research specialist. When given a topic, provide comprehensive, factual information with key insights and data points. Be thorough but concise.",
	}

	writerAgent := aigentic.Agent{
		Model:        model,
		Name:         "Writer",
		Description:  "Expert at creating clear, engaging written content",
		Instructions: "You are a professional writer. Create well-structured, engaging content based on the information provided. Use clear language and proper formatting.",
	}

	coordinator := aigentic.Agent{
		Model:        model,
		Name:         "ProjectManager",
		Description:  "Coordinates research and writing tasks to produce high-quality content",
		Instructions: "You manage a team of specialists. First, delegate research tasks to the Researcher. Then, use the research findings to have the Writer create the final content. Coordinate the work to ensure high quality output.",
		Agents:       []aigentic.Agent{researchAgent, writerAgent},
		LogLevel:     slog.LevelInfo,
	}

	response, err := coordinator.Execute("Create a brief article about the benefits of renewable energy, focusing on solar and wind power.")
	if err != nil {
		log.Fatalf("Error: %v", err)
	}
	fmt.Printf("Final Article:\n%s\n\n", response)

	fmt.Println("✅ Example completed successfully!")
}
