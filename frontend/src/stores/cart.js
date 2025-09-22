import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import api from '@/services/api'

export const useCartStore = defineStore('cart', () => {
  // 狀態
  const items = ref([])
  const loading = ref(false)
  const error = ref(null)
  const lastSync = ref(null)

  // 計算屬性
  const totalPrice = computed(() => {
    return items.value.reduce((total, item) => {
      return total + (item.price * item.quantity)
    }, 0)
  })

  const itemCount = computed(() => {
    return items.value.reduce((count, item) => {
      return count + item.quantity
    }, 0)
  })

  const isEmpty = computed(() => {
    return items.value.length === 0
  })

  const hasErrors = computed(() => {
    return error.value !== null
  })

  // 方法
  const setLoading = (value) => {
    loading.value = value
  }

  const setError = (err) => {
    error.value = err
  }

  const clearError = () => {
    error.value = null
  }

  // 獲取購物車
  const getCart = async () => {
    setLoading(true)
    clearError()
    
    try {
      const response = await api.get('/api/cart')
      items.value = response.data.items || []
      lastSync.value = new Date()
      return response.data
    } catch (err) {
      const errorMessage = err.response?.data?.error || '獲取購物車失敗'
      setError(errorMessage)
      throw new Error(errorMessage)
    } finally {
      setLoading(false)
    }
  }

  // 添加商品到購物車
  const addToCart = async (productId, quantity = 1) => {
    setLoading(true)
    clearError()
    
    try {
      const response = await api.post('/api/cart/items', {
        product_id: productId,
        quantity: quantity
      })
      
      // 更新本地狀態而不是重新獲取購物車
      const existingItem = items.value.find(item => item.product_id === productId)
      if (existingItem) {
        existingItem.quantity += quantity
      } else {
        // 如果商品不在購物車中，需要獲取商品信息
        // 這裡暫時添加一個基本結構，實際應該從商品列表獲取
        items.value.push({
          product_id: productId,
          quantity: quantity,
          price: 0, // 暫時設為0，實際應該從商品信息獲取
          product: null // 暫時設為null，實際應該從商品信息獲取
        })
      }
      
      // 購物車摘要會通過計算屬性自動更新
      
      return response.data
    } catch (err) {
      const errorMessage = err.response?.data?.error || '添加商品到購物車失敗'
      setError(errorMessage)
      throw new Error(errorMessage)
    } finally {
      setLoading(false)
    }
  }

  // 更新購物車商品數量
  const updateQuantity = async (productId, quantity) => {
    setLoading(true)
    clearError()
    
    try {
      const response = await api.put(`/api/cart/items/${productId}`, {
        quantity: quantity
      })
      
      // 更新本地狀態
      const itemIndex = items.value.findIndex(item => item.product_id === productId)
      if (itemIndex !== -1) {
        if (quantity === 0) {
          items.value.splice(itemIndex, 1)
        } else {
          items.value[itemIndex].quantity = quantity
        }
      }
      
      return response.data
    } catch (err) {
      const errorMessage = err.response?.data?.error || '更新購物車失敗'
      setError(errorMessage)
      throw new Error(errorMessage)
    } finally {
      setLoading(false)
    }
  }

  // 從購物車移除商品
  const removeFromCart = async (productId) => {
    setLoading(true)
    clearError()
    
    try {
      const response = await api.delete(`/api/cart/items/${productId}`)
      
      // 更新本地狀態
      const itemIndex = items.value.findIndex(item => item.product_id === productId)
      if (itemIndex !== -1) {
        items.value.splice(itemIndex, 1)
      }
      
      return response.data
    } catch (err) {
      const errorMessage = err.response?.data?.error || '移除商品失敗'
      setError(errorMessage)
      throw new Error(errorMessage)
    } finally {
      setLoading(false)
    }
  }

  // 清空購物車
  const clearCart = async () => {
    setLoading(true)
    clearError()
    
    try {
      const response = await api.delete('/api/cart')
      
      // 清空本地狀態
      items.value = []
      
      return response.data
    } catch (err) {
      const errorMessage = err.response?.data?.error || '清空購物車失敗'
      setError(errorMessage)
      throw new Error(errorMessage)
    } finally {
      setLoading(false)
    }
  }

  // 獲取購物車摘要
  const getCartSummary = async () => {
    setLoading(true)
    clearError()
    
    try {
      const response = await api.get('/api/cart/summary')
      return response.data
    } catch (err) {
      const errorMessage = err.response?.data?.error || '獲取購物車摘要失敗'
      setError(errorMessage)
      throw new Error(errorMessage)
    } finally {
      setLoading(false)
    }
  }

  // 獲取購物車商品數量
  const getCartItemCount = async () => {
    try {
      const response = await api.get('/api/cart/count')
      return response.data.item_count
    } catch (err) {
      console.error('獲取購物車商品數量失敗:', err)
      return 0
    }
  }


  // 檢查商品是否在購物車中
  const isInCart = (productId) => {
    return items.value.some(item => item.product_id === productId)
  }

  // 獲取購物車中特定商品的數量
  const getItemQuantity = (productId) => {
    const item = items.value.find(item => item.product_id === productId)
    return item ? item.quantity : 0
  }




  return {
    // 狀態
    items,
    loading,
    error,
    lastSync,
    
    // 計算屬性
    totalPrice,
    itemCount,
    isEmpty,
    hasErrors,
    
    // 方法
    setLoading,
    setError,
    clearError,
    getCart,
    addToCart,
    updateQuantity,
    removeFromCart,
    clearCart,
    getCartSummary,
    getCartItemCount,
    isInCart,
    getItemQuantity
  }
})
