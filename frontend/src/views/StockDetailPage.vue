<template>
  <div class="stock-detail-page">
    <!-- 回上一頁按鈕 -->
    <div class="back-button-container">
      <button @click="goBack" class="back-button">
        <span class="back-icon">←</span>
        回上一頁
      </button>
    </div>

    <!-- 載入狀態 -->
    <div v-if="loading" class="loading-container">
      <div class="loading-spinner"></div>
      <p>載入股票資訊中...</p>
    </div>

    <!-- 股票不存在 -->
    <div v-else-if="!loading && !stock" class="error-container">
      <div class="error-message">
        <h2>股票不存在</h2>
        <p>找不到您要查看的股票資訊</p>
        <button @click="goBack" class="back-button">返回股票列表</button>
      </div>
    </div>

    <!-- 股票詳情 -->
    <div v-else-if="stock" class="stock-detail">
      <!-- 股票基本資訊 -->
      <div class="stock-header">
        <div class="stock-info">
          <h1 class="stock-name">{{ stock.name }}</h1>
          <div class="stock-code">{{ stock.code }}</div>
          <div class="stock-category">{{ getCategoryName(stock.category) }}</div>
          <div class="stock-market">
            <span :class="['market-badge', stock.market.toLowerCase()]">
              {{ stock.market === 'TSE' ? '上市' : '上櫃' }}
            </span>
          </div>
        </div>
        
        <div v-if="stock.price" class="price-info">
          <div class="current-price">
            <span class="price-value">{{ formatPrice(stock.price.price) }}</span>
            <span class="price-unit">元</span>
          </div>
          <div class="price-change">
            <span :class="['change-value', getChangeClass(stock.price.change)]">
              {{ formatChange(stock.price.change) }}
            </span>
            <span :class="['change-percent', getChangeClass(stock.price.change)]">
              {{ formatPercent(stock.price.change_percent) }}
            </span>
          </div>
        </div>
        <div v-else class="no-price-info">
          <div class="no-price-message">暫無價格數據</div>
        </div>
      </div>

      <!-- 最後更新時間 -->
      <div v-if="lastUpdateTime" class="last-update">
        <small>最後更新：{{ lastUpdateTime.toLocaleString() }}</small>
      </div>

      <!-- 價格資訊卡片 -->
      <div v-if="stock.price" class="price-cards">
        <div class="price-card">
          <div class="card-title">開盤價</div>
          <div class="card-value">{{ formatPrice(stock.price.open_price) }}</div>
        </div>
        <div class="price-card">
          <div class="card-title">最高價</div>
          <div class="card-value high">{{ formatPrice(stock.price.high_price) }}</div>
        </div>
        <div class="price-card">
          <div class="card-title">最低價</div>
          <div class="card-value low">{{ formatPrice(stock.price.low_price) }}</div>
        </div>
        <div class="price-card">
          <div class="card-title">昨收價</div>
          <div class="card-value">{{ formatPrice(stock.price.close_price) }}</div>
        </div>
        <div class="price-card">
          <div class="card-title">成交量</div>
          <div class="card-value">{{ formatVolume(stock.price.volume) }}</div>
        </div>
        <div class="price-card">
          <div class="card-title">成交金額</div>
          <div class="card-value">{{ formatAmount(stock.price.amount) }}</div>
        </div>
      </div>

      <!-- 圖表區域 -->
      <div class="chart-section">
        <div class="section-header">
          <h2>價格走勢圖</h2>
          <div class="chart-controls">
            <button 
              v-for="period in chartPeriods" 
              :key="period.value"
              @click="selectedPeriod = period.value"
              :class="['period-btn', { active: selectedPeriod === period.value }]"
            >
              {{ period.label }}
            </button>
          </div>
        </div>
        <div class="chart-container">
          <div class="chart-placeholder">
            <p>圖表功能開發中...</p>
            <p>當前選擇：{{ getPeriodLabel(selectedPeriod) }}</p>
          </div>
        </div>
      </div>

      <!-- 技術指標區域 -->
      <div class="indicators-section">
        <h2>技術指標</h2>
        <div class="indicators-grid">
          <div class="indicator-card">
            <div class="indicator-title">RSI (14)</div>
            <div class="indicator-value">--</div>
            <div class="indicator-status">計算中</div>
          </div>
          <div class="indicator-card">
            <div class="indicator-title">MACD</div>
            <div class="indicator-value">--</div>
            <div class="indicator-status">計算中</div>
          </div>
          <div class="indicator-card">
            <div class="indicator-title">KD</div>
            <div class="indicator-value">--</div>
            <div class="indicator-status">計算中</div>
          </div>
          <div class="indicator-card">
            <div class="indicator-title">MA5</div>
            <div class="indicator-value">--</div>
            <div class="indicator-status">計算中</div>
          </div>
        </div>
      </div>

      <!-- 新聞資訊區域 -->
      <div class="news-section">
        <h2>相關新聞</h2>
        <div class="news-list">
          <div class="news-item">
            <div class="news-title">新聞功能開發中...</div>
            <div class="news-time">2024-01-01 12:00</div>
          </div>
        </div>
      </div>
    </div>

    <!-- 錯誤狀態 -->
    <div v-else class="error-container">
      <div class="error-icon">⚠️</div>
      <h2>股票不存在</h2>
      <p>找不到指定的股票資訊</p>
      <button @click="goBack" class="back-btn">返回列表</button>
    </div>

    <!-- AI 股票助手按鈕 -->
    <button class="ai-stock-chat-btn" @click="toggleStockChatWindow" title="AI 股票助手">
      🤖
      <span class="ai-btn-text">AI助手</span>
    </button>

    <!-- AI 股票聊天窗口 -->
    <AIChatWindow 
      v-if="showStockChatWindow"
      @close="toggleStockChatWindow"
      :stock-context="stock && stock.code ? stock : null"
    />
  </div>
