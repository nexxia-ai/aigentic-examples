package main

import (
	"context"
	"errors"
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

// getLogLevel determines the log level based on environment
func getLogLevel() slog.Level {
	env := os.Getenv("ENV")
	switch env {
	case "production", "prod":
		return slog.LevelWarn // Only warnings and errors in production
	case "development", "dev":
		return slog.LevelDebug // Verbose logging for development
	default:
		return slog.LevelInfo // Default to info level
	}
}

// Example 1: Robust Error Handling
// Demonstrates handling LLM failures, tool errors, and timeouts
func example1RobustErrorHandling() {
	fmt.Println("=== Example 1: Robust Error Handling ===")
	fmt.Println("Demonstrating various error scenarios and recovery patterns")
	fmt.Println()

	apiKey := getAPIKey()
	model := openai.NewModel("gpt-4o-mini", apiKey)

	// Create a tool that can fail in different ways
	type CheckStatusInput struct {
		Service string `json:"service" description:"Service name to check"`
	}

	unreliableTool := aigentic.NewTool(
		"check_status",
		"Checks system status (may fail)",
		func(run *aigentic.AgentRun, input CheckStatusInput) (string, error) {
			// Simulate different error scenarios
			if input.Service == "database" {
				// Return error that will be handled gracefully
				return "", errors.New("database connection failed: timeout after 5s. This is a recoverable error")
			}

			if input.Service == "api" {
				// Critical error - return as Go error
				return "", errors.New("API service check failed: internal error")
			}

			// Success case
			return fmt.Sprintf("%s service is healthy", input.Service), nil
		},
	)

	agent := aigentic.Agent{
		Model:        model,
		Name:         "ErrorHandlingAgent",
		Description:  "An agent that demonstrates error handling",
		Instructions: "Check the status of multiple services and report back. If a service check fails, try to provide a helpful message about what went wrong.",
		AgentTools:   []aigentic.AgentTool{unreliableTool},
		Retries:      2, // Retry up to 2 times on failure
		LogLevel:     slog.LevelInfo,
	}

	// Test with a service that will fail gracefully
	fmt.Println("Testing with recoverable error (database):")
	response, err := agent.Execute("Check the status of the database service")
	if err != nil {
		fmt.Printf("Agent execution failed: %v\n", err)
	} else {
		fmt.Printf("Response: %s\n", response)
	}
	fmt.Println()

	// Test with a service that succeeds
	fmt.Println("Testing with successful service (cache):")
	response, err = agent.Execute("Check the status of the cache service")
	if err != nil {
		fmt.Printf("Agent execution failed: %v\n", err)
	} else {
		fmt.Printf("Response: %s\n", response)
	}
	fmt.Println()

	// Test with a critical error scenario
	fmt.Println("Testing with critical error (api):")
	response, err = agent.Execute("Check the status of the api service")
	if err != nil {
		fmt.Printf("Agent execution failed (as expected): %v\n", err)
	} else {
		fmt.Printf("Response: %s\n", response)
	}
	fmt.Println()
}

// Example 2: Using Trace for Debugging
// Demonstrates how to use traces to debug agent behavior
func example2TraceDebugging() {
	fmt.Println("=== Example 2: Using Trace for Debugging ===")
	fmt.Println("Traces help you understand what your agent is doing")
	fmt.Println()

	apiKey := getAPIKey()
	model := openai.NewModel("gpt-4o-mini", apiKey)

	// Create trace with custom configuration
	trace := aigentic.NewTrace()

	agent := aigentic.Agent{
		Model:        model,
		Name:         "TracedAgent",
		Description:  "An agent with tracing enabled",
		Instructions: "You are a helpful assistant. Explain quantum computing in simple terms.",
		Trace:        trace, // Enable tracing
		LogLevel:     slog.LevelDebug,
	}

	response, err := agent.Execute("Explain quantum computing in 2 sentences")
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	fmt.Printf("Response: %s\n\n", response)

	// Trace files are automatically saved to temp directory
	traceDir := "/tmp/aigentic-traces" // Default trace directory
	fmt.Printf("Trace saved to: %s/trace-%s.txt\n", traceDir, trace.SessionID)
	fmt.Println("Review the trace file to see:")
	fmt.Println("- All LLM interactions")
	fmt.Println("- Tool calls and responses")
	fmt.Println("- Timing information")
	fmt.Println("- Context sent to the model")
	fmt.Println()
}

// Example 3: MaxLLMCalls to Prevent Runaway Loops
// Demonstrates protecting against infinite agent loops
func example3MaxLLMCalls() {
	fmt.Println("=== Example 3: MaxLLMCalls Limit ===")
	fmt.Println("Preventing runaway agent loops with MaxLLMCalls")
	fmt.Println()

	apiKey := getAPIKey()
	model := openai.NewModel("gpt-4o-mini", apiKey)

	// Create a tool that encourages repeated calls
	recursiveTool := aigentic.NewTool(
		"generate_number",
		"Generates a random-looking number",
		func(run *aigentic.AgentRun, input struct{}) (string, error) {
			// Each call returns a different result, potentially causing repeated calls
			return fmt.Sprintf("Generated number: %d", time.Now().Unix()%100), nil
		},
	)

	agent := aigentic.Agent{
		Model:        model,
		Name:         "LimitedAgent",
		Description:  "An agent with LLM call limits",
		Instructions: "Generate a number using the tool and report it.",
		AgentTools:   []aigentic.AgentTool{recursiveTool},
		MaxLLMCalls:  3, // Limit to 3 LLM calls maximum
		LogLevel:     slog.LevelInfo,
	}

	response, err := agent.Execute("Generate a number for me")
	if err != nil {
		fmt.Printf("Agent stopped due to limits: %v\n", err)
	} else {
		fmt.Printf("Response: %s\n", response)
	}
	fmt.Println()
	fmt.Println("Note: MaxLLMCalls protects against:")
	fmt.Println("- Infinite tool call loops")
	fmt.Println("- Excessive API costs")
	fmt.Println("- Long-running agent processes")
	fmt.Println()
}

// Example 4: Retries Configuration
// Demonstrates automatic retry behavior on failures
func example4Retries() {
	fmt.Println("=== Example 4: Retry Configuration ===")
	fmt.Println("Automatic retries for transient failures")
	fmt.Println()

	apiKey := getAPIKey()
	model := openai.NewModel("gpt-4o-mini", apiKey)

	callCount := 0
	flakyTool := aigentic.NewTool(
		"flaky_service",
		"A service that fails sometimes but eventually succeeds",
		func(run *aigentic.AgentRun, input struct{}) (string, error) {
			callCount++
			fmt.Printf("  [Tool call #%d]\n", callCount)

			// Fail on first call, succeed on retry
			if callCount == 1 {
				return "", errors.New("temporary network error")
			}

			return "Service responded successfully!", nil
		},
	)

	agent := aigentic.Agent{
		Model:        model,
		Name:         "ResilientAgent",
		Description:  "An agent with retry logic",
		Instructions: "Call the flaky service and report the result.",
		AgentTools:   []aigentic.AgentTool{flakyTool},
		Retries:      3, // Retry up to 3 times
		LogLevel:     slog.LevelInfo,
	}

	response, err := agent.Execute("Call the flaky service")
	if err != nil {
		fmt.Printf("Failed after retries: %v\n", err)
	} else {
		fmt.Printf("Response: %s\n", response)
	}
	fmt.Println()
	fmt.Println("Retry strategies help with:")
	fmt.Println("- Transient network errors")
	fmt.Println("- Rate limiting (with backoff)")
	fmt.Println("- Temporary service outages")
	fmt.Println()
}

// Example 5: LogLevel Configuration
// Demonstrates different log levels for different environments
func example5LogLevels() {
	fmt.Println("=== Example 5: Log Level Configuration ===")
	fmt.Println("Adjust logging verbosity for different environments")
	fmt.Println()

	apiKey := getAPIKey()
	model := openai.NewModel("gpt-4o-mini", apiKey)

	currentLogLevel := getLogLevel()
	fmt.Printf("Current environment: %s\n", os.Getenv("ENV"))
	fmt.Printf("Log level: %v\n\n", currentLogLevel)

	// Development agent - verbose logging
	devAgent := aigentic.Agent{
		Model:        model,
		Name:         "DevAgent",
		Description:  "Agent with development logging",
		Instructions: "You are a helpful assistant.",
		LogLevel:     slog.LevelDebug, // Shows all details
	}

	// Production agent - minimal logging
	prodAgent := aigentic.Agent{
		Model:        model,
		Name:         "ProdAgent",
		Description:  "Agent with production logging",
		Instructions: "You are a helpful assistant.",
		LogLevel:     slog.LevelWarn, // Only warnings and errors
	}

	// Use the appropriate agent based on environment
	agent := devAgent
	if currentLogLevel == slog.LevelWarn {
		agent = prodAgent
	}

	response, err := agent.Execute("What is 2+2?")
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	fmt.Printf("Response: %s\n\n", response)
	fmt.Println("Log levels guide:")
	fmt.Println("- Debug: Development, debugging, troubleshooting")
	fmt.Println("- Info: Default, general operational info")
	fmt.Println("- Warn: Production, only warnings and errors")
	fmt.Println("- Error: Critical issues only")
	fmt.Println()
}

// Example 6: Context Cancellation and Graceful Shutdown
// Demonstrates handling timeouts and cancellations
func example6ContextCancellation() {
	fmt.Println("=== Example 6: Context Cancellation ===")
	fmt.Println("Graceful shutdown and timeout handling")
	fmt.Println()

	apiKey := getAPIKey()
	model := openai.NewModel("gpt-4o-mini", apiKey)

	// Create a session with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	session := aigentic.NewSession(ctx)

	agent := aigentic.Agent{
		Model:        model,
		Name:         "TimeoutAgent",
		Description:  "Agent with timeout handling",
		Instructions: "You are a helpful assistant.",
		Session:      session,
		LogLevel:     slog.LevelInfo,
	}

	// Start the agent
	run, err := agent.Start("Explain machine learning briefly")
	if err != nil {
		fmt.Printf("Failed to start agent: %v\n", err)
		return
	}

	// Wait for completion with timeout monitoring
	fmt.Println("Agent running with 30s timeout...")
	response, err := run.Wait(0)
	if err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			fmt.Printf("Agent timed out: %v\n", err)
		} else if ctx.Err() == context.Canceled {
			fmt.Printf("Agent was cancelled: %v\n", err)
		} else {
			fmt.Printf("Agent error: %v\n", err)
		}
		return
	}

	fmt.Printf("Response: %s\n\n", response)

	fmt.Println("Context cancellation enables:")
	fmt.Println("- Request timeouts")
	fmt.Println("- Graceful shutdown on SIGTERM")
	fmt.Println("- User cancellation (Ctrl+C)")
	fmt.Println("- Resource cleanup")
	fmt.Println()

	// Clean up
	session.Cancel()
}

