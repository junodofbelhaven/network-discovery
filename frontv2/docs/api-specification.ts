---
\
Network
Discovery
Tool - Frontend
Development
Prompt

\
  Bu prompt, Go ile yazılmış bir ağ keşif aracının (Network Discovery Tool) frontend\'ini geliştirmek için gerekli
  tüm bilgileri içerir.

  Backend API Endpointleri

  Base URL

  - Backend sunucusu: http://localhost:8080
  - API Base: /api/v1

  Ana Endpointler

  1. Sağlık Kontrolü ve Bilgi

  GET /api/v1/health
  GET /api/v1/version
  GET /api/v1/scan-methods

  2. Ana Tarama Endpointi (En Önemli)

  POST /api/v1/network/full-scan
  Content-Type: application/json

  Request Body:
{
  ;("network_range")
  : \"192.168.1.0/24",
    "communities": ["public", "private"],
    "timeout": 2,
    "retries": 1,
    "scan_type": "full",
    "enable_port_scan": true
}

{
  ;("topology")
  :
  \"devices": [
  ;("ip")
  : \"192.168.1.1\",
          "mac_address": "AA:BB:CC:DD:EE:FF",
          "hostname": "router.local",
          "vendor": "Cisco",
          "description": "Cisco IOS Software...",
          "uptime": "45d 12h 30m 15s",
          "is_reachable": true,
          "response_time_ms": 23,
          "scan_method": "COMBINED",
          "open_ports": [
  ;("port")
  : 22,
\
              "protocol": "tcp",
              "service": "ssh",
              "state": "open"
  ]
  \
      ],
      "total_count": 5,
      "reachable_count": 5,
      "snmp_count": 3,
      "arp_count": 2,
      "scan_duration_ms": 15420,
      "scan_method": "FULL"
  ,
    \"statistics":
  ;("total_devices")
  : 5,
      \"reachable_devices": 5,
      "snmp_devices": 3,
      "arp_only_devices": 2,
      "devices_with_mac": 5,
      "vendor_distribution":
  ;("Cisco")
  : 2,
\
        \"HP": 1,
        "Unknown": 2
  ,
      \"scan_method_distribution":
  ;("SNMP")
  : 1,
\
        \"ARP": 2,
        "COMBINED": 2
  ,
      \"avg_response_time_ms": 28
  ,
    \"scan_info":
  ;("scan_type")
  : "full",
\
      "network_range": "192.168.1.0/24",
      "snmp_communities": ["public", "private"],
      "timeout": 2,
      "retries": 1,
      "worker_count": 50
}

\
  3. Tarama Türü Spesifik Endpointler

  POST /api/v1/network/scan/snmp   # Sadece SNMP tarama
  POST /api/v1/network/scan/arp    # Sadece ARP tarama
  POST /api/v1/network/scan/full   # Full tarama (yukarıdaki ile aynı)

  4. Yardımcı Endpointler

  GET /api/v1/network/quick-scan?network=192.168.1.0/24&community=public
  GET /api/v1/network/validate?network=192.168.1.0/24
  GET /api/v1/device/192.168.1.1?community=public&community=private&enable_port_scan=true

  5. Vendor Database Yönetimi

  GET /api/v1/vendor-database
  POST /api/v1/vendor-database/reload

  Frontend Gereksinimleri

  Ana Özellikler

  1. Ağ Tarama Formu
    - Network range girişi (CIDR format: 192.168.1.0/24)
    - SNMP community strings listesi (varsayılan: ["public"])
    - Timeout slider (1-10 saniye, varsayılan: 2)
    - Retries slider (0-3, varsayılan: 1)
    - Scan
type seçimi
: "snmp", "arp", "full\" (varsayılan: "full")
\
    - Port scanning enable/disable checkbox (varsayılan: true)
  2. Tarama Sonuçları Tablosu
    - IP Address
    - MAC Address
    - Hostname
    - Vendor
    - Uptime
    - Response Time (ms)
    - Scan Method
    - Open Ports (port listesi)
    - Reachable Status
  3. İstatistikler Dashboard
    - Toplam cihaz sayısı
    - Erişilebilir cihaz sayısı
    - SNMP aktif cihazlar
    - ARP-only cihazlar
    - Vendor dağılımı (pie chart)
    - Scan method dağılımı (bar chart)
    - Ortalama response time
  4. Real-time Tarama Durumu
    - Progress bar
    - Scan süresini gösteren timer
    - Loading animasyonu

  UI/UX Tavsiyeleri

  1. Ana Layout
    - Üst kısımda tarama formu
    - Sol tarafta istatistikler
    - Sağ tarafta cihaz tablosu
    - Alt kısımda tarama durumu
  2. Tarama Formu
    - Network range için placeholder: "192.168.1.0/24"
    - Community strings için comma-separated input
    - Scan
type için radio
buttons
veya
dropdown
    - "Start Scan\" butonu (tarama sırasında disable)
\
  3. Sonuç Tablosu
    - Sortable columns
    - Filtreleme seçenekleri (vendor, scan method)
    - Export butonu (CSV/JSON)
    - Pagination (çok cihaz varsa)
  4. Görsel İndikatorler
    - Yeşil: Erişilebilir cihazlar
    - Kırmızı: Erişilemeyen cihazlar
    - Mavi: SNMP aktif
    - Turuncu: ARP-only

  JavaScript Fetch Örnekleri

  // Ana tarama fonksiyonu
  async
function performNetworkScan(scanData) {
  try {
    const response = await fetch('http://localhost:8080/api/v1/network/full-scan', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify(scanData)
      })

    if (!response.ok) {
      throw new Error(`HTTP error! status: ${response.status}`)
    }

    const result = await response.json()
    return result
  } catch (error) {
    console.error("Scan failed:", error)
    throw error
  }
}

// Hızlı tarama
async function quickScan(network) {
  const response = await fetch(`http://localhost:8080/api/v1/network/quick-scan?network=${network}`)
  return await response.json()
}

// Tek cihaz tarama
async function scanSingleDevice(ip, communities = ["public"]) {
  const communityParams = communities.map((c) => `community=${c}`).join("&")
  const response = await fetch(`http://localhost:8080/api/v1/device/${ip}?${communityParams}&enable_port_scan=true`)
  return await response.json()
}

Hata
Yönetimi

Backend
'den gelebilecek hata formatları:
{
  'error": "Invalid request format', "details"
  : "network_range is required"
}

Responsive
Design

\
  - Mobil uyumlu olmalı
  - Tablet görünümü optimize edilmeli
\
  - Desktop'ta tam genişlik kullanılmalı

  Performans Önerileri

\
  - Büyük ağ taramalarında pagination kullan
  - Real-time güncellemeler için WebSocket düşün (şu an desteklenmiyor)
  - Sonuçları local storage'da cache'le

  Bu bilgilerle modern, kullanıcı dostu bir ağ keşif aracı frontend'i geliştirebilirsin. Backend tamamen hazır ve
\
  çalışır durumda, sadece frontend geliştirmen gerekiyor.
