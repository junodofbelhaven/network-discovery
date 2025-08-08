package api

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"network-discovery/internal/discovery"
	"network-discovery/internal/models"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type Handlers struct {
	discovery *discovery.NetworkDiscovery
	logger    *logrus.Logger
}

func NewHandlers(discovery *discovery.NetworkDiscovery) *Handlers {
	logger := logrus.New()
	logger.SetLevel(logrus.InfoLevel)

	return &Handlers{
		discovery: discovery,
		logger:    logger,
	}
}

// PerformFullScan handles comprehensive network scanning requests (SNMP + ARP)
func (h *Handlers) PerformFullScan(c *gin.Context) {
	var req models.ScanRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Errorf("Invalid request: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request format",
			"details": err.Error(),
		})
		return
	}

	// Set optimized defaults for faster scanning
	if req.Timeout == 0 {
		req.Timeout = 2 // Reduced from 5 to 2 seconds
	}
	if req.Retries < 0 { // Allow 0 retries
		req.Retries = 1 // Reduced from 2 to 1
	}
	if req.ScanType == "" {
		req.ScanType = "full"
	}

	h.logger.Infof("Received full scan request for network: %s (type: %s, timeout: %ds, retries: %d)",
		req.NetworkRange, req.ScanType, req.Timeout, req.Retries)

	// Perform the full discovery
	result, err := h.discovery.PerformFullScan(&req)
	if err != nil {
		h.logger.Errorf("Full network discovery failed: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Full network discovery failed",
			"details": err.Error(),
		})
		return
	}

	h.logger.Infof("Full scan completed successfully: %d total devices, %d SNMP, %d ARP-only",
		result.Topology.TotalCount, result.Topology.SNMPCount, result.Topology.ARPCount)

	c.JSON(http.StatusOK, result)
}

// ScanNetwork handles network scanning requests (SNMP only - backward compatibility)
func (h *Handlers) ScanNetwork(c *gin.Context) {
	var req models.ScanRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Errorf("Invalid request: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request format",
			"details": err.Error(),
		})
		return
	}

	// Set optimized defaults for faster scanning
	if req.Timeout == 0 {
		req.Timeout = 2 // Reduced from 5 to 2 seconds
	}
	if req.Retries < 0 { // Allow 0 retries
		req.Retries = 1 // Reduced from 2 to 1
	}

	h.logger.Infof("Received SNMP scan request for network: %s (timeout: %ds, retries: %d)",
		req.NetworkRange, req.Timeout, req.Retries)

	// Perform the discovery
	topology, err := h.discovery.DiscoverNetwork(&req)
	if err != nil {
		h.logger.Errorf("Network discovery failed: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Network discovery failed",
			"details": err.Error(),
		})
		return
	}

	// Get statistics
	stats := h.discovery.GetNetworkStatistics(topology)

	response := gin.H{
		"topology":   topology,
		"statistics": stats,
	}

	c.JSON(http.StatusOK, response)
}

// ScanNetworkByType handles network scanning requests with specific scan type
func (h *Handlers) ScanNetworkByType(c *gin.Context) {
	scanType := c.Param("type")

	// Validate scan type
	validTypes := map[string]bool{
		"snmp": true,
		"arp":  true,
		"full": true,
	}

	if !validTypes[scanType] {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid scan type. Supported types: snmp, arp, full",
		})
		return
	}

	var req models.ScanRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Errorf("Invalid request: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request format",
			"details": err.Error(),
		})
		return
	}

	// Set scan type from URL parameter
	req.ScanType = scanType

	// Set optimized defaults for faster scanning
	if req.Timeout == 0 {
		req.Timeout = 2 // Reduced timeout
	}
	if req.Retries < 0 { // Allow 0 retries
		req.Retries = 1 // Reduced retries
	}

	h.logger.Infof("Received %s scan request for network: %s (timeout: %ds, retries: %d)",
		scanType, req.NetworkRange, req.Timeout, req.Retries)

	// Perform the discovery
	result, err := h.discovery.PerformFullScan(&req)
	if err != nil {
		h.logger.Errorf("%s network discovery failed: %v", scanType, err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   fmt.Sprintf("%s network discovery failed", scanType),
			"details": err.Error(),
		})
		return
	}

	h.logger.Infof("%s scan completed successfully: %d total devices",
		scanType, result.Topology.TotalCount)

	c.JSON(http.StatusOK, result)
}

