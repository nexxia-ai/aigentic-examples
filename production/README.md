# Production-Ready Patterns Example

This example demonstrates production-ready patterns for building reliable, maintainable, and robust AI agents with aigentic. Learn the essential patterns that distinguish a proof-of-concept from a production system.

## What You'll Learn

- **Error Handling**: Handle LLM failures, tool errors, and timeouts gracefully
- **Trace & Debug**: Use traces to understand and debug agent behavior
- **Safety Limits**: Prevent runaway loops with MaxLLMCalls
- **Retry Logic**: Configure automatic retries for transient failures
- **Logging**: Adjust log levels for different environments
- **Cancellation**: Handle timeouts and graceful shutdown with context
- **Production Setup**: Combine all patterns into production-ready code

## Running the Example

```bash
# Set your OpenAI API key
export OPENAI_API_KEY=your_api_key_here

# Run with default (info) log level
go run main.go

# Run in development mode (verbose logging)
ENV=dev go run main.go

# Run in production mode (minimal logging)
ENV=prod go run main.go

# Or run from examples directory
go run github.com/nexxia-ai/aigentic-examples/production@latest
```

## Examples Demonstrated

### Example 1: Robust Error Handling

Learn to handle three types of errors:

1. **Tool Errors** - Graceful degradation when tools fail
2. **Critical Errors** - Proper propagation of fatal errors
3. **Transient Errors** - Recoverable errors with helpful messages

```go
agent := aigentic.Agent{
    Model:    model,
    Retries:  2, // Retry failed operations
    LogLevel: slog.LevelInfo,
}

response, err := agent.Execute("Check system status")
if err != nil {
    // Handle different error scenarios
    log.Printf("Agent failed: %v", err)
    return
}
```

### Example 2: Using Trace for Debugging

Traces capture complete execution history for debugging:

```go
trace := aigentic.NewTrace()

agent := aigentic.Agent{
    Model:    model,
    Trace:    trace, // Enable tracing
    LogLevel: slog.LevelDebug,
}

response, err := agent.Execute("Your query")

// Review trace file to see:
// - All LLM interactions
// - Tool calls and responses
// - Timing information
// - Context sent to model
fmt.Printf("Trace session ID: %s\n", trace.SessionID)
```

### Example 3: MaxLLMCalls Limit

Prevent infinite loops and control costs:

```go
agent := aigentic.Agent{
    Model:       model,
    MaxLLMCalls: 10, // Maximum 10 LLM calls per execution
}

// Agent will stop after 10 calls even if not finished
response, err := agent.Execute("Complex task")
```

**What it protects against:**
- Infinite tool call loops
- Excessive API costs
- Long-running processes
- Agent stuck states

### Example 4: Retry Configuration

Handle transient failures automatically:

```go
agent := aigentic.Agent{
    Model:   model,
    Retries: 3, // Retry up to 3 times on failure
}

// Agent will automatically retry on:
// - Network errors
// - Temporary service outages
// - Rate limiting
```

### Example 5: Log Level Configuration

Adjust logging verbosity by environment:

```go
// Development - verbose logging
devAgent := aigentic.Agent{
    LogLevel: slog.LevelDebug,
}

// Production - minimal logging
prodAgent := aigentic.Agent{
    LogLevel: slog.LevelWarn,
}
```

**Log Level Guide:**
- `Debug` - Development, debugging, troubleshooting
- `Info` - Default, general operational information
- `Warn` - Production, only warnings and errors
- `Error` - Critical issues only

### Example 6: Context Cancellation

Handle timeouts and graceful shutdown:

```go
// Create session with timeout
ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
defer cancel()

session := aigentic.NewSession(ctx)

agent := aigentic.Agent{
    Model:   model,
    Session: session,
}

run, err := agent.Start("Your query")
response, err := run.Wait(0)

if ctx.Err() == context.DeadlineExceeded {
    log.Printf("Agent timed out")
} else if ctx.Err() == context.Canceled {
    log.Printf("Agent was cancelled")
}

// Cleanup
session.Cancel()
```

**Use cases:**
- Request timeouts
- Graceful shutdown (SIGTERM)
- User cancellation (Ctrl+C)
- Resource cleanup

### Example 7: Comprehensive Production Setup

All patterns combined:

