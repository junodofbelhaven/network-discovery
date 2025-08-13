<template>
  <div class="scan-info">
    <div class="info-header">
      <h3 class="info-title">
        <span class="title-icon">📚</span>
        Scan Types Explained
      </h3>
      <div class="header-decoration"></div>
    </div>
    
    <div class="scan-methods">
      <div class="scan-method full-scan">
        <div class="method-header">
          <div class="method-icon">🔍</div>
          <div class="method-title">
            <strong>Full Scan (SNMP + ARP)</strong>
            <span class="method-badge premium">Premium</span>
          </div>
        </div>
        <p class="method-description">
          Combines both SNMP and ARP scanning for the most comprehensive
          device discovery. Provides detailed device information plus
          broad network coverage with advanced analytics.
        </p>
        <div class="method-features">
          <div class="feature-item">
            <span class="feature-icon">✅</span>
            Complete network topology
          </div>
          <div class="feature-item">
            <span class="feature-icon">✅</span>
            Device details & uptime
          </div>
          <div class="feature-item">
            <span class="feature-icon">✅</span>
            Vendor identification
          </div>
        </div>
      </div>
      
      <div class="scan-method snmp-scan">
        <div class="method-header">
          <div class="method-icon">📡</div>
          <div class="method-title">
            <strong>SNMP Only</strong>
            <span class="method-badge snmp">Detailed</span>
          </div>
        </div>
        <p class="method-description">
          Discovers only SNMP-enabled devices with rich information.
          Provides detailed data like hostname, description, vendor,
          system uptime, and performance metrics.
        </p>
        <div class="method-features">
          <div class="feature-item">
            <span class="feature-icon">📊</span>
            System information
          </div>
          <div class="feature-item">
            <span class="feature-icon">⏱️</span>
            Response times
          </div>
          <div class="feature-item">
            <span class="feature-icon">🏷️</span>
            Device hostnames
          </div>
        </div>
      </div>
      
      <div class="scan-method arp-scan">
        <div class="method-header">
          <div class="method-icon">🌐</div>
          <div class="method-title">
            <strong>ARP Only</strong>
            <span class="method-badge arp">Fast</span>
          </div>
        </div>
        <p class="method-description">
          Discovers all IP-enabled devices on the network segment quickly.
          Provides IP and MAC addresses with vendor identification
          based on MAC address database.
        </p>
        <div class="method-features">
          <div class="feature-item">
            <span class="feature-icon">⚡</span>
            Quick discovery
          </div>
          <div class="feature-item">
            <span class="feature-icon">🔗</span>
            MAC addresses
          </div>
          <div class="feature-item">
            <span class="feature-icon">🏭</span>
            Hardware vendors
          </div>
        </div>
      </div>
    </div>
    
    <div class="scan-tips">
      <div class="tips-header">
        <span class="tips-icon">💡</span>
        <h4>Pro Tips</h4>
      </div>
      <div class="tips-content">
        <div class="tip-item">
          <span class="tip-bullet">•</span>
          Use <strong>Full Scan</strong> for comprehensive network auditing and documentation
        </div>
        <div class="tip-item">
          <span class="tip-bullet">•</span>
          Choose <strong>SNMP Only</strong> when you need detailed device information and performance data
        </div>
        <div class="tip-item">
          <span class="tip-bullet">•</span>
          Select <strong>ARP Only</strong> for quick device enumeration and basic network mapping
        </div>
      </div>
    </div>
  </div>
</template>

<script>
export default {
  name: 'ScanTypes'
}
</script>

<style scoped>
.scan-info {
  margin-top: 30px;
  padding: 30px;
  background: rgba(232, 235, 247, 0.95); 
  backdrop-filter: blur(10px);
  border-radius: 20px;
  border: 1px solid rgba(172, 190, 216, 0.2); 
  box-shadow: 0 8px 32px rgba(0, 0, 0, 0.05);
  position: relative;
  overflow: hidden;
}

