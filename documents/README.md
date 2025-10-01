# Document Processing Example

This example demonstrates how to work with documents in aigentic agents by analyzing an employment contract.

## What You'll Learn

- Embedding documents directly in agent context
- Processing text documents for information extraction
- Extracting structured data from documents
- Best practices for document analysis

## Example Demonstrated

### Contract Analysis
Analyze an employment contract by embedding the full text in the agent's context and extracting key information.

**Use case**: Contract analysis, legal document review, information extraction

## Running the Example

```bash
# Set your OpenAI API key
export OPENAI_API_KEY=your_api_key_here

# Run from the examples directory
go run github.com/nexxia-ai/aigentic-examples/documents@latest

# Or run locally
cd documents
go run main.go
```

## Key Concepts

### Documents vs DocumentReferences

The aigentic framework provides two ways to work with documents:

#### Documents (Embedded)

```go
agent := aigentic.Agent{
    Documents: []*document.Document{doc1, doc2},
}
```

**Characteristics**:
- Document content is embedded directly in the agent's context
- Content is always available to the LLM
- Counts against the model's context window
- Best for small to medium documents

**When to use**:
- Document is small enough to fit comfortably in context (< 10KB)
- Agent needs full access to document content
- Processing a single or few documents
- Analyzing structured data from documents
- Examples: contracts, receipts, short reports

#### DocumentReferences (On-Demand)

```go
agent := aigentic.Agent{
    DocumentReferences: []*document.Document{doc1, doc2, doc3, doc4},
}
```

**Characteristics**:
- Documents are referenced but not embedded
- LLM uses built-in tools to retrieve content when needed
- Only fetches relevant portions
- Doesn't count against initial context window
- Best for large document collections

**When to use**:
- Multiple large documents (knowledge base, documentation)
- Agent should decide what to retrieve
- Context window optimization is important
- Search-like interactions across documents
- Examples: technical documentation, FAQ databases, product manuals

### Document Creation Methods

#### 1. In-Memory Documents

Create documents from in-memory data:

```go
doc := document.NewInMemoryDocument(
    "doc_id",           // Unique identifier
    "filename.txt",     // Display name
    []byte(content),    // Document content as bytes
    nil,                // Source document (for chunks)
)
```

**Use cases**:
- Generated content
- API responses
- Temporary documents
- Testing

#### 2. Local File Store

Load documents from filesystem with lazy loading:

```go
store := document.NewLocalStore("/path/to/documents")
doc, err := store.Open(ctx, "report.pdf")

// Content is loaded only when accessed
content, err := doc.Bytes()
```

**Use cases**:
- Existing document repositories
- File-based workflows
- Large documents (lazy loading)
- Production systems

#### 3. Custom Document Stores

Implement the `DocumentStore` interface for custom backends:

```go
type DocumentStore interface {
    Open(ctx context.Context, filePath string) (*Document, error)
    Close(ctx context.Context) error
}
```

**Use cases**:
- Cloud storage (S3, GCS, Azure Blob)
- Databases
- CMS systems
- Custom document management systems

### Document Properties

Each document has the following properties:

```go
type Document struct {
    Filename    string    // Original filename
    FilePath    string    // Path or identifier
    FileSize    int64     // Size in bytes
    MimeType    string    // MIME type (e.g., "text/plain", "image/png")
    CreatedAt   time.Time // Creation timestamp

    // Chunking metadata (when document is split)
    SourceDoc   *Document // Original document (if this is a chunk)
    ChunkIndex  int       // Index of this chunk
    TotalChunks int       // Total number of chunks
    StartChar   int       // Start position in source
    EndChar     int       // End position in source
    PageNumber  int       // Page number (for PDFs)
}
```

## Document Processing Patterns

### Pattern 1: Single Document Analysis

```go
doc := document.NewInMemoryDocument("id", "contract.txt", data, nil)

agent := aigentic.Agent{
    Documents: []*document.Document{doc},
    Instructions: "Analyze this contract and extract key terms",
}
```

