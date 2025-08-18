<!-- ===== DeviceList.vue ===== -->
<template>
  <div class="device-list" v-if="devices.length > 0">
    <!-- Header Section -->
    <div class="dashboard-header">
      <div class="header-content">
        <div class="header-title">
          <span class="header-icon">🌐</span>
          <h2>Network Discovery Dashboard</h2>
          <p class="header-subtitle">Real-time device monitoring and analysis</p>
        </div>
        <div class="stats-container">
          <div class="stat-card">
            <div class="stat-value">{{ devices.length }}</div>
            <div class="stat-label">Total Devices</div>
          </div>
          <div class="stat-card">
            <div class="stat-value">{{ onlineDevices }}</div>
            <div class="stat-label">Online</div>
          </div>
          <div class="stat-card">
            <div class="stat-value">{{ snmpDevices.length }}</div>
            <div class="stat-label">SNMP</div>
          </div>
          <div class="stat-card">
            <div class="stat-value">{{ arpDevices.length }}</div>
            <div class="stat-label">ARP</div>
          </div>
        </div>
      </div>
    </div>

    <!-- Filter Tabs -->
    <div class="filter-tabs">
      <button
        :class="['filter-tab', { active: currentFilter === 'all' }]"
        @click="setFilter('all')"
      >
        <span class="tab-icon">📊</span>
        All Devices
        <span class="tab-count">{{ devices.length }}</span>
      </button>
      <button
        :class="['filter-tab', { active: currentFilter === 'snmp' }]"
        @click="setFilter('snmp')"
      >
        <span class="tab-icon">📡</span>
        SNMP Protocol
        <span class="tab-count">{{ snmpDevices.length }}</span>
      </button>
      <button
        :class="['filter-tab', { active: currentFilter === 'arp' }]"
        @click="setFilter('arp')"
      >
        <span class="tab-icon">🌐</span>
        ARP Protocol
        <span class="tab-count">{{ arpDevices.length }}</span>
      </button>
    </div>

    <!-- Device Cards Grid -->
    <div class="cards-grid">
      <transition-group name="card-animation">
        <div
          v-for="device in filteredDevices"
          :key="device.ip"
          :class="['device-card', device.scan_method === 'ARP' ? 'arp-card' : 'snmp-card']"
        >
          <!-- Card Header with IP, Hostname and Status -->
          <div class="card-header">
            <div class="card-title-section">
              <div class="card-ip">
                <span class="protocol-icon">
                  {{ device.scan_method === 'SNMP' ? '📡' : '🌐' }}
                </span>
                {{ device.ip }}
              </div>
              <div v-if="device.hostname" class="card-hostname">
                <span class="hostname-icon">🏷️</span>
                {{ device.hostname }}
              </div>
            </div>
            <span :class="['status-badge', device.scan_method === 'SNMP' ? 'online' : 'active']">
              <span class="status-dot"></span>
              {{ device.scan_method === 'SNMP' ? 'Online' : 'Active' }}
            </span>
          </div>

          <!-- Card Body with Device Information -->
          <div class="card-body">
            <!-- Vendor Information -->
            <div v-if="device.vendor" class="info-row">
              <span class="info-label">Vendor:</span>
              <span class="info-value">{{ device.vendor }}</span>
            </div>

            <!-- MAC Address -->
            <div v-if="device.mac_address" class="info-row">
              <span class="info-label">MAC:</span>
              <span class="info-value mac-address">{{ device.mac_address }}</span>
            </div>

            <!-- System Description (Full SNMP info) -->
            <div v-if="device.description" class="info-row">
              <span class="info-label">System:</span>
              <span class="info-value system-desc">{{ device.description }}</span>
            </div>

            <!-- Model Information -->
            <div v-if="device.model" class="info-row">
              <span class="info-label">Model:</span>
              <span class="info-value">{{ device.model }}</span>
            </div>

            <!-- Version/Firmware -->
            <div v-if="device.version" class="info-row">
              <span class="info-label">Version:</span>
              <span class="info-value">{{ device.version }}</span>
            </div>

            <!-- Contact Information -->
            <div v-if="device.contact" class="info-row">
              <span class="info-label">Contact:</span>
              <span class="info-value">{{ device.contact }}</span>
            </div>

            <!-- Location -->
            <div v-if="device.location" class="info-row">
              <span class="info-label">Location:</span>
              <span class="info-value">{{ device.location }}</span>
            </div>

            <!-- Uptime -->
            <div v-if="device.uptime" class="info-row">
              <span class="info-label">Uptime:</span>
              <span class="info-value uptime">{{ device.uptime }}</span>
            </div>

            <!-- Last Seen -->
            <div v-if="device.last_seen" class="info-row">
              <span class="info-label">Last Seen:</span>
              <span class="info-value">{{ formatDate(device.last_seen) }}</span>
            </div>
          </div>

          <!-- Card Footer with Protocol Type and Response Time -->
          <div class="card-footer">
            <span class="device-type-badge">
              <span class="badge-icon">
                {{ device.scan_method === 'SNMP' ? '📡' : '🌐' }}
              </span>
              {{ device.scan_method }} Protocol
            </span>
            <div class="footer-stats">
              <span v-if="device.response_time_ms" class="response-time">
                <span class="time-icon">⚡</span>
                {{ device.response_time_ms }}ms
              </span>
              <span v-if="device.is_reachable" class="reachable-status">
                <span class="check-icon">✓</span>
                Reachable
              </span>
            </div>
          </div>
        </div>
      </transition-group>
    </div>

    <!-- Empty State for Filtered Results -->
    <div v-if="filteredDevices.length === 0 && devices.length > 0" class="empty-state">
      <div class="empty-icon">🔍</div>
      <h3 class="empty-title">No {{ currentFilter }} devices found</h3>
      <p class="empty-text">Try selecting a different filter</p>
    </div>
  </div>

  <!-- Initial Empty State -->
  <div v-else class="empty-state initial">
    <div class="empty-icon">📡</div>
    <h2 class="empty-title">No Devices Discovered Yet</h2>
    <p class="empty-text">Start a network scan to discover devices on your network</p>
  </div>
