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
	"github.com/nexxia-ai/aigentic/ai"
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

// createCalculatorTool demonstrates a mathematical calculator tool
func createCalculatorTool() aigentic.AgentTool {
	return aigentic.AgentTool{
		Name:        "calculator",
		Description: "Performs basic mathematical calculations. Supports +, -, *, /, sqrt, and ^ (power) operations.",
		InputSchema: map[string]interface{}{
			"type": "object",
			"properties": map[string]interface{}{
				"expression": map[string]interface{}{
					"type":        "string",
					"description": "Mathematical expression to evaluate (e.g., '2 + 2', '10 * 5', 'sqrt 16', '2 ^ 3')",
				},
			},
			"required": []string{"expression"},
		},
		Execute: func(run *aigentic.AgentRun, args map[string]interface{}) (*ai.ToolResult, error) {
			expr, ok := args["expression"].(string)
			if !ok {
				return &ai.ToolResult{
					Content: []ai.ToolContent{{
						Type:    "text",
						Content: "Error: expression must be a string",
					}},
					Error: true,
				}, nil
			}

			result, err := evaluateExpression(expr)
			if err != nil {
				return &ai.ToolResult{
					Content: []ai.ToolContent{{
						Type:    "text",
						Content: fmt.Sprintf("Error evaluating expression: %v", err),
					}},
					Error: true,
				}, nil
			}

			return &ai.ToolResult{
				Content: []ai.ToolContent{{
					Type:    "text",
					Content: fmt.Sprintf("Result: %v", result),
				}},
			}, nil
		},
	}
}

