# Network Discovery - Kod Sadeleştirme Önerileri

## 🎯 Ana Hedefler
- Kod tekrarını azaltma
- Karmaşıklığı düşürme
- Performansı artırma
- Bakımı kolaylaştırma

## 🔧 Önerilen Değişiklikler

### 1. Logger Yapısını Sadeleştirme
```go
// Önerilen: Tek bir global logger manager
type LoggerManager struct {
    logger *logrus.Logger
}

func NewLoggerManager(level logrus.Level) *LoggerManager {
    logger := logrus.New()
    logger.SetLevel(level)
    logger.SetFormatter(&logrus.JSONFormatter{})
    return &LoggerManager{logger: logger}
}
```

### 2. Factory Pattern ile Constructor Sadeleştirme
```go
// Mevcut durumda çok fazla constructor var:
// - NewClient() 
// - NewClientWithLogger()
// - NewScanner()
// - NewScannerWithLogger()

// Önerilen: Tek factory function
type Config struct {
    Timeout    time.Duration
    Retries    int
    Logger     *logrus.Logger
    MaxWorkers int
}

func NewSNMPClient(cfg Config) *Client {
    if cfg.Logger == nil {
        cfg.Logger = logrus.New()
    }
    return &Client{
        timeout: cfg.Timeout,
        retries: cfg.Retries,
        logger:  cfg.Logger,
    }
}
```

### 3. Error Handling Sadeleştirme
```go
// Mevcut: Her OID için ayrı error handling
// Önerilen: Bulk error handling

func (c *Client) queryBasicOIDs(client *gosnmp.GoSNMP, device *models.Device) {
    oids := []string{OIDSysDescr, OIDSysName, OIDSysContact, OIDSysLocation, OIDSysUptime}
    
    result, err := client.Get(oids)
    if err != nil {
        c.logger.Debugf("Bulk OID query failed: %v", err)
        return
    }
    
    for i, variable := range result.Variables {
        c.parseOIDVariable(oids[i], variable, device)
    }
}
```

### 4. Vendor Detection Birleştirme
```go
// Önerilen: Unified vendor detection
type VendorDetector struct {
    snmpPatterns map[string]string
    ouiDatabase  map[string]string
}

func (vd *VendorDetector) DetectVendor(device *models.Device) string {
    // 1. Try MAC OUI first (more reliable)
    if device.MACAddress != "" {
        if vendor := vd.detectFromMAC(device.MACAddress); vendor != "" {
            return vendor
        }
    }
    
    // 2. Fallback to SNMP description
    if device.Description != "" {
        return vd.detectFromDescription(device.Description)
    }
    
    return "Unknown"
}
```

### 5. Scanner Interface Standardization
```go
// Önerilen: Common interface
type NetworkScanner interface {
    Scan(networkRange string, options ScanOptions) (*models.NetworkTopology, error)
}

type ScanOptions struct {
    Communities    []string
    Timeout        time.Duration
    Retries        int
    EnablePortScan bool
    ScanType       string
}

// Unified scanner
type UnifiedScanner struct {
    snmpClient  *snmp.Client
    arpClient   *arp.Scanner
    portClient  *ports.Scanner
    logger      *logrus.Logger
}
```

### 6. Configuration Simplification
```go
// Mevcut: Birden fazla yerde hardcoded değerler
// Önerilen: Centralized configuration

type ScannerConfig struct {
    DefaultCommunities []string        `yaml:"default_communities"`
    DefaultTimeout     time.Duration   `yaml:"default_timeout"`
    DefaultRetries     int             `yaml:"default_retries"`
    MaxWorkers         int             `yaml:"max_workers"`
    EnablePortScan     bool            `yaml:"enable_port_scan"`
    VendorDBPath       string          `yaml:"vendor_db_path"`
}

func LoadConfig(path string) (*ScannerConfig, error) {
    // YAML config loading logic
}
```

### 7. Dependency Injection
```go
// Önerilen: Dependency container
type Container struct {
    Config      *ScannerConfig
    Logger      *logrus.Logger
    VendorMgr   *VendorManager
    SNMPClient  *snmp.Client
    ARPScanner  *arp.Scanner
    PortScanner *ports.Scanner
}

func NewContainer(configPath string) (*Container, error) {
    config, err := LoadConfig(configPath)
    if err != nil {
        return nil, err
    }
    
    logger := NewLogger(config.LogLevel)
    
    return &Container{
        Config:      config,
        Logger:      logger,
        VendorMgr:   NewVendorManager(config.VendorDBPath, logger),
        SNMPClient:  NewSNMPClient(config, logger),
        ARPScanner:  NewARPScanner(config, logger),
        PortScanner: NewPortScanner(config, logger),
    }, nil
}
```

## 🚀 Faydalar

### Performance Improvements
- Daha az memory allocation
- Efficient bulk operations
- Reduced goroutine overhead

### Code Quality
- %30 daha az kod satırı
- Unified interfaces
- Better testability
- Cleaner error handling

### Maintainability
- Single source of truth for config
- Standardized patterns
- Easier to extend
- Better documentation

## 📋 Implementation Plan

1. **Phase 1**: Logger ve Config unification
2. **Phase 2**: Scanner interface standardization  
3. **Phase 3**: Vendor detection merge
4. **Phase 4**: Error handling simplification
5. **Phase 5**: Dependency injection implementation

## 🧪 Testing Strategy

- Unit test coverage'ı artırma
- Integration testler ekleme
- Performance benchmarking
- Memory profiling

## 📈 Expected Results

- **Code Lines**: 2000+ → 1400-1500 satır
- **Complexity**: Cyclomatic complexity %25 düşüş
- **Performance**: %15-20 hız artışı
- **Memory**: %10-15 daha az RAM kullanımı