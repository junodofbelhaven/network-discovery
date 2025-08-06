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

  // Show loading
  document.getElementById("loading").style.display = "block";
  document.getElementById("results").style.display = "none";
  document.getElementById("error").style.display = "none";

  // Disable form
  const submitBtn = document.querySelector(".btn");
  submitBtn.disabled = true;
  submitBtn.textContent = "Scanning...";

  try {
    const response = await fetch(`${API_BASE}/network/scan`, {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify({
        network_range: networkRange,
        communities: communities,
        timeout: timeout,
        retries: retries,
      }),
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
  const { topology, statistics } = data;

  // Display statistics
  const statsContainer = document.getElementById("stats");
  statsContainer.innerHTML = `
                <div class="stat-card">
                    <div class="stat-number">${statistics.total_devices}</div>
                    <div>Total Devices</div>
                </div>
                <div class="stat-card">
                    <div class="stat-number">${
                      statistics.reachable_devices
                    }</div>
                    <div>Reachable</div>
                </div>
                <div class="stat-card">
                    <div class="stat-number">${
                      statistics.unreachable_devices
                    }</div>
                    <div>Unreachable</div>
                </div>
                <div class="stat-card">
                    <div class="stat-number">${(
                      topology.scan_duration_ms / 1000
                    ).toFixed(1)}s</div>
                    <div>Scan Duration</div>
                </div>
            `;

  // Display devices
  const devicesContainer = document.getElementById("devicesGrid");
  devicesContainer.innerHTML = topology.devices
    .map(
      (device) => `
                <div class="device-card">
                    <div class="device-header">
                        <div class="device-ip">${device.ip}</div>
                        <div class="device-status ${
                          device.is_reachable
                            ? "status-reachable"
                            : "status-unreachable"
                        }">
                            ${device.is_reachable ? "Reachable" : "Unreachable"}
                        </div>
                    </div>
                    ${
                      device.is_reachable
                        ? `
                        <div class="device-info">
                            <div class="info-item">
                                <div class="info-label">Hostname</div>
                                <div class="info-value">${
                                  device.hostname || "N/A"
                                }</div>
                            </div>
                            <div class="info-item">
                                <div class="info-label">Vendor</div>
                                <div class="info-value">${
                                  device.vendor || "Unknown"
                                }</div>
                            </div>
                            <div class="info-item">
                                <div class="info-label">Description</div>
                                <div class="info-value">${
                                  device.description || "N/A"
                                }</div>
                            </div>
                            <div class="info-item">
                                <div class="info-label">Uptime</div>
                                <div class="info-value">${
                                  device.uptime || "N/A"
                                }</div>
                            </div>
                            <div class="info-item">
                                <div class="info-label">Response Time</div>
                                <div class="info-value">${
                                  device.response_time_ms
                                }ms</div>
                            </div>
                            <div class="info-item">
                                <div class="info-label">Location</div>
                                <div class="info-value">${
                                  device.location || "N/A"
                                }</div>
                            </div>
                        </div>
                    `
                        : ""
                    }
                </div>
            `
    )
    .join("");

  document.getElementById("results").style.display = "block";
}

function showError(message) {
  const errorContainer = document.getElementById("error");
  errorContainer.textContent = message;
  errorContainer.style.display = "block";
}

// Auto-fill local network on page load
window.addEventListener("load", () => {
  // Try to detect local network range
  const networkInput = document.getElementById("networkRange");
  if (!networkInput.value) {
    networkInput.value = "192.168.1.0/24";
  }
});
