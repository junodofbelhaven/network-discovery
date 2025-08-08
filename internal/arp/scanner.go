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
		cmd = exec.Command("ping", "-n", "1", "-w", "500", ip) // Reduced from 1000ms to 500ms
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

	s.logger.Debugf("Extracting vendor for MAC: %s, OUI: %s", macAddress, oui)

	// Common OUI mappings (expanded and corrected)
	ouiVendors := map[string]string{
		// VMware
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

		// Xen
		"00163E": "Xen",

		// Intel
		"001B21": "Intel",
		"7085C2": "Intel",
		"8086F2": "Intel",
		"A45E60": "Intel",
		"009027": "Intel",
		"0015F2": "Intel",

		// Dell
		"001372": "Dell",
		"001422": "Dell",
		"0018F3": "Dell",
		"002564": "Dell",
		"00B0D0": "Dell",
		"001E4F": "Dell",
		"F8B156": "Dell",
		"3417EB": "Dell",

		// Raspberry Pi
		"B827EB": "Raspberry Pi",
		"DCA632": "Raspberry Pi",
		"E45F01": "Raspberry Pi",
		"DC2632": "Raspberry Pi",
		"28CD4C": "Raspberry Pi",
		"2C3AE8": "Raspberry Pi",

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
		"5855CA": "Apple",
		"BC92B6": "Apple",
		"48A195": "Apple",

		// Cisco
		"0050F2": "Cisco",
		"000142": "Cisco",
		"0060B0": "Cisco",
		"000A41": "Cisco",
		"001BD4": "Cisco",
		"001C58": "Cisco",
		"002155": "Cisco",
		"5087B5": "Cisco",

		// HP
		"001F29": "HP",
		"002608": "HP",
		"002655": "HP",
		"70106F": "HP",
		"009C02": "HP",
		"001E0B": "HP",
		"001CC4": "HP",

		// Netgear
		"001E2A": "Netgear",
		"0026F2": "Netgear",
		"002713": "Netgear",
		"84D47E": "Netgear",
		"A021B7": "Netgear",

		// D-Link
		"001195": "D-Link",
		"001346": "D-Link",
		"0015E9": "D-Link",
		"001CF0": "D-Link",
		"14D64D": "D-Link",

		// TP-Link
		"141877": "TP-Link",
		"50C7BF": "TP-Link",
		"C04A00": "TP-Link",
		"E894F6": "TP-Link",
		"F4F26D": "TP-Link",
		"FC7516": "TP-Link",

		// MikroTik
		"E748B7": "MikroTik",
		"64D154": "MikroTik",
		"48A9C2": "MikroTik",
		"CC2DE0": "MikroTik",
		"6C3B6B": "MikroTik",

		// Ubiquiti
		"043E37": "Ubiquiti",
		"68D79A": "Ubiquiti",
		"80EA96": "Ubiquiti",
		"F09FC2": "Ubiquiti",
		"78A050": "Ubiquiti",
		"F4E2C5": "Ubiquiti",

		// Samsung
		"001D25":  "Samsung",
		"002454":  "Samsung",
		"0025E5":  "Samsung",
		"34E6AD":  "Samsung",
		"8C77120": "Samsung",
		"C85195":  "Samsung",

		// Huawei
		"001E10": "Huawei",
		"002EC7": "Huawei",
		"0025CE": "Huawei",
		"34E318": "Huawei",
		"4C549D": "Huawei",
		"6C92BF": "Huawei",
		"ACE215": "Huawei",

		// Realtek
		"52540A": "Realtek",
		"001DD8": "Realtek",
		"0019E0": "Realtek",
		"E0469A": "Realtek",

		// Microsoft
		"B499BA": "Microsoft",

		// Asus
		"001E8C":  "Asus",
		"0026184": "Asus",
		"2C56DC":  "Asus",
		"04925A":  "Asus",
		"C860008": "Asus",

		// Qualcomm
		"001A8A": "Qualcomm",
		"002719": "Qualcomm",
		"8CFDB0": "Qualcomm",

		// Broadcom
		"001018": "Broadcom",
		"002067": "Broadcom",
		"ACF1DF": "Broadcom",
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
