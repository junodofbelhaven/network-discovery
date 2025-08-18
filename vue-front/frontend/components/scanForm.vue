<!-- ===== scanForm.vue ===== -->
<template>
  <form @submit.prevent="startScan" class="scan-form-modern">
    <!-- Form Header -->
    <div class="form-header">
      <h3 class="form-title">
        <span class="form-icon">🔍</span>
        Network Scanner Configuration
      </h3>
      <p class="form-subtitle">Configure and initiate network discovery</p>
    </div>

    <!-- Scan Mode Selection -->
    <div class="form-section">
      <div class="section-title">Scan Mode</div>
      <div class="radio-group">
        <label class="radio-card">
          <input type="radio" v-model="scanMode" value="full-scan" class="radio-input" />
          <div class="radio-content">
            <span class="radio-icon">🌐</span>
            <div class="radio-text">
              <div class="radio-title">Network Range Scan</div>
              <div class="radio-subtitle">Scan entire network subnet</div>
            </div>
          </div>
        </label>

        <label class="radio-card">
          <input type="radio" v-model="scanMode" value="single-ip" class="radio-input" />
          <div class="radio-content">
            <span class="radio-icon">🎯</span>
            <div class="radio-text">
              <div class="radio-title">Single Device</div>
              <div class="radio-subtitle">Target specific IP address</div>
            </div>
          </div>
        </label>
      </div>
    </div>

    <!-- Network Range Input (for full scan) -->
    <div v-if="scanMode === 'full-scan'" class="form-group-modern">
      <label for="networkRange" class="input-label">
        <span class="label-icon">📍</span>
        Network Range (CIDR)
      </label>
      <div class="input-container">
        <input
          type="text"
          id="networkRange"
          v-model="networkRange"
          placeholder="192.168.1.0/24"
          class="modern-input-field"
          required
        />
        <div class="input-helper">Enter network in CIDR notation</div>
      </div>
    </div>

    <!-- Single IP Input (for single scan) -->
    <div v-if="scanMode === 'single-ip'" class="form-group-modern">
      <label for="targetIp" class="input-label">
        <span class="label-icon">🎯</span>
        Target IP Address
      </label>
      <div class="input-container">
        <input
          type="text"
          id="targetIp"
          v-model="targetIp"
          placeholder="192.168.1.100"
          class="modern-input-field"
          required
        />
        <div class="input-helper">Enter specific device IP</div>
      </div>
    </div>

    <!-- Scan Type Selection -->
    <div class="form-group-modern">
      <label for="scanType" class="input-label">
        <span class="label-icon">⚙️</span>
        Scan Protocol
      </label>
      <div class="select-container">
        <select id="scanType" v-model="scanType" class="modern-select-field">
          <option value="full">Full Scan (SNMP + ARP)</option>
          <option value="snmp">SNMP Protocol Only</option>
          <option value="arp">ARP Protocol Only</option>
        </select>
        <div class="select-arrow">
          <svg width="12" height="8" viewBox="0 0 12 8">
            <path d="M1 1l5 5 5-5" stroke="currentColor" stroke-width="2" fill="none" />
          </svg>
        </div>
      </div>
    </div>

    <!-- SNMP Communities -->
    <div v-if="scanType === 'full' || scanType === 'snmp'" class="form-group-modern">
      <label for="communities" class="input-label">
        <span class="label-icon">🔐</span>
        SNMP Communities
      </label>
      <div class="input-container">
        <input
          type="text"
          id="communities"
          v-model="communities"
          placeholder="public, private, community"
          class="modern-input-field"
        />
        <div class="input-helper">Comma-separated community strings</div>
      </div>
    </div>

    <!-- Advanced Settings -->
    <div class="advanced-settings">
      <button type="button" @click="showAdvanced = !showAdvanced" class="advanced-toggle">
        <span class="toggle-icon">{{ showAdvanced ? '▼' : '▶' }}</span>
        Advanced Settings
      </button>

      <transition name="slide">
        <div v-if="showAdvanced" class="advanced-content">
          <div class="settings-grid">
            <div class="form-group-modern">
              <label for="timeout" class="input-label">
                <span class="label-icon">⏱️</span>
                Timeout (seconds)
              </label>
              <input
                type="number"
                id="timeout"
                v-model="timeout"
                min="1"
                max="30"
                class="modern-input-field"
              />
            </div>

            <div class="form-group-modern">
              <label for="retries" class="input-label">
                <span class="label-icon">🔄</span>
                Retries
              </label>
              <input
                type="number"
                id="retries"
                v-model="retries"
                min="0"
                max="5"
                class="modern-input-field"
              />
            </div>
          </div>
        </div>
      </transition>
    </div>

    <!-- Submit Button -->
    <button type="submit" class="scan-button" :disabled="loading">
      <div class="button-content">
        <span v-if="!loading" class="button-icon">🚀</span>
        <div v-if="loading" class="button-spinner"></div>
        <span class="button-text">
          {{ loading ? 'Scanning Network...' : 'Start Discovery' }}
        </span>
      </div>
      <div class="button-glow"></div>
    </button>
  </form>

  <!-- Loading Overlay -->
  <transition name="fade">
    <div v-if="loading" class="scanning-overlay">
      <div class="scanning-content">
        <div class="scanning-animation">
          <div class="radar-sweep"></div>
          <div class="radar-dot"></div>
          <div class="radar-dot"></div>
          <div class="radar-dot"></div>
        </div>
        <h3 class="scanning-title">Discovering Network Devices</h3>
        <p class="scanning-subtitle">{{ scanningMessage }}</p>
        <div class="progress-bar-container">
          <div class="progress-bar-fill"></div>
        </div>
      </div>
    </div>
  </transition>

  <!-- Device List Component -->
  <device-list :devices="devices"></device-list>
