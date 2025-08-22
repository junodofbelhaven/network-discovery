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
		v1.GET("/scan-methods", handlers.GetScanMethods)

		// Vendor database endpoints
		vendor := v1.Group("/vendor-database")
		{
			vendor.GET("/", handlers.GetVendorDatabase)
			vendor.POST("/reload", handlers.ReloadVendorDatabase)
		}

		// Network discovery endpoints
		network := v1.Group("/network")
		{
			// Full scan endpoint - primary endpoint for comprehensive discovery
			network.POST("/full-scan", handlers.PerformFullScan)

			// Scan by type - allows specifying scan method in URL
			network.POST("/scan/:type", handlers.ScanNetworkByType)

			// Utility endpoints
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
	router.Static("/assets", "./frontend-build/dist/assets")
	router.LoadHTMLGlob("./frontend-build/dist/index.html")

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
			"features": []string{
				"SNMP Network Discovery",
				"ARP Network Scanning",
				"Full Network Scan (SNMP + ARP)",
				"Device Information Extraction",
				"MAC Address Resolution",
				"Vendor Identification",
			},
			"endpoints": gin.H{
				"health":       "GET  /api/v1/health",
				"version":      "GET  /api/v1/version",
				"scan_methods": "GET  /api/v1/scan-methods",
				"full_scan":    "POST /api/v1/network/full-scan",
				"scan_by_type": "POST /api/v1/network/scan/{type}",
				"legacy_scan":  "POST /api/v1/network/scan",
				"quick_scan":   "GET  /api/v1/network/quick-scan?network=<CIDR>",
				"validate":     "GET  /api/v1/network/validate?network=<CIDR>",
				"scan_device":  "GET  /api/v1/device/<IP>",
			},
			"scan_types": []string{"snmp", "arp", "full"},
			"examples": gin.H{
				"full_scan": gin.H{
					"url":         "POST /api/v1/network/full-scan",
					"description": "Comprehensive network discovery using both SNMP and ARP",
					"body": gin.H{
						"network_range": "192.168.1.0/24",
						"communities":   []string{"public", "private"},
						"timeout":       5,
						"retries":       2,
						"scan_type":     "full",
					},
				},
				"snmp_scan": gin.H{
					"url":         "POST /api/v1/network/scan/snmp",
					"description": "SNMP-only network discovery",
					"body": gin.H{
						"network_range": "192.168.1.0/24",
						"communities":   []string{"public", "private"},
						"timeout":       5,
						"retries":       2,
					},
				},
				"arp_scan": gin.H{
					"url":         "POST /api/v1/network/scan/arp",
					"description": "ARP-only network discovery",
					"body": gin.H{
						"network_range": "192.168.1.0/24",
						"timeout":       5,
						"retries":       2,
					},
				},
				"quick_scan":  "GET /api/v1/network/quick-scan?network=192.168.1.0/24",
				"scan_device": "GET /api/v1/device/192.168.1.1?community=public&community=private",
			},
			"response_format": gin.H{
				"full_scan_response": gin.H{
					"topology": gin.H{
						"devices":          "array of discovered devices",
						"total_count":      "total number of devices found",
						"reachable_count":  "number of reachable devices",
						"snmp_count":       "number of SNMP-enabled devices",
						"arp_count":        "number of ARP-only devices",
						"scan_method":      "scan method used (SNMP/ARP/FULL)",
						"scan_duration_ms": "scan duration in milliseconds",
					},
					"statistics": gin.H{
						"vendor_distribution":      "device vendor breakdown",
						"scan_method_distribution": "how devices were discovered",
						"devices_with_mac":         "number of devices with MAC addresses",
						"response_time_stats":      "network response statistics",
					},
					"scan_info": gin.H{
						"scan_type":        "type of scan performed",
						"network_range":    "scanned network range",
						"snmp_communities": "SNMP communities used",
						"worker_count":     "number of concurrent workers",
					},
				},
				"device_properties": gin.H{
					"ip":            "IP address",
					"mac_address":   "MAC address (if available)",
					"hostname":      "device hostname",
					"description":   "SNMP system description",
					"vendor":        "detected vendor",
					"uptime":        "system uptime",
					"response_time": "response time in milliseconds",
					"scan_method":   "discovery method (SNMP/ARP/COMBINED)",
				},
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
