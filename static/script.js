const API_BASE = "/api/v1";

document.getElementById("scanForm").addEventListener("submit", async (e) => {
  e.preventDefault();

  const networkRange = document.getElementById("networkRange").value;
  const communities = document
    .getElementById("communities")
    .value.split(",")
    .map((c) => c.trim());
  const timeout = parseInt(document.getElementById("timeout").value);
  const retries = parseInt(document.getElementById("retries").value);
  const scanType = document.getElementById("scanType").value;

  // Show loading
  document.getElementById("loading").style.display = "block";
  document.getElementById("results").style.display = "none";
  document.getElementById("error").style.display = "none";

  // Disable form
  const submitBtn = document.querySelector(".btn");
  submitBtn.disabled = true;
  submitBtn.textContent = "Scanning...";

  try {
    // Choose endpoint based on scan type
    let endpoint;
    let requestBody = {
      network_range: networkRange,
      timeout: timeout,
      retries: retries,
    };

    if (scanType === "arp") {
      endpoint = `${API_BASE}/network/scan/arp`;
    } else if (scanType === "snmp") {
      endpoint = `${API_BASE}/network/scan/snmp`;
      requestBody.communities = communities;
    } else {
      // Full scan (default)
      endpoint = `${API_BASE}/network/full-scan`;
      requestBody.communities = communities;
      requestBody.scan_type = "full";
    }

    const response = await fetch(endpoint, {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify(requestBody),
    });

    if (!response.ok) {
      throw new Error(`HTTP error! status: ${response.status}`);
    }

    const data = await response.json();
    displayResults(data);
  } catch (error) {
    showError(`Scan failed: ${error.message}`);
  } finally {
    // Hide loading and re-enable form
    document.getElementById("loading").style.display = "none";
    submitBtn.disabled = false;
    submitBtn.textContent = "Start Network Scan";
  }
});

function displayResults(data) {
  // Handle both new full scan format and legacy format
  const topology = data.topology || data.topology;
  const statistics = data.statistics || data.statistics;
  const scanInfo = data.scan_info;

  // Display statistics
  const statsContainer = document.getElementById("stats");
  statsContainer.innerHTML = `
    <div class="stat-card">
      <div class="stat-number">${statistics.total_devices}</div>
      <div>Total Devices</div>
    </div>
    <div class="stat-card">
      <div class="stat-number">${statistics.reachable_devices}</div>
      <div>Reachable</div>
    </div>
    ${
      statistics.snmp_devices !== undefined
        ? `
    <div class="stat-card">
      <div class="stat-number">${statistics.snmp_devices}</div>
      <div>SNMP Devices</div>
    </div>
    `
        : ""
    }
    ${
      statistics.arp_only_devices !== undefined
        ? `
    <div class="stat-card">
      <div class="stat-number">${statistics.arp_only_devices}</div>
      <div>ARP Only</div>
    </div>
    `
        : ""
    }
    ${
      statistics.devices_with_mac !== undefined
        ? `
    <div class="stat-card">
      <div class="stat-number">${statistics.devices_with_mac}</div>
      <div>With MAC</div>
    </div>
    `
        : ""
    }
    <div class="stat-card">
      <div class="stat-number">${(topology.scan_duration_ms / 1000).toFixed(
        1
      )}s</div>
      <div>Scan Duration</div>
    </div>
  `;

  // Add scan method info if available
  if (scanInfo) {
    const scanInfoDiv = document.createElement("div");
    scanInfoDiv.className = "scan-info";
    scanInfoDiv.innerHTML = `
      <div class="scan-method-badge ${scanInfo.scan_type}">
        ${getScanMethodLabel(scanInfo.scan_type)} Scan
      </div>
    `;
    statsContainer.appendChild(scanInfoDiv);
  }

  // Display devices
  const devicesContainer = document.getElementById("devicesGrid");
  devicesContainer.innerHTML = topology.devices
    .map((device) => createDeviceCard(device))
    .join("");

  document.getElementById("results").style.display = "block";
}

function createDeviceCard(device) {
  const scanMethodBadge = getScanMethodBadge(device.scan_method);

  return `
    <div class="device-card">
      <div class="device-header">
        <div class="device-ip">${device.ip}</div>
        <div class="device-badges">
          ${scanMethodBadge}
          <div class="device-status ${
            device.is_reachable ? "status-reachable" : "status-unreachable"
          }">
            ${device.is_reachable ? "Reachable" : "Unreachable"}
          </div>
        </div>
      </div>
      ${device.is_reachable ? createDeviceInfo(device) : ""}
    </div>
  `;
}

function createDeviceInfo(device) {
  return `
    <div class="device-info">
      ${
        device.hostname
          ? `
      <div class="info-item">
        <div class="info-label">Hostname</div>
        <div class="info-value">${device.hostname}</div>
      </div>
      `
          : ""
      }
      
      ${
        device.mac_address
          ? `
      <div class="info-item">
        <div class="info-label">MAC Address</div>
        <div class="info-value mac-address">${device.mac_address}</div>
      </div>
      `
          : ""
      }
      
      <div class="info-item">
        <div class="info-label">Vendor</div>
        <div class="info-value">${device.vendor || "Unknown"}</div>
      </div>
      
      ${
        device.description
          ? `
      <div class="info-item">
        <div class="info-label">Description</div>
        <div class="info-value description">${device.description}</div>
      </div>
      `
          : ""
      }
      
      ${
        device.uptime
          ? `
      <div class="info-item">
        <div class="info-label">Uptime</div>
        <div class="info-value">${device.uptime}</div>
      </div>
      `
          : ""
      }
      
      <div class="info-item">
        <div class="info-label">Response Time</div>
        <div class="info-value">${device.response_time_ms}ms</div>
      </div>
      
      ${
        device.location
          ? `
      <div class="info-item">
        <div class="info-label">Location</div>
        <div class="info-value">${device.location}</div>
      </div>
      `
          : ""
      }
    </div>
  `;
}

function getScanMethodBadge(scanMethod) {
  const badges = {
    SNMP: '<div class="scan-method-badge snmp">SNMP</div>',
    ARP: '<div class="scan-method-badge arp">ARP</div>',
    COMBINED: '<div class="scan-method-badge combined">SNMP+ARP</div>',
  };
  return badges[scanMethod] || "";
}

function getScanMethodLabel(scanType) {
  const labels = {
    full: "Full",
    snmp: "SNMP",
    arp: "ARP",
  };
  return labels[scanType] || "Unknown";
}

function showError(message) {
  const errorContainer = document.getElementById("error");
  errorContainer.textContent = message;
  errorContainer.style.display = "block";
}

// Handle scan type change to show/hide community field
document.getElementById("scanType").addEventListener("change", (e) => {
  const communitiesGroup = document.querySelector(".communities-group");
  if (e.target.value === "arp") {
    communitiesGroup.style.display = "none";
  } else {
    communitiesGroup.style.display = "block";
  }
});

// Auto-fill local network on page load
window.addEventListener("load", () => {
  // Try to detect local network range
  const networkInput = document.getElementById("networkRange");
  if (!networkInput.value) {
    networkInput.value = "192.168.1.0/24";
  }
});
