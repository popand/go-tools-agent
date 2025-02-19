# Go Tools Agent

A Go implementation of a Tools Agent system, similar to the n8n LangChain Tools Agent. This implementation provides a framework for creating AI agents that can use tools to accomplish tasks.

## Features

- Tool-based agent system with support for multiple tools:
  - Calculator: Perform basic mathematical operations
  - HTTP Request: Make HTTP requests to external APIs
  - Wikipedia: Search and retrieve information from Wikipedia
  - Code Execution: Execute code snippets in various languages
- Advanced JSON schema validation with type checking and range validation
- Memory management for context persistence
- Output parsing and validation
- Context-aware execution with timeout handling
- Configurable system messages
- Environment-based configuration
- Proper error handling and reporting

## Requirements

- Go 1.21 or later
- OpenAI API key (GPT-4 access required)
- Python (for code execution tool)

## Installation

1. Clone the repository:
```bash
git clone https://github.com/yourusername/go-tools-agent.git
cd go-tools-agent
```

2. Install dependencies:
```bash
go mod download
```

3. Set up configuration:
```bash
# Copy the example environment file
cp .env.example .env

# Edit .env with your configuration
nano .env  # or use your preferred editor
```

The `.env` file should contain:
```env
# Required
OPENAI_API_KEY=your-api-key-here

# Optional
# MAX_ITERATIONS=5
# SYSTEM_MESSAGE="You are a helpful assistant that can perform calculations, make HTTP requests, search Wikipedia, and execute code."
```

## Usage

1. Run the example:
```bash
go run cmd/main.go
```

The application will automatically load the configuration from your `.env` file and execute a sample task that demonstrates all available tools.

## Example Output

```json
{
  "final_output": {
    "response": "Here are the results for your tasks:\n\n1. The result of 15 divided by 3 and multiplied by 4 is 20.\n\n2. Information about Go programming language from Wikipedia...\n\n3. GitHub repository information: golang/go has 125,938 stars and 17,886 forks.\n\n4. Python script output: Hello from Python!",
    "confidence": 1.0
  },
  "steps": [
    {
      "action": "calculator",
      "input": {
        "operation": "divide",
        "a": 15,
        "b": 3
      },
      "output": {
        "result": 5
      },
      "timestamp": 1739983078
    },
    // Additional steps omitted for brevity...
  ]
}
```

## Available Tools

### Calculator Tool
- Performs basic mathematical operations (add, subtract, multiply, divide)
- Input validation and error handling
- JSON schema-compliant input/output

### HTTP Request Tool
- Supports GET, POST, PUT, DELETE methods
- Headers and query parameters support
- Response includes status code, headers, and body

### Wikipedia Tool
- Search Wikipedia articles
- Retrieve article content and metadata
- Returns article URL and extract

### Code Execution Tool
- Execute code snippets in various languages
- Currently supports Python
- Captures output and exit code
- Secure execution environment

## Adding New Tools

To add a new tool:

1. Create a new file in the `internal/tools` directory
2. Implement the tool following the pattern in existing tools
3. Add the tool to the agent configuration in `main.go`

Example:
```go
// Create your tool
name, desc, schema, handler := tools.NewYourTool()
tools = append(tools, agent.Tool{
    Name:        name,
    Description: desc,
    Schema:      schema,
    Handler:     handler,
})
```

## Architecture

The system consists of several components:

- **Agent**: Core logic for tool selection and execution
  - Manages conversation with OpenAI API
  - Handles tool calls and responses
  - Maintains conversation context
- **Memory**: State management system
  - Persists conversation history
  - Maintains context between calls
- **Parser**: Output formatting and validation
  - JSON schema validation
  - Type checking and range validation
  - Consistent output formatting
- **Tools**: Individual tool implementations
  - Standardized interface
  - Input validation
  - Error handling
- **Config**: Environment and configuration management
  - Environment variable loading
  - Default configuration
  - Runtime configuration

## Security Notes

- Never commit your `.env` file to version control
- Keep your API keys secure and rotate them regularly
- Use environment-specific `.env` files for different deployments
- Be cautious with the code execution tool in production environments
- Implement rate limiting for API calls
- Validate and sanitize all inputs

## License

MIT License

## Contributing

1. Fork the repository
2. Create your feature branch
3. Commit your changes
4. Push to the branch
5. Create a new Pull Request 