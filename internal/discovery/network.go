package discovery

import (
	"fmt"
	"time"

	"network-discovery/internal/models"
	"network-discovery/internal/ports"
	"network-discovery/internal/scanner"
	"network-discovery/internal/snmp"

	"github.com/sirupsen/logrus"
)

type NetworkDiscovery struct {
	fullScanner *scanner.FullScanner
	logger      *logrus.Logger

	// Default SNMP communities to try
	defaultCommunities []string

	// Default timeout and retries
	defaultTimeout time.Duration
	defaultRetries int
	maxWorkers     int
}

func NewNetworkDiscovery() *NetworkDiscovery {
	logger := logrus.New()
	logger.SetLevel(logrus.InfoLevel)

	// Create SNMP client with default settings
	client := snmp.NewClient(time.Second*5, 2)

	// Create full scanner with 50 concurrent workers
	fullScanner := scanner.NewFullScanner(client, 50)

	return &NetworkDiscovery{
		fullScanner: fullScanner,
		logger:      logger,
		defaultCommunities: []string{
			"public",
			"private",
			"community",
			"admin",
		},
		defaultTimeout: time.Second * 5,
		defaultRetries: 2,
		maxWorkers:     50,
	}
}

func NewNetworkDiscoveryWithLogLevel(level logrus.Level) *NetworkDiscovery {
	logger := logrus.New()
	logger.SetLevel(level)
	logger.SetFormatter(&logrus.JSONFormatter{})

	// Create SNMP client with default settings and custom logger
	client := snmp.NewClientWithLogger(time.Second*5, 2, logger)

	// Create full scanner with 50 concurrent workers and custom logger
	fullScanner := scanner.NewFullScannerWithLogger(client, 50, logger)

	return &NetworkDiscovery{
		fullScanner: fullScanner,
		logger:      logger,
		defaultCommunities: []string{
			"public",
			"private",
			"community",
			"admin",
		},
		defaultTimeout: time.Second * 5,
		defaultRetries: 2,
		maxWorkers:     50,
	}
}

// PerformFullScan performs comprehensive network discovery using both SNMP and ARP
func (nd *NetworkDiscovery) PerformFullScan(req *models.ScanRequest) (*models.FullScanResult, error) {
	nd.logger.Infof("Starting full network discovery for range: %s", req.NetworkRange)

	// Use provided communities or default ones
	communities := req.Communities
	if len(communities) == 0 {
		communities = nd.defaultCommunities
		nd.logger.Infof("Using default communities: %v", communities)
	} else {
		nd.logger.Infof("Using provided communities: %v", communities)
	}

	// Update client settings if provided
	if req.Timeout > 0 || req.Retries > 0 {
		timeout := nd.defaultTimeout
		retries := nd.defaultRetries

		if req.Timeout > 0 {
			timeout = time.Duration(req.Timeout) * time.Second
			nd.logger.Infof("Using custom timeout: %d seconds", req.Timeout)
		}
		if req.Retries > 0 {
			retries = req.Retries
			nd.logger.Infof("Using custom retries: %d", req.Retries)
		}

		// Create new client with updated settings
		var client *snmp.Client
		if nd.logger.Level == logrus.DebugLevel {
			client = snmp.NewClientWithLogger(timeout, retries, nd.logger)
		} else {
			client = snmp.NewClient(timeout, retries)
		}
		nd.fullScanner = scanner.NewFullScannerWithLogger(client, nd.maxWorkers, nd.logger)
	}

	// Perform the scan based on scan type
	var topology *models.NetworkTopology
	var err error

	// Configure port scan toggle (default true)
	enablePortScan := true
	if req.EnablePortScan != nil {
		enablePortScan = *req.EnablePortScan
	}
	nd.fullScanner.SetPortScanEnabled(enablePortScan)

	switch req.ScanType {
	case "snmp":
		topology, err = nd.fullScanner.PerformSNMPScan(req.NetworkRange, communities)
	case "arp":
		topology, err = nd.fullScanner.PerformARPScan(req.NetworkRange)
	case "full", "":
		topology, err = nd.fullScanner.PerformFullScan(req.NetworkRange, communities)
	default:
		return nil, fmt.Errorf("invalid scan type: %s. Supported types: snmp, arp, full", req.ScanType)
	}

	if err != nil {
		return nil, fmt.Errorf("network scan failed: %v", err)
	}

	nd.logger.Infof("Discovery completed. Found %d devices (%d reachable, %d SNMP, %d ARP-only)",
		topology.TotalCount, topology.ReachableCount, topology.SNMPCount, topology.ARPCount)

	// Generate statistics
	statistics := nd.GetNetworkStatistics(topology)

	// Create scan info
	scanInfo := models.ScanInfo{
		ScanType:        req.ScanType,
		NetworkRange:    req.NetworkRange,
		SNMPCommunities: communities,
		Timeout:         req.Timeout,
		Retries:         req.Retries,
		WorkerCount:     nd.maxWorkers,
	}

	result := &models.FullScanResult{
		Topology:   topology,
		Statistics: statistics,
		ScanInfo:   scanInfo,
	}

	return result, nil
}

