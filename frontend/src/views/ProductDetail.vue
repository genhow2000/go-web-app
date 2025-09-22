<template>
  <div class="product-detail">
    <Header />
    
    <div class="container">
      <div v-if="loading" class="loading">
        <div class="spinner"></div>
        <p>載入商品詳情中...</p>
      </div>
      
      <div v-else-if="!product" class="error-state">
        <div class="error-icon">❌</div>
        <h3>商品不存在</h3>
        <p>找不到指定的商品，請檢查商品ID是否正確。</p>
        <router-link to="/" class="btn btn-primary">
          返回首頁
        </router-link>
      </div>

      <div v-else class="product-content">
        <!-- 麵包屑導航 -->
        <nav class="breadcrumb">
          <router-link to="/">首頁</router-link>
          <span class="separator">></span>
          <router-link :to="`/category/${encodeURIComponent(product.category)}`">
            {{ product.category }}
          </router-link>
          <span class="separator">></span>
          <span class="current">{{ product.name }}</span>
        </nav>

        <div class="product-main">
          <!-- 商品圖片區域 -->
          <div class="product-gallery">
            <div class="main-image">
              <img 
                v-if="product.image_url" 
                :src="product.image_url" 
                :alt="product.name"
                class="product-image"
                @error="handleImageError"
              >
              <div v-else class="image-placeholder">
                <div class="placeholder-icon">📦</div>
                <p>暫無圖片</p>
              </div>
            </div>
            
            <!-- 商品標籤 -->
            <div class="product-badges">
              <span v-if="product.is_featured" class="badge badge-featured">精選商品</span>
              <span v-if="product.is_on_sale" class="badge badge-sale">特價商品</span>
              <span v-if="!product.is_active" class="badge badge-inactive">已下架</span>
            </div>
          </div>

          <!-- 商品信息區域 -->
          <div class="product-info">
            <h1 class="product-title">{{ product.name }}</h1>
            
            <div class="product-meta">
              <div class="category">
                <span class="label">分類：</span>
                <router-link :to="`/category/${encodeURIComponent(product.category)}`" class="category-link">
                  {{ product.category }}
                </router-link>
              </div>
              <div v-if="product.brand" class="brand">
                <span class="label">品牌：</span>
                <span>{{ product.brand }}</span>
              </div>
              <div v-if="product.sku" class="sku">
                <span class="label">商品編號：</span>
                <span>{{ product.sku }}</span>
              </div>
            </div>

            <div class="product-description">
              <h3>商品描述</h3>
              <p>{{ product.description }}</p>
            </div>

            <div class="product-stats">
              <div class="stat-item">
                <span class="stat-label">瀏覽次數</span>
                <span class="stat-value">{{ product.view_count || 0 }}</span>
              </div>
              <div class="stat-item">
                <span class="stat-label">銷售數量</span>
                <span class="stat-value">{{ product.sales_count || 0 }}</span>
              </div>
              <div v-if="product.rating > 0" class="stat-item">
                <span class="stat-label">評分</span>
                <span class="stat-value">
                  <span class="rating">{{ '★'.repeat(Math.floor(product.rating)) }}{{ '☆'.repeat(5 - Math.floor(product.rating)) }}</span>
                  ({{ product.review_count || 0 }} 評價)
                </span>
              </div>
            </div>
          </div>

          <!-- 購買區域 -->
          <div class="purchase-section">
            <div class="price-section">
              <div class="current-price">
                <span class="currency">NT$</span>
                <span class="amount">{{ product.price.toLocaleString() }}</span>
              </div>
              <div v-if="product.original_price && product.original_price > product.price" class="original-price">
                <span class="currency">NT$</span>
                <span class="amount">{{ product.original_price.toLocaleString() }}</span>
                <span class="discount">
                  省 {{ ((product.original_price - product.price) / product.original_price * 100).toFixed(0) }}%
                </span>
              </div>
            </div>

            <div class="stock-info">
              <span class="label">庫存：</span>
              <span :class="['stock', product.stock > 0 ? 'in-stock' : 'out-of-stock']">
                {{ product.stock > 0 ? `${product.stock} 件` : '缺貨' }}
              </span>
            </div>

            <div class="quantity-selector">
              <label for="quantity">數量：</label>
              <div class="quantity-controls">
                <button @click="decreaseQuantity" :disabled="quantity <= 1" class="btn-quantity">-</button>
                <input 
                  v-model.number="quantity" 
                  type="number" 
                  id="quantity"
                  min="1" 
                  :max="product.stock"
                  class="quantity-input"
                >
                <button @click="increaseQuantity" :disabled="quantity >= product.stock" class="btn-quantity">+</button>
              </div>
            </div>

            <div class="action-buttons">
              <AddToCartButton 
                :product="product"
                :show-quantity-selector="false"
                variant="primary"
                @added-to-cart="handleAddedToCart"
                @error="handleCartError"
              />
              
              <button 
                @click="toggleFavorite" 
                class="btn btn-outline btn-large"
                :class="{ 'favorited': isFavorited }"
              >
                <span class="btn-icon">{{ isFavorited ? '❤️' : '🤍' }}</span>
                {{ isFavorited ? '已收藏' : '收藏' }}
              </button>
            </div>

            <div class="product-features">
              <div class="feature-item">
                <span class="feature-icon">🚚</span>
                <span>免費配送</span>
              </div>
              <div class="feature-item">
                <span class="feature-icon">🔄</span>
                <span>7天退換</span>
              </div>
              <div class="feature-item">
                <span class="feature-icon">🛡️</span>
                <span>品質保證</span>
              </div>
            </div>
          </div>
        </div>

        <!-- 相關商品推薦 -->
        <div class="related-products">
          <h2>相關商品推薦</h2>
          <div v-if="relatedLoading" class="loading">
            <div class="spinner"></div>
            <p>載入相關商品中...</p>
          </div>
          <div v-else-if="relatedProducts.length > 0" class="products-grid">
            <ProductCard 
              v-for="relatedProduct in relatedProducts" 
              :key="relatedProduct.id"
              :product="relatedProduct"
              @view="viewProduct"
              @add-to-cart="addToCart"
              @toggle-favorite="toggleFavorite"
            />
          </div>
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
import AddToCartButton from '@/components/cart/AddToCartButton.vue'
import api from '@/services/api'

