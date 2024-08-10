package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	VertexAIEndpoint string
	AnthropicAPIKey  string
	LogLevel         string
}

func Load() (*Config, error) {
	if err := godotenv.Load(); err != nil {
		return nil, err
	}

	return &Config{
		VertexAIEndpoint: os.Getenv("VERTEX_AI_ENDPOINT"),
		AnthropicAPIKey:  os.Getenv("ANTHROPIC_API_KEY"),
		LogLevel:         getEnv("LOG_LEVEL", "info"),
	}, nil
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
