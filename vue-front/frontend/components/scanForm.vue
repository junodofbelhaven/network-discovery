
<template>
  <form class="scan-form" @submit.prevent="startScan">
    <!-- Scan Mode -->
    <div class="form-group">
      <label for="scanMode" class="modern-label">
        <span class="label-icon">🎯</span>
        Scan Mode
      </label>
      <div class="input-wrapper">
        <select id="scanMode" v-model="scanMode" class="modern-select">
          <option value="full-scan">Full Network Scan</option>
          <option value="single-ip">Single IP Scan</option>
        </select>
        <div class="select-arrow">
          <svg width="12" height="8" viewBox="0 0 12 8">
            <path d="M1 1l5 5 5-5" stroke="currentColor" stroke-width="2" fill="none"/>
          </svg>
        </div>
      </div>
    </div>

    <!-- Network Range -->
    <div class="form-group" v-if="scanMode === 'full-scan'">
      <label for="networkRange" class="modern-label">
        <span class="label-icon">🌐</span>
        Network Range (CIDR)
      </label>
      <div class="input-wrapper">
        <input
          type="text"
          id="networkRange"
          placeholder="ex, 192.168.1.0/24"
          v-model="networkRange"
          class="modern-input"
        />
        <div class="input-focus-line"></div>
      </div>
    </div>

    <!-- Single IP -->
    <div class="form-group" v-if="scanMode === 'single-ip'">
      <label for="targetIp" class="modern-label">
        <span class="label-icon">🎯</span>
        Target IP
      </label>
      <div class="input-wrapper">
        <input
          type="text"
          id="targetIp"
          placeholder="ex, 192.168.1.10"
          v-model="targetIp"
          class="modern-input"
        />
        <div class="input-focus-line"></div>
      </div>
    </div>

    <!-- Scan Type (only for full scan) -->
    <div class="form-group" v-if="scanMode === 'full-scan'">
      <label for="scanType" class="modern-label">
        <span class="label-icon">⚡</span>
        Scan Type
      </label>
      <div class="input-wrapper">
        <select id="scanType" v-model="scanType" class="modern-select">
          <option value="full">Full Scan (SNMP + ARP)</option>
          <option value="snmp">SNMP Only</option>
          <option value="arp">ARP Only</option>
        </select>
        <div class="select-arrow">
          <svg width="12" height="8" viewBox="0 0 12 8">
            <path d="M1 1l5 5 5-5" stroke="currentColor" stroke-width="2" fill="none"/>
          </svg>
        </div>
      </div>
    </div>

    <!-- SNMP Communities -->
    <div class="form-group" v-if="scanMode === 'full-scan'">
      <label for="communities" class="modern-label">
        <span class="label-icon">🔐</span>
        SNMP Communities
      </label>
      <div class="input-wrapper">
        <input
          type="text"
          id="communities"
          placeholder="public,private,community"
          v-model="communities"
          class="modern-input"
        />
        <div class="input-focus-line"></div>
      </div>
    </div>

    <button type="submit" class="modern-btn" :disabled="loading">
      <span class="btn-content">
        <span class="btn-icon" v-if="!loading">🚀</span>
        <span class="btn-text">{{ loading ? 'Scanning...' : 'Start Scan' }}</span>
      </span>

    </button>
  </form>

  <div v-if="loading" class="loading-container">
    <div class="modern-spinner">
      <div class="spinner-ring"></div>
      <div class="spinner-ring"></div>
      <div class="spinner-ring"></div>
    </div>
    <p class="loading-text">Scanning network...</p>
    <div class="loading-progress">
      <div class="progress-bar"></div>
    </div>
  </div>

  <device-list :devices="devices"></device-list>
</template>

<script>
import DeviceList from "./DeviceList.vue";

export default {
  components: {
    "device-list": DeviceList,
  },
  data() {
    return {
      loading: false,
      scanMode: "full-scan",
      networkRange: "",
      targetIp: "",
      scanType: "full",
      communities: "public",
      timeout: 2,
      retries: 0,
      devices: [],
    };
  },
  methods: {
    async startScan() {
      this.loading = true;
      try {
        if (this.scanMode === "full-scan") {
          const payload = {
            network_range: this.networkRange,
            scan_type: this.scanType,
            communities: this.communities
              .split(",")
              .map((c) => c.trim())
              .filter((c) => c.length > 0),
            timeout: this.timeout,
            retries: this.retries,
          };

          const res = await fetch("http://localhost:8080/api/v1/network/full-scan", {
            method: "POST",
            headers: { "Content-Type": "application/json" },
            body: JSON.stringify(payload),
          });

          if (!res.ok) throw new Error("Full scan failed");
          const data = await res.json();
          this.devices = data.topology?.devices || [];
        } else if (this.scanMode === "single-ip") {
          const communitiesArray = this.communities
            .split(",")
            .map((c) => c.trim())
            .filter((c) => c.length > 0);

          const query = communitiesArray
            .map((c) => `community=${encodeURIComponent(c)}`)
            .join("&");

          const url = `http://localhost:8080/api/v1/device/${encodeURIComponent(
            this.targetIp.trim()
          )}?${query}`;

          const res = await fetch(url);

          if (!res.ok) throw new Error("Single IP scan failed");

          const data = await res.json();

          this.devices = [data.device];
        }
      } catch (err) {
        alert(err.message);
        console.error(err);
      } finally {
        this.loading = false;
      }
    },
  },
};
</script>

