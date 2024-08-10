package middleware

import (
	"net/http"
	"vertexai-anthropic-proxy/config"
	"vertexai-anthropic-proxy/utils"
)

func AuthMiddleware(cfg *config.Config) func(http.HandlerFunc) http.HandlerFunc {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			logger := utils.GetLogger()
			apiKey := r.Header.Get("X-API-Key")

			logger.Infof("Received API Key: %s", apiKey)
			logger.Infof("Expected API Key: %s", cfg.AnthropicProxyAPIKey)
			logger.Infof("Keys match: %v", apiKey == cfg.AnthropicProxyAPIKey)

			if apiKey == "" || apiKey != cfg.AnthropicProxyAPIKey {
				logger.Warn("Unauthorized access attempt")
				utils.RespondWithError(w, http.StatusUnauthorized, "Unauthorized")
				return
			}

			next.ServeHTTP(w, r)
		}
	}
}
