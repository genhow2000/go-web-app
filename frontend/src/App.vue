<template>
  <div id="app">
    <router-view />
  </div>
</template>

<script>
import { useAuthStore } from '@/stores/auth'
import { onMounted } from 'vue'

export default {
  name: 'App',
  setup() {
    const authStore = useAuthStore()
    
    onMounted(async () => {
      // 初始化認證狀態
      await authStore.initAuth()
    })
  }
}
</script>

<style>
#app {
  font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, 'Helvetica Neue', Arial, sans-serif;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
}

* {
  margin: 0;
  padding: 0;
  box-sizing: border-box;
}

body {
  line-height: 1.6;
  color: #333;
  background-color: #f8f9fa;
}
</style>
