package arp

import (
	"fmt"
	"net"
	"os/exec"
	"regexp"
	"runtime"
	"strconv"
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
		// Linux/macOS format: "192.168.1.1              ether   aa:bb:cc:dd:ee:ff   C                     eth0"
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
	if len(macAddress) < 8 {
		s.logger.Debugf("MAC address too short for vendor detection: %s", macAddress)
		return "Unknown"
	}

	// Extract OUI (first 3 octets) and normalize
	oui := strings.ReplaceAll(macAddress[:8], ":", "")
	oui = strings.ToUpper(oui)

	s.logger.Debugf("Processing MAC: %s -> OUI: %s", macAddress, oui)

	// Check if it's a locally administered address (2nd bit of first octet is 1)
	if len(oui) >= 2 {
		firstOctet := oui[:2]
		if val, err := strconv.ParseInt(firstOctet, 16, 64); err == nil {
			if val&0x02 != 0 { // Locally administered bit is set
				s.logger.Debugf("MAC %s is locally administered (virtual/random)", macAddress)
				return "Virtual/Random"
			}
		}
	}

	// Extensive OUI mappings database - real vendor assignments
	ouiVendors := map[string]string{
		// Specific devices from your scan
		"000121": "Cabletron Systems",
		"C403A8": "Shenzhen Coship Electronics",
		"5C3A45": "Liteon Technology Corporation",
		"600308": "Apple",

		// VMware (for testing)
		"005056": "VMware",
		"000C29": "VMware",
		"001C14": "VMware",
		"005069": "VMware",

		// VirtualBox
		"080027": "Oracle VirtualBox",
		"0A0027": "Oracle VirtualBox",

		// QEMU/KVM
		"525400": "QEMU/KVM",
		"021C42": "QEMU/KVM",

		// Raspberry Pi
		"B827EB": "Raspberry Pi Foundation",
		"DCA632": "Raspberry Pi Trading",
		"E45F01": "Raspberry Pi Foundation",
		"DC2632": "Raspberry Pi Trading",
		"28CD4C": "Raspberry Pi Trading",
		"2C3AE8": "Raspberry Pi Trading",

		// Apple
		"28CDC1": "Apple",
		"40CBC0": "Apple",
		"3C0754": "Apple",
		"001FF3": "Apple",
		"D481D7": "Apple",
		"A4B1C1": "Apple",
		"0023DF": "Apple",
		"F01898": "Apple",
		"ACDE48": "Apple",
		"843835": "Apple",
		"9801A7": "Apple",
		"6C4008": "Apple",
		"8863DF": "Apple",
		"784F43": "Apple",
		"0026BB": "Apple",
		"003EE1": "Apple",
		"040CCE": "Apple",
		"8C5877": "Apple",

		// Intel
		"001B21": "Intel Corporation",
		"7085C2": "Intel Corporation",
		"8086F2": "Intel Corporation",
		"A45E60": "Intel Corporation",
		"009027": "Intel Corporation",
		"0015F2": "Intel Corporation",

		// Cisco
		"0050F2": "Cisco Systems",
		"000142": "Cisco Systems",
		"0060B0": "Cisco Systems",
		"000A41": "Cisco Systems",
		"001BD4": "Cisco Systems",
		"001C58": "Cisco Systems",
		"002155": "Cisco Systems",

		// Dell
		"001372": "Dell Inc.",
		"001422": "Dell Inc.",
		"0018F3": "Dell Inc.",
		"002564": "Dell Inc.",
		"00B0D0": "Dell Inc.",
		"001E4F": "Dell Inc.",

		// HP
		"001F29": "Hewlett Packard",
		"002608": "Hewlett Packard",
		"002655": "Hewlett Packard",
		"70106F": "Hewlett Packard",

		// Samsung
		"001D25": "Samsung Electronics",
		"002454": "Samsung Electronics",
		"0025E5": "Samsung Electronics",
		"34E6AD": "Samsung Electronics",

		// Huawei
		"001E10": "Huawei Technologies",
		"002EC7": "Huawei Technologies",
		"0025CE": "Huawei Technologies",
		"34E318": "Huawei Technologies",

		// TP-Link
		"141877": "TP-Link Technologies",
		"50C7BF": "TP-Link Technologies",
		"C04A00": "TP-Link Technologies",
		"E894F6": "TP-Link Technologies",

		// Netgear
		"001E2A": "Netgear",
		"0026F2": "Netgear",
		"002713": "Netgear",
		"84D47E": "Netgear",

		// D-Link
		"001195": "D-Link Corporation",
		"001346": "D-Link Corporation",
		"0015E9": "D-Link Corporation",
		"001CF0": "D-Link Corporation",

		// Xiaomi
		"8CFDB0": "Xiaomi Communications",
		"64B473": "Xiaomi Communications",
		"F8A45F": "Xiaomi Communications",

		// Microsoft
		"7CD1C3": "Microsoft Corporation",

		// Other common manufacturers
		"001B63": "Hon Hai Precision Ind.",
		"E0469A": "Realtek Semiconductor",
		"001DD8": "Realtek Semiconductor",
	}

	// Check exact OUI match (6 hex digits)
	if vendor, exists := ouiVendors[oui]; exists {
		s.logger.Debugf("Vendor found for OUI %s: %s", oui, vendor)
		return vendor
	}

	s.logger.Debugf("No vendor found for OUI: %s", oui)
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
