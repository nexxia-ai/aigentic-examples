package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
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

// createSendEmailTool demonstrates a simple tool that requires approval before execution
func createSendEmailTool() aigentic.AgentTool {
	return aigentic.AgentTool{
		Name:            "send_email",
		Description:     "Sends an email to a recipient with subject and body. Requires approval before sending.",
		RequireApproval: true,
		InputSchema: map[string]interface{}{
			"type": "object",
			"properties": map[string]interface{}{
				"to": map[string]interface{}{
					"type":        "string",
					"description": "Email recipient address",
				},
				"subject": map[string]interface{}{
					"type":        "string",
					"description": "Email subject line",
				},
				"body": map[string]interface{}{
					"type":        "string",
					"description": "Email body content",
				},
			},
			"required": []string{"to", "subject", "body"},
		},
		Execute: func(run *aigentic.AgentRun, args map[string]interface{}) (*ai.ToolResult, error) {
			to := args["to"].(string)
			subject := args["subject"].(string)
			_ = args["body"].(string) // body variable extracted for completeness

			// Simulate sending email
			time.Sleep(500 * time.Millisecond)

			return &ai.ToolResult{
				Content: []ai.ToolContent{{
					Type:    "text",
					Content: fmt.Sprintf("Email successfully sent to %s with subject '%s'", to, subject),
				}},
			}, nil
		},
	}
}

// createDeleteFileTool demonstrates a destructive operation that requires approval
func createDeleteFileTool() aigentic.AgentTool {
	return aigentic.AgentTool{
		Name:            "delete_file",
		Description:     "Deletes a file from the filesystem. This is a destructive operation that requires approval.",
		RequireApproval: true,
		InputSchema: map[string]interface{}{
			"type": "object",
			"properties": map[string]interface{}{
				"filepath": map[string]interface{}{
					"type":        "string",
					"description": "Path to the file to delete",
				},
				"reason": map[string]interface{}{
					"type":        "string",
					"description": "Reason for deleting the file",
				},
			},
			"required": []string{"filepath", "reason"},
		},
		Execute: func(run *aigentic.AgentRun, args map[string]interface{}) (*ai.ToolResult, error) {
			filepath := args["filepath"].(string)
			reason := args["reason"].(string)

			// Simulate file deletion (don't actually delete)
			time.Sleep(300 * time.Millisecond)

			return &ai.ToolResult{
				Content: []ai.ToolContent{{
					Type:    "text",
					Content: fmt.Sprintf("File '%s' has been deleted. Reason: %s", filepath, reason),
				}},
			}, nil
		},
	}
}

// createTransferMoneyTool demonstrates a financial transaction with validation and approval
func createTransferMoneyTool() aigentic.AgentTool {
	return aigentic.AgentTool{
		Name:            "transfer_money",
		Description:     "Transfers money from one account to another. Requires approval for amounts over $100.",
		RequireApproval: true,
		InputSchema: map[string]interface{}{
			"type": "object",
			"properties": map[string]interface{}{
				"from_account": map[string]interface{}{
					"type":        "string",
					"description": "Source account number",
				},
				"to_account": map[string]interface{}{
					"type":        "string",
					"description": "Destination account number",
				},
				"amount": map[string]interface{}{
					"type":        "number",
					"description": "Amount to transfer in USD",
				},
				"memo": map[string]interface{}{
					"type":        "string",
					"description": "Optional memo for the transaction",
				},
			},
			"required": []string{"from_account", "to_account", "amount"},
		},
		Validate: func(run *aigentic.AgentRun, args map[string]interface{}) (aigentic.ValidationResult, error) {
			amount, ok := args["amount"].(float64)
			if !ok {
				return aigentic.ValidationResult{
					Values:  args,
					Message: "Invalid amount format",
					ValidationErrors: []error{
						fmt.Errorf("amount must be a number"),
					},
				}, nil
			}

			// Add validation warnings for large amounts
			var message string
			if amount > 10000 {
				message = fmt.Sprintf("WARNING: Large transaction amount: $%.2f", amount)
			} else if amount > 1000 {
				message = fmt.Sprintf("CAUTION: Moderate transaction amount: $%.2f", amount)
			} else {
				message = fmt.Sprintf("Transaction amount: $%.2f", amount)
			}

			return aigentic.ValidationResult{
				Values:  args,
				Message: message,
			}, nil
		},
		Execute: func(run *aigentic.AgentRun, args map[string]interface{}) (*ai.ToolResult, error) {
			fromAccount := args["from_account"].(string)
			toAccount := args["to_account"].(string)
			amount := args["amount"].(float64)
			memo := ""
			if m, ok := args["memo"].(string); ok {
				memo = m
			}

			// Simulate money transfer
			time.Sleep(1 * time.Second)

			result := fmt.Sprintf("Successfully transferred $%.2f from %s to %s", amount, fromAccount, toAccount)
			if memo != "" {
				result += fmt.Sprintf(" (Memo: %s)", memo)
			}

			return &ai.ToolResult{
				Content: []ai.ToolContent{{
					Type:    "text",
					Content: result,
				}},
			}, nil
		},
	}
}

