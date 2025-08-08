package api

import (
	"fmt"
	"net/http"

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

	// Set defaults if not provided
	if req.Timeout == 0 {
		req.Timeout = 5
	}
	if req.Retries == 0 {
		req.Retries = 2
	}
	if req.ScanType == "" {
		req.ScanType = "full"
	}

	h.logger.Infof("Received full scan request for network: %s (type: %s)", req.NetworkRange, req.ScanType)

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

	// Set defaults if not provided
	if req.Timeout == 0 {
		req.Timeout = 5
	}
	if req.Retries == 0 {
		req.Retries = 2
	}

	h.logger.Infof("Received SNMP scan request for network: %s", req.NetworkRange)

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

	// Set defaults if not provided
	if req.Timeout == 0 {
		req.Timeout = 5
	}
	if req.Retries == 0 {
		req.Retries = 2
	}

	h.logger.Infof("Received %s scan request for network: %s", scanType, req.NetworkRange)

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

// QuickScan handles quick network scanning requests
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
		},
		"arp": gin.H{
			"name":         "ARP Scan",
			"description":  "Discovers devices using ARP (Address Resolution Protocol). Finds all devices that respond to ping and have ARP entries.",
			"requirements": []string{"Devices must be on the same network segment", "Devices must respond to ping"},
			"advantages":   []string{"Discovers all IP-enabled devices", "No special configuration required", "Fast discovery"},
			"limitations":  []string{"Limited device information", "Only provides IP and MAC addresses", "May miss some devices behind firewalls"},
		},
		"full": gin.H{
			"name":         "Full Scan (SNMP + ARP)",
			"description":  "Combines both SNMP and ARP scanning methods for comprehensive network discovery. Provides the most complete view of network devices.",
			"requirements": []string{"Network access to target range"},
			"advantages":   []string{"Most comprehensive discovery", "Combines detailed SNMP info with broad ARP coverage", "Merges MAC addresses for SNMP devices"},
			"limitations":  []string{"Takes longer than individual scans", "Higher network traffic"},
		},
	}

	c.JSON(http.StatusOK, gin.H{
		"scan_methods": methods,
		"default":      "full",
		"recommended":  "full",
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
			"vendor_detection": "Device manufacturer identification",
		},
	})
}
