package metrics

import (
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// PrometheusMiddleware coleta métricas HTTP
func PrometheusMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.FullPath()
		if path == "" {
			path = "unknown"
		}

		// Medir tamanho da requisição
		var requestSize float64
		if c.Request.ContentLength > 0 {
			requestSize = float64(c.Request.ContentLength)
		}

		c.Next()

		// Coletar métricas
		status := strconv.Itoa(c.Writer.Status())
		method := c.Request.Method
		duration := time.Since(start).Seconds()
		responseSize := float64(c.Writer.Size())

		HttpRequestsTotal.WithLabelValues(method, path, status).Inc()
		HttpRequestDuration.WithLabelValues(method, path).Observe(duration)
		HttpRequestSize.WithLabelValues(method, path).Observe(requestSize)
		HttpResponseSize.WithLabelValues(method, path).Observe(responseSize)
	}
}

// MetricsHandler retorna o handler para o endpoint /metrics
func MetricsHandler() gin.HandlerFunc {
	return gin.WrapH(promhttp.Handler())
}

// Helper functions para uso nos handlers

// RecordDBQuery registra uma query no banco
func RecordDBQuery(operation, table string, duration time.Duration, err error) {
	status := "success"
	if err != nil {
		status = "error"
	}
	DbQueryTotal.WithLabelValues(operation, table, status).Inc()
	DbQueryDuration.WithLabelValues(operation, table).Observe(duration.Seconds())
}

// RecordCacheHit registra um acerto de cache
func RecordCacheHit(cacheType string) {
	CacheHits.WithLabelValues(cacheType).Inc()
}

// RecordCacheMiss registra um miss de cache
func RecordCacheMiss(cacheType string) {
	CacheMisses.WithLabelValues(cacheType).Inc()
}