// createDatabaseQueryTool demonstrates a read-only tool that doesn't require approval
func createDatabaseQueryTool() aigentic.AgentTool {
	return aigentic.AgentTool{
		Name:            "query_database",
		Description:     "Queries the database for information. Read-only operation, no approval needed.",
		RequireApproval: false, // Explicitly false for demonstration
		InputSchema: map[string]interface{}{
			"type": "object",
			"properties": map[string]interface{}{
				"query": map[string]interface{}{
					"type":        "string",
					"description": "SQL query to execute",
				},
			},
			"required": []string{"query"},
		},
		Execute: func(run *aigentic.AgentRun, args map[string]interface{}) (*ai.ToolResult, error) {
			query := args["query"].(string)

			// Simulate database query
			time.Sleep(200 * time.Millisecond)

			return &ai.ToolResult{
				Content: []ai.ToolContent{{
					Type:    "text",
					Content: fmt.Sprintf("Query executed: %s\nResults: [{'id': 1, 'name': 'Sample Data'}]", query),
				}},
			}, nil
		},
	}
}

// simulateApprovalUI simulates a user interface for approval decisions
func simulateApprovalUI(e *aigentic.ApprovalEvent) bool {
	fmt.Println("\n" + strings.Repeat("=", 70))
	fmt.Println("APPROVAL REQUIRED")
	fmt.Println(strings.Repeat("=", 70))
	fmt.Printf("Tool: %s\n", e.ToolName)
	fmt.Printf("Approval ID: %s\n", e.ApprovalID)

	if e.ValidationResult.Message != "" {
		fmt.Printf("Validation: %s\n", e.ValidationResult.Message)
	}

	if len(e.ValidationResult.ValidationErrors) > 0 {
		fmt.Println("\nValidation Errors:")
		for _, err := range e.ValidationResult.ValidationErrors {
			fmt.Printf("  - %v\n", err)
		}
	}

	if args, ok := e.ValidationResult.Values.(map[string]interface{}); ok {
		fmt.Println("\nParameters:")
		for key, value := range args {
			fmt.Printf("  %s: %v\n", key, value)
		}
	}

	fmt.Println(strings.Repeat("=", 70))
	fmt.Print("Approve this action? (y/n): ")

	reader := bufio.NewReader(os.Stdin)
	response, _ := reader.ReadString('\n')
	response = strings.TrimSpace(strings.ToLower(response))

	approved := response == "y" || response == "yes"
	if approved {
		fmt.Println("✓ Action APPROVED")
	} else {
		fmt.Println("✗ Action REJECTED")
	}
	fmt.Println(strings.Repeat("=", 70) + "\n")

	return approved
}

