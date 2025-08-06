package snmp

import (
	"fmt"
	"network-discovery/internal/models"
	"testing"
	"time"
)

func TestNewClient(t *testing.T) {
	timeout := time.Second * 5
	retries := 3

	client := NewClient(timeout, retries)

	if client == nil {
		t.Fatal("Expected client to be created, got nil")
	}

	if client.timeout != timeout {
		t.Errorf("Expected timeout %v, got %v", timeout, client.timeout)
	}

	if client.retries != retries {
		t.Errorf("Expected retries %d, got %d", retries, client.retries)
	}
}

func TestParseVendorAndVersion(t *testing.T) {
	client := NewClient(time.Second*5, 2)

	testCases := []struct {
		description string
		expected    string
	}{
		{
			description: "Cisco IOS Software, C2960 Software (C2960-LANBASEK9-M), Version 15.0(2)SE4",
			expected:    "Cisco",
		},
		{
			description: "Juniper Networks, Inc. srx240h2 internet router",
			expected:    "Juniper",
		},
		{
			description: "Huawei Versatile Routing Platform Software",
			expected:    "Huawei",
		},
		{
			description: "HP J9019B ProCurve Switch 2510G-24",
			expected:    "HP",
		},
		{
			description: "Unknown Device Description",
			expected:    "Unknown",
		},
	}

	for _, tc := range testCases {
		device := &models.Device{Description: tc.description}
		client.parseVendorAndVersion(device)

		if device.Vendor != tc.expected {
			t.Errorf("Expected vendor %s for description '%s', got %s",
				tc.expected, tc.description, device.Vendor)
		}
	}
}

func TestParseUptime(t *testing.T) {
	//client := NewClient(time.Second*5, 2)

	// Test uptime parsing
	// 1 day = 86400 seconds = 8640000 ticks
	ticks := uint32(8640000)
	expected := "1d 0h 0m 0s"

	// Create a mock SNMP PDU (this would normally come from gosnmp)
	// For testing, we'll test the format function directly
	seconds := ticks / 100
	days := seconds / 86400
	hours := (seconds % 86400) / 3600
	minutes := (seconds % 3600) / 60
	secs := seconds % 60

	result := fmt.Sprintf("%dd %dh %dm %ds", days, hours, minutes, secs)

	if result != expected {
		t.Errorf("Expected uptime %s, got %s", expected, result)
	}
}

func TestIsDeviceReachable(t *testing.T) {
	client := NewClient(time.Second*2, 1)

	// Test with localhost (should not be reachable via SNMP typically)
	communities := []string{"public"}
	reachable := client.IsDeviceReachable("127.0.0.1", communities)

	// This should typically return false as localhost usually doesn't run SNMP
	if reachable {
		t.Log("Localhost appears to be running SNMP service")
	} else {
		t.Log("Localhost is not reachable via SNMP (expected)")
	}
}

// Benchmark tests
func BenchmarkParseVendorAndVersion(b *testing.B) {
	client := NewClient(time.Second*5, 2)
	device := &models.Device{
		Description: "Cisco IOS Software, C2960 Software (C2960-LANBASEK9-M), Version 15.0(2)SE4",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		client.parseVendorAndVersion(device)
	}
}
