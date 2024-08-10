package client

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os/exec"
	"strings"
	"testing"

	"vertexai-anthropic-proxy/config"
)

func TestVertexAIIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	cfg := config.LoadConfig()

	// Get access token
	cmd := exec.Command("gcloud", "auth", "print-access-token")
	output, err := cmd.Output()
	if err != nil {
		t.Fatalf("Failed to get access token: %v", err)
	}
	accessToken := strings.TrimSpace(string(output))

	// Prepare request
	url := fmt.Sprintf("%s/v1/projects/%s/locations/%s/publishers/anthropic/models/%s:streamRawPredict",
		cfg.VertexAIEndpoint, cfg.VertexAIProjectID, cfg.VertexAILocation, cfg.AnthropicModel)

	reqBody := map[string]interface{}{
		"anthropic_version": cfg.AnthropicVersion,
		"messages": []map[string]string{
			{"role": "user", "content": "Hey Claude!"},
		},
		"max_tokens": 100,
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		t.Fatalf("Failed to marshal request body: %v", err)
	}

	req, err := http.NewRequestWithContext(context.Background(), "POST", url, strings.NewReader(string(jsonData)))
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	req.Header.Set("Authorization", "Bearer "+accessToken)
	req.Header.Set("Content-Type", "application/json")

	// Send request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("Failed to send request: %v", err)
	}
	defer resp.Body.Close()

	// Read response
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("Failed to read response body: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Unexpected status code: %d, body: %s", resp.StatusCode, string(body))
	}

	// Parse response
	var vertexResp map[string]interface{}
	err = json.Unmarshal(body, &vertexResp)
	if err != nil {
		t.Fatalf("Failed to parse response body: %v", err)
	}

	// Check response
	content, ok := vertexResp["content"].([]interface{})
	if !ok || len(content) == 0 {
		t.Fatalf("Unexpected response format: %v", vertexResp)
	}

	textContent, ok := content[0].(map[string]interface{})["text"].(string)
	if !ok || textContent == "" {
		t.Fatalf("Unexpected text content: %v", content[0])
	}

	t.Logf("Received response: %s", textContent)
}