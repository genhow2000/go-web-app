<template>
  <nav class="navbar">
    <div class="nav-container">
      <router-link to="/" class="logo">🏪 阿和商城</router-link>
      
      <ul class="nav-links">
        <li><a href="#home">首頁</a></li>
        <li><a href="#categories">分類</a></li>
        <li><a href="#products">商品</a></li>
        <li><router-link to="/stock-market">阿和台股站</router-link></li>
        <li><router-link to="/stocks">股票列表</router-link></li>
        <li><a href="#about">關於我們</a></li>
        <li><router-link to="/tech-showcase">技術展示</router-link></li>
      </ul>
      
      <div class="nav-actions">
        <!-- 版本號顯示 -->
        <div class="version-info">
          <span class="version-text">v{{ version }}</span>
        </div>
        
        <!-- 購物車圖標 -->
        <CartIcon 
          v-if="isAuthenticated && user?.role === 'customer'"
          @mouseenter="handleMouseEnter"
          @mouseleave="handleMouseLeave"
        />
        
        <div v-if="!isAuthenticated" class="login-dropdown">
          <button class="btn btn-outline" @click="toggleDropdown">登入/註冊</button>
          <div v-show="showDropdown" class="dropdown-menu">
            <router-link to="/customer/login" @click="closeDropdown">客戶登入</router-link>
            <router-link to="/merchant/login" @click="closeDropdown">商戶登入</router-link>
            <router-link to="/admin/login" @click="closeDropdown">管理員登入</router-link>
            <hr style="margin: 8px 0; border: none; border-top: 1px solid #e2e8f0;">
            <router-link to="/register" @click="closeDropdown">註冊帳號</router-link>
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
import { computed, ref, onMounted, onUnmounted } from 'vue'
import { useAuthStore } from '@/stores/auth'
import { useRouter } from 'vue-router'
import CartIcon from '@/components/cart/CartIcon.vue'
import api from '@/services/api'

export default {
  name: 'Header',
  components: {
    CartIcon
  },
  setup() {
    const authStore = useAuthStore()
    const router = useRouter()
    const showDropdown = ref(false)
    const version = ref('2.0.0')
    
    const isAuthenticated = computed(() => authStore.isAuthenticated)
    const user = computed(() => authStore.user)
    
    const toggleDropdown = () => {
      showDropdown.value = !showDropdown.value
    }
    
    const closeDropdown = () => {
      showDropdown.value = false
    }

    // 購物車圖標事件處理
    const handleMouseEnter = () => {
      // 購物車圖標的懸停事件
    }

    const handleMouseLeave = () => {
      // 購物車圖標的離開事件
    }
    
    const handleLogout = async () => {
      console.log('用戶點擊登出按鈕')
      console.log('開始登出流程...')
      try {
        await authStore.logout()
        console.log('登出成功，準備跳轉')
        showDropdown.value = false
        router.push('/')
      } catch (error) {
        console.error('登出過程中發生錯誤:', error)
        // 即使登出失敗，也要跳轉到首頁
        showDropdown.value = false
        router.push('/')
      }
    }
    
    // 獲取版本號
    const loadVersion = async () => {
      try {
        const response = await api.get('/api/version/short')
        if (response.data.success) {
          version.value = response.data.data.version
        }
      } catch (error) {
        console.error('獲取版本號失敗:', error)
        // 保持默認版本號
      }
    }
    
    // 點擊外部關閉下拉菜單
    const handleClickOutside = (event) => {
      const dropdown = event.target.closest('.login-dropdown')
      if (!dropdown && showDropdown.value) {
        showDropdown.value = false
      }
    }
    
    onMounted(() => {
      document.addEventListener('click', handleClickOutside)
      loadVersion()
    })
    
    onUnmounted(() => {
      document.removeEventListener('click', handleClickOutside)
    })
    
    return {
      isAuthenticated,
      user,
      showDropdown,
      version,
      toggleDropdown,
      closeDropdown,
      handleMouseEnter,
      handleMouseLeave,
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

.version-info {
  display: flex;
  align-items: center;
  margin-right: 0.5rem;
}

.version-text {
  font-size: 0.75rem;
  color: rgba(255, 255, 255, 0.8);
  background: rgba(255, 255, 255, 0.1);
  padding: 0.25rem 0.5rem;
  border-radius: 4px;
  font-family: 'Courier New', monospace;
  font-weight: 500;
  border: 1px solid rgba(255, 255, 255, 0.2);
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
