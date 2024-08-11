package middleware

import (
	"net/http"
	"strings"
	"vertexai-anthropic-proxy/config"
	"vertexai-anthropic-proxy/utils"
)

func AuthMiddleware(cfg *config.Config) func(http.HandlerFunc) http.HandlerFunc {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			logger := utils.GetLogger()
			apiKey := r.Header.Get("Authorization")

			// Remove "Bearer " prefix if present
			apiKey = strings.TrimPrefix(apiKey, "Bearer ")

			if apiKey == "" {
				apiKey = r.Header.Get("X-API-Key")
			}

			logger.Infof("Received API Key: %s", apiKey)

			if apiKey == "" || (apiKey != cfg.AnthropicProxyAPIKey && apiKey != cfg.OpenAIProxyAPIKey) {
				logger.Warn("Unauthorized access attempt")
				utils.RespondWithError(w, http.StatusUnauthorized, "Unauthorized")
				return
			}

			next.ServeHTTP(w, r)
		}
	}
}