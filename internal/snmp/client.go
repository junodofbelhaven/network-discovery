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

	// Interface OIDs for MAC address detection
	OIDIfPhysAddress = "1.3.6.1.2.1.2.2.1.6" // Interface physical address (MAC)
	OIDIfDescr       = "1.3.6.1.2.1.2.2.1.2" // Interface description
	OIDIfType        = "1.3.6.1.2.1.2.2.1.3" // Interface type
	OIDIfAdminStatus = "1.3.6.1.2.1.2.2.1.7" // Interface admin status
	OIDIfOperStatus  = "1.3.6.1.2.1.2.2.1.8" // Interface operational status
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

func NewClientWithLogger(timeout time.Duration, retries int, logger *logrus.Logger) *Client {
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

	c.logger.Debugf("Starting SNMP query for %s with communities: %v", ip, communities)

	for i, community := range communities {
		c.logger.Debugf("Trying community %d/%d: '%s' for %s", i+1, len(communities), community, ip)

		if err := c.queryWithCommunity(ip, community, device); err == nil {
			device.IsReachable = true
			device.Community = community
			device.ResponseTime = time.Since(start).Milliseconds()
			device.ScanMethod = "SNMP"

			c.logger.Infof("Successfully queried device %s with community '%s'", ip, community)
			c.logger.Debugf("Device details: hostname='%s', description='%s', vendor='%s', mac='%s'",
				device.Hostname, device.Description, device.Vendor, device.MACAddress)

			return device, nil
		} else {
			c.logger.Debugf("Failed with community '%s': %v", community, err)
		}
	}

	device.ResponseTime = time.Since(start).Milliseconds()
	c.logger.Debugf("Failed to query device %s with any community", ip)
	return device, fmt.Errorf("failed to query device %s with any community", ip)
}

func (c *Client) queryWithCommunity(ip, community string, device *models.Device) error {
	c.logger.Debugf("Attempting SNMP connection to %s with community '%s'", ip, community)

	// Create SNMP client
	client := &gosnmp.GoSNMP{
		Target:    ip,
		Port:      161,
		Community: community,
		Version:   gosnmp.Version2c,
		Timeout:   c.timeout,
		Retries:   c.retries,
	}

	c.logger.Debugf("SNMP client config: Target=%s, Port=161, Community=%s, Timeout=%v, Retries=%d",
		ip, community, c.timeout, c.retries)

	err := client.Connect()
	if err != nil {
		c.logger.Debugf("SNMP connection failed to %s: %v", ip, err)
		return fmt.Errorf("failed to connect to %s: %v", ip, err)
	}
	defer client.Conn.Close()

	c.logger.Debugf("SNMP connection established to %s", ip)

	// Query basic system OIDs
	oidQueries := map[string]string{
		OIDSysDescr:    "System Description",
		OIDSysName:     "System Name",
		OIDSysContact:  "System Contact",
		OIDSysLocation: "System Location",
		OIDSysUptime:   "System Uptime",
	}

	for oid, name := range oidQueries {
		c.logger.Debugf("Querying %s (%s)", name, oid)

		result, err := client.Get([]string{oid})
		if err != nil {
			c.logger.Debugf("Failed to query %s: %v", name, err)
			continue // Skip this OID but continue with others
		}

		if len(result.Variables) == 0 {
			c.logger.Debugf("No variables returned for %s", name)
			continue
		}

		variable := result.Variables[0]
		c.logger.Debugf("Variable Type for %s: %v", name, variable.Type)

		if variable.Type == gosnmp.NoSuchObject || variable.Type == gosnmp.NoSuchInstance {
			c.logger.Debugf("No such object/instance for %s", name)
			continue
		}

		// Parse the result
		switch oid {
		case OIDSysDescr:
			device.Description = c.parseString(variable)
			c.logger.Debugf("System Description parsed: '%s'", device.Description)
			// Safely parse vendor and version
			func() {
				defer func() {
					if r := recover(); r != nil {
						c.logger.Errorf("Error in parseVendorAndVersion: %v", r)
					}
				}()
				c.parseVendorAndVersion(device)
			}()
		case OIDSysName:
			device.Hostname = c.parseString(variable)
			c.logger.Debugf("System Name parsed: '%s'", device.Hostname)
		case OIDSysContact:
			device.Contact = c.parseString(variable)
			c.logger.Debugf("System Contact parsed: '%s'", device.Contact)
		case OIDSysLocation:
			device.Location = c.parseString(variable)
			c.logger.Debugf("System Location parsed: '%s'", device.Location)
		case OIDSysUptime:
			device.Uptime = c.parseUptime(variable)
			c.logger.Debugf("System Uptime parsed: '%s'", device.Uptime)
		}
	}

	// Try to get MAC address from interface table
	c.getMACAddress(client, device)

	// Check if we got at least some data
	if device.Description == "" && device.Hostname == "" && device.Contact == "" && device.Location == "" && device.Uptime == "" {
		return fmt.Errorf("no SNMP data retrieved from %s", ip)
	}

	c.logger.Debugf("Successfully retrieved SNMP data from %s", ip)
	return nil
}