</template>

<script>
import DeviceList from './DeviceList.vue'

export default {
  name: 'ScanForm',
  components: {
    'device-list': DeviceList,
  },
  data() {
    return {
      loading: false,
      scanMode: 'full-scan',
      networkRange: '192.168.1.0/24',
      targetIp: '',
      scanType: 'full',
      communities: 'public',
      timeout: 2,
      retries: 0,
      devices: [],
      showAdvanced: false,
      scanningMessage: 'Initializing scan...',
    }
  },
  methods: {
    async startScan() {
      this.loading = true
      this.devices = []
      this.updateScanningMessage()

      try {
        if (this.scanMode === 'full-scan') {
          const payload = {
            network_range: this.networkRange,
            scan_type: this.scanType,
            communities: this.communities
              .split(',')
              .map((c) => c.trim())
              .filter((c) => c.length > 0),
            timeout: this.timeout,
            retries: this.retries,
          }

          const res = await fetch('http://localhost:8080/api/v1/network/full-scan', {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify(payload),
          })

          if (!res.ok) throw new Error('Network scan failed')

          const data = await res.json()
          this.devices = data.topology?.devices || []
        } else if (this.scanMode === 'single-ip') {
          const communitiesArray = this.communities
            .split(',')
            .map((c) => c.trim())
            .filter((c) => c.length > 0)

          const query = communitiesArray.map((c) => `community=${encodeURIComponent(c)}`).join('&')

          const url = `http://localhost:8080/api/v1/device/${encodeURIComponent(
            this.targetIp.trim(),
          )}?${query}`

          const res = await fetch(url)
          if (!res.ok) throw new Error('Device scan failed')

          const data = await res.json()
          this.devices = [data.device]
        }

        this.$emit('scan-complete', this.devices)
      } catch (err) {
        console.error('Scan error:', err)
        alert(`Scan failed: ${err.message}`)
      } finally {
        this.loading = false
      }
    },

    updateScanningMessage() {
      const messages = [
        'Initializing scan...',
        'Discovering devices...',
        'Querying SNMP agents...',
        'Resolving MAC addresses...',
        'Gathering device information...',
        'Analyzing network topology...',
      ]

      let index = 0
      const interval = setInterval(() => {
        if (!this.loading) {
          clearInterval(interval)
          return
        }
        this.scanningMessage = messages[index % messages.length]
        index++
      }, 2000)
    },
  },
}
</script>

