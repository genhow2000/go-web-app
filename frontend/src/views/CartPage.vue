<template>
  <div class="cart-page">
    <Header />
    
    <div class="container">
      <!-- 頁面標題 -->
      <div class="page-header">
        <h1>購物車</h1>
        <p class="page-subtitle">查看和管理您的購物車商品</p>
      </div>

      <!-- 加載狀態 -->
      <div v-if="loading" class="loading-container">
        <div class="spinner"></div>
        <p>載入購物車中...</p>
      </div>

      <!-- 錯誤狀態 -->
      <div v-else-if="hasErrors" class="error-container">
        <div class="error-icon">⚠️</div>
        <h3>載入購物車失敗</h3>
        <p>{{ error }}</p>
        <button @click="retryLoadCart" class="btn-retry">重試</button>
      </div>

      <!-- 空購物車 -->
      <div v-else-if="isEmpty" class="empty-cart">
        <div class="empty-icon">🛒</div>
        <h3>購物車是空的</h3>
        <p>您還沒有添加任何商品到購物車</p>
        <router-link to="/" class="btn-shopping">開始購物</router-link>
      </div>

      <!-- 購物車內容 -->
      <div v-else class="cart-content">
        <div class="cart-main">
          <!-- 購物車商品列表 -->
          <div class="cart-items">
            <h2>購物車商品 ({{ itemCount }} 件)</h2>
            <div class="items-list">
              <CartItem
                v-for="item in items"
                :key="item.id"
                :item="item"
                @update-quantity="handleUpdateQuantity"
                @remove-item="handleRemoveItem"
              />
            </div>
          </div>
        </div>

        <!-- 購物車摘要 -->
        <div class="cart-sidebar">
          <CartSummary
            :items="items"
            :total-price="totalPrice"
            :item-count="itemCount"
            @clear-cart="handleClearCart"
            @checkout="handleCheckout"
          />
        </div>
      </div>
    </div>

    <Footer />
  </div>
</template>

<script>
import { ref, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useCartStore } from '@/stores/cart'
import { useAuthStore } from '@/stores/auth'
import Header from '@/components/common/Header.vue'
import Footer from '@/components/common/Footer.vue'
import CartItem from '@/components/cart/CartItem.vue'
import CartSummary from '@/components/cart/CartSummary.vue'

export default {
  name: 'CartPage',
  components: {
    Header,
    Footer,
    CartItem,
    CartSummary
  },
  setup() {
    const router = useRouter()
    const cartStore = useCartStore()
    const authStore = useAuthStore()

    // 計算屬性
    const items = computed(() => cartStore.items)
    const totalPrice = computed(() => cartStore.totalPrice)
    const itemCount = computed(() => cartStore.itemCount)
    const loading = computed(() => cartStore.loading)
    const error = computed(() => cartStore.error)
    const isEmpty = computed(() => cartStore.isEmpty)
    const hasErrors = computed(() => cartStore.hasErrors)

    // 載入購物車
    const loadCart = async () => {
      try {
        await cartStore.getCart()
      } catch (err) {
        console.error('載入購物車失敗:', err)
        // 如果是認證錯誤，不要重新拋出，避免觸發登出
        if (err.response?.status === 401) {
          console.log('認證失敗，跳轉到登入頁面')
          router.push('/customer/login')
          return
        }
      }
    }

    // 重試載入購物車
    const retryLoadCart = () => {
      loadCart()
    }

    // 更新商品數量
    const handleUpdateQuantity = async (productId, quantity) => {
      try {
        await cartStore.updateQuantity(productId, quantity)
      } catch (err) {
        console.error('更新商品數量失敗:', err)
        // 可以添加錯誤提示
      }
    }

    // 移除商品
    const handleRemoveItem = async (productId) => {
      try {
        await cartStore.removeFromCart(productId)
      } catch (err) {
        console.error('移除商品失敗:', err)
        // 可以添加錯誤提示
      }
    }

    // 清空購物車
    const handleClearCart = async () => {
      if (confirm('確定要清空購物車嗎？')) {
        try {
          await cartStore.clearCart()
        } catch (err) {
          console.error('清空購物車失敗:', err)
          // 可以添加錯誤提示
        }
      }
    }

    // 結算
    const handleCheckout = () => {
      // 檢查是否已登入
      if (!authStore.isAuthenticated) {
        router.push('/customer/login')
        return
      }

      // 這裡可以跳轉到結算頁面
      // router.push('/checkout')
      alert('結算功能尚未實現')
    }

    // 組件掛載時載入購物車
    onMounted(() => {
      // 檢查是否已登入
      if (!authStore.isAuthenticated) {
        router.push('/customer/login')
        return
      }
      
      loadCart()
    })

    return {
      items,
      totalPrice,
      itemCount,
      loading,
      error,
      isEmpty,
      hasErrors,
      loadCart,
      retryLoadCart,
      handleUpdateQuantity,
      handleRemoveItem,
      handleClearCart,
      handleCheckout
    }
  }
}
</script>