</template>

<script>
export default {
  name: 'DeviceList',
  props: {
    devices: {
      type: Array,
      required: true,
    },
  },
  data() {
    return {
      currentFilter: 'all',
    }
  },
  computed: {
    snmpDevices() {
      return this.devices.filter(
        (device) => device.scan_method === 'SNMP' || device.scan_method === 'COMBINED',
      )
    },
    arpDevices() {
      return this.devices.filter((device) => device.scan_method === 'ARP')
    },
    onlineDevices() {
      return this.devices.filter((device) => device.is_reachable).length
    },
    filteredDevices() {
      switch (this.currentFilter) {
        case 'snmp':
          return this.snmpDevices
        case 'arp':
          return this.arpDevices
        default:
          return this.devices
      }
    },
  },
  methods: {
    setFilter(filter) {
      this.currentFilter = filter
    },
    formatDate(dateString) {
      if (!dateString) return 'N/A'
      const date = new Date(dateString)
      return date.toLocaleString()
    },
  },
}
</script>

<style scoped>
/* ===== ROOT VARIABLES ===== */
:root {
  --primary-gradient: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  --snmp-color: #667eea;
  --arp-color: #48bb78;
  --bg-primary: #0f0e17;
  --bg-secondary: #1a1928;
  --card-bg: #232136;
  --text-primary: #fffffe;
  --text-secondary: #b8b5ff;
  --text-muted: #94a1b2;
  --success-color: #48bb78;
  --warning-color: #f6ad55;
  --border-radius: 16px;
  --card-shadow: 0 10px 30px rgba(0, 0, 0, 0.3);
}

/* ===== DEVICE LIST CONTAINER ===== */
.device-list {
  margin-top: 30px;
  max-width: 1600px;
  margin-left: auto;
  margin-right: auto;
  padding: 0 20px;
}

/* ===== DASHBOARD HEADER ===== */
.dashboard-header {
  background: var(--primary-gradient);
  padding: 2rem;
  border-radius: var(--border-radius);
  margin-bottom: 2rem;
  box-shadow: var(--card-shadow);
}

.header-content {
  display: flex;
  justify-content: space-between;
  align-items: center;
  flex-wrap: wrap;
  gap: 2rem;
}

.header-title {
  flex: 1;
  min-width: 300px;
}

.header-title h2 {
  font-size: 1.75rem;
  margin: 0.5rem 0;
  color: var(--text-primary);
  display: flex;
  align-items: center;
  gap: 0.75rem;
}

.header-icon {
  font-size: 2rem;
}

