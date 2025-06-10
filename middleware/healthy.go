package middleware

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// HealthStatus represents the health status of a component
type HealthStatus string

const (
	HealthStatusHealthy   HealthStatus = "healthy"
	HealthStatusUnhealthy HealthStatus = "unhealthy"
	HealthStatusDegraded  HealthStatus = "degraded"
)

// HealthCheck represents a single health check
type HealthCheck struct {
	Name        string                 `json:"name"`
	Status      HealthStatus           `json:"status"`
	Message     string                 `json:"message,omitempty"`
	LastChecked time.Time              `json:"last_checked"`
	Duration    time.Duration          `json:"duration"`
	Metadata    map[string]interface{} `json:"metadata,omitempty"`
}

// HealthResponse represents the overall health response
type HealthResponse struct {
	Status    HealthStatus           `json:"status"`
	Timestamp time.Time              `json:"timestamp"`
	Service   string                 `json:"service"`
	Version   string                 `json:"version"`
	Uptime    time.Duration          `json:"uptime"`
	Checks    map[string]HealthCheck `json:"checks,omitempty"`
	Metadata  map[string]interface{} `json:"metadata,omitempty"`
}

// HealthChecker is a function that performs a health check
type HealthChecker func() HealthCheck

// HealthConfig holds configuration for health check middleware
type HealthConfig struct {
	ServiceName    string
	ServiceVersion string
	StartTime      time.Time
	Checkers       map[string]HealthChecker
	HealthPath     string
	ReadinessPath  string
	LivenessPath   string
}

// DefaultHealthConfig returns default health configuration
func DefaultHealthConfig(serviceName, version string) HealthConfig {
	return HealthConfig{
		ServiceName:    serviceName,
		ServiceVersion: version,
		StartTime:      time.Now(),
		Checkers:       make(map[string]HealthChecker),
		HealthPath:     "/health",
		ReadinessPath:  "/ready",
		LivenessPath:   "/live",
	}
}

// HealthMiddleware creates health check endpoints
func HealthMiddleware(config HealthConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		switch c.Request.URL.Path {
		case config.HealthPath:
			handleHealthCheck(c, config, true) // Detailed health check
		case config.ReadinessPath:
			handleReadinessCheck(c, config) // Readiness check
		case config.LivenessPath:
			handleLivenessCheck(c, config) // Liveness check
		default:
			c.Next()
		}
	}
}

// SimpleHealthMiddleware creates a simple health endpoint
func SimpleHealthMiddleware(serviceName, version string) gin.HandlerFunc {
	config := DefaultHealthConfig(serviceName, version)

	return func(c *gin.Context) {
		if c.Request.URL.Path == "/health" {
			handleHealthCheck(c, config, false) // Simple health check
		} else {
			c.Next()
		}
	}
}

// handleHealthCheck handles the main health check endpoint
func handleHealthCheck(c *gin.Context, config HealthConfig, detailed bool) {
	response := HealthResponse{
		Timestamp: time.Now(),
		Service:   config.ServiceName,
		Version:   config.ServiceVersion,
		Uptime:    time.Since(config.StartTime),
		Status:    HealthStatusHealthy,
	}

	if detailed && len(config.Checkers) > 0 {
		response.Checks = runHealthChecks(config.Checkers)
		response.Status = determineOverallStatus(response.Checks)
	}

	statusCode := http.StatusOK
	if response.Status == HealthStatusUnhealthy {
		statusCode = http.StatusServiceUnavailable
	}

	c.JSON(statusCode, response)
}

// handleReadinessCheck handles readiness probe (Kubernetes)
func handleReadinessCheck(c *gin.Context, config HealthConfig) {
	checks := runHealthChecks(config.Checkers)
	status := determineOverallStatus(checks)

	response := gin.H{
		"status":    status,
		"timestamp": time.Now(),
		"service":   config.ServiceName,
	}

	statusCode := http.StatusOK
	if status == HealthStatusUnhealthy {
		statusCode = http.StatusServiceUnavailable
	}

	c.JSON(statusCode, response)
}

// handleLivenessCheck handles liveness probe (Kubernetes)
func handleLivenessCheck(c *gin.Context, config HealthConfig) {
	// Liveness is usually just "is the service running?"
	c.JSON(http.StatusOK, gin.H{
		"status":    HealthStatusHealthy,
		"timestamp": time.Now(),
		"service":   config.ServiceName,
		"uptime":    time.Since(config.StartTime),
	})
}

