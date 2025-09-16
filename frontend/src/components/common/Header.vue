<template>
  <nav class="navbar">
    <div class="nav-container">
      <router-link to="/" class="logo">🏪 阿和商城</router-link>
      
      <ul class="nav-links">
        <li><a href="#home">首頁</a></li>
        <li><a href="#categories">分類</a></li>
        <li><a href="#products">商品</a></li>
        <li><a href="#about">關於我們</a></li>
        <li><router-link to="/tech-showcase">技術展示</router-link></li>
      </ul>
      
      <div class="nav-actions">
        <div v-if="!isAuthenticated" class="login-dropdown">
          <a href="#" class="btn btn-outline">登入/註冊</a>
          <div class="dropdown-menu">
            <router-link to="/customer/login">客戶登入</router-link>
            <router-link to="/merchant/login">商戶登入</router-link>
            <router-link to="/admin/login">管理員登入</router-link>
            <hr style="margin: 8px 0; border: none; border-top: 1px solid #e2e8f0;">
            <router-link to="/register">註冊帳號</router-link>
          </div>
        </div>
        
        <div v-else class="user-menu">
          <span class="user-name">歡迎，{{ user?.name }}</span>
          <button @click="handleLogout" class="logout-btn">登出</button>
        </div>
        
        <router-link to="/merchant/register" class="btn btn-success">我要開店</router-link>
      </div>
    </div>
  </nav>
</template>

<script>
import { computed } from 'vue'
import { useAuthStore } from '@/stores/auth'
import { useRouter } from 'vue-router'

export default {
  name: 'Header',
  setup() {
    const authStore = useAuthStore()
    const router = useRouter()
    
    const isAuthenticated = computed(() => authStore.isAuthenticated)
    const user = computed(() => authStore.user)
    
    const handleLogout = async () => {
      if (confirm('確定要登出嗎？')) {
        await authStore.logout()
        router.push('/')
      }
    }
    
    return {
      isAuthenticated,
      user,
      handleLogout
    }
  }
}
</script>

<style scoped>
.navbar {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
  padding: 1rem 0;
  position: sticky;
  top: 0;
  z-index: 1000;
  box-shadow: 0 2px 10px rgba(0,0,0,0.1);
}

.nav-container {
  max-width: 1200px;
  margin: 0 auto;
  padding: 0 20px;
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.logo {
  font-size: 1.8rem;
  font-weight: bold;
  text-decoration: none;
  color: white;
}

.nav-links {
  display: flex;
  list-style: none;
  gap: 2rem;
}

.nav-links a {
  color: white;
  text-decoration: none;
  font-weight: 500;
  transition: opacity 0.3s;
}

.nav-links a:hover {
  opacity: 0.8;
}

.nav-actions {
  display: flex;
  gap: 1rem;
  align-items: center;
}

.btn {
  padding: 0.5rem 1.5rem;
  border-radius: 25px;
  text-decoration: none;
  font-weight: 600;
  transition: all 0.3s;
  border: 2px solid transparent;
}

.btn-outline {
  background: transparent;
  color: white;
  border-color: white;
}

.btn-outline:hover {
  background: white;
  color: #667eea;
}

.btn-success {
  background: #48bb78;
  color: white;
}

.btn-success:hover {
  background: #38a169;
  color: white;
  border-color: #38a169;
}

.login-dropdown {
  position: relative;
  display: inline-block;
}

.dropdown-menu {
  display: none;
  position: absolute;
  top: 100%;
  left: 0;
  background: white;
  min-width: 160px;
  box-shadow: 0 8px 16px rgba(0,0,0,0.1);
  border-radius: 8px;
  z-index: 1000;
  margin-top: 5px;
}

.dropdown-menu a {
  color: #333;
  padding: 12px 16px;
  text-decoration: none;
  display: block;
  transition: background-color 0.3s;
  border-radius: 8px;
  margin: 4px;
}

.dropdown-menu a:hover {
  background-color: #f1f1f1;
}

.login-dropdown:hover .dropdown-menu {
  display: block;
}

.user-menu {
  display: flex;
  align-items: center;
  gap: 1rem;
}

.user-name {
  color: white;
  font-weight: 500;
}

.logout-btn {
  background: #e53e3e;
  color: white;
  border: none;
  padding: 8px 16px;
  border-radius: 6px;
  cursor: pointer;
  font-size: 0.9rem;
}

.logout-btn:hover {
  background: #c53030;
}

@media (max-width: 768px) {
  .nav-links {
    display: none;
  }
}
</style>
