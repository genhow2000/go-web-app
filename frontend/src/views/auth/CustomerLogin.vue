<template>
  <div class="login-page">
    <div class="login-container">
      <div class="login-header">
        <h1>客戶登入</h1>
        <p>歡迎回來，請登入您的客戶帳戶</p>
      </div>

      <!-- Demo 帳密提示 -->
      <div class="demo-credentials">
        <h4>🎯 Demo 測試帳號</h4>
        <div class="credential-item">
          <span class="label">電子郵件：</span>
          <code class="credential-value">customer@example.com</code>
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
      <div v-if="successMessage" class="success-message">
        {{ successMessage }}
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
          {{ loading ? '登入中...' : '登入' }}
        </button>
      </form>

      <div class="divider">
        <span>或</span>
      </div>

      <!-- LINE登入按鈕 -->
      <div class="oauth-login">
        <button @click="handleLineLogin" class="btn-line-login" :disabled="loading">
          <svg width="20" height="20" viewBox="0 0 24 24" fill="currentColor">
            <path d="M19.365 9.863c.349 0 .63.285.63.631 0 .345-.281.63-.63.63H17.61v1.125h1.755c.349 0 .63.283.63.63 0 .344-.281.629-.63.629h-2.386c-.345 0-.627-.285-.627-.629V8.108c0-.345.282-.63.63-.63h2.386c.349 0 .63.285.63.63 0 .346-.281.63-.63.63H17.61v1.125h1.755zm-3.855 3.016c0 .27-.174.51-.432.596-.064.021-.133.031-.199.031-.211 0-.391-.09-.51-.25l-2.443-3.317v2.94c0 .344-.279.629-.631.629-.346 0-.626-.285-.626-.629V8.108c0-.27.173-.51.43-.595.06-.023.136-.033.194-.033.195 0 .375.104.495.254l2.462 3.33V8.108c0-.345.282-.63.63-.63.345 0 .63.285.63.63v4.771zm-5.741 0c0 .344-.282.629-.631.629-.345 0-.626-.285-.626-.629V8.108c0-.345.281-.63.63-.63.346 0 .63.285.63.63v4.771zm-2.466.629H4.917c-.345 0-.63-.285-.63-.629V8.108c0-.345.285-.63.63-.63.348 0 .63.285.63.63v4.141h1.756c.348 0 .629.283.629.63 0 .344-.281.629-.629.629M24 10.314C24 4.943 18.615.572 12 .572S0 4.943 0 10.314c0 4.811 4.27 8.842 10.035 9.608.391.082.923.258 1.058.59.12.301.079.766.038 1.08l-.164 1.02c-.045.301-.24 1.186 1.049.645 1.291-.539 6.916-4.078 9.436-6.975C23.176 14.393 24 12.458 24 10.314"/>
          </svg>
          {{ loading ? '跳轉中...' : '使用 LINE 登入' }}
        </button>
      </div>

      <div class="divider">
        <span>其他登入方式</span>
      </div>

      <div class="role-links">
        <router-link to="/merchant/login" class="role-link">商戶登入</router-link>
        <router-link to="/admin/login" class="role-link">管理員登入</router-link>
      </div>

      <div class="register-link">
        <p>還沒有帳戶？ <router-link to="/register">立即註冊</router-link></p>
      </div>
    </div>
  </div>
</template>

<script>
import { ref, reactive } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth'

export default {
  name: 'CustomerLogin',
  setup() {
    const router = useRouter()
    const authStore = useAuthStore()
    
    const form = reactive({
      email: '',
      password: ''
    })
    
    const errorMessage = ref('')
    const successMessage = ref('')
    const loading = ref(false)

    const handleLogin = async () => {
      errorMessage.value = ''
      successMessage.value = ''
      loading.value = true

      try {
        console.log('開始登入流程...', { email: form.email, role: 'customer' })
        
        const result = await authStore.login({
          email: form.email,
          password: form.password
        }, 'customer')

        console.log('登入結果:', result)

        if (result.success) {
          successMessage.value = result.data.message || '登入成功！'
          console.log('登入成功，準備跳轉到儀表板...')
          setTimeout(() => {
            console.log('執行路由跳轉...')
            router.push('/customer/dashboard')
          }, 1000)
        } else {
          console.error('登入失敗:', result.error)
          errorMessage.value = result.error || '登入失敗，請重試'
        }
      } catch (error) {
        console.error('登入過程發生錯誤:', error)
        errorMessage.value = '網路錯誤，請重試'
      } finally {
        loading.value = false
      }
    }

    const handleLineLogin = () => {
      loading.value = true
      // 直接跳轉到後端的 LINE 登入端點
      window.location.href = '/auth/line'
    }

    const fillDemoCredentials = () => {
      form.email = 'customer@example.com'
      form.password = '111111'
    }

    return {
      form,
      errorMessage,
      successMessage,
      loading,
      handleLogin,
      handleLineLogin,
      fillDemoCredentials
    }
  }
}
</script>

