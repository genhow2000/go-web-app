<template>
  <div class="stock-list-page">
    <!-- 頁面標題 -->
    <div class="page-header">
      <h1>台股列表</h1>
      <p class="page-subtitle">即時股票資訊與分頁瀏覽</p>
    </div>

    <!-- 篩選和搜尋區域 -->
    <div class="filter-section">
      <div class="filter-row">
        <!-- 搜尋框 -->
        <div class="search-box">
          <input 
            v-model="searchKeyword" 
            @input="onSearchInput"
            type="text" 
            placeholder="搜尋股票代碼或名稱..."
            class="search-input"
          >
          <button @click="performSearch" class="search-btn">搜尋</button>
        </div>

        <!-- 分類篩選 -->
        <div class="filter-group">
          <label>產業分類：</label>
          <select v-model="selectedCategory" @change="onFilterChange" class="filter-select">
            <option value="">全部</option>
            <option v-for="category in categories" :key="category.code" :value="category.code">
              {{ category.name }}
            </option>
          </select>
        </div>

        <!-- 市場篩選 -->
        <div class="filter-group">
          <label>市場：</label>
          <select v-model="selectedMarket" @change="onFilterChange" class="filter-select">
            <option value="">全部</option>
            <option value="TSE">上市</option>
            <option value="OTC">上櫃</option>
          </select>
        </div>

        <!-- 排序選項 -->
        <div class="filter-group">
          <label>排序：</label>
          <select v-model="sortBy" @change="onSortChange" class="filter-select">
            <option value="code">股票代碼</option>
            <option value="name">股票名稱</option>
            <option value="price">股價</option>
            <option value="change_percent">漲跌幅</option>
            <option value="volume">成交量</option>
          </select>
          <select v-model="sortOrder" @change="onSortChange" class="filter-select">
            <option value="asc">升序</option>
            <option value="desc">降序</option>
          </select>
        </div>
      </div>
    </div>

    <!-- 股票列表 -->
    <div class="stock-list-container">
      <!-- 載入狀態 -->
      <div v-if="loading" class="loading-container">
        <div class="loading-spinner"></div>
        <p>載入中...</p>
      </div>

      <!-- 股票表格 -->
      <div v-else class="stock-table-container">
        <table class="stock-table">
          <thead>
            <tr>
              <th>股票代碼</th>
              <th>股票名稱</th>
              <th>產業分類</th>
              <th>市場</th>
              <th>現價</th>
              <th>漲跌</th>
              <th>漲跌幅</th>
              <th>成交量</th>
              <th>成交金額</th>
              <th>操作</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="stock in stocks" :key="stock.code" class="stock-row">
              <td class="stock-code">{{ stock.code }}</td>
              <td class="stock-name">{{ stock.name }}</td>
              <td class="stock-category">{{ getCategoryName(stock.category) }}</td>
              <td class="stock-market">
                <span :class="['market-badge', stock.market.toLowerCase()]">
                  {{ stock.market === 'TSE' ? '上市' : '上櫃' }}
                </span>
              </td>
              <td class="stock-price">
                <span v-if="stock.price" class="price-value">
                  {{ formatPrice(stock.price.price) }}
                </span>
                <span v-else class="no-data">--</span>
              </td>
              <td class="stock-change">
                <span v-if="stock.price" :class="['change-value', getChangeClass(stock.price.change)]">
                  {{ formatChange(stock.price.change) }}
                </span>
                <span v-else class="no-data">--</span>
              </td>
              <td class="stock-change-percent">
                <span v-if="stock.price" :class="['change-percent', getChangeClass(stock.price.change)]">
                  {{ formatPercent(stock.price.change_percent) }}
                </span>
                <span v-else class="no-data">--</span>
              </td>
              <td class="stock-volume">
                <span v-if="stock.price" class="volume-value">
                  {{ formatVolume(stock.price.volume) }}
                </span>
                <span v-else class="no-data">--</span>
              </td>
              <td class="stock-amount">
                <span v-if="stock.price" class="amount-value">
                  {{ formatAmount(stock.price.amount) }}
                </span>
                <span v-else class="no-data">--</span>
              </td>
              <td class="stock-actions">
                <button @click="viewStockDetail(stock.code)" class="view-btn">
                  查看詳情
                </button>
              </td>
            </tr>
          </tbody>
        </table>

        <!-- 無數據提示 -->
        <div v-if="stocks.length === 0" class="no-data-container">
          <p>沒有找到符合條件的股票</p>
        </div>
      </div>
    </div>

    <!-- 分頁組件 -->
    <div v-if="!loading && pagination.total_pages > 1" class="pagination-container">
      <StockPagination 
        :current-page="pagination.current_page"
        :total-pages="pagination.total_pages"
        :total-count="pagination.total_count"
        :per-page="pagination.per_page"
        @page-change="onPageChange"
        @per-page-change="onPerPageChange"
      />
    </div>

    <!-- 統計資訊 -->
    <div v-if="!loading" class="stats-container">
      <div class="stats-item">
        <span class="stats-label">總共</span>
        <span class="stats-value">{{ pagination.total_count }}</span>
        <span class="stats-label">支股票</span>
      </div>
      <div class="stats-item">
        <span class="stats-label">第</span>
        <span class="stats-value">{{ pagination.current_page }}</span>
        <span class="stats-label">頁，共</span>
        <span class="stats-value">{{ pagination.total_pages }}</span>
        <span class="stats-label">頁</span>
      </div>
    </div>
  </div>
