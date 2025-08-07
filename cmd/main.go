package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"network-discovery/internal/api"
	"network-discovery/internal/discovery"

	"github.com/sirupsen/logrus"
)

var (
	port     = flag.String("port", "8080", "Server port")
	host     = flag.String("host", "0.0.0.0", "Server host")
	logLevel = flag.String("log-level", "info", "Log level (debug, info, warn, error)")
)

func main() {
	flag.Parse()

	// Setup logger
	logger := logrus.New()
	level, err := logrus.ParseLevel(*logLevel)
	if err != nil {
		logger.Fatal("Invalid log level")
	}
	logger.SetLevel(level)
	logger.SetFormatter(&logrus.JSONFormatter{})

	logger.Infof("Starting Network Discovery Service with log level: %s", *logLevel)

	// Create network discovery service with custom log level
	networkDiscovery := discovery.NewNetworkDiscoveryWithLogLevel(level)

	// Setup routes
	router := api.SetupRoutes(networkDiscovery)

	// Create HTTP server
	server := &http.Server{
		Addr:         fmt.Sprintf("%s:%s", *host, *port),
		Handler:      router,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Start server in a goroutine
	go func() {
		logger.Infof("Server starting on %s:%s", *host, *port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatalf("Server failed to start: %v", err)
		}
		fmt.Printf("log leve: %s", *logLevel)
	}()

	// Print startup information
	printStartupInfo(*host, *port)

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("Shutting down server...")

	// Give outstanding requests 30 seconds to complete
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		logger.Fatalf("Server forced to shutdown: %v", err)
	}

	logger.Info("Server exited")
}

func printStartupInfo(host, port string) {
	fmt.Printf(`
🌐 Network Discovery Service Started Successfully!

📡 Server Information:
   • Host: %s
   • Port: %s
   • URL: http://%s:%s

🔗 API Endpoints:
   • Health Check:    GET  /api/v1/health
   • Version Info:    GET  /api/v1/version
   • Network Scan:    POST /api/v1/network/scan
   • Quick Scan:      GET  /api/v1/network/quick-scan?network=<CIDR>
   • Validate Range:  GET  /api/v1/network/validate?network=<CIDR>
   • Device Scan:     GET  /api/v1/device/<IP>

📋 Example Usage:
   curl -X POST http://%s:%s/api/v1/network/scan \
     -H "Content-Type: application/json" \
     -d '{
       "network_range": "192.168.1.0/24",
       "communities": ["public", "private"],
       "timeout": 5,
       "retries": 2
     }'

   curl "http://%s:%s/api/v1/network/quick-scan?network=192.168.1.0/24"
   
   curl "http://%s:%s/api/v1/device/192.168.1.1?community=public"

🛠️  Features:
   • SNMP v2c protocol support
   • Concurrent network scanning
   • Device information discovery
   • Vendor detection
   • Response time measurement
   • Network topology analysis

📝 Logs: Check console output for detailed scanning information
🔧 Configuration: Use command line flags to customize settings

Ready to discover your network! 🚀
`, host, port, host, port, host, port, host, port, host, port)
}