// runHealthChecks executes all registered health checkers
func runHealthChecks(checkers map[string]HealthChecker) map[string]HealthCheck {
	checks := make(map[string]HealthCheck)

	for name, checker := range checkers {
		start := time.Now()
		check := checker()
		check.Duration = time.Since(start)
		check.LastChecked = time.Now()
		checks[name] = check
	}

	return checks
}

// determineOverallStatus determines the overall status based on individual checks
func determineOverallStatus(checks map[string]HealthCheck) HealthStatus {
	if len(checks) == 0 {
		return HealthStatusHealthy
	}

	hasUnhealthy := false
	hasDegraded := false

	for _, check := range checks {
		switch check.Status {
		case HealthStatusUnhealthy:
			hasUnhealthy = true
		case HealthStatusDegraded:
			hasDegraded = true
		}
	}

	if hasUnhealthy {
		return HealthStatusUnhealthy
	}
	if hasDegraded {
		return HealthStatusDegraded
	}

	return HealthStatusHealthy
}

// AddHealthChecker adds a health checker to the configuration
func (c *HealthConfig) AddHealthChecker(name string, checker HealthChecker) {
	if c.Checkers == nil {
		c.Checkers = make(map[string]HealthChecker)
	}
	c.Checkers[name] = checker
}

// DatabaseHealthChecker creates a health checker for database connectivity
func DatabaseHealthChecker(pingFunc func() error) HealthChecker {
	return func() HealthCheck {
		start := time.Now()
		err := pingFunc()
		duration := time.Since(start)

		if err != nil {
			return HealthCheck{
				Name:     "database",
				Status:   HealthStatusUnhealthy,
				Message:  err.Error(),
				Duration: duration,
			}
		}

		status := HealthStatusHealthy
		if duration > 100*time.Millisecond {
			status = HealthStatusDegraded
		}

		return HealthCheck{
			Name:     "database",
			Status:   status,
			Duration: duration,
			Metadata: map[string]interface{}{
				"response_time_ms": duration.Milliseconds(),
			},
		}
	}
}

// ExternalServiceHealthChecker creates a health checker for external services
func ExternalServiceHealthChecker(name, url string, timeout time.Duration) HealthChecker {
	return func() HealthCheck {
		start := time.Now()

		// This is a placeholder - in real implementation you'd make HTTP request
		// For now, simulate a check
		time.Sleep(10 * time.Millisecond) // Simulate network call
		duration := time.Since(start)

		return HealthCheck{
			Name:     name,
			Status:   HealthStatusHealthy,
			Duration: duration,
			Metadata: map[string]interface{}{
				"url":              url,
				"timeout_ms":       timeout.Milliseconds(),
				"response_time_ms": duration.Milliseconds(),
			},
		}
	}
}

// MemoryHealthChecker creates a health checker for memory usage
func MemoryHealthChecker(maxMemoryMB int64) HealthChecker {
	return func() HealthCheck {
		// This is a placeholder - in real implementation you'd check actual memory usage
		currentMemoryMB := int64(50) // Simulate current memory usage

		status := HealthStatusHealthy
		message := "Memory usage within normal range"

		if currentMemoryMB > maxMemoryMB {
			status = HealthStatusUnhealthy
			message = "Memory usage exceeded threshold"
		} else if currentMemoryMB > maxMemoryMB*80/100 {
			status = HealthStatusDegraded
			message = "Memory usage approaching threshold"
		}

		return HealthCheck{
			Name:    "memory",
			Status:  status,
			Message: message,
			Metadata: map[string]interface{}{
				"current_memory_mb": currentMemoryMB,
				"max_memory_mb":     maxMemoryMB,
				"usage_percent":     (currentMemoryMB * 100) / maxMemoryMB,
			},
		}
	}
}

// DiskHealthChecker creates a health checker for disk usage
func DiskHealthChecker(path string, maxUsagePercent int) HealthChecker {
	return func() HealthCheck {
		// This is a placeholder - in real implementation you'd check actual disk usage
		currentUsagePercent := 45 // Simulate current disk usage

		status := HealthStatusHealthy
		message := "Disk usage within normal range"

		if currentUsagePercent > maxUsagePercent {
			status = HealthStatusUnhealthy
			message = "Disk usage exceeded threshold"
		} else if currentUsagePercent > maxUsagePercent*80/100 {
			status = HealthStatusDegraded
			message = "Disk usage approaching threshold"
		}

		return HealthCheck{
			Name:    "disk",
			Status:  status,
			Message: message,
			Metadata: map[string]interface{}{
				"path":              path,
				"usage_percent":     currentUsagePercent,
				"max_usage_percent": maxUsagePercent,
			},
		}
	}
}
