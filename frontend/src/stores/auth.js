import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import api from '@/services/api'

export const useAuthStore = defineStore('auth', () => {
  const user = ref(null)
  const token = ref(localStorage.getItem('authToken'))
  const loading = ref(false)

  const isAuthenticated = computed(() => !!token.value && !!user.value)

  // 登入
  const login = async (credentials, role = 'customer') => {
    loading.value = true
    try {
      // 根據角色選擇正確的登入端點
      const loginEndpoint = `/${role}/login`
      const response = await api.post(loginEndpoint, credentials)
      const { token: newToken, user: userData } = response.data
      
      token.value = newToken
      user.value = userData
      localStorage.setItem('authToken', newToken)
      
      return { success: true, data: response.data }
    } catch (error) {
      return { 
        success: false, 
        error: error.response?.data?.error || '登入失敗' 
      }
    } finally {
      loading.value = false
    }
  }

  // 登出
  const logout = async () => {
    try {
      await api.post('/logout')
    } catch (error) {
      console.error('登出請求失敗:', error)
    } finally {
      user.value = null
      token.value = null
      localStorage.removeItem('authToken')
    }
  }

  // 獲取用戶信息
  const fetchUserProfile = async () => {
    if (!token.value || !user.value) return
    
    try {
      // 根據用戶角色選擇正確的個人資料端點
      const role = user.value.role
      const profileEndpoint = `/${role}/profile`
      const response = await api.get(profileEndpoint)
      user.value = response.data
    } catch (error) {
      console.error('獲取用戶信息失敗:', error)
      // 如果token無效，清除認證狀態
      if (error.response?.status === 401) {
        logout()
      }
    }
  }

  // 初始化認證狀態
  const initAuth = async () => {
    if (token.value) {
      await fetchUserProfile()
    }
  }

  return {
    user,
    token,
    loading,
    isAuthenticated,
    login,
    logout,
    fetchUserProfile,
    initAuth
  }
})
