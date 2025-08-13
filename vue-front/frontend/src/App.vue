<template>
  <div class="bodydiv">
    <div class="background-pattern"></div>
    <div class="content-wrapper">
      <header-bar></header-bar>
      <input-place></input-place>
      <footer-bar></footer-bar>
    </div>
  </div>
</template>

<script>
import headerBar from '/components/headerBar.vue'
import inputPlace from '../components/inputPlace.vue'
import footerBar from "../components/footerBar.vue"

export default{
  components : {
    "header-bar" : headerBar,
    "input-place" : inputPlace,
    "footer-bar" : footerBar
  },
  data(){
    return{

    }
  }
}
</script>

<style>
* {
  margin: 0;
  padding: 0;
  box-sizing: border-box;
}

.bodydiv {
  font-family: "Inter", "Segoe UI", Tahoma, Geneva, Verdana, sans-serif;
  background: linear-gradient(135deg, #ACBED8 0%, #E8EBF7 50%, #F2D398 100%);
  min-height: 100vh;
  color: #333;
  position: relative;
  overflow-x: hidden;
}

.background-pattern {
  position: fixed;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  background: 
    radial-gradient(circle at 20% 80%, rgba(172, 190, 216, 0.3) 0%, transparent 50%),
    radial-gradient(circle at 80% 20%, rgba(242, 211, 152, 0.3) 0%, transparent 50%),
    radial-gradient(circle at 40% 40%, rgba(172, 190, 216, 0.2) 0%, transparent 50%);
  animation: float 20s ease-in-out infinite;
  z-index: -1;
}

.content-wrapper {
  position: relative;
  z-index: 1;
  backdrop-filter: blur(1px);
}

@keyframes float {
  0%, 100% {
    transform: translateY(0px) rotate(0deg);
  }
  33% {
    transform: translateY(-20px) rotate(1deg);
  }
  66% {
    transform: translateY(10px) rotate(-1deg);
  }
}

/* Custom scrollbar */
::-webkit-scrollbar {
  width: 8px;
}

::-webkit-scrollbar-track {
  background: rgba(232, 235, 247, 0.3);
  border-radius: 10px;
}

::-webkit-scrollbar-thumb {
  background: linear-gradient(135deg, #ACBED8, #F2D398);
  border-radius: 10px;
}


/* Global button animations */
button {
  position: relative;
  overflow: hidden;
}

button::before {
  content: '';
  position: absolute;
  top: 50%;
  left: 50%;
  width: 0;
  height: 0;
  background: rgba(242, 211, 152, 0.2);
  border-radius: 50%;
  transform: translate(-50%, -50%);
  transition: width 0.6s, height 0.6s;
}

button:hover::before {
  width: 300px;
  height: 300px;
}
</style>
