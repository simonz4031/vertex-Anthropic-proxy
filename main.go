package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"vertexai-anthropic-proxy/config"
	"vertexai-anthropic-proxy/handlers"
	"vertexai-anthropic-proxy/middleware"
	"vertexai-anthropic-proxy/utils"
)

func main() {
	cfg := config.LoadConfig()

	// Initialize logger
	utils.InitLogger("info")
	logger := utils.GetLogger()

	// Set log flags to include file name and line number
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	// Set up routes with middleware
	http.HandleFunc("/v1/messages", middleware.AuthMiddleware(cfg)(handlers.HandleMessages(cfg)))
	http.HandleFunc("/v1/chat/completions", middleware.AuthMiddleware(cfg)(handlers.HandleMessages(cfg))) // Using the same handler for now
	http.HandleFunc("/set-log-level", handlers.HandleSetLogLevel)
	http.HandleFunc("/refresh-credentials", handlers.HandleRefreshCredentials)

	// Log configuration
	logger.Infof("Starting server with configuration:")
	logger.Infof("Vertex AI Project ID: %s", cfg.VertexAIProjectID)
	logger.Infof("Vertex AI Region: %s", cfg.VertexAIRegion)
	logger.Infof("Vertex AI Endpoint: %s", cfg.VertexAIEndpoint)
	logger.Infof("ANTHROPIC_API_KEY: %s", cfg.AnthropicProxyAPIKey[:5]+"...") // Log only the first 5 characters for security

	// Get port from environment variable
	port := os.Getenv("PORT")
	if port == "" {
		port = "8070" // Default port if not specified
	}

	// Start server
	addr := fmt.Sprintf(":%s", port)
	logger.Infof("Server listening on %s", addr)
	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}