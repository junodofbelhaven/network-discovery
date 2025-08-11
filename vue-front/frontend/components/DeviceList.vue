<template>
    <div class="device-list">
        <h2>Found Devices ({{ devices.length }})</h2>
    <div class="cards">
        <div v-for="device in devices" :key="device.ip" class="device-card">
        <p><strong>IP:</strong> {{ device.ip }}</p>
        <p v-if="device.mac_address"><strong>MAC:</strong> {{ device.mac_address }}</p>
        <p v-if="device.vendor"><strong>Vendor:</strong> {{ device.vendor }}</p>
        <p v-if="device.scan_method"><strong>Method:</strong> {{ device.scan_method }}</p>
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
};
</script>

<style scoped>
.device-list {
    margin-top: 20px;
}

.cards {
    display: flex;
    flex-wrap: wrap;
    gap: 16px;
}

.device-card {
    background: #f5f5f5;
    border-radius: 8px;
    padding: 16px;
    width: 220px;
    box-shadow: 0 2px 5px rgba(0,0,0,0.1);
    transition: box-shadow 0.3s ease, transform 0.3s ease;
    cursor: pointer;
}

.device-card:hover {
    box-shadow: 0 6px 15px rgba(0,0,0,0.25);
    transform: translateY(-5px);
}
</style>
