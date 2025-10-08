# 🌐 Network Discovery Tool

A modern network discovery tool written in Go that analyzes network topology using SNMP and ARP protocols. Automatically discovers all devices on your network and collects detailed information.

## ✨ Features

- 🔍 **Full Network Scan**: Comprehensive network discovery with SNMP + ARP combination
- 📡 **SNMP v2c Support**: SNMP discovery with detailed device information
- 🌐 **ARP Scanning**: Discovery of all IP-enabled devices
- ⚡ **High Performance**: Fast scanning with 50 concurrent workers
- 🏷️ **Vendor Detection**: Vendor recognition with JSON-based OUI database
- 📱 **MAC Address Resolution**: Hardware address identification
- 🔍 **Port Scanning**: Detection and display of open ports
- ⏱️ **Response Time Measurement**: Measures network latency for each device
- 🌐 **REST API**: Easy integration with RESTful web services
- 💻 **Web Interface**: User-friendly web-based control panel
- 📊 **Detailed Reporting**: Network statistics and device inventory

## 🚀 Quick Start

### Prerequisites

- Go 1.23 or higher
- Git
- nmap

### Installation

```bash
# Clone the project
git clone https://github.com/junodofbelhaven/network-discovery.git
cd network-discovery

# Install dependencies
go mod tidy

# Run the application
go run cmd/main.go
```

## Web GUI

## <img width="1755" height="1640" alt="image" src="https://github.com/user-attachments/assets/d677405d-c5a7-4123-ab5a-2fe279879ac8" />



## <img width="789" height="688" alt="image" src="https://github.com/user-attachments/assets/690256c3-496a-4bce-99b6-ce28d9127c3f" />


## 📖 API Documentation

### Main Endpoints

| Method | Endpoint                         | Description                |
| ------ | -------------------------------- | -------------------------- |
| GET    | `/api/v1/health`                 | Service health check       |
| GET    | `/api/v1/version`                | Version information        |
| GET    | `/api/v1/scan-methods`           | Scan methods information   |
| POST   | `/api/v1/network/full-scan`      | Full scan (SNMP + ARP)     |
| POST   | `/api/v1/network/scan/snmp`      | SNMP scan only             |
| POST   | `/api/v1/network/scan/arp`       | ARP scan only              |
| POST   | `/api/v1/network/scan/full`      | Full scan (alternative)    |
| POST   | `/api/v1/network/scan`           | Legacy SNMP scan           |
| GET    | `/api/v1/network/quick-scan`     | Quick device discovery     |
| GET    | `/api/v1/network/validate`       | Network range validation   |
| GET    | `/api/v1/device/{ip}`            | Single device scan         |
| GET    | `/api/v1/vendor-database`        | Vendor database info       |
| POST   | `/api/v1/vendor-database/reload` | Reload vendor database     |

### Full Network Scan (Main Endpoint)

**POST** `/api/v1/network/full-scan`

```json
{
  "network_range": "192.168.1.0/24",
  "communities": ["public", "private"],
  "timeout": 2,
  "retries": 1,
  "scan_type": "full"
}
```

**Response:**

```json
{
  "topology": {
    "devices": [
      {
        "ip": "192.168.1.1",
        "mac_address": "AA:BB:CC:DD:EE:FF",
        "hostname": "router.local",
        "vendor": "Cisco",
        "description": "Cisco IOS Software...",
        "uptime": "45d 12h 30m 15s",
        "is_reachable": true,
        "response_time_ms": 23,
        "scan_method": "COMBINED"
      }
    ],
    "total_count": 5,
    "reachable_count": 5,
    "snmp_count": 3,
    "arp_count": 2,
    "scan_duration_ms": 15420,
    "scan_method": "FULL"
  },
  "statistics": {
    "total_devices": 5,
    "reachable_devices": 5,
    "snmp_devices": 3,
    "arp_only_devices": 2,
    "devices_with_mac": 5,
    "vendor_distribution": {
      "Cisco": 2,
      "HP": 1,
      "Unknown": 2
    },
    "scan_method_distribution": {
      "SNMP": 1,
      "ARP": 2,
      "COMBINED": 2
    },
    "avg_response_time_ms": 28
  },
  "scan_info": {
    "scan_type": "full",
    "network_range": "192.168.1.0/24",
    "snmp_communities": ["public", "private"],
    "timeout": 2,
    "retries": 1,
    "worker_count": 50
  }
}
```

### Type-Specific Scanning

**POST** `/api/v1/network/scan/snmp` (SNMP Only)

```json
{
  "network_range": "192.168.1.0/24",
  "communities": ["public", "private"],
  "timeout": 2,
  "retries": 1
}
```

**POST** `/api/v1/network/scan/arp` (ARP Only)

```json
{
  "network_range": "192.168.1.0/24",
  "timeout": 2,
  "retries": 1
}
```

### Quick Scan

**GET** `/api/v1/network/quick-scan?network=192.168.1.0/24&community=public`

```json
{
  "reachable_ips": ["192.168.1.1", "192.168.1.10", "192.168.1.20"],
  "count": 3
}
```

### Single Device Scan

**GET** `/api/v1/device/192.168.1.1?community=public&community=private`

```json
{
  "device": {
    "ip": "192.168.1.1",
    "mac_address": "AA:BB:CC:DD:EE:FF",
    "hostname": "router.local",
    "vendor": "Cisco",
    "model": "2960",
    "version": "15.0(2)SE4",
    "description": "Cisco IOS Software, C2960...",
    "contact": "admin@company.com",
    "location": "Server Room",
    "uptime": "45d 12h 30m 15s",
    "is_reachable": true,
    "response_time_ms": 23,
    "scan_method": "SNMP",
    "last_seen": "2024-01-15T10:30:00Z"
  }
}
```

