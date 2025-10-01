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
	type SendEmailInput struct {
		To      string `json:"to" description:"Email recipient address"`
		Subject string `json:"subject" description:"Email subject line"`
		Body    string `json:"body" description:"Email body content"`
	}

	emailTool := aigentic.NewTool(
		"send_email",
		"Sends an email to a recipient with subject and body. Requires approval before sending.",
		func(run *aigentic.AgentRun, input SendEmailInput) (string, error) {
			// Simulate sending email
			time.Sleep(500 * time.Millisecond)
			return fmt.Sprintf("Email successfully sent to %s with subject '%s'", input.To, input.Subject), nil
		},
	)
	emailTool.RequireApproval = true
	return emailTool
}

// createDeleteFileTool demonstrates a destructive operation that requires approval
func createDeleteFileTool() aigentic.AgentTool {
	type DeleteFileInput struct {
		Filepath string `json:"filepath" description:"Path to the file to delete"`
		Reason   string `json:"reason" description:"Reason for deleting the file"`
	}

	deleteTool := aigentic.NewTool(
		"delete_file",
		"Deletes a file from the filesystem. This is a destructive operation that requires approval.",
		func(run *aigentic.AgentRun, input DeleteFileInput) (string, error) {
			// Simulate file deletion (don't actually delete)
			time.Sleep(300 * time.Millisecond)
			return fmt.Sprintf("File '%s' has been deleted. Reason: %s", input.Filepath, input.Reason), nil
		},
	)
	deleteTool.RequireApproval = true
	return deleteTool
}

// createTransferMoneyTool demonstrates a financial transaction with validation and approval
func createTransferMoneyTool() aigentic.AgentTool {
	type TransferMoneyInput struct {
		FromAccount string  `json:"from_account" description:"Source account number"`
		ToAccount   string  `json:"to_account" description:"Destination account number"`
		Amount      float64 `json:"amount" description:"Amount to transfer in USD"`
		Memo        string  `json:"memo,omitempty" description:"Optional memo for the transaction"`
	}

	transferTool := aigentic.NewTool(
		"transfer_money",
		"Transfers money from one account to another. Requires approval for amounts over $100.",
		func(run *aigentic.AgentRun, input TransferMoneyInput) (string, error) {
			// Simulate money transfer
			time.Sleep(1 * time.Second)

			result := fmt.Sprintf("Successfully transferred $%.2f from %s to %s", input.Amount, input.FromAccount, input.ToAccount)
			if input.Memo != "" {
				result += fmt.Sprintf(" (Memo: %s)", input.Memo)
			}

			return result, nil
		},
	)
	transferTool.RequireApproval = true

	// Add custom validation for large amounts
	transferTool.Validate = func(run *aigentic.AgentRun, args map[string]interface{}) (aigentic.ValidationResult, error) {
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
	}

	return transferTool
}

// createDatabaseQueryTool demonstrates a read-only tool that doesn't require approval
func createDatabaseQueryTool() aigentic.AgentTool {
	type DatabaseQueryInput struct {
		Query string `json:"query" description:"SQL query to execute"`
	}

	queryTool := aigentic.NewTool(
		"query_database",
		"Queries the database for information. Read-only operation, no approval needed.",
		func(run *aigentic.AgentRun, input DatabaseQueryInput) (string, error) {
			// Simulate database query
			time.Sleep(200 * time.Millisecond)
			return fmt.Sprintf("Query executed: %s\nResults: [{'id': 1, 'name': 'Sample Data'}]", input.Query), nil
		},
	)
	queryTool.RequireApproval = false // Explicitly false for demonstration
	return queryTool
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
