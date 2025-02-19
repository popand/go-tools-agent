package agent

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/sashabaranov/go-openai"
)

// ToolsAgent represents the main agent implementation
type ToolsAgent struct {
	config AgentConfig
	client *openai.Client
	memory Memory
	parser OutputParser
}

// NewToolsAgent creates a new instance of ToolsAgent
func NewToolsAgent(config AgentConfig, client *openai.Client, memory Memory, parser OutputParser) *ToolsAgent {
	return &ToolsAgent{
		config: config,
		client: client,
		memory: memory,
		parser: parser,
	}
}

// Execute runs the agent with the given input
func (a *ToolsAgent) Execute(ctx context.Context, input string) (*AgentResponse, error) {
	var steps []AgentStep
	var finalOutput json.RawMessage

	// Load memory if available
	var memoryContent []byte
	if a.memory != nil {
		var err error
		memoryContent, err = a.memory.LoadMemory(ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to load memory: %w", err)
		}
	}

	// Prepare messages for the model
	messages := []openai.ChatCompletionMessage{
		{
			Role:    openai.ChatMessageRoleSystem,
			Content: a.config.SystemMessage,
		},
		{
			Role:    openai.ChatMessageRoleUser,
			Content: input,
		},
	}

	// Add memory context if available
	if len(memoryContent) > 0 {
		messages = append(messages, openai.ChatCompletionMessage{
			Role:    openai.ChatMessageRoleSystem,
			Content: fmt.Sprintf("Previous context: %s", string(memoryContent)),
		})
	}

	// Main execution loop
	for iteration := 0; iteration < a.config.MaxIterations; iteration++ {
		// Check context cancellation
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
		}

		// Prepare tool choices for the model
		tools := make([]openai.Tool, len(a.config.Tools))
		for i, tool := range a.config.Tools {
			tools[i] = openai.Tool{
				Type: openai.ToolTypeFunction,
				Function: openai.FunctionDefinition{
					Name:        tool.Name,
					Description: tool.Description,
					Parameters:  tool.Schema,
				},
			}
		}

		// Create chat completion request
		req := openai.ChatCompletionRequest{
			Model:    openai.GPT4TurboPreview,
			Messages: messages,
			Tools:    tools,
		}

		// Get model response
		resp, err := a.client.CreateChatCompletion(ctx, req)
		if err != nil {
			return nil, fmt.Errorf("failed to get model response: %w", err)
		}

		// Process tool calls
		if len(resp.Choices) > 0 && len(resp.Choices[0].Message.ToolCalls) > 0 {
			for _, toolCall := range resp.Choices[0].Message.ToolCalls {
				step := AgentStep{
					Action:    toolCall.Function.Name,
					Input:     json.RawMessage(toolCall.Function.Arguments),
					Timestamp: time.Now().Unix(),
				}

				// Find and execute the tool
				var toolOutput json.RawMessage
				for _, tool := range a.config.Tools {
					if tool.Name == toolCall.Function.Name {
						output, err := tool.Handler(ctx, step.Input)
						if err != nil {
							step.Error = err.Error()
						} else {
							step.Output = output
							toolOutput = output
						}
						break
					}
				}

				steps = append(steps, step)

				// Add tool result to messages
				messages = append(messages, openai.ChatCompletionMessage{
					Role:    openai.ChatMessageRoleAssistant,
					Content: "",
					ToolCalls: []openai.ToolCall{
						{
							ID:   toolCall.ID,
							Type: openai.ToolTypeFunction,
							Function: openai.FunctionCall{
								Name:      toolCall.Function.Name,
								Arguments: string(step.Input),
							},
						},
					},
				})

				if toolOutput != nil {
					messages = append(messages, openai.ChatCompletionMessage{
						Role:       openai.ChatMessageRoleTool,
						Content:    string(toolOutput),
						Name:       toolCall.Function.Name,
						ToolCallID: toolCall.ID,
					})
				} else {
					messages = append(messages, openai.ChatCompletionMessage{
						Role:       openai.ChatMessageRoleTool,
						Content:    "Error: Tool execution failed",
						Name:       toolCall.Function.Name,
						ToolCallID: toolCall.ID,
					})
				}
			}
		} else {
			// No more tool calls, we have the final output
			if len(resp.Choices) > 0 {
				// Ensure the output is in JSON format
				finalOutput = json.RawMessage(fmt.Sprintf(`{"response": %q, "confidence": 1.0}`, resp.Choices[0].Message.Content))
				break
			}
		}
	}

	// Parse final output if parser is available
	if a.parser != nil && len(finalOutput) > 0 {
		parsedOutput, err := a.parser.Parse(finalOutput)
		if err != nil {
			return nil, fmt.Errorf("failed to parse output: %w", err)
		}
		finalOutput = parsedOutput
	}

	// Save to memory if available
	if a.memory != nil && len(finalOutput) > 0 {
		if err := a.memory.SaveMemory(ctx, finalOutput); err != nil {
			return nil, fmt.Errorf("failed to save memory: %w", err)
		}
	}

	response := &AgentResponse{
		FinalOutput: finalOutput,
	}

	if a.config.ReturnIntermediateSteps {
		response.Steps = steps
	}

	return response, nil
}