</template>

<script>
import { ref, onMounted, computed, watch, onUnmounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import api from '@/services/api'
import AIChatWindow from '@/components/chat/AIChatWindow.vue'

export default {
  name: 'StockDetailPage',
  components: {
    AIChatWindow
  },
  setup() {
    const route = useRoute()
    const router = useRouter()
    
    // 響應式數據
    const stock = ref(null)
    const loading = ref(true)
    const selectedPeriod = ref('1d')
    const categories = ref([])
    const autoUpdateInterval = ref(null)
    const lastUpdateTime = ref(null)
    const showStockChatWindow = ref(false)

    // 圖表週期選項
    const chartPeriods = ref([
      { label: '1日', value: '1d' },
      { label: '1週', value: '1w' },
      { label: '1月', value: '1m' },
      { label: '3月', value: '3m' },
      { label: '6月', value: '6m' },
      { label: '1年', value: '1y' }
    ])

    // 載入股票分類
    const loadCategories = async () => {
      try {
        const response = await api.get('/api/stock/categories')
        categories.value = response.data.data || []
      } catch (error) {
        console.error('載入股票分類失敗:', error)
      }
    }

    // 載入股票詳情
    const loadStockDetail = async (showLoading = true) => {
      const code = route.params.code
      if (!code) {
        loading.value = false
        return
      }

      if (showLoading) {
        loading.value = true
      }

      try {
        const response = await api.get(`/api/stock/stocks/${code}`)
        stock.value = response.data.data
        lastUpdateTime.value = new Date()
      } catch (error) {
        console.error('載入股票詳情失敗:', error)
        stock.value = null
      } finally {
        loading.value = false
      }
    }

    // 開始自動更新
    const startAutoUpdate = () => {
      // 每5秒更新一次
      autoUpdateInterval.value = setInterval(() => {
        loadStockDetail(false) // 不顯示loading
      }, 5000)
    }

    // 停止自動更新
    const stopAutoUpdate = () => {
      if (autoUpdateInterval.value) {
        clearInterval(autoUpdateInterval.value)
        autoUpdateInterval.value = null
      }
    }

    // 返回列表
    const goBack = () => {
      router.push('/stocks')
    }

    // 切換AI聊天窗口
    const toggleStockChatWindow = () => {
      showStockChatWindow.value = !showStockChatWindow.value
    }

    // 格式化函數
    const formatPrice = (price) => {
      if (price === null || price === undefined || price === 0) return '--'
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
      if (volume === null || volume === undefined || volume === 0) return '--'
      if (volume >= 100000000) {
        return `${(volume / 100000000).toFixed(1)}億`
      } else if (volume >= 10000) {
        return `${(volume / 10000).toFixed(1)}萬`
      }
      return volume.toLocaleString()
    }

    const formatAmount = (amount) => {
      if (amount === null || amount === undefined || amount === 0) return '--'
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

    const getPeriodLabel = (value) => {
      const period = chartPeriods.value.find(p => p.value === value)
      return period ? period.label : value
    }

    // 監聽路由參數變化
    watch(() => route.params.code, (newCode) => {
      if (newCode) {
        loadStockDetail()
      }
    })

    // 組件掛載時載入數據
    onMounted(() => {
      loadCategories()
      loadStockDetail()
      startAutoUpdate()
    })

    // 組件卸載時停止自動更新
    onUnmounted(() => {
      stopAutoUpdate()
    })

    return {
      stock,
      loading,
      selectedPeriod,
      chartPeriods,
      lastUpdateTime,
      showStockChatWindow,
      loadStockDetail,
      goBack,
      toggleStockChatWindow,
      formatPrice,
      formatChange,
      formatPercent,
      formatVolume,
      formatAmount,
      getChangeClass,
      getCategoryName,
      getPeriodLabel
    }
  }
}
</script>

<style scoped>
.stock-detail-page {
  max-width: 1200px;
  margin: 0 auto;
  padding: 20px;
}

.loading-container {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 100px 20px;
}

.loading-spinner {
  width: 50px;
  height: 50px;
  border: 5px solid #e2e8f0;
  border-top: 5px solid #4299e1;
  border-radius: 50%;
  animation: spin 1s linear infinite;
  margin-bottom: 20px;
}

@keyframes spin {
  0% { transform: rotate(0deg); }
  100% { transform: rotate(360deg); }
}

.stock-detail {
  display: flex;
  flex-direction: column;
  gap: 30px;
}

.stock-header {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  background: white;
  padding: 30px;
  border-radius: 12px;
  box-shadow: 0 4px 20px rgba(0,0,0,0.1);
}

.stock-info {
  flex: 1;
}

.stock-name {
  font-size: 2.5rem;
  font-weight: 700;
  color: #2d3748;
  margin-bottom: 10px;
}

.stock-code {
  font-size: 1.2rem;
  font-weight: 600;
  color: #4a5568;
  font-family: 'Courier New', monospace;
  margin-bottom: 8px;
}

.stock-category {
  font-size: 1rem;
  color: #718096;
  margin-bottom: 8px;
}

.market-badge {
  padding: 6px 12px;
  border-radius: 16px;
  font-size: 14px;
  font-weight: 600;
}

.market-badge.tse {
  background: #e6fffa;
  color: #234e52;
}

.market-badge.otc {
  background: #fef5e7;
  color: #744210;
}

.price-info {
  text-align: right;
}

.current-price {
  display: flex;
  align-items: baseline;
  gap: 8px;
  margin-bottom: 10px;
}

.price-value {
  font-size: 3rem;
  font-weight: 700;
  color: #2d3748;
}

.price-unit {
  font-size: 1.2rem;
  color: #718096;
}

.price-change {
  display: flex;
  align-items: center;
  gap: 10px;
  justify-content: flex-end;
}

.change-value {
  font-size: 1.2rem;
  font-weight: 600;
}

.change-percent {
  font-size: 1rem;
  font-weight: 600;
}

.change-value.positive,
.change-percent.positive {
  color: #e53e3e;
}

.change-value.negative,
.change-percent.negative {
  color: #38a169;
}

.change-value.neutral,
.change-percent.neutral {
  color: #718096;
}

.price-cards {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
  gap: 20px;
}

.price-card {
  background: white;
  padding: 20px;
  border-radius: 12px;
  box-shadow: 0 2px 10px rgba(0,0,0,0.1);
  text-align: center;
}

.card-title {
  font-size: 14px;
  color: #718096;
  margin-bottom: 8px;
  font-weight: 500;
}

.card-value {
  font-size: 1.5rem;
  font-weight: 700;
  color: #2d3748;
}

.card-value.high {
  color: #e53e3e;
}

.card-value.low {
  color: #38a169;
}

.chart-section,
.indicators-section,
.news-section {
  background: white;
  padding: 30px;
  border-radius: 12px;
  box-shadow: 0 2px 10px rgba(0,0,0,0.1);
}

.section-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
}

.section-header h2 {
  font-size: 1.5rem;
  font-weight: 600;
  color: #2d3748;
  margin: 0;
}

.chart-controls {
  display: flex;
  gap: 8px;
}

.period-btn {
  padding: 8px 16px;
  border: 2px solid #e2e8f0;
  background: white;
  color: #4a5568;
  border-radius: 8px;
  cursor: pointer;
  font-size: 14px;
  font-weight: 500;
  transition: all 0.3s ease;
}

.period-btn:hover {
  border-color: #4299e1;
  color: #4299e1;
}

.period-btn.active {
  background: #4299e1;
  border-color: #4299e1;
  color: white;
}

.chart-container {
  height: 400px;
  border: 2px dashed #e2e8f0;
  border-radius: 8px;
  display: flex;
  align-items: center;
  justify-content: center;
}

.chart-placeholder {
  text-align: center;
  color: #718096;
}

.indicators-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
  gap: 20px;
}