```go
// Setup
ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
defer cancel()

session := aigentic.NewSession(ctx)
trace := aigentic.NewTrace()

// Production-ready agent
agent := aigentic.Agent{
    Model:       model,
    Session:     session,
    Trace:       trace,
    Retries:     2,
    MaxLLMCalls: 10,
    LogLevel:    slog.LevelWarn,
    AgentTools:  []aigentic.AgentTool{productionTool},
}

// Execute with error handling
response, err := agent.Execute("Production query")
if err != nil {
    if ctx.Err() == context.DeadlineExceeded {
        log.Printf("ERROR: Timeout")
    } else if ctx.Err() == context.Canceled {
        log.Printf("ERROR: Cancelled")
    } else {
        log.Printf("ERROR: %v", err)
    }
    // Log to monitoring system
    return
}

// Success - log metrics
log.Printf("Success: session=%s", trace.SessionID)

// Cleanup
session.Cancel()
```

## Production Considerations

### Error Handling Strategy

**1. Tool Errors vs System Errors**
- Tool errors: Use `ToolResult` with `Error: true` for graceful handling
- System errors: Return Go errors for critical failures

```go
Execute: func(run *aigentic.AgentRun, args map[string]interface{}) (*ai.ToolResult, error) {
    // Graceful error - agent can continue
    if validationFailed {
        return &ai.ToolResult{
            Content: []ai.ToolContent{{
                Type:    "text",
                Content: "Validation failed: invalid input",
            }},
            Error: true,
        }, nil
    }

    // Critical error - agent should stop
    if criticalFailure {
        return nil, errors.New("database connection lost")
    }

    // Success
    return &ai.ToolResult{...}, nil
}
```

**2. Error Recovery**
- Use retries for transient errors
- Log all errors with context
- Degrade gracefully when possible
- Provide clear error messages to users

### Monitoring & Observability

**1. Traces**
- Enable traces in production for debugging
- Store traces in centralized logging system
- Use trace files to diagnose issues
- Review traces for prompt optimization

**2. Metrics to Track**
- Agent execution time
- LLM call count per execution
- Error rates by type
- Tool call success/failure rates
- Timeout occurrences

**3. Logging Best Practices**
```go
// Structured logging with context
log.Printf("Agent execution failed",
    "agent", agent.Name,
    "error", err,
    "duration", time.Since(start),
    "llm_calls", callCount,
)
```

### Performance Optimization

**1. Timeouts**
```go
// Set appropriate timeouts based on use case
// Quick queries: 10-30 seconds
// Complex tasks: 1-3 minutes
// Background jobs: 5-10 minutes

ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
defer cancel()
```

**2. Resource Limits**
```go
// Prevent resource exhaustion
agent := aigentic.Agent{
    MaxLLMCalls: 10, // Adjust based on expected complexity
}
```

**3. Connection Pooling**
- Reuse model instances across requests
- Implement connection pooling for APIs
- Cache frequently accessed data

### Cost Management

**1. Set Hard Limits**
```go
agent := aigentic.Agent{
    MaxLLMCalls: 5, // Control costs per execution
}
```

**2. Choose Appropriate Models**
```go
// Quick tasks - use smaller models
model := openai.NewModel("gpt-4o-mini", apiKey)

// Complex reasoning - use larger models
model := openai.NewModel("gpt-4o", apiKey)
```

**3. Monitor Usage**
- Track LLM API costs per agent
- Set budget alerts
- Review traces to optimize prompts

### Security Considerations

**1. API Key Management**
```go
// Never hardcode API keys
apiKey := os.Getenv("OPENAI_API_KEY")

// Use secrets management in production
// - AWS Secrets Manager
// - HashiCorp Vault
// - Kubernetes Secrets
```

**2. Input Validation**
```go
Execute: func(run *aigentic.AgentRun, args map[string]interface{}) (*ai.ToolResult, error) {
    // Always validate inputs
    param, ok := args["param"].(string)
    if !ok {
        return &ai.ToolResult{Error: true, ...}, nil
    }

    // Sanitize user input
    if !isValid(param) {
        return &ai.ToolResult{Error: true, ...}, nil
    }
}
```

**3. Rate Limiting**
- Implement rate limiting per user/tenant
- Protect against abuse
- Handle rate limit errors gracefully

### Deployment Best Practices

**1. Environment Configuration**
```bash
# Development
export ENV=dev
export LOG_LEVEL=debug
export TIMEOUT=60s

# Production
export ENV=prod
export LOG_LEVEL=warn
export TIMEOUT=30s
```

**2. Health Checks**
```go
func healthCheck() error {
    // Test model connectivity
    _, err := model.Ping()
    return err
}
```

**3. Graceful Shutdown**
```go
// Handle SIGTERM for graceful shutdown
sigChan := make(chan os.Signal, 1)
signal.Notify(sigChan, syscall.SIGTERM, syscall.SIGINT)

go func() {
    <-sigChan
    log.Println("Shutting down gracefully...")
    session.Cancel()
    os.Exit(0)
}()
```

## Debugging Techniques

### 1. Enable Verbose Logging
```bash
ENV=dev go run main.go
```

