package arp

import (
	"fmt"
	"net"
	"os/exec"
	"regexp"
	"runtime"
	"strings"
	"sync"
	"time"

	"network-discovery/internal/models"

	"github.com/sirupsen/logrus"
)

type Scanner struct {
	logger        *logrus.Logger
	maxWorkers    int
	vendorManager *VendorManager
}

func NewScanner(maxWorkers int) *Scanner {
	logger := logrus.New()
	logger.SetLevel(logrus.InfoLevel)

	return &Scanner{
		logger:        logger,
		maxWorkers:    maxWorkers,
		vendorManager: NewVendorManager("", logger), // Default config path
	}
}

func NewScannerWithLogger(maxWorkers int, logger *logrus.Logger) *Scanner {
	return &Scanner{
		logger:        logger,
		maxWorkers:    maxWorkers,
		vendorManager: NewVendorManager("", logger), // Default config path
	}
}

func NewScannerWithConfig(maxWorkers int, logger *logrus.Logger, configPath string) *Scanner {
	return &Scanner{
		logger:        logger,
		maxWorkers:    maxWorkers,
		vendorManager: NewVendorManager(configPath, logger),
	}
}

// ScanNetwork performs ARP scan on the given network range
func (s *Scanner) ScanNetwork(networkRange string) ([]*models.Device, error) {
	start := time.Now()
	s.logger.Infof("Starting ARP scan for range: %s", networkRange)

	// Parse network range
	ips, err := s.parseNetworkRange(networkRange)
	if err != nil {
		return nil, fmt.Errorf("failed to parse network range: %v", err)
	}

	s.logger.Infof("ARP scanning %d IP addresses", len(ips))

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

	s.logger.Infof("Starting %d workers for ARP scanning", workers)

	for i := 0; i < workers; i++ {
		wg.Add(1)
		go s.worker(ipChan, resultChan, &wg)
	}

	// Wait for all workers to complete
	go func() {
		wg.Wait()
		close(resultChan)
	}()

	// Collect results
	var devices []*models.Device
	for device := range resultChan {
		if device != nil {
			devices = append(devices, device)
			s.logger.Infof("Found ARP device: %s (%s) - %s", device.IP, device.MACAddress, device.Vendor)
		}
	}

	scanDuration := time.Since(start)
	s.logger.Infof("ARP scan completed in %v. Found %d devices", scanDuration, len(devices))

	return devices, nil
}

func (s *Scanner) worker(ipChan <-chan string, resultChan chan<- *models.Device, wg *sync.WaitGroup) {
	defer wg.Done()

	for ip := range ipChan {
		s.logger.Debugf("Scanning IP: %s", ip)

		device := s.scanSingleIP(ip)
		if device != nil {
			resultChan <- device
		}
	}
}

// scanSingleIP performs ARP scan for a single IP
func (s *Scanner) scanSingleIP(ip string) *models.Device {
	start := time.Now()
	s.logger.Debugf("Starting ARP scan for IP: %s", ip)

	// First try to ping the IP to see if it's reachable
	if !s.pingIP(ip) {
		s.logger.Debugf("IP %s is not reachable via ping", ip)
		return nil
	}
	s.logger.Debugf("IP %s responded to ping", ip)

	// Get MAC address using ARP
	macAddress, err := s.getARPEntry(ip)
	if err != nil {
		s.logger.Debugf("Failed to get ARP entry for %s: %v", ip, err)
		return nil
	}

	if macAddress == "" {
		s.logger.Debugf("No ARP entry found for %s", ip)
		return nil
	}

	s.logger.Debugf("Found MAC address for %s: %s", ip, macAddress)

	// Get vendor information
	vendor := s.getVendorFromMAC(macAddress)
	s.logger.Debugf("Vendor detection for %s (MAC: %s): %s", ip, macAddress, vendor)

	device := &models.Device{
		IP:           ip,
		MACAddress:   macAddress,
		LastSeen:     time.Now(),
		IsReachable:  true,
		ResponseTime: time.Since(start).Milliseconds(),
		Vendor:       vendor,
		ScanMethod:   "ARP",
	}

	s.logger.Debugf("ARP scan successful for %s: MAC=%s, Vendor=%s", ip, macAddress, vendor)
	return device
}

// pingIP checks if an IP is reachable via ping
func (s *Scanner) pingIP(ip string) bool {
	var cmd *exec.Cmd

	switch runtime.GOOS {
	case "windows":
		cmd = exec.Command("ping", "-n", "1", "-w", "500", ip)
	default: // Linux, macOS
		cmd = exec.Command("ping", "-c", "1", "-W", "1", ip)
	}

	err := cmd.Run()
	return err == nil
}

// getARPEntry retrieves MAC address from ARP table
func (s *Scanner) getARPEntry(ip string) (string, error) {
	var cmd *exec.Cmd

	switch runtime.GOOS {
	case "windows":
		cmd = exec.Command("arp", "-a", ip)
	default: // Linux, macOS
		cmd = exec.Command("arp", "-n", ip)
	}

	output, err := cmd.Output()
	if err != nil {
		return "", err
	}

	return s.parseARPOutput(string(output), ip)
}

// parseARPOutput parses ARP command output to extract MAC address
func (s *Scanner) parseARPOutput(output, targetIP string) (string, error) {
	lines := strings.Split(output, "\n")
	s.logger.Debugf("Parsing ARP output for %s:\n%s", targetIP, output)

	switch runtime.GOOS {
	case "windows":
		// Windows format: "  192.168.1.1           aa-bb-cc-dd-ee-ff     dynamic"
		macRegex := regexp.MustCompile(`([0-9a-fA-F]{2}[-:]){5}[0-9a-fA-F]{2}`)
		for _, line := range lines {
			if strings.Contains(line, targetIP) {
				s.logger.Debugf("Processing line: %s", line)
				matches := macRegex.FindStringSubmatch(line)
				if len(matches) > 0 {
					// Convert to standard format (XX:XX:XX:XX:XX:XX)
					mac := strings.ReplaceAll(matches[0], "-", ":")
					mac = strings.ToUpper(mac)
					s.logger.Debugf("Found MAC for %s: %s", targetIP, mac)
					return mac, nil
				}
			}
		}
	default:
		macRegex := regexp.MustCompile(`([0-9a-fA-F]{2}:){5}[0-9a-fA-F]{2}`)
		for _, line := range lines {
			if strings.Contains(line, targetIP) {
				s.logger.Debugf("Processing line: %s", line)
				matches := macRegex.FindStringSubmatch(line)
				if len(matches) > 0 {
					mac := strings.ToUpper(matches[0])
					s.logger.Debugf("Found MAC for %s: %s", targetIP, mac)
					return mac, nil
				}
			}
		}
	}

	s.logger.Debugf("No MAC address found in ARP output for %s", targetIP)
	return "", fmt.Errorf("MAC address not found for IP %s", targetIP)
}

// getVendorFromMAC attempts to identify vendor from MAC address OUI
func (s *Scanner) getVendorFromMAC(macAddress string) string {
	return s.vendorManager.GetVendor(macAddress)
}

// GetVendorDatabaseInfo returns information about the vendor database
func (s *Scanner) GetVendorDatabaseInfo() map[string]interface{} {
	return s.vendorManager.GetDatabaseInfo()
}

// ReloadVendorDatabase reloads the vendor database from file
func (s *Scanner) ReloadVendorDatabase() error {
	return s.vendorManager.ReloadDatabase()
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

	s.logger.Debugf("Generated %d IPs from range %s for ARP scan", len(ips), networkRange)
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