.indicator-card {
  background: #f7fafc;
  padding: 20px;
  border-radius: 8px;
  text-align: center;
}

.indicator-title {
  font-size: 14px;
  color: #4a5568;
  margin-bottom: 8px;
  font-weight: 500;
}

.indicator-value {
  font-size: 1.2rem;
  font-weight: 600;
  color: #2d3748;
  margin-bottom: 4px;
}

.indicator-status {
  font-size: 12px;
  color: #718096;
}

.news-list {
  display: flex;
  flex-direction: column;
  gap: 15px;
}

.news-item {
  padding: 15px;
  background: #f7fafc;
  border-radius: 8px;
  border-left: 4px solid #4299e1;
}

.news-title {
  font-size: 16px;
  font-weight: 500;
  color: #2d3748;
  margin-bottom: 5px;
}

.news-time {
  font-size: 14px;
  color: #718096;
}

.error-container {
  text-align: center;
  padding: 100px 20px;
}

.error-icon {
  font-size: 4rem;
  margin-bottom: 20px;
}

.error-container h2 {
  font-size: 2rem;
  color: #2d3748;
  margin-bottom: 10px;
}

.error-container p {
  color: #718096;
  margin-bottom: 30px;
}

.back-btn {
  padding: 12px 24px;
  background: #4299e1;
  color: white;
  border: none;
  border-radius: 8px;
  cursor: pointer;
  font-size: 16px;
  font-weight: 500;
  transition: background-color 0.3s ease;
}

