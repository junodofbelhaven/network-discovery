<template>
  <div class="container">
    <footer class="modern-footer">
      <div class="footer-glow"></div>
      <div class="footer-content">
        <div class="footer-main">
          <h3 class="footer-title">
            Network Discovery Excellence
          </h3>
          <p class="footer-subtitle">
            Discover and analyze network devices using SNMP and ARP protocols
          </p>
        </div>
        
        <div class="footer-divider"></div>
        
        <div class="footer-credits">
          <div class="credits-content">
            <span class="copyright">© 2025 Network Discovery Tool</span>
            <div class="authors">
              <span class="by-text">Crafted with</span>
              <span class="tech-icon vue" title="Vue.js">Vue</span>
              <span class="tech-icon go" title="Go">Go</span>
              <span class="tech-icon js" title="JavaScript">JS</span>
              <span class="by-text">by</span>
              <div class="author-names">
                <span class="author-name">Ozan</span>
                <span class="author-separator">&</span>
                <span class="author-name">Anıl</span>
              </div>
            </div>
          </div>
        </div>
        
        <div class="footer-wave">
          <svg viewBox="0 0 1200 120" xmlns="http://www.w3.org/2000/svg">
            <path d="M0,60 C150,100 350,0 600,60 C850,120 1050,20 1200,60 V120 H0 V60 Z" fill="rgba(242, 211, 152, 0.15)"/>
          </svg>
        </div>
      </div>
    </footer>
  </div>
</template>

<script>
export default {
  name: 'FooterBar'
}
</script>

<style scoped>
.container {
  max-width: 1400px;
  margin: 0 auto;
  padding: 40px 20px 20px;
}

.modern-footer {
  position: relative;
  text-align: center;
  color: #3a3a3a;
  background: linear-gradient(135deg, #ACBED8, #E8EBF7);
  border-radius: 20px;
  padding: 40px 30px 20px;
  border: 1px solid rgba(172, 190, 216, 0.5);
  box-shadow: 
    0 8px 32px rgba(172, 190, 216, 0.3),
    inset 0 1px 0 rgba(255, 255, 255, 0.3);
  overflow: hidden;
}

.footer-glow {
  position: absolute;
  top: -40%;
  left: 50%;
  transform: translateX(-50%);
  width: 180px;
  height: 180px;
  background: radial-gradient(circle, rgba(172, 140, 210, 0.3) 0%, transparent 70%);
  animation: pulse-glow 4s ease-in-out infinite alternate;
  pointer-events: none;
  z-index: 1;
}

.footer-content {
  position: relative;
  z-index: 2;
}

.footer-main {
  margin-bottom: 25px;
}

.footer-title {
  font-size: 1.3rem;
  font-weight: 700;
  margin-bottom: 10px;
  color: #4a4a4a;
}

.footer-subtitle {
  font-size: 0.95rem;
  opacity: 0.85;
  font-weight: 400;
  letter-spacing: 0.3px;
  line-height: 1.4;
  color: #565656;
}

.footer-divider {
  width: 80px;
  height: 2px;
  background: linear-gradient(90deg, transparent, rgba(172, 140, 210, 1), transparent);
  margin: 25px auto;
}

.footer-credits {
  display: flex;
  justify-content:center;
  align-items: center;
  flex-wrap: wrap;
  gap: 20px;
  margin-bottom: 20px;
}

.credits-content {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 8px;
  text-align: center;
}

.copyright {
  font-size: 0.85rem;
  opacity: 0.7;
  font-weight: 400;
  color: #5c5c5c;
}

.authors {
  display: flex;
  align-items: center;
  gap: 6px;
  font-size: 0.9rem;
  color: #3a3a3a;
}

.by-text {
  opacity: 0.8;
  font-weight: 300;
}

.author-names {
  display: flex;
  align-items: center;
  gap: 8px;
  color: #3a3a3a;
}

.author-name {
  font-weight: 600;
  cursor: pointer;
  transition: all 0.3s ease;
  color: #5a5a5a;
}

.author-separator {
  opacity: 0.6;
  font-weight: 300;
  color: #6d6d6d;
}

.tech-icon {
  padding: 4px 8px;
  border-radius: 6px;
  font-size: 0.7rem;
  font-weight: 600;
  text-transform: uppercase;
  letter-spacing: 0.5px;
  transition: all 0.3s ease;
  cursor: pointer;
  user-select: none;
  color: #5a5a5a;
  border: 1px solid transparent;
}

.tech-icon.vue {
  background: rgba(172, 190, 216, 0.3);
  color: #2d6a9f;
  border-color: rgba(45, 106, 159, 0.3);
}

.tech-icon.go {
  background: rgba(242, 211, 152, 0.3);
  color: #b58e40;
  border-color: rgba(181, 142, 64, 0.3);
}

.tech-icon.js {
  background: rgba(255, 223, 93, 0.3); 
  color: #d4af37; 
  border-color: rgba(212, 175, 55, 0.4);
}


.footer-wave {
  position: absolute;
  bottom: 0;
  left: 0;
  width: 100%;
  height: 60px;
  overflow: hidden;
  opacity: 0.2;
  z-index: 0;
}

.footer-wave svg {
  width: 100%;
  height: 100%;
  animation: wave-flow 8s ease-in-out infinite;
}

@keyframes pulse-glow {
  0% {
    opacity: 0.15;
    transform: translateX(-50%) scale(1);
  }
  100% {
    opacity: 0.35;
    transform: translateX(-50%) scale(1.1);
  }
}

@keyframes wave-flow {
  0%, 100% {
    transform: translateX(0);
  }
  50% {
    transform: translateX(-20px);
  }
}

@media (max-width: 768px) {
  .container {
    padding: 30px 15px 15px;
  }
  
  .modern-footer {
    padding: 30px 20px 15px;
  }
  
  .footer-credits {
    flex-direction: column;
    gap: 15px;
    text-align: center;
  }
  
  .credits-content,
  .footer-tech {
    text-align: center;
  }
  
  .footer-title {
    font-size: 1.1rem;
  }
  
  .authors {
    justify-content: center;
    flex-wrap: wrap;
  }
}

@media (max-width: 480px) {
  .modern-footer {
    padding: 25px 15px 10px;
  }
  
  .footer-title {
    font-size: 1rem;
  }
  
  .footer-subtitle {
    font-size: 0.9rem;
    padding: 0 10px;
  }
  
  .author-names {
    flex-wrap: wrap;
    justify-content: center;
  }
}
</style>
