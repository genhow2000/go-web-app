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
    return !!token.value && !!user.value
  })

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
    try {
      await api.post('/logout')
    } catch (error) {
      console.error('登出請求失敗:', error)
      // 即使請求失敗，也要清除本地狀態
    } finally {
      user.value = null
      token.value = null
      localStorage.removeItem('authToken')
    }
  }

  // 獲取用戶信息
  const fetchUserProfile = async () => {
    if (!token.value) return
    
    try {
      // 如果沒有用戶信息，先嘗試從 token 中解析基本信息
      if (!user.value) {
        try {
          // 解析 JWT token 獲取基本信息
          const base64Url = token.value.split('.')[1]
          const base64 = base64Url.replace(/-/g, '+').replace(/_/g, '/')
          const payload = JSON.parse(atob(base64))
          user.value = {
            id: payload.id,
            email: payload.email,
            name: payload.name,
            role: payload.role,
            is_active: payload.is_active
          }
        } catch (error) {
          console.error('JWT token 解析失敗:', error)
          // 如果解析失敗，清除 token
          token.value = null
          localStorage.removeItem('authToken')
          return
        }
      }
      
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
    // 檢查URL參數中的token（用於第三方登入）
    const urlParams = new URLSearchParams(window.location.search)
    const urlToken = urlParams.get('token')
    
    // 重新檢查 token（可能從 cookie 中獲取）
    const localStorageToken = localStorage.getItem('authToken')
    const cookieToken = getTokenFromCookie()
    const currentToken = urlToken || localStorageToken || cookieToken
    
    if (currentToken && currentToken !== token.value) {
      token.value = currentToken
      // 同步到localStorage，確保前端狀態一致
      localStorage.setItem('authToken', currentToken)
      
      // 如果是從URL參數獲取的token，清除URL參數
      if (urlToken) {
        const url = new URL(window.location)
        url.searchParams.delete('token')
        window.history.replaceState({}, '', url)
      }
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
