package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/go-tools-agent/internal/agent"
	"github.com/go-tools-agent/internal/config"
	"github.com/go-tools-agent/internal/memory"
	"github.com/go-tools-agent/internal/parser"
	calculator "github.com/go-tools-agent/internal/tools/calculator"
	"github.com/go-tools-agent/internal/tools/code"
	httpTool "github.com/go-tools-agent/internal/tools/http"
	"github.com/go-tools-agent/internal/tools/wikipedia"
	"github.com/sashabaranov/go-openai"
	httpSwagger "github.com/swaggo/http-swagger"
)

// @title Go Tools Agent API
// @version 1.0
// @description A powerful, extensible Go-based agent framework that combines OpenAI's GPT-4 capabilities with practical tools.
// @license.name MIT
// @license.url https://opensource.org/licenses/MIT
// @host localhost:8080
// @BasePath /
type ExecuteRequest struct {
	Input string `json:"input"`
}

type ExecuteResponse struct {
	Result *agent.AgentResponse `json:"result,omitempty"`
	Error  string               `json:"error,omitempty"`
}

func main() {
	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Initialize OpenAI client
	client := openai.NewClient(cfg.OpenAIAPIKey)

	// Create memory storage
	mem := memory.NewInMemoryStorage()

	// Create output parser with a schema
	outputSchema := map[string]interface{}{
		"response": map[string]interface{}{
			"type": "string",
		},
		"confidence": map[string]interface{}{
			"type":    "number",
			"minimum": 0,
			"maximum": 1,
		},
	}
	parser := parser.NewJSONOutputParser(outputSchema)

	// Create tools
	var tools []agent.Tool

	// Add calculator tool
	name, desc, schema, handler := calculator.NewCalculatorTool()
	tools = append(tools, agent.Tool{
		Name:        name,
		Description: desc,
		Schema:      schema,
		Handler:     handler,
	})

	// Add HTTP request tool
	name, desc, schema, handler = httpTool.NewHTTPRequestTool()
	tools = append(tools, agent.Tool{
		Name:        name,
		Description: desc,
		Schema:      schema,
		Handler:     handler,
	})

	// Add Wikipedia tool
	name, desc, schema, handler = wikipedia.NewWikipediaTool()
	tools = append(tools, agent.Tool{
		Name:        name,
		Description: desc,
		Schema:      schema,
		Handler:     handler,
	})

	// Add code execution tool
	name, desc, schema, handler = code.NewCodeExecutionTool()
	tools = append(tools, agent.Tool{
		Name:        name,
		Description: desc,
		Schema:      schema,
		Handler:     handler,
	})

	// Configure the agent
	agentConfig := agent.AgentConfig{
		SystemMessage:           "You are a helpful assistant that can perform calculations, make HTTP requests, search Wikipedia, and execute code.",
		MaxIterations:           5,
		ReturnIntermediateSteps: true,
		Tools:                   tools,
	}

	// Create the agent
	toolsAgent := agent.NewToolsAgent(agentConfig, client, mem, parser)

	// Serve Swagger documentation
	http.HandleFunc("/swagger/doc.json", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "api/swagger.yaml")
	})
	http.HandleFunc("/swagger/", httpSwagger.Handler(
		httpSwagger.URL("/swagger/doc.json"),
	))

	// Create HTTP handler for execute endpoint
	http.HandleFunc("/execute", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var req ExecuteRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		// Create context with timeout
		ctx, cancel := context.WithTimeout(r.Context(), 30*time.Second)
		defer cancel()

		// Execute the agent
		response, err := toolsAgent.Execute(ctx, req.Input)

		// Prepare response
		executeResponse := ExecuteResponse{}
		if err != nil {
			executeResponse.Error = err.Error()
		} else {
			executeResponse.Result = response
		}

		// Send response
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(executeResponse); err != nil {
			http.Error(w, "Failed to encode response", http.StatusInternalServerError)
			return
		}
	})

	// Start server
	log.Printf("Server started:\n- API: http://localhost:8080\n- Swagger UI: http://localhost:8080/swagger/")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