<style scoped>
.cart-page {
  min-height: 100vh;
  background: #f8f9fa;
}

.container {
  max-width: 1200px;
  margin: 0 auto;
  padding: 20px;
}

.page-header {
  text-align: center;
  margin-bottom: 40px;
}

.page-header h1 {
  font-size: 2.5rem;
  color: #2d3748;
  margin-bottom: 10px;
}

.page-subtitle {
  color: #718096;
  font-size: 1.1rem;
}

.loading-container {
  text-align: center;
  padding: 60px 20px;
  color: #718096;
}

.spinner {
  display: inline-block;
  width: 40px;
  height: 40px;
  border: 4px solid #f3f3f3;
  border-top: 4px solid #667eea;
  border-radius: 50%;
  animation: spin 1s linear infinite;
  margin-bottom: 20px;
}

@keyframes spin {
  0% { transform: rotate(0deg); }
  100% { transform: rotate(360deg); }
}

.error-container {
  text-align: center;
  padding: 60px 20px;
  background: white;
  border-radius: 10px;
  box-shadow: 0 2px 10px rgba(0,0,0,0.1);
}

.error-icon {
  font-size: 4rem;
  margin-bottom: 20px;
}

.error-container h3 {
  color: #e53e3e;
  margin-bottom: 10px;
}

.btn-retry {
  background: #667eea;
  color: white;
  border: none;
  padding: 12px 24px;
  border-radius: 6px;
  cursor: pointer;
  font-size: 1rem;
  margin-top: 20px;
}

.btn-retry:hover {
  background: #5a67d8;
}

.empty-cart {
  text-align: center;
  padding: 80px 20px;
  background: white;
  border-radius: 10px;
  box-shadow: 0 2px 10px rgba(0,0,0,0.1);
}

.empty-icon {
  font-size: 5rem;
  margin-bottom: 20px;
}

.empty-cart h3 {
  color: #4a5568;
  margin-bottom: 10px;
  font-size: 1.5rem;
}

.empty-cart p {
  color: #718096;
  margin-bottom: 30px;
}

.btn-shopping {
  display: inline-block;
  background: #667eea;
  color: white;
  text-decoration: none;
  padding: 12px 24px;
  border-radius: 6px;
  font-size: 1rem;
  transition: background 0.3s;
}

.btn-shopping:hover {
  background: #5a67d8;
}

.cart-content {
  display: grid;
  grid-template-columns: 1fr 350px;
  gap: 30px;
}

.cart-main {
  background: white;
  border-radius: 10px;
  padding: 30px;
  box-shadow: 0 2px 10px rgba(0,0,0,0.1);
}

.cart-items h2 {
  color: #2d3748;
  margin-bottom: 20px;
  font-size: 1.5rem;
}

.items-list {
  display: flex;
  flex-direction: column;
  gap: 20px;
}

.cart-sidebar {
  position: sticky;
  top: 20px;
  height: fit-content;
}

/* 響應式設計 */
@media (max-width: 768px) {
  .cart-content {
    grid-template-columns: 1fr;
    gap: 20px;
  }
  
  .cart-sidebar {
    position: static;
  }
  
  .container {
    padding: 15px;
  }
  
  .page-header h1 {
    font-size: 2rem;
  }
}
</style>