**Best for**: Focused analysis of a single document

### Pattern 2: Multi-Document Comparison

```go
agent := aigentic.Agent{
    Documents: []*document.Document{doc1, doc2, doc3},
    Instructions: "Compare these reports and identify trends",
}
```

**Best for**: Comparative analysis, trend identification

### Pattern 3: Document Library with Search

```go
agent := aigentic.Agent{
    DocumentReferences: []*document.Document{docs...},
    Instructions: "Search the documentation to answer user questions",
}
```

**Best for**: Knowledge bases, documentation assistants

### Pattern 4: Mixed Approach

```go
agent := aigentic.Agent{
    Documents: []*document.Document{priorityDoc},
    DocumentReferences: []*document.Document{referenceDoc1, referenceDoc2},
    Instructions: "Focus on the priority document, consult references as needed",
}
```

**Best for**: Main document with supporting references

## Supported Document Types

### Text Documents
- Plain text (`.txt`)
- Markdown (`.md`)
- JSON (`.json`)
- CSV (`.csv`)
- Source code files

### Images
- PNG (`.png`)
- JPEG (`.jpg`, `.jpeg`)
- GIF (`.gif`)
- WebP (`.webp`)

**Note**: Requires vision-capable models (e.g., GPT-4 Vision, Claude with vision)

### PDFs
- PDF documents (`.pdf`)
- Multi-page support
- Chunking for large PDFs

### Office Documents
- Word documents (`.docx`)
- Excel spreadsheets (`.xlsx`)
- PowerPoint presentations (`.pptx`)

**Note**: May require preprocessing to extract text

## Best Practices

### 1. Choose the Right Approach

| Scenario | Recommendation |
|----------|----------------|
| Single small document (< 10KB) | `Documents` (embedded) |
| Multiple small documents (< 5 docs) | `Documents` (embedded) |
| Large document collection (> 5 docs) | `DocumentReferences` |
| Document library / knowledge base | `DocumentReferences` |
| Need full document context | `Documents` (embedded) |
| Agent should search/filter | `DocumentReferences` |

### 2. Optimize Context Usage

```go
// BAD: Embedding too many large documents
agent := aigentic.Agent{
    Documents: []*document.Document{
        largeDoc1, // 100KB
        largeDoc2, // 150KB
        largeDoc3, // 200KB
    },
}

// GOOD: Use references for large collections
agent := aigentic.Agent{
    DocumentReferences: []*document.Document{
        largeDoc1, largeDoc2, largeDoc3,
    },
}
```

### 3. Provide Clear Instructions

```go
agent := aigentic.Agent{
    Documents: []*document.Document{contractDoc},
    Instructions: `You are a legal assistant.
        Analyze employment contracts focusing on:
        - Compensation and benefits
        - Term and termination clauses
        - Non-compete agreements
        Extract information clearly and cite specific clauses.`,
}
```

### 4. Handle Errors Gracefully

```go
doc, err := store.Open(ctx, "report.pdf")
if err != nil {
    log.Printf("Failed to load document: %v", err)
    // Handle error appropriately
    return
}

content, err := doc.Bytes()
if err != nil {
    log.Printf("Failed to read document content: %v", err)
    // Handle error
    return
}
```

### 5. Use Lazy Loading for Large Files

```go
// LocalStore loads content only when needed
store := document.NewLocalStore("/path/to/docs")
doc, _ := store.Open(ctx, "large_file.pdf")

// Content not loaded yet - metadata only
fmt.Printf("Size: %d bytes\n", doc.FileSize)

// Content loaded on first access
content, _ := doc.Bytes()
```

### 6. Structure Your Queries

```go
// BAD: Vague query
response, _ := agent.Execute("What's in this document?")

// GOOD: Specific query
response, _ := agent.Execute(
    "Extract the total amount, date, and payment method from this receipt.",
)
```

