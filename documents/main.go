package main

import (
	"fmt"
	"log"
	"os"

	"github.com/nexxia-ai/aigentic"
	openai "github.com/nexxia-ai/aigentic-openai"
	"github.com/nexxia-ai/aigentic/document"
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

func main() {
	utils.LoadEnvFile("../.env")

	fmt.Println("Document Processing with Aigentic")
	fmt.Println("===================================")
	fmt.Println()

	model := openai.NewModel("gpt-4o-mini", getAPIKey())

	contractText := `
EMPLOYMENT CONTRACT

This Employment Agreement is entered into on January 15, 2024, between:

EMPLOYER: TechCorp Industries, Inc.
EMPLOYEE: Jane Smith

1. POSITION AND DUTIES
The Employee is hired as a Senior Software Engineer and will report to the Engineering Director.

2. COMPENSATION
   - Base Salary: $145,000 per year
   - Annual Bonus: Up to 20% of base salary based on performance
   - Stock Options: 5,000 shares vesting over 4 years

3. BENEFITS
   - Health Insurance: Comprehensive medical, dental, and vision
   - Retirement: 401(k) with 5% company match
   - Paid Time Off: 20 days per year
   - Remote Work: Hybrid schedule (3 days office, 2 days remote)

4. TERM
This agreement begins on February 1, 2024, and continues indefinitely unless terminated.

5. TERMINATION
Either party may terminate with 30 days written notice.

Signed: _____________________
Date: January 15, 2024
`

	contractDoc := document.NewInMemoryDocument(
		"contract_001",
		"employment_contract.txt",
		[]byte(contractText),
		nil,
	)

	agent := aigentic.Agent{
		Model:        model,
		Name:         "ContractAnalyzer",
		Description:  "Analyzes employment contracts and extracts key information",
		Instructions: "You are a legal assistant specializing in employment contracts. Analyze the provided contract and extract key information clearly and accurately.",
		Documents:    []*document.Document{contractDoc},
	}

	response, err := agent.Execute("What is the base salary and what benefits are included in this employment contract?")
	if err != nil {
		log.Fatalf("Error: %v", err)
	}
	fmt.Printf("Analysis:\n%s\n\n", response)

	fmt.Println("✅ Example completed successfully!")
}