<style scoped>
/* ===== FORM CONTAINER ===== */
.scan-form-modern {
  background: linear-gradient(135deg, #232136 0%, #1a1928 100%);
  border-radius: 20px;
  padding: 2rem;
  box-shadow: 0 20px 40px rgba(0, 0, 0, 0.3);
  border: 1px solid rgba(184, 181, 255, 0.1);
  position: relative;
  overflow: hidden;
}

.scan-form-modern::before {
  content: '';
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  height: 3px;
  background: linear-gradient(90deg, #667eea, #764ba2, #667eea);
  background-size: 200% 100%;
  animation: gradient-shift 3s ease infinite;
}

@keyframes gradient-shift {
  0% {
    background-position: 0% 50%;
  }
  50% {
    background-position: 100% 50%;
  }
  100% {
    background-position: 0% 50%;
  }
}

/* ===== FORM HEADER ===== */
.form-header {
  margin-bottom: 2rem;
  text-align: center;
}

.form-title {
  color: #fffffe;
  font-size: 1.5rem;
  font-weight: 700;
  margin: 0 0 0.5rem 0;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 0.75rem;
}

.form-icon {
  font-size: 1.75rem;
}

.form-subtitle {
  color: #94a1b2;
  font-size: 0.95rem;
  margin: 0;
}

/* ===== FORM SECTIONS ===== */
.form-section {
  margin-bottom: 2rem;
}

.section-title {
  color: #b8b5ff;
  font-size: 0.9rem;
  font-weight: 600;
  text-transform: uppercase;
  letter-spacing: 1px;
  margin-bottom: 1rem;
}

/* ===== RADIO GROUP ===== */
.radio-group {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
  gap: 1rem;
}

.radio-card {
  position: relative;
  cursor: pointer;
}

.radio-input {
  position: absolute;
  opacity: 0;
}

.radio-content {
  display: flex;
  align-items: center;
  gap: 1rem;
  padding: 1rem;
  background: rgba(184, 181, 255, 0.05);
  border: 2px solid rgba(184, 181, 255, 0.2);
  border-radius: 12px;
  transition: all 0.3s ease;
}

.radio-input:checked + .radio-content {
  background: rgba(102, 126, 234, 0.2);
  border-color: #667eea;
  box-shadow: 0 0 20px rgba(102, 126, 234, 0.3);
}

.radio-icon {
  font-size: 1.5rem;
}

.radio-text {
  flex: 1;
}

.radio-title {
  color: #fffffe;
  font-weight: 600;
  font-size: 0.95rem;
  margin-bottom: 0.25rem;
}

.radio-subtitle {
  color: #94a1b2;
  font-size: 0.8rem;
}

/* ===== FORM GROUPS ===== */
.form-group-modern {
  margin-bottom: 1.5rem;
}

.input-label {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  color: #b8b5ff;
  font-weight: 600;
  font-size: 0.95rem;
  margin-bottom: 0.75rem;
}

.label-icon {
  font-size: 1.1rem;
}

.input-container {
  position: relative;
}

.modern-input-field {
  width: 100%;
  padding: 1rem 1.25rem;
  background: rgba(184, 181, 255, 0.05);
  border: 2px solid rgba(184, 181, 255, 0.2);
  border-radius: 12px;
  color: #fffffe;
  font-size: 1rem;
  transition: all 0.3s ease;
}

.modern-input-field:focus {
  outline: none;
  border-color: #667eea;
  background: rgba(102, 126, 234, 0.1);
  box-shadow: 0 0 20px rgba(102, 126, 234, 0.3);
}

.modern-input-field::placeholder {
  color: #94a1b2;
}

.input-helper {
  margin-top: 0.5rem;
  font-size: 0.8rem;
  color: #94a1b2;
}

/* ===== SELECT FIELD ===== */
.select-container {
  position: relative;
}

.modern-select-field {
  width: 100%;
  padding: 1rem 3rem 1rem 1.25rem;
  background: rgba(184, 181, 255, 0.05);
  border: 2px solid rgba(184, 181, 255, 0.2);
  border-radius: 12px;
  color: #fffffe;
  font-size: 1rem;
  appearance: none;
  cursor: pointer;
  transition: all 0.3s ease;
}

.modern-select-field:focus {
  outline: none;
  border-color: #667eea;
  background: rgba(102, 126, 234, 0.1);
  box-shadow: 0 0 20px rgba(102, 126, 234, 0.3);
}

.select-arrow {
  position: absolute;
  right: 1.25rem;
  top: 50%;
  transform: translateY(-50%);
  pointer-events: none;
  color: #b8b5ff;
  transition: transform 0.3s ease;
}

.modern-select-field:focus + .select-arrow {
  transform: translateY(-50%) rotate(180deg);
}

/* ===== ADVANCED SETTINGS ===== */
.advanced-settings {
  margin: 2rem 0;
}

.advanced-toggle {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  background: transparent;
  border: none;
  color: #b8b5ff;
  font-size: 0.95rem;
  font-weight: 600;
  cursor: pointer;
  padding: 0.5rem;
  transition: all 0.3s ease;
}

.advanced-toggle:hover {
  color: #667eea;
}

.toggle-icon {
  font-size: 0.8rem;
  transition: transform 0.3s ease;
}

.advanced-content {
  margin-top: 1rem;
  padding: 1rem;
  background: rgba(184, 181, 255, 0.05);
  border-radius: 12px;
  border: 1px solid rgba(184, 181, 255, 0.1);
}

.settings-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(150px, 1fr));
  gap: 1rem;
}

/* ===== SCAN BUTTON ===== */
.scan-button {
  width: 100%;
  padding: 1.25rem 2rem;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  border: none;
  border-radius: 12px;
  color: #fffffe;
  font-size: 1.1rem;
  font-weight: 700;
  cursor: pointer;
  position: relative;
  overflow: hidden;
  transition: all 0.3s ease;
  text-transform: uppercase;
  letter-spacing: 1px;
  margin-top: 2rem;
}