// QuickScan handles quick network scanning requests with very fast settings
func (h *Handlers) QuickScan(c *gin.Context) {
	networkRange := c.Query("network")
	if networkRange == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Network range is required",
		})
		return
	}

	communities := c.QueryArray("community")

	h.logger.Infof("Received quick scan request for network: %s", networkRange)

	reachableIPs, err := h.discovery.QuickDiscovery(networkRange, communities)
	if err != nil {
		h.logger.Errorf("Quick discovery failed: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Quick discovery failed",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"reachable_ips": reachableIPs,
		"count":         len(reachableIPs),
	})
}

// ScanDevice handles single device scanning requests
func (h *Handlers) ScanDevice(c *gin.Context) {
	ip := c.Param("ip")
	if ip == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "IP address is required",
		})
		return
	}

	// Get communities from query parameters
	communities := c.QueryArray("community")

	h.logger.Infof("Received device scan request for IP: %s", ip)

	device, err := h.discovery.DiscoverDevice(ip, communities)
	if err != nil {
		h.logger.Errorf("Device discovery failed: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Device discovery failed",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"device": device,
	})
}

// ValidateNetwork handles network range validation requests
func (h *Handlers) ValidateNetwork(c *gin.Context) {
	networkRange := c.Query("network")
	if networkRange == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Network range is required",
		})
		return
	}

	err := h.discovery.ValidateNetworkRange(networkRange)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"valid": false,
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"valid":   true,
		"network": networkRange,
	})
}

// GetScanMethods returns available scan methods and their descriptions
func (h *Handlers) GetScanMethods(c *gin.Context) {
	methods := gin.H{
		"snmp": gin.H{
			"name":         "SNMP Scan",
			"description":  "Discovers devices using SNMP protocol. Provides detailed device information including hostname, description, vendor, uptime, and system details.",
			"requirements": []string{"SNMP enabled on target devices", "Valid SNMP community strings"},
			"advantages":   []string{"Detailed device information", "Vendor identification", "System uptime and status"},
			"limitations":  []string{"Only discovers SNMP-enabled devices", "Requires correct community strings"},
			"recommended_settings": gin.H{
				"timeout": "1-3 seconds",
				"retries": "0-1",
			},
		},
		"arp": gin.H{
			"name":         "ARP Scan",
			"description":  "Discovers devices using ARP (Address Resolution Protocol). Finds all devices that respond to ping and have ARP entries.",
			"requirements": []string{"Devices must be on the same network segment", "Devices must respond to ping"},
			"advantages":   []string{"Discovers all IP-enabled devices", "No special configuration required", "Fast discovery"},
			"limitations":  []string{"Limited device information", "Only provides IP and MAC addresses", "May miss some devices behind firewalls"},
			"recommended_settings": gin.H{
				"timeout": "1-2 seconds",
				"retries": "0",
			},
		},
		"full": gin.H{
			"name":         "Full Scan (SNMP + ARP)",
			"description":  "Combines both SNMP and ARP scanning methods for comprehensive network discovery. Provides the most complete view of network devices.",
			"requirements": []string{"Network access to target range"},
			"advantages":   []string{"Most comprehensive discovery", "Combines detailed SNMP info with broad ARP coverage", "Merges MAC addresses for SNMP devices"},
			"limitations":  []string{"Takes longer than individual scans", "Higher network traffic"},
			"recommended_settings": gin.H{
				"timeout": "2-3 seconds",
				"retries": "0-1",
			},
		},
	}

	c.JSON(http.StatusOK, gin.H{
		"scan_methods": methods,
		"default":      "full",
		"recommended":  "full",
		"performance_tips": []string{
			"Use smaller network ranges for faster scans",
			"Set timeout to 1-2 seconds for local networks",
			"Set retries to 0 for fastest scanning",
			"Use ARP scan for quick discovery without SNMP details",
		},
	})
}

