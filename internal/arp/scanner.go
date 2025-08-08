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
	logger     *logrus.Logger
	maxWorkers int
}

func NewScanner(maxWorkers int) *Scanner {
	logger := logrus.New()
	logger.SetLevel(logrus.InfoLevel)

	return &Scanner{
		logger:     logger,
		maxWorkers: maxWorkers,
	}
}

func NewScannerWithLogger(maxWorkers int, logger *logrus.Logger) *Scanner {
	return &Scanner{
		logger:     logger,
		maxWorkers: maxWorkers,
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
			s.logger.Infof("Found ARP device: %s (%s)", device.IP, device.MACAddress)
		}
	}

	scanDuration := time.Since(start)
	s.logger.Infof("ARP scan completed in %v. Found %d devices", scanDuration, len(devices))

	return devices, nil
}

func (s *Scanner) worker(ipChan <-chan string, resultChan chan<- *models.Device, wg *sync.WaitGroup) {
	defer wg.Done()

	for ip := range ipChan {
		s.logger.Debugf("ARP scanning IP: %s", ip)

		device := s.scanSingleIP(ip)
		if device != nil {
			resultChan <- device
		}
	}
}

// scanSingleIP performs ARP scan for a single IP
func (s *Scanner) scanSingleIP(ip string) *models.Device {
	start := time.Now()

	// First try to ping the IP to see if it's reachable
	if !s.pingIP(ip) {
		s.logger.Debugf("IP %s is not reachable via ping", ip)
		return nil
	}

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

	device := &models.Device{
		IP:           ip,
		MACAddress:   macAddress,
		LastSeen:     time.Now(),
		IsReachable:  true,
		ResponseTime: time.Since(start).Milliseconds(),
		Vendor:       s.getVendorFromMAC(macAddress),
		ScanMethod:   "ARP",
	}

	s.logger.Debugf("ARP scan successful for %s: MAC=%s", ip, macAddress)
	return device
}

// pingIP checks if an IP is reachable via ping
func (s *Scanner) pingIP(ip string) bool {
	var cmd *exec.Cmd

	switch runtime.GOOS {
	case "windows":
		cmd = exec.Command("ping", "-n", "1", "-w", "1000", ip)
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

	switch runtime.GOOS {
	case "windows":
		// Windows format: "  192.168.1.1           aa-bb-cc-dd-ee-ff     dynamic"
		macRegex := regexp.MustCompile(`([0-9a-fA-F]{2}[-:]){5}[0-9a-fA-F]{2}`)
		for _, line := range lines {
			if strings.Contains(line, targetIP) {
				matches := macRegex.FindStringSubmatch(line)
				if len(matches) > 0 {
					// Convert Windows format (aa-bb-cc-dd-ee-ff) to standard format
					mac := strings.ReplaceAll(matches[0], "-", ":")
					return strings.ToUpper(mac), nil
				}
			}
		}
	default:
		// Linux/macOS format: "192.168.1.1              ether   aa:bb:cc:dd:ee:ff   C                     eth0"
		macRegex := regexp.MustCompile(`([0-9a-fA-F]{2}:){5}[0-9a-fA-F]{2}`)
		for _, line := range lines {
			if strings.Contains(line, targetIP) {
				matches := macRegex.FindStringSubmatch(line)
				if len(matches) > 0 {
					return strings.ToUpper(matches[0]), nil
				}
			}
		}
	}

	return "", fmt.Errorf("MAC address not found for IP %s", targetIP)
}

// getVendorFromMAC attempts to identify vendor from MAC address OUI
func (s *Scanner) getVendorFromMAC(macAddress string) string {
	if len(macAddress) < 8 {
		return "Unknown"
	}

	// Extract OUI (first 3 octets)
	oui := strings.ReplaceAll(macAddress[:8], ":", "")
	oui = strings.ToUpper(oui)

	// Common OUI mappings (this could be expanded with a full OUI database)
	ouiVendors := map[string]string{
		"00:50:56": "VMware",
		"00:0C:29": "VMware",
		"00:05:69": "VMware",
		"08:00:27": "Oracle VirtualBox",
		"52:54:00": "QEMU/KVM",
		"00:16:3E": "Xen",
		"00:1B:21": "Intel",
		"00:13:72": "Dell",
		"00:14:22": "Dell",
		"B8:27:EB": "Raspberry Pi",
		"DC:A6:32": "Raspberry Pi",
		"E4:5F:01": "Raspberry Pi",
		"28:CD:C1": "Apple",
		"40:CB:C0": "Apple",
		"3C:07:54": "Apple",
		"00:1F:F3": "Apple",
		"D4:81:D7": "Apple",
		"A4:B1:C1": "Apple",
		"00:23:DF": "Apple",
		"F0:18:98": "Apple",
		"AC:DE:48": "Apple",
		"84:38:35": "Apple",
		"98:01:A7": "Apple",
		"6C:40:08": "Apple",
		"88:63:DF": "Apple",
		"78:4F:43": "Apple",
		"00:26:BB": "Apple",
		"00:3E:E1": "Apple",
		"04:0C:CE": "Apple",
		"8C:58:77": "Apple",
		"00:50:C2": "IEEE Registration Authority",
		"00:00:01": "Xerox",
		"AA:00:04": "DEC",
		"02:07:01": "Racal-Interlan",
	}

	// Check exact match first
	if vendor, exists := ouiVendors[macAddress[:8]]; exists {
		return vendor
	}

	// Check first 6 characters (3 octets without colons)
	if len(oui) >= 6 {
		ouiKey := oui[:2] + ":" + oui[2:4] + ":" + oui[4:6]
		if vendor, exists := ouiVendors[ouiKey]; exists {
			return vendor
		}
	}

	return "Unknown"
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