// runExample1 demonstrates simple tool approval with email
func runExample1() {
	fmt.Println("\n=== Example 1: Simple Email Approval ===")
	fmt.Println("This example shows a basic approval workflow for sending an email.\n")

	model := openai.NewModel("gpt-4o-mini", getAPIKey())

	agent := aigentic.Agent{
		Model:        model,
		Name:         "EmailAgent",
		Description:  "An agent that can send emails with approval",
		Instructions: "You can send emails using the send_email tool. Always use the tool when asked to send an email.",
		AgentTools: []aigentic.AgentTool{
			createSendEmailTool(),
		},
		Stream: true,
	}

	run, err := agent.Start("Send an email to john@example.com with subject 'Project Update' and body 'The project is on track and will be completed by end of week.'")
	if err != nil {
		log.Fatalf("Failed to start agent: %v", err)
	}

	var fullResponse string
	for event := range run.Next() {
		switch e := event.(type) {
		case *aigentic.ContentEvent:
			fmt.Print(e.Content)
			fullResponse += e.Content
		case *aigentic.ApprovalEvent:
			approved := simulateApprovalUI(e)
			run.Approve(e.ApprovalID, approved)
		case *aigentic.ToolEvent:
			fmt.Printf("\n[Tool executed: %s]\n", e.ToolName)
		case *aigentic.ErrorEvent:
			log.Printf("Error: %v", e.Err)
		}
	}

	fmt.Printf("\n\nFinal Response: %s\n", fullResponse)
}

// runExample2 demonstrates file deletion approval
func runExample2() {
	fmt.Println("\n=== Example 2: File Deletion Approval ===")
	fmt.Println("This example shows approval for destructive operations.\n")

	model := openai.NewModel("gpt-4o-mini", getAPIKey())

	agent := aigentic.Agent{
		Model:        model,
		Name:         "FileAgent",
		Description:  "An agent that can manage files with approval",
		Instructions: "You can delete files using the delete_file tool. Always provide a clear reason for deletion.",
		AgentTools: []aigentic.AgentTool{
			createDeleteFileTool(),
		},
		Stream: true,
	}

	run, err := agent.Start("Delete the file /tmp/old-logs.txt because it contains outdated information from last year.")
	if err != nil {
		log.Fatalf("Failed to start agent: %v", err)
	}

	for event := range run.Next() {
		switch e := event.(type) {
		case *aigentic.ContentEvent:
			fmt.Print(e.Content)
		case *aigentic.ApprovalEvent:
			approved := simulateApprovalUI(e)
			run.Approve(e.ApprovalID, approved)
		case *aigentic.ToolEvent:
			fmt.Printf("\n[Tool executed: %s]\n", e.ToolName)
		case *aigentic.ErrorEvent:
			log.Printf("Error: %v", e.Err)
		}
	}
}

// runExample3 demonstrates financial transaction with validation
func runExample3() {
	fmt.Println("\n=== Example 3: Financial Transaction with Validation ===")
	fmt.Println("This example shows approval with custom validation logic.\n")

	model := openai.NewModel("gpt-4o-mini", getAPIKey())

	agent := aigentic.Agent{
		Model:        model,
		Name:         "BankingAgent",
		Description:  "An agent that can perform financial transactions with approval",
		Instructions: "You can transfer money using the transfer_money tool. Always verify the amounts and accounts.",
		AgentTools: []aigentic.AgentTool{
			createTransferMoneyTool(),
		},
		Stream: true,
	}

	run, err := agent.Start("Transfer $5,000 from account 123-456-789 to account 987-654-321 with memo 'Monthly payment'")
	if err != nil {
		log.Fatalf("Failed to start agent: %v", err)
	}

	for event := range run.Next() {
		switch e := event.(type) {
		case *aigentic.ContentEvent:
			fmt.Print(e.Content)
		case *aigentic.ApprovalEvent:
			approved := simulateApprovalUI(e)
			run.Approve(e.ApprovalID, approved)
		case *aigentic.ToolEvent:
			fmt.Printf("\n[Tool executed: %s]\n", e.ToolName)
		case *aigentic.ErrorEvent:
			log.Printf("Error: %v", e.Err)
		}
	}
}

