package handlers

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"vertexai-anthropic-proxy/config"
	"vertexai-anthropic-proxy/translation"
	"vertexai-anthropic-proxy/client"
	"vertexai-anthropic-proxy/utils"
)

// MockVertexAIClient is used for testing
var MockVertexAIClient func(cfg *config.Config, vertexReq *translation.VertexAIRequest) ([]byte, error)

func HandleMessages(cfg *config.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logger := utils.GetLogger()

		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			logger.Errorf("Error reading request body: %v", err)
			http.Error(w, "Error reading request body", http.StatusBadRequest)
			return
		}

		var anthropicReq translation.AnthropicRequest
		if err := json.Unmarshal(body, &anthropicReq); err != nil {
			logger.Errorf("Error parsing request body: %v", err)
			http.Error(w, "Error parsing request body", http.StatusBadRequest)
			return
		}

		logger.Infof("Received request: %+v", anthropicReq)

		vertexReq, err := translation.AnthropicToVertexAI(anthropicReq)
		if err != nil {
			logger.Errorf("Error translating request: %v", err)
			http.Error(w, "Error translating request", http.StatusInternalServerError)
			return
		}

		resp, err := client.SendToVertexAI(cfg, &vertexReq)
		if err != nil {
			logger.Errorf("Error sending request to Vertex AI: %v", err)
			http.Error(w, "Error processing request", http.StatusInternalServerError)
			return
		}

		var vertexAIResponse translation.VertexAIResponse
		if err := json.Unmarshal(resp, &vertexAIResponse); err != nil {
			logger.Errorf("Error parsing Vertex AI response: %v", err)
			http.Error(w, "Error processing response", http.StatusInternalServerError)
			return
		}

		anthropicResp, err := translation.VertexAIToAnthropic(vertexAIResponse)
		if err != nil {
			logger.Errorf("Error translating Vertex AI response: %v", err)
			http.Error(w, "Error processing response", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(anthropicResp)
	}
}

func HandleChatCompletions(cfg *config.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		// For now, we'll use the same logic as HandleMessages
		// In the future, you might want to implement specific chat completion logic
		HandleMessages(cfg).ServeHTTP(w, r)
	}
}