package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type HealthResponse struct {
	Status    string    `json:"status"`
	Timestamp time.Time `json:"timestamp"`
	Version   string    `json:"version"`
}

func HealthCheck() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusOK, HealthResponse{
			Status:    "ok",
			Timestamp: time.Now(),
			Version:   "1.0.0", // This should come from build info in production
		})
	}
}

func SetupHealthRoutes(router *gin.Engine) {
	router.GET("/health", HealthCheck())
	router.GET("/ready", HealthCheck()) // Readiness probe
	router.GET("/live", HealthCheck())  // Liveness probe
} 