// Example 7: Comprehensive Production Setup
// Combines all best practices into one production-ready example
func example7ComprehensiveSetup() {
	fmt.Println("=== Example 7: Comprehensive Production Setup ===")
	fmt.Println("Combining all production best practices")
	fmt.Println()

	apiKey := getAPIKey()
	model := openai.NewModel("gpt-4o-mini", apiKey)

	// Create production-grade session
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	session := aigentic.NewSession(ctx)
	trace := aigentic.NewTrace()

	// Production-ready tool with comprehensive error handling
	type FetchDataInput struct {
		Query string `json:"query" description:"The data query to execute"`
	}

	dataTool := aigentic.NewTool(
		"fetch_data",
		"Fetches data from a database or API",
		func(run *aigentic.AgentRun, input FetchDataInput) (string, error) {
			// Check context for cancellation
			if run.Session().Context.Err() != nil {
				return "", run.Session().Context.Err()
			}

			// Simulate data fetch
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
		LogLevel:     getLogLevel(),
	}

	// Execute with comprehensive error handling
	fmt.Println("Executing production agent...")
	response, err := agent.Execute("Fetch data for user activity in the last 24 hours")

	if err != nil {
		// Handle different error types
		if ctx.Err() == context.DeadlineExceeded {
			log.Printf("ERROR: Agent execution timed out")
		} else if ctx.Err() == context.Canceled {
			log.Printf("ERROR: Agent execution was cancelled")
		} else {
			log.Printf("ERROR: Agent execution failed: %v", err)
		}
		// In production: log to monitoring system, increment error metrics
		return
	}

	fmt.Printf("Response: %s\n\n", response)

	// Log success metrics
	fmt.Println("Production agent completed successfully!")
	traceDir := "/tmp/aigentic-traces"
	fmt.Printf("Trace available at: %s/trace-%s.txt\n", traceDir, trace.SessionID)
	fmt.Println()
	fmt.Println("Production checklist:")
	fmt.Println("✓ Context timeout configured")
	fmt.Println("✓ Trace enabled for debugging")
	fmt.Println("✓ Retry logic configured")
	fmt.Println("✓ MaxLLMCalls limit set")
	fmt.Println("✓ Environment-based log levels")
	fmt.Println("✓ Comprehensive error handling")
	fmt.Println("✓ Resource cleanup (defer cancel)")
	fmt.Println()

	// Cleanup
	session.Cancel()
}

func main() {
	utils.LoadEnvFile("../.env")

	fmt.Println("==================================================")
	fmt.Println("  Aigentic Production-Ready Patterns")
	fmt.Println("==================================================")
	fmt.Println()
	fmt.Println("This example demonstrates production-ready patterns")
	fmt.Println("for building reliable, maintainable AI agents.")
	fmt.Println()

	// Run all examples
	example1RobustErrorHandling()
	fmt.Println("--------------------------------------------------")
	fmt.Println()

	example2TraceDebugging()
	fmt.Println("--------------------------------------------------")
	fmt.Println()

	example3MaxLLMCalls()
	fmt.Println("--------------------------------------------------")
	fmt.Println()

	example4Retries()
	fmt.Println("--------------------------------------------------")
	fmt.Println()

	example5LogLevels()
	fmt.Println("--------------------------------------------------")
	fmt.Println()

	example6ContextCancellation()
	fmt.Println("--------------------------------------------------")
	fmt.Println()

	example7ComprehensiveSetup()
	fmt.Println("--------------------------------------------------")
	fmt.Println()

	fmt.Println("✅ All production examples completed!")
	fmt.Println()
	fmt.Println("Next steps:")
	fmt.Println("- Review the trace files in your temp directory")
	fmt.Println("- Experiment with different ENV values (dev, prod)")
	fmt.Println("- Try adjusting timeout durations")
	fmt.Println("- Integrate with your monitoring system")
}
