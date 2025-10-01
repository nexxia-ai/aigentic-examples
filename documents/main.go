package main

import (
	"context"
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

	// Example 1: Analyzing an embedded text document
	fmt.Println("=== Example 1: Embedded Text Document ===")
	fmt.Println("Embedding a text document directly in the agent's context")
	fmt.Println()

	// Create a sample contract text
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

	// Create an in-memory document with the contract text
	contractDoc := document.NewInMemoryDocument(
		"contract_001",
		"employment_contract.txt",
		[]byte(contractText),
		nil,
	)

	// Create agent with embedded document
	contractAgent := aigentic.Agent{
		Model:        model,
		Name:         "ContractAnalyzer",
		Description:  "Analyzes employment contracts and extracts key information",
		Instructions: "You are a legal assistant specializing in employment contracts. Analyze the provided contract and extract key information clearly and accurately.",
		Documents:    []*document.Document{contractDoc}, // Document is embedded in context
	}

	response, err := contractAgent.Execute("What is the base salary and what benefits are included in this employment contract?")
	if err != nil {
		log.Fatalf("Error: %v", err)
	}
	fmt.Printf("Analysis:\n%s\n\n", response)

	// Example 2: Processing an embedded image with receipt data
	fmt.Println("=== Example 2: Embedded Image Document ===")
	fmt.Println("Processing a receipt image for data extraction")
	fmt.Println()

	// Create a simple representation of receipt data
	// In a real scenario, this would be actual image bytes (PNG, JPG, etc.)
	receiptText := `
=================================
      GROCERY MART
   123 Main Street
   Anytown, ST 12345
   Tel: (555) 123-4567
=================================

Date: 2024-10-01      Time: 14:23
Cashier: Alice        Register: 3

ITEMS PURCHASED:
--------------------------------
Organic Bananas      $3.99
Whole Wheat Bread    $4.50
2% Milk (1 Gallon)   $5.25
Free Range Eggs      $6.99
Fresh Spinach        $3.49
Chicken Breast       $12.99
Olive Oil            $8.99
--------------------------------
Subtotal:           $46.20
Tax (8%):            $3.70
--------------------------------
TOTAL:              $49.90
--------------------------------

Payment Method: Credit Card
Card: **** **** **** 1234

Thank you for shopping!
=================================
`

	// Create document from text (in production, use actual image data)
	receiptDoc := document.NewInMemoryDocument(
		"receipt_001",
		"receipt.txt",
		[]byte(receiptText),
		nil,
	)

	receiptAgent := aigentic.Agent{
		Model:        model,
		Name:         "ReceiptProcessor",
		Description:  "Extracts structured data from receipt images",
		Instructions: "You are a receipt processing system. Extract key information including total amount, date, items, and payment method. Format your response as structured data.",
		Documents:    []*document.Document{receiptDoc},
	}

	response, err = receiptAgent.Execute("Extract the total amount, date, and list the top 3 most expensive items from this receipt.")
	if err != nil {
		log.Fatalf("Error: %v", err)
	}
	fmt.Printf("Extracted Data:\n%s\n\n", response)

	// Example 3: Using DocumentReferences for large documents
	fmt.Println("=== Example 3: Document References (On-Demand Loading) ===")
	fmt.Println("Using document references allows the LLM to decide when to retrieve content")
	fmt.Println()

	// Create multiple documents
	technicalDoc := document.NewInMemoryDocument(
		"tech_spec",
		"technical_specification.txt",
		[]byte(`
TECHNICAL SPECIFICATION - Cloud Storage API

Version: 2.1.0
Last Updated: 2024-09-15

OVERVIEW:
The Cloud Storage API provides RESTful endpoints for file management operations.

AUTHENTICATION:
- OAuth 2.0 with Bearer tokens
- API keys for server-to-server communication
- Token expiration: 3600 seconds

ENDPOINTS:

1. Upload File
   POST /api/v2/files/upload
   Headers: Authorization: Bearer {token}
   Body: multipart/form-data
   Max Size: 5GB

2. Download File
   GET /api/v2/files/{file_id}
   Response: Binary stream

3. List Files
   GET /api/v2/files
   Query Params: page, limit, filter
   Response: JSON array

RATE LIMITS:
- Free Tier: 100 requests/hour
- Pro Tier: 1000 requests/hour
- Enterprise: Unlimited

ERROR CODES:
- 401: Unauthorized
- 413: File too large
- 429: Rate limit exceeded
`),
		nil,
	)

	userGuideDoc := document.NewInMemoryDocument(
		"user_guide",
		"user_guide.txt",
		[]byte(`
USER GUIDE - Cloud Storage Platform

GETTING STARTED:

1. Create an Account
   Visit www.cloudstorage.example.com/signup
   Provide email and create a strong password

2. Download Desktop App
   Available for Windows, Mac, and Linux
   Install and sign in with your credentials

3. Upload Your First File
   - Drag and drop files into the app
   - Or click "Upload" button
   - Files sync automatically to the cloud

SHARING FILES:

1. Right-click on any file
2. Select "Share"
3. Enter recipient email addresses
4. Set permissions (view, edit, or comment)
5. Click "Send Invitation"

MOBILE ACCESS:

Download our mobile apps:
- iOS: Available on App Store
- Android: Available on Google Play

Features:
- Upload photos automatically
- Offline access to important files
- Scan documents with your camera

PRICING:

Free Plan: 10GB storage
Pro Plan: 1TB storage - $9.99/month
Enterprise: Unlimited - Contact sales
`),
		nil,
	)

	// Create agent with document references (not embedded)
	// The LLM will use built-in tools to retrieve documents only when needed
	docsAgent := aigentic.Agent{
		Model:        model,
		Name:         "DocumentationAssistant",
		Description:  "Helps users find information in product documentation",
		Instructions: "You have access to technical specifications and user guides. Retrieve and reference the appropriate documents to answer user questions accurately. Only retrieve documents when needed to answer the question.",
		DocumentReferences: []*document.Document{technicalDoc, userGuideDoc}, // Referenced, not embedded
	}

	response, err = docsAgent.Execute("What are the API rate limits for different pricing tiers?")
	if err != nil {
		log.Fatalf("Error: %v", err)
	}
	fmt.Printf("Answer:\n%s\n\n", response)

	// Example 4: Multi-document analysis
	fmt.Println("=== Example 4: Multi-Document Analysis ===")
	fmt.Println("Comparing and analyzing multiple documents together")
	fmt.Println()

	// Create quarterly reports
	q1Report := document.NewInMemoryDocument(
		"q1_2024",
		"q1_2024_report.txt",
		[]byte(`
Q1 2024 FINANCIAL REPORT

Revenue: $2.4M (+15% YoY)
Expenses: $1.8M
Net Profit: $600K
Profit Margin: 25%

Key Metrics:
- New Customers: 1,200
- Customer Retention: 92%
- Average Deal Size: $2,000

Highlights:
- Launched new product line
- Expanded to 3 new markets
- Hired 15 new employees
`),
		nil,
	)

	q2Report := document.NewInMemoryDocument(
		"q2_2024",
		"q2_2024_report.txt",
		[]byte(`
Q2 2024 FINANCIAL REPORT

Revenue: $2.8M (+17% YoY, +16.7% QoQ)
Expenses: $2.0M
Net Profit: $800K
Profit Margin: 28.6%

Key Metrics:
- New Customers: 1,500
- Customer Retention: 94%
- Average Deal Size: $2,200

Highlights:
- Product line gained traction
- Partnership with Fortune 500 company
- Opened new office in Austin
- Hired 20 new employees
`),
		nil,
	)

	// Embed both documents for comparison
	analysisAgent := aigentic.Agent{
		Model:        model,
		Name:         "FinancialAnalyst",
		Description:  "Analyzes financial reports and identifies trends",
		Instructions: "You are a financial analyst. Compare multiple quarterly reports, identify trends, and provide insights on business performance.",
		Documents:    []*document.Document{q1Report, q2Report},
	}

	response, err = analysisAgent.Execute("Compare Q1 and Q2 2024 performance. What are the key trends and improvements?")
	if err != nil {
		log.Fatalf("Error: %v", err)
	}
	fmt.Printf("Financial Analysis:\n%s\n\n", response)

	// Example 5: Loading documents from local filesystem
	fmt.Println("=== Example 5: Loading Documents from Local Store ===")
	fmt.Println("Using LocalStore to load documents from the filesystem")
	fmt.Println()

	// Create a document store pointing to the testdata directory
	// First, let's create a test file
	testDataDir := "./testdata"
	err = os.MkdirAll(testDataDir, 0755)
	if err != nil {
		log.Fatalf("Error creating testdata directory: %v", err)
	}

	sampleContent := `
PRODUCT REQUIREMENTS DOCUMENT

Product Name: TaskFlow Pro
Version: 1.0
Date: October 2024

OVERVIEW:
TaskFlow Pro is a team productivity application that helps teams manage projects and tasks efficiently.

TARGET USERS:
- Small to medium-sized teams (5-50 people)
- Remote and hybrid work environments
- Project managers and team leads

CORE FEATURES:

1. Task Management
   - Create, assign, and track tasks
   - Set due dates and priorities
   - Add comments and attachments

2. Project Views
   - Kanban boards
   - Gantt charts
   - Calendar view
   - List view

3. Team Collaboration
   - Real-time updates
   - @mentions and notifications
   - File sharing
   - Activity feed

4. Reporting
   - Team productivity metrics
   - Project progress tracking
   - Time tracking reports
   - Custom dashboards

SUCCESS METRICS:
- User adoption: 80% of team within 30 days
- Daily active usage: >60%
- Task completion rate improvement: >20%
- User satisfaction score: >4.5/5
`

	sampleFile := testDataDir + "/product_requirements.txt"
	err = os.WriteFile(sampleFile, []byte(sampleContent), 0644)
	if err != nil {
		log.Fatalf("Error writing sample file: %v", err)
	}

	// Create a local document store
	store := document.NewLocalStore(testDataDir)

	// Load document from store (lazy loading - content not loaded yet)
	ctx := context.Background()
	localDoc, err := store.Open(ctx, "product_requirements.txt")
	if err != nil {
		log.Fatalf("Error opening document: %v", err)
	}

	// Create agent with the document
	productAgent := aigentic.Agent{
		Model:        model,
		Name:         "ProductAnalyst",
		Description:  "Analyzes product requirements and provides insights",
		Instructions: "You are a product analyst. Review product requirements documents and provide clear summaries and recommendations.",
		Documents:    []*document.Document{localDoc},
	}

	response, err = productAgent.Execute("What are the core features and success metrics for TaskFlow Pro?")
	if err != nil {
		log.Fatalf("Error: %v", err)
	}
	fmt.Printf("Product Analysis:\n%s\n\n", response)

	// Clean up
	err = store.Close(ctx)
	if err != nil {
		log.Printf("Warning: Error closing store: %v", err)
	}

	fmt.Println("All document examples completed successfully!")
	fmt.Println()
	fmt.Println("Key Takeaways:")
	fmt.Println("- Documents: Embedded directly in the agent's context (best for small documents)")
	fmt.Println("- DocumentReferences: Retrieved on-demand by the LLM (best for large document sets)")
	fmt.Println("- LocalStore: Load documents from filesystem with lazy loading")
	fmt.Println("- Multi-document: Compare and analyze multiple documents together")
	fmt.Println("- Use cases: Contracts, receipts, technical docs, financial reports, and more")
}