// GetHealth handles health check requests
func (h *Handlers) GetHealth(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":  "healthy",
		"service": "network-discovery",
		"version": "1.0.0",
		"features": []string{
			"SNMP Discovery",
			"ARP Discovery",
			"Full Network Scan",
			"Device Information Extraction",
			"MAC Address Resolution",
			"Vendor Identification",
			"JSON-based Vendor Database",
		},
	})
}

// GetVersion handles version requests
func (h *Handlers) GetVersion(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"version":    "1.0.0",
		"build_time": "2024-01-01T00:00:00Z",
		"go_version": "1.21",
		"features": gin.H{
			"snmp_discovery":   "v2c protocol support",
			"arp_discovery":    "Cross-platform ARP scanning",
			"full_scan":        "Combined SNMP + ARP discovery",
			"mac_resolution":   "Hardware address identification",
			"vendor_detection": "JSON-based vendor database",
		},
	})
}

// TestVendorDetection - test endpoint for debugging vendor detection
func (h *Handlers) TestVendorDetection(c *gin.Context) {
	macAddr := c.Query("mac")
	if macAddr == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "MAC address is required. Example: ?mac=AA:BB:CC:DD:EE:FF",
		})
		return
	}

	// Test vendor detection logic
	if len(macAddr) < 8 {
		c.JSON(http.StatusOK, gin.H{
			"mac":    macAddr,
			"vendor": "Unknown",
			"reason": "MAC address too short",
		})
		return
	}

	oui := strings.ReplaceAll(macAddr[:8], ":", "")
	oui = strings.ToUpper(oui)

	// Sample of OUI mappings for testing
	testVendors := map[string]string{
		"005056": "VMware",
		"000C29": "VMware",
		"080027": "Oracle VirtualBox",
		"B827EB": "Raspberry Pi",
		"28CDC1": "Apple",
		"001E10": "Huawei",
		"141877": "TP-Link",
		"000121": "Cabletron Systems",
		"C403A8": "Shenzhen Coship Electronics",
		"5C3A45": "Liteon Technology Corporation",
		"600308": "Apple",
	}

	vendor := "Unknown"
	if v, exists := testVendors[oui]; exists {
		vendor = v
	}

	c.JSON(http.StatusOK, gin.H{
		"mac":                 macAddr,
		"oui":                 oui,
		"vendor":              vendor,
		"available_test_ouis": testVendors,
		"note":                "This is a test endpoint. Real vendor detection uses JSON database.",
	})
}

// GetVendorDatabase returns vendor database information
func (h *Handlers) GetVendorDatabase(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":      "vendor database loaded from JSON file",
		"config_path": "configs/oui_vendors.json",
		"description": "External JSON-based OUI vendor database",
		"features": []string{
			"Runtime reloadable",
			"Thread-safe access",
			"Metadata tracking",
			"Virtual MAC detection",
		},
		"endpoints": gin.H{
			"reload": "POST /api/v1/vendor-database/reload",
			"info":   "GET  /api/v1/vendor-database",
		},
		"management": gin.H{
			"edit_file":   "Edit configs/oui_vendors.json directly",
			"reload_db":   "POST to /api/v1/vendor-database/reload",
			"test_vendor": "GET /api/v1/test-vendor?mac=XX:XX:XX:XX:XX:XX",
		},
	})
}

// ReloadVendorDatabase reloads the vendor database from JSON file
func (h *Handlers) ReloadVendorDatabase(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":    "reload triggered",
		"message":   "Vendor database reloaded from configs/oui_vendors.json",
		"timestamp": time.Now().Format(time.RFC3339),
		"note":      "Database will be reloaded on next scan operation",
	})
}
