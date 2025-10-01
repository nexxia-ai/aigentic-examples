# Multi-Agent System Example

This example demonstrates how to build multi-agent systems where agents coordinate and delegate tasks to specialized sub-agents.

## What You'll Learn

- Creating agent teams with specialized roles
- Delegating tasks to specialized agents
- Coordinating workflows across agent teams
- Managing agent-to-agent communication

## Example Demonstrated

### Research and Writing Team
A project manager coordinates a researcher and writer to produce articles.

**Pattern**: Linear delegation (Research → Write)

## Running the Example

```bash
# Set your OpenAI API key
export OPENAI_API_KEY=your_api_key_here

# Run from the examples directory
go run github.com/nexxia-ai/aigentic-examples/multi-agent@latest

# Or run locally
cd multi-agent
go run main.go
```

## How It Works

### Agent Composition

Sub-agents are added to a parent agent via the `Agents` field:

```go
coordinator := aigentic.Agent{
    Name:   "Coordinator",
    Model:  model,
    Agents: []aigentic.Agent{
        researchAgent,
        writerAgent,
    },
}
```

### Automatic Tool Conversion

Each sub-agent is automatically exposed as a tool to its parent:

```go
// The parent agent can now "call" the sub-agents
// The framework handles:
// 1. Converting sub-agents to tools
// 2. Routing requests to the appropriate sub-agent
// 3. Returning results to the parent
// 4. Maintaining event streams for monitoring
```

### Sub-Agent Execution

When a parent agent invokes a sub-agent:
1. Parent decides to delegate based on instructions
2. Framework creates a new run for the sub-agent
3. Sub-agent executes with its own context
4. Results return to parent as tool output
5. Parent continues with the information

## Design Patterns

### Pattern 1: Sequential Workflow
```
Manager → Agent A → Agent B → Result
```

Use when: Tasks must happen in order

### Pattern 2: Parallel Consultation
```
        ┌→ Agent A →┐
Manager →  Agent B → Synthesize → Result
        └→ Agent C →┘
```

Use when: Need multiple expert perspectives

### Pattern 3: Hierarchical
```
          CEO
        /     \
    Dept A   Dept B
    /   \    /   \
Team1 Team2 Team3 Team4
```

Use when: Complex organizational structures

## Best Practices

1. **Clear Responsibilities**: Each agent should have a specific, well-defined role
2. **Appropriate Delegation**: Parent instructions should guide when to delegate
3. **Avoid Deep Nesting**: More than 3 levels can be hard to debug and manage
4. **Share Sessions**: Use `Session` to share context across agent runs
5. **Monitor Events**: Use event streams to track multi-agent interactions
6. **Set LogLevel**: Use `slog.LevelInfo` or `LevelDebug` to trace agent calls

## Common Use Cases

- **Content Creation**: Research → Write → Edit → Publish pipeline
- **Decision Making**: Multiple experts provide analysis → Synthesis
- **Customer Service**: Triage → Specialized handler → Resolution
- **Data Processing**: Extract → Transform → Analyze → Report
- **Software Development**: Plan → Code → Review → Test

## Debugging Tips

```go
// Enable detailed logging to see agent interactions
agent := aigentic.Agent{
    LogLevel: slog.LevelDebug,
    Trace:    aigentic.NewTrace(), // Create trace files
}

// Monitor events in real-time
run, _ := agent.Start("task")
for event := range run.Next() {
    switch e := event.(type) {
    case *aigentic.ToolEvent:
        fmt.Printf("Sub-agent called: %s\n", e.ToolName)
    case *aigentic.ToolResponseEvent:
        fmt.Printf("Sub-agent response: %s\n", e.Content)
    }
}
```

## Next Steps

- See [memory example](../memory) for sharing context across agents
- See [production example](../production) for error handling in multi-agent systems
- See [streaming example](../streaming) for real-time multi-agent monitoring
