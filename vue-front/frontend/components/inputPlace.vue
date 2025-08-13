<template> 
  <div class="card-container">
    <div class="card">
      <scan-form></scan-form> 
    </div>
  </div>
</template>

<script>
import scanForm from './scanForm.vue';
export default {
  name: 'InputPlace',
  components: {
    "scan-form": scanForm
  }
}
</script>

<style scoped>
.card-container {
  max-width: 1400px;
  margin: 0 auto;
  padding: 0 20px;
}

.card {
  position: relative;
  background: #E8EBF7;
  border-radius: 20px;
  padding: 40px;
  margin-bottom: 40px;
  box-shadow: 
    0 20px 40px rgba(172, 190, 216, 0.3),
    0 15px 25px rgba(242, 211, 152, 0.15); 
  border: 1px solid rgba(172, 190, 216, 0.4); 
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
  overflow: hidden;
}

.card::before {
  content: '';
  position: absolute;
  top: 0;
  left: -100%;
  width: 100%;
  height: 100%;
  background: linear-gradient(
    90deg,
    transparent,
    rgba(242, 211, 152, 0.4),
    transparent
  );
}



@media (max-width: 768px) {
  .card-container {
    padding: 0 15px;
  }
  
  .card {
    padding: 25px;
    border-radius: 15px;
  }
}

@media (max-width: 480px) {
  .card {
    padding: 20px;
    margin-bottom: 20px;
  }
}
</style>
