<template>
  <div class="merchant-dashboard">
    <Header />
    
    <div class="container">
      <div class="welcome-section">
        <h1 class="welcome-title">歡迎來到商戶後台</h1>
        <p class="welcome-subtitle">管理您的商戶業務，查看統計數據，處理訂單和客戶</p>
      </div>

      <div class="stats-grid">
        <div class="stat-card">
          <div class="stat-icon orders">📦</div>
          <div class="stat-title">今日訂單</div>
          <div class="stat-value">24</div>
        </div>
        <div class="stat-card">
          <div class="stat-icon revenue">💰</div>
          <div class="stat-title">今日營收</div>
          <div class="stat-value">$2,450</div>
        </div>
        <div class="stat-card">
          <div class="stat-icon customers">👥</div>
          <div class="stat-title">活躍客戶</div>
          <div class="stat-value">156</div>
        </div>
        <div class="stat-card">
          <div class="stat-icon products">📦</div>
          <div class="stat-title">商品總數</div>
          <div class="stat-value">89</div>
        </div>
      </div>

      <div class="features-grid">
        <div class="feature-card">
          <h3 class="feature-title">訂單管理</h3>
          <p class="feature-desc">查看和管理所有訂單，處理訂單狀態更新</p>
          <a href="#" class="feature-btn" @click="showComingSoon('訂單管理')">管理訂單</a>
        </div>
        <div class="feature-card">
          <h3 class="feature-title">商品管理</h3>
          <p class="feature-desc">添加、編輯和刪除商品，管理庫存</p>
          <router-link to="/merchant/products" class="feature-btn">管理商品</router-link>
        </div>
        <div class="feature-card">
          <h3 class="feature-title">客戶管理</h3>
          <p class="feature-desc">查看客戶信息，管理客戶關係</p>
          <a href="#" class="feature-btn" @click="showComingSoon('客戶管理')">管理客戶</a>
        </div>
        <div class="feature-card">
          <h3 class="feature-title">財務報表</h3>
          <p class="feature-desc">查看銷售報表，分析營收趨勢</p>
          <a href="#" class="feature-btn" @click="showComingSoon('財務報表')">查看報表</a>
        </div>
        <div class="feature-card">
          <h3 class="feature-title">設定</h3>
          <p class="feature-desc">管理商戶設定，更新個人信息</p>
          <a href="#" class="feature-btn" @click="showComingSoon('商戶設定')">商戶設定</a>
        </div>
        <div class="feature-card">
          <h3 class="feature-title">支援</h3>
          <p class="feature-desc">聯繫技術支援，查看幫助文檔</p>
          <a href="#" class="feature-btn" @click="showComingSoon('技術支援')">取得支援</a>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import Header from '@/components/common/Header.vue'

export default {
  name: 'MerchantDashboard',
  components: {
    Header
  },
  setup() {
    const showComingSoon = (featureName) => {
      alert(`🚧 ${featureName} 功能正在開發中，敬請期待！\n\n目前可用的功能：\n✅ 商品管理 - 完整功能\n✅ 用戶認證 - 完整功能\n\n即將推出：\n⏳ 訂單管理\n⏳ 客戶管理\n⏳ 財務報表\n⏳ 商戶設定`)
    }

    return {
      showComingSoon
    }
  }
}
</script>

<style scoped>
.merchant-dashboard {
  font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif;
  background: #f7fafc;
  color: #2d3748;
  min-height: 100vh;
}

.container {
  max-width: 1200px;
  margin: 0 auto;
  padding: 20px;
}

.welcome-section {
  background: white;
  border-radius: 15px;
  padding: 30px;
  margin-bottom: 30px;
  box-shadow: 0 4px 15px rgba(0,0,0,0.1);
}

.welcome-title {
  font-size: 2rem;
  color: #2d3748;
  margin-bottom: 10px;
}

.welcome-subtitle {
  color: #718096;
  font-size: 1.1rem;
}

.stats-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(250px, 1fr));
  gap: 20px;
  margin-bottom: 30px;
}

.stat-card {
  background: white;
  border-radius: 15px;
  padding: 25px;
  box-shadow: 0 4px 15px rgba(0,0,0,0.1);
  transition: transform 0.3s ease;
}

.stat-card:hover {
  transform: translateY(-5px);
}

.stat-icon {
  width: 50px;
  height: 50px;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 1.5rem;
  margin-bottom: 15px;
}

.stat-icon.orders { background: linear-gradient(135deg, #4299e1, #3182ce); }
.stat-icon.revenue { background: linear-gradient(135deg, #48bb78, #38a169); }
.stat-icon.customers { background: linear-gradient(135deg, #ed8936, #dd6b20); }
.stat-icon.products { background: linear-gradient(135deg, #9f7aea, #805ad5); }

.stat-title {
  font-size: 0.9rem;
  color: #718096;
  margin-bottom: 5px;
}

.stat-value {
  font-size: 2rem;
  font-weight: bold;
  color: #2d3748;
}

.features-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(300px, 1fr));
  gap: 20px;
}

.feature-card {
  background: white;
  border-radius: 15px;
  padding: 25px;
  box-shadow: 0 4px 15px rgba(0,0,0,0.1);
  transition: transform 0.3s ease;
}

.feature-card:hover {
  transform: translateY(-5px);
}

.feature-title {
  font-size: 1.2rem;
  font-weight: bold;
  color: #2d3748;
  margin-bottom: 10px;
}

.feature-desc {
  color: #718096;
  margin-bottom: 20px;
}

.feature-btn {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
  border: none;
  padding: 10px 20px;
  border-radius: 8px;
  cursor: pointer;
  text-decoration: none;
  display: inline-block;
  transition: transform 0.3s ease;
}

.feature-btn:hover {
  transform: translateY(-2px);
}

@media (max-width: 768px) {
  .stats-grid {
    grid-template-columns: 1fr;
  }
}
</style>
