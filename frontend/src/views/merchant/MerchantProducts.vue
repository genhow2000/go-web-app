<template>
  <div class="merchant-products">
    <Header />
    
    <div class="container">
      <div class="page-header">
        <h1>商品管理</h1>
        <router-link to="/merchant/products/create" class="btn btn-primary">
          <i class="fas fa-plus"></i>
          新增商品
        </router-link>
      </div>

      <div class="stats-cards">
        <div class="stat-card">
          <div class="stat-icon">📦</div>
          <div class="stat-info">
            <div class="stat-number">{{ productStats.total }}</div>
            <div class="stat-label">總商品數</div>
          </div>
        </div>
        <div class="stat-card">
          <div class="stat-icon">✅</div>
          <div class="stat-info">
            <div class="stat-number">{{ productStats.active }}</div>
            <div class="stat-label">上架商品</div>
          </div>
        </div>
        <div class="stat-card">
          <div class="stat-icon">⏸️</div>
          <div class="stat-info">
            <div class="stat-number">{{ productStats.inactive }}</div>
            <div class="stat-label">下架商品</div>
          </div>
        </div>
      </div>

      <div class="filters">
        <div class="filter-group">
          <label>狀態篩選：</label>
          <select v-model="filters.status" @change="loadProducts">
            <option value="">全部</option>
            <option value="active">上架</option>
            <option value="inactive">下架</option>
          </select>
        </div>
        <div class="filter-group">
          <label>搜尋：</label>
          <input 
            v-model="filters.search"
            type="text" 
            placeholder="搜尋商品名稱..."
            @input="searchProducts"
          >
        </div>
      </div>

      <div class="products-table">
        <div v-if="loading" class="loading">
          <div class="spinner"></div>
          <p>載入商品中...</p>
        </div>
        
        <div v-else-if="products.length === 0" class="empty-state">
          <div class="empty-icon">📦</div>
          <h3>暫無商品</h3>
          <p>您還沒有添加任何商品，點擊上方按鈕開始添加商品吧！</p>
          <router-link to="/merchant/products/create" class="btn btn-primary">
            新增商品
          </router-link>
        </div>

        <div v-else class="products-grid">
          <div v-for="product in products" :key="product.id" class="product-card">
            <div class="product-image">
              <img 
                v-if="product.image_url" 
                :src="product.image_url" 
                :alt="product.name"
                class="product-img"
                @error="handleImageError"
              >
              <div v-else class="image-placeholder">📦</div>
            </div>
            <div class="product-info">
              <h3 class="product-name">{{ product.name }}</h3>
              <p class="product-description">{{ product.description }}</p>
              <div class="product-meta">
                <span class="product-price">NT$ {{ product.price.toLocaleString() }}</span>
                <span :class="['product-status', product.is_active ? 'active' : 'inactive']">
                  {{ product.is_active ? '上架' : '下架' }}
                </span>
              </div>
            </div>
            <div class="product-actions">
              <router-link 
                :to="`/merchant/products/${product.id}/edit`" 
                class="btn btn-sm btn-outline"
              >
                編輯
              </router-link>
              <button 
                @click="toggleProductStatus(product)"
                :class="['btn', 'btn-sm', product.is_active ? 'btn-warning' : 'btn-success']"
              >
                {{ product.is_active ? '下架' : '上架' }}
              </button>
              <button 
                @click="deleteProduct(product.id)"
                class="btn btn-sm btn-danger"
              >
                刪除
              </button>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import { ref, reactive, onMounted } from 'vue'
import Header from '@/components/common/Header.vue'
import api from '@/services/api'

export default {
  name: 'MerchantProducts',
  components: {
    Header
  },
  setup() {
    const products = ref([])
    const loading = ref(false)
    const productStats = ref({
      total: 0,
      active: 0,
      inactive: 0
    })

    const filters = reactive({
      status: '',
      search: ''
    })

    const loadProducts = async () => {
      loading.value = true
      try {
        const params = new URLSearchParams()
        if (filters.status) params.append('status', filters.status)
        if (filters.search) params.append('search', filters.search)
        
        const response = await api.get(`/merchant/api/products?${params}`)
        products.value = response.data.products || response.data
      } catch (error) {
        console.error('載入商品失敗:', error)
        // 使用模擬數據
        products.value = [
          {
            id: 1,
            name: '測試商品 1',
            description: '這是一個測試商品',
            price: 2999,
            is_active: true
          },
          {
            id: 2,
            name: '測試商品 2',
            description: '這是另一個測試商品',
            price: 4999,
            is_active: false
          }
        ]
      } finally {
        loading.value = false
      }
    }

    const loadStats = async () => {
      try {
        const response = await api.get('/merchant/api/products/stats')
        productStats.value = response.data
      } catch (error) {
        console.error('載入統計數據失敗:', error)
        // 使用模擬數據
        productStats.value = {
          total: 2,
          active: 1,
          inactive: 1
        }
      }
    }

    const searchProducts = () => {
      loadProducts()
    }

    const toggleProductStatus = async (product) => {
      try {
        await api.put(`/merchant/api/products/${product.id}/toggle-status`)
        product.is_active = !product.is_active
        loadStats() // 重新載入統計數據
      } catch (error) {
        console.error('切換商品狀態失敗:', error)
        alert('操作失敗，請重試')
      }
    }

    const handleImageError = (event) => {
      // 當圖片載入失敗時，隱藏圖片並顯示佔位符
      event.target.style.display = 'none'
      const placeholder = event.target.nextElementSibling
      if (placeholder) {
        placeholder.style.display = 'flex'
      }
    }

    const deleteProduct = async (productId) => {
      if (confirm('確定要刪除這個商品嗎？此操作無法復原。')) {
        try {
          await api.delete(`/merchant/api/products/${productId}`)
          loadProducts()
          loadStats()
        } catch (error) {
          console.error('刪除商品失敗:', error)
          alert('刪除失敗，請重試')
        }
      }
    }

    onMounted(() => {
      loadProducts()
      loadStats()
    })

    return {
      products,
      loading,
      productStats,
      filters,
      loadProducts,
      searchProducts,
      toggleProductStatus,
      handleImageError,
      deleteProduct
    }
  }
}
</script>

