package main

import (
	"fmt"
	"log"
	"math"
	"os"
	"strconv"
	"strings"
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

func createCalculatorTool() aigentic.AgentTool {
	type CalculatorInput struct {
		Expression string `json:"expression" description:"Mathematical expression to evaluate (e.g., '2 + 2', '10 * 5', 'sqrt 16', '2 ^ 3')"`
	}

	return aigentic.NewTool(
		"calculator",
		"Performs basic mathematical calculations. Supports +, -, *, /, sqrt, and ^ (power) operations.",
		func(run *aigentic.AgentRun, input CalculatorInput) (string, error) {
			result, err := evaluateExpression(input.Expression)
			if err != nil {
				return "", fmt.Errorf("error evaluating expression: %v", err)
			}
			return fmt.Sprintf("Result: %v", result), nil
		},
	)
}

func evaluateExpression(expr string) (float64, error) {
	expr = strings.TrimSpace(expr)

	if strings.HasPrefix(expr, "sqrt") {
		numStr := strings.TrimSpace(strings.TrimPrefix(expr, "sqrt"))
		num, err := strconv.ParseFloat(numStr, 64)
		if err != nil {
			return 0, fmt.Errorf("invalid number for sqrt: %v", err)
		}
		return math.Sqrt(num), nil
	}

	for _, op := range []string{"+", "-", "*", "/", "^"} {
		if strings.Contains(expr, op) {
			parts := strings.Split(expr, op)
			if len(parts) != 2 {
				return 0, fmt.Errorf("invalid expression format")
			}

			left, err := strconv.ParseFloat(strings.TrimSpace(parts[0]), 64)
			if err != nil {
				return 0, fmt.Errorf("invalid left operand: %v", err)
			}

			right, err := strconv.ParseFloat(strings.TrimSpace(parts[1]), 64)
			if err != nil {
				return 0, fmt.Errorf("invalid right operand: %v", err)
			}

			switch op {
			case "+":
				return left + right, nil
			case "-":
				return left - right, nil
			case "*":
				return left * right, nil
			case "/":
				if right == 0 {
					return 0, fmt.Errorf("division by zero")
				}
				return left / right, nil
			case "^":
				return math.Pow(left, right), nil
			}
		}
	}

	return 0, fmt.Errorf("unsupported expression format")
}

func createTimeTool() aigentic.AgentTool {
	type TimeInput struct {
		Timezone string `json:"timezone" description:"IANA timezone name (e.g., 'America/New_York', 'Europe/London', 'Asia/Tokyo')"`
	}

	return aigentic.NewTool(
		"get_current_time",
		"Gets the current time in a specified timezone",
		func(run *aigentic.AgentRun, input TimeInput) (string, error) {
			loc, err := time.LoadLocation(input.Timezone)
			if err != nil {
				return "", fmt.Errorf("invalid timezone '%s'. Use IANA timezone names like 'America/New_York'", input.Timezone)
			}

			currentTime := time.Now().In(loc)
			timeStr := currentTime.Format("Monday, January 2, 2006 at 3:04 PM MST")

			return fmt.Sprintf("Current time in %s: %s", input.Timezone, timeStr), nil
		},
	)
}

func main() {
	utils.LoadEnvFile("../.env")

	fmt.Println("🛠️  Aigentic Tool Integration Example")
	fmt.Println("=====================================")
	fmt.Println()

	model := openai.NewModel("gpt-4o-mini", getAPIKey())

	agent := aigentic.Agent{
		Model:        model,
		Name:         "ToolAssistant",
		Description:  "An assistant with access to calculator and time tools",
		Instructions: "You are a helpful assistant with access to calculator and time tools. Use them to help users with calculations and time queries. Always use the appropriate tool when the user asks for something you can handle with a tool.",
		AgentTools: []aigentic.AgentTool{
			createCalculatorTool(),
			createTimeTool(),
		},
	}

	response, err := agent.Execute("What is 15 multiplied by 23, plus 100? Also, what time is it in New York?")
	if err != nil {
		log.Fatalf("Error: %v", err)
	}
	fmt.Printf("Response: %s\n\n", response)

	fmt.Println("✅ Example completed successfully!")
}