### 2. Review Trace Files
```bash
# Find trace files
ls -la /tmp/aigentic-traces/

# Review specific trace
cat /tmp/aigentic-traces/trace_20240101120000.txt
```

### 3. Test Error Scenarios
```go
// Simulate failures in development
if os.Getenv("ENV") == "dev" {
    // Inject test errors
}
```

### 4. Monitor LLM Calls
```go
agent := aigentic.Agent{
    LogLevel:    slog.LevelDebug, // See all LLM interactions
    MaxLLMCalls: 5,                // Limit for testing
}
```

## Common Patterns

### Pattern 1: Request Handler with Timeout
```go
func handleRequest(w http.ResponseWriter, r *http.Request) {
    ctx, cancel := context.WithTimeout(r.Context(), 30*time.Second)
    defer cancel()

    session := aigentic.NewSession(ctx)
    defer session.Cancel()

    agent := aigentic.Agent{
        Model:       model,
        Session:     session,
        MaxLLMCalls: 10,
        LogLevel:    slog.LevelWarn,
    }

    response, err := agent.Execute(r.FormValue("query"))
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    json.NewEncoder(w).Encode(map[string]string{
        "response": response,
    })
}
```

### Pattern 2: Background Job with Monitoring
```go
func processJob(job Job) error {
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
    defer cancel()

    session := aigentic.NewSession(ctx)
    defer session.Cancel()

    trace := aigentic.NewTrace()

    agent := aigentic.Agent{
        Model:       model,
        Session:     session,
        Trace:       trace,
        Retries:     3,
        MaxLLMCalls: 20,
        LogLevel:    slog.LevelInfo,
    }

    start := time.Now()
    response, err := agent.Execute(job.Query)
    duration := time.Since(start)

    // Log metrics
    metrics.RecordJobDuration(duration)
    metrics.RecordJobResult(err == nil)

    if err != nil {
        log.Printf("Job failed: session=%s error=%v", trace.SessionID, err)
        return err
    }

    log.Printf("Job succeeded: session=%s duration=%v", trace.SessionID, duration)
    return nil
}
```

### Pattern 3: Multi-Tenant System
```go
func createTenantAgent(tenantID string) aigentic.Agent {
    // Per-tenant configuration
    config := getTenantConfig(tenantID)

    return aigentic.Agent{
        Model:       model,
        MaxLLMCalls: config.MaxCalls,
        Retries:     config.Retries,
        LogLevel:    config.LogLevel,
    }
}
```

## Testing Production Patterns

### Unit Testing
```go
func TestAgentWithError(t *testing.T) {
    agent := aigentic.Agent{
        Model:    testModel,
        Retries:  2,
        LogLevel: slog.LevelDebug,
    }

    _, err := agent.Execute("test query")
    assert.NoError(t, err)
}
```

### Integration Testing
```go
func TestProductionSetup(t *testing.T) {
    ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
    defer cancel()

    session := aigentic.NewSession(ctx)
    defer session.Cancel()

    agent := aigentic.Agent{
        Model:       model,
        Session:     session,
        MaxLLMCalls: 5,
        Retries:     2,
    }

    response, err := agent.Execute("integration test query")
    assert.NoError(t, err)
    assert.NotEmpty(t, response)
}
```

### Load Testing
```go
func BenchmarkAgentExecution(b *testing.B) {
    agent := aigentic.Agent{
        Model:       model,
        MaxLLMCalls: 3,
    }

    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        _, err := agent.Execute("benchmark query")
        if err != nil {
            b.Fatal(err)
        }
    }
}
```

## Next Steps

- See [tools example](../tools) for production-ready tool patterns
- See [memory example](../memory) for persistent state management
- See [multi-agent example](../multi-agent) for distributed systems
- See [approval example](../approval) for human-in-the-loop workflows

## Production Checklist

Before deploying to production, ensure:

- [ ] Context timeouts configured appropriately
- [ ] Trace enabled for debugging
- [ ] Retry logic configured for transient failures
- [ ] MaxLLMCalls set to prevent runaway loops
- [ ] Log level set based on environment
- [ ] Error handling for all failure scenarios
- [ ] API keys stored securely (not hardcoded)
- [ ] Monitoring and alerting configured
- [ ] Health checks implemented
- [ ] Graceful shutdown handlers registered
- [ ] Resource cleanup (defer statements)
- [ ] Rate limiting implemented
- [ ] Cost tracking enabled
- [ ] Load testing completed
- [ ] Backup/rollback plan documented

## Support

For issues or questions:
- GitHub Issues: https://github.com/nexxia-ai/aigentic/issues
- Documentation: https://github.com/nexxia-ai/aigentic
- Examples: https://github.com/nexxia-ai/aigentic-examples
