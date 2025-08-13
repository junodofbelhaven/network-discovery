# 🌐 Network Discovery Tool

SNMP ve ARP protokolleri ile ağ topolojisini analiz eden, Go ile yazılmış modern bir ağ keşif aracı. Ağınızdaki tüm cihazları otomatik olarak keşfeder ve detaylı bilgilerini toplar.

## ✨ Özellikler

- 🔍 **Full Network Scan**: SNMP + ARP kombinasyonu ile kapsamlı ağ keşfi
- 📡 **SNMP v2c Desteği**: Detaylı cihaz bilgileri ile SNMP keşfi
- 🌐 **ARP Tarama**: Tüm IP-etkin cihazları keşfetme
- ⚡ **Yüksek Performans**: 50 eşzamanlı worker ile hızlı tarama
- 🏷️ **Vendor Algılama**: JSON tabanlı OUI veritabanı ile vendor tanıma
- 📱 **MAC Adresi Çözümleme**: Donanım adresi tanımlama
- ⏱️ **Yanıt Süresi Ölçümü**: Her cihaz için ağ gecikmesini ölçer
- 🌐 **REST API**: RESTful web servisleri ile kolay entegrasyon
- 💻 **Web Arayüzü**: Kullanıcı dostu web tabanlı kontrol paneli
- 📊 **Detaylı Raporlama**: Ağ istatistikleri ve cihaz envanteri

## 🚀 Hızlı Başlangıç

### Ön Gereksinimler

- Go 1.23 veya üzeri
- Docker (opsiyonel)
- Git

### Kurulum

```bash
# Projeyi klonla
git clone https://github.com/junodofbelhaven/network-discovery.git
cd network-discovery

# Bağımlılıkları yükle
go mod tidy

# Uygulamayı çalıştır
go run cmd/main.go
```

## Web GUI

## <img width="1755" height="1586" alt="image" src="https://github.com/user-attachments/assets/579b2e48-e90e-4626-8f8d-71ed4cb25da1" />


## 📖 API Dokümantasyonu

### Temel Endpoints

| Method | Endpoint                         | Açıklama                   |
| ------ | -------------------------------- | -------------------------- |
| GET    | `/api/v1/health`                 | Servis durumu kontrolü     |
| GET    | `/api/v1/version`                | Versiyon bilgisi           |
| GET    | `/api/v1/scan-methods`           | Tarama yöntemleri bilgisi  |
| POST   | `/api/v1/network/full-scan`      | Full scan (SNMP + ARP)     |
| POST   | `/api/v1/network/scan/snmp`      | Sadece SNMP taraması       |
| POST   | `/api/v1/network/scan/arp`       | Sadece ARP taraması        |
| POST   | `/api/v1/network/scan/full`      | Full scan (alternatif)     |
| POST   | `/api/v1/network/scan`           | Legacy SNMP taraması       |
| GET    | `/api/v1/network/quick-scan`     | Hızlı cihaz keşfi          |
| GET    | `/api/v1/network/validate`       | Ağ aralığı doğrulama       |
| GET    | `/api/v1/device/{ip}`            | Tek cihaz taraması         |
| GET    | `/api/v1/vendor-database`        | Vendor veritabanı bilgisi  |
| POST   | `/api/v1/vendor-database/reload` | Vendor veritabanı yenileme |

### Full Network Scan (Ana Endpoint)

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

**Yanıt:**

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

### Tip Özelinde Tarama

**POST** `/api/v1/network/scan/snmp` (Sadece SNMP)

```json
{
  "network_range": "192.168.1.0/24",
  "communities": ["public", "private"],
  "timeout": 2,
  "retries": 1
}
```

**POST** `/api/v1/network/scan/arp` (Sadece ARP)

```json
{
  "network_range": "192.168.1.0/24",
  "timeout": 2,
  "retries": 1
}
```

### Hızlı Tarama

**GET** `/api/v1/network/quick-scan?network=192.168.1.0/24&community=public`

```json
{
  "reachable_ips": ["192.168.1.1", "192.168.1.10", "192.168.1.20"],
  "count": 3
}
```

### Tek Cihaz Taraması

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

### Tarama Yöntemleri Bilgisi

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

### Vendor Veritabanı Yönetimi

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

### Ağ Aralığı Doğrulama

**GET** `/api/v1/network/validate?network=192.168.1.0/24`

```json
{
  "valid": true,
  "network": "192.168.1.0/24"
}
```

## 🛠️ Geliştirme

### Proje Yapısı

```
network-discovery/
├── cmd/                    # Ana uygulama
│   └── main.go
├── internal/               # İç paketler
│   ├── api/               # HTTP handlers ve routes
│   ├── discovery/         # Ağ keşif servisleri
│   ├── models/            # Veri modelleri
│   ├── snmp/              # SNMP istemcisi
│   └── arp/               # ARP tarayıcı ve vendor yönetimi
├── vue-front/             # Frontend kaynak kodları
│   └── frontend/          # Vue.js uygulaması
├── frontend-build/        # Derlenmiş web arayüzü
│   └── dist/              # Statik dosyalar
├── configs/               # Konfigürasyon dosyaları
│   └── oui_vendors.json  # Vendor veritabanı
├── config.yaml            # Ana konfigürasyon
├── go.mod                 # Go modül tanımı
├── go.sum                 # Go bağımlılık sağlama toplamları
└── README.md             # Dokümantasyon
```

