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
        <div class="input-helper">Enter specific device IP (SNMP scan only)</div>
      </div>
    </div>

    <!-- Scan Type Selection (only for full scan) -->
    <div v-if="scanMode === 'full-scan'" class="form-group-modern">
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
    <div
      v-if="
        (scanMode === 'full-scan' && (scanType === 'full' || scanType === 'snmp')) ||
        scanMode === 'single-ip'
      "
      class="form-group-modern"
    >
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
          {{ loading ? 'Scanning...' : 'Start Discovery' }}
        </span>
      </div>
      <div class="button-glow"></div>
    </button>
  </form>

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
    }
  },
  methods: {
    async startScan() {
      this.loading = true
      this.devices = []

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
  },
}
</script>

<style scoped>
/* ===== FORM CONTAINER ===== */
.scan-form-modern {
  background: linear-gradient(135deg, #1e293b 0%, #1a2331 100%);
  border-radius: 20px;
  padding: 2rem;
  box-shadow: 0 20px 40px rgba(0, 0, 0, 0.3);
  border: 1px solid rgba(59, 130, 246, 0.2);
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
  background: linear-gradient(90deg, #3b82f6, #10b981, #3b82f6);
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
  color: #f1f5f9;
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
  color: #cbd5e1;
  font-size: 0.95rem;
  margin: 0;
}

/* ===== FORM SECTIONS ===== */
.form-section {
  margin-bottom: 2rem;
}

.section-title {
  color: #3b82f6;
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
  background: rgba(59, 130, 246, 0.08);
  border: 2px solid rgba(59, 130, 246, 0.2);
  border-radius: 12px;
  transition: all 0.3s ease;
}

.radio-input:checked + .radio-content {
  background: rgba(59, 130, 246, 0.15);
  border-color: #3b82f6;
  box-shadow: 0 0 20px rgba(59, 130, 246, 0.3);
}

.radio-icon {
  font-size: 1.5rem;
}

.radio-text {
  flex: 1;
}

.radio-title {
  color: #f1f5f9;
  font-weight: 600;
  font-size: 0.95rem;
  margin-bottom: 0.25rem;
}

.radio-subtitle {
  color: #cbd5e1;
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
  color: #3b82f6;
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
  background: rgba(59, 130, 246, 0.08);
  border: 2px solid rgba(59, 130, 246, 0.2);
  border-radius: 12px;
  color: #f1f5f9;
  font-size: 1rem;
  transition: all 0.3s ease;
}

.modern-input-field:focus {
  outline: none;
  border-color: #3b82f6;
  background: rgba(59, 130, 246, 0.12);
  box-shadow: 0 0 20px rgba(59, 130, 246, 0.3);
}

.modern-input-field::placeholder {
  color: #94a3b8;
}

.input-helper {
  margin-top: 0.5rem;
  font-size: 0.8rem;
  color: #cbd5e1;
}

/* ===== SELECT FIELD ===== */
.select-container {
  position: relative;
}

.modern-select-field {
  width: 100%;
  padding: 1rem 3rem 1rem 1.25rem;
  background: rgba(59, 130, 246, 0.08);
  border: 2px solid rgba(59, 130, 246, 0.2);
  border-radius: 12px;
  color: #f1f5f9;
  font-size: 1rem;
  appearance: none;
  cursor: pointer;
  transition: all 0.3s ease;
}

.modern-select-field:focus {
  outline: none;
  border-color: #3b82f6;
  background: rgba(59, 130, 246, 0.12);
  box-shadow: 0 0 20px rgba(59, 130, 246, 0.3);
}

.modern-select-field option {
  background: #1e293b;
  color: #f1f5f9;
}

.select-arrow {
  position: absolute;
  right: 1.25rem;
  top: 50%;
  transform: translateY(-50%);
  pointer-events: none;
  color: #3b82f6;
  transition: transform 0.3s ease;
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
  color: #3b82f6;
  font-size: 0.95rem;
  font-weight: 600;
  cursor: pointer;
  padding: 0.5rem;
  transition: all 0.3s ease;
}

.advanced-toggle:hover {
  color: #10b981;
}

.toggle-icon {
  font-size: 0.8rem;
  transition: transform 0.3s ease;
}

.advanced-content {
  margin-top: 1rem;
  padding: 1rem;
  background: rgba(59, 130, 246, 0.08);
  border-radius: 12px;
  border: 1px solid rgba(59, 130, 246, 0.2);
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
  background: linear-gradient(135deg, #3b82f6 0%, #adb3cc 100%);
  border: none;
  border-radius: 12px;
  color: #f1f5f9;
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
  box-shadow: 0 10px 30px rgba(59, 130, 246, 0.4);
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
  border-top-color: #f1f5f9;
  border-radius: 50%;
  animation: spin 0.8s linear infinite;
}

@keyframes spin {
  0% {
    transform: rotate(0deg);
  }
  100% {
    transform: rotate(360deg);
  }
}

.button-glow {
  position: absolute;
  top: 50%;
  left: 50%;
  transform: translate(-50%, -50%);
  width: 100%;
  height: 100%;
  background: radial-gradient(circle, rgba(255, 255, 255, 0.2) 0%, transparent 70%);
  opacity: 0;
  transition: opacity 0.3s ease;
}

.scan-button:hover .button-glow {
  opacity: 1;
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