export default {
  name: 'ProductDetail',
  components: {
    Header,
    Footer,
    ProductCard,
    AddToCartButton
  },
  setup() {
    const route = useRoute()
    const router = useRouter()
    
    const product = ref(null)
    const relatedProducts = ref([])
    const loading = ref(true)
    const relatedLoading = ref(false)
    const quantity = ref(1)
    const isFavorited = ref(false)

    // 載入商品詳情
    const loadProduct = async () => {
      loading.value = true
      try {
        const productId = route.params.id
        const response = await api.get(`/api/products/${productId}`)
        product.value = response.data
        
        // 載入相關商品
        loadRelatedProducts()
      } catch (error) {
        console.error('載入商品失敗:', error)
        product.value = null
      } finally {
        loading.value = false
      }
    }

    // 載入相關商品
    const loadRelatedProducts = async () => {
      if (!product.value) return
      
      relatedLoading.value = true
      try {
        const response = await api.get(`/api/products/category/${encodeURIComponent(product.value.category)}?limit=4`)
        // 過濾掉當前商品
        relatedProducts.value = response.data.filter(p => p.id !== product.value.id)
      } catch (error) {
        console.error('載入相關商品失敗:', error)
        relatedProducts.value = []
      } finally {
        relatedLoading.value = false
      }
    }

    // 圖片錯誤處理
    const handleImageError = (event) => {
      event.target.style.display = 'none'
      const placeholder = event.target.nextElementSibling
      if (placeholder) {
        placeholder.style.display = 'flex'
      }
    }

    // 數量控制
    const increaseQuantity = () => {
      if (quantity.value < product.value.stock) {
        quantity.value++
      }
    }

    const decreaseQuantity = () => {
      if (quantity.value > 1) {
        quantity.value--
      }
    }

    // 加入購物車（保留原有方法以備用）
    const addToCart = () => {
      if (!product.value.is_active || product.value.stock <= 0) {
        alert('商品無法購買')
        return
      }
      
      // TODO: 實現加入購物車邏輯
      alert(`已將 ${quantity.value} 件「${product.value.name}」加入購物車！`)
    }

    // 處理購物車按鈕事件
    const handleAddedToCart = (data) => {
      console.log('商品已加入購物車:', data)
      // 可以在這裡添加成功提示
    }

    const handleCartError = (error) => {
      console.error('購物車錯誤:', error)
      // 可以在這裡添加錯誤提示
    }

    // 切換收藏
    const toggleFavorite = () => {
      isFavorited.value = !isFavorited.value
      // TODO: 實現收藏功能邏輯
      alert(isFavorited.value ? '已加入收藏！' : '已取消收藏！')
    }

    // 查看商品詳情
    const viewProduct = (productId) => {
      router.push(`/product/${productId}`)
    }

    onMounted(() => {
      loadProduct()
    })

    return {
      product,
      relatedProducts,
      loading,
      relatedLoading,
      quantity,
      isFavorited,
      handleImageError,
      increaseQuantity,
      decreaseQuantity,
      addToCart,
      handleAddedToCart,
      handleCartError,
      toggleFavorite,
      viewProduct
    }
  }
}
</script>

