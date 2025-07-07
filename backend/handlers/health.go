package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type HealthHandler struct{}

func NewHealthHandler() *HealthHandler {
	return &HealthHandler{}
}

func (h *HealthHandler) Check(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status": "ok",
	})
}

// HealthCheck handles basic health check requests
func HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status": "ok",
	})
}

// ReadinessCheck handles readiness check requests
func ReadinessCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status": "ready",
	})
}

func SetupHealthRoutes(router *gin.Engine) {
	router.GET("/health", HealthCheck)
	router.GET("/ready", ReadinessCheck) // Readiness probe
	router.GET("/live", HealthCheck)  // Liveness probe
} 