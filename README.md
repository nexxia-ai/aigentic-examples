# Aigentic Examples

Runnable examples for the [aigentic](https://github.com/nexxia-ai/aigentic) framework - a declarative AI agent framework for Go.

## Quick Start

1. **Set your OpenAI API key:**
   ```bash
   export OPENAI_API_KEY=sk-...
   ```

2. **Run any example directly:**
   ```bash
   go run github.com/nexxia-ai/aigentic-examples/simple@latest
   ```

3. **Or clone and run locally:**
   ```bash
   git clone https://github.com/nexxia-ai/aigentic-examples
   cd aigentic-examples/simple
   go run main.go
   ```

---

## Examples by Category

### 🚀 Getting Started

#### [simple/](simple/)
**Basic agent usage** - Your first aigentic agent
Learn: Agent creation, simple Q&A, basic configuration

```bash
go run github.com/nexxia-ai/aigentic-examples/simple@latest
```

#### [streaming/](streaming/)
**Real-time streaming responses** - Live content generation
Learn: Streaming mode, event handling, progress updates

```bash
go run github.com/nexxia-ai/aigentic-examples/streaming@latest
```

---

### 🛠️ Tool Integration

#### [tools/](tools/)
**Custom tool creation** - Calculator, weather, and time tools
Learn: Tool schemas, execution logic, parameter validation, error handling

```bash
go run github.com/nexxia-ai/aigentic-examples/tools@latest
```

#### [mcp/](mcp/)
**Model Context Protocol** - MCP server integration
Learn: MCP tools, external integrations, protocol usage

```bash
go run github.com/nexxia-ai/aigentic-examples/mcp@latest
```

---

### 👥 Multi-Agent Systems

#### [multi-agent/](multi-agent/)
**Agent teams and coordination** - Hierarchical agent delegation
Learn: Sub-agents, team coordination, expert panels, organizational hierarchies

```bash
go run github.com/nexxia-ai/aigentic-examples/multi-agent@latest
```

---

### 💾 Memory & Context

#### [memory/](memory/)
**Persistent memory system** - Run, session, and plan memory
Learn: Memory compartments, context persistence, shared state across agents

```bash
go run github.com/nexxia-ai/aigentic-examples/memory@latest
```

---

### 📄 Document Processing

#### [documents/](documents/)
**Document analysis** - Text, images, PDFs
Learn: Embedded documents, document references, multi-document analysis

```bash
go run github.com/nexxia-ai/aigentic-examples/documents@latest
```

---

### 🔒 Human-in-the-Loop

#### [approval/](approval/)
**Approval workflows** - Human oversight for sensitive operations
Learn: RequireApproval, ApprovalEvent, validation, timeout handling

```bash
go run github.com/nexxia-ai/aigentic-examples/approval@latest
```

---

### 🏭 Production Patterns

#### [production/](production/)
**Production-ready patterns** - Error handling, monitoring, debugging
Learn: Robust error handling, tracing, limits, retries, context cancellation

```bash
go run github.com/nexxia-ai/aigentic-examples/production@latest
```

#### [benchmark/](benchmark/)
**Performance benchmarking** - Testing agent performance
Learn: Benchmarking techniques, performance metrics

```bash
go run github.com/nexxia-ai/aigentic-examples/benchmark@latest
```

---

## Learning Path

### Beginner Track
1. **simple** - Understand basic agent creation
2. **streaming** - Learn event-driven patterns
3. **tools** - Create custom tool integrations

### Intermediate Track
4. **multi-agent** - Build coordinated agent teams
5. **memory** - Implement persistent context
6. **documents** - Process multi-modal documents

### Advanced Track
7. **approval** - Add human-in-the-loop workflows
8. **production** - Deploy production-ready agents
9. **mcp** - Integrate external protocols

---

## Example Structure

Each example includes:
- **main.go** - Fully commented, runnable code
- **README.md** - Comprehensive documentation
- **go.mod** - Independent module (easy versioning)

---

## Running Examples Locally

### Clone the Repository
```bash
git clone https://github.com/nexxia-ai/aigentic-examples
cd aigentic-examples
```

### Set Environment Variables
Create a `.env` file in the root:
```bash
OPENAI_API_KEY=your_openai_api_key_here
```

Or export directly:
```bash
export OPENAI_API_KEY=your_openai_api_key_here
```

### Run an Example
```bash
cd simple
go run main.go
```

---

## Key Concepts Covered

### Agent Configuration
- Declarative agent setup
- Model selection (OpenAI, Ollama, Gemini)
- Instructions and descriptions
- Streaming vs. blocking execution

### Tool Integration
- Custom tool creation
- Input schema definition
- Tool validation and error handling
- Human approval workflows

### Multi-Agent Coordination
- Hierarchical delegation
- Expert panels
- Team coordination patterns
- Agent-to-agent communication

### Memory Management
- Run memory (temporary)
- Session memory (persistent)
- Plan memory (complex workflows)
- Shared memory across agents

### Document Processing
- Embedded documents
- Document references (lazy loading)
- Multi-modal support (text, images, PDFs)
- Multi-document analysis

### Production Patterns
- Error handling and recovery
- Tracing and debugging
- Rate limiting and cost control
- Context cancellation
- Monitoring and observability

---

## Common Use Cases

| Use Case | Examples to Study |
|----------|-------------------|
| **Chatbot / Assistant** | simple, streaming, memory |
| **Research & Writing** | multi-agent, tools, documents |
| **Data Analysis** | tools, documents, memory |
| **Workflow Automation** | approval, multi-agent, production |
| **Content Generation** | streaming, multi-agent, memory |
| **Document Processing** | documents, tools, memory |
| **Decision Support** | multi-agent, approval, memory |

---

## Requirements

- **Go**: 1.21 or higher
- **API Keys**: OpenAI API key (or Ollama for local models)

---

## Getting Help

- 📖 **Framework Docs**: [github.com/nexxia-ai/aigentic](https://github.com/nexxia-ai/aigentic)
- 💬 **Issues**: [Report issues](https://github.com/nexxia-ai/aigentic-examples/issues)
- 🌟 **Star us**: Support the project on GitHub

---

## Contributing

We welcome contributions! To add a new example:

1. Fork the repository
2. Create a new directory with your example
3. Include `main.go`, `go.mod`, and `README.md`
4. Follow the existing example structure
5. Submit a pull request

---

## License

MIT License - see [LICENSE](LICENSE) file for details.

---

## Related Projects

- **[aigentic](https://github.com/nexxia-ai/aigentic)** - Main framework
- **[aigentic-openai](https://github.com/nexxia-ai/aigentic-openai)** - OpenAI provider
- **[aigentic-ollama](https://github.com/nexxia-ai/aigentic-ollama)** - Ollama provider
- **[aigentic-google](https://github.com/nexxia-ai/aigentic-google)** - Google Gemini provider