<style scoped>
.scan-form {
  display: grid;
  gap: 25px;
  margin-bottom: 40px;
}

.form-group {
  display: flex;
  flex-direction: column;
  position: relative;
}

.modern-label {
  font-weight: 600;
  margin-bottom: 12px;
  color: #4a5568;
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 0.95rem;
  letter-spacing: 0.3px;
}

.label-icon {
  font-size: 1.1rem;
}

.input-wrapper {
  position: relative;
}

.modern-input,
.modern-select {
  width: 100%;
  padding: 16px 18px;
  border: 2px solid #acbed8; 
  border-radius: 12px;
  font-size: 16px;
  background: #e8ebf7; 
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
  box-shadow: 0 1px 3px rgba(172, 190, 216, 0.2);
  font-family: inherit;
  color: #2d3748;
}

.modern-select {
  appearance: none;
  background-image: none;
  cursor: pointer;
  padding-right: 45px;
}

.select-arrow {
  position: absolute;
  right: 15px;
  top: 50%;
  transform: translateY(-50%);
  pointer-events: none;
  color: #4a5568;
  transition: transform 0.3s ease;
}

.modern-select:focus + .select-arrow {
  transform: translateY(-50%) rotate(180deg);
}

.modern-input:focus,
.modern-select:focus {
  outline: none;
  border-color: #f2d398; 
  box-shadow: 
    0 0 0 3px rgba(242, 211, 152, 0.3),
    0 4px 12px rgba(172, 190, 216, 0.3);
  transform: translateY(-1px);
}

.modern-btn {
  background: linear-gradient(135deg, #f2d398 0%, #acbed8 100%);
  color: #2d3748;
  border: none;
  border-radius: 12px;
  padding: 18px 32px;
  font-size: 16px;
  font-weight: 600;
  cursor: pointer;
  position: relative;
  overflow: hidden;
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
  text-transform: uppercase;
  letter-spacing: 1px;
  box-shadow: 0 4px 15px rgba(172, 190, 216, 0.5);
  min-height: 60px;
}

.btn-content {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 8px;
  position: relative;
  z-index: 2;
}

.btn-icon {
  font-size: 1.4rem;
}



.modern-btn:hover {
  transform: translateY(-2px);
  box-shadow: 0 8px 25px rgba(172, 190, 216, 0.6);
}



.modern-btn:disabled {
  opacity: 0.7;
  cursor: not-allowed;
  transform: none;
  box-shadow: 0 2px 8px rgba(172, 190, 216, 0.2);
}

.loading-container {
  display: flex;
  flex-direction: column;
  align-items: center;
  padding: 40px 20px;
  background: #e8ebf7; 
  border-radius: 16px;
  margin-bottom: 30px;
  box-shadow: 0 4px 20px rgba(172, 190, 216, 0.2);
}

.modern-spinner {
  position: relative;
  width: 60px;
  height: 60px;
  margin-bottom: 20px;
}

.spinner-ring {
  position: absolute;
  width: 60px;
  height: 60px;
  border: 3px solid transparent;
  border-top-color: #acbed8; 
  border-radius: 50%;
  animation: spin 1.5s linear infinite;
}

.spinner-ring:nth-child(2) {
  width: 40px;
  height: 40px;
  top: 10px;
  left: 10px;
  border-top-color: #f2d398; 
  animation-duration: 1s;
  animation-direction: reverse;
}

.spinner-ring:nth-child(3) {
  width: 20px;
  height: 20px;
  top: 20px;
  left: 20px;
  border-top-color: #acbed8;
  animation-duration: 0.8s;
}

.loading-text {
  color: #4a5568;
  font-weight: 500;
  margin-bottom: 20px;
  font-size: 1.1rem;
}

.loading-progress {
  width: 200px;
  height: 4px;
  background: #acbed8;
  border-radius: 2px;
  overflow: hidden;
}

.progress-bar {
  width: 100%;
  height: 100%;
  background: linear-gradient(90deg, #acbed8, #f2d398);
  border-radius: 2px;
  animation: loading-progress 2s ease-in-out infinite;
}

@keyframes spin {
  0% { transform: rotate(0deg); }
  100% { transform: rotate(360deg); }
}

@keyframes loading-progress {
  0% { transform: translateX(-100%); }
  50% { transform: translateX(0%); }
  100% { transform: translateX(100%); }
}

@media (max-width: 768px) {
  .scan-form {
    gap: 20px;
  }
  
  .modern-input,
  .modern-select {
    padding: 14px 16px;
    font-size: 16px;
  }
  
  .modern-btn {
    padding: 16px 28px;
    font-size: 15px;
  }
}
</style>