// runExample4 demonstrates mixed tools with and without approval
func runExample4() {
	fmt.Println("\n=== Example 4: Mixed Tools with Selective Approval ===")
	fmt.Println("This example shows combining tools that require approval with those that don't.\n")

	model := openai.NewModel("gpt-4o-mini", getAPIKey())

	agent := aigentic.Agent{
		Model:        model,
		Name:         "MixedAgent",
		Description:  "An agent with both approved and non-approved tools",
		Instructions: "You have access to database queries (no approval) and money transfers (requires approval). Use them as needed.",
		AgentTools: []aigentic.AgentTool{
			createDatabaseQueryTool(),
			createTransferMoneyTool(),
		},
		Stream: true,
	}

	run, err := agent.Start("First, query the database for account 123-456-789 balance. Then transfer $250 to account 999-888-777.")
	if err != nil {
		log.Fatalf("Failed to start agent: %v", err)
	}

	for event := range run.Next() {
		switch e := event.(type) {
		case *aigentic.ContentEvent:
			fmt.Print(e.Content)
		case *aigentic.ApprovalEvent:
			approved := simulateApprovalUI(e)
			run.Approve(e.ApprovalID, approved)
		case *aigentic.ToolEvent:
			if e.ToolName == "query_database" {
				fmt.Printf("\n[Database query executed - no approval needed]\n")
			} else {
				fmt.Printf("\n[Tool executed: %s]\n", e.ToolName)
			}
		case *aigentic.ErrorEvent:
			log.Printf("Error: %v", e.Err)
		}
	}
}

// runAutomatedExample demonstrates automated approval for testing/demos
func runAutomatedExample() {
	fmt.Println("\n=== Automated Example: Auto-Approve for Testing ===")
	fmt.Println("This example shows how to automatically approve requests for testing purposes.\n")

	model := openai.NewModel("gpt-4o-mini", getAPIKey())

	agent := aigentic.Agent{
		Model:        model,
		Name:         "AutoApproveAgent",
		Description:  "An agent that auto-approves all actions",
		Instructions: "Send emails and perform operations as requested.",
		AgentTools: []aigentic.AgentTool{
			createSendEmailTool(),
			createTransferMoneyTool(),
		},
		Stream: true,
	}

	run, err := agent.Start("Send an email to team@example.com about the meeting tomorrow.")
	if err != nil {
		log.Fatalf("Failed to start agent: %v", err)
	}

	for event := range run.Next() {
		switch e := event.(type) {
		case *aigentic.ContentEvent:
			fmt.Print(e.Content)
		case *aigentic.ApprovalEvent:
			// Automated approval for testing
			fmt.Printf("\n[AUTO-APPROVED: %s]\n", e.ToolName)
			run.Approve(e.ApprovalID, true)
		case *aigentic.ToolEvent:
			fmt.Printf("\n[Tool executed: %s]\n", e.ToolName)
		case *aigentic.ErrorEvent:
			log.Printf("Error: %v", e.Err)
		}
	}
}

func main() {
	utils.LoadEnvFile("../.env")

	fmt.Println("Human-in-the-Loop Approval Examples")
	fmt.Println("====================================")
	fmt.Println()
	fmt.Println("This example demonstrates various approval workflows for sensitive operations.")
	fmt.Println()

	// Check if running in non-interactive mode
	if len(os.Args) > 1 && os.Args[1] == "--auto" {
		runAutomatedExample()
		return
	}

	fmt.Println("Available Examples:")
	fmt.Println("1. Simple Email Approval")
	fmt.Println("2. File Deletion Approval")
	fmt.Println("3. Financial Transaction with Validation")
	fmt.Println("4. Mixed Tools with Selective Approval")
	fmt.Println("5. Run All Examples")
	fmt.Println("6. Automated Example (auto-approve)")
	fmt.Println()
	fmt.Print("Select an example (1-6): ")

	reader := bufio.NewReader(os.Stdin)
	choice, _ := reader.ReadString('\n')
	choice = strings.TrimSpace(choice)

	switch choice {
	case "1":
		runExample1()
	case "2":
		runExample2()
	case "3":
		runExample3()
	case "4":
		runExample4()
	case "5":
		runExample1()
		runExample2()
		runExample3()
		runExample4()
	case "6":
		runAutomatedExample()
	default:
		fmt.Println("Invalid choice. Running Example 1 by default.")
		runExample1()
	}

	fmt.Println("\n✅ Example completed successfully!")
}
