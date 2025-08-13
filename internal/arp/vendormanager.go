package arp

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
)

// VendorDatabase represents the structure of the OUI vendors JSON file
type VendorDatabase struct {
	OUIVendors map[string]string `json:"oui_vendors"`
	Metadata   VendorMetadata    `json:"metadata"`
	mu         sync.RWMutex      // For thread-safe access
	lastLoaded time.Time
}

// VendorMetadata contains metadata about the vendor database
type VendorMetadata struct {
	Version      string   `json:"version"`
	LastUpdated  string   `json:"last_updated"`
	Description  string   `json:"description"`
	TotalEntries int      `json:"total_entries"`
	Sources      []string `json:"sources"`
}

// VendorManager manages the vendor database
type VendorManager struct {
	database   *VendorDatabase
	configPath string
	logger     *logrus.Logger
}

// NewVendorManager creates a new vendor manager
func NewVendorManager(configPath string, logger *logrus.Logger) *VendorManager {
	if configPath == "" {
		configPath = "configs/oui_vendors.json"
	}

	vm := &VendorManager{
		configPath: configPath,
		logger:     logger,
		database:   &VendorDatabase{},
	}

	// Load the database
	if err := vm.LoadDatabase(); err != nil {
		logger.Errorf("Failed to load vendor database: %v", err)
		// Initialize with empty database as fallback
		vm.database = &VendorDatabase{
			OUIVendors: make(map[string]string),
			Metadata: VendorMetadata{
				Version:     "1.0.0",
				Description: "Fallback empty vendor database",
			},
		}
	}

	return vm
}

// LoadDatabase loads the vendor database from JSON file
func (vm *VendorManager) LoadDatabase() error {
	vm.logger.Infof("Loading vendor database from: %s", vm.configPath)

	// Check if file exists
	if _, err := os.Stat(vm.configPath); os.IsNotExist(err) {
		return fmt.Errorf("vendor database file not found: %s", vm.configPath)
	}

	// Read the JSON file
	data, err := os.ReadFile(vm.configPath)
	if err != nil {
		return fmt.Errorf("failed to read vendor database file: %v", err)
	}

	// Parse JSON
	var db VendorDatabase
	if err := json.Unmarshal(data, &db); err != nil {
		return fmt.Errorf("failed to parse vendor database JSON: %v", err)
	}

	// Update the database with thread safety
	vm.database.mu.Lock()
	vm.database.OUIVendors = db.OUIVendors
	vm.database.Metadata = db.Metadata
	vm.database.lastLoaded = time.Now()
	vm.database.mu.Unlock()

	vm.logger.Infof("Loaded vendor database: version %s, %d entries",
		db.Metadata.Version, len(db.OUIVendors))

	return nil
}

// ReloadDatabase reloads the vendor database from file
func (vm *VendorManager) ReloadDatabase() error {
	vm.logger.Info("Reloading vendor database...")
	return vm.LoadDatabase()
}

// GetVendor returns the vendor for a given MAC address
func (vm *VendorManager) GetVendor(macAddress string) string {
	if len(macAddress) < 8 {
		vm.logger.Debugf("MAC address too short for vendor detection: %s", macAddress)
		return "Unknown"
	}

	// Extract OUI (first 3 octets) and normalize
	oui := strings.ReplaceAll(macAddress[:8], ":", "")
	oui = strings.ToUpper(oui)

	vm.logger.Debugf("Processing MAC: %s -> OUI: %s", macAddress, oui)

	// Check if it's a locally administered address (2nd bit of first octet is 1)
	if len(oui) >= 2 {
		firstOctet := oui[:2]
		if val, err := strconv.ParseInt(firstOctet, 16, 64); err == nil {
			if val&0x02 != 0 { // Locally administered bit is set
				vm.logger.Debugf("MAC %s is locally administered (virtual)", macAddress)
				return "Virtual"
			}
		}
	}

	// Look up vendor in database
	vm.database.mu.RLock()
	vendor, exists := vm.database.OUIVendors[oui]
	vm.database.mu.RUnlock()

	if exists {
		vm.logger.Debugf("Vendor found for OUI %s: %s", oui, vendor)
		return vendor
	}

	vm.logger.Debugf("No vendor found for OUI: %s", oui)
	return "Unknown"
}

// AddVendor adds a new OUI-vendor mapping (runtime only, not persisted)
func (vm *VendorManager) AddVendor(oui, vendor string) {
	oui = strings.ToUpper(strings.ReplaceAll(oui, ":", ""))

	vm.database.mu.Lock()
	vm.database.OUIVendors[oui] = vendor
	vm.database.mu.Unlock()

	vm.logger.Infof("Added vendor mapping: %s -> %s", oui, vendor)
}

// GetDatabaseInfo returns information about the loaded database
func (vm *VendorManager) GetDatabaseInfo() map[string]interface{} {
	vm.database.mu.RLock()
	defer vm.database.mu.RUnlock()

	return map[string]interface{}{
		"version":       vm.database.Metadata.Version,
		"last_updated":  vm.database.Metadata.LastUpdated,
		"description":   vm.database.Metadata.Description,
		"total_entries": len(vm.database.OUIVendors),
		"sources":       vm.database.Metadata.Sources,
		"config_path":   vm.configPath,
		"last_loaded":   vm.database.lastLoaded.Format(time.RFC3339),
	}
}

// GetAllVendors returns all OUI-vendor mappings (for debugging)
func (vm *VendorManager) GetAllVendors() map[string]string {
	vm.database.mu.RLock()
	defer vm.database.mu.RUnlock()

	// Return a copy to prevent external modification
	vendors := make(map[string]string)
	for oui, vendor := range vm.database.OUIVendors {
		vendors[oui] = vendor
	}
	return vendors
}

// SaveDatabase saves the current database to JSON file
func (vm *VendorManager) SaveDatabase() error {
	vm.database.mu.RLock()
	defer vm.database.mu.RUnlock()

	// Update metadata
	vm.database.Metadata.TotalEntries = len(vm.database.OUIVendors)
	vm.database.Metadata.LastUpdated = time.Now().Format("2006-01-02")

	// Marshal to JSON with indentation
	data, err := json.MarshalIndent(vm.database, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal vendor database: %v", err)
	}

	// Ensure directory exists
	dir := filepath.Dir(vm.configPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create config directory: %v", err)
	}

	// Write to file
	if err := os.WriteFile(vm.configPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write vendor database file: %v", err)
	}

	vm.logger.Infof("Saved vendor database to: %s", vm.configPath)
	return nil
}
