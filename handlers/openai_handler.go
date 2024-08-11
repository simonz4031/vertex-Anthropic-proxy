package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"vertexai-anthropic-proxy/client"
	"vertexai-anthropic-proxy/config"
	"vertexai-anthropic-proxy/translation"
	"vertexai-anthropic-proxy/utils"
)

func HandleOpenAIMessages(cfg *config.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logger := utils.GetLogger()

		logger.Info("Received request to OpenAI endpoint")

		// Read the request body
		body, err := io.ReadAll(r.Body)
		if err != nil {
			logger.Errorf("Error reading request body: %v", err)
			http.Error(w, "Error reading request", http.StatusBadRequest)
			return
		}
		defer r.Body.Close()

		logger.Infof("Request body: %s", string(body))

		// Parse the OpenAI request
		var openAIReq translation.OpenAIRequest
		if err := json.Unmarshal(body, &openAIReq); err != nil {
			logger.Errorf("Error parsing request: %v", err)
			http.Error(w, "Error parsing request", http.StatusBadRequest)
			return
		}

		logger.Info("Parsed OpenAI request successfully")

		// Translate OpenAI request to Anthropic request
		anthropicReq := translation.OpenAIToAnthropic(openAIReq)

		// Translate Anthropic request to Vertex AI request
		vertexAIReq, err := translation.AnthropicToVertexAI(anthropicReq)
		if err != nil {
			logger.Errorf("Error translating Anthropic request to Vertex AI: %v", err)
			http.Error(w, "Error processing request", http.StatusInternalServerError)
			return
		}

		logger.Info("Translated request to Vertex AI format")

		if openAIReq.Stream {
			// Set headers for SSE
			w.Header().Set("Content-Type", "text/event-stream")
			w.Header().Set("Cache-Control", "no-cache")
			w.Header().Set("Connection", "keep-alive")

			// Create a channel to receive streaming responses
			responseChan := make(chan []byte)

			// Start a goroutine to send the request to Vertex AI and write responses to the channel
			go func() {
				err := client.SendToVertexAIStream(cfg, &vertexAIReq, responseChan)
				if err != nil {
					logger.Errorf("Error sending request to Vertex AI: %v", err)
					close(responseChan)
				}
			}()

			// Stream responses back to the client
			for response := range responseChan {
				fmt.Fprintf(w, "data: %s\n\n", response)
				w.(http.Flusher).Flush()
			}

			// Send the final SSE message
			fmt.Fprintf(w, "data: [DONE]\n\n")
			w.(http.Flusher).Flush()
		} else {
			// Send request to Vertex AI
			responseStream, err := client.SendToVertexAI(cfg, &vertexAIReq)
			if err != nil {
				logger.Errorf("Error sending request to Vertex AI: %v", err)
				http.Error(w, "Error processing request", http.StatusInternalServerError)
				return
			}
			defer responseStream.Close()

			logger.Info("Received response from Vertex AI")

			// Read the entire response
			responseBody, err := io.ReadAll(responseStream)
			if err != nil {
				logger.Errorf("Error reading response from Vertex AI: %v", err)
				http.Error(w, "Error processing response", http.StatusInternalServerError)
				return
			}

			logger.Infof("Response body from Vertex AI: %s", string(responseBody))

			// Parse the response
			var vertexAIResp translation.VertexAIResponse
			err = json.Unmarshal(responseBody, &vertexAIResp)
			if err != nil {
				logger.Errorf("Error parsing response: %v", err)
				http.Error(w, "Error processing response", http.StatusInternalServerError)
				return
			}

			// Translate Vertex AI response to OpenAI response
			openAIResp := translation.VertexAIToOpenAI(vertexAIResp, openAIReq.Model)

			// Send the response
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(openAIResp)

			logger.Info("Finished sending response to client")
		}
	}
}