.header-subtitle {
  color: rgba(255, 255, 255, 0.9);
  font-size: 0.95rem;
  margin: 0;
}

/* ===== STATS CONTAINER ===== */
.stats-container {
  display: flex;
  gap: 1.5rem;
  flex-wrap: wrap;
}

.stat-card {
  background: rgba(255, 255, 255, 0.1);
  backdrop-filter: blur(10px);
  padding: 1rem 1.5rem;
  border-radius: 12px;
  min-width: 100px;
  text-align: center;
}

.stat-value {
  font-size: 1.75rem;
  font-weight: 700;
  color: var(--text-primary);
}

.stat-label {
  font-size: 0.875rem;
  color: rgba(255, 255, 255, 0.8);
  text-transform: uppercase;
  letter-spacing: 0.5px;
  margin-top: 0.25rem;
}

/* ===== FILTER TABS ===== */
.filter-tabs {
  display: flex;
  gap: 1rem;
  margin-bottom: 2rem;
  background: var(--bg-secondary);
  padding: 0.5rem;
  border-radius: var(--border-radius);
  flex-wrap: wrap;
}

.filter-tab {
  flex: 1;
  min-width: 150px;
  padding: 0.875rem 1.25rem;
  background: transparent;
  border: 2px solid transparent;
  color: var(--text-secondary);
  font-size: 0.95rem;
  font-weight: 600;
  cursor: pointer;
  border-radius: 12px;
  transition: all 0.3s ease;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 0.5rem;
}

.filter-tab:hover {
  background: rgba(184, 181, 255, 0.1);
  border-color: rgba(184, 181, 255, 0.3);
}

.filter-tab.active {
  background: var(--primary-gradient);
  color: var(--text-primary);
  border-color: transparent;
}

.tab-icon {
  font-size: 1.1rem;
}

.tab-count {
  display: inline-block;
  margin-left: 0.25rem;
  padding: 0.125rem 0.5rem;
  background: rgba(255, 255, 255, 0.2);
  border-radius: 10px;
  font-size: 0.875rem;
}

/* ===== CARDS GRID ===== */
.cards-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(400px, 1fr));
  gap: 1.5rem;
  margin-bottom: 2rem;
}

@media (max-width: 900px) {
  .cards-grid {
    grid-template-columns: repeat(auto-fill, minmax(350px, 1fr));
  }
}

@media (max-width: 768px) {
  .cards-grid {
    grid-template-columns: 1fr;
  }
}

/* ===== DEVICE CARD ===== */
.device-card {
  background: var(--card-bg);
  border-radius: var(--border-radius);
  padding: 1.5rem;
  box-shadow: var(--card-shadow);
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
  position: relative;
  overflow: hidden;
  border-left: 4px solid var(--snmp-color);
}

.device-card.arp-card {
  border-left-color: var(--arp-color);
}

.device-card:hover {
  transform: translateY(-4px);
  box-shadow: 0 20px 40px rgba(0, 0, 0, 0.4);
}

.device-card::before {
  content: '';
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  height: 2px;
  background: var(--primary-gradient);
  transform: scaleX(0);
  transition: transform 0.3s ease;
}

.device-card:hover::before {
  transform: scaleX(1);
}

/* ===== CARD HEADER ===== */
.card-header {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  margin-bottom: 1.25rem;
  padding-bottom: 1rem;
  border-bottom: 1px solid rgba(184, 181, 255, 0.1);
}

.card-title-section {
  flex: 1;
}

.card-ip {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  font-size: 1.25rem;
  font-weight: 700;
  color: var(--text-primary);
  margin-bottom: 0.25rem;
}

.protocol-icon {
  font-size: 1.25rem;
}

.card-hostname {
  display: flex;
  align-items: center;
  gap: 0.375rem;
  color: var(--text-secondary);
  font-size: 0.95rem;
  font-weight: 500;
}

.hostname-icon {
  font-size: 0.9rem;
}

/* ===== STATUS BADGE ===== */
.status-badge {
  display: inline-flex;
  align-items: center;
  gap: 0.375rem;
  padding: 0.375rem 0.875rem;
  border-radius: 20px;
  font-size: 0.75rem;
  font-weight: 600;
  text-transform: uppercase;
  letter-spacing: 0.5px;
}

.status-badge.online {
  background: rgba(72, 187, 120, 0.2);
  color: var(--success-color);
}

