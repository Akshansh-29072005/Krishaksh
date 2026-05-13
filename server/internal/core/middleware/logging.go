package middleware

import (
	"bytes"
	"encoding/json"
	"io"
	"strings"

	"github.com/aarcsx/krisho-backend/internal/observability"
	"github.com/gin-gonic/gin"
)

// BodyLogWriter wraps gin.ResponseWriter to capture response body
type BodyLogWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w BodyLogWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

func (w BodyLogWriter) WriteString(s string) (int, error) {
	w.body.WriteString(s)
	return w.ResponseWriter.WriteString(s)
}

// RequestResponseLogger logs important request/response payloads with sensitive data filtering
func RequestResponseLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Skip logging for health checks
		if strings.Contains(c.Request.URL.Path, "/health") || strings.Contains(c.Request.URL.Path, "/ready") {
			c.Next()
			return
		}

		// Capture request body
		var reqBody []byte
		if c.Request.Body != nil {
			reqBody, _ = io.ReadAll(c.Request.Body)
			// Restore body for downstream handlers
			c.Request.Body = io.NopCloser(bytes.NewReader(reqBody))
		}

		// Log request with sensitive data filtered
		reqBodyStr := string(reqBody)
		if len(reqBody) > 0 && isJSON(reqBody) {
			reqBodyStr = filterSensitiveData(string(reqBody))
		}

		// Capture response body
		blw := &BodyLogWriter{body: &bytes.Buffer{}, ResponseWriter: c.Writer}
		c.Writer = blw

		c.Next()

		// Log response
		respBodyStr := blw.body.String()
		if len(respBodyStr) > 0 && isJSON([]byte(respBodyStr)) {
			respBodyStr = filterSensitiveData(respBodyStr)
		}

		// Log important endpoints
		if isImportantEndpoint(c.Request.URL.Path, c.Request.Method) {
			observability.InitLogger().Info("request_response",
				"path", c.Request.URL.Path,
				"method", c.Request.Method,
				"status", c.Writer.Status(),
				"request_body", truncateString(reqBodyStr, 500),
				"response_body", truncateString(respBodyStr, 500),
				"request_id", c.GetString("request_id"),
			)
		}
	}
}

func isJSON(data []byte) bool {
	var js interface{}
	return json.Unmarshal(data, &js) == nil
}

func filterSensitiveData(jsonStr string) string {
	sensitiveKeys := []string{
		"password", "token", "secret", "api_key", "authorization",
		"razorpay_key", "razorpay_secret", "aws_secret", "gemini_api_key",
	}

	var obj map[string]interface{}
	if err := json.Unmarshal([]byte(jsonStr), &obj); err == nil {
		for _, key := range sensitiveKeys {
			if _, exists := obj[key]; exists {
				obj[key] = "***REDACTED***"
			}
		}

		if filtered, err := json.Marshal(obj); err == nil {
			return string(filtered)
		}
	}

	return jsonStr
}

func isImportantEndpoint(path, method string) bool {
	importantPaths := []string{
		"/api/v1/auth/",
		"/api/v1/orders/",
		"/api/v1/payments/",
		"/api/v1/users/",
		"/api/v1/admin/",
	}

	for _, p := range importantPaths {
		if strings.Contains(path, p) {
			return true
		}
	}
	return false
}

func truncateString(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen] + "...[truncated]"
}
