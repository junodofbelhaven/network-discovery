<template>
    <div class="scan-info">
        <h3>Scan Types Explained:</h3>
        <div class="scan-methods">
            <div class="scan-method">
            <strong>🔍 Full Scan (SNMP + ARP):</strong>
                <p>
                Combines both SNMP and ARP scanning for the most comprehensive
                device discovery. Provides detailed device information plus
                broad network coverage.
                </p>
            </div>
            <div class="scan-method">
            <strong>📡 SNMP Only:</strong>
                <p>
                Discovers only SNMP-enabled devices. Provides detailed
                information like hostname, description, vendor, and uptime.
                </p>
            </div>
            <div class="scan-method">
                <strong>🌐 ARP Only:</strong>
                <p>
                Discovers all IP-enabled devices on the network segment.
                Provides IP and MAC addresses with vendor identification.
                </p>
            </div>
        </div>
    </div>
</template>


<script>
</script>


<style scoped>
.scan-info {
    margin-top: 20px;
    padding: 20px;
    background: #f8f9fa;
    border-radius: 8px;
    border-left: 4px solid #515877;
}

.scan-info h3 {
    color: #515877;
    margin-bottom: 15px;
}

.scan-methods {
    display: grid;
    gap: 12px;
}

.scan-method {
    padding: 12px;
    background: white;
    border-radius: 6px;
    border: 1px solid #e1e5e9;
}

.scan-method strong {
    color: #515877;
    display: block;
    margin-bottom: 5px;
}

.scan-method p {
    font-size: 0.9em;
    color: #666;
    margin: 0;
}
</style>