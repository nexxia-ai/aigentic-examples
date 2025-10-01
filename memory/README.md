### Memory System Example

This example demonstrates aigentic's compartmentalized memory system for building agents that remember context across conversations and coordinate through shared state.

## What You'll Learn

- Using run memory for temporary task state
- Implementing session memory for persistent context
- Managing complex plans with plan memory
- Sharing memory across multi-agent systems
- Best practices for memory usage

## Memory Compartments

Aigentic provides three types of memory:

### 1. Run Memory
- **Scope**: Single agent execution (one `Execute()` or `Start()` call)
- **Persistence**: Cleared after agent run completes
- **Use Cases**: Task progress, intermediate calculations, temporary state
- **Access**: Automatically included in context (no tool call needed)

### 2. Session Memory
- **Scope**: All agent runs within the same session
- **Persistence**: Lasts for the session lifetime
- **Use Cases**: User preferences, conversation context, shared information
- **Access**: Must retrieve using `get_memory` tool

### 3. Plan Memory
- **Scope**: Complex multi-step plans
- **Persistence**: Configurable (run-level or session-level)
- **Use Cases**: Project plans, workflows, progress tracking
- **Access**: Must retrieve using `get_memory` tool

## Running the Example

```bash
# Set your OpenAI API key
export OPENAI_API_KEY=your_api_key_here

# Run from the examples directory
go run github.com/nexxia-ai/aigentic-examples/memory@latest

# Or run locally
cd memory
go run main.go
```

## Examples Demonstrated

### Example 1: Run Memory
Tracking multi-step task progress within a single agent execution.

### Example 2: Session Memory
Personal assistant that remembers user information across conversations.

### Example 3: Plan Memory
Project planner that creates, tracks, and updates complex plans.

### Example 4: Shared Memory in Multi-Agent Systems
Research and writing team coordinating through shared session memory.

## How to Use Memory

### Enable Memory

```go
agent := aigentic.Agent{
    Model:  model,
    Memory: memory.NewMemory(), // Enable memory system
}
```

### Agent Uses Memory Tools

The agent has access to three memory tools:

1. **save_memory** - Save information to a compartment
2. **get_memory** - Retrieve information from a compartment
3. **clear_memory** - Clear information from a compartment

The LLM decides when to use these tools based on your instructions.

### Guide the Agent

```go
Instructions: `
You are a personal assistant. Remember important information using session memory.
When a user shares preferences or personal details, save them with:
save_memory(compartment="session", content="...", category="preference")

To recall information later, use:
get_memory(compartment="session", category="preference")
`
```

## Compartment Selection Guide

| Scenario | Compartment | Why |
|----------|-------------|-----|
| Intermediate calculation results | Run | Temporary, not needed after task completes |
| User name and preferences | Session | Persist across conversations |
| Shopping cart contents | Session | Persist during shopping session |
| Multi-step project plan | Plan | Complex workflow that needs tracking |
| Current task step number | Run | Only relevant to current execution |
| User's language preference | Session | Applies to all interactions |

## Sharing Memory Across Agents

Create a shared session and memory instance:

```go
session := aigentic.NewSession(context.Background())
mem := memory.NewMemory()

agent1 := aigentic.Agent{
    Session: session,
    Memory:  mem,
}

agent2 := aigentic.Agent{
    Session: session,
    Memory:  mem,
}

// Both agents can access the same session memory
```

## Memory Configuration

### Custom Storage Location

```go
config := &memory.MemoryConfig{
    StoragePath: "/custom/path/memory.json",
}
mem := memory.NewMemoryWithConfig(config)
```

### Programmatic Access

```go
// Save entry directly (for testing or initialization)
entry := &memory.MemoryEntry{
    Content:  "User prefers dark mode",
    Category: "ui_preference",
    Priority: 5,
    Tags:     []string{"ui", "preference"},
}
err := mem.SaveEntry(memory.SessionMemory, entry)
```

## Best Practices

1. **Clear Instructions**: Tell the agent when and what to remember
2. **Use Categories**: Organize memory with categories for easy retrieval
3. **Appropriate Compartment**: Choose the right memory scope for your use case
4. **Session Management**: Create sessions for logical groupings of interactions
5. **Run Memory Auto-Clears**: Don't save important long-term data in run memory
6. **Test Memory Behavior**: Verify memory persistence across agent runs

## Common Patterns

### Pattern 1: Conversational Agent
```go
agent := aigentic.Agent{
    Session: session,
    Memory:  memory.NewMemory(),
    Instructions: "Remember user preferences and context using session memory",
}
```

### Pattern 2: Task Tracker
```go
agent := aigentic.Agent{
    Memory: memory.NewMemory(),
    Instructions: "Track task progress using run memory",
}
```

### Pattern 3: Project Manager
```go
agent := aigentic.Agent{
    Session: session,
    Memory:  memory.NewMemory(),
    Instructions: "Create and track plans using plan memory",
}
```

## Debugging Memory

```go
// Enable debug logging to see memory operations
agent := aigentic.Agent{
    LogLevel: slog.LevelDebug,
    Memory:   memory.NewMemory(),
}

// Memory operations will log:
// - When save_memory is called
// - When get_memory is called
// - What compartment is accessed
// - Memory content being saved/retrieved
```

## Next Steps

- See [multi-agent example](../multi-agent) for team coordination with shared memory
- See [production example](../production) for handling memory errors gracefully
- See [streaming example](../streaming) for real-time memory operations
