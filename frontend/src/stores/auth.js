import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import api from '@/services/api'

// 從 cookie 中獲取 token 的輔助函數
const getTokenFromCookie = () => {
  const cookies = document.cookie.split(';')
  for (let cookie of cookies) {
    const [name, value] = cookie.trim().split('=')
    if (name === 'auth_token') {
      return value
    }
  }
  return null
}

export const useAuthStore = defineStore('auth', () => {
  const user = ref(null)
  const token = ref(localStorage.getItem('authToken') || getTokenFromCookie())
  const loading = ref(false)

  const isAuthenticated = computed(() => {
    const result = !!token.value && !!user.value
    console.log('認證狀態計算:', { 
      token: !!token.value, 
      user: !!user.value, 
      isAuthenticated: result 
    })
    return result
  })

  // 登入
  const login = async (credentials, role = 'customer') => {
    loading.value = true
    try {
      // 根據角色選擇正確的登入端點
      const loginEndpoint = `/${role}/login`
      console.log('發送登入請求到:', loginEndpoint, credentials)
      
      const response = await api.post(loginEndpoint, credentials)
      console.log('API響應:', response.data)
      
      const { token: newToken, user: userData } = response.data
      
      token.value = newToken
      user.value = userData
      localStorage.setItem('authToken', newToken)
      
      console.log('認證狀態已更新:', { token: !!newToken, user: !!userData })
      
      return { success: true, data: response.data }
    } catch (error) {
      console.error('登入API錯誤:', error)
      console.error('錯誤詳情:', error.response?.data)
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
    console.log('開始登出流程...')
    try {
      await api.post('/logout')
      console.log('登出請求成功')
    } catch (error) {
      console.error('登出請求失敗:', error)
      // 即使請求失敗，也要清除本地狀態
    } finally {
      console.log('清除本地認證狀態...')
      user.value = null
      token.value = null
      localStorage.removeItem('authToken')
      console.log('登出完成，認證狀態:', { 
        isAuthenticated: !!token.value && !!user.value,
        token: !!token.value,
        user: !!user.value
      })
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
    // 重新檢查 token（可能從 cookie 中獲取）
    const currentToken = localStorage.getItem('authToken') || getTokenFromCookie()
    if (currentToken && currentToken !== token.value) {
      token.value = currentToken
    }
    
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
