# Vertex AI Anthropic Proxy

This project provides a proxy server that allows you to use the Anthropic Claude API format while leveraging Google Cloud's Vertex AI backend. It translates requests from the Anthropic format to Vertex AI format and vice versa.

## Table of Contents

- [Setup](#setup)
- [Configuration](#configuration)
- [API Endpoints](#api-endpoints)
- [Usage Examples](#usage-examples)
- [Development](#development)

## Setup

1. Clone the repository:
   ```
   git clone https://github.com/yourusername/vertex-anthropic-proxy.git
   cd vertex-anthropic-proxy
   ```

2. Install dependencies:
   ```
   go mod tidy
   ```

3. Set up your environment variables by creating a `.env` file in the project root:
   ```
   PORT=8070
   VERTEX_AI_PROJECT_ID=your-project-id
   VERTEX_AI_REGION=your-region
   MODEL=claude-3-5-sonnet@20240620
   ANTHROPIC_API_KEY=your-api-key
   ```

4. Build the project:
   ```
   go build
   ```

5. Run the server:
   ```
   ./vertexai-anthropic-proxy
   ```

## Configuration

The following environment variables are required:

- `PORT`: The port on which the server will listen (default: 8070)
- `VERTEX_AI_PROJECT_ID`: Your Google Cloud project ID
- `VERTEX_AI_REGION`: The region for Vertex AI (e.g., us-east5)
- `MODEL`: The Claude model to use (e.g., claude-3-5-sonnet@20240620)
- `ANTHROPIC_API_KEY`: Your Anthropic API key

## API Endpoints

### POST /v1/messages

This endpoint accepts requests in the Anthropic Claude API format and returns responses in the same format.

Request body:
```json
{
  "model": "claude-3-5-sonnet",
  "messages": [
    {"role": "user", "content": "Your message here"}
  ],
  "max_tokens": 100
}
```

Response body:
```json
{
  "content": "Claude's response",
  "model": "claude-3-5-sonnet-20240620",
  "usage": {
    "input_tokens": 10,
    "output_tokens": 20
  }
}
```

## Usage Examples

### cURL

```bash
curl -X POST http://localhost:8070/v1/messages \
  -H "Content-Type: application/json" \
  -d '{
    "model": "claude-3-5-sonnet",
    "messages": [
      {"role": "user", "content": "Hello, Claude! How are you today?"}
    ],
    "max_tokens": 100
  }'
```

### Python

```python
import requests
import json

url = "http://localhost:8070/v1/messages"
headers = {"Content-Type": "application/json"}
data = {
    "model": "claude-3-5-sonnet",
    "messages": [
        {"role": "user", "content": "Hello, Claude! How are you today?"}
    ],
    "max_tokens": 100
}

response = requests.post(url, headers=headers, data=json.dumps(data))
print(response.json())
```

## Development

To contribute to this project:

1. Fork the repository
2. Create a new branch for your feature
3. Implement your changes
4. Write or update tests as necessary
5. Submit a pull request

Please ensure your code adheres to the existing style and passes all tests before submitting a pull request.