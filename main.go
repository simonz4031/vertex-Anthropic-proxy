package main

import (
	"log"
	"net/http"

	"github.com/yourusername/vertexai-anthropic-proxy/config"
	"github.com/yourusername/vertexai-anthropic-proxy/handlers"
	"github.com/yourusername/vertexai-anthropic-proxy/utils"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	logger := utils.NewLogger(cfg.LogLevel)

	http.HandleFunc("/v1/messages", handlers.AnthropicHandler(cfg, logger))
	http.HandleFunc("/v1/messages:stream", handlers.AnthropicStreamHandler(cfg, logger))

	logger.Info("Starting server on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		logger.Fatal("Failed to start server", "error", err)
	}
}
