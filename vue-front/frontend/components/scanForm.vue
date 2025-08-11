<template>
  <form class="scan-form" id="scanForm" @submit.prevent="startScan">
    
    <div class="form-group">
      <label for="networkRange">Network Range (CIDR)</label>
      <input
        type="text"
        id="networkRange"
        placeholder="ex, 192.168.1.0/24"
        required
        v-model="networkRange"
      />
    </div>

    <div class="form-group">
      <label for="scanType">Scan Type</label>
      <select id="scanType" v-model="scanType">
        <option value="full">Full Scan (SNMP + ARP) - Recommended</option>
        <option value="snmp">SNMP Only</option>
        <option value="arp">ARP Only</option>
      </select>
    </div>

    <div class="form-group communities-group">
      <label for="communities">SNMP Communities (comma-separated)</label>
      <input
        type="text"
        id="communities"
        placeholder="public,private,community"
        v-model="communities"
      />
      <small class="form-help">Used for SNMP discovery. Try common communities like 'public', 'private'</small>
    </div>

    <div style="display: grid; grid-template-columns: 1fr 1fr; gap: 20px">
      <div class="form-group">
        <label for="timeout">Timeout (seconds)</label>
        <input
          type="number"
          id="timeout"
          min="1"
          max="10"
          v-model.number="timeout"
        />
        <small class="form-help">Lower values = faster scan</small>
      </div>

      <div class="form-group">
        <label for="retries">Retries</label>
        <input
          type="number"
          id="retries"
          min="0"
          max="3"
          v-model.number="retries"
        />
        <small class="form-help">0 = fastest, 1-2 = more reliable</small>
      </div>
    </div>

    <button type="submit" class="btn" :disabled="loading">
      {{ loading ? 'Scanning...' : 'Start Network Scan' }}
    </button>
  </form>

  <scanTypes></scanTypes>

  <div class="loading" id="loading" v-if="loading">
    <div class="spinner"></div>
    <p>Scanning network... This may take a few minutes.</p>
  </div>

  <DeviceList :devices="devices"></DeviceList>
</template>

<script>
import scanTypes from "./scanTypes.vue";
import DeviceList from "./DeviceList.vue";

export default {
  components: {
    scanTypes,
    DeviceList,
  },
  data() {
    return {
      loading: false,
      networkRange: "",
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

      const communitiesArray = this.communities
        .split(",")
        .map((c) => c.trim())
        .filter((c) => c.length > 0);

      const payload = {
        network_range: this.networkRange,
        scan_type: this.scanType,
        communities: communitiesArray,
        timeout: this.timeout,
        retries: this.retries,
      };

      try {
        const response = await fetch(
          "http://localhost:8080/api/v1/network/full-scan",
          {
            method: "POST",
            headers: {
              "Content-Type": "application/json",
            },
            body: JSON.stringify(payload),
          }
        );

        if (!response.ok) {
          const errorData = await response.json();
          alert("Scan failed: " + (errorData.error || "Unknown error"));
          this.loading = false;
          return;
        }

        const data = await response.json();
        console.log("Scan result:", data);

        this.devices = data.topology?.devices || [];
      } catch (error) {
        alert("Network error or backend not reachable.");
        console.error(error);
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
  gap: 20px;
  margin-bottom: 30px;
}

.form-group {
  display: flex;
  flex-direction: column;
}

label {
  font-weight: 600;
  margin-bottom: 8px;
  color: #555;
}

.form-help {
  font-size: 0.85em;
  color: #666;
  margin-top: 4px;
  font-style: italic;
}

input,
select,
button {
  padding: 12px;
  border: 2px solid #e1e5e9;
  border-radius: 8px;
  font-size: 16px;
  transition: all 0.3s ease;
}

input:focus,
select:focus {
  outline: none;
  border-color: #515877;
  box-shadow: 0 0 0 3px rgba(102, 126, 234, 0.1);
}

.btn {
  background: linear-gradient(135deg, #515877 0%, #bbbabb 100%);
  color: white;
  border: none;
  cursor: pointer;
  font-weight: 600;
  text-transform: uppercase;
  letter-spacing: 1px;
  transition: all 0.3s ease;
}

.btn:hover {
  transform: translateY(-2px);
  box-shadow: 0 5px 15px rgba(102, 126, 234, 0.4);
}

.btn:disabled {
  opacity: 0.6;
  cursor: not-allowed;
  transform: none;
}

.loading {
  text-align: center;
  padding: 20px;
}

.spinner {
  border: 4px solid #f3f3f3;
  border-top: 4px solid #515877;
  border-radius: 50%;
  width: 40px;
  height: 40px;
  animation: spin 1s linear infinite;
  margin: 0 auto 10px;
}

@keyframes spin {
  0% {
    transform: rotate(0deg);
  }
  100% {
    transform: rotate(360deg);
  }
}
</style>
