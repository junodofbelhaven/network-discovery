package discovery

import (
	"fmt"
	"time"

	"network-discovery/internal/models"
	"network-discovery/internal/snmp"

	"github.com/sirupsen/logrus"
)

type NetworkDiscovery struct {
	scanner *snmp.Scanner
	logger  *logrus.Logger

	// Default SNMP communities to try
	defaultCommunities []string

	// Default timeout and retries
	defaultTimeout time.Duration
	defaultRetries int
}

func NewNetworkDiscovery() *NetworkDiscovery {
	logger := logrus.New()
	logger.SetLevel(logrus.InfoLevel)

	// Create SNMP client with default settings
	client := snmp.NewClient(time.Second*5, 2)

	// Create scanner with 50 concurrent workers
	scanner := snmp.NewScanner(client, 50)

	return &NetworkDiscovery{
		scanner: scanner,
		logger:  logger,
		defaultCommunities: []string{
			"public",
			"private",
			"community",
			"admin",
		},
		defaultTimeout: time.Second * 5,
		defaultRetries: 2,
	}
}

func NewNetworkDiscoveryWithLogLevel(level logrus.Level) *NetworkDiscovery {
	logger := logrus.New()
	logger.SetLevel(level)
	logger.SetFormatter(&logrus.JSONFormatter{})

	// Create SNMP client with default settings and custom logger
	client := snmp.NewClientWithLogger(time.Second*5, 2, logger)

	// Create scanner with 50 concurrent workers and custom logger
	scanner := snmp.NewScannerWithLogger(client, 50, logger)

	return &NetworkDiscovery{
		scanner: scanner,
		logger:  logger,
		defaultCommunities: []string{
			"public",
			"private",
			"community",
			"admin",
		},
		defaultTimeout: time.Second * 5,
		defaultRetries: 2,
	}
}

func (nd *NetworkDiscovery) DiscoverNetwork(req *models.ScanRequest) (*models.NetworkTopology, error) {
	nd.logger.Infof("Starting network discovery for range: %s", req.NetworkRange)

	// Use provided communities or default ones
	communities := req.Communities
	if len(communities) == 0 {
		communities = nd.defaultCommunities
		nd.logger.Infof("Using default communities: %v", communities)
	} else {
		nd.logger.Infof("Using provided communities: %v", communities)
	}

	// Update client settings if provided
	if req.Timeout > 0 {
		nd.logger.Infof("Using custom timeout: %d seconds", req.Timeout)
		var client *snmp.Client
		if nd.logger.Level == logrus.DebugLevel {
			client = snmp.NewClientWithLogger(time.Duration(req.Timeout)*time.Second, req.Retries, nd.logger)
		} else {
			client = snmp.NewClient(time.Duration(req.Timeout)*time.Second, req.Retries)
		}
		nd.scanner = snmp.NewScannerWithLogger(client, 50, nd.logger)
	}

	// Perform the scan
	topology, err := nd.scanner.ScanNetwork(req.NetworkRange, communities)
	if err != nil {
		return nil, fmt.Errorf("network scan failed: %v", err)
	}

	nd.logger.Infof("Discovery completed. Found %d devices (%d reachable)",
		topology.TotalCount, topology.ReachableCount)

	return topology, nil
}

func (nd *NetworkDiscovery) DiscoverDevice(ip string, communities []string) (*models.Device, error) {
	nd.logger.Infof("Discovering single device: %s", ip)

	if len(communities) == 0 {
		communities = nd.defaultCommunities
	}

	device, err := nd.scanner.ScanSingleDevice(ip, communities)
	if err != nil {
		return nil, fmt.Errorf("device discovery failed: %v", err)
	}

	return device, nil
}

func (nd *NetworkDiscovery) QuickDiscovery(networkRange string, communities []string) ([]string, error) {
	nd.logger.Infof("Starting quick discovery for range: %s", networkRange)

	if len(communities) == 0 {
		communities = nd.defaultCommunities
	}

	reachableIPs, err := nd.scanner.QuickScan(networkRange, communities)
	if err != nil {
		return nil, fmt.Errorf("quick discovery failed: %v", err)
	}

	nd.logger.Infof("Quick discovery found %d reachable devices", len(reachableIPs))
	return reachableIPs, nil
}

func (nd *NetworkDiscovery) GetNetworkStatistics(topology *models.NetworkTopology) map[string]interface{} {
	stats := make(map[string]interface{})

	// Basic counts
	stats["total_devices"] = topology.TotalCount
	stats["reachable_devices"] = topology.ReachableCount
	stats["unreachable_devices"] = topology.TotalCount - topology.ReachableCount

	// Vendor distribution
	vendorCount := make(map[string]int)
	for _, device := range topology.Devices {
		if device.IsReachable && device.Vendor != "" {
			vendorCount[device.Vendor]++
		} else if device.IsReachable {
			vendorCount["Unknown"]++
		}
	}
	stats["vendor_distribution"] = vendorCount

	// Response time statistics
	var responseTimes []int64
	var totalResponseTime int64
	for _, device := range topology.Devices {
		if device.IsReachable {
			responseTimes = append(responseTimes, device.ResponseTime)
			totalResponseTime += device.ResponseTime
		}
	}

	if len(responseTimes) > 0 {
		stats["avg_response_time_ms"] = totalResponseTime / int64(len(responseTimes))

		// Find min and max response times
		minTime, maxTime := responseTimes[0], responseTimes[0]
		for _, time := range responseTimes {
			if time < minTime {
				minTime = time
			}
			if time > maxTime {
				maxTime = time
			}
		}
		stats["min_response_time_ms"] = minTime
		stats["max_response_time_ms"] = maxTime
	}

	// Scan performance
	stats["scan_duration_ms"] = topology.ScanDuration
	stats["scan_time"] = topology.ScanTime.Format(time.RFC3339)

	return stats
}

func (nd *NetworkDiscovery) ValidateNetworkRange(networkRange string) error {
	_, err := nd.scanner.QuickScan(networkRange, []string{"public"})
	if err != nil {
		return fmt.Errorf("invalid network range: %v", err)
	}
	return nil
}
