<template>
  <div class="stock-market-page">
    <!-- 頁面標題 -->
    <div class="page-header">
      <div class="header-content">
        <div class="header-left">
          <button @click="goHome" class="home-button">
            <span class="home-icon">🏠</span>
            回首頁
          </button>
        </div>
        <div class="header-center">
          <h1>阿和台股站</h1>
          <p class="subtitle">即時台股資訊，投資理財好幫手</p>
        </div>
        <div class="header-right">
          <!-- 預留空間，保持對稱 -->
        </div>
      </div>
    </div>

    <!-- 市場概覽 -->
    <div class="market-overview">
      <div class="overview-cards">
        <div class="overview-card">
          <div class="card-icon">📈</div>
          <div class="card-content">
            <h3>加權指數</h3>
            <div class="price">{{ formatIndex(marketData.taiex) }}</div>
            <div class="change" :class="marketData.taiexChange >= 0 ? 'positive' : 'negative'">
              {{ marketData.taiexChange >= 0 ? '+' : '' }}{{ formatChange(marketData.taiexChange) }}
              ({{ marketData.taiexChangePercent >= 0 ? '+' : '' }}{{ formatChangePercent(marketData.taiexChangePercent) }}%)
            </div>
          </div>
        </div>
        
        <div class="overview-card">
          <div class="card-icon">🏭</div>
          <div class="card-content">
            <h3>上市股票</h3>
            <div class="price">{{ marketData.listedCount || '--' }}</div>
            <div class="change">支股票</div>
          </div>
        </div>
        
        <div class="overview-card">
          <div class="card-icon">📊</div>
          <div class="card-content">
            <h3>上櫃指數</h3>
            <div class="price">{{ formatIndex(marketData.otcIndex) }}</div>
            <div class="change" :class="marketData.otcChange >= 0 ? 'positive' : 'negative'">
              {{ marketData.otcChange >= 0 ? '+' : '' }}{{ formatChange(marketData.otcChange) }}
              ({{ marketData.otcChangePercent >= 0 ? '+' : '' }}{{ formatChangePercent(marketData.otcChangePercent) }}%)
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- 快速導航 -->
    <div class="quick-nav">
      <h2>快速導航</h2>
      <div class="nav-buttons">
        <router-link to="/stocks" class="nav-button">
          <div class="nav-icon">📊</div>
          <div class="nav-text">股票列表</div>
        </router-link>
        <router-link to="/stocks?category=ELECTRONICS" class="nav-button">
          <div class="nav-icon">💻</div>
          <div class="nav-text">電子股</div>
        </router-link>
        <router-link to="/stocks?category=FINANCE" class="nav-button">
          <div class="nav-icon">🏦</div>
          <div class="nav-text">金融股</div>
        </router-link>
        <router-link to="/stocks?category=INDUSTRY" class="nav-button">
          <div class="nav-icon">🏭</div>
          <div class="nav-text">傳產股</div>
        </router-link>
      </div>
    </div>

    <!-- 熱門股票 -->
    <div class="hot-stocks">
      <h2>熱門股票</h2>
      <div class="stocks-grid" v-if="hotStocks.length > 0">
        <div 
          v-for="stock in hotStocks" 
          :key="stock.code"
          class="stock-card"
          @click="goToStockDetail(stock.code)"
        >
          <div class="stock-header">
            <div class="stock-code">{{ stock.code }}</div>
            <div class="stock-name">{{ stock.name }}</div>
          </div>
          <div class="stock-price">
            <div class="current-price">{{ formatPrice(stock.price?.price) }}</div>
            <div 
              class="price-change"
              :class="stock.price?.change >= 0 ? 'positive' : 'negative'"
            >
              {{ stock.price?.change >= 0 ? '+' : '' }}{{ formatPrice(stock.price?.change) }}
              ({{ stock.price?.changePercent >= 0 ? '+' : '' }}{{ formatPrice(stock.price?.changePercent) }}%)
            </div>
          </div>
        </div>
      </div>
      <div v-else class="loading">
        <div class="spinner"></div>
        <p>載入中...</p>
      </div>
    </div>

    <!-- 市場資訊 -->
    <div class="market-info">
      <h2>市場資訊</h2>
      <div class="info-cards">
        <div class="info-card">
          <h3>交易時間</h3>
          <p>週一至週五 09:00-13:30</p>
          <p class="note">連續交易，無午休</p>
        </div>
        <div class="info-card">
          <h3>漲跌幅限制</h3>
          <p>一般股票：±10%</p>
          <p class="note">特殊股票：±5%</p>
        </div>
        <div class="info-card">
          <h3>最小交易單位</h3>
          <p>1張 = 1000股</p>
          <p class="note">零股交易：1股起</p>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'

