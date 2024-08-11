package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	VertexAIProjectID    string
	VertexAIRegion       string
	VertexAIEndpoint     string
	AnthropicModel       string
	AnthropicProxyAPIKey string
	OpenAIProxyAPIKey    string
}

func LoadConfig() *Config {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	cfg := &Config{
		VertexAIProjectID:    os.Getenv("VERTEX_AI_PROJECT_ID"),
		VertexAIRegion:       os.Getenv("VERTEX_AI_REGION"),
		VertexAIEndpoint:     os.Getenv("VERTEX_AI_ENDPOINT"),
		AnthropicModel:       os.Getenv("MODEL"),
		AnthropicProxyAPIKey: os.Getenv("ANTHROPIC_PROXY_API_KEY"),
		OpenAIProxyAPIKey:    os.Getenv("OPENAI_PROXY_API_KEY"),
	}

	if cfg.VertexAIEndpoint == "" {
		log.Fatal("VERTEX_AI_ENDPOINT is not set in the environment")
	}

	return cfg
}