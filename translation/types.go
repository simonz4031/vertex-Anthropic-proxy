package translation

type AnthropicRequest struct {
    Model     string    `json:"model"`
    Messages  []Message `json:"messages"`
    MaxTokens int       `json:"max_tokens"`
    Stream    bool      `json:"stream"`
}

type Message struct {
    Role    string      `json:"role"`
    Content interface{} `json:"content"`
}

type VertexAIRequest struct {
    AnthropicVersion string    `json:"anthropic_version"`
    Messages         []Message `json:"messages"`
    MaxTokens        int       `json:"max_tokens"`
}

type VertexAIResponse struct {
    ID         string    `json:"id"`
    Type       string    `json:"type"`
    Role       string    `json:"role"`
    Content    []Content `json:"content"`
    Model      string    `json:"model"`
    StopReason string    `json:"stop_reason"`
    Usage      Usage     `json:"usage"`
}

type Content struct {
    Type string `json:"type"`
    Text string `json:"text"`
}

type Usage struct {
    InputTokens  int `json:"input_tokens"`
    OutputTokens int `json:"output_tokens"`
}