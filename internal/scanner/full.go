// internal/scanner/full.go
package scanner

import (
	"fmt"
	"sync"
	"time"

	"network-discovery/internal/arp"
	"network-discovery/internal/models"
	"network-discovery/internal/snmp"

	"github.com/sirupsen/logrus"
)

type FullScanner struct {
	snmpScanner *snmp.Scanner
	arpScanner  *arp.Scanner
	logger      *logrus.Logger
}

func NewFullScanner(snmpClient *snmp.Client, maxWorkers int) *FullScanner {
	logger := logrus.New()
	logger.SetLevel(logrus.InfoLevel)

	return &FullScanner{
		snmpScanner: snmp.NewScanner(snmpClient, maxWorkers),
		arpScanner:  arp.NewScanner(maxWorkers),
		logger:      logger,
	}
}

func NewFullScannerWithLogger(snmpClient *snmp.Client, maxWorkers int, logger *logrus.Logger) *FullScanner {
	return &FullScanner{
		snmpScanner: snmp.NewScannerWithLogger(snmpClient, maxWorkers, logger),
		arpScanner:  arp.NewScannerWithLogger(maxWorkers, logger),
		logger:      logger,
	}
}

// PerformFullScan performs both SNMP and ARP scans and merges the results
func (fs *FullScanner) PerformFullScan(networkRange string, communities []string) (*models.NetworkTopology, error) {
	start := time.Now()
	fs.logger.Infof("Starting full scan (SNMP + ARP) for range: %s", networkRange)

	// Channels for concurrent scanning
	snmpChan := make(chan []*models.Device, 1)
	arpChan := make(chan []*models.Device, 1)
	errorChan := make(chan error, 2)

	var wg sync.WaitGroup

	// Start SNMP scan in goroutine
	wg.Add(1)
	go func() {
		defer wg.Done()
		fs.logger.Info("Starting SNMP scan...")

		topology, err := fs.snmpScanner.ScanNetwork(networkRange, communities)
		if err != nil {
			fs.logger.Errorf("SNMP scan failed: %v", err)
			errorChan <- fmt.Errorf("SNMP scan failed: %v", err)
			snmpChan <- []*models.Device{}
			return
		}

		// Convert topology devices to device pointers
		var devices []*models.Device
		for i := range topology.Devices {
			topology.Devices[i].ScanMethod = "SNMP"
			devices = append(devices, &topology.Devices[i])
		}

		fs.logger.Infof("SNMP scan completed: found %d devices", len(devices))
		snmpChan <- devices
	}()

	// Start ARP scan in goroutine
	wg.Add(1)
	go func() {
		defer wg.Done()
		fs.logger.Info("Starting ARP scan...")

		devices, err := fs.arpScanner.ScanNetwork(networkRange)
		if err != nil {
			fs.logger.Errorf("ARP scan failed: %v", err)
			errorChan <- fmt.Errorf("ARP scan failed: %v", err)
			arpChan <- []*models.Device{}
			return
		}

		fs.logger.Infof("ARP scan completed: found %d devices", len(devices))
		arpChan <- devices
	}()

	// Wait for both scans to complete
	go func() {
		wg.Wait()
		close(snmpChan)
		close(arpChan)
		close(errorChan)
	}()

	// Collect results
	var snmpDevices []*models.Device
	var arpDevices []*models.Device
	var scanErrors []error

	// Wait for results
	for snmpChan != nil || arpChan != nil {
		select {
		case devices, ok := <-snmpChan:
			if !ok {
				snmpChan = nil
			} else {
				snmpDevices = devices
			}
		case devices, ok := <-arpChan:
			if !ok {
				arpChan = nil
			} else {
				arpDevices = devices
			}
		case err := <-errorChan:
			if err != nil {
				scanErrors = append(scanErrors, err)
			}
		}
	}

	// Merge results
	mergedDevices := fs.mergeDevices(snmpDevices, arpDevices)

	scanDuration := time.Since(start)

	// Count different types of devices
	snmpCount := 0
	arpOnlyCount := 0
	for _, device := range mergedDevices {
		if device.ScanMethod == "SNMP" || device.ScanMethod == "COMBINED" {
			snmpCount++
		}
		if device.ScanMethod == "ARP" {
			arpOnlyCount++
		}
	}

	topology := &models.NetworkTopology{
		Devices:        mergedDevices,
		TotalCount:     len(mergedDevices),
		ReachableCount: len(mergedDevices), // All devices in full scan are reachable
		SNMPCount:      snmpCount,
		ARPCount:       arpOnlyCount,
		ScanTime:       start,
		ScanDuration:   scanDuration.Milliseconds(),
		ScanMethod:     "FULL",
	}

	fs.logger.Infof("Full scan completed in %v. Found %d total devices (%d SNMP, %d ARP-only)",
		scanDuration, len(mergedDevices), snmpCount, arpOnlyCount)

	// If there were errors but we still got some results, log warnings
	if len(scanErrors) > 0 {
		for _, err := range scanErrors {
			fs.logger.Warnf("Scan error: %v", err)
		}
	}

	return topology, nil
}

