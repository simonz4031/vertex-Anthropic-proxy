package translation

import (
	"log"
)

func AnthropicToVertexAI(ar AnthropicRequest) (VertexAIRequest, error) {
    log.Printf("Received Anthropic request: %+v", ar)

    normalizedModelName := NormalizeModelName(ar.Model)
    log.Printf("Normalized model name: %s", normalizedModelName)

    maxTokens := ar.MaxTokens
    if maxTokens == 0 {
        maxTokens = 1000 // Default value if not provided
    }

    vertexAIReq := VertexAIRequest{
        AnthropicVersion: "vertex-2023-10-16",
        Messages:         ar.Messages,
        System:           ar.System,
        MaxTokens:        maxTokens,
        Stream:           ar.Stream,
    }

    log.Printf("Translated to Vertex AI request: %+v", vertexAIReq)

    return vertexAIReq, nil
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