package middleware

import (
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/aarcsx/krisho-backend/internal/observability"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func RequestContext() gin.HandlerFunc {
	return func(c *gin.Context) {
		rid := c.GetHeader("X-Request-ID")
		if rid == "" {
			rid = uuid.NewString()
		}
		trace := c.GetHeader("X-Trace-ID")
		if trace == "" {
			trace = uuid.NewString()
		}
		c.Writer.Header().Set("X-Request-ID", rid)
		c.Writer.Header().Set("X-Trace-ID", trace)
		c.Set("request_id", rid)
		c.Set("trace_id", trace)
		ctx := observability.WithTraceID(observability.WithRequestID(c.Request.Context(), rid), trace)
		c.Request = c.Request.WithContext(ctx)
		c.Next()
	}
}

func RequestTiming() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()
		d := time.Since(start)
		status := c.Writer.Status()
		observability.M.Inc(fmt.Sprintf("http_requests_total:%s:%d", c.FullPath(), status))
		observability.M.Add("http_request_duration_ms_total", d.Milliseconds())
		observability.InitLogger().Info("http_request",
			"path", c.FullPath(), "method", c.Request.Method, "status", status,
			"duration_ms", d.Milliseconds(), "request_id", c.GetString("request_id"), "trace_id", c.GetString("trace_id"))
	}
}

func ErrorTracker() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
		if len(c.Errors) > 0 {
			for _, e := range c.Errors {
				observability.M.Inc("errors_total")
				observability.InitLogger().Error("request_error", "error", e.Error(), "path", c.FullPath(), "request_id", c.GetString("request_id"), "trace_id", c.GetString("trace_id"))
			}
		}
	}
}

func UploadSizeLimit(max int64) gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.ContentLength > max {
			c.AbortWithStatusJSON(http.StatusRequestEntityTooLarge, gin.H{"success": false, "message": "payload too large"})
			return
		}
		c.Next()
	}
}

func ValidationGuard() gin.HandlerFunc {
	return func(c *gin.Context) {
		if strings.Contains(strings.ToLower(c.Request.URL.Path), "..") {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"success": false, "message": "invalid path"})
			return
		}
		c.Next()
	}
}

func RecoveryJSON() gin.HandlerFunc {
	return gin.CustomRecovery(func(c *gin.Context, rec interface{}) {
		observability.M.Inc("panic_total")
		observability.InitLogger().Error("panic", "value", rec, "path", c.FullPath(), "request_id", c.GetString("request_id"))
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"success": false, "message": "internal server error"})
	})
}

var (
	rlMu      sync.Mutex
	rlBuckets = map[string]*bucket{}
)

type bucket struct {
	tokens int
	last   time.Time
}

func RateLimit(perMinute int) gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := c.ClientIP()
		now := time.Now()
		rlMu.Lock()
		b := rlBuckets[ip]
		if b == nil {
			b = &bucket{tokens: perMinute - 1, last: now}
			rlBuckets[ip] = b
			rlMu.Unlock()
			c.Next()
			return
		}
		elapsed := now.Sub(b.last)
		if elapsed >= time.Minute {
			b.tokens = perMinute
			b.last = now
		}
		if b.tokens <= 0 {
			rlMu.Unlock()
			observability.M.Inc("rate_limited_total")
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{"success": false, "message": "rate limited"})
			return
		}
		b.tokens--
		rlMu.Unlock()
		c.Next()
	}
}

func AbusePrevention() gin.HandlerFunc {
	return func(c *gin.Context) {
		ua := strings.ToLower(c.GetHeader("User-Agent"))
		if strings.Contains(ua, "sqlmap") || strings.Contains(ua, "nikto") {
			observability.M.Inc("abuse_blocked_total")
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"success": false, "message": "request blocked"})
			return
		}
		c.Next()
	}
}
