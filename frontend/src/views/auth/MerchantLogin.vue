<template>
  <div class="login-page">
    <a href="/" class="back-home">← 返回首頁</a>
    
    <div class="login-container">
      <div class="logo">🏪</div>
      <h1 class="title">商戶登入</h1>
      <p class="subtitle">登入您的商戶後台管理系統</p>
      
      <div class="role-badge">商戶專用</div>
      
      <!-- Demo 帳密提示 -->
      <div class="demo-credentials">
        <h4>🎯 Demo 測試帳號</h4>
        <div class="credential-item">
          <span class="label">電子郵件：</span>
          <code class="credential-value">merchant@example.com</code>
          <button @click="fillDemoCredentials" class="fill-btn">填入</button>
        </div>
        <div class="credential-item">
          <span class="label">密碼：</span>
          <code class="credential-value">111111</code>
          <button @click="fillDemoCredentials" class="fill-btn">填入</button>
        </div>
      </div>
      
      <div v-if="errorMessage" class="error-message">
        {{ errorMessage }}
      </div>
      
      <form @submit.prevent="handleLogin">
        <div class="form-group">
          <label for="email">電子郵件</label>
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
          <label for="password">密碼</label>
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
          {{ loading ? '登入中...' : '登入商戶後台' }}
        </button>
      </form>
      
      <div class="links">
        <router-link to="/merchant/register">註冊商戶帳號</router-link>
        <router-link to="/customer/login">客戶登入</router-link>
        <router-link to="/admin/login">管理員登入</router-link>
      </div>
    </div>
  </div>
</template>

<script>
import { ref, reactive } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth'

export default {
  name: 'MerchantLogin',
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
        }, 'merchant')

        if (result.success) {
          router.push('/merchant/dashboard')
        } else {
          errorMessage.value = result.error || '登入失敗'
        }
      } catch (error) {
        errorMessage.value = '網路錯誤，請稍後再試'
      } finally {
        loading.value = false
      }
    }

    const fillDemoCredentials = () => {
      form.email = 'merchant@example.com'
      form.password = '111111'
    }

    return {
      form,
      errorMessage,
      loading,
      handleLogin,
      fillDemoCredentials
    }
  }
}
</script>

<style scoped>
.login-page {
  font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
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
  background: #e6fffa;
  color: #234e52;
  padding: 4px 12px;
  border-radius: 20px;
  font-size: 0.8rem;
  margin-bottom: 20px;
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
  border-color: #667eea;
}

.form-group input:disabled {
  background-color: #f5f5f5;
  cursor: not-allowed;
}

.login-btn {
  width: 100%;
  padding: 12px;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
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
  color: #667eea;
  text-decoration: none;
  font-size: 0.9rem;
}

.links a:hover {
  text-decoration: underline;
}

.demo-credentials {
  background: #f8f9ff;
  border: 2px solid #e1e5e9;
  border-radius: 8px;
  padding: 1rem;
  margin-bottom: 1.5rem;
}

.demo-credentials h4 {
  margin: 0 0 0.75rem 0;
  color: #333;
  font-size: 0.9rem;
  font-weight: 600;
}

.credential-item {
  display: flex;
  align-items: center;
  margin-bottom: 0.5rem;
  gap: 0.5rem;
}

.credential-item:last-child {
  margin-bottom: 0;
}

.credential-item .label {
  font-size: 0.85rem;
  color: #666;
  min-width: 60px;
}

.credential-value {
  background: #e1e5e9;
  padding: 0.25rem 0.5rem;
  border-radius: 4px;
  font-family: 'Courier New', monospace;
  font-size: 0.8rem;
  color: #333;
  flex: 1;
}

.fill-btn {
  background: #667eea;
  color: white;
  border: none;
  padding: 0.25rem 0.75rem;
  border-radius: 4px;
  font-size: 0.75rem;
  cursor: pointer;
  transition: background-color 0.3s ease;
}

.fill-btn:hover {
  background: #5a67d8;
}

@media (max-width: 480px) {
  .login-container {
    margin: 20px;
    padding: 30px 20px;
  }
}
</style>