<style scoped>
.login-page {
  font-family: 'Segoe UI', Tahoma, Geneva, Verdana, sans-serif;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  min-height: 100vh;
  display: flex;
  align-items: center;
  justify-content: center;
}

.login-container {
  background: white;
  padding: 2rem;
  border-radius: 10px;
  box-shadow: 0 15px 35px rgba(0, 0, 0, 0.1);
  width: 100%;
  max-width: 400px;
}

.login-header {
  text-align: center;
  margin-bottom: 2rem;
}

.login-header h1 {
  color: #333;
  margin-bottom: 0.5rem;
  font-size: 1.8rem;
}

.login-header p {
  color: #666;
  font-size: 0.9rem;
}

.form-group {
  margin-bottom: 1.5rem;
}

.form-group label {
  display: block;
  margin-bottom: 0.5rem;
  color: #333;
  font-weight: 500;
}

.form-group input {
  width: 100%;
  padding: 0.75rem;
  border: 2px solid #e1e5e9;
  border-radius: 5px;
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
  padding: 0.75rem;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
  border: none;
  border-radius: 5px;
  font-size: 1rem;
  font-weight: 500;
  cursor: pointer;
  transition: transform 0.2s ease;
}

.login-btn:hover:not(:disabled) {
  transform: translateY(-2px);
}

.login-btn:disabled {
  opacity: 0.6;
  cursor: not-allowed;
  transform: none;
}

.divider {
  text-align: center;
  margin: 1.5rem 0;
  position: relative;
  color: #666;
}

.divider::before {
  content: '';
  position: absolute;
  top: 50%;
  left: 0;
  right: 0;
  height: 1px;
  background: #e1e5e9;
}

.divider span {
  background: white;
  padding: 0 1rem;
}

.role-links {
  display: flex;
  justify-content: space-between;
  margin-top: 1rem;
}

.role-link {
  color: #667eea;
  text-decoration: none;
  font-size: 0.9rem;
  padding: 0.5rem;
  border-radius: 5px;
  transition: background-color 0.3s ease;
}

.role-link:hover {
  background-color: #f8f9ff;
}

.register-link {
  text-align: center;
  margin-top: 1.5rem;
}

.register-link a {
  color: #667eea;
  text-decoration: none;
  font-weight: 500;
}

.register-link a:hover {
  text-decoration: underline;
}

.error-message {
  background-color: #fee;
  color: #c33;
  padding: 0.75rem;
  border-radius: 5px;
  margin-bottom: 1rem;
  border: 1px solid #fcc;
}

.success-message {
  background-color: #efe;
  color: #3c3;
  padding: 0.75rem;
  border-radius: 5px;
  margin-bottom: 1rem;
  border: 1px solid #cfc;
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

.oauth-login {
  margin: 1rem 0;
}

.btn-line-login {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 0.5rem;
  width: 100%;
  padding: 0.75rem;
  background: #00B900;
  color: white;
  text-decoration: none;
  border: none;
  border-radius: 5px;
  font-weight: 500;
  cursor: pointer;
  transition: background-color 0.3s ease, transform 0.2s ease;
}

.btn-line-login:hover:not(:disabled) {
  background: #009900;
  transform: translateY(-2px);
}

.btn-line-login:disabled {
  background: #ccc;
  cursor: not-allowed;
  transform: none;
}

.btn-line-login svg {
  flex-shrink: 0;
}
</style>
