<template>
  <div class="device-list" v-if="devices.length > 0">
    <div class="results-header">
      <div class="results-title">
        <span class="results-icon">📊</span>
        <h2>Network Scan Results</h2>
        <div class="device-count">{{ devices.length }} {{ devices.length === 1 ? 'device' : 'devices' }} found</div>
      </div>
      
      <!-- Summary Stats -->
      <div class="scan-summary">
        <div class="stat-item snmp-stat" v-if="snmpDevices.length > 0">
          <div class="stat-number">{{ snmpDevices.length }}</div>
          <div class="stat-label">SNMP Devices</div>
        </div>
        <div class="stat-item arp-stat" v-if="arpDevices.length > 0">
          <div class="stat-number">{{ arpDevices.length }}</div>
          <div class="stat-label">ARP Devices</div>
        </div>
      </div>
    </div>

    <!-- Device List -->
    <div class="devices-container">
      <!-- SNMP Devices Section -->
      <div v-if="snmpDevices.length > 0" class="device-section snmp-section">
        <div class="section-header snmp-header">
          <div class="section-title">
            <span class="section-icon">📡</span>
            <span>SNMP Devices</span>
            <span class="protocol-badge snmp-badge">SNMP Protocol</span>
          </div>
          <div class="section-count">{{ snmpDevices.length }} devices</div>
        </div>
        
        <div class="device-rows">
          <div 
            v-for="device in snmpDevices" 
            :key="device.ip" 
            class="device-row snmp-row"
          >
            <div class="device-row-main">
              <div class="device-primary-info">
                <div class="device-ip">
                  <span class="ip-icon">🌐</span>
                  {{ device.ip }}
                </div>
                <div class="device-hostname" v-if="device.hostname">
                  <span class="hostname-icon">🏷️</span>
                  {{ device.hostname }}
                </div>
                <div class="device-status online">
                  <span class="status-dot"></span>
                  Online
                </div>
              </div>
              
              <div class="device-secondary-info">
                <div class="info-item" v-if="device.vendor">
                  <span class="info-label">Vendor:</span>
                  <span class="info-value">{{ device.vendor }}</span>
                </div>
                <div class="info-item" v-if="device.response_time">
                  <span class="info-label">Response:</span>
                  <span class="info-value response-time">{{ device.response_time }}ms</span>
                </div>
                <div class="info-item" v-if="device.uptime">
                  <span class="info-label">Uptime:</span>
                  <span class="info-value uptime">{{ device.uptime }}</span>
                </div>
              </div>
            </div>
            
            <div class="device-row-details" v-if="device.description || device.mac_address">
              <div class="detail-item" v-if="device.description">
                <span class="detail-icon">📋</span>
                <span class="detail-text">{{ device.description }}</span>
              </div>
              <div class="detail-item" v-if="device.mac_address">
                <span class="detail-icon">🔗</span>
                <span class="detail-text">MAC: {{ device.mac_address }}</span>
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- ARP Devices Section -->
      <div v-if="arpDevices.length > 0" class="device-section arp-section">
        <div class="section-header arp-header">
          <div class="section-title">
            <span class="section-icon">🌐</span>
            <span>ARP Devices</span>
            <span class="protocol-badge arp-badge">ARP Protocol</span>
          </div>
          <div class="section-count">{{ arpDevices.length }} devices</div>
        </div>
        
        <div class="device-rows">
          <div 
            v-for="device in arpDevices" 
            :key="device.ip" 
            class="device-row arp-row"
          >
            <div class="device-row-main">
              <div class="device-primary-info">
                <div class="device-ip">
                  <span class="ip-icon">🌐</span>
                  {{ device.ip }}
                </div>
                <div class="device-status online">
                  <span class="status-dot"></span>
                  Active
                </div>
              </div>
              
              <div class="device-secondary-info">
                <div class="info-item" v-if="device.vendor">
                  <span class="info-label">Vendor:</span>
                  <span class="info-value">{{ device.vendor }}</span>
                </div>
                <div class="info-item" v-if="device.mac_address">
                  <span class="info-label">MAC:</span>
                  <span class="info-value mac-address">{{ device.mac_address }}</span>
                </div>
              </div>
            </div>
            
            
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
export default {
  name: "DeviceList",
  props: {
    devices: {
      type: Array,
      required: true,
    },
  },
  computed: {
    snmpDevices() {
      return this.devices.filter(device => device.scan_method === 'SNMP');
    },
    arpDevices() {
      return this.devices.filter(device => device.scan_method === 'ARP');
    }
  }
};
</script>

