<template> 
    <div class="card">
        <scan-form></scan-form> 
    </div>
</template>


<script>
import scanForm from './scanForm.vue';
export default {
    components: {
    "scan-form" : scanForm
    }
}
</script>


<style scoped>
.card {
    background: white;
    border-radius: 15px;
    padding: 30px;
    margin-bottom: 20px;
    box-shadow: 0 10px 30px rgba(0, 0, 0, 0.2);
    backdrop-filter: blur(10px);
    max-width: 1400px;
    margin: 0 auto;
    padding: 20px;
}
</style>