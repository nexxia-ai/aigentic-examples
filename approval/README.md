# Human-in-the-Loop Approval Example

This example demonstrates how to implement human-in-the-loop approval workflows for sensitive operations in aigentic. When agents need to perform critical actions like sending emails, deleting files, or transferring money, you can require explicit human approval before execution.

## What You'll Learn

- How to enable approval for specific tools using `RequireApproval`
- Handling `ApprovalEvent` in the agent event stream
- Implementing approval UI/workflows
- Using validation functions to provide context for approval decisions
- Combining tools that require approval with tools that don't
- Timeout handling and rejection flows
- Best practices for building production approval systems

## Examples Demonstrated

### 1. Simple Email Approval
Shows basic approval workflow for sending emails. The agent requests to send an email, and the human must approve before it's sent.

**Use Case**: Preventing accidental or unauthorized email communications

### 2. File Deletion Approval
Demonstrates approval for destructive operations that can't be easily undone.

**Use Case**: Protecting critical files from accidental deletion

### 3. Financial Transaction with Validation
Shows advanced approval with custom validation logic that provides warnings for large amounts.

**Use Case**: Preventing unauthorized or erroneous financial transactions

### 4. Mixed Tools with Selective Approval
Demonstrates combining tools that require approval (transfers) with tools that don't (read-only queries).

**Use Case**: Balancing security with efficiency - only approve risky operations

## Running the Example

```bash
# Set your OpenAI API key
export OPENAI_API_KEY=your_api_key_here

# Run from the examples directory
go run github.com/nexxia-ai/aigentic-examples/approval@latest

# Or run locally (interactive mode)
cd approval
go run main.go

# Run in automated mode (auto-approves all requests)
go run main.go --auto
```

## How It Works

### 1. Enable Approval on Tools

Set `RequireApproval: true` on any tool that needs human oversight:

```go
// Define input struct with JSON tags
type SendEmailInput struct {
    To      string `json:"to" description:"Email recipient"`
    Subject string `json:"subject" description:"Email subject"`
    Body    string `json:"body" description:"Email body"`
}

// Create tool using aigentic.NewTool
emailTool := aigentic.NewTool(
    "send_email",
    "Sends an email to a recipient",
    func(run *aigentic.AgentRun, input SendEmailInput) (string, error) {
        // Tool implementation...
        return "Email sent successfully", nil
    },
)
emailTool.RequireApproval = true // Enable approval
```

### 2. Handle Approval Events

When a tool requires approval, an `ApprovalEvent` is sent to your event handler:

```go
run, err := agent.Start("Send an email to john@example.com")
if err != nil {
    log.Fatal(err)
}

for event := range run.Next() {
    switch e := event.(type) {
    case *aigentic.ApprovalEvent:
        // Present approval UI to user
        fmt.Printf("Approve %s? (y/n): ", e.ToolName)

        // Get user decision
        var response string
        fmt.Scanln(&response)
        approved := response == "y"

        // Send approval decision back to agent
        run.Approve(e.ApprovalID, approved)

    case *aigentic.ContentEvent:
        fmt.Print(e.Content)

    case *aigentic.ToolEvent:
        fmt.Printf("\n[Tool executed: %s]\n", e.ToolName)

    case *aigentic.ErrorEvent:
        log.Printf("Error: %v", e.Err)
    }
}
```

### 3. ApprovalEvent Structure

The `ApprovalEvent` contains important information for making approval decisions:

```go
type ApprovalEvent struct {
    RunID            string            // Unique run identifier
    ApprovalID       string            // Unique approval request ID
    ToolName         string            // Name of tool requesting approval
    ValidationResult ValidationResult  // Results from validation function
}

type ValidationResult struct {
    Values           any      // Tool parameters (usually map[string]interface{})
    Message          string   // Custom validation message
    ValidationErrors []error  // Any validation errors
}
```

### 4. Custom Validation Logic

Add validation logic to provide context and warnings for approval decisions:

```go
transferTool := aigentic.AgentTool{
    Name:            "transfer_money",
    RequireApproval: true,
    // ... schema definition ...
    Validate: func(run *aigentic.AgentRun, args map[string]interface{}) (aigentic.ValidationResult, error) {
        amount := args["amount"].(float64)

        var message string
        if amount > 10000 {
            message = fmt.Sprintf("WARNING: Large transaction: $%.2f", amount)
        } else if amount > 1000 {
            message = fmt.Sprintf("CAUTION: Moderate transaction: $%.2f", amount)
        } else {
            message = fmt.Sprintf("Transaction amount: $%.2f", amount)
        }

        return aigentic.ValidationResult{
            Values:  args,
            Message: message,
        }, nil
    },
    Execute: func(run *aigentic.AgentRun, args map[string]interface{}) (*ai.ToolResult, error) {
        // Perform the actual transfer...
        return &ai.ToolResult{
            Content: []ai.ToolContent{{
                Type:    "text",
                Content: "Transfer completed successfully",
            }},
        }, nil
    },
}
```

## Approval Workflows

### Simple Approval Flow