.scan-button:hover:not(:disabled) {
  transform: translateY(-2px);
  box-shadow: 0 10px 30px rgba(102, 126, 234, 0.4);
}

.scan-button:disabled {
  opacity: 0.7;
  cursor: not-allowed;
  transform: none;
}

.button-content {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 0.75rem;
  position: relative;
  z-index: 1;
}

.button-icon {
  font-size: 1.3rem;
}

.button-spinner {
  width: 20px;
  height: 20px;
  border: 3px solid rgba(255, 255, 255, 0.3);
  border-top-color: #fffffe;
  border-radius: 50%;
  animation: spin 0.8s linear infinite;
}

.button-glow {
  position: absolute;
  top: 50%;
  left: 50%;
  transform: translate(-50%, -50%);
  width: 100%;
  height: 100%;
  background: radial-gradient(circle, rgba(255, 255, 255, 0.3) 0%, transparent 70%);
  opacity: 0;
  transition: opacity 0.3s ease;
}

.scan-button:hover .button-glow {
  opacity: 1;
}

/* ===== SCANNING OVERLAY ===== */
.scanning-overlay {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(15, 14, 23, 0.95);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 1000;
  backdrop-filter: blur(10px);
}

.scanning-content {
  text-align: center;
  max-width: 400px;
}

.scanning-animation {
  width: 120px;
  height: 120px;
  margin: 0 auto 2rem;
  position: relative;
}

.radar-sweep {
  width: 100%;
  height: 100%;
  border: 3px solid rgba(102, 126, 234, 0.3);
  border-radius: 50%;
  position: relative;
}

.radar-sweep::before {
  content: '';
  position: absolute;
  top: 50%;
  left: 50%;
  width: 100%;
  height: 100%;
  background: conic-gradient(from 0deg, transparent 0deg, #667eea 30deg, transparent 60deg);
  border-radius: 50%;
  transform: translate(-50%, -50%);
  animation: radar-sweep 2s linear infinite;
}

@keyframes radar-sweep {
  0% {
    transform: translate(-50%, -50%) rotate(0deg);
  }
  100% {
    transform: translate(-50%, -50%) rotate(360deg);
  }
}

.radar-dot {
  position: absolute;
  width: 8px;
  height: 8px;
  background: #48bb78;
  border-radius: 50%;
  animation: radar-ping 2s ease-out infinite;
}

.radar-dot:nth-child(2) {
  top: 20%;
  left: 30%;
  animation-delay: 0.3s;
}

.radar-dot:nth-child(3) {
  top: 60%;
  left: 70%;
  animation-delay: 0.6s;
}

.radar-dot:nth-child(4) {
  top: 40%;
  left: 80%;
  animation-delay: 0.9s;
}

@keyframes radar-ping {
  0% {
    opacity: 0;
    transform: scale(0);
  }
  50% {
    opacity: 1;
    transform: scale(1);
  }
  100% {
    opacity: 0;
    transform: scale(1.5);
  }
}

.scanning-title {
  color: #fffffe;
  font-size: 1.5rem;
  font-weight: 700;
  margin-bottom: 0.5rem;
}

.scanning-subtitle {
  color: #94a1b2;
  font-size: 0.95rem;
  margin-bottom: 2rem;
}

.progress-bar-container {
  width: 100%;
  height: 4px;
  background: rgba(184, 181, 255, 0.2);
  border-radius: 2px;
  overflow: hidden;
}

.progress-bar-fill {
  height: 100%;
  background: linear-gradient(90deg, #667eea, #764ba2);
  animation: progress-animation 2s ease-in-out infinite;
}

@keyframes progress-animation {
  0% {
    width: 0%;
  }
  50% {
    width: 100%;
  }
  100% {
    width: 0%;
  }
}

/* ===== ANIMATIONS ===== */
@keyframes spin {
  0% {
    transform: rotate(0deg);
  }
  100% {
    transform: rotate(360deg);
  }
}

.slide-enter-active,
.slide-leave-active {
  transition: all 0.3s ease;
}

.slide-enter-from {
  opacity: 0;
  transform: translateY(-10px);
}

.slide-leave-to {
  opacity: 0;
  transform: translateY(-10px);
}

.fade-enter-active,
.fade-leave-active {
  transition: opacity 0.3s ease;
}

.fade-enter-from,
.fade-leave-to {
  opacity: 0;
}

/* ===== RESPONSIVE ===== */
@media (max-width: 768px) {
  .scan-form-modern {
    padding: 1.5rem;
  }

  .form-title {
    font-size: 1.25rem;
  }

  .radio-group {
    grid-template-columns: 1fr;
  }

  .settings-grid {
    grid-template-columns: 1fr;
  }
}
</style>