</template>

<script>
import { ref, onMounted, watch } from 'vue'
import { useRouter } from 'vue-router'
import StockPagination from '@/components/stock/StockPagination.vue'
import api from '@/services/api'

export default {
  name: 'StockListPage',
  components: {
    StockPagination
  },
  setup() {
    const router = useRouter()
    
    // 響應式數據
    const stocks = ref([])
    const categories = ref([])
    const loading = ref(false)
    const searchKeyword = ref('')
    const selectedCategory = ref('')
    const selectedMarket = ref('')
    const sortBy = ref('code')
    const sortOrder = ref('asc')
    const pagination = ref({
      current_page: 1,
      per_page: 20,
      total_pages: 1,
      total_count: 0,
      has_next: false,
      has_prev: false
    })

    // 載入股票分類
    const loadCategories = async () => {
      try {
        const response = await api.get('/api/stock/categories')
        categories.value = response.data.data || []
      } catch (error) {
        console.error('載入股票分類失敗:', error)
      }
    }

    // 載入股票列表
    const loadStocks = async () => {
      loading.value = true
      try {
        const params = {
          page: pagination.value.current_page,
          limit: pagination.value.per_page,
          category: selectedCategory.value,
          market: selectedMarket.value,
          search: searchKeyword.value,
          sort_by: sortBy.value,
          sort_order: sortOrder.value
        }

        const response = await api.get('/api/stock/stocks', { params })
        const data = response.data.data
        
        stocks.value = data.stocks || []
        pagination.value = data.pagination || pagination.value
      } catch (error) {
        console.error('載入股票列表失敗:', error)
        stocks.value = []
      } finally {
        loading.value = false
      }
    }

    // 搜尋處理
    const onSearchInput = () => {
      // 防抖處理
      clearTimeout(searchTimeout)
      searchTimeout = setTimeout(() => {
        pagination.value.current_page = 1
        loadStocks()
      }, 500)
    }

    let searchTimeout = null

    const performSearch = () => {
      pagination.value.current_page = 1
      loadStocks()
    }

    // 篩選變更
    const onFilterChange = () => {
      pagination.value.current_page = 1
      loadStocks()
    }

    // 排序變更
    const onSortChange = () => {
      pagination.value.current_page = 1
      loadStocks()
    }

    // 分頁變更
    const onPageChange = (page) => {
      pagination.value.current_page = page
      loadStocks()
    }

    // 每頁筆數變更
    const onPerPageChange = (perPage) => {
      pagination.value.per_page = perPage
      pagination.value.current_page = 1
      loadStocks()
    }

    // 查看股票詳情
    const viewStockDetail = (code) => {
      router.push(`/stock/${code}`)
    }

    // 格式化函數
    const formatPrice = (price) => {
      if (price === null || price === undefined) return '--'
      return price.toFixed(2)
    }

    const formatChange = (change) => {
      if (change === null || change === undefined) return '--'
      return change > 0 ? `+${change.toFixed(2)}` : change.toFixed(2)
    }

    const formatPercent = (percent) => {
      if (percent === null || percent === undefined) return '--'
      return percent > 0 ? `+${percent.toFixed(2)}%` : `${percent.toFixed(2)}%`
    }

    const formatVolume = (volume) => {
      if (volume === null || volume === undefined) return '--'
      if (volume >= 100000000) {
        return `${(volume / 100000000).toFixed(1)}億`
      } else if (volume >= 10000) {
        return `${(volume / 10000).toFixed(1)}萬`
      }
      return volume.toLocaleString()
    }

    const formatAmount = (amount) => {
      if (amount === null || amount === undefined) return '--'
      if (amount >= 100000000) {
        return `${(amount / 100000000).toFixed(1)}億`
      } else if (amount >= 10000) {
        return `${(amount / 10000).toFixed(1)}萬`
      }
      return amount.toLocaleString()
    }

    const getChangeClass = (change) => {
      if (change > 0) return 'positive'
      if (change < 0) return 'negative'
      return 'neutral'
    }

    const getCategoryName = (categoryCode) => {
      const category = categories.value.find(cat => cat.code === categoryCode)
      return category ? category.name : categoryCode
    }

    // 監聽路由參數變化
    watch(() => router.currentRoute.value.query, (newQuery) => {
      if (newQuery.category) selectedCategory.value = newQuery.category
      if (newQuery.market) selectedMarket.value = newQuery.market
      if (newQuery.search) searchKeyword.value = newQuery.search
      loadStocks()
    }, { immediate: true })

    // 組件掛載時載入數據
    onMounted(() => {
      loadCategories()
      loadStocks()
    })

    return {
      stocks,
      categories,
      loading,
      searchKeyword,
      selectedCategory,
      selectedMarket,
      sortBy,
      sortOrder,
      pagination,
      onSearchInput,
      performSearch,
      onFilterChange,
      onSortChange,
      onPageChange,
      onPerPageChange,
      viewStockDetail,
      formatPrice,
      formatChange,
      formatPercent,
      formatVolume,
      formatAmount,
      getChangeClass,
      getCategoryName
    }
  }
}
</script>