// mergeDevices merges SNMP and ARP scan results, combining devices found by both methods
func (fs *FullScanner) mergeDevices(snmpDevices, arpDevices []*models.Device) []models.Device {
	deviceMap := make(map[string]*models.Device)

	// Add SNMP devices first
	for _, device := range snmpDevices {
		deviceMap[device.IP] = device
	}

	// Process ARP devices
	for _, arpDevice := range arpDevices {
		if existingDevice, exists := deviceMap[arpDevice.IP]; exists {
			// Device found in both SNMP and ARP scans - merge information
			fs.logger.Debugf("Merging SNMP and ARP data for device: %s", arpDevice.IP)

			// Add MAC address from ARP if not already present
			if existingDevice.MACAddress == "" {
				existingDevice.MACAddress = arpDevice.MACAddress
			}

			// Update vendor if not detected via SNMP but detected via MAC
			if (existingDevice.Vendor == "" || existingDevice.Vendor == "Unknown") &&
				(arpDevice.Vendor != "" && arpDevice.Vendor != "Unknown") {
				existingDevice.Vendor = arpDevice.Vendor
			}

			// Update scan method to indicate combined data
			existingDevice.ScanMethod = "COMBINED"

		} else {
			// Device only found via ARP
			fs.logger.Debugf("Adding ARP-only device: %s", arpDevice.IP)
			deviceMap[arpDevice.IP] = arpDevice
		}
	}

	// Try to get MAC addresses for SNMP devices that don't have them
	fs.enhanceSNMPDevicesWithMAC(deviceMap)

	// Convert map to slice
	var result []models.Device
	for _, device := range deviceMap {
		result = append(result, *device)
	}

	return result
}

// enhanceSNMPDevicesWithMAC attempts to get MAC addresses for SNMP devices
func (fs *FullScanner) enhanceSNMPDevicesWithMAC(deviceMap map[string]*models.Device) {
	for ip, device := range deviceMap {
		if device.ScanMethod == "SNMP" && device.MACAddress == "" {
			fs.logger.Debugf("Attempting to get MAC address for SNMP device: %s", ip)

			// Try to get MAC via ARP for this specific IP
			if macAddr := fs.getMACForIP(ip); macAddr != "" {
				device.MACAddress = macAddr
				device.ScanMethod = "COMBINED"
				fs.logger.Debugf("Added MAC address %s to SNMP device %s", macAddr, ip)
			}
		}
	}
}

// getMACForIP attempts to get MAC address for a specific IP using ARP
func (fs *FullScanner) getMACForIP(ip string) string {
	// Create a temporary ARP scanner for single IP lookup
	tempScanner := arp.NewScannerWithLogger(1, fs.logger)

	// This is a simplified approach - in a full implementation,
	// you might want to use a more direct ARP lookup method
	devices, err := tempScanner.ScanNetwork(ip + "/32")
	if err != nil || len(devices) == 0 {
		return ""
	}

	return devices[0].MACAddress
}

// PerformSNMPScan performs only SNMP scan
func (fs *FullScanner) PerformSNMPScan(networkRange string, communities []string) (*models.NetworkTopology, error) {
	fs.logger.Infof("Starting SNMP-only scan for range: %s", networkRange)

	topology, err := fs.snmpScanner.ScanNetwork(networkRange, communities)
	if err != nil {
		return nil, err
	}

	// Update scan method for all devices
	for i := range topology.Devices {
		topology.Devices[i].ScanMethod = "SNMP"
	}

	topology.ScanMethod = "SNMP"
	topology.SNMPCount = topology.ReachableCount
	topology.ARPCount = 0

	return topology, nil
}

// PerformARPScan performs only ARP scan
func (fs *FullScanner) PerformARPScan(networkRange string) (*models.NetworkTopology, error) {
	start := time.Now()
	fs.logger.Infof("Starting ARP-only scan for range: %s", networkRange)

	devices, err := fs.arpScanner.ScanNetwork(networkRange)
	if err != nil {
		return nil, err
	}

	// Convert to regular devices slice
	var deviceSlice []models.Device
	for _, device := range devices {
		deviceSlice = append(deviceSlice, *device)
	}

	scanDuration := time.Since(start)

	topology := &models.NetworkTopology{
		Devices:        deviceSlice,
		TotalCount:     len(deviceSlice),
		ReachableCount: len(deviceSlice),
		SNMPCount:      0,
		ARPCount:       len(deviceSlice),
		ScanTime:       start,
		ScanDuration:   scanDuration.Milliseconds(),
		ScanMethod:     "ARP",
	}

	return topology, nil
}