// evaluateExpression is a simple expression evaluator
func evaluateExpression(expr string) (float64, error) {
	expr = strings.TrimSpace(expr)

	// Handle sqrt
	if strings.HasPrefix(expr, "sqrt") {
		numStr := strings.TrimSpace(strings.TrimPrefix(expr, "sqrt"))
		num, err := strconv.ParseFloat(numStr, 64)
		if err != nil {
			return 0, fmt.Errorf("invalid number for sqrt: %v", err)
		}
		return math.Sqrt(num), nil
	}

	// Handle basic operations
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

// createWeatherTool demonstrates a mock weather API tool
func createWeatherTool() aigentic.AgentTool {
	return aigentic.AgentTool{
		Name:        "get_weather",
		Description: "Gets the current weather for a specified city",
		InputSchema: map[string]interface{}{
			"type": "object",
			"properties": map[string]interface{}{
				"city": map[string]interface{}{
					"type":        "string",
					"description": "The city name to get weather for",
				},
				"units": map[string]interface{}{
					"type":        "string",
					"description": "Temperature units: 'celsius' or 'fahrenheit'",
					"enum":        []string{"celsius", "fahrenheit"},
				},
			},
			"required": []string{"city"},
		},
		Execute: func(run *aigentic.AgentRun, args map[string]interface{}) (*ai.ToolResult, error) {
			city, ok := args["city"].(string)
			if !ok {
				return &ai.ToolResult{
					Content: []ai.ToolContent{{
						Type:    "text",
						Content: "Error: city must be a string",
					}},
					Error: true,
				}, nil
			}

			units := "celsius"
			if u, ok := args["units"].(string); ok {
				units = u
			}

			// Mock weather data
			weather := mockWeatherData(city, units)

			return &ai.ToolResult{
				Content: []ai.ToolContent{{
					Type:    "text",
					Content: weather,
				}},
			}, nil
		},
	}
}

func mockWeatherData(city, units string) string {
	// Simple mock data based on city name hash
	temp := 20 + (len(city) % 15)
	if units == "fahrenheit" {
		temp = (temp * 9 / 5) + 32
	}

	conditions := []string{"Sunny", "Cloudy", "Rainy", "Partly Cloudy"}
	condition := conditions[len(city)%len(conditions)]

	unit := "°C"
	if units == "fahrenheit" {
		unit = "°F"
	}

	return fmt.Sprintf("Current weather in %s: %s, %d%s", city, condition, temp, unit)
}

// createTimeTool demonstrates a time utility tool
func createTimeTool() aigentic.AgentTool {
	return aigentic.AgentTool{
		Name:        "get_current_time",
		Description: "Gets the current time in a specified timezone",
		InputSchema: map[string]interface{}{
			"type": "object",
			"properties": map[string]interface{}{
				"timezone": map[string]interface{}{
					"type":        "string",
					"description": "IANA timezone name (e.g., 'America/New_York', 'Europe/London', 'Asia/Tokyo')",
				},
			},
			"required": []string{"timezone"},
		},
		Execute: func(run *aigentic.AgentRun, args map[string]interface{}) (*ai.ToolResult, error) {
			timezone, ok := args["timezone"].(string)
			if !ok {
				return &ai.ToolResult{
					Content: []ai.ToolContent{{
						Type:    "text",
						Content: "Error: timezone must be a string",
					}},
					Error: true,
				}, nil
			}

			loc, err := time.LoadLocation(timezone)
			if err != nil {
				return &ai.ToolResult{
					Content: []ai.ToolContent{{
						Type:    "text",
						Content: fmt.Sprintf("Error: invalid timezone '%s'. Use IANA timezone names like 'America/New_York'", timezone),
					}},
					Error: true,
				}, nil
			}

			currentTime := time.Now().In(loc)
			timeStr := currentTime.Format("Monday, January 2, 2006 at 3:04 PM MST")

			return &ai.ToolResult{
				Content: []ai.ToolContent{{
					Type:    "text",
					Content: fmt.Sprintf("Current time in %s: %s", timezone, timeStr),
				}},
			}, nil
		},
	}
}

func main() {
	utils.LoadEnvFile("../.env")

	fmt.Println("🛠️  Aigentic Tool Integration Examples")
	fmt.Println("======================================")
	fmt.Println()

	model := openai.NewModel("gpt-4o-mini", getAPIKey())

	// Agent with multiple tools
	agent := aigentic.Agent{
		Model:        model,
		Name:         "ToolAssistant",
		Description:  "An assistant with access to calculator, weather, and time tools",
		Instructions: "You are a helpful assistant with access to various tools. Use them to help users with calculations, weather information, and time queries. Always use the appropriate tool when the user asks for something you can handle with a tool.",
		AgentTools: []aigentic.AgentTool{
			createCalculatorTool(),
			createWeatherTool(),
			createTimeTool(),
		},
	}

	// Example 1: Calculator
	fmt.Println("=== Example 1: Calculator Tool ===")
	response, err := agent.Execute("What is 15 multiplied by 23, plus 100?")
	if err != nil {
		log.Fatalf("Error: %v", err)
	}
	fmt.Printf("Response: %s\n\n", response)

	// Example 2: Weather
	fmt.Println("=== Example 2: Weather Tool ===")
	response, err = agent.Execute("What's the weather like in Tokyo and London? Show both in celsius.")
	if err != nil {
		log.Fatalf("Error: %v", err)
	}
	fmt.Printf("Response: %s\n\n", response)

	// Example 3: Time
	fmt.Println("=== Example 3: Time Tool ===")
	response, err = agent.Execute("What time is it right now in New York and Sydney?")
	if err != nil {
		log.Fatalf("Error: %v", err)
	}
	fmt.Printf("Response: %s\n\n", response)

	// Example 4: Multiple tools in one query
	fmt.Println("=== Example 4: Using Multiple Tools ===")
	response, err = agent.Execute("Calculate the square root of 144, tell me the weather in Paris, and what time it is in Tokyo.")
	if err != nil {
		log.Fatalf("Error: %v", err)
	}
	fmt.Printf("Response: %s\n\n", response)

	fmt.Println("✅ All tool examples completed successfully!")
}