### Komut Satırı Parametreleri

| Parametre    | Açıklama              | Varsayılan                 |
| ------------ | --------------------- | -------------------------- |
| `-port`      | HTTP sunucu portu     | `8080`                     |
| `-host`      | HTTP sunucu host'u    | `0.0.0.0`                  |
| `-log-level` | Log seviyesi          | `debug`                    |
| `-config`    | Vendor config dosyası | `configs/oui_vendors.json` |

### Ortam Değişkenleri

| Değişken       | Açıklama               | Varsayılan |
| -------------- | ---------------------- | ---------- |
| `SERVER_PORT`  | HTTP sunucu portu      | `8080`     |
| `LOG_LEVEL`    | Log seviyesi           | `info`     |
| `SNMP_TIMEOUT` | SNMP timeout           | `5s`       |
| `MAX_WORKERS`  | Maksimum worker sayısı | `50`       |

## 📊 Desteklenen Cihazlar

### Vendor Desteği

- ✅ **Cisco**: IOS, NX-OS, IOS-XE, ASA
- ✅ **Juniper**: JunOS (SRX, MX, EX, QFX serisi)
- ✅ **Huawei**: VRP (S5700, S6700, CloudEngine)
- ✅ **HP/HPE**: ProCurve, Aruba
- ✅ **Dell**: PowerConnect, Force10, OS10
- ✅ **MikroTik**: RouterOS, RouterBoard
- ✅ **Ubiquiti**: UniFi, EdgeMax
- ✅ **Fortinet**: FortiGate, FortiOS
- ✅ **Palo Alto**: PA serisi
- ✅ **Netgear**: ProSafe serisi
- ✅ **D-Link**: DGS, DES serisi
- ✅ **TP-Link**: Managed switch'ler
- ✅ **Apple**: Mac cihazları
- ✅ **Intel**: Network kartları
- ✅ **VMware**: Sanal makineler
- ✅ **Raspberry Pi**: IoT cihazları

### SNMP Bilgileri

Uygulama aşağıdaki SNMP OID'lerini kullanır:

- `1.3.6.1.2.1.1.1.0` - System Description
- `1.3.6.1.2.1.1.5.0` - System Name
- `1.3.6.1.2.1.1.4.0` - System Contact
- `1.3.6.1.2.1.1.6.0` - System Location
- `1.3.6.1.2.1.1.3.0` - System Uptime
- `1.3.6.1.2.1.2.2.1.6` - Interface Physical Address

## 🔒 Güvenlik

### SNMP Community Strings

SNMP community string'leri hassas bilgilerdir. Üretim ortamında:

- SNMP'yi sadece güvenli ağlarda kullanın
- Varsayılan community string'leri değiştirin
- Mümkünse SNMPv3 kullanın (gelecek sürümlerde)

## 🐛 Sorun Giderme

### Yaygın Sorunlar

**Cihazlar keşfedilmiyor:**

- SNMP servisinin aktif olduğunu kontrol edin
- Community string'lerin doğru olduğunu doğrulayın
- Firewall kurallarını kontrol edin (UDP 161 portu)
- SNMP servisinin güvenlik kısmından hangi bağlantılardan bağlantı kabul ettiğini kontrol edin

**ARP tarama çalışmıyor:**

- Ping komutunun sistem üzerinde mevcut olduğunu kontrol edin
- ARP komutunun sistem üzerinde mevcut olduğunu kontrol edin
- Hedef cihazların aynı ağ segmentinde olduğunu kontrol edin

**Yavaş tarama:**

- Worker sayısını artırın (`max_workers`)
- Timeout değerini azaltın
- Retry sayısını azaltın

**Memory kullanımı yüksek:**

- Worker sayısını azaltın
- Tarama aralığını küçültün

### Debug Modu

```bash
# Debug logları ile çalıştır
./network-discovery -log-level=debug

# Belirli bir cihazı test et
curl "http://localhost:8080/api/v1/device/192.168.1.1?community=public"

# Full scan test et
curl -X POST http://localhost:8080/api/v1/network/full-scan \
  -H "Content-Type: application/json" \
  -d '{"network_range":"192.168.1.0/24","scan_type":"full"}'
```

### Log Analizi

```bash
# Başarılı taramaları filtrele
grep "Successfully queried device" /var/log/network-discovery.log

# Hata mesajlarını görüntüle
grep "ERROR" /var/log/network-discovery.log

# Performans metrikleri
grep "Scan completed" /var/log/network-discovery.log
```

## 📝 Lisans

Bu proje MIT lisansı altında lisanslanmıştır. Detaylar için `LICENSE` dosyasına bakın.

## Planlanan geliştirme

- 🐳 **Docker Desteği**: Kolay kurulum ve deployment
- 🔒 **SNMPv3 Desteği**: Authorization ve encryption desteği

##

⭐ **Bu projeyi beğendiyseniz yıldız vermeyi unutmayın!**
