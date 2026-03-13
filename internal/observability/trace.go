package observability

import (
	"context"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
)

type traceCtxKey struct{}

func TraceIDFromContext(ctx context.Context) string {
	if v := ctx.Value(traceCtxKey{}); v != nil {
		if s, ok := v.(string); ok {
			return s
		}
	}
	return ""
}

func WithTraceID(ctx context.Context, id string) context.Context {
	return context.WithValue(ctx, traceCtxKey{}, id)
}

func TracingMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		traceID := c.GetHeader("X-Request-ID")
		if traceID == "" {
			traceID = uuid.New().String()
		}
		c.Writer.Header().Set("X-Request-ID", traceID)
		c.Request = c.Request.WithContext(WithTraceID(c.Request.Context(), traceID))

		start := time.Now()
		method := c.Request.Method
		path := c.FullPath()
		if path == "" {
			path = c.Request.URL.Path
		}

		log.Info().
			Str("trace_id", traceID).
			Str("method", method).
			Str("path", path).
			Str("client_ip", c.ClientIP()).
			Msg("request_start")

		c.Next()

		log.Info().
			Str("trace_id", traceID).
			Str("method", method).
			Str("path", path).
			Int("status", c.Writer.Status()).
			Int64("duration_ms", time.Since(start).Milliseconds()).
			Msg("request_complete")
	}
}
