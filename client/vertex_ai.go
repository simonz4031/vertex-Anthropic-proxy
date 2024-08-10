package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"

	"vertexai-anthropic-proxy/config"
	"vertexai-anthropic-proxy/translation"
)

func SendToVertexAI(cfg *config.Config, req *translation.VertexAIRequest) (io.ReadCloser, error) {
	ctx := context.Background()

	credentials, err := google.FindDefaultCredentials(ctx, "https://www.googleapis.com/auth/cloud-platform")
	if err != nil {
		log.Printf("Error finding default credentials: %v", err)
		return nil, err
	}

	client := oauth2.NewClient(ctx, credentials.TokenSource)

	url := fmt.Sprintf("%s/v1/projects/%s/locations/%s/publishers/anthropic/models/%s:streamRawPredict",
		cfg.VertexAIEndpoint, cfg.VertexAIProjectID, cfg.VertexAIRegion, cfg.AnthropicModel)

	jsonData, err := json.Marshal(req)
	if err != nil {
		log.Printf("Error marshaling request: %v", err)
		return nil, err
	}

	log.Printf("Sending request to Vertex AI: %s", url)
	log.Printf("Request body: %s", string(jsonData))

	httpReq, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		log.Printf("Error creating request: %v", err)
		return nil, err
	}

	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(httpReq)
	if err != nil {
		log.Printf("Error sending request: %v", err)
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		resp.Body.Close()
		log.Printf("Vertex AI returned non-OK status: %d, body: %s", resp.StatusCode, string(body))
		return nil, fmt.Errorf("Vertex AI returned non-OK status: %d, body: %s", resp.StatusCode, string(body))
	}

	return resp.Body, nil
}