<style scoped>
.stock-list-page {
  max-width: 1400px;
  margin: 0 auto;
  padding: 20px;
}

.page-header {
  text-align: center;
  margin-bottom: 30px;
}

.page-header h1 {
  color: #2d3748;
  font-size: 2.5rem;
  margin-bottom: 10px;
}

.page-subtitle {
  color: #718096;
  font-size: 1.1rem;
}

.filter-section {
  background: white;
  border-radius: 12px;
  padding: 20px;
  margin-bottom: 20px;
  box-shadow: 0 2px 10px rgba(0,0,0,0.1);
}

.filter-row {
  display: flex;
  gap: 20px;
  align-items: center;
  flex-wrap: wrap;
}

.search-box {
  display: flex;
  gap: 10px;
  flex: 1;
  min-width: 300px;
}

.search-input {
  flex: 1;
  padding: 10px 15px;
  border: 2px solid #e2e8f0;
  border-radius: 8px;
  font-size: 14px;
  transition: border-color 0.3s ease;
}

.search-input:focus {
  outline: none;
  border-color: #4299e1;
}

.search-btn {
  padding: 10px 20px;
  background: #4299e1;
  color: white;
  border: none;
  border-radius: 8px;
  cursor: pointer;
  font-size: 14px;
  transition: background-color 0.3s ease;
}

.search-btn:hover {
  background: #3182ce;
}

.filter-group {
  display: flex;
  align-items: center;
  gap: 8px;
}

.filter-group label {
  font-weight: 500;
  color: #4a5568;
  white-space: nowrap;
}