```
Agent decides to use tool
        ↓
Validation runs (if defined)
        ↓
ApprovalEvent sent to handler
        ↓
Human reviews and approves/rejects
        ↓
run.Approve(id, decision) called
        ↓
If approved: Tool executes
If rejected: Error returned to agent
        ↓
Agent continues with result
```

### Timeout Handling

By default, approval requests timeout after 5 minutes. When a timeout occurs:

```go
case *aigentic.ErrorEvent:
    if strings.Contains(e.Err.Error(), "approval timeout") {
        fmt.Println("Approval request timed out")
        // Handle timeout (e.g., notify user, log event)
    }
```

### Rejection Flow

When an approval is rejected:

```go
run.Approve(e.ApprovalID, false) // Reject the approval
```

The agent receives an error indicating the operation was denied and can adjust accordingly.

## When to Use Approval Workflows

### Critical Operations
- Financial transactions
- Data deletion or modification
- Email/message sending
- Account changes
- System configuration updates

### Compliance Requirements
- Operations requiring audit trails
- Actions needing dual authorization
- Regulatory compliance scenarios
- Privacy-sensitive operations

### Risk Mitigation
- Preventing AI hallucination impact
- Protecting against prompt injection
- Safeguarding production systems
- Controlling external API costs

## Building Production Approval UIs

### Web-Based Approval

```go
// Store pending approvals
type PendingApproval struct {
    ApprovalID string
    ToolName   string
    Parameters map[string]interface{}
    Timestamp  time.Time
}

var pendingApprovals sync.Map

// In your event handler
case *aigentic.ApprovalEvent:
    approval := PendingApproval{
        ApprovalID: e.ApprovalID,
        ToolName:   e.ToolName,
        Parameters: e.ValidationResult.Values.(map[string]interface{}),
        Timestamp:  time.Now(),
    }
    pendingApprovals.Store(e.ApprovalID, approval)

    // Send notification to web UI
    notifyWebUI(approval)

// In your web handler
func handleApprovalDecision(w http.ResponseWriter, r *http.Request) {
    approvalID := r.FormValue("approval_id")
    approved := r.FormValue("decision") == "approve"

    if val, ok := pendingApprovals.Load(approvalID); ok {
        run := val.(PendingApproval)
        run.Approve(approvalID, approved)
        pendingApprovals.Delete(approvalID)
    }
}
```

### Slack/Teams Integration

```go
case *aigentic.ApprovalEvent:
    // Send approval request to Slack
    message := fmt.Sprintf(
        "Approval needed for %s\nParameters: %v\nReact with ✅ to approve or ❌ to reject",
        e.ToolName,
        e.ValidationResult.Values,
    )

    messageID := sendSlackMessage(channel, message)

    // Wait for reaction in a goroutine
    go func() {
        reaction := waitForReaction(messageID)
        approved := reaction == "✅"
        run.Approve(e.ApprovalID, approved)
    }()
```

### Mobile Push Notifications

```go
case *aigentic.ApprovalEvent:
    // Send push notification
    notification := PushNotification{
        Title: fmt.Sprintf("Approve %s?", e.ToolName),
        Body:  fmt.Sprintf("Parameters: %v", e.ValidationResult.Values),
        Data: map[string]string{
            "approval_id": e.ApprovalID,
            "run_id":     e.RunID,
        },
    }

    sendPushNotification(userToken, notification)
```

## Best Practices

### 1. Clear Tool Descriptions
Provide detailed descriptions so the LLM uses tools appropriately:

```go
Description: "Sends an email to a recipient. Use this tool when the user explicitly asks to send an email. Requires approval before sending."
```

### 2. Comprehensive Validation
Use validation to provide context for approval decisions:

```go
Validate: func(run *aigentic.AgentRun, args map[string]interface{}) (aigentic.ValidationResult, error) {
    // Check parameter validity
    // Provide warnings for risky operations
    // Return helpful messages for approval UI
}
```

### 3. Meaningful Approval IDs
Approval IDs are automatically generated and unique. Store them with context:

```go
type ApprovalRecord struct {
    ID        string
    RunID     string
    ToolName  string
    Timestamp time.Time
    Approved  bool
    User      string
}
```

### 4. Audit Logging
Log all approval decisions for compliance:

```go
case *aigentic.ApprovalEvent:
    logApprovalRequest(e)

    approved := getUserDecision(e)
    run.Approve(e.ApprovalID, approved)

    logApprovalDecision(e.ApprovalID, approved, currentUser)
```

### 5. Timeout Configuration
Handle timeouts gracefully:

```go
// Check for pending approvals periodically
go func() {
    ticker := time.NewTicker(30 * time.Second)
    for range ticker.C {
        checkForTimedOutApprovals()
    }
}()
```

### 6. Testing with Auto-Approval
For automated testing, auto-approve all requests:

```go
func TestAgentWithApproval(t *testing.T) {
    run, _ := agent.Start("perform action")

    for event := range run.Next() {
        if e, ok := event.(*aigentic.ApprovalEvent); ok {
            run.Approve(e.ApprovalID, true) // Auto-approve for testing
        }
    }
}
```

### 7. Granular Approval Rules
Implement rule-based auto-approval for low-risk operations:

