package api

import (
	"net/http"
	"time"

	"network-discovery/internal/discovery"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func SetupRoutes(discovery *discovery.NetworkDiscovery) *gin.Engine {
	// Create Gin router
	gin.SetMode(gin.ReleaseMode)
	router := gin.New()

	// Add middleware
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	// CORS middleware
	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	config.AllowMethods = []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"}
	config.AllowHeaders = []string{"Origin", "Content-Length", "Content-Type", "Authorization"}
	router.Use(cors.New(config))

	// Create handlers
	handlers := NewHandlers(discovery)

	// API versioning
	v1 := router.Group("/api/v1")
	{
		// Health and status endpoints
		v1.GET("/health", handlers.GetHealth)
		v1.GET("/version", handlers.GetVersion)

		// Network discovery endpoints
		network := v1.Group("/network")
		{
			network.POST("/scan", handlers.ScanNetwork)
			network.GET("/quick-scan", handlers.QuickScan)
			network.GET("/validate", handlers.ValidateNetwork)
		}

		// Device endpoints
		device := v1.Group("/device")
		{
			device.GET("/:ip", handlers.ScanDevice)
		}
	}

	// Serve static files (if needed for frontend)
	router.Static("/static", "./static")
	router.LoadHTMLGlob("templates/*")

	router.GET("/index", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", gin.H{
			"title": "Network Discovery UI",
		})
	})

	// Default route for API documentation
	router.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Network Discovery API",
			"version": "1.0.0",
			"endpoints": gin.H{
				"health":       "GET /api/v1/health",
				"version":      "GET /api/v1/version",
				"scan_network": "POST /api/v1/network/scan",
				"quick_scan":   "GET /api/v1/network/quick-scan?network=<CIDR>",
				"validate":     "GET /api/v1/network/validate?network=<CIDR>",
				"scan_device":  "GET /api/v1/device/<IP>",
			},
			"examples": gin.H{
				"scan_network": gin.H{
					"url": "POST /api/v1/network/scan",
					"body": gin.H{
						"network_range": "192.168.1.0/24",
						"communities":   []string{"public", "private"},
						"timeout":       5,
						"retries":       2,
					},
				},
				"quick_scan":  "GET /api/v1/network/quick-scan?network=192.168.1.0/24",
				"scan_device": "GET /api/v1/device/192.168.1.1?community=public&community=private",
			},
		})
	})

	return router
}

// LoggingMiddleware creates a custom logging middleware
func LoggingMiddleware() gin.HandlerFunc {
	logger := logrus.New()
	logger.SetLevel(logrus.InfoLevel)

	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery

		// Process request
		c.Next()

		// Log request
		end := time.Now()
		latency := end.Sub(start)

		if raw != "" {
			path = path + "?" + raw
		}

		logger.WithFields(logrus.Fields{
			"status_code": c.Writer.Status(),
			"latency":     latency,
			"client_ip":   c.ClientIP(),
			"method":      c.Request.Method,
			"path":        path,
		}).Info("HTTP Request")
	}
}

// RateLimitMiddleware creates a simple rate limiting middleware
func RateLimitMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Simple rate limiting can be implemented here
		// For production, consider using redis-based rate limiting
		c.Next()
	}
}
