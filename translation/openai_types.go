package translation

import (
	"time"
)

type OpenAIRequest struct {
	Model     string    `json:"model"`
	Messages  []Message `json:"messages"`
	MaxTokens int       `json:"max_tokens"`
	Stream    bool      `json:"stream"`
}

type OpenAIResponse struct {
	ID      string   `json:"id"`
	Object  string   `json:"object"`
	Created int64    `json:"created"`
	Model   string   `json:"model"`
	Usage   Usage    `json:"usage"`
	Choices []Choice `json:"choices"`
}

type Choice struct {
	Message      Message `json:"message"`
	FinishReason string  `json:"finish_reason"`
	Index        int     `json:"index"`
}

func OpenAIToAnthropic(openAIReq OpenAIRequest) AnthropicRequest {
	maxTokens := openAIReq.MaxTokens
	if maxTokens == 0 {
		maxTokens = 1000 // Default value if not provided
	}

	var systemMessage string
	var userMessages []Message

	for _, msg := range openAIReq.Messages {
		if msg.Role == "system" {
			systemMessage = msg.Content.(string)
		} else {
			userMessages = append(userMessages, msg)
		}
	}

	return AnthropicRequest{
		Model:     "claude-3-5-sonnet@20240620",
		Messages:  userMessages,
		System:    systemMessage,
		MaxTokens: maxTokens,
		Stream:    openAIReq.Stream,
	}
}

func VertexAIToOpenAI(vertexAIResp VertexAIResponse, model string) OpenAIResponse {
	return OpenAIResponse{
		ID:      vertexAIResp.ID,
		Object:  "chat.completion",
		Created: time.Now().Unix(),
		Model:   model,
		Usage:   vertexAIResp.Usage,
		Choices: []Choice{
			{
				Message: Message{
					Role:    "assistant",
					Content: vertexAIResp.Content[0].Text,
				},
				FinishReason: vertexAIResp.StopReason,
				Index:        0,
			},
		},
	}
}