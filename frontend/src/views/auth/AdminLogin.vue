<template>
  <div class="login-page">
    <a href="/" class="back-home">← 返回首頁</a>
    
    <div class="login-container">
      <div class="logo">👑</div>
      <h1 class="title">管理員登入</h1>
      <p class="subtitle">系統管理員專用登入</p>
      
      <div class="role-badge">管理員專用</div>
      
      <div class="security-notice">
        ⚠️ 此為系統管理員專用登入，請確保您有適當的權限
      </div>
      
      <div v-if="errorMessage" class="error-message">
        {{ errorMessage }}
      </div>
      
      <form @submit.prevent="handleLogin">
        <div class="form-group">
          <label for="email">管理員電子郵件</label>
          <input 
            v-model="form.email"
            type="email" 
            id="email" 
            name="email" 
            required
            :disabled="loading"
          >
        </div>
        
        <div class="form-group">
          <label for="password">管理員密碼</label>
          <input 
            v-model="form.password"
            type="password" 
            id="password" 
            name="password" 
            required
            :disabled="loading"
          >
        </div>
        
        <button type="submit" class="login-btn" :disabled="loading">
          {{ loading ? '登入中...' : '登入管理員後台' }}
        </button>
      </form>
      
      <div class="links">
        <router-link to="/merchant/login">商戶登入</router-link>
        <router-link to="/admin/db/login">資料庫管理</router-link>
      </div>
    </div>
  </div>
</template>

<script>
import { ref, reactive } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth'

export default {
  name: 'AdminLogin',
  setup() {
    const router = useRouter()
    const authStore = useAuthStore()
    
    const form = reactive({
      email: '',
      password: ''
    })
    
    const errorMessage = ref('')
    const loading = ref(false)

    const handleLogin = async () => {
      errorMessage.value = ''
      loading.value = true

      try {
        const result = await authStore.login({
          email: form.email,
          password: form.password
        }, 'admin')

        if (result.success) {
          router.push('/admin/dashboard')
        } else {
          errorMessage.value = result.error || '登入失敗'
        }
      } catch (error) {
        errorMessage.value = '網路錯誤，請稍後再試'
      } finally {
        loading.value = false
      }
    }

    return {
      form,
      errorMessage,
      loading,
      handleLogin
    }
  }
}
</script>

<style scoped>
.login-page {
  font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif;
  background: linear-gradient(135deg, #e53e3e 0%, #9f7aea 100%);
  min-height: 100vh;
  display: flex;
  align-items: center;
  justify-content: center;
  position: relative;
}

.back-home {
  position: absolute;
  top: 20px;
  left: 20px;
  color: white;
  text-decoration: none;
  font-size: 1.1rem;
  display: flex;
  align-items: center;
}

.back-home:hover {
  text-decoration: underline;
}

.login-container {
  background: white;
  border-radius: 20px;
  box-shadow: 0 20px 40px rgba(0,0,0,0.1);
  padding: 40px;
  width: 100%;
  max-width: 400px;
  text-align: center;
}

.logo {
  font-size: 2.5rem;
  margin-bottom: 10px;
  color: #4a5568;
}

.title {
  font-size: 1.8rem;
  color: #2d3748;
  margin-bottom: 10px;
  font-weight: bold;
}

.subtitle {
  color: #718096;
  margin-bottom: 30px;
  font-size: 1rem;
}

.role-badge {
  display: inline-block;
  background: #fed7d7;
  color: #c53030;
  padding: 4px 12px;
  border-radius: 20px;
  font-size: 0.8rem;
  margin-bottom: 20px;
}

.security-notice {
  background: #fef5e7;
  color: #744210;
  padding: 12px;
  border-radius: 8px;
  margin-bottom: 20px;
  font-size: 0.9rem;
}

.form-group {
  margin-bottom: 20px;
  text-align: left;
}

.form-group label {
  display: block;
  margin-bottom: 8px;
  color: #4a5568;
  font-weight: 500;
}

.form-group input {
  width: 100%;
  padding: 12px 16px;
  border: 2px solid #e2e8f0;
  border-radius: 10px;
  font-size: 1rem;
  transition: border-color 0.3s ease;
}

.form-group input:focus {
  outline: none;
  border-color: #e53e3e;
}

.form-group input:disabled {
  background-color: #f5f5f5;
  cursor: not-allowed;
}

.login-btn {
  width: 100%;
  padding: 12px;
  background: linear-gradient(135deg, #e53e3e 0%, #9f7aea 100%);
  color: white;
  border: none;
  border-radius: 10px;
  font-size: 1rem;
  font-weight: bold;
  cursor: pointer;
  transition: transform 0.3s ease;
}

.login-btn:hover:not(:disabled) {
  transform: translateY(-2px);
}

.login-btn:disabled {
  opacity: 0.6;
  cursor: not-allowed;
  transform: none;
}

.error-message {
  background: #fed7d7;
  color: #c53030;
  padding: 12px;
  border-radius: 8px;
  margin-bottom: 20px;
}

.links {
  margin-top: 20px;
  display: flex;
  justify-content: space-between;
  flex-wrap: wrap;
  gap: 10px;
}

.links a {
  color: #e53e3e;
  text-decoration: none;
  font-size: 0.9rem;
}

.links a:hover {
  text-decoration: underline;
}

@media (max-width: 480px) {
  .login-container {
    margin: 20px;
    padding: 30px 20px;
  }
}
</style>