.scan-info::before {
  content: '';
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  height: 4px;
  background: linear-gradient(90deg, #ACBED8, #F2D398);
}

.info-header {
  text-align: center;
  margin-bottom: 30px;
  position: relative;
}

.info-title {
  color: #333; 
  margin-bottom: 10px;
  font-size: 1.5rem;
  font-weight: 700;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 12px;
}

.title-icon {
  font-size: 1.8rem;
}

.header-decoration {
  width: 60px;
  height: 3px;
  background: linear-gradient(90deg, #ACBED8, #F2D398);
  margin: 0 auto;
  border-radius: 2px;
}

.scan-methods {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(320px, 1fr));
  gap: 25px;
  margin-bottom: 30px;
}

.scan-method {
  padding: 25px;
  border-radius: 16px;
  border: 2px solid transparent;
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
  position: relative;
  overflow: hidden;
  background: rgba(172, 190, 216, 0.1); 
}

.scan-method::before {
  content: '';
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  background: linear-gradient(135deg, transparent, rgba(242, 211, 152, 0.1), transparent);
  transform: translateX(-100%);
  transition: transform 0.6s ease;
}

.scan-method:hover::before {
  transform: translateX(100%);
}

.full-scan {
  border-color: rgba(242, 211, 152, 0.3);
}

.full-scan:hover {
  border-color: #F2D398;
  transform: translateY(-3px);
  box-shadow: 0 10px 25px rgba(242, 211, 152, 0.3);
}

.snmp-scan {
  border-color: rgba(172, 190, 216, 0.3);
}

.snmp-scan:hover {
  border-color: #ACBED8;
  transform: translateY(-3px);
  box-shadow: 0 10px 25px rgba(172, 190, 216, 0.3);
}

.arp-scan {
  border-color: rgba(172, 190, 216, 0.3);
}

.arp-scan:hover {
  border-color: #ACBED8;
  transform: translateY(-3px);
  box-shadow: 0 10px 25px rgba(172, 190, 216, 0.3);
}

.method-header {
  display: flex;
  align-items: center;
  gap: 15px;
  margin-bottom: 15px;
}

.method-icon {
  font-size: 2rem;
  width: 60px;
  height: 60px;
  display: flex;
  align-items: center;
  justify-content: center;
  background: rgba(242, 211, 152, 0.8);
  border-radius: 12px;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.05);
}

.method-title {
  flex: 1;
  display: flex;
  align-items: center;
  gap: 10px;
  flex-wrap: wrap;
}

.method-title strong {
  color: #333;
  font-size: 1.1rem;
  font-weight: 700;
}

.method-badge {
  padding: 4px 10px;
  border-radius: 12px;
  font-size: 0.7rem;
  font-weight: 600;
  text-transform: uppercase;
  letter-spacing: 0.5px;
  color: white;
}

.method-badge.premium {
  background: linear-gradient(135deg, #ACBED8, #F2D398);
}

.method-badge.snmp {
  background: linear-gradient(135deg, #ACBED8, #E8EBF7);
}

.method-badge.arp {
  background: linear-gradient(135deg, #F2D398, #ACBED8);
}

.method-description {
  font-size: 0.95rem;
  color: #4a5568;
  line-height: 1.6;
  margin-bottom: 20px;
}

.method-features {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.feature-item {
  display: flex;
  align-items: center;
  gap: 10px;
  font-size: 0.9rem;
  color: #4a5568;
  font-weight: 500;
}

.feature-icon {
  font-size: 0.8rem;
}

.scan-tips {
  background: linear-gradient(135deg, rgba(242, 211, 152, 0.1), rgba(172, 190, 216, 0.05));
  border-radius: 16px;
  padding: 25px;
  border-left: 4px solid #F2D398;
}

.tips-header {
  display: flex;
  align-items: center;
  gap: 10px;
  margin-bottom: 15px;
}

.tips-icon {
  font-size: 1.3rem;
  color: #F2D398;
}

.tips-header h4 {
  color: #333;
  font-size: 1.1rem;
  font-weight: 700;
  margin: 0;
}

.tips-content {
  display: flex;
  flex-direction: column;
  gap: 10px;
}

.tip-item {
  display: flex;
  align-items: flex-start;
  gap: 12px;
  font-size: 0.9rem;
  color: #4a5568;
  line-height: 1.5;
}

.tip-bullet {
  color: #F2D398;
  font-weight: 700;
  font-size: 1.2rem;
}

@media (max-width: 768px) {
  .scan-info {
    padding: 20px;
  }
  
  .scan-methods {
    grid-template-columns: 1fr;
    gap: 20px;
  }
  
  .scan-method {
    padding: 20px;
  }
  
  .method-header {
    gap: 12px;
  }
  
  .method-icon {
    width: 50px;
    height: 50px;
    font-size: 1.5rem;
  }
  
  .method-title strong {
    font-size: 1rem;
  }
  
  .info-title {
    font-size: 1.3rem;
    flex-direction: column;
    gap: 8px;
  }
}

@media (max-width: 480px) {
  .scan-methods {
    grid-template-columns: 1fr;
  }
  
  .method-header {
    flex-direction: column;
    text-align: center;
    gap: 15px;
  }
  
  .method-title {
    flex-direction: column;
    gap: 8px;
  }
}
</style>
