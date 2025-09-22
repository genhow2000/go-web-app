<template>
  <div class="homepage">
    <Header />
    
    <!-- 英雄區域 -->
    <section class="hero" id="home">
      <div class="hero-content">
        <h1>DEMO發現精選商品</h1>
        <p>優質商品，優惠價格，讓您享受最好的購物體驗</p>
        
        <!-- AI 聊天提示 -->
        <div class="ai-chat-hint hero-hint">
          <div class="hint-content">
            <div class="hint-icon">🤖</div>
            <div class="hint-text">
              <h4>需要購物建議？</h4>
              <p>點擊左下角的 AI 助手（真實ＡＩ），我可以為您推薦商品、回答問題或協助您找到最適合的商品！</p>
            </div>
            <button class="hint-btn" @click="toggleChatWindow">立即諮詢</button>
          </div>
        </div>
        
        <div class="search-bar">
          <input 
            v-model="searchKeyword"
            type="text" 
            placeholder="搜尋您想要的商品..." 
            @keypress.enter="searchProducts"
          >
          <button @click="searchProducts">搜尋</button>
        </div>
      </div>
    </section>

    <!-- 分類導航 -->
    <section class="categories" id="categories">
      <div class="container">
        <h2 class="section-title">商品分類</h2>
        <div class="category-grid">
          <div v-if="loadingCategories" class="loading">
            <div class="spinner"></div>
            <p>載入分類中...</p>
          </div>
          <div 
            v-else
            v-for="category in categories" 
            :key="category"
            class="category-card"
            @click="filterByCategory(category)"
          >
            <div class="category-icon">{{ getCategoryIcon(category) }}</div>
            <div class="category-name">{{ category }}</div>
          </div>
        </div>
      </div>
    </section>

    <!-- 促銷橫幅 -->
    <section class="promo-banner">
      <div class="promo-content">
        <h2 class="promo-title">🎉 新年特惠活動</h2>
        <p class="promo-description">全場商品8折優惠，滿額免運費，限時優惠不容錯過！</p>
        <div class="promo-code">優惠碼：NEWYEAR2025</div>
      </div>
    </section>

    <!-- 精選商品 -->
    <section class="featured-products" id="products">
      <div class="container">
        <h2 class="section-title">精選商品</h2>
        <div class="product-grid">
          <div v-if="loadingProducts" class="loading">
            <div class="spinner"></div>
            <p>載入商品中...</p>
          </div>
          <ProductCard 
            v-else
            v-for="product in products" 
            :key="product.id"
            :product="product"
            @view="viewProduct"
            @toggle-favorite="toggleFavorite"
          />
        </div>
      </div>
    </section>

    <!-- 特色服務 -->
    <section class="features">
      <div class="container">
        <h2 class="section-title">為什麼選擇我們</h2>
        <div class="feature-grid">
          <div class="feature-item">
            <div class="feature-icon">🚚</div>
            <h3 class="feature-title">快速配送</h3>
            <p class="feature-description">24小時內發貨，3-5天送達，讓您快速收到心儀商品</p>
          </div>
          <div class="feature-item">
            <div class="feature-icon">🔒</div>
            <h3 class="feature-title">安全支付</h3>
            <p class="feature-description">多種支付方式，銀行級安全加密，保障您的資金安全</p>
          </div>
          <div class="feature-item">
            <div class="feature-icon">💎</div>
            <h3 class="feature-title">品質保證</h3>
            <p class="feature-description">嚴選優質商品，7天無理由退換，讓您買得放心</p>
          </div>
          <div class="feature-item">
            <div class="feature-icon">🎯</div>
            <h3 class="feature-title">精準推薦</h3>
            <p class="feature-description">基於AI的智能推薦，為您推薦最適合的商品</p>
          </div>
        </div>
      </div>
    </section>

    <!-- 頁腳 -->
    <Footer />

    <!-- AI 聊天按鈕 -->
    <button class="ai-chat-btn" @click="toggleChatWindow" title="AI 助手">
      🤖
      <span v-if="hasNewMessage" class="notification-badge">1</span>
    </button>

    <!-- 技術展示按鈕 -->
    <button class="tech-showcase-btn" @click="goToTechShowcase" title="技術展示">
      🚀
    </button>

    <!-- AI 聊天窗口 -->
    <AIChatWindow 
      v-if="showChatWindow"
      @close="toggleChatWindow"
    />
  </div>
</template>

<script>
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import Header from '@/components/common/Header.vue'
import Footer from '@/components/common/Footer.vue'
import ProductCard from '@/components/product/ProductCard.vue'
import AIChatWindow from '@/components/chat/AIChatWindow.vue'
import api from '@/services/api'

export default {
  name: 'HomePage',
  components: {
    Header,
    Footer,
    ProductCard,
    AIChatWindow
  },
  setup() {
    const router = useRouter()
    
    // 響應式數據
    const searchKeyword = ref('')
    const categories = ref([])
    const products = ref([])
    const loadingCategories = ref(false)
    const loadingProducts = ref(false)
    const showChatWindow = ref(false)
    const hasNewMessage = ref(false)

    // 分類圖標映射
    const categoryIcons = {
      '電子產品': '📱',
      '服飾': '👕',
      '家居': '🏠',
      '美妝': '💄',
      '運動': '⚽',
      '食品': '🍎',
      '圖書': '📚',
      '其他': '📦'
    }

    // 獲取分類圖標
    const getCategoryIcon = (category) => {
      return categoryIcons[category] || '📦'
    }

    // 載入分類
    const loadCategories = async () => {
      loadingCategories.value = true
      try {
        const response = await api.get('/api/categories')
        categories.value = response.data
      } catch (error) {
        console.error('載入分類失敗:', error)
      } finally {
        loadingCategories.value = false
      }
    }

    // 載入商品
    const loadProducts = async () => {
      loadingProducts.value = true
      try {
        const response = await api.get('/api/products/featured?limit=8')
        products.value = response.data
      } catch (error) {
        console.error('載入商品失敗:', error)
      } finally {
        loadingProducts.value = false
      }
    }

    // 搜尋商品
    const searchProducts = () => {
      if (searchKeyword.value.trim()) {
        router.push(`/search?q=${encodeURIComponent(searchKeyword.value)}`)
      }
    }

    // 按分類篩選
    const filterByCategory = (category) => {
      router.push(`/category/${encodeURIComponent(category)}`)
    }

    // 查看商品詳情
    const viewProduct = (productId) => {
      router.push(`/product/${productId}`)
    }


    // 切換收藏
    const toggleFavorite = (productId) => {
      // TODO: 實現收藏功能邏輯
      alert('已加入收藏！')
    }

    // 切換聊天窗口
    const toggleChatWindow = () => {
      showChatWindow.value = !showChatWindow.value
      if (showChatWindow.value) {
        hasNewMessage.value = false
      }
    }

    // 前往技術展示頁面
    const goToTechShowcase = () => {
      router.push('/tech-showcase')
    }

    // 組件掛載時載入數據
    onMounted(() => {
      loadCategories()
      loadProducts()
    })

    return {
      searchKeyword,
      categories,
      products,
      loadingCategories,
      loadingProducts,
      showChatWindow,
      hasNewMessage,
      getCategoryIcon,
      searchProducts,
      filterByCategory,
      viewProduct,
      toggleFavorite,
      toggleChatWindow,
      goToTechShowcase
    }
  }
}
</script>

<style scoped>
/* 這裡包含所有原有的CSS樣式 */
/* 為了簡潔，我將在單獨的CSS文件中定義 */
@import '@/assets/styles/homepage.css';
</style>