.status-badge.active {
  background: rgba(246, 173, 85, 0.2);
  color: var(--warning-color);
}

.status-dot {
  width: 6px;
  height: 6px;
  border-radius: 50%;
  background: currentColor;
  animation: pulse 2s infinite;
}

@keyframes pulse {
  0%,
  100% {
    opacity: 1;
  }
  50% {
    opacity: 0.5;
  }
}

/* ===== CARD BODY ===== */
.card-body {
  display: flex;
  flex-direction: column;
  gap: 0.875rem;
}

.info-row {
  display: flex;
  align-items: flex-start;
  gap: 0.75rem;
  font-size: 0.9rem;
}

.info-label {
  font-weight: 600;
  color: var(--text-secondary);
  min-width: 80px;
  flex-shrink: 0;
}

.info-value {
  color: var(--text-primary);
  font-weight: 400;
  word-break: break-word;
  flex: 1;
}

.info-value.mac-address {
  font-family: 'Courier New', monospace;
  font-size: 0.85rem;
  background: rgba(184, 181, 255, 0.1);
  padding: 0.25rem 0.5rem;
  border-radius: 6px;
  letter-spacing: 0.5px;
}

.info-value.system-desc {
  font-size: 0.85rem;
  line-height: 1.4;
}

.info-value.uptime {
  color: var(--success-color);
  font-weight: 500;
}

/* ===== CARD FOOTER ===== */
.card-footer {
  margin-top: 1.25rem;
  padding-top: 1rem;
  border-top: 1px solid rgba(184, 181, 255, 0.1);
  display: flex;
  justify-content: space-between;
  align-items: center;
  flex-wrap: wrap;
  gap: 0.75rem;
}

.device-type-badge {
  display: inline-flex;
  align-items: center;
  gap: 0.375rem;
  padding: 0.375rem 0.875rem;
  border-radius: 8px;
  font-size: 0.8rem;
  font-weight: 600;
  background: rgba(102, 126, 234, 0.2);
  color: var(--snmp-color);
}

.arp-card .device-type-badge {
  background: rgba(72, 187, 120, 0.2);
  color: var(--arp-color);
}

.badge-icon {
  font-size: 0.9rem;
}

.footer-stats {
  display: flex;
  align-items: center;
  gap: 1rem;
}

.response-time {
  display: flex;
  align-items: center;
  gap: 0.25rem;
  font-size: 0.85rem;
  color: var(--text-secondary);
  font-weight: 500;
}

.reachable-status {
  display: flex;
  align-items: center;
  gap: 0.25rem;
  font-size: 0.85rem;
  color: var(--success-color);
  font-weight: 500;
}

.time-icon,
.check-icon {
  font-size: 0.9rem;
}

/* ===== EMPTY STATE ===== */
.empty-state {
  text-align: center;
  padding: 4rem 2rem;
  background: var(--bg-secondary);
  border-radius: var(--border-radius);
  margin: 2rem 0;
}

.empty-state.initial {
  background: var(--card-bg);
  border: 2px dashed rgba(184, 181, 255, 0.3);
}

.empty-icon {
  font-size: 4rem;
  margin-bottom: 1rem;
  opacity: 0.5;
}

.empty-title {
  font-size: 1.5rem;
  margin-bottom: 0.5rem;
  color: var(--text-primary);
}

.empty-text {
  color: var(--text-muted);
  font-size: 0.95rem;
}

/* ===== ANIMATIONS ===== */
.card-animation-enter-active,
.card-animation-leave-active {
  transition: all 0.3s ease;
}

.card-animation-enter-from {
  opacity: 0;
  transform: translateY(20px);
}

.card-animation-leave-to {
  opacity: 0;
  transform: translateY(-20px);
}

/* ===== RESPONSIVE ADJUSTMENTS ===== */
@media (max-width: 768px) {
  .header-content {
    flex-direction: column;
    align-items: flex-start;
  }

  .stats-container {
    width: 100%;
    justify-content: space-between;
  }

  .stat-card {
    flex: 1;
    min-width: 80px;
    padding: 0.75rem 1rem;
  }

  .filter-tabs {
    flex-direction: column;
  }

  .filter-tab {
    width: 100%;
  }

  .device-card {
    padding: 1.25rem;
  }

  .card-footer {
    flex-direction: column;
    align-items: flex-start;
  }
}
</style>
