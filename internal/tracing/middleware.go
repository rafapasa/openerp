package tracing

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
)

// TracingMiddleware é o middleware HTTP para tracing
func TracingMiddleware(serviceName string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Extrair contexto de propagação
		ctx := otel.GetTextMapPropagator().Extract(
			c.Request.Context(),
			propagation.HeaderCarrier(c.Request.Header),
		)

		// Iniciar novo span
		spanName := fmt.Sprintf("%s %s", c.Request.Method, c.Request.URL.Path)
		ctx, span := Tracer.Start(ctx, spanName,
			trace.WithAttributes(
				attribute.String("http.method", c.Request.Method),
				attribute.String("http.url", c.Request.URL.String()),
				attribute.String("http.user_agent", c.Request.UserAgent()),
				attribute.String("http.client_ip", c.ClientIP()),
			),
		)
		defer span.End()

		// Adicionar ao contexto para uso nos handlers
		c.Set("tracing_context", ctx)
		c.Set("span", span)

		c.Next()

		// Registrar status da resposta
		span.SetAttributes(
			attribute.Int("http.status_code", c.Writer.Status()),
			attribute.Int("http.response_size", c.Writer.Size()),
		)
	}
}

// Helper functions para uso nos handlers

// GetSpan recupera o span do contexto
func GetSpan(c *gin.Context) trace.Span {
	if span, exists := c.Get("span"); exists {
		return span.(trace.Span)
	}
	return nil
}

// AddErrorToSpan adiciona um erro ao span
func AddErrorToSpan(c *gin.Context, err error) {
	if span := GetSpan(c); span != nil {
		span.RecordError(err)
		span.SetAttributes(attribute.Bool("error", true))
	}
}

// AddEvent adiciona um evento ao span
func AddEvent(c *gin.Context, name string, attrs ...attribute.KeyValue) {
	if span := GetSpan(c); span != nil {
		span.AddEvent(name, trace.WithAttributes(attrs...))
	}
}
