package models

import "time"

// Device represents a network device discovered via SNMP
type Device struct {
	IP           string    `json:"ip"`
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
}

// NetworkTopology represents the overall network topology
type NetworkTopology struct {
	Devices        []Device  `json:"devices"`
	TotalCount     int       `json:"total_count"`
	ReachableCount int       `json:"reachable_count"`
	ScanTime       time.Time `json:"scan_time"`
	ScanDuration   int64     `json:"scan_duration_ms"`
}

// ScanRequest represents a network scan request
type ScanRequest struct {
	NetworkRange string   `json:"network_range" binding:"required"` // e.g., "192.168.1.0/24"
	Communities  []string `json:"communities"`                      // SNMP communities to try
	Timeout      int      `json:"timeout"`                          // Timeout in seconds
	Retries      int      `json:"retries"`                          // Number of retries
}
