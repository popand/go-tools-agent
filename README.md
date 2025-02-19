# Go Tools Agent

A Go implementation of a Tools Agent system, similar to the n8n LangChain Tools Agent. This implementation provides a framework for creating AI agents that can use tools to accomplish tasks.

## Features

- Tool-based agent system
- Memory management
- Output parsing and validation
- Context-aware execution
- Configurable system messages
- Support for multiple tools
- JSON schema validation
- Environment-based configuration

## Requirements

- Go 1.21 or later
- OpenAI API key

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
# SYSTEM_MESSAGE="You are a helpful assistant that can perform calculations."
```

## Usage

1. Run the example:
```bash
go run cmd/main.go
```

The application will automatically load the configuration from your `.env` file.

## Example Output

```json
{
  "final_output": {
    "response": "I'll help you with that calculation. First, I'll divide 15 by 3, which gives us 5. Then, I'll multiply that result by 4, which gives us 20.",
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
      "timestamp": 1645567890
    },
    {
      "action": "calculator",
      "input": {
        "operation": "multiply",
        "a": 5,
        "b": 4
      },
      "output": {
        "result": 20
      },
      "timestamp": 1645567891
    }
  ]
}
```

## Adding New Tools

To add a new tool:

1. Create a new file in the `internal/tools` directory
2. Implement the tool following the pattern in `calculator.go`
3. Add the tool to the agent configuration in `main.go`

Example:
```go
// Create your tool
name, desc, schema, handler := tools.NewYourTool()
yourTool := agent.Tool{
    Name:        name,
    Description: desc,
    Schema:      schema,
    Handler:     handler,
}

// Add to config
config := agent.AgentConfig{
    // ...
    Tools: []agent.Tool{calculatorTool, yourTool},
}
```

## Architecture

The system consists of several components:

- **Agent**: Core logic for tool selection and execution
- **Memory**: State management system
- **Parser**: Output formatting and validation
- **Tools**: Individual tool implementations
- **Config**: Environment and configuration management

## Security Notes

- Never commit your `.env` file to version control
- Keep your API keys secure and rotate them regularly
- Use environment-specific `.env` files for different deployments

## License

MIT License

## Contributing

1. Fork the repository
2. Create your feature branch
3. Commit your changes
4. Push to the branch
5. Create a new Pull Request 