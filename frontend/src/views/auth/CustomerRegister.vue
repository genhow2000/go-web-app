<template>
  <div class="register-page">
    <a href="/" class="back-home">← 返回首頁</a>
    
    <div class="register-container">
      <div class="logo">👤</div>
      <h1 class="title">客戶註冊</h1>
      <p class="subtitle">創建您的客戶帳戶，開始購物之旅</p>
      
      <div class="role-badge">客戶專用</div>
      
      <div v-if="errorMessage" class="error-message">
        {{ errorMessage }}
      </div>
      
      <div v-if="successMessage" class="success-message">
        {{ successMessage }}
      </div>
      
      <form @submit.prevent="handleRegister">
        <div class="form-group">
          <label for="name">姓名 *</label>
          <input 
            v-model="form.name"
            type="text" 
            id="name" 
            name="name" 
            required
            :disabled="loading"
          >
        </div>
        
        <div class="form-group">
          <label for="email">電子郵件 *</label>
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
          <label for="password">密碼 *</label>
          <input 
            v-model="form.password"
            type="password" 
            id="password" 
            name="password" 
            required
            minlength="6"
            :disabled="loading"
          >
          <small>密碼至少需要6個字符</small>
        </div>
        
        <div class="form-group">
          <label for="confirmPassword">確認密碼 *</label>
          <input 
            v-model="form.confirmPassword"
            type="password" 
            id="confirmPassword" 
            name="confirmPassword" 
            required
            :disabled="loading"
          >
        </div>
        
        <div class="form-group">
          <label for="phone">電話號碼</label>
          <input 
            v-model="form.phone"
            type="tel" 
            id="phone" 
            name="phone" 
            :disabled="loading"
          >
        </div>
        
        <div class="form-group">
          <label for="address">地址</label>
          <input 
            v-model="form.address"
            type="text" 
            id="address" 
            name="address" 
            :disabled="loading"
          >
        </div>
        
        <div class="form-group">
          <label for="birthDate">生日</label>
          <input 
            v-model="form.birthDate"
            type="date" 
            id="birthDate" 
            name="birthDate" 
            :disabled="loading"
          >
        </div>
        
        <div class="form-group">
          <label for="gender">性別</label>
          <select 
            v-model="form.gender"
            id="gender" 
            name="gender" 
            :disabled="loading"
          >
            <option value="">請選擇</option>
            <option value="男">男</option>
            <option value="女">女</option>
            <option value="其他">其他</option>
          </select>
        </div>
        
        <div class="form-group checkbox-group">
          <label class="checkbox-label">
            <input 
              v-model="form.agreeTerms"
              type="checkbox" 
              required
              :disabled="loading"
            >
            <span class="checkmark"></span>
            我同意 <a href="#" @click.prevent="showTerms">服務條款</a> 和 <a href="#" @click.prevent="showPrivacy">隱私政策</a>
          </label>
        </div>
        
        <button type="submit" class="register-btn" :disabled="loading">
          {{ loading ? '註冊中...' : '註冊客戶帳戶' }}
        </button>
      </form>
      
      <div class="links">
        <router-link to="/customer/login">已有客戶帳戶？登入</router-link>
        <router-link to="/merchant/login">商戶登入</router-link>
        <router-link to="/admin/login">管理員登入</router-link>
      </div>
    </div>
  </div>
</template>

<script>
import { ref, reactive } from 'vue'
import { useRouter } from 'vue-router'
import api from '@/services/api'

export default {
  name: 'CustomerRegister',
  setup() {
    const router = useRouter()
    
    const form = reactive({
      name: '',
      email: '',
      password: '',
      confirmPassword: '',
      phone: '',
      address: '',
      birthDate: '',
      gender: '',
      agreeTerms: false
    })
    
    const errorMessage = ref('')
    const successMessage = ref('')
    const loading = ref(false)

    const handleRegister = async () => {
      errorMessage.value = ''
      successMessage.value = ''
      
      // 驗證密碼
      if (form.password !== form.confirmPassword) {
        errorMessage.value = '密碼確認不一致'
        return
      }
      
      if (form.password.length < 6) {
        errorMessage.value = '密碼至少需要6個字符'
        return
      }
      
      if (!form.agreeTerms) {
        errorMessage.value = '請同意服務條款和隱私政策'
        return
      }
      
      loading.value = true

      try {
        const response = await api.post('/customer/register', {
          name: form.name,
          email: form.email,
          password: form.password,
          phone: form.phone || null,
          address: form.address || null,
          birth_date: form.birthDate || null,
          gender: form.gender || null,
          role: 'customer'
        })

        if (response.data) {
          successMessage.value = '註冊成功！請登入您的客戶帳戶'
          setTimeout(() => {
            router.push('/customer/login')
          }, 2000)
        }
      } catch (error) {
        errorMessage.value = error.response?.data?.error || '註冊失敗，請稍後再試'
      } finally {
        loading.value = false
      }
    }
    
    const showTerms = () => {
      alert('服務條款內容將在此顯示')
    }
    
    const showPrivacy = () => {
      alert('隱私政策內容將在此顯示')
    }

    return {
      form,
      errorMessage,
      successMessage,
      loading,
      handleRegister,
      showTerms,
      showPrivacy
    }
  }
}
</script>

<style scoped>
.register-page {
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

.register-container {
  background: white;
  border-radius: 20px;
  box-shadow: 0 20px 40px rgba(0,0,0,0.1);
  padding: 40px;
  width: 100%;
  max-width: 500px;
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

.form-group input,
.form-group select {
  width: 100%;
  padding: 12px 16px;
  border: 2px solid #e2e8f0;
  border-radius: 10px;
  font-size: 1rem;
  transition: border-color 0.3s ease;
}

.form-group input:focus,
.form-group select:focus {
  outline: none;
  border-color: #667eea;
}

.form-group input:disabled,
.form-group select:disabled {
  background-color: #f5f5f5;
  cursor: not-allowed;
}

.form-group small {
  color: #718096;
  font-size: 0.8rem;
  margin-top: 4px;
  display: block;
}

.checkbox-group {
  margin-bottom: 25px;
}

.checkbox-label {
  display: flex;
  align-items: center;
  cursor: pointer;
  font-size: 0.9rem;
  color: #4a5568;
}

.checkbox-label input[type="checkbox"] {
  width: auto;
  margin-right: 8px;
}

.checkbox-label a {
  color: #667eea;
  text-decoration: none;
}

.checkbox-label a:hover {
  text-decoration: underline;
}

.register-btn {
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

.register-btn:hover:not(:disabled) {
  transform: translateY(-2px);
}

.register-btn:disabled {
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

.success-message {
  background: #c6f6d5;
  color: #22543d;
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

@media (max-width: 480px) {
  .register-container {
    margin: 20px;
    padding: 30px 20px;
  }
  
  .links {
    flex-direction: column;
    align-items: center;
  }
}
</style>
