package client

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"vertexai-anthropic-proxy/config"
	"vertexai-anthropic-proxy/translation"
)

func TestSendToVertexAI(t *testing.T) {
	// Create a mock server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check the request method
		if r.Method != http.MethodPost {
			t.Errorf("Expected POST request, got %s", r.Method)
		}

		// Check the request headers
		if r.Header.Get("Content-Type") != "application/json" {
			t.Errorf("Expected Content-Type: application/json, got %s", r.Header.Get("Content-Type"))
		}

		// Parse the request body
		var reqBody translation.VertexAIRequest
		err := json.NewDecoder(r.Body).Decode(&reqBody)
		if err != nil {
			t.Errorf("Error decoding request body: %v", err)
		}

		// Check the request body
		if len(reqBody.Instances) == 0 {
			t.Errorf("Expected non-empty instances in request body")
		}

		// Send a mock response
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(translation.VertexAIResponse{
			Predictions: []translation.Prediction{
				{Content: "This is a mock response from Vertex AI"},
			},
		})
	}))
	defer server.Close()

	// Create a test config
	cfg := &config.Config{
		VertexAIProjectID: "test-project",
		VertexAILocation:  "us-central1",
		VertexAIEndpoint:  server.URL, // Use the mock server URL
		AnthropicModel:    "claude-3-5-sonnet@20240620",
		AnthropicVersion:  "vertex-2023-10-16",
	}

	// Create a test request
	req := &translation.VertexAIRequest{
		Instances: []translation.Instance{
			{
				Messages: []translation.Message{
					{Role: "user", Content: "Hello, how are you?"},
				},
			},
		},
		Parameters: translation.Parameters{
			MaxOutputTokens: 100,
			Temperature:     0.7,
			TopP:            0.95,
			TopK:            40,
		},
	}

	// Send the request to the mock Vertex AI server
	resp, err := SendToVertexAI(cfg, req)
	if err != nil {
		t.Fatalf("Error sending request to Vertex AI: %v", err)
	}

	// Parse the response
	var vertexResp translation.VertexAIResponse
	err = json.Unmarshal(resp, &vertexResp)
	if err != nil {
		t.Fatalf("Error parsing Vertex AI response: %v", err)
	}

	// Check the response
	if len(vertexResp.Predictions) == 0 {
		t.Errorf("Expected non-empty predictions in response")
	}
	if vertexResp.Predictions[0].Content != "This is a mock response from Vertex AI" {
		t.Errorf("Unexpected response content: %s", vertexResp.Predictions[0].Content)
	}
}