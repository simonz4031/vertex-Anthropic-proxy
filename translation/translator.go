package translation

import (
    "fmt"
    "log"
    "strings"
)

func AnthropicToVertexAI(ar AnthropicRequest) (VertexAIRequest, error) {
    log.Printf("Received Anthropic request: %+v", ar)

    normalizedModelName := NormalizeModelName(ar.Model)
    log.Printf("Normalized model name: %s", normalizedModelName)

    // Process messages to ensure content is always a string
    processedMessages := make([]Message, len(ar.Messages))
    for i, msg := range ar.Messages {
        processedMsg := Message{Role: msg.Role}
        switch content := msg.Content.(type) {
        case string:
            processedMsg.Content = content
        case []interface{}:
            // Join array elements into a single string
            var parts []string
            for _, part := range content {
                if contentMap, ok := part.(map[string]interface{}); ok {
                    if text, ok := contentMap["text"].(string); ok {
                        parts = append(parts, text)
                    }
                }
            }
            processedMsg.Content = strings.Join(parts, " ")
        default:
            return VertexAIRequest{}, fmt.Errorf("unsupported content type for message %d", i)
        }
        processedMessages[i] = processedMsg
    }

    vertexReq := VertexAIRequest{
        AnthropicVersion: "vertex-2023-10-16",
        Messages:         processedMessages,
        MaxTokens:        ar.MaxTokens,
    }

    log.Printf("Translated to Vertex AI request: %+v", vertexReq)

    return vertexReq, nil
}

func VertexAIToAnthropic(vr VertexAIResponse) (map[string]interface{}, error) {
    if len(vr.Content) == 0 {
        return nil, nil
    }

    content := ""
    for _, c := range vr.Content {
        if c.Type == "text" {
            content += c.Text
        }
    }

    anthropicResp := map[string]interface{}{
        "content": content,
        "model":   vr.Model,
        "usage":   vr.Usage,
    }

    return anthropicResp, nil
}