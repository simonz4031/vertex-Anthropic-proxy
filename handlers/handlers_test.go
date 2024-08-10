package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"vertexai-anthropic-proxy/config"
	"vertexai-anthropic-proxy/translation"
	"vertexai-anthropic-proxy/utils"
)

func TestHandleMessages(t *testing.T) {
	// Initialize the logger
	utils.InitLogger("info")

	// Create a mock config
	mockConfig := &config.Config{
		VertexAIProjectID:    "test-project",
		VertexAIRegion:       "us-central1",
		VertexAIEndpoint:     "https://test-endpoint.com",
		AnthropicProxyAPIKey: "test-api-key",
	}

	// Add this at the beginning of the TestHandleMessages function
	MockVertexAIClient = func(cfg *config.Config, vertexReq *translation.VertexAIRequest) ([]byte, error) {
		// This is a mock response. In a real scenario, you'd validate the input and return appropriate responses.
		return []byte(`{"predictions": [{"content": "This is a mock response from Vertex AI"}]}`), nil
	}

	// Make sure to reset the mock at the end of the test
	defer func() {
		MockVertexAIClient = nil
	}()

	tests := []struct {
		name           string
		inputJSON      string
		expectedStatus int
		expectedBody   string
	}{
		{
			name: "Valid request",
			inputJSON: `{
				"model": "claude-v1",
				"messages": [
					{"role": "user", "content": "Hello, how are you?"}
				],
				"max_tokens": 100
			}`,
			expectedStatus: http.StatusOK,
			expectedBody:   `{"content":"This is a mock response from Vertex AI"}`,
		},
		{
			name:           "Invalid JSON",
			inputJSON:      `{"invalid": "json"`,
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "Error parsing request body",
		},
		// Add more test cases here
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a request with the test input
			req, err := http.NewRequest("POST", "/v1/messages", bytes.NewBufferString(tt.inputJSON))
			if err != nil {
				t.Fatal(err)
			}

			// Create a ResponseRecorder to record the response
			rr := httptest.NewRecorder()

			// Call the handler function
			handler := HandleMessages(mockConfig)
			handler.ServeHTTP(rr, req)

			// Check the status code
			if status := rr.Code; status != tt.expectedStatus {
				t.Errorf("handler returned wrong status code: got %v want %v", status, tt.expectedStatus)
			}

			// Check the response body
			if tt.expectedStatus == http.StatusOK {
				var got, want map[string]interface{}
				json.Unmarshal(rr.Body.Bytes(), &got)
				json.Unmarshal([]byte(tt.expectedBody), &want)

				if !jsonEqual(got, want) {
					t.Errorf("handler returned unexpected body: got %v want %v", got, want)
				}
			} else {
				// For error cases, just check if the expected message is contained in the response
				if !strings.Contains(rr.Body.String(), tt.expectedBody) {
					t.Errorf("handler returned unexpected body: got %v want it to contain %v", rr.Body.String(), tt.expectedBody)
				}
			}
		})
	}
}

// Helper function to compare JSON objects
func jsonEqual(a, b map[string]interface{}) bool {
	return string(mustMarshalJSON(a)) == string(mustMarshalJSON(b))
}

func mustMarshalJSON(v interface{}) []byte {
	b, _ := json.Marshal(v)
	return b
}
