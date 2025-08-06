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
- 🐳 **Docker Desteği**: Kolay kurulum ve deployment
- 📊 **Detaylı Raporlama**: Ağ istatistikleri ve cihaz envantери

## 🚀 Hızlı Başlangıç

### Ön Gereksinimler

- Go 1.21 veya üzeri
- Docker (opsiyonel)
- Git

### Kurulum

```bash
# Projeyi klonla
git clone <repository-url>
cd network-discovery

# Bağımlılıkları yükle
go mod tidy

# Uygulamayı çalıştır
make dev
```

### Docker ile Çalıştırma

```bash
# Docker image'ını oluştur ve çalıştır
make docker-compose-up
```

## 📖 API Dokümantasyonu

### Temel Endpoints

| Method | Endpoint                                    | Açıklama               |
| ------ | ------------------------------------------- | ---------------------- |
| GET    | `/api/v1/health`                            | Servis durumu kontrolü |
| GET    | `/api/v1/version`                           | Versiyon bilgisi       |
| POST   | `/api/v1/network/scan`                      | Ağ taraması başlat     |
| GET    | `/api/v1/network/quick-scan`                | Hızlı cihaz keşfi      |
| GET    | `/api/v1/device/{ip}?community={community}` | Tek cihaz taraması     |

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

### Komutlar

```bash
# Geliştirme modunda ç
```
