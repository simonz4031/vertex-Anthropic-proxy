package handlers

import (
    "encoding/json"
    "fmt"
    "io"
    "net/http"
    "vertexai-anthropic-proxy/config"
    "vertexai-anthropic-proxy/translation"
    "vertexai-anthropic-proxy/client"
    "vertexai-anthropic-proxy/utils"
)

type ContentItem struct {
    Type string `json:"type"`
    Text string `json:"text"`
}

type AnthropicResponse struct {
    ID            string        `json:"id"`
    Type          string        `json:"type"`
    Role          string        `json:"role"`
    Model         string        `json:"model"`
    Content       []ContentItem `json:"content"`
    StopReason    string        `json:"stop_reason"`
    StopSequence  interface{}   `json:"stop_sequence"`
    Usage         Usage         `json:"usage"`
}

type Usage struct {
    InputTokens  int `json:"input_tokens"`
    OutputTokens int `json:"output_tokens"`
}

func HandleMessages(cfg *config.Config) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        logger := utils.GetLogger()

        logger.Info("Received request to /v1/messages")

        // Read the request body
        body, err := io.ReadAll(r.Body)
        if err != nil {
            logger.Errorf("Error reading request body: %v", err)
            http.Error(w, "Error reading request", http.StatusBadRequest)
            return
        }
        defer r.Body.Close()

        logger.Infof("Request body: %s", string(body))

        // Parse the Anthropic request
        var anthropicReq translation.AnthropicRequest
        if err := json.Unmarshal(body, &anthropicReq); err != nil {
            logger.Errorf("Error parsing request: %v", err)
            http.Error(w, "Error parsing request", http.StatusBadRequest)
            return
        }

        logger.Info("Parsed Anthropic request successfully")

        // Translate Anthropic request to Vertex AI request
        vertexAIReq, err := translation.AnthropicToVertexAI(anthropicReq)
        if err != nil {
            logger.Errorf("Error translating Anthropic request to Vertex AI: %v", err)
            http.Error(w, "Error processing request", http.StatusInternalServerError)
            return
        }

        logger.Info("Translated request to Vertex AI format")

        // Send request to Vertex AI
        responseStream, err := client.SendToVertexAI(cfg, &vertexAIReq)
        if err != nil {
            logger.Errorf("Error sending request to Vertex AI: %v", err)
            http.Error(w, fmt.Sprintf("Error processing request: %v", err), http.StatusInternalServerError)
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
        var response map[string]interface{}
        err = json.Unmarshal(responseBody, &response)
        if err != nil {
            logger.Errorf("Error parsing response: %v", err)
            http.Error(w, "Error processing response", http.StatusInternalServerError)
            return
        }

        // Extract the content
        content := response["content"].([]interface{})[0].(map[string]interface{})["text"].(string)

        // Prepare the response in Anthropic's format
        anthropicResponse := AnthropicResponse{
            ID:         response["id"].(string),
            Type:       "message",
            Role:       "assistant",
            Model:      response["model"].(string),
            Content:    []ContentItem{{Type: "text", Text: content}},
            StopReason: response["stop_reason"].(string),
            Usage: Usage{
                InputTokens:  int(response["usage"].(map[string]interface{})["input_tokens"].(float64)),
                OutputTokens: int(response["usage"].(map[string]interface{})["output_tokens"].(float64)),
            },
        }

        // Add BOM at the beginning of the response
        w.Write([]byte{0xEF, 0xBB, 0xBF})

        if anthropicReq.Stream {
            // Streaming response
            w.Header().Set("Content-Type", "text/event-stream; charset=utf-8")
            w.Header().Set("Cache-Control", "no-cache")
            w.Header().Set("Connection", "keep-alive")

            // Send the response in chunks
            chunks := splitResponse(content, 10) // Split into 10 chunks, adjust as needed
            for i, chunk := range chunks {
                event := map[string]interface{}{
                    "type": "content_block_delta",
                    "delta": map[string]interface{}{
                        "type": "text_delta",
                        "text": chunk,
                    },
                    "index": i,
                }
                eventJSON, _ := json.Marshal(event)
                fmt.Fprintf(w, "data: %s\n\n", eventJSON)
                w.(http.Flusher).Flush()
            }

            // Send the final event
            finalEvent := map[string]interface{}{
                "type": "message_stop",
            }
            finalEventJSON, _ := json.Marshal(finalEvent)
            fmt.Fprintf(w, "data: %s\n\n", finalEventJSON)
            w.(http.Flusher).Flush()
        } else {
            // Non-streaming response
            w.Header().Set("Content-Type", "application/json; charset=utf-8")
            jsonResponse, err := json.Marshal(anthropicResponse)
            if err != nil {
                logger.Errorf("Error marshaling response: %v", err)
                http.Error(w, "Internal server error", http.StatusInternalServerError)
                return
            }
            w.Write(jsonResponse)
        }

        logger.Info("Finished sending response to client")
    }
}

func splitResponse(response string, chunks int) []string {
    var result []string
    length := len(response)
    chunkSize := length / chunks
    for i := 0; i < length; i += chunkSize {
        end := i + chunkSize
        if end > length {
            end = length
        }
        result = append(result, response[i:end])
    }
    return result
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

func HandleRefreshCredentials(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Since we've removed the RefreshCredentials function, we'll just return a success message.
	// In a production environment, you might want to implement a different refresh mechanism or remove this endpoint.
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Credentials refresh is not necessary with the current implementation"})
}

func HandleSetLogLevel(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var request struct {
		Level string `json:"level"`
	}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	utils.SetLogLevel(request.Level)
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Log level updated successfully"})
}