.filter-select {
  padding: 8px 12px;
  border: 2px solid #e2e8f0;
  border-radius: 6px;
  font-size: 14px;
  background: white;
  cursor: pointer;
}

.filter-select:focus {
  outline: none;
  border-color: #4299e1;
}

.stock-list-container {
  background: white;
  border-radius: 12px;
  overflow: hidden;
  box-shadow: 0 2px 10px rgba(0,0,0,0.1);
}

.loading-container {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 60px 20px;
}

.loading-spinner {
  width: 40px;
  height: 40px;
  border: 4px solid #e2e8f0;
  border-top: 4px solid #4299e1;
  border-radius: 50%;
  animation: spin 1s linear infinite;
  margin-bottom: 20px;
}

@keyframes spin {
  0% { transform: rotate(0deg); }
  100% { transform: rotate(360deg); }
}

.stock-table-container {
  overflow-x: auto;
}

.stock-table {
  width: 100%;
  border-collapse: collapse;
  font-size: 14px;
}

.stock-table th {
  background: #f7fafc;
  padding: 15px 12px;
  text-align: left;
  font-weight: 600;
  color: #4a5568;
  border-bottom: 2px solid #e2e8f0;
  white-space: nowrap;
}

.stock-table td {
  padding: 12px;
  border-bottom: 1px solid #e2e8f0;
  vertical-align: middle;
}

.stock-row:hover {
  background: #f7fafc;
}

.stock-code {
  font-weight: 600;
  color: #2d3748;
  font-family: 'Courier New', monospace;
}

.stock-name {
  font-weight: 500;
  color: #2d3748;
}

.stock-category {
  color: #718096;
  font-size: 13px;
}

.market-badge {
  padding: 4px 8px;
  border-radius: 12px;
  font-size: 12px;
  font-weight: 500;
}

.market-badge.tse {
  background: #e6fffa;
  color: #234e52;
}

.market-badge.otc {
  background: #fef5e7;
  color: #744210;
}

.price-value {
  font-weight: 600;
  font-size: 15px;
}

.change-value {
  font-weight: 600;
}

.change-value.positive {
  color: #e53e3e;
}

.change-value.negative {
  color: #38a169;
}

.change-value.neutral {
  color: #718096;
}

.change-percent {
  font-weight: 600;
}

.change-percent.positive {
  color: #e53e3e;
}

.change-percent.negative {
  color: #38a169;
}

.change-percent.neutral {
  color: #718096;
}

.volume-value, .amount-value {
  font-family: 'Courier New', monospace;
  font-size: 13px;
}

.no-data {
  color: #a0aec0;
  font-style: italic;
}

.view-btn {
  padding: 6px 12px;
  background: #4299e1;
  color: white;
  border: none;
  border-radius: 6px;
  cursor: pointer;
  font-size: 12px;
  transition: background-color 0.3s ease;
}

.view-btn:hover {
  background: #3182ce;
}

.no-data-container {
  text-align: center;
  padding: 60px 20px;
  color: #718096;
}

.pagination-container {
  margin-top: 30px;
  display: flex;
  justify-content: center;
}

.stats-container {
  display: flex;
  justify-content: center;
  gap: 30px;
  margin-top: 20px;
  padding: 15px;
  background: #f7fafc;
  border-radius: 8px;
}

.stats-item {
  display: flex;
  align-items: center;
  gap: 5px;
}

.stats-label {
  color: #718096;
  font-size: 14px;
}

.stats-value {
  color: #2d3748;
  font-weight: 600;
  font-size: 16px;
}

/* 響應式設計 */
@media (max-width: 768px) {
  .filter-row {
    flex-direction: column;
    align-items: stretch;
  }
  
  .search-box {
    min-width: auto;
  }
  
  .filter-group {
    justify-content: space-between;
  }
  
  .stock-table {
    font-size: 12px;
  }
  
  .stock-table th,
  .stock-table td {
    padding: 8px 6px;
  }
  
  .stats-container {
    flex-direction: column;
    gap: 10px;
    text-align: center;
  }
}
</style>