// getMACAddress attempts to retrieve MAC address from SNMP interface table
func (c *Client) getMACAddress(client *gosnmp.GoSNMP, device *models.Device) {
	c.logger.Debugf("Attempting to get MAC address for device %s", device.IP)

	// Walk the interface physical address table
	err := client.Walk(OIDIfPhysAddress, func(pdu gosnmp.SnmpPDU) error {
		if pdu.Type == gosnmp.OctetString && pdu.Value != nil {
			if macBytes, ok := pdu.Value.([]byte); ok && len(macBytes) == 6 {
				// Convert bytes to MAC address string
				macAddr := fmt.Sprintf("%02X:%02X:%02X:%02X:%02X:%02X",
					macBytes[0], macBytes[1], macBytes[2],
					macBytes[3], macBytes[4], macBytes[5])

				// Skip empty MAC addresses (00:00:00:00:00:00)
				if macAddr != "00:00:00:00:00:00" {
					device.MACAddress = macAddr
					c.logger.Debugf("Found MAC address: %s", macAddr)
					return fmt.Errorf("stop_walk") // Stop walking after finding first valid MAC
				}
			}
		}
		return nil
	})

	if err != nil && err.Error() != "stop_walk" {
		c.logger.Debugf("Failed to walk interface table for MAC address: %v", err)
	}

	if device.MACAddress == "" {
		c.logger.Debugf("No MAC address found via SNMP for device %s", device.IP)
	}
}

func (c *Client) parseString(variable gosnmp.SnmpPDU) string {
	defer func() {
		if r := recover(); r != nil {
			c.logger.Errorf("Error parsing SNMP variable: %v", r)
		}
	}()

	switch variable.Type {
	case gosnmp.OctetString:
		if variable.Value != nil {
			if bytes, ok := variable.Value.([]byte); ok {
				return strings.TrimSpace(string(bytes))
			}
		}
		return ""
	case gosnmp.Integer:
		if variable.Value != nil {
			return fmt.Sprintf("%d", variable.Value)
		}
		return ""
	default:
		if variable.Value != nil {
			return strings.TrimSpace(fmt.Sprintf("%v", variable.Value))
		}
		return ""
	}
}

func (c *Client) parseUptime(variable gosnmp.SnmpPDU) string {
	defer func() {
		if r := recover(); r != nil {
			c.logger.Errorf("Error parsing uptime: %v", r)
		}
	}()

	switch variable.Type {
	case gosnmp.TimeTicks:
		if variable.Value != nil {
			if ticks, ok := variable.Value.(uint32); ok {
				seconds := ticks / 100

				days := seconds / 86400
				hours := (seconds % 86400) / 3600
				minutes := (seconds % 3600) / 60
				secs := seconds % 60

				return fmt.Sprintf("%dd %dh %dm %ds", days, hours, minutes, secs)
			}
		}
		return ""
	default:
		if variable.Value != nil {
			return fmt.Sprintf("%v", variable.Value)
		}
		return ""
	}
}

func (c *Client) parseVendorAndVersion(device *models.Device) {
	if device.Description == "" {
		c.logger.Debugf("No description to parse for vendor detection")
		return
	}

	desc := strings.ToLower(device.Description)
	c.logger.Debugf("Parsing description for vendor: %s", desc)

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
		"microsoft": "Microsoft",
		"windows":   "Microsoft",
		"linux":     "Linux",
		"ubuntu":    "Ubuntu",
		"centos":    "CentOS",
		"redhat":    "Red Hat",
	}

	for pattern, vendor := range vendors {
		if strings.Contains(desc, pattern) {
			device.Vendor = vendor
			c.logger.Debugf("Vendor detected: %s (pattern: %s)", vendor, pattern)
			break
		}
	}

	if device.Vendor == "" {
		device.Vendor = "Unknown"
		c.logger.Debugf("No vendor detected, set to Unknown")
	}

	// Try to extract version information
	lines := strings.Split(device.Description, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		lineLower := strings.ToLower(line)
		if strings.Contains(lineLower, "version") {
			device.Version = line
			c.logger.Debugf("Version found: %s", line)
			break
		}
	}

	if device.Version == "" {
		c.logger.Debugf("No version information found")
	}
}

// IsDeviceReachable checks if a device responds to SNMP
func (c *Client) IsDeviceReachable(ip string, communities []string) bool {
	c.logger.Debugf("Checking if device %s is reachable with communities: %v", ip, communities)

	for _, community := range communities {
		c.logger.Debugf("Testing reachability with community: %s", community)

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
			c.logger.Debugf("Connection failed for %s with %s: %v", ip, community, err)
			continue
		}

		c.logger.Debugf("Connection successful, testing SNMP query...")
		result, err := client.Get([]string{OIDSysDescr})
		client.Conn.Close()

		if err == nil && len(result.Variables) > 0 {
			variable := result.Variables[0]
			if variable.Type != gosnmp.NoSuchObject && variable.Type != gosnmp.NoSuchInstance {
				c.logger.Debugf("Device %s is reachable with community %s", ip, community)
				return true
			}
		}

		if err != nil {
			c.logger.Debugf("SNMP query failed for %s: %v", ip, err)
		}
	}

	c.logger.Debugf("Device %s is not reachable with any community", ip)
	return false
}