export default {
  name: 'StockMarketPage',
  setup() {
    const router = useRouter()
    const hotStocks = ref([])
    const marketData = ref({
      taiex: null,
      taiexChange: null,
      taiexChangePercent: null,
      otcIndex: null,
      otcChange: null,
      otcChangePercent: null,
      listedCount: null,
      totalAmount: null
    })

    // 獲取熱門股票
    const fetchHotStocks = async () => {
      try {
        const response = await fetch('/api/stock/stocks?page=1&limit=8&sort_by=volume&sort_order=desc')
        const data = await response.json()
        if (data.success) {
          hotStocks.value = data.data.stocks
        }
      } catch (error) {
        console.error('獲取熱門股票失敗:', error)
      }
    }

    // 獲取市場數據
    const fetchMarketData = async () => {
      try {
        const response = await fetch('/api/stock/market-stats')
        const data = await response.json()
        if (data.success) {
          marketData.value = {
            taiex: data.data.taiex,
            taiexChange: data.data.taiexChange,
            taiexChangePercent: data.data.taiexChangePercent,
            otcIndex: data.data.otc_index,
            otcChange: data.data.otc_change,
            otcChangePercent: data.data.otc_change_percent,
            listedCount: data.data.total_count,
            totalAmount: data.data.total_amount
          }
        }
      } catch (error) {
        console.error('獲取市場數據失敗:', error)
        // 如果API失敗，使用模擬數據
        marketData.value = {
          taiex: 17500.25,
          taiexChange: 125.50,
          taiexChangePercent: 0.72,
          otcIndex: 220.0,
          otcChange: 2.5,
          otcChangePercent: 1.15,
          listedCount: 44,
          totalAmount: 2850.75
        }
      }
    }

    // 格式化價格
    const formatPrice = (price) => {
      if (price === null || price === undefined) return '--'
      return Number(price).toFixed(2)
    }

    // 格式化指數（小數點第二位）
    const formatIndex = (index) => {
      if (index === null || index === undefined) return '--'
      return Number(index).toFixed(2)
    }

    // 格式化漲跌點數（整數）
    const formatChange = (change) => {
      if (change === null || change === undefined) return '--'
      return Math.round(Number(change))
    }

    // 格式化漲跌幅（小數點第二位）
    const formatChangePercent = (percent) => {
      if (percent === null || percent === undefined) return '--'
      return Number(percent).toFixed(2)
    }

    // 格式化金額
    const formatAmount = (amount) => {
      if (amount === null || amount === undefined) return '--'
      return Number(amount).toFixed(2)
    }

    // 跳轉到股票詳情
    const goToStockDetail = (code) => {
      router.push(`/stocks?code=${code}`)
    }

    // 回首頁
    const goHome = () => {
      router.push('/')
    }

    onMounted(() => {
      fetchHotStocks()
      fetchMarketData()
    })

    return {
      hotStocks,
      marketData,
      formatPrice,
      formatIndex,
      formatChange,
      formatChangePercent,
      formatAmount,
      goToStockDetail,
      goHome
    }
  }
}
</script>

<style scoped>
.stock-market-page {
  max-width: 1200px;
  margin: 0 auto;
  padding: 20px;
}

.page-header {
  margin-bottom: 40px;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
  border-radius: 12px;
  padding: 2rem 0;
}

.header-content {
  display: flex;
  align-items: center;
  justify-content: space-between;
  max-width: 1200px;
  margin: 0 auto;
  padding: 0 2rem;
}

.header-left, .header-right {
  flex: 1;
}

.header-center {
  flex: 2;
  text-align: center;
}

.page-header h1 {
  font-size: 2.5rem;
  color: white;
  margin-bottom: 10px;
  text-shadow: 2px 2px 4px rgba(0,0,0,0.3);
}

.subtitle {
  font-size: 1.2rem;
  color: rgba(255, 255, 255, 0.9);
  margin: 0;
}

/* 回首頁按鈕 */
.home-button {
  display: inline-flex;
  align-items: center;
  gap: 0.5rem;
  padding: 0.75rem 1.5rem;
  background: rgba(255, 255, 255, 0.2);
  border: 2px solid rgba(255, 255, 255, 0.3);
  border-radius: 8px;
  color: white;
  text-decoration: none;
  font-size: 0.9rem;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.3s ease;
  backdrop-filter: blur(10px);
}