## Common Use Cases

### Legal Document Analysis

```go
agent := aigentic.Agent{
    Name: "LegalAssistant",
    Instructions: "Analyze legal documents and extract key clauses",
    Documents: []*document.Document{contractDoc},
}

response, _ := agent.Execute("Summarize the termination conditions")
```

### Receipt/Invoice Processing

```go
agent := aigentic.Agent{
    Name: "ReceiptProcessor",
    Instructions: "Extract structured data from receipts and invoices",
    Documents: []*document.Document{receiptImage},
}

response, _ := agent.Execute("Extract all line items and the total amount")
```

### Technical Documentation Assistant

```go
agent := aigentic.Agent{
    Name: "DocsAssistant",
    Instructions: "Help users find information in technical documentation",
    DocumentReferences: []*document.Document{apiDocs, userGuide, faq},
}

response, _ := agent.Execute("How do I authenticate API requests?")
```

### Financial Report Analysis

```go
agent := aigentic.Agent{
    Name: "FinancialAnalyst",
    Instructions: "Analyze financial reports and identify trends",
    Documents: []*document.Document{q1Report, q2Report, q3Report},
}

response, _ := agent.Execute("What are the revenue trends across quarters?")
```

### Document Summarization

```go
agent := aigentic.Agent{
    Name: "Summarizer",
    Instructions: "Create concise summaries of documents",
    Documents: []*document.Document{longReport},
}

response, _ := agent.Execute("Provide a 3-paragraph summary of this report")
```

## Performance Considerations

### Context Window Management

Most LLMs have context limits (e.g., GPT-4: 128K tokens, Claude: 200K tokens).

**Rule of thumb**: 1 token ≈ 4 characters for English text

```go
// A 50KB document ≈ 12,500 tokens
// A 500KB document ≈ 125,000 tokens (near GPT-4 limit)
```

### Document Size Guidelines

| Document Size | Approach | Notes |
|---------------|----------|-------|
| < 10KB | Embed | Fits easily in context |
| 10-100KB | Embed (carefully) | Monitor context usage |
| 100KB-1MB | Reference | Let LLM retrieve sections |
| > 1MB | Chunk or Reference | Consider preprocessing |

### Optimization Strategies

1. **Chunking**: Split large documents into smaller pieces
2. **Preprocessing**: Extract text from complex formats (PDF, DOCX)
3. **Summarization**: Create summaries for very large documents
4. **Selective loading**: Use references to load only relevant sections
5. **Caching**: Cache frequently accessed documents

## Troubleshooting

### Document Not Loading

```go
// Check if loader is set
if doc.loader == nil {
    log.Fatal("Document has no loader")
}

// Verify file exists
if _, err := os.Stat(filePath); os.IsNotExist(err) {
    log.Fatal("File does not exist")
}
```

### Context Window Exceeded

```
Error: context length exceeded
```

**Solution**: Use `DocumentReferences` instead of `Documents`

### Image Not Recognized

**Issue**: Using text-only model with image documents

**Solution**: Use vision-capable models:
- OpenAI: `gpt-4-vision-preview`, `gpt-4o`
- Anthropic: Claude 3 models (Opus, Sonnet, Haiku)

### Poor Extraction Quality

**Issue**: Agent not extracting information accurately

**Solutions**:
1. Provide more specific instructions
2. Include examples in your query
3. Ask for structured output format (JSON, table)
4. Use a more capable model

## Next Steps

- See [tools example](../tools) for combining documents with custom tools
- See [multi-agent example](../multi-agent) for specialized document processing teams
- See [memory example](../memory) for maintaining document context across sessions
- See [production example](../production) for error handling and monitoring

## Additional Resources

- [Aigentic Documentation](https://github.com/nexxia-ai/aigentic)
- [Document Package Reference](https://github.com/nexxia-ai/aigentic/tree/main/document)
- [Vision Models Guide](https://platform.openai.com/docs/guides/vision)
