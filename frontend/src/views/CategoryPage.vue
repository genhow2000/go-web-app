<template>
  <div class="category-page">
    <Header />
    
    <div class="container">
      <div class="page-header">
        <h1>{{ categoryName }} 商品</h1>
        <p>探索 {{ categoryName }} 分類下的精選商品</p>
      </div>

      <div v-if="loading" class="loading">
        <div class="spinner"></div>
        <p>載入商品中...</p>
      </div>
      
      <div v-else-if="products.length === 0" class="empty-state">
        <div class="empty-icon">📦</div>
        <h3>暫無商品</h3>
        <p>此分類下暫時沒有商品，請查看其他分類</p>
        <router-link to="/" class="btn btn-primary">
          返回首頁
        </router-link>
      </div>

      <div v-else class="products-section">
        <div class="products-grid">
          <ProductCard 
            v-for="product in products" 
            :key="product.id"
            :product="product"
            @view="viewProduct"
            @add-to-cart="addToCart"
            @toggle-favorite="toggleFavorite"
          />
        </div>
        
        <div v-if="hasMore" class="load-more">
          <button @click="loadMore" class="btn btn-outline" :disabled="loadingMore">
            {{ loadingMore ? '載入中...' : '載入更多' }}
          </button>
        </div>
      </div>
    </div>
    
    <Footer />
  </div>
</template>

<script>
import { ref, onMounted, computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import Header from '@/components/common/Header.vue'
import Footer from '@/components/common/Footer.vue'
import ProductCard from '@/components/product/ProductCard.vue'
import api from '@/services/api'

export default {
  name: 'CategoryPage',
  components: {
    Header,
    Footer,
    ProductCard
  },
  setup() {
    const route = useRoute()
    const router = useRouter()
    
    const products = ref([])
    const loading = ref(true)
    const loadingMore = ref(false)
    const currentPage = ref(1)
    const hasMore = ref(true)
    const limit = 12

    const categoryName = computed(() => {
      return decodeURIComponent(route.params.category || '')
    })

    // 載入商品
    const loadProducts = async (page = 1, append = false) => {
      if (page === 1) {
        loading.value = true
      } else {
        loadingMore.value = true
      }

      try {
        const offset = (page - 1) * limit
        const response = await api.get(`/api/products/category/${categoryName.value}?limit=${limit}&offset=${offset}`)
        
        if (append) {
          products.value = [...products.value, ...response.data]
        } else {
          products.value = response.data
        }
        
        // 檢查是否還有更多商品
        hasMore.value = response.data.length === limit
        currentPage.value = page
      } catch (error) {
        console.error('載入商品失敗:', error)
        if (!append) {
          products.value = []
        }
      } finally {
        loading.value = false
        loadingMore.value = false
      }
    }

    // 載入更多商品
    const loadMore = () => {
      loadProducts(currentPage.value + 1, true)
    }

    // 查看商品詳情
    const viewProduct = (productId) => {
      router.push(`/product/${productId}`)
    }

    // 加入購物車
    const addToCart = (productId) => {
      // TODO: 實現加入購物車邏輯
      alert('商品已加入購物車！')
    }

    // 切換收藏
    const toggleFavorite = (productId) => {
      // TODO: 實現收藏功能邏輯
      alert('已加入收藏！')
    }

    onMounted(() => {
      loadProducts()
    })

    return {
      products,
      loading,
      loadingMore,
      hasMore,
      categoryName,
      loadMore,
      viewProduct,
      addToCart,
      toggleFavorite
    }
  }
}
</script>

<style scoped>
.category-page {
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
  background: white;
  padding: 40px 20px;
  border-radius: 10px;
  box-shadow: 0 2px 10px rgba(0,0,0,0.1);
}

.page-header h1 {
  margin: 0 0 10px 0;
  color: #2d3748;
  font-size: 2.5rem;
}

.page-header p {
  margin: 0;
  color: #718096;
  font-size: 1.1rem;
}

.loading {
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

.empty-state {
  text-align: center;
  padding: 60px 20px;
  color: #718096;
  background: white;
  border-radius: 10px;
  box-shadow: 0 2px 10px rgba(0,0,0,0.1);
}

.empty-icon {
  font-size: 4rem;
  margin-bottom: 20px;
}

.empty-state h3 {
  margin-bottom: 10px;
  color: #4a5568;
}

.empty-state p {
  margin-bottom: 30px;
}

.products-section {
  background: white;
  border-radius: 10px;
  box-shadow: 0 2px 10px rgba(0,0,0,0.1);
  padding: 30px;
}

.products-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(280px, 1fr));
  gap: 30px;
  margin-bottom: 40px;
}

.load-more {
  text-align: center;
  padding-top: 20px;
  border-top: 1px solid #e2e8f0;
}

.btn {
  padding: 12px 24px;
  border-radius: 6px;
  text-decoration: none;
  font-size: 14px;
  font-weight: 500;
  border: none;
  cursor: pointer;
  transition: all 0.3s;
  display: inline-flex;
  align-items: center;
  gap: 8px;
}

.btn-primary {
  background: #667eea;
  color: white;
}

.btn-primary:hover {
  background: #5a6fd8;
}

.btn-outline {
  background: transparent;
  color: #667eea;
  border: 1px solid #667eea;
}

.btn-outline:hover:not(:disabled) {
  background: #667eea;
  color: white;
}

.btn-outline:disabled {
  opacity: 0.6;
  cursor: not-allowed;
}

@media (max-width: 768px) {
  .products-grid {
    grid-template-columns: repeat(auto-fill, minmax(250px, 1fr));
    gap: 20px;
  }
  
  .page-header h1 {
    font-size: 2rem;
  }
}
</style>
