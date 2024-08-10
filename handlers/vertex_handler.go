package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"

	"github.com/yourusername/vertexai-anthropic-proxy/config"
	"github.com/yourusername/vertexai-anthropic-proxy/translation"
	"github.com/yourusername/vertexai-anthropic-proxy/utils"
)

func callVertexAI(cfg *config.Config, request translation.VertexAIRequest) (*translation.VertexAIResponse, error) {
	jsonData, err := json.Marshal(request)
	if err != nil {
		return nil, err
	}

	resp, err := http.Post(cfg.VertexAIEndpoint, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var vertexResponse translation.VertexAIResponse
	if err := json.NewDecoder(resp.Body).Decode(&vertexResponse); err != nil {
		return nil, err
	}

	return &vertexResponse, nil
}