```go
case *aigentic.ApprovalEvent:
    // Auto-approve small transactions
    if e.ToolName == "transfer_money" {
        params := e.ValidationResult.Values.(map[string]interface{})
        amount := params["amount"].(float64)

        if amount < 100 {
            run.Approve(e.ApprovalID, true) // Auto-approve small amounts
            return
        }
    }

    // Require human approval for others
    approved := getUserDecision(e)
    run.Approve(e.ApprovalID, approved)
```

## Combining with Other Features

### With Streaming
Approval works seamlessly with streaming responses:

```go
agent := aigentic.Agent{
    Model:  model,
    Stream: true, // Enable streaming
    AgentTools: []aigentic.AgentTool{
        approvalTool,
    },
}

run, _ := agent.Start("task")
for event := range run.Next() {
    switch e := event.(type) {
    case *aigentic.ContentEvent:
        fmt.Print(e.Content) // Stream content as it generates
    case *aigentic.ApprovalEvent:
        approved := getApproval(e)
        run.Approve(e.ApprovalID, approved)
    }
}
```

### With Multi-Agent Systems
Sub-agents can have their own approval requirements:

```go
coordinator := aigentic.Agent{
    Name: "Coordinator",
    Agents: []aigentic.Agent{
        {
            Name: "BankingAgent",
            AgentTools: []aigentic.AgentTool{
                transferTool, // Requires approval
            },
        },
    },
}
```

### With Memory
Track approval history in agent memory:

```go
case *aigentic.ApprovalEvent:
    // Store in memory
    session.Memory.Store("approvals", []ApprovalRecord{
        {
            ToolName:  e.ToolName,
            Timestamp: time.Now(),
            Approved:  approved,
        },
    })
```

## Error Handling

### Approval Timeout
```go
case *aigentic.ErrorEvent:
    if strings.Contains(e.Err.Error(), "timeout") {
        // Notify user
        // Log timeout
        // Cleanup pending state
    }
```

### Validation Errors
```go
Validate: func(run *aigentic.AgentRun, args map[string]interface{}) (aigentic.ValidationResult, error) {
    if /* validation fails */ {
        return aigentic.ValidationResult{
            Values: args,
            ValidationErrors: []error{
                fmt.Errorf("invalid parameter: %v", param),
            },
        }, nil
    }
}
```

### Tool Execution Errors
```go
Execute: func(run *aigentic.AgentRun, args map[string]interface{}) (*ai.ToolResult, error) {
    result, err := performOperation(args)
    if err != nil {
        return &ai.ToolResult{
            Content: []ai.ToolContent{{
                Type:    "text",
                Content: fmt.Sprintf("Operation failed: %v", err),
            }},
            Error: true, // Mark as error
        }, nil
    }
    return result, nil
}
```

## Common Patterns

### Pattern 1: Always Approve
For development/testing - auto-approve everything:

```go
case *aigentic.ApprovalEvent:
    run.Approve(e.ApprovalID, true)
```

### Pattern 2: Rule-Based Approval
Auto-approve based on parameters:

```go
case *aigentic.ApprovalEvent:
    approved := evaluateApprovalRules(e)
    run.Approve(e.ApprovalID, approved)
```

### Pattern 3: Async Approval
Approve from a different goroutine or service:

```go
case *aigentic.ApprovalEvent:
    go func() {
        // Send to approval queue
        decision := waitForApprovalDecision(e.ApprovalID)
        run.Approve(e.ApprovalID, decision)
    }()
```

### Pattern 4: Multi-Level Approval
Require multiple approvers:

```go
case *aigentic.ApprovalEvent:
    approvers := []string{"manager", "finance"}
    allApproved := true

    for _, approver := range approvers {
        approved := getApprovalFrom(approver, e)
        if !approved {
            allApproved = false
            break
        }
    }

    run.Approve(e.ApprovalID, allApproved)
```

## Security Considerations

1. **Authentication**: Verify user identity before accepting approval decisions
2. **Authorization**: Check if user has permission to approve specific operations
3. **Audit Trail**: Log all approval requests and decisions with timestamps
4. **Rate Limiting**: Prevent approval request spam
5. **Encryption**: Protect sensitive parameters in transit
6. **Session Management**: Tie approvals to authenticated sessions
7. **Timeout**: Set reasonable timeouts to prevent indefinite waits

## Performance Tips

1. **Async Processing**: Don't block the main thread waiting for approval
2. **Caching**: Cache approval rules for repeated operations
3. **Batch Approvals**: Allow approving multiple similar operations at once
4. **Background Workers**: Use worker pools for approval processing
5. **Database Indexing**: Index approval records for fast lookups

## Next Steps

- See [tools example](../tools) for creating custom tools
- See [streaming example](../streaming) for real-time event handling
- See [production example](../production) for building robust production systems
- See [multi-agent example](../multi-agent) for coordinating multiple agents with approvals

## Additional Resources

- [aigentic Documentation](https://github.com/nexxia-ai/aigentic)
- [Tool Development Guide](../tools)
- [Event Handling Guide](../streaming)
- [Production Best Practices](../production)
