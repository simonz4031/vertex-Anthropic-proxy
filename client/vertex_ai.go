package client

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"

	"vertexai-anthropic-proxy/config"
	"vertexai-anthropic-proxy/translation"
	"github.com/google/uuid"
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

	resp, err := client.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		log.Printf("Error sending request to Vertex AI: %v", err)
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

func SendToVertexAIStream(cfg *config.Config, req *translation.VertexAIRequest, responseChan chan<- []byte) error {
	ctx := context.Background()

	credentials, err := google.FindDefaultCredentials(ctx, "https://www.googleapis.com/auth/cloud-platform")
	if err != nil {
		log.Printf("Error finding default credentials: %v", err)
		return err
	}

	client := oauth2.NewClient(ctx, credentials.TokenSource)

	url := fmt.Sprintf("%s/v1/projects/%s/locations/%s/publishers/anthropic/models/%s:streamRawPredict",
		cfg.VertexAIEndpoint, cfg.VertexAIProjectID, cfg.VertexAIRegion, cfg.AnthropicModel)

	jsonData, err := json.Marshal(req)
	if err != nil {
		log.Printf("Error marshaling request: %v", err)
		return err
	}

	log.Printf("Sending streaming request to Vertex AI: %s", url)
	log.Printf("Request body: %s", string(jsonData))

	resp, err := client.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		log.Printf("Error sending request to Vertex AI: %v", err)
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		log.Printf("Vertex AI returned non-OK status: %d, body: %s", resp.StatusCode, string(body))
		return fmt.Errorf("Vertex AI returned non-OK status: %d, body: %s", resp.StatusCode, string(body))
	}

	reader := bufio.NewReader(resp.Body)
	var currentEvent string
	var currentData []byte

	for {
		line, err := reader.ReadBytes('\n')
		if err != nil {
			if err == io.EOF {
				break
			}
			log.Printf("Error reading response: %v", err)
			return err
		}

		// Trim whitespace and newline characters
		line = bytes.TrimSpace(line)

		// Skip empty lines
		if len(line) == 0 {
			continue
		}

		// Check for SSE prefix
		if bytes.HasPrefix(line, []byte("event: ")) {
			currentEvent = string(bytes.TrimPrefix(line, []byte("event: ")))
		} else if bytes.HasPrefix(line, []byte("data: ")) {
			currentData = bytes.TrimPrefix(line, []byte("data: "))

			// Process the event
			switch currentEvent {
			case "message_start", "ping", "content_block_start", "content_block_stop":
				// Ignore these events
			case "content_block_delta":
				// Parse the JSON content
				var event struct {
					Delta struct {
						Text string `json:"text"`
					} `json:"delta"`
				}
				if err := json.Unmarshal(currentData, &event); err != nil {
					log.Printf("Error parsing JSON: %v", err)
					continue
				}

				// Format as OpenAI-compatible JSON event
				openAIEvent := map[string]interface{}{
					"id":      "chatcmpl-" + uuid.New().String(),
					"object":  "chat.completion.chunk",
					"created": time.Now().Unix(),
					"model":   "gpt-3.5-turbo-0613", // or whatever model name you want to use
					"choices": []map[string]interface{}{
						{
							"delta": map[string]string{
								"content": event.Delta.Text,
							},
							"index":        0,
							"finish_reason": nil,
						},
					},
				}

				jsonData, err := json.Marshal(openAIEvent)
				if err != nil {
					log.Printf("Error marshaling JSON: %v", err)
					continue
				}

				// Send the formatted JSON event
				responseChan <- jsonData
			case "message_delta":
				// Ignore this event for now
			case "message_stop":
				// End of the message
				close(responseChan)
				return nil
			default:
				log.Printf("Unexpected event type: %s", currentEvent)
			}
		} else {
			log.Printf("Unexpected line format: %s", string(line))
		}
	}

	return nil
}