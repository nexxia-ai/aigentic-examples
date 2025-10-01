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
calculatorTool := aigentic.AgentTool{
    Name:        "calculator",
    Description: "Performs mathematical calculations",
    InputSchema: map[string]interface{}{
        "type": "object",
        "properties": map[string]interface{}{
            "expression": map[string]interface{}{
                "type":        "string",
                "description": "Mathematical expression to evaluate",
            },
        },
        "required": []string{"expression"},
    },
    Execute: func(run *aigentic.AgentRun, args map[string]interface{}) (*ai.ToolResult, error) {
        expr := args["expression"].(string)
        result := evaluateExpression(expr)

        return &ai.ToolResult{
            Content: []ai.ToolContent{{
                Type:    "text",
                Content: fmt.Sprintf("Result: %v", result),
            }},
        }, nil
    },
}
```

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
