package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"

	"github.com/ankuragarwal/fernfs/backend/handlers"
	"github.com/ankuragarwal/fernfs/backend/middleware"
	"github.com/ankuragarwal/fernfs/backend/storage/local"
)

func main() {
	// Initialize logger
	logger, err := zap.NewProduction()
	if err != nil {
		log.Fatalf("Failed to create logger: %v", err)
	}
	defer logger.Sync()

	// Set Gin mode
	gin.SetMode(gin.ReleaseMode)

	// Create router
	router := gin.New()

	// Add middleware
	router.Use(middleware.Logger(logger))
	router.Use(middleware.Metrics())
	router.Use(gin.Recovery())

	// Add metrics endpoint
	router.GET("/metrics", gin.WrapH(promhttp.Handler()))

	// Add health check endpoint
	healthHandler := handlers.NewHealthHandler()
	router.GET("/health", healthHandler.Check)

	// Initialize storage adapter
	dataDir := os.Getenv("DATA_DIR")
	if dataDir == "" {
		dataDir = "data"
	}

	storageAdapter, err := local.NewAdapter(dataDir)
	if err != nil {
		logger.Fatal("Failed to create storage adapter",
			zap.String("data_dir", dataDir),
			zap.Error(err),
		)
	}

	// Initialize and register file handler
	fileHandler := handlers.NewFileHandler(storageAdapter)
	fileHandler.RegisterRoutes(router)

	// Start server
	addr := os.Getenv("ADDR")
	if addr == "" {
		addr = ":8080"
	}

	logger.Info("Starting server",
		zap.String("address", addr),
		zap.String("mode", gin.Mode()),
	)

	if err := router.Run(addr); err != nil {
		logger.Fatal("Failed to start server",
			zap.Error(err),
		)
	}
} 