.home-button:hover {
  background: rgba(255, 255, 255, 0.3);
  border-color: rgba(255, 255, 255, 0.5);
  transform: translateY(-2px);
  box-shadow: 0 4px 12px rgba(0,0,0,0.2);
}

.home-icon {
  font-size: 1.1rem;
}

.market-overview {
  margin-bottom: 40px;
}

.overview-cards {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(250px, 1fr));
  gap: 20px;
  margin-bottom: 20px;
}

.overview-card {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
  padding: 25px;
  border-radius: 15px;
  display: flex;
  align-items: center;
  box-shadow: 0 8px 25px rgba(0,0,0,0.1);
  transition: transform 0.3s ease;
}

.overview-card:hover {
  transform: translateY(-5px);
}

.card-icon {
  font-size: 2.5rem;
  margin-right: 20px;
}

.card-content h3 {
  margin: 0 0 10px 0;
  font-size: 1.1rem;
  opacity: 0.9;
}

.card-content .price {
  font-size: 1.8rem;
  font-weight: bold;
  margin-bottom: 5px;
}

.card-content .change {
  font-size: 0.9rem;
  opacity: 0.8;
}

.quick-nav {
  margin-bottom: 40px;
}

.quick-nav h2 {
  color: #2c3e50;
  margin-bottom: 20px;
}

.nav-buttons {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(150px, 1fr));
  gap: 15px;
}

.nav-button {
  background: white;
  border: 2px solid #e9ecef;
  border-radius: 10px;
  padding: 20px;
  text-decoration: none;
  color: #2c3e50;
  text-align: center;
  transition: all 0.3s ease;
  display: flex;
  flex-direction: column;
  align-items: center;
}

.nav-button:hover {
  border-color: #3498db;
  background: #f8f9fa;
  transform: translateY(-2px);
  box-shadow: 0 4px 15px rgba(0,0,0,0.1);
}

.nav-icon {
  font-size: 2rem;
  margin-bottom: 10px;
}

.nav-text {
  font-weight: 500;
}

.hot-stocks {
  margin-bottom: 40px;
}

.hot-stocks h2 {
  color: #2c3e50;
  margin-bottom: 20px;
}

.stocks-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
  gap: 15px;
}

.stock-card {
  background: white;
  border: 1px solid #e9ecef;
  border-radius: 10px;
  padding: 20px;
  cursor: pointer;
  transition: all 0.3s ease;
  box-shadow: 0 2px 10px rgba(0,0,0,0.05);
}

.stock-card:hover {
  border-color: #3498db;
  transform: translateY(-2px);
  box-shadow: 0 4px 20px rgba(0,0,0,0.1);
}

.stock-header {
  margin-bottom: 15px;
}

.stock-code {
  font-size: 1.2rem;
  font-weight: bold;
  color: #2c3e50;
}

.stock-name {
  font-size: 0.9rem;
  color: #7f8c8d;
  margin-top: 5px;
}

.stock-price {
  text-align: right;
}

.current-price {
  font-size: 1.3rem;
  font-weight: bold;
  color: #2c3e50;
  margin-bottom: 5px;
}

.price-change {
  font-size: 0.9rem;
  font-weight: 500;
}

.price-change.positive {
  color: #e74c3c;
}

.price-change.negative {
  color: #27ae60;
}

.market-info {
  margin-bottom: 40px;
}

.market-info h2 {
  color: #2c3e50;
  margin-bottom: 20px;
}

.info-cards {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(250px, 1fr));
  gap: 20px;
}

.info-card {
  background: #f8f9fa;
  border-radius: 10px;
  padding: 25px;
  border-left: 4px solid #3498db;
}

.info-card h3 {
  color: #2c3e50;
  margin-bottom: 15px;
  font-size: 1.2rem;
}

.info-card p {
  margin: 8px 0;
  color: #555;
}

.info-card .note {
  font-size: 0.9rem;
  color: #7f8c8d;
  font-style: italic;
}

.loading {
  text-align: center;
  padding: 40px;
}

.spinner {
  width: 40px;
  height: 40px;
  border: 4px solid #f3f3f3;
  border-top: 4px solid #3498db;
  border-radius: 50%;
  animation: spin 1s linear infinite;
  margin: 0 auto 20px;
}

@keyframes spin {
  0% { transform: rotate(0deg); }
  100% { transform: rotate(360deg); }
}

@media (max-width: 768px) {
  .page-header h1 {
    font-size: 2rem;
  }
  
  .overview-cards {
    grid-template-columns: 1fr;
  }
  
  .nav-buttons {
    grid-template-columns: repeat(2, 1fr);
  }
  
  .stocks-grid {
    grid-template-columns: repeat(2, 1fr);
  }
}
</style>
