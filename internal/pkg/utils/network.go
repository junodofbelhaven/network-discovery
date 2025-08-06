package utils

import (
	"fmt"
	"net"
	"regexp"
	"strings"
)

// IsValidIP checks if the given string is a valid IP address
func IsValidIP(ip string) bool {
	return net.ParseIP(ip) != nil
}

// IsValidCIDR checks if the given string is a valid CIDR notation
func IsValidCIDR(cidr string) bool {
	_, _, err := net.ParseCIDR(cidr)
	return err == nil
}

// GetLocalNetworks returns the local network ranges
func GetLocalNetworks() ([]string, error) {
	var networks []string

	interfaces, err := net.Interfaces()
	if err != nil {
		return nil, err
	}

	for _, iface := range interfaces {
		if iface.Flags&net.FlagUp == 0 || iface.Flags&net.FlagLoopback != 0 {
			continue
		}

		addrs, err := iface.Addrs()
		if err != nil {
			continue
		}

		for _, addr := range addrs {
			if ipnet, ok := addr.(*net.IPNet); ok {
				if ipnet.IP.To4() != nil {
					networks = append(networks, ipnet.String())
				}
			}
		}
	}

	return networks, nil
}

// ParseVendorFromDescription extracts vendor information from SNMP description
func ParseVendorFromDescription(description string) string {
	desc := strings.ToLower(description)

	vendorPatterns := map[string][]string{
		"Cisco":      {"cisco", "ios", "catalyst", "nexus"},
		"Juniper":    {"juniper", "junos", "srx", "mx", "ex"},
		"Huawei":     {"huawei", "vrp", "s5700", "s6700"},
		"HP":         {"hp", "hewlett", "packard", "procurve", "aruba"},
		"Dell":       {"dell", "powerconnect", "force10"},
		"Netgear":    {"netgear", "prosafe"},
		"D-Link":     {"d-link", "dgs", "des"},
		"TP-Link":    {"tp-link", "tl-", "archer"},
		"MikroTik":   {"mikrotik", "routeros", "routerboard"},
		"Ubiquiti":   {"ubiquiti", "unifi", "edgemax"},
		"Fortinet":   {"fortinet", "fortigate", "fortios"},
		"Palo Alto":  {"palo alto", "pa-", "panorama"},
		"SonicWall":  {"sonicwall", "sonicpoints"},
		"Watchguard": {"watchguard", "firebox"},
		"Meraki":     {"meraki", "cisco meraki"},
	}

	for vendor, patterns := range vendorPatterns {
		for _, pattern := range patterns {
			if strings.Contains(desc, pattern) {
				return vendor
			}
		}
	}

	return "Unknown"
}

// FormatUptime formats uptime ticks into human readable format
func FormatUptime(ticks uint32) string {
	seconds := ticks / 100

	days := seconds / 86400
	hours := (seconds % 86400) / 3600
	minutes := (seconds % 3600) / 60
	secs := seconds % 60

	if days > 0 {
		return fmt.Sprintf("%dd %dh %dm %ds", days, hours, minutes, secs)
	} else if hours > 0 {
		return fmt.Sprintf("%dh %dm %ds", hours, minutes, secs)
	} else if minutes > 0 {
		return fmt.Sprintf("%dm %ds", minutes, secs)
	} else {
		return fmt.Sprintf("%ds", secs)
	}
}

// SanitizeHostname cleans up hostname strings
func SanitizeHostname(hostname string) string {
	// Remove null bytes and control characters
	re := regexp.MustCompile(`[\x00-\x1f\x7f]`)
	cleaned := re.ReplaceAllString(hostname, "")

	// Trim whitespace
	cleaned = strings.TrimSpace(cleaned)

	return cleaned
}

// GetNetworkSize calculates the number of hosts in a CIDR range
func GetNetworkSize(cidr string) (int, error) {
	_, ipNet, err := net.ParseCIDR(cidr)
	if err != nil {
		return 0, err
	}

	ones, bits := ipNet.Mask.Size()
	hostBits := bits - ones

	// Subtract network and broadcast addresses
	size := (1 << hostBits) - 2
	if size < 0 {
		size = 0
	}

	return size, nil
}

// IsPrivateIP checks if an IP address is in a private range
func IsPrivateIP(ip string) bool {
	parsedIP := net.ParseIP(ip)
	if parsedIP == nil {
		return false
	}

	privateRanges := []string{
		"10.0.0.0/8",
		"172.16.0.0/12",
		"192.168.0.0/16",
		"127.0.0.0/8",
		"169.254.0.0/16",
	}

	for _, cidr := range privateRanges {
		_, network, _ := net.ParseCIDR(cidr)
		if network.Contains(parsedIP) {
			return true
		}
	}

	return false
}
