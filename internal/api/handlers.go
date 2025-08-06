package api

import (
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

// ScanNetwork handles network scanning requests
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

	h.logger.Infof("Received scan request for network: %s", req.NetworkRange)

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

// GetHealth handles health check requests
func (h *Handlers) GetHealth(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":  "healthy",
		"service": "network-discovery",
		"version": "1.0.0",
	})
}

// GetVersion handles version requests
func (h *Handlers) GetVersion(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"version":    "1.0.0",
		"build_time": "2024-01-01T00:00:00Z",
		"go_version": "1.21",
	})
}
