# 🌐 Network Discovery Tool

SNMP protokolü ile ağ topolojisini analiz eden, Go ile yazılmış modern bir ağ keşif aracı. Ağınızdaki tüm SNMP-etkin cihazları otomatik olarak keşfeder ve detaylı bilgilerini toplar.

## ✨ Özellikler

- 🔍 **Otomatik Ağ Keşfi**: CIDR notasyonu ile belirtilen ağ aralıklarını tarar
- 📡 **SNMP v2c Desteği**: Standart SNMP protokolü ile cihaz bilgilerini toplar
- ⚡ **Yüksek Performans**: 50 eşzamanlı worker ile hızlı tarama
- 🏷️ **Vendor Algılama**: Cisco, Juniper, Huawei, HP gibi popüler markaları otomatik tanır
- ⏱️ **Yanıt Süresi Ölçümü**: Her cihaz için ağ gecikmesini ölçer
- 🌐 **REST API**: RESTful web servisleri ile kolay entegrasyon
- 💻 **Web Arayüzü**: Kullanıcı dostu web tabanlı kontrol paneli
- 📊 **Detaylı Raporlama**: Ağ istatistikleri ve cihaz envantери

## 🚀 Hızlı Başlangıç

### Ön Gereksinimler

- Go 1.21 veya üzeri
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

## Web GUI'ına bağlanma

## <img width="1775" height="691" alt="image" src="https://github.com/user-attachments/assets/e027deab-6732-4690-ad4c-49f9371ec50d" />

Uygulama şuanda GUI olarak sadece belirtilen ip aralığında scan yapmaya yarayan bir web page sunuyor. GUI'a bağlanmak için;

- Uygulamayı başlatın
- Herhangi bir web browserdan localhost:{port}/index sayfasına girin (Belirtilmediği sürece varsayılan port 8080 olarak başlar)

## 📖 API Dokümantasyonu

### Temel Endpoints

| Method | Endpoint                     | Açıklama               |
| ------ | ---------------------------- | ---------------------- |
| GET    | `/api/v1/health`             | Servis durumu kontrolü |
| GET    | `/api/v1/version`            | Versiyon bilgisi       |
| POST   | `/api/v1/network/scan`       | Ağ taraması başlat     |
| GET    | `/api/v1/network/quick-scan` | Hızlı cihaz keşfi      |
| GET    | `/api/v1/device/{ip}`        | Tek cihaz taraması     |

### Ağ Taraması

**POST** `/api/v1/network/scan`

```json
{
  "network_range": "192.168.1.0/24",
  "communities": ["public", "private"],
  "timeout": 5,
  "retries": 2
}
```

**Yanıt:**

```json
{
  "topology": {
    "devices": [
      {
        "ip": "192.168.1.1",
        "hostname": "router.local",
        "vendor": "Cisco",
        "description": "Cisco IOS Software...",
        "uptime": "45d 12h 30m 15s",
        "is_reachable": true,
        "response_time_ms": 23
      }
    ],
    "total_count": 5,
    "reachable_count": 3,
    "scan_duration_ms": 15420
  },
  "statistics": {
    "vendor_distribution": {
      "Cisco": 2,
      "HP": 1
    },
    "avg_response_time_ms": 28
  }
}
```

### Hızlı Tarama

**GET** `/api/v1/network/quick-scan?network=192.168.1.0/24`

```json
{
  "reachable_ips": ["192.168.1.1", "192.168.1.10", "192.168.1.20"],
  "count": 3
}
```

### Tek Cihaz Taraması

**GET** `/api/v1/device/192.168.1.1?community=public`

```json
{
  "device": {
    "ip": "192.168.1.1",
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
    "last_seen": "2024-01-15T10:30:00Z"
  }
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
│   └── snmp/              # SNMP istemcisi
├── pkg/                   # Genel yardımcı paketler
├── frontend/              # Web arayüzü
├── config.yaml            # Konfigürasyon
├── docker-compose.yml     # Docker yapılandırması
├── Dockerfile             # Container tanımı
├── Makefile              # Build ve deployment
└── README.md             # Dokümantasyon
```

## 🔧 Konfigürasyon

Uygulama `config.yaml` dosyası ile yapılandırılabilir:

```yaml
server:
  host: "0.0.0.0"
  port: 8080

snmp:
  timeout: 5s
  retries: 2
  default_communities:
    - "public"
    - "private"
    - "community"

scanning:
  max_workers: 50
  max_scan_duration: 10m

logging:
  level: "info"
  format: "json"
```

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

### SNMP Bilgileri

Uygulama aşağıdaki SNMP OID'lerini kullanır:

- `1.3.6.1.2.1.1.1.0` - System Description
- `1.3.6.1.2.1.1.5.0` - System Name
- `1.3.6.1.2.1.1.4.0` - System Contact
- `1.3.6.1.2.1.1.6.0` - System Location
- `1.3.6.1.2.1.1.3.0` - System Uptime

## 🔒 Güvenlik

### SNMP Community Strings

SNMP community string'leri hassas bilgilerdir. Üretim ortamında:

- SNMPyi sadece güvenli ağlarda kullanın 



## 🐛 Sorun Giderme

### Yaygın Sorunlar

**Cihazlar keşfedilmiyor:**

- SNMP servisinin aktif olduğunu kontrol edin
- Community string'lerin doğru olduğunu doğrulayın
- Firewall kurallarını kontrol edin (UDP 161 portu)
- SNMP servisinin güvenlik kısmından hangi bağlantılardan bağlantı kabul ettiğini kontrol edin

**Yavaş tarama:**

- Worker sayısını artırın (`max_workers`)
- Timeout değerini azaltın

**Memory kullanımı yüksek:**

- Worker sayısını azaltın
- Tarama aralığını küçültün

### Debug Modu

```bash
# Debug logları ile çalıştır
./network-discovery -log-level=debug

# Belirli bir cihazı test et
curl "http://localhost:8080/api/v1/device/192.168.1.1?community=public"
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
- ❗️ **IMCP ile Ping**: SNMP desteklemeyen cihazlara ping yollayarak ağı tarama

##

⭐ **Bu projeyi beğendiyseniz yıldız vermeyi unutmayın!**






