package routes

import (
	"BlockCertify/internal/config"
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// HealthRoutes registers /healthz and /readyz on the root router (no prefix).
// These must be reachable by Kubernetes probes without auth.
func HealthRoutes(r *gin.Engine) {
	// Liveness: pod is alive (not deadlocked/crashed). Never checks external deps.
	r.GET("/healthz", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	// Readiness: pod can serve traffic. Fails → k8s stops routing to this pod.
	r.GET("/readyz", func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(c.Request.Context(), 2*time.Second)
		defer cancel()

		checks := gin.H{}
		ready := true

		// PostgreSQL
		if sqlDB, err := config.DB.DB(); err != nil || sqlDB.PingContext(ctx) != nil {
			checks["postgres"] = "fail"
			ready = false
		} else {
			checks["postgres"] = "ok"
		}

		// Redis
		if err := config.RedisClient.Ping(ctx).Err(); err != nil {
			checks["redis"] = "fail"
			ready = false
		} else {
			checks["redis"] = "ok"
		}

		status := http.StatusOK
		if !ready {
			status = http.StatusServiceUnavailable
		}
		c.JSON(status, gin.H{"status": map[bool]string{true: "ready", false: "not ready"}[ready], "checks": checks})
	})
}
