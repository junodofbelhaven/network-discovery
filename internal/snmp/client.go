package snmp

import (
	"fmt"
	"strings"
	"time"

	"network-discovery/internal/models"

	"github.com/gosnmp/gosnmp"
	"github.com/sirupsen/logrus"
)

// Standard SNMP OIDs
const (
	OIDSysDescr    = "1.3.6.1.2.1.1.1.0" // System description
	OIDSysName     = "1.3.6.1.2.1.1.5.0" // System name
	OIDSysContact  = "1.3.6.1.2.1.1.4.0" // System contact
	OIDSysLocation = "1.3.6.1.2.1.1.6.0" // System location
	OIDSysUptime   = "1.3.6.1.2.1.1.3.0" // System uptime
)

type Client struct {
	timeout time.Duration
	retries int
	logger  *logrus.Logger
}

func NewClient(timeout time.Duration, retries int) *Client {
	logger := logrus.New()
	logger.SetLevel(logrus.InfoLevel)

	return &Client{
		timeout: timeout,
		retries: retries,
		logger:  logger,
	}
}

// QueryDevice queries a single device using SNMP
func (c *Client) QueryDevice(ip string, communities []string) (*models.Device, error) {
	device := &models.Device{
		IP:          ip,
		LastSeen:    time.Now(),
		IsReachable: false,
	}

	start := time.Now()

	for _, community := range communities {
		if err := c.queryWithCommunity(ip, community, device); err == nil {
			device.IsReachable = true
			device.Community = community
			device.ResponseTime = time.Since(start).Milliseconds()
			c.logger.Infof("Successfully queried device %s with community '%s'", ip, community)
			return device, nil
		}
	}

	device.ResponseTime = time.Since(start).Milliseconds()
	return device, fmt.Errorf("failed to query device %s with any community", ip)
}

func (c *Client) queryWithCommunity(ip, community string, device *models.Device) error {
	// Create SNMP client
	client := &gosnmp.GoSNMP{
		Target:    ip,
		Port:      161,
		Community: community,
		Version:   gosnmp.Version2c,
		Timeout:   c.timeout,
		Retries:   c.retries,
	}

	err := client.Connect()
	if err != nil {
		return fmt.Errorf("failed to connect to %s: %v", ip, err)
	}
	defer client.Conn.Close()

	// Query multiple OIDs at once
	oids := []string{
		OIDSysDescr,
		OIDSysName,
		OIDSysContact,
		OIDSysLocation,
		OIDSysUptime,
	}

	result, err := client.Get(oids)
	if err != nil {
		return fmt.Errorf("SNMP GET failed for %s: %v", ip, err)
	}

	// Parse results
	for _, variable := range result.Variables {
		switch variable.Name {
		case OIDSysDescr:
			device.Description = c.parseString(variable)
			c.parseVendorAndVersion(device)
		case OIDSysName:
			device.Hostname = c.parseString(variable)
		case OIDSysContact:
			device.Contact = c.parseString(variable)
		case OIDSysLocation:
			device.Location = c.parseString(variable)
		case OIDSysUptime:
			device.Uptime = c.parseUptime(variable)
		}
	}

	return nil
}

func (c *Client) parseString(variable gosnmp.SnmpPDU) string {
	switch variable.Type {
	case gosnmp.OctetString:
		return string(variable.Value.([]byte))
	default:
		return fmt.Sprintf("%v", variable.Value)
	}
}

func (c *Client) parseUptime(variable gosnmp.SnmpPDU) string {
	switch variable.Type {
	case gosnmp.TimeTicks:
		ticks := variable.Value.(uint32)
		seconds := ticks / 100

		days := seconds / 86400
		hours := (seconds % 86400) / 3600
		minutes := (seconds % 3600) / 60
		secs := seconds % 60

		return fmt.Sprintf("%dd %dh %dm %ds", days, hours, minutes, secs)
	default:
		return fmt.Sprintf("%v", variable.Value)
	}
}

func (c *Client) parseVendorAndVersion(device *models.Device) {
	desc := strings.ToLower(device.Description)

	// Common vendor detection patterns
	vendors := map[string]string{
		"cisco":     "Cisco",
		"juniper":   "Juniper",
		"huawei":    "Huawei",
		"hp":        "HP",
		"dell":      "Dell",
		"netgear":   "Netgear",
		"d-link":    "D-Link",
		"tp-link":   "TP-Link",
		"mikrotik":  "MikroTik",
		"ubiquiti":  "Ubiquiti",
		"fortinet":  "Fortinet",
		"palo alto": "Palo Alto",
	}

	for pattern, vendor := range vendors {
		if strings.Contains(desc, pattern) {
			device.Vendor = vendor
			break
		}
	}

	// Try to extract version information
	lines := strings.Split(device.Description, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.Contains(strings.ToLower(line), "version") {
			device.Version = line
			break
		}
	}
}

// IsDeviceReachable checks if a device responds to SNMP
func (c *Client) IsDeviceReachable(ip string, communities []string) bool {
	for _, community := range communities {
		client := &gosnmp.GoSNMP{
			Target:    ip,
			Port:      161,
			Community: community,
			Version:   gosnmp.Version2c,
			Timeout:   time.Second * 2, // Quick check
			Retries:   1,
		}

		err := client.Connect()
		if err != nil {
			continue
		}
		defer client.Conn.Close()

		_, err = client.Get([]string{OIDSysDescr})
		if err == nil {
			return true
		}
	}
	return false
}