### Scan Methods Information

**GET** `/api/v1/scan-methods`

```json
{
  "scan_methods": {
    "snmp": {
      "name": "SNMP Scan",
      "description": "Discovers devices using SNMP protocol...",
      "recommended_settings": {
        "timeout": "1-3 seconds",
        "retries": "0-1"
      }
    },
    "arp": {
      "name": "ARP Scan",
      "description": "Discovers devices using ARP protocol...",
      "recommended_settings": {
        "timeout": "1-2 seconds",
        "retries": "0"
      }
    },
    "full": {
      "name": "Full Scan (SNMP + ARP)",
      "description": "Combines both SNMP and ARP scanning methods...",
      "recommended_settings": {
        "timeout": "2-3 seconds",
        "retries": "0-1"
      }
    }
  },
  "default": "full",
  "recommended": "full"
}
```

### Vendor Database Management

**GET** `/api/v1/vendor-database`

```json
{
  "status": "vendor database loaded from JSON file",
  "config_path": "configs/oui_vendors.json",
  "description": "External JSON-based OUI vendor database"
}
```

**POST** `/api/v1/vendor-database/reload`

```json
{
  "status": "reload triggered",
  "message": "Vendor database reloaded from configs/oui_vendors.json",
  "timestamp": "2024-01-15T10:30:00Z"
}
```

### Network Range Validation

**GET** `/api/v1/network/validate?network=192.168.1.0/24`

```json
{
  "valid": true,
  "network": "192.168.1.0/24"
}
```

## 🛠️ Development

### Project Structure

```
network-discovery/
├── cmd/                    # Main application
│   └── main.go
├── internal/               # Internal packages
│   ├── api/               # HTTP handlers and routes
│   ├── discovery/         # Network discovery services
│   ├── models/            # Data models
│   ├── snmp/              # SNMP client
│   └── arp/               # ARP scanner and vendor management
├── frontend-build/        # Compiled web interface
│   └── dist/              # Static frontend files
├── configs/               # Configuration files
│   └── oui_vendors.json  # Vendor database
├── config.yaml            # Main configuration
├── go.mod                 # Go module definition
├── go.sum                 # Go dependency checksums
└── README.md             # Documentation
```

### Command Line Parameters

| Parameter    | Description           | Default                    |
| ------------ | --------------------- | -------------------------- |
| `-port`      | HTTP server port      | `8080`                     |
| `-host`      | HTTP server host      | `0.0.0.0`                  |
| `-log-level` | Log level             | `debug`                    |
| `-config`    | Vendor config file    | `configs/oui_vendors.json` |

### Environment Variables

| Variable       | Description            | Default |
| -------------- | ---------------------- | ------- |
| `SERVER_PORT`  | HTTP server port       | `8080`  |
| `LOG_LEVEL`    | Log level              | `info`  |
| `SNMP_TIMEOUT` | SNMP timeout           | `5s`    |
| `MAX_WORKERS`  | Maximum worker count   | `50`    |

### SNMP Information

The application uses the following SNMP OIDs:

- `1.3.6.1.2.1.1.1.0` - System Description
- `1.3.6.1.2.1.1.5.0` - System Name
- `1.3.6.1.2.1.1.4.0` - System Contact
- `1.3.6.1.2.1.1.6.0` - System Location
- `1.3.6.1.2.1.1.3.0` - System Uptime
- `1.3.6.1.2.1.2.2.1.6` - Interface Physical Address

## 🔒 Security

### SNMP Community Strings

SNMP community strings are sensitive information. In production:

- Use SNMP only on secure networks
- Change default community strings
- Use SNMPv3 when possible (coming in future releases)

## 🐛 Troubleshooting

### Common Issues

**Devices not being discovered:**

- Verify SNMP service is active
- Confirm community strings are correct
- Check firewall rules (UDP port 161)
- Verify SNMP service security settings for allowed connections

**ARP scan not working:**

- Verify ping command is available on the system
- Verify ARP command is available on the system
- Ensure target devices are on the same network segment

**Slow scanning:**

- Increase worker count (`max_workers`)
- Decrease timeout value
- Reduce retry count

**High memory usage:**

- Reduce worker count
- Narrow the scan range

### Debug Mode

```bash
# Run with debug logs
./network-discovery -log-level=debug

# Test a specific device
curl "http://localhost:8080/api/v1/device/192.168.1.1?community=public"

# Test full scan
curl -X POST http://localhost:8080/api/v1/network/full-scan \
  -H "Content-Type: application/json" \
  -d '{"network_range":"192.168.1.0/24","scan_type":"full"}'
```

### Log Analysis

```bash
# Filter successful scans
grep "Successfully queried device" /var/log/network-discovery.log

# View error messages
grep "ERROR" /var/log/network-discovery.log

# Performance metrics
grep "Scan completed" /var/log/network-discovery.log
```

## 📝 License

This project is licensed under the MIT License. See the `LICENSE` file for details.

## Planned Development

- 🐳 **Docker Support**: Easy installation and deployment
- 🔒 **SNMPv3 Support**: Authorization and encryption support
- 🔍 **Vulnerability Scanning**: Security vulnerability detection and assessment

##

⭐ **If you like this project, don't forget to give it a star!**

