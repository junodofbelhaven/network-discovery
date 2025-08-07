package snmp

import (
	"fmt"
	"net"
	"sync"
	"time"

	"network-discovery/internal/models"

	"github.com/sirupsen/logrus"
)

type Scanner struct {
	client     *Client
	logger     *logrus.Logger
	maxWorkers int
}

func NewScanner(client *Client, maxWorkers int) *Scanner {
	logger := logrus.New()
	logger.SetLevel(logrus.InfoLevel)

	return &Scanner{
		client:     client,
		logger:     logger,
		maxWorkers: maxWorkers,
	}
}

func NewScannerWithLogger(client *Client, maxWorkers int, logger *logrus.Logger) *Scanner {
	return &Scanner{
		client:     client,
		logger:     logger,
		maxWorkers: maxWorkers,
	}
}

func (s *Scanner) ScanNetwork(networkRange string, communities []string) (*models.NetworkTopology, error) {
	start := time.Now()

	s.logger.Infof("Starting network scan for range: %s", networkRange)

	// Parse network range
	ips, err := s.parseNetworkRange(networkRange)
	if err != nil {
		return nil, fmt.Errorf("failed to parse network range: %v", err)
	}

	s.logger.Infof("Scanning %d IP addresses", len(ips))

	// Create channels for work distribution
	ipChan := make(chan string, len(ips))
	resultChan := make(chan *models.Device, len(ips))

	// Add IPs to channel
	for _, ip := range ips {
		ipChan <- ip
	}
	close(ipChan)

	// Start workers
	var wg sync.WaitGroup
	workers := s.maxWorkers
	if workers > len(ips) {
		workers = len(ips)
	}

	s.logger.Infof("Starting %d workers for scanning", workers)

	for i := 0; i < workers; i++ {
		wg.Add(1)
		go s.worker(ipChan, resultChan, communities, &wg)
	}

	// Wait for all workers to complete
	go func() {
		wg.Wait()
		close(resultChan)
	}()

	// Collect results
	var devices []models.Device
	reachableCount := 0

	for device := range resultChan {
		if device != nil {
			devices = append(devices, *device)
			if device.IsReachable {
				reachableCount++
				s.logger.Infof("Found SNMP device: %s (%s)", device.IP, device.Hostname)
			}
		}
	}

	scanDuration := time.Since(start)

	topology := &models.NetworkTopology{
		Devices:        devices,
		TotalCount:     len(devices),
		ReachableCount: reachableCount,
		ScanTime:       start,
		ScanDuration:   scanDuration.Milliseconds(),
	}

	s.logger.Infof("Scan completed in %v. Found %d devices (%d reachable)",
		scanDuration, len(devices), reachableCount)

	return topology, nil
}

func (s *Scanner) worker(ipChan <-chan string, resultChan chan<- *models.Device, communities []string, wg *sync.WaitGroup) {
	defer wg.Done()

	for ip := range ipChan {
		s.logger.Debugf("Scanning IP: %s", ip)

		device, err := s.client.QueryDevice(ip, communities)
		if err != nil {
			s.logger.Debugf("Failed to query %s: %v", ip, err)
			// Sadece reachable olanları ekleyelim
			// resultChan <- nil // Unreachable cihazları ekleme
		} else {
			s.logger.Debugf("Successfully queried %s", ip)
			resultChan <- device
		}
	}
}

func (s *Scanner) ScanSingleDevice(ip string, communities []string) (*models.Device, error) {
	s.logger.Infof("Scanning single device: %s", ip)

	device, err := s.client.QueryDevice(ip, communities)
	if err != nil {
		return nil, fmt.Errorf("failed to scan device %s: %v", ip, err)
	}

	return device, nil
}

func (s *Scanner) QuickScan(networkRange string, communities []string) ([]string, error) {
	ips, err := s.parseNetworkRange(networkRange)
	if err != nil {
		return nil, fmt.Errorf("failed to parse network range: %v", err)
	}

	var reachableIPs []string
	var mu sync.Mutex

	ipChan := make(chan string, len(ips))
	for _, ip := range ips {
		ipChan <- ip
	}
	close(ipChan)

	var wg sync.WaitGroup
	workers := s.maxWorkers
	if workers > len(ips) {
		workers = len(ips)
	}

	for i := 0; i < workers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for ip := range ipChan {
				if s.client.IsDeviceReachable(ip, communities) {
					mu.Lock()
					reachableIPs = append(reachableIPs, ip)
					mu.Unlock()
					s.logger.Debugf("Quick scan found reachable device: %s", ip)
				}
			}
		}()
	}

	wg.Wait()

	s.logger.Infof("Quick scan found %d reachable devices", len(reachableIPs))
	return reachableIPs, nil
}

func (s *Scanner) parseNetworkRange(networkRange string) ([]string, error) {
	_, ipNet, err := net.ParseCIDR(networkRange)
	if err != nil {
		return nil, fmt.Errorf("invalid CIDR notation: %v", err)
	}

	var ips []string

	// Generate all IPs in the network
	for ip := ipNet.IP.Mask(ipNet.Mask); ipNet.Contains(ip); s.incrementIP(ip) {
		// Skip network and broadcast addresses
		if !s.isNetworkOrBroadcast(ip, ipNet) {
			ips = append(ips, ip.String())
		}
	}

	s.logger.Debugf("Generated %d IPs from range %s", len(ips), networkRange)
	return ips, nil
}

func (s *Scanner) incrementIP(ip net.IP) {
	for j := len(ip) - 1; j >= 0; j-- {
		ip[j]++
		if ip[j] > 0 {
			break
		}
	}
}

func (s *Scanner) isNetworkOrBroadcast(ip net.IP, ipNet *net.IPNet) bool {
	// Check if it's the network address
	if ip.Equal(ipNet.IP.Mask(ipNet.Mask)) {
		return true
	}

	// Check if it's the broadcast address
	broadcast := make(net.IP, len(ipNet.IP))
	copy(broadcast, ipNet.IP.Mask(ipNet.Mask))

	for i := range broadcast {
		broadcast[i] |= ^ipNet.Mask[i]
	}

	return ip.Equal(broadcast)
}