// DiscoverNetwork performs SNMP-only network discovery (backward compatibility)
func (nd *NetworkDiscovery) DiscoverNetwork(req *models.ScanRequest) (*models.NetworkTopology, error) {
	nd.logger.Infof("Starting SNMP network discovery for range: %s", req.NetworkRange)

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
		nd.fullScanner = scanner.NewFullScannerWithLogger(client, nd.maxWorkers, nd.logger)
	}

	// Configure port scan toggle (default true)
	enablePortScan := true
	if req.EnablePortScan != nil {
		enablePortScan = *req.EnablePortScan
	}
	nd.fullScanner.SetPortScanEnabled(enablePortScan)

	// Perform SNMP-only scan
	topology, err := nd.fullScanner.PerformSNMPScan(req.NetworkRange, communities)
	if err != nil {
		return nil, fmt.Errorf("network scan failed: %v", err)
	}

	nd.logger.Infof("Discovery completed. Found %d devices (%d reachable)",
		topology.TotalCount, topology.ReachableCount)

	return topology, nil
}

func (nd *NetworkDiscovery) DiscoverDevice(ip string, communities []string, enablePortScan bool) (*models.Device, error) {
	nd.logger.Infof("Discovering single device: %s", ip)

	if len(communities) == 0 {
		communities = nd.defaultCommunities
	}

	// Create a temporary SNMP client for single device query
	client := snmp.NewClientWithLogger(nd.defaultTimeout, nd.defaultRetries, nd.logger)
	device, err := client.QueryDevice(ip, communities)
	if err != nil {
		return nil, fmt.Errorf("device discovery failed: %v", err)
	}

	// Best-effort port scan for the single device
	_ = scanner.NewFullScannerWithLogger(client, nd.maxWorkers, nd.logger) // ensure consistency
	portScanner := ports.NewScannerWithLogger(5, nd.logger)
	if device != nil {
		if portsInfo, err := portScanner.ScanHost(device.IP); err == nil {
			device.OpenPorts = portsInfo
		} else {
			nd.logger.Debugf("Port scan failed for %s: %v", device.IP, err)
		}
	}

	return device, nil
}

func (nd *NetworkDiscovery) QuickDiscovery(networkRange string, communities []string) ([]string, error) {
	nd.logger.Infof("Starting quick discovery for range: %s", networkRange)

	if len(communities) == 0 {
		communities = nd.defaultCommunities
	}

	// Use SNMP scanner for quick discovery
	client := snmp.NewClientWithLogger(nd.defaultTimeout, nd.defaultRetries, nd.logger)
	scanner := snmp.NewScannerWithLogger(client, nd.maxWorkers, nd.logger)

	reachableIPs, err := scanner.QuickScan(networkRange, communities)
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
	stats["snmp_devices"] = topology.SNMPCount
	stats["arp_only_devices"] = topology.ARPCount

	// Vendor distribution
	vendorCount := make(map[string]int)
	scanMethodCount := make(map[string]int)
	macAddressCount := 0

	for _, device := range topology.Devices {
		if device.IsReachable {
			// Count vendors
			if device.Vendor != "" {
				vendorCount[device.Vendor]++
			} else {
				vendorCount["Unknown"]++
			}

			// Count scan methods
			scanMethodCount[device.ScanMethod]++

			// Count devices with MAC addresses
			if device.MACAddress != "" {
				macAddressCount++
			}
		}
	}

	stats["vendor_distribution"] = vendorCount
	stats["scan_method_distribution"] = scanMethodCount
	stats["devices_with_mac"] = macAddressCount

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
	stats["scan_method"] = topology.ScanMethod

	return stats
}

func (nd *NetworkDiscovery) ValidateNetworkRange(networkRange string) error {
	client := snmp.NewClientWithLogger(nd.defaultTimeout, nd.defaultRetries, nd.logger)
	scanner := snmp.NewScannerWithLogger(client, nd.maxWorkers, nd.logger)

	_, err := scanner.QuickScan(networkRange, []string{"public"})
	if err != nil {
		return fmt.Errorf("invalid network range: %v", err)
	}
	return nil
}
