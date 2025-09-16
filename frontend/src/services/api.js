import axios from 'axios'
import { useAuthStore } from '@/stores/auth'

// 創建axios實例
const api = axios.create({
  baseURL: import.meta.env.DEV ? '' : '', // 開發環境使用相對路徑（Vite代理），生產環境也使用相對路徑（同域）
  timeout: 10000,
  headers: {
    'Content-Type': 'application/json'
  }
})

// 請求攔截器 - 添加認證token
api.interceptors.request.use(
  (config) => {
    const authStore = useAuthStore()
    if (authStore.token) {
      config.headers.Authorization = `Bearer ${authStore.token}`
    }
    return config
  },
  (error) => {
    return Promise.reject(error)
  }
)

// 響應攔截器 - 處理認證錯誤
api.interceptors.response.use(
  (response) => {
    return response
  },
  (error) => {
    if (error.response?.status === 401) {
      // Token過期或無效，清除認證狀態
      const authStore = useAuthStore()
      authStore.logout()
    }
    return Promise.reject(error)
  }
)

export default api
