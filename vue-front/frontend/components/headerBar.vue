<template>
  <div class="container">
    <header class="modern-header">
      <div class="header-glow"></div>
      <div class="header-content">
        <div class="icon-wrapper">
          <div class="rotating-icon">🌐</div>
        </div>
        <h1 class="glitch-text">Network Discovery Tool</h1>
        <p class="subtitle">Discover and analyze network devices using SNMP and ARP protocols</p>
        <div class="pulse-line"></div>
      </div>
    </header>
  </div>
</template>

<script>
export default {
  name: 'HeaderBar'
}
</script>

<style scoped>
.container {
  max-width: 1400px;
  margin: 0 auto;
  padding: 40px 20px;
}

.modern-header {
  position: relative;
  text-align: center;
  color: #ACBED8;
  margin-bottom: 40px;
  overflow: hidden;
}

.header-glow {
  position: absolute;
  top: 50%;
  left: 50%;
  transform: translate(-50%, -50%);
  width: 300px;
  height: 300px;
  background: radial-gradient(circle, rgba(242, 211, 152, 0.2) 0%, transparent 70%); 
  animation: pulse-glow 3s ease-in-out infinite alternate;
}

.header-content {
  position: relative;
  z-index: 2;
}

.icon-wrapper {
  margin-bottom: 20px;
}

.rotating-icon {
  font-size: 4rem;
  display: inline-block;
  animation: rotate-float 8s ease-in-out infinite;
  filter: drop-shadow(0 0 15px rgba(172, 190, 216, 0.4)); 
  color: #ACBED8;
}

.glitch-text {
  font-size: clamp(2rem, 5vw, 3.5rem);
  font-weight: 700;
  margin-bottom: 15px;
  text-shadow: 
    2px 2px 4px rgba(30, 30, 50, 0.5), 
    0 0 25px rgba(123, 112, 210, 0.6); 
  background: linear-gradient(45deg, #9CA3E1, #A2B2FF, #F2D398); 
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
  animation: text-shimmer 3s ease-in-out infinite;
  position: relative;
}

.subtitle {
  font-size: 1.2rem;
  opacity: 0.9;
  font-weight: 300;
  letter-spacing: 0.5px;
  margin-bottom: 30px;
  animation: fade-in-up 1s ease-out 0.5s both;
  color: #A2B2FF;
}

.pulse-line {
  width: 100px;
  height: 2px;
  background: linear-gradient(90deg, transparent, #A2B2FF, transparent); 
  margin: 0 auto;
  animation: pulse-width 2s ease-in-out infinite alternate;
}

@keyframes pulse-glow {
  0% {
    opacity: 0.25;
    transform: translate(-50%, -50%) scale(1);
  }
  100% {
    opacity: 0.45;
    transform: translate(-50%, -50%) scale(1.05);
  }
}
@keyframes rotate-float {
  0%, 100% {
    transform: rotate(0deg) translateY(0px);
  }
  25% {
    transform: rotate(5deg) translateY(-10px);
  }
  50% {
    transform: rotate(0deg) translateY(-5px);
  }
  75% {
    transform: rotate(-5deg) translateY(-15px);
  }
}

@keyframes text-shimmer {
  0%, 100% {
    background-position: 0% 50%;
  }
  50% {
    background-position: 100% 50%;
  }
}

@keyframes fade-in-up {
  from {
    opacity: 0;
    transform: translateY(30px);
  }
  to {
    opacity: 0.9;
    transform: translateY(0);
  }
}

@keyframes pulse-width {
  0% {
    width: 80px;
    opacity: 0.5;
  }
  100% {
    width: 120px;
    opacity: 1;
  }
}

@media (max-width: 768px) {
  .container {
    padding: 20px 15px;
  }
  
  .rotating-icon {
    font-size: 3rem;
  }
  
  .subtitle {
    font-size: 1rem;
    padding: 0 20px;
  }
}
</style>
