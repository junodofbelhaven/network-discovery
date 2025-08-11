<template>
  <div class="bodydiv">
    <header-bar></header-bar>
    <input-place></input-place>
    <footer-bar></footer-bar>
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
  font-family: "Segoe UI", Tahoma, Geneva, Verdana, sans-serif;
  background: linear-gradient(135deg, #515877 0%, #646364 100%);
  min-height: 100vh;
  color: #333;
}




</style>