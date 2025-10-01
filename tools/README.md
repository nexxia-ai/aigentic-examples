# Tool Integration Example

This example demonstrates how to create and use custom tools with aigentic agents.

## What You'll Learn

- Creating custom tools with schema definitions
- Implementing tool execution logic
- Handling tool parameters and validation
- Using multiple tools in a single agent
- Error handling in tools

## Tools Demonstrated

1. **Calculator Tool** - Performs mathematical operations
2. **Weather Tool** - Retrieves weather information (mock data)
3. **Time Tool** - Gets current time in different timezones

## Running the Example

```bash
# Set your OpenAI API key
export OPENAI_API_KEY=your_api_key_here

# Run from the examples directory
go run github.com/nexxia-ai/aigentic-examples/tools@latest

# Or run locally
cd tools
go run main.go
```

## Key Concepts

### Tool Structure

Each tool consists of:
- **Name**: Unique identifier for the tool
- **Description**: What the tool does (helps the LLM decide when to use it)
- **InputSchema**: JSON schema defining the tool's parameters
- **Execute**: Function that implements the tool's logic

### Example Tool Definition

```go
// Define input struct with JSON tags for automatic schema generation
type CalculatorInput struct {
    Expression string `json:"expression" description:"Mathematical expression to evaluate"`
}

// Create tool using aigentic.NewTool for type safety and automatic schema generation
calculatorTool := aigentic.NewTool(
    "calculator",
    "Performs mathematical calculations",
    func(run *aigentic.AgentRun, input CalculatorInput) (string, error) {
        result, err := evaluateExpression(input.Expression)
        if err != nil {
            return "", fmt.Errorf("error evaluating expression: %v", err)
        }
        return fmt.Sprintf("Result: %v", result), nil
    },
)
```

### Benefits of the New Approach

- **Type Safety**: Input parameters are strongly typed
- **Automatic Schema Generation**: JSON schema is generated from struct tags
- **Cleaner Code**: No manual schema definitions or type assertions
- **Better Error Handling**: Return errors directly instead of wrapping in ToolResult

## Best Practices

1. **Clear Descriptions**: Write detailed descriptions so the LLM knows when to use the tool
2. **Schema Validation**: Define comprehensive input schemas with types and constraints
3. **Error Handling**: Return meaningful error messages in ToolResult with `Error: true`
4. **Type Safety**: Validate parameter types before using them
5. **Documentation**: Include examples of valid inputs in descriptions

## Next Steps

- See [approval example](../approval) for adding human approval to tools
- See [production example](../production) for robust error handling patterns
- See [mcp example](../mcp) for Model Context Protocol integration
