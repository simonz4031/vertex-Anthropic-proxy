package translation

import (
	"reflect"
	"testing"
)

func TestAnthropicToVertexAI(t *testing.T) {
	tests := []struct {
		name    string
		input   AnthropicRequest
		want    VertexAIRequest
		wantErr bool
	}{
		{
			name: "Basic conversion",
			input: AnthropicRequest{
				Model: "claude-v1",
				Messages: []Message{
					{Role: "user", Content: "Hello, how are you?"},
				},
				MaxTokens: 100,
			},
			want: VertexAIRequest{
				Instances: []Instance{
					{
						Messages: []Message{
							{Role: "user", Content: "Hello, how are you?"},
						},
					},
				},
				Parameters: Parameters{
					MaxOutputTokens: 100,
					Temperature:     0.7,
					TopP:            0.95,
					TopK:            40,
				},
			},
			wantErr: false,
		},
		// Add more test cases here
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := AnthropicToVertexAI(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("AnthropicToVertexAI() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("AnthropicToVertexAI() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestVertexAIToAnthropic(t *testing.T) {
	tests := []struct {
		name    string
		input   VertexAIResponse
		want    map[string]interface{}
		wantErr bool
	}{
		{
			name: "Basic conversion",
			input: VertexAIResponse{
				Predictions: []Prediction{
					{Content: "Hello! I'm doing well, thank you for asking. How can I assist you today?"},
				},
			},
			want: map[string]interface{}{
				"content": "Hello! I'm doing well, thank you for asking. How can I assist you today?",
			},
			wantErr: false,
		},
		{
			name: "Empty predictions",
			input: VertexAIResponse{
				Predictions: []Prediction{},
			},
			want:    nil,
			wantErr: true,
		},
		// Add more test cases here
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := VertexAIToAnthropic(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("VertexAIToAnthropic() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("VertexAIToAnthropic() = %v, want %v", got, tt.want)
			}
		})
	}
}