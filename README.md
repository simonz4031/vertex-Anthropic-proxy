# Vertex AI Anthropic Proxy

This project provides a proxy server that translates requests from Anthropic's Claude API and OpenAI's API format to Google Cloud's Vertex AI, specifically for the Claude 3.5 Sonnet model.

## Table of Contents

- [Setup](#setup)
- [Configuration](#configuration)
- [API Endpoints](#api-endpoints)
- [Usage Examples](#usage-examples)
- [Development](#development)
- [Debugging](#debugging)
- [Refreshing Google Credentials](#refreshing-google-credentials)
- [Docker Deployment](#docker-deployment)
- [Security Considerations](#security-considerations)

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
   OPENAI_PROXY_API_KEY=your-openai-proxy-api-key
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
- `OPENAI_PROXY_API_KEY`: Your OpenAI proxy API key

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

### POST /v1/chat/completions

This endpoint accepts requests in both Anthropic Claude API and OpenAI API formats, and returns responses in the corresponding format.

OpenAI API Request body:
```json
{
  "model": "gpt-3.5-turbo",
  "messages": [
    {"role": "user", "content": "Your message here"}
  ]
}
```

OpenAI API Response body:
```json
{
  "id": "chatcmpl-123",
  "object": "chat.completion",
  "created": 1677652288,
  "model": "gpt-3.5-turbo-0613",
  "choices": [{
    "index": 0,
    "message": {
      "role": "assistant",
      "content": "Claude's response"
    },
    "finish_reason": "stop"
  }]
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

## Debugging

This project uses dynamic log levels, allowing you to change the logging verbosity at runtime. By default, the log level is set to "info". To change the log level:

1. While the server is running, send a POST request to the `/set-log-level` endpoint:

   ```
   curl -X POST http://localhost:8070/set-log-level -H "Content-Type: application/json" -d '{"level":"debug"}'
   ```

   This will set the log level to "debug", enabling more verbose logging.

2. To revert to the default "info" level:

   ```
   curl -X POST http://localhost:8070/set-log-level -H "Content-Type: application/json" -d '{"level":"info"}'
   ```

Available log levels are: "debug", "info", "warn", "error", "dpanic", "panic", and "fatal".

Note: Changing the log level affects all subsequent log messages. Use debug logging judiciously in production environments as it may impact performance.

## Refreshing Google Credentials

If you need to refresh the Google credentials without restarting the service, you can use the `/refresh-credentials` endpoint:

```
curl -X POST http://localhost:8070/refresh-credentials
```

This will attempt to refresh the Google credentials used by the service. If successful, it will return a 200 OK status with a success message. If there's an error, it will return an appropriate error status and message.

Note: This endpoint should be secured in production environments to prevent unauthorized access.

## Setting up Google Cloud Credentials

When running this service, especially in a containerized environment, you need to provide Google Cloud credentials. There are two main ways to do this:

1. **Using a Service Account Key File:**
   - Create a service account in your Google Cloud Console.
   - Download the JSON key file for this service account.
   - Set the `GOOGLE_APPLICATION_CREDENTIALS` environment variable to the path of this JSON file.

   Example:
   ```
   export GOOGLE_APPLICATION_CREDENTIALS="/path/to/your/service-account-key.json"
   ```

   When using Docker, you can mount this file into the container and set the environment variable in your Dockerfile or docker-compose file.

2. **Using Google Cloud Compute Engine Default Credentials:**
   If you're running this service on a Google Cloud Compute Engine instance, you can use the instance's default service account. Make sure the instance has the necessary permissions to access Vertex AI.

## Docker Deployment

This service can be easily deployed using Docker. Here are the steps to build and run the Docker container:

1. **Build the Docker image:**

   Navigate to the project directory and run:

   ```
   docker build -t vertexai-anthropic-proxy .
   ```

2. **Run the Docker container:**

   ```
   docker run -p 8070:8070 -v /path/to/your/service-account-key.json:/etc/secrets/service-account-key.json vertexai-anthropic-proxy
   ```

   Replace `/path/to/your/service-account-key.json` with the actual path to your Google Cloud service account key file.

3. **Environment Variables:**

   You can pass environment variables to the container using the `-e` flag. For example:

   ```
   docker run -p 8070:8070 \
     -v /path/to/your/service-account-key.json:/etc/secrets/service-account-key.json \
     -e VERTEX_AI_PROJECT_ID=your-project-id \
     -e VERTEX_AI_REGION=your-region \
     -e VERTEX_AI_ENDPOINT=your-endpoint \
     -e MODEL=your-model \
     -e ANTHROPIC_PROXY_API_KEY=your-api-key \
     -e OPENAI_PROXY_API_KEY=your-openai-proxy-api-key \
     vertexai-anthropic-proxy
   ```

4. **Accessing the Service:**

   Once the container is running, you can access the service at `http://localhost:8070`.

5. **Refreshing Credentials:**

   To refresh the Google Cloud credentials, you can use the `/refresh-credentials` endpoint:

   ```
   curl -X POST http://localhost:8070/refresh-credentials
   ```

   This will attempt to reload the credentials from the mounted service account key file.

## Security Considerations

- Never commit your service account key to version control.
- In production environments, consider using more secure methods to provide credentials, such as using Google Cloud's built-in service account when running on Google Cloud Platform, or using a secrets management system.
- Ensure that the `/refresh-credentials` endpoint is properly secured to prevent unauthorized access.