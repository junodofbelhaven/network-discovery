package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"runtime"
	"syscall"
	"time"

	"network-discovery/internal/api"
	"network-discovery/internal/discovery"

	"github.com/sirupsen/logrus"
)

var (
	port       = flag.String("port", "8080", "Server port")
	host       = flag.String("host", "0.0.0.0", "Server host")
	logLevel   = flag.String("log-level", "info", "Log level (debug, info, warn, error)")
	configPath = flag.String("config", "configs/oui_vendors.json", "Path to OUI vendors JSON file")
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

	// Create HTTP server with increased timeouts for long scans
	server := &http.Server{
		Addr:         fmt.Sprintf("%s:%s", *host, *port),
		Handler:      router,
		ReadTimeout:  5 * time.Minute,  // Increased from 30s to 5 minutes
		WriteTimeout: 5 * time.Minute,  // Increased from 30s to 5 minutes
		IdleTimeout:  10 * time.Minute, // Increased from 60s to 10 minutes
	}

	// Start server in a goroutine
	go func() {
		logger.Infof("Server starting on %s:%s", *host, *port)

		go openBrowser("http://localhost:8080/index")

		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatalf("Server failed to start: %v", err)
		}
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

func openBrowser(url string) {
	var err error

	switch runtime.GOOS {
	case "windows":
		err = exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
	case "linux":
		err = exec.Command("xdg-open", url).Start()
	case "darwin":
		err = exec.Command("open", url).Start()
	default:
		err = fmt.Errorf("unsupported platform")
	}
	if err != nil {
		fmt.Println("Failed to open browser:", err)
	}
}

func printStartupInfo(host, port string) {
	fmt.Printf(`
🌐 Network Discovery Service Started Successfully!

📡 Server Information:
   • Host: %s
   • Port: %s
   • URL: http://%s:%s

🔗 API Endpoints:
   • Health Check:       GET  /api/v1/health
   • Version Info:       GET  /api/v1/version
   • Scan Methods:       GET  /api/v1/scan-methods
   • Full Network Scan:  POST /api/v1/network/full-scan
   • SNMP Scan:          POST /api/v1/network/scan/snmp
   • ARP Scan:           POST /api/v1/network/scan/arp
   • Legacy SNMP Scan:   POST /api/v1/network/scan
   • Quick Scan:         GET  /api/v1/network/quick-scan?network=<CIDR>
   • Validate Range:     GET  /api/v1/network/validate?network=<CIDR>
   • Device Scan:        GET  /api/v1/device/<IP>


   📋 Example Usage (Windows Command Prompt):

   🔍 Full Scan (SNMP + ARP):
   curl -X POST http://%s:%s/api/v1/network/full-scan ^
     -H "Content-Type: application/json" ^
     -d "{
       \"network_range\": \"192.168.1.0/24\",
       \"communities\": [\"public\", \"private\"],
       \"timeout\": 5,
       \"retries\": 2,
       \"scan_type\": \"full\"
     }"

   📡 SNMP Only Scan:
   curl -X POST http://%s:%s/api/v1/network/scan/snmp ^
     -H "Content-Type: application/json" ^
     -d "{
       \"network_range\": \"192.168.1.0/24\",
       \"communities\": [\"public\", \"private\"],
       \"timeout\": 5,
       \"retries\": 2
     }"

   🌐 ARP Only Scan:
   curl -X POST http://%s:%s/api/v1/network/scan/arp ^
     -H "Content-Type: application/json" ^
     -d "{
       \"network_range\": \"192.168.1.0/24\",
       \"timeout\": 5,
       \"retries\": 2
     }"
	 
   ⚡ Quick Scan:
   curl "http://%s:%s/api/v1/network/quick-scan?network=192.168.1.0/24"
   
   🖥️  Single Device:
   curl "http://%s:%s/api/v1/device/192.168.1.1?community=public"

🛠️  Features:
   • 🔍 Full Network Discovery (SNMP + ARP)
   • 📡 SNMP v2c protocol support
   • 🌐 ARP-based device detection
   • 🏃 Concurrent network scanning
   • 📄 Device information discovery
   • 🏭 Vendor detection (SNMP description + MAC OUI)
   • 🕒 Response time measurement
   • 🗺️  Network topology analysis
   • 🔗 MAC address resolution
   • 📊 Comprehensive statistics

📝 Scan Types:
   • Full Scan: Combines SNMP and ARP for maximum device discovery
   • SNMP Scan: Detailed information from SNMP-enabled devices
   • ARP Scan: Broad discovery of all IP-enabled devices

🌐 Web Interface: http://%s:%s/index

📝 Logs: Check console output for detailed scanning information
🔧 Configuration: Use command line flags to customize settings

Ready to discover your network! 🚀
`, host, port, host, port, host, port, host, port, host, port, host, port, host, port, host, port)
}
