package main

import (
	"fmt"
	"log"
	"log/slog"
	"os"

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

func main() {
	utils.LoadEnvFile("../.env")

	fmt.Println("👥 Aigentic Multi-Agent System Examples")
	fmt.Println("========================================")
	fmt.Println()

	model := openai.NewModel("gpt-4o-mini", getAPIKey())

	// Example 1: Simple delegation - Research and Writing Team
	fmt.Println("=== Example 1: Research and Writing Team ===")
	fmt.Println("A coordinator agent delegates to specialized researcher and writer agents")
	fmt.Println()

	researchAgent := aigentic.Agent{
		Model:        model,
		Name:         "Researcher",
		Description:  "Expert at gathering and analyzing information on any topic",
		Instructions: "You are a research specialist. When given a topic, provide comprehensive, factual information with key insights and data points. Be thorough but concise.",
	}

	writerAgent := aigentic.Agent{
		Model:        model,
		Name:         "Writer",
		Description:  "Expert at creating clear, engaging written content",
		Instructions: "You are a professional writer. Create well-structured, engaging content based on the information provided. Use clear language and proper formatting.",
	}

	coordinator := aigentic.Agent{
		Model:        model,
		Name:         "ProjectManager",
		Description:  "Coordinates research and writing tasks to produce high-quality content",
		Instructions: "You manage a team of specialists. First, delegate research tasks to the Researcher. Then, use the research findings to have the Writer create the final content. Coordinate the work to ensure high quality output.",
		Agents:       []aigentic.Agent{researchAgent, writerAgent},
		LogLevel:     slog.LevelInfo,
	}

	response, err := coordinator.Execute("Create a brief article about the benefits of renewable energy, focusing on solar and wind power.")
	if err != nil {
		log.Fatalf("Error: %v", err)
	}
	fmt.Printf("Final Article:\n%s\n\n", response)

	// Example 2: Expert Panel - Multiple specialists working together
	fmt.Println("=== Example 2: Expert Panel ===")
	fmt.Println("Multiple domain experts collaborate on a complex problem")
	fmt.Println()

	techExpert := aigentic.Agent{
		Model:        model,
		Name:         "TechExpert",
		Description:  "Expert in technology, software architecture, and engineering best practices",
		Instructions: "You are a senior technical architect. Provide insights on technical feasibility, architecture, and implementation considerations.",
	}

	businessExpert := aigentic.Agent{
		Model:        model,
		Name:         "BusinessExpert",
		Description:  "Expert in business strategy, market analysis, and ROI",
		Instructions: "You are a business strategist. Evaluate business value, market fit, costs, and return on investment.",
	}

	uxExpert := aigentic.Agent{
		Model:        model,
		Name:         "UXExpert",
		Description:  "Expert in user experience, design, and usability",
		Instructions: "You are a UX designer. Focus on user needs, usability, accessibility, and overall user experience.",
	}

	panelLead := aigentic.Agent{
		Model:        model,
		Name:         "PanelLead",
		Description:  "Facilitates expert panel discussions and synthesizes recommendations",
		Instructions: "You lead an expert panel. Consult each expert (TechExpert, BusinessExpert, UXExpert) to gather their perspectives, then synthesize their insights into a comprehensive recommendation. Present different viewpoints and a balanced conclusion.",
		Agents:       []aigentic.Agent{techExpert, businessExpert, uxExpert},
		LogLevel:     slog.LevelInfo,
	}

	response, err = panelLead.Execute("Should we build a mobile app or a progressive web app (PWA) for our new product?")
	if err != nil {
		log.Fatalf("Error: %v", err)
	}
	fmt.Printf("Expert Panel Recommendation:\n%s\n\n", response)

	// Example 3: Hierarchical Team - Manager with specialized sub-teams
	fmt.Println("=== Example 3: Hierarchical Organization ===")
	fmt.Println("A CEO coordinates department heads who manage their own teams")
	fmt.Println()

	dataAnalyst := aigentic.Agent{
		Model:        model,
		Name:         "DataAnalyst",
		Description:  "Analyzes data and creates insights",
		Instructions: "You analyze data and provide statistical insights. Be data-driven and precise.",
	}

	dataScienceHead := aigentic.Agent{
		Model:        model,
		Name:         "DataScienceHead",
		Description:  "Leads the data science team",
		Instructions: "You manage data analytics. Delegate to your DataAnalyst when needed and provide strategic data insights.",
		Agents:       []aigentic.Agent{dataAnalyst},
	}

	marketingSpecialist := aigentic.Agent{
		Model:        model,
		Name:         "MarketingSpecialist",
		Description:  "Creates marketing strategies and campaigns",
		Instructions: "You develop marketing strategies based on data and market trends.",
	}

	marketingHead := aigentic.Agent{
		Model:        model,
		Name:         "MarketingHead",
		Description:  "Leads the marketing department",
		Instructions: "You manage marketing initiatives. Delegate to your MarketingSpecialist and coordinate with data teams.",
		Agents:       []aigentic.Agent{marketingSpecialist},
	}

	ceo := aigentic.Agent{
		Model:        model,
		Name:         "CEO",
		Description:  "Chief executive who coordinates all departments",
		Instructions: "You coordinate different department heads (DataScienceHead, MarketingHead) to make strategic decisions. Gather insights from each department and make informed recommendations.",
		Agents:       []aigentic.Agent{dataScienceHead, marketingHead},
		LogLevel:     slog.LevelInfo,
	}

	response, err = ceo.Execute("We're launching a new product next quarter. What should be our go-to-market strategy?")
	if err != nil {
		log.Fatalf("Error: %v", err)
	}
	fmt.Printf("Strategic Recommendation:\n%s\n\n", response)

	fmt.Println("✅ All multi-agent examples completed successfully!")
	fmt.Println()
	fmt.Println("Key Takeaways:")
	fmt.Println("- Agents can delegate tasks to specialized sub-agents")
	fmt.Println("- Sub-agents are exposed as tools to parent agents")
	fmt.Println("- Build complex hierarchies for sophisticated workflows")
	fmt.Println("- Each agent maintains its own expertise and instructions")
}