<style scoped>
.product-detail {
  min-height: 100vh;
  background: #f8f9fa;
}

.container {
  max-width: 1200px;
  margin: 0 auto;
  padding: 20px;
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

.error-state {
  text-align: center;
  padding: 60px 20px;
  color: #718096;
  background: white;
  border-radius: 10px;
  box-shadow: 0 2px 10px rgba(0,0,0,0.1);
}

.error-icon {
  font-size: 4rem;
  margin-bottom: 20px;
}

.error-state h3 {
  margin-bottom: 10px;
  color: #4a5568;
}

.breadcrumb {
  margin-bottom: 20px;
  font-size: 14px;
  color: #718096;
}

.breadcrumb a {
  color: #667eea;
  text-decoration: none;
}

.breadcrumb a:hover {
  text-decoration: underline;
}

.separator {
  margin: 0 8px;
}

.current {
  color: #4a5568;
  font-weight: 500;
}

.product-main {
  display: grid;
  grid-template-columns: 1fr 1fr 1fr;
  gap: 40px;
  margin-bottom: 60px;
  background: white;
  border-radius: 10px;
  box-shadow: 0 2px 10px rgba(0,0,0,0.1);
  padding: 30px;
}

.product-gallery {
  position: relative;
}

.main-image {
  width: 100%;
  height: 400px;
  border-radius: 8px;
  overflow: hidden;
  background: #f7fafc;
  display: flex;
  align-items: center;
  justify-content: center;
  margin-bottom: 20px;
}

.product-image {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.image-placeholder {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  color: #a0aec0;
}

.placeholder-icon {
  font-size: 3rem;
  margin-bottom: 10px;
}

.product-badges {
  display: flex;
  gap: 10px;
  flex-wrap: wrap;
}

.badge {
  padding: 4px 12px;
  border-radius: 20px;
  font-size: 12px;
  font-weight: 500;
}

.badge-featured {
  background: #e6fffa;
  color: #234e52;
}

.badge-sale {
  background: #fed7d7;
  color: #742a2a;
}

.badge-inactive {
  background: #e2e8f0;
  color: #4a5568;
}

.product-info {
  padding: 0 20px;
}

.product-title {
  font-size: 2rem;
  font-weight: 700;
  color: #2d3748;
  margin-bottom: 20px;
  line-height: 1.2;
}

.product-meta {
  margin-bottom: 30px;
}

.product-meta > div {
  margin-bottom: 8px;
  font-size: 14px;
}

.label {
  color: #718096;
  font-weight: 500;
}

.category-link {
  color: #667eea;
  text-decoration: none;
}

.category-link:hover {
  text-decoration: underline;
}

.product-description h3 {
  font-size: 1.2rem;
  color: #2d3748;
  margin-bottom: 10px;
}

.product-description p {
  color: #4a5568;
  line-height: 1.6;
  margin-bottom: 30px;
}

.product-stats {
  display: flex;
  gap: 20px;
  flex-wrap: wrap;
}

.stat-item {
  display: flex;
  flex-direction: column;
  align-items: center;
  padding: 15px;
  background: #f7fafc;
  border-radius: 8px;
  min-width: 80px;
}

.stat-label {
  font-size: 12px;
  color: #718096;
  margin-bottom: 5px;
}

.stat-value {
  font-size: 16px;
  font-weight: 600;
  color: #2d3748;
}

.rating {
  color: #f6ad55;
}

.purchase-section {
  border-left: 1px solid #e2e8f0;
  padding-left: 30px;
}

.price-section {
  margin-bottom: 20px;
}

.current-price {
  display: flex;
  align-items: baseline;
  margin-bottom: 10px;
}

.currency {
  font-size: 1.2rem;
  color: #e53e3e;
  font-weight: 500;
}

.amount {
  font-size: 2.5rem;
  color: #e53e3e;
  font-weight: 700;
  margin-left: 5px;
}

.original-price {
  display: flex;
  align-items: center;
  gap: 10px;
}

.original-price .currency,
.original-price .amount {
  font-size: 1rem;
  color: #a0aec0;
  text-decoration: line-through;
}

.discount {
  background: #fed7d7;
  color: #742a2a;
  padding: 2px 8px;
  border-radius: 12px;
  font-size: 12px;
  font-weight: 500;
}

.stock-info {
  margin-bottom: 20px;
  font-size: 14px;
}

.stock.in-stock {
  color: #38a169;
  font-weight: 500;
}

.stock.out-of-stock {
  color: #e53e3e;
  font-weight: 500;
}

.quantity-selector {
  margin-bottom: 30px;
}

.quantity-selector label {
  display: block;
  margin-bottom: 10px;
  font-weight: 500;
  color: #4a5568;
}

.quantity-controls {
  display: flex;
  align-items: center;
  gap: 10px;
}

.btn-quantity {
  width: 40px;
  height: 40px;
  border: 1px solid #e2e8f0;
  background: white;
  border-radius: 6px;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 18px;
  font-weight: 500;
  color: #4a5568;
}

.btn-quantity:hover:not(:disabled) {
  background: #f7fafc;
  border-color: #cbd5e0;
}

.btn-quantity:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.quantity-input {
  width: 80px;
  height: 40px;
  border: 1px solid #e2e8f0;
  border-radius: 6px;
  text-align: center;
  font-size: 16px;
  font-weight: 500;
}

.action-buttons {
  display: flex;
  flex-direction: column;
  gap: 15px;
  margin-bottom: 30px;
}

.btn {
  padding: 15px 24px;
  border-radius: 8px;
  font-size: 16px;
  font-weight: 500;
  border: none;
  cursor: pointer;
  transition: all 0.3s;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 8px;
}

.btn-large {
  padding: 18px 24px;
  font-size: 18px;
}

.btn-primary {
  background: #667eea;
  color: white;
}

.btn-primary:hover:not(:disabled) {
  background: #5a6fd8;
}

.btn-primary:disabled {
  background: #a0aec0;
  cursor: not-allowed;
}

.btn-outline {
  background: transparent;
  color: #667eea;
  border: 2px solid #667eea;
}

.btn-outline:hover {
  background: #667eea;
  color: white;
}

.btn-outline.favorited {
  background: #667eea;
  color: white;
}

.btn-icon {
  font-size: 18px;
}

.product-features {
  display: flex;
  flex-direction: column;
  gap: 10px;
}

.feature-item {
  display: flex;
  align-items: center;
  gap: 10px;
  font-size: 14px;
  color: #4a5568;
}

.feature-icon {
  font-size: 16px;
}

.related-products {
  background: white;
  border-radius: 10px;
  box-shadow: 0 2px 10px rgba(0,0,0,0.1);
  padding: 30px;
}

.related-products h2 {
  margin-bottom: 30px;
  color: #2d3748;
  font-size: 1.5rem;
}

.products-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(250px, 1fr));
  gap: 20px;
}

@media (max-width: 768px) {
  .product-main {
    grid-template-columns: 1fr;
    gap: 20px;
  }
  
  .purchase-section {
    border-left: none;
    border-top: 1px solid #e2e8f0;
    padding-left: 0;
    padding-top: 20px;
  }
  
  .action-buttons {
    flex-direction: column;
  }
  
  .products-grid {
    grid-template-columns: repeat(auto-fill, minmax(200px, 1fr));
  }
}
</style>
