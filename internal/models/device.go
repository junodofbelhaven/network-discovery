package models

import "time"

// Device represents a network device discovered via SNMP or ARP
type Device struct {
	IP           string    `json:"ip"`
	MACAddress   string    `json:"mac_address,omitempty"` // MAC address from ARP or SNMP
	Hostname     string    `json:"hostname"`
	Description  string    `json:"description"`
	Contact      string    `json:"contact"`
	Location     string    `json:"location"`
	Uptime       string    `json:"uptime"`
	Vendor       string    `json:"vendor"`
	Model        string    `json:"model"`
	Version      string    `json:"version"`
	Community    string    `json:"-"` // SNMP community string (hidden from JSON)
	LastSeen     time.Time `json:"last_seen"`
	IsReachable  bool      `json:"is_reachable"`
	ResponseTime int64     `json:"response_time_ms"`
	ScanMethod   string    `json:"scan_method"` // "SNMP", "ARP", or "COMBINED"
}

// NetworkTopology represents the overall network topology
type NetworkTopology struct {
	Devices        []Device  `json:"devices"`
	TotalCount     int       `json:"total_count"`
	ReachableCount int       `json:"reachable_count"`
	SNMPCount      int       `json:"snmp_count"` // Number of SNMP-enabled devices
	ARPCount       int       `json:"arp_count"`  // Number of ARP-only devices
	ScanTime       time.Time `json:"scan_time"`
	ScanDuration   int64     `json:"scan_duration_ms"`
	ScanMethod     string    `json:"scan_method"` // "SNMP", "ARP", or "FULL"
}

// ScanRequest represents a network scan request
type ScanRequest struct {
	NetworkRange string   `json:"network_range" binding:"required"` // e.g., "192.168.1.0/24"
	Communities  []string `json:"communities"`                      // SNMP communities to try
	Timeout      int      `json:"timeout"`                          // Timeout in seconds
	Retries      int      `json:"retries"`                          // Number of retries
	ScanType     string   `json:"scan_type"`                        // "snmp", "arp", or "full"
}

// FullScanResult represents the result of a full scan (SNMP + ARP)
type FullScanResult struct {
	Topology   *NetworkTopology       `json:"topology"`
	Statistics map[string]interface{} `json:"statistics"`
	ScanInfo   ScanInfo               `json:"scan_info"`
}

// ScanInfo provides detailed information about the scan
type ScanInfo struct {
	ScanType        string   `json:"scan_type"`
	NetworkRange    string   `json:"network_range"`
	SNMPCommunities []string `json:"snmp_communities,omitempty"`
	Timeout         int      `json:"timeout"`
	Retries         int      `json:"retries"`
	WorkerCount     int      `json:"worker_count"`
}
