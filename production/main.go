package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"os"
	"time"

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

	fmt.Println("Production-Ready Agent Example")
	fmt.Println("==============================")
	fmt.Println()

	model := openai.NewModel("gpt-4o-mini", getAPIKey())

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	session := aigentic.NewSession(ctx)
	trace := aigentic.NewTrace()

	type FetchDataInput struct {
		Query string `json:"query" description:"The data query to execute"`
	}

	dataTool := aigentic.NewTool(
		"fetch_data",
		"Fetches data from a database or API",
		func(run *aigentic.AgentRun, input FetchDataInput) (string, error) {
			if run.Session().Context.Err() != nil {
				return "", run.Session().Context.Err()
			}

			time.Sleep(100 * time.Millisecond)
			return fmt.Sprintf("Query '%s' returned 42 results", input.Query), nil
		},
	)

	agent := aigentic.Agent{
		Model:        model,
		Name:         "ProductionAgent",
		Description:  "A production-ready agent with all safeguards",
		Instructions: "You are a data assistant. Help users query their data safely and efficiently.",
		Session:      session,
		AgentTools:   []aigentic.AgentTool{dataTool},
		Trace:        trace,
		Retries:      2,
		MaxLLMCalls:  10,
		LogLevel:     slog.LevelInfo,
	}

	fmt.Println("Executing production agent...")
	response, err := agent.Execute("Fetch data for user activity in the last 24 hours")

	if err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			log.Printf("ERROR: Agent execution timed out")
		} else if ctx.Err() == context.Canceled {
			log.Printf("ERROR: Agent execution was cancelled")
		} else {
			log.Printf("ERROR: Agent execution failed: %v", err)
		}
		return
	}

	fmt.Printf("Response: %s\n\n", response)

	fmt.Println("Production agent completed successfully!")
	traceDir := "/tmp/aigentic-traces"
	fmt.Printf("Trace available at: %s/trace-%s.txt\n", traceDir, trace.SessionID)
	fmt.Println()
	fmt.Println("Production features enabled:")
	fmt.Println("✓ Context timeout configured")
	fmt.Println("✓ Trace enabled for debugging")
	fmt.Println("✓ Retry logic configured")
	fmt.Println("✓ MaxLLMCalls limit set")
	fmt.Println("✓ Comprehensive error handling")
	fmt.Println("✓ Resource cleanup")

	session.Cancel()
	fmt.Println("\n✅ Example completed successfully!")
}