<style scoped>
.merchant-products {
  min-height: 100vh;
  background: #f8f9fa;
}

.container {
  max-width: 1200px;
  margin: 0 auto;
  padding: 20px;
}

.page-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 30px;
  background: white;
  padding: 20px;
  border-radius: 10px;
  box-shadow: 0 2px 10px rgba(0,0,0,0.1);
}

.page-header h1 {
  margin: 0;
  color: #2d3748;
}

.stats-cards {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
  gap: 20px;
  margin-bottom: 30px;
}

.stat-card {
  background: white;
  padding: 20px;
  border-radius: 10px;
  box-shadow: 0 2px 10px rgba(0,0,0,0.1);
  display: flex;
  align-items: center;
  gap: 15px;
}

.stat-icon {
  font-size: 2rem;
}

.stat-number {
  font-size: 1.8rem;
  font-weight: bold;
  color: #2d3748;
}

.stat-label {
  color: #718096;
  font-size: 0.9rem;
}

.filters {
  background: white;
  padding: 20px;
  border-radius: 10px;
  box-shadow: 0 2px 10px rgba(0,0,0,0.1);
  margin-bottom: 30px;
  display: flex;
  gap: 20px;
  align-items: center;
}

.filter-group {
  display: flex;
  align-items: center;
  gap: 10px;
}

.filter-group label {
  font-weight: 500;
  color: #4a5568;
}

.filter-group select,
.filter-group input {
  padding: 8px 12px;
  border: 1px solid #e2e8f0;
  border-radius: 6px;
  font-size: 14px;
}

.filter-group input {
  min-width: 200px;
}

.products-table {
  background: white;
  border-radius: 10px;
  box-shadow: 0 2px 10px rgba(0,0,0,0.1);
  padding: 20px;
}

.loading {
  text-align: center;
  padding: 40px;
  color: #718096;
}

.empty-state {
  text-align: center;
  padding: 60px 20px;
  color: #718096;
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

.products-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(300px, 1fr));
  gap: 20px;
}

.product-card {
  border: 1px solid #e2e8f0;
  border-radius: 10px;
  padding: 20px;
  transition: all 0.3s;
}

.product-card:hover {
  box-shadow: 0 4px 15px rgba(0,0,0,0.1);
  transform: translateY(-2px);
}

.product-image {
  width: 100%;
  height: 150px;
  background: #f7fafc;
  border-radius: 8px;
  display: flex;
  align-items: center;
  justify-content: center;
  margin-bottom: 15px;
}

.image-placeholder {
  font-size: 3rem;
  color: #a0aec0;
}

.product-img {
  width: 100%;
  height: 100%;
  object-fit: cover;
  border-radius: 8px;
}

.product-name {
  font-size: 1.1rem;
  font-weight: 600;
  color: #2d3748;
  margin-bottom: 8px;
}

.product-description {
  color: #718096;
  font-size: 0.9rem;
  margin-bottom: 15px;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
}

.product-meta {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 15px;
}

.product-price {
  font-size: 1.2rem;
  font-weight: bold;
  color: #e53e3e;
}

.product-status {
  padding: 4px 8px;
  border-radius: 12px;
  font-size: 0.8rem;
  font-weight: 500;
}

.product-status.active {
  background: #c6f6d5;
  color: #22543d;
}

.product-status.inactive {
  background: #fed7d7;
  color: #742a2a;
}

.product-actions {
  display: flex;
  gap: 8px;
  flex-wrap: wrap;
}

.btn {
  padding: 8px 16px;
  border-radius: 6px;
  text-decoration: none;
  font-size: 0.9rem;
  font-weight: 500;
  border: none;
  cursor: pointer;
  transition: all 0.3s;
  display: inline-flex;
  align-items: center;
  gap: 5px;
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

.btn-outline:hover {
  background: #667eea;
  color: white;
}

.btn-success {
  background: #48bb78;
  color: white;
}

.btn-success:hover {
  background: #38a169;
}

.btn-warning {
  background: #ed8936;
  color: white;
}

.btn-warning:hover {
  background: #dd6b20;
}

.btn-danger {
  background: #e53e3e;
  color: white;
}

.btn-danger:hover {
  background: #c53030;
}

.btn-sm {
  padding: 6px 12px;
  font-size: 0.8rem;
}

@media (max-width: 768px) {
  .page-header {
    flex-direction: column;
    gap: 15px;
    align-items: flex-start;
  }
  
  .filters {
    flex-direction: column;
    align-items: flex-start;
  }
  
  .filter-group input {
    min-width: 100%;
  }
  
  .products-grid {
    grid-template-columns: 1fr;
  }
  
  .product-actions {
    justify-content: center;
  }
}
</style>