<style scoped>
.device-list {
  margin-top: 30px;
  max-width: 1400px;
  margin-left: auto;
  margin-right: auto;
  padding: 0 20px;
}

.results-header {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  margin-bottom: 30px;
  padding: 25px;
  background: rgba(64, 56, 74, 0.85); 
  backdrop-filter: blur(10px);
  border-radius: 16px;
  box-shadow: 0 4px 20px rgba(100, 80, 130, 0.4);
  border: 1px solid rgba(120, 90, 150, 0.4);
}

.results-title {
  display: flex;
  align-items: center;
  gap: 12px;
}

.results-icon {
  font-size: 1.5rem;
}

.results-title h2 {
  color: #e0d8ff; 
  font-size: 1.5rem;
  font-weight: 700;
  margin: 0;
}

.device-count {
  background: linear-gradient(135deg, #8a6eff, #764ba2);
  color: white;
  padding: 6px 12px;
  border-radius: 20px;
  font-size: 0.85rem;
  font-weight: 600;
}

.scan-summary {
  display: flex;
  gap: 20px;
}

.stat-item {
  text-align: center;
  padding: 15px 20px;
  border-radius: 12px;
  min-width: 80px;
  transition: transform 0.3s ease;
}

.stat-item:hover {
  transform: translateY(-2px);
}

.stat-item.snmp-stat {
  background: linear-gradient(135deg, #667eea, #764ba2); 
  color: white;
}

.stat-item.arp-stat {
  background: linear-gradient(135deg, #48bb78, #2f855a); 
  color: white;
}

.stat-number {
  font-size: 1.8rem;
  font-weight: 700;
  margin-bottom: 4px;
}

.stat-label {
  font-size: 0.8rem;
  font-weight: 500;
  opacity: 0.9;
}

.devices-container {
  display: flex;
  flex-direction: column;
  gap: 30px;
}

.device-section {
  background: rgba(64, 56, 74, 0.85); 
  backdrop-filter: blur(10px);
  border-radius: 16px;
  overflow: hidden;
  box-shadow: 0 4px 20px rgba(100, 80, 130, 0.2);
  border: 1px solid rgba(120, 90, 150, 0.3);
}

.section-header {
  padding: 20px 25px;
  display: flex;
  justify-content: space-between;
  align-items: center;
  border-bottom: 2px solid rgba(0, 0, 0, 0.05);
}

.snmp-header {
  background: linear-gradient(135deg, rgba(102, 126, 234, 0.2), rgba(118, 75, 162, 0.1));
  border-left: 4px solid #8a6eff;
}

.arp-header {
  background: linear-gradient(135deg, rgba(72, 187, 120, 0.15), rgba(47, 133, 90, 0.08));
  border-left: 4px solid #48bb78;
}

.section-title {
  display: flex;
  align-items: center;
  gap: 12px;
  font-size: 1.2rem;
  font-weight: 600;
  color: #dcd6ff;
}

.section-icon {
  font-size: 1.3rem;
}

.protocol-badge {
  padding: 4px 10px;
  border-radius: 12px;
  font-size: 0.75rem;
  font-weight: 600;
  text-transform: uppercase;
  letter-spacing: 0.5px;
}

.snmp-badge {
  background: rgba(138, 110, 255, 0.25);
  color: #6f57c7;
}

.arp-badge {
  background: rgba(72, 187, 120, 0.25);
  color: #2f855a;
}

.section-count {
  color: #bdb6e8;
  font-size: 0.9rem;
  font-weight: 500;
}

.device-rows {
  display: flex;
  flex-direction: column;
}

.device-row {
  position: relative;
  padding: 20px 25px;
  border-bottom: 1px solid rgba(0, 0, 0, 0.05);
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
  cursor: pointer;
}

.device-row:last-child {
  border-bottom: none;
}

.device-row:hover {
  background: rgba(138, 110, 255, 0.12);
  transform: translateX(5px);
}

.snmp-row {
  border-left: 3px solid transparent;
}

.snmp-row:hover {
  border-left-color: #8a6eff;
  background: rgba(138, 110, 255, 0.1);
}

.arp-row {
  border-left: 3px solid transparent;
}

.arp-row:hover {
  border-left-color: #48bb78;
  background: rgba(72, 187, 120, 0.1);
}

.device-row-main {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  margin-bottom: 12px;
}

.device-primary-info {
  display: flex;
  align-items: center;
  gap: 20px;
  flex-wrap: wrap;
}

.device-ip {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 1.1rem;
  font-weight: 600;
  color: #dcd6ff;
}

.ip-icon {
  font-size: 1rem;
}

.device-hostname {
  display: flex;
  align-items: center;
  gap: 6px;
  color: #b1aaff;
  font-weight: 500;
}

.hostname-icon {
  font-size: 0.9rem;
}

.device-status {
  display: flex;
  align-items: center;
  gap: 6px;
  padding: 4px 10px;
  border-radius: 12px;
  font-size: 0.8rem;
  font-weight: 600;
  text-transform: uppercase;
}

.device-status.online {
  background: rgba(72, 187, 120, 0.15);
  color: #8acc9b;
}

.status-dot {
  width: 6px;
  height: 6px;
  border-radius: 50%;
  background: #48bb78;
  animation: pulse-dot 2s infinite;
}

.device-secondary-info {
  display: flex;
  gap: 20px;
  flex-wrap: wrap;
  align-items: center;
}

.info-item {
  display: flex;
  align-items: center;
  gap: 6px;
  font-size: 0.9rem;
}

.info-label {
  color: #b0a8d6;
  font-weight: 500;
}

.info-value {
  color: #dcd6ff;
  font-weight: 600;
}

.response-time {
  color: #66bb6a;
}

.uptime {
  color: #7c5ec4;
}

.mac-address {
  font-family: 'Monaco', 'Menlo', 'Courier New', monospace;
  font-size: 0.85rem;
}

.device-row-details {
  display: flex;
  flex-direction: column;
  gap: 8px;
  margin-top: 8px;
  padding-top: 12px;
  border-top: 1px solid rgba(0, 0, 0, 0.05);
}

.detail-item {
  display: flex;
  align-items: flex-start;
  gap: 8px;
  font-size: 0.9rem;
  color: #b1aaff;
}

.detail-icon {
  font-size: 0.8rem;
  margin-top: 2px;
}

.detail-text {
  flex: 1;
}

.device-protocol-indicator {
  position: absolute;
  top: 15px;
  right: 20px;
  padding: 6px 12px;
  border-radius: 8px;
  font-size: 0.75rem;
  font-weight: 700;
  text-transform: uppercase;
  letter-spacing: 0.5px;
}

.snmp-indicator {
  background: linear-gradient(135deg, #8a6eff, #764ba2);
  color: #e0d8ff;
}

.arp-indicator {
  background: linear-gradient(135deg, #48bb78, #2f855a);
  color: white;
}

@keyframes pulse-dot {
  0%, 100% {
    opacity: 1;
    transform: scale(1);
  }
  50% {
    opacity: 0.5;
    transform: scale(1.2);
  }
}

@media (max-width: 768px) {
  .device-list {
    padding: 0 15px;
  }
  
  .results-header {
    flex-direction: column;
    gap: 20px;
    align-items: stretch;
    padding: 20px;
  }
  
  .results-title {
    justify-content: center;
    flex-wrap: wrap;
  }
  
  .scan-summary {
    justify-content: center;
  }
  
  .device-row {
    padding: 15px 20px;
  }
  
  .device-row-main {
    flex-direction: column;
    gap: 12px;
  }
  
  .device-primary-info {
    gap: 12px;
  }
  
  .device-secondary-info {
    gap: 12px;
  }
  
  .device-protocol-indicator {
    position: static;
    align-self: flex-start;
    margin-top: 10px;
  }
}

@media (max-width: 480px) {
  .results-title h2 {
    font-size: 1.2rem;
  }
  
  .device-ip {
    font-size: 1rem;
  }
  
  .section-title {
    font-size: 1rem;
  }
  
  .stat-item {
    padding: 10px 15px;
    min-width: 70px;
  }
  
  .stat-number {
    font-size: 1.4rem;
  }
}
</style>
