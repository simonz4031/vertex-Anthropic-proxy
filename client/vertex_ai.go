package client

import (
    "bytes"
    "encoding/json"
    "fmt"
    "io/ioutil"
    "log"
    "net/http"
    "os/exec"

    "vertexai-anthropic-proxy/config"
    "vertexai-anthropic-proxy/translation"
)

func SendToVertexAI(cfg *config.Config, req *translation.VertexAIRequest) ([]byte, error) {
    url := cfg.VertexAIEndpoint
    if url == "" {
        return nil, fmt.Errorf("Vertex AI endpoint URL is empty")
    }

    jsonData, err := json.Marshal(req)
    if err != nil {
        return nil, fmt.Errorf("error marshaling request: %v", err)
    }

    log.Printf("Sending request to Vertex AI: %s", url)
    log.Printf("Request body: %s", string(jsonData))

    httpReq, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
    if err != nil {
        return nil, fmt.Errorf("error creating request: %v", err)
    }

    httpReq.Header.Set("Content-Type", "application/json")
    
    // Get access token using gcloud
    accessToken, err := getAccessToken()
    if err != nil {
        return nil, fmt.Errorf("error getting access token: %v", err)
    }
    httpReq.Header.Set("Authorization", "Bearer "+accessToken)

    client := &http.Client{}
    resp, err := client.Do(httpReq)
    if err != nil {
        return nil, fmt.Errorf("error sending request: %v", err)
    }
    defer resp.Body.Close()

    body, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        return nil, fmt.Errorf("error reading response body: %v", err)
    }

    if resp.StatusCode != http.StatusOK {
        return nil, fmt.Errorf("Vertex AI returned non-OK status: %d, body: %s", resp.StatusCode, string(body))
    }

    return body, nil
}

func getAccessToken() (string, error) {
    cmd := exec.Command("gcloud", "auth", "print-access-token")
    output, err := cmd.Output()
    if err != nil {
        return "", err
    }
    return string(bytes.TrimSpace(output)), nil
}