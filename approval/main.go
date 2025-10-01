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
			time.Sleep(500 * time.Millisecond)
			return fmt.Sprintf("Email successfully sent to %s with subject '%s'", input.To, input.Subject), nil
		},
	)
	emailTool.RequireApproval = true
	return emailTool
}

func simulateApprovalUI(e *aigentic.ApprovalEvent) bool {
	fmt.Println("\n" + strings.Repeat("=", 70))
	fmt.Println("APPROVAL REQUIRED")
	fmt.Println(strings.Repeat("=", 70))
	fmt.Printf("Tool: %s\n", e.ToolName)
	fmt.Printf("Approval ID: %s\n", e.ApprovalID)

	if e.ValidationResult.Message != "" {
		fmt.Printf("Validation: %s\n", e.ValidationResult.Message)
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

func main() {
	utils.LoadEnvFile("../.env")

	fmt.Println("Human-in-the-Loop Approval Example")
	fmt.Println("==================================")
	fmt.Println()

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
	fmt.Println("\n✅ Example completed successfully!")
}
