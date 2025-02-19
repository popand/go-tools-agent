package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// Config holds all configuration values
type Config struct {
	OpenAIAPIKey string
}

// LoadConfig loads configuration from environment variables and .env file
func LoadConfig() (*Config, error) {
	// Try to load .env file if it exists
	if err := loadEnvFile(); err != nil {
		// Don't return error as .env file is optional
		fmt.Printf("Note: %v\n", err)
	}

	// Get API key from environment
	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		return nil, fmt.Errorf("OPENAI_API_KEY environment variable is required")
	}

	return &Config{
		OpenAIAPIKey: apiKey,
	}, nil
}

// loadEnvFile loads environment variables from .env file
func loadEnvFile() error {
	// Get the current working directory
	dir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get working directory: %w", err)
	}

	// Look for .env file
	envPath := filepath.Join(dir, ".env")
	content, err := os.ReadFile(envPath)
	if err != nil {
		return fmt.Errorf("failed to read .env file: %w", err)
	}

	// Parse each line
	lines := strings.Split(string(content), "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}

		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])

		// Remove quotes if present
		value = strings.Trim(value, `"'`)

		// Set environment variable if not already set
		if os.Getenv(key) == "" {
			os.Setenv(key, value)
		}
	}

	return nil
} 