.back-btn:hover {
  background: #3182ce;
}

/* 響應式設計 */
@media (max-width: 768px) {
  .stock-header {
    flex-direction: column;
    gap: 20px;
  }

  .price-info {
    text-align: left;
  }

  .price-change {
    justify-content: flex-start;
  }

  .section-header {
    flex-direction: column;
    gap: 15px;
    align-items: flex-start;
  }

  .chart-controls {
    flex-wrap: wrap;
  }

  .price-cards {
    grid-template-columns: repeat(auto-fit, minmax(150px, 1fr));
  }

  .indicators-grid {
    grid-template-columns: repeat(auto-fit, minmax(150px, 1fr));
  }
}

/* 回上一頁按鈕 */
.back-button-container {
  margin-bottom: 1rem;
}

.back-button {
  display: inline-flex;
  align-items: center;
  gap: 0.5rem;
  padding: 0.75rem 1.5rem;
  background: #f8f9fa;
  border: 1px solid #dee2e6;
  border-radius: 8px;
  color: #495057;
  text-decoration: none;
  font-size: 0.9rem;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.2s ease;
}

.back-button:hover {
  background: #e9ecef;
  border-color: #adb5bd;
  color: #212529;
}

.back-icon {
  font-size: 1.1rem;
  font-weight: bold;
}

/* AI 股票助手按鈕 */
.ai-stock-chat-btn {
  position: fixed;
  bottom: 30px;
  right: 30px;
  width: 60px;
  height: 60px;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  border: none;
  border-radius: 50%;
  color: white;
  font-size: 24px;
  cursor: pointer;
  box-shadow: 0 4px 20px rgba(102, 126, 234, 0.4);
  transition: all 0.3s ease;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  z-index: 1000;
}

.ai-stock-chat-btn:hover {
  transform: translateY(-3px) scale(1.05);
  box-shadow: 0 6px 25px rgba(102, 126, 234, 0.6);
}

.ai-btn-text {
  font-size: 10px;
  font-weight: 600;
  margin-top: 2px;
  line-height: 1;
}

/* 響應式設計 */
@media (max-width: 768px) {
  .ai-stock-chat-btn {
    bottom: 20px;
    right: 20px;
    width: 50px;
    height: 50px;
    font-size: 20px;
  }
  
  .ai-btn-text {
    font-size: 8px;
  }
}
</style>
