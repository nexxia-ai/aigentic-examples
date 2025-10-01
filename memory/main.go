package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/nexxia-ai/aigentic"
	openai "github.com/nexxia-ai/aigentic-openai"
	"github.com/nexxia-ai/aigentic/memory"
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

	fmt.Println("💾 Aigentic Memory System Examples")
	fmt.Println("===================================")
	fmt.Println()

	model := openai.NewModel("gpt-4o-mini", getAPIKey())

	// Example 1: Run Memory - Temporary memory within a single agent run
	fmt.Println("=== Example 1: Run Memory ===")
	fmt.Println("Run memory persists across LLM calls within a single agent execution")
	fmt.Println()

	runMemoryAgent := aigentic.Agent{
		Model:        model,
		Name:         "TaskManager",
		Description:  "An agent that tracks task progress using run memory",
		Instructions: "You help users complete multi-step tasks. Use run memory to track progress and intermediate results. Save important state using save_memory tool with compartment 'run'.",
		Memory:       memory.NewMemory(),
	}

	response, err := runMemoryAgent.Execute("I need to complete 3 tasks: 1) research AI agents, 2) write a summary, 3) create a presentation. Let's start with task 1.")
	if err != nil {
		log.Fatalf("Error: %v", err)
	}
	fmt.Printf("Response: %s\n\n", response)

	// Example 2: Session Memory - Persistent memory across multiple agent runs
	fmt.Println("=== Example 2: Session Memory ===")
	fmt.Println("Session memory persists across multiple agent executions in the same session")
	fmt.Println()

	// Create a session that will be shared across agent runs
	session := aigentic.NewSession(context.Background())

	sessionMemoryAgent := aigentic.Agent{
		Model:        model,
		Name:         "PersonalAssistant",
		Description:  "A personal assistant that remembers user preferences and context",
		Instructions: "You are a personal assistant. Remember user preferences, past interactions, and important information using session memory. Use save_memory tool with compartment 'session' for information that should persist across conversations.",
		Session:      session,
		Memory:       memory.NewMemory(),
	}

	// First conversation
	fmt.Println("First conversation:")
	response, err = sessionMemoryAgent.Execute("My name is Alice and I prefer morning meetings. I'm working on a project about renewable energy.")
	if err != nil {
		log.Fatalf("Error: %v", err)
	}
	fmt.Printf("Response: %s\n\n", response)

	// Second conversation - agent should remember from first conversation
	fmt.Println("Second conversation (new agent run, same session):")
	response, err = sessionMemoryAgent.Execute("What's my name and what project am I working on?")
	if err != nil {
		log.Fatalf("Error: %v", err)
	}
	fmt.Printf("Response: %s\n\n", response)

	// Third conversation - testing preference recall
	fmt.Println("Third conversation:")
	response, err = sessionMemoryAgent.Execute("When do I prefer to have meetings?")
	if err != nil {
		log.Fatalf("Error: %v", err)
	}
	fmt.Printf("Response: %s\n\n", response)

	// Example 3: Plan Memory - For complex multi-step workflows
	fmt.Println("=== Example 3: Plan Memory ===")
	fmt.Println("Plan memory helps track complex multi-step plans and their progress")
	fmt.Println()

	plannerSession := aigentic.NewSession(context.Background())

	plannerAgent := aigentic.Agent{
		Model:        model,
		Name:         "ProjectPlanner",
		Description:  "An agent that creates and tracks complex project plans",
		Instructions: "You help create detailed project plans and track their progress. Use plan memory (compartment 'plan') to store multi-step plans, milestones, and progress updates. Break complex projects into clear steps.",
		Session:      plannerSession,
		Memory:       memory.NewMemory(),
	}

	// Create a plan
	fmt.Println("Creating a project plan:")
	response, err = plannerAgent.Execute("Create a plan for launching a new mobile app. Save the plan so we can track progress.")
	if err != nil {
		log.Fatalf("Error: %v", err)
	}
	fmt.Printf("Response: %s\n\n", response)

	// Check plan progress
	fmt.Println("Checking plan status:")
	response, err = plannerAgent.Execute("What's the status of our app launch plan? Retrieve it from memory.")
	if err != nil {
		log.Fatalf("Error: %v", err)
	}
	fmt.Printf("Response: %s\n\n", response)

	// Update plan progress
	fmt.Println("Updating plan progress:")
	response, err = plannerAgent.Execute("We completed the design phase. Update the plan to reflect this progress.")
	if err != nil {
		log.Fatalf("Error: %v", err)
	}
	fmt.Printf("Response: %s\n\n", response)

	// Example 4: Memory with Multi-Agent Systems
	fmt.Println("=== Example 4: Shared Session Memory in Multi-Agent System ===")
	fmt.Println("Multiple agents can share session memory for coordination")
	fmt.Println()

	sharedSession := aigentic.NewSession(context.Background())
	sharedMemory := memory.NewMemory()

	researchAgent := aigentic.Agent{
		Model:        model,
		Name:         "Researcher",
		Description:  "Research specialist that stores findings in shared memory",
		Instructions: "You research topics and store findings in session memory for other team members to use.",
		Session:      sharedSession,
		Memory:       sharedMemory,
	}

	writerAgent := aigentic.Agent{
		Model:        model,
		Name:         "Writer",
		Description:  "Content writer that uses research from shared memory",
		Instructions: "You write content based on research stored in session memory. Retrieve research findings before writing.",
		Session:      sharedSession,
		Memory:       sharedMemory,
	}

	coordinator := aigentic.Agent{
		Model:        model,
		Name:         "TeamLead",
		Description:  "Coordinates research and writing team with shared memory",
		Instructions: "You coordinate a research and writing team. They share session memory, so research findings are available to the writer. Delegate appropriately.",
		Session:      sharedSession,
		Memory:       sharedMemory,
		Agents:       []aigentic.Agent{researchAgent, writerAgent},
	}

	response, err = coordinator.Execute("Research the top 3 benefits of electric vehicles, then write a short article about them.")
	if err != nil {
		log.Fatalf("Error: %v", err)
	}
	fmt.Printf("Final Result:\n%s\n\n", response)

	fmt.Println("✅ All memory examples completed successfully!")
	fmt.Println()
	fmt.Println("Key Takeaways:")
	fmt.Println("- Run Memory: Temporary state within a single agent execution")
	fmt.Println("- Session Memory: Persistent across agent runs in the same session")
	fmt.Println("- Plan Memory: Track complex multi-step workflows")
	fmt.Println("- Shared Memory: Multiple agents can coordinate via shared session")
}
