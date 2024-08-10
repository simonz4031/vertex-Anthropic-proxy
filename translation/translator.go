package translation

import (
    "log"
)

func AnthropicToVertexAI(ar AnthropicRequest) (VertexAIRequest, error) {
    // Log the incoming Anthropic request
    log.Printf("Received Anthropic request: %+v", ar)

    vertexReq := VertexAIRequest{
        AnthropicVersion: "vertex-2023-10-16",
        Messages:         ar.Messages,
        MaxTokens:        ar.MaxTokens,
    }

    // Log the outgoing Vertex AI request
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