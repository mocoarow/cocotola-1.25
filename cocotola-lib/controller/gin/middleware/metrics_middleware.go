package middleware

import (
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	// Total HTTP requests labeled by status, method, and path
	httpRequestsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{ //nolint:exhaustruct
			Name: "http_requests_total",
			Help: "Total number of HTTP requests",
		},
		[]string{"method", "path", "status"},
	)

	// HTTP request duration labeled by method and path
	httpRequestDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{ //nolint:exhaustruct
			Name:    "http_request_duration_seconds",
			Help:    "HTTP request duration in seconds",
			Buckets: prometheus.DefBuckets, // Default buckets: 0.005, 0.01, 0.025, 0.05, 0.1, 0.25, 0.5, 1, 2.5, 5, 10
		},
		[]string{"method", "path"},
	)
)

// PrometheusMiddleware returns a Gin middleware that collects HTTP metrics for Prometheus.
func PrometheusMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Record request start time
		start := time.Now()

		// Run the next handler in the chain
		c.Next()

		// Calculate how long the handler took
		duration := time.Since(start).Seconds()

		// Persist metrics with normalized labels
		path := c.FullPath()
		if path == "" {
			path = c.Request.URL.Path
		}
		method := c.Request.Method
		status := strconv.Itoa(c.Writer.Status())

		// Count the request
		httpRequestsTotal.WithLabelValues(method, path, status).Inc()

		// Record the request duration
		httpRequestDuration.WithLabelValues(method, path).Observe(duration)
	}
}
