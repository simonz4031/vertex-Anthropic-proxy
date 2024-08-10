package config

import (
    "fmt"
    "log"
    "os"

    "github.com/joho/godotenv"
)

type Config struct {
    Port               string
    VertexAIProjectID  string
    VertexAIRegion     string
    VertexAIEndpoint   string
    AnthropicAPIKey    string
    Model              string
}

func LoadConfig() *Config {
    err := godotenv.Load()
    if err != nil {
        log.Printf("Error loading .env file: %v", err)
    }

    cfg := &Config{
        Port:              getEnvOrDefault("PORT", "8070"),
        VertexAIProjectID: getEnvOrFatal("VERTEX_AI_PROJECT_ID"),
        VertexAIRegion:    getEnvOrFatal("VERTEX_AI_REGION"),
        AnthropicAPIKey:   getEnvOrFatal("ANTHROPIC_API_KEY"),
        Model:             getEnvOrFatal("MODEL"),
    }

    cfg.VertexAIEndpoint = fmt.Sprintf(
        "https://%s-aiplatform.googleapis.com/v1/projects/%s/locations/%s/publishers/anthropic/models/%s:streamRawPredict",
        cfg.VertexAIRegion,
        cfg.VertexAIProjectID,
        cfg.VertexAIRegion,
        cfg.Model,
    )

    log.Printf("Loaded configuration: %+v", cfg)

    return cfg
}

func getEnvOrDefault(key, defaultValue string) string {
    value, exists := os.LookupEnv(key)
    if !exists {
        return defaultValue
    }
    return value
}

func getEnvOrFatal(key string) string {
    value, exists := os.LookupEnv(key)
    if !exists {
        log.Fatalf("Environment variable %s is not set", key)
    }
    return value
}