<template>
  <div class="customer-dashboard">
    <Header />
    
    <div class="container">
      <div class="welcome-section">
        <h2>歡迎回來！</h2>
        <p>這是您的個人客戶儀表板，您可以在這裡管理您的帳戶和查看相關信息。</p>
      </div>

      <div class="stats-grid">
        <div class="stat-card">
          <div class="stat-number">{{ userStats.totalOrders }}</div>
          <div class="stat-label">總訂單數</div>
        </div>
        <div class="stat-card">
          <div class="stat-number">{{ userStats.activeChats }}</div>
          <div class="stat-label">活躍對話</div>
        </div>
        <div class="stat-card">
          <div class="stat-number">{{ userStats.loginCount }}</div>
          <div class="stat-label">登入次數</div>
        </div>
        <div class="stat-card">
          <div class="stat-number">{{ userStats.accountAge }}</div>
          <div class="stat-label">帳戶天數</div>
        </div>
      </div>

      <div class="dashboard-grid">
        <div class="dashboard-card">
          <div class="card-header">
            <div class="card-icon">📦</div>
            <div class="card-title">我的訂單</div>
          </div>
          <div class="card-content">
            查看和管理您的所有訂單，包括訂單狀態、配送信息等。
          </div>
          <a href="#" class="card-action">查看訂單</a>
        </div>

        <div class="dashboard-card">
          <div class="card-header">
            <div class="card-icon">💬</div>
            <div class="card-title">客服對話</div>
          </div>
          <div class="card-content">
            與客服人員進行即時對話，獲得幫助和支持。
          </div>
          <a href="/api/chat/conversations" class="card-action">開始對話</a>
        </div>

        <div class="dashboard-card">
          <div class="card-header">
            <div class="card-icon">👤</div>
            <div class="card-title">個人資料</div>
          </div>
          <div class="card-content">
            更新您的個人信息、密碼和偏好設置。
          </div>
          <a href="#" class="card-action">編輯資料</a>
        </div>

        <div class="dashboard-card">
          <div class="card-header">
            <div class="card-icon">🔔</div>
            <div class="card-title">通知中心</div>
          </div>
          <div class="card-content">
            查看系統通知、訂單更新和重要消息。
          </div>
          <a href="#" class="card-action">查看通知</a>
        </div>
      </div>

      <div class="recent-activity">
        <h3>最近活動</h3>
        <div class="activity-item">
          <div class="activity-icon">✓</div>
          <div class="activity-content">
            <div class="activity-title">成功登入</div>
            <div class="activity-time">剛剛</div>
          </div>
        </div>
        <div class="activity-item">
          <div class="activity-icon">📧</div>
          <div class="activity-content">
            <div class="activity-title">歡迎使用客戶系統</div>
            <div class="activity-time">今天</div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import { ref, onMounted, computed } from 'vue'
import Header from '@/components/common/Header.vue'
import { useAuthStore } from '@/stores/auth'
import api from '@/services/api'

export default {
  name: 'CustomerDashboard',
  components: {
    Header
  },
  setup() {
    const authStore = useAuthStore()
    
    const userStats = ref({
      totalOrders: 0,
      activeChats: 0,
      loginCount: 0,
      accountAge: 0
    })

    const user = computed(() => authStore.user)

    const loadUserInfo = async () => {
      if (!authStore.isAuthenticated) {
        return
      }
      
      try {
        const response = await api.get('/customer/profile')
        const userData = response.data
        
        userStats.value.loginCount = userData.login_count || 0
        
        // 計算帳戶天數
        if (userData.created_at) {
          const createdDate = new Date(userData.created_at)
          const now = new Date()
          const diffTime = Math.abs(now - createdDate)
          const diffDays = Math.ceil(diffTime / (1000 * 60 * 60 * 24))
          userStats.value.accountAge = diffDays
        }
      } catch (error) {
        console.error('載入用戶信息失敗:', error)
      }
    }

    onMounted(() => {
      if (authStore.isAuthenticated) {
        loadUserInfo()
      } else {
        // 如果認證狀態未準備好，等待一下再重試
        setTimeout(() => {
          if (authStore.isAuthenticated) {
            loadUserInfo()
          }
        }, 1000)
      }
    })

    return {
      userStats,
      user
    }
  }
}
</script>

<style scoped>
.customer-dashboard {
  font-family: 'Segoe UI', Tahoma, Geneva, Verdana, sans-serif;
  background-color: #f5f7fa;
  line-height: 1.6;
  min-height: 100vh;
}

.container {
  max-width: 1200px;
  margin: 2rem auto;
  padding: 0 2rem;
}

.welcome-section {
  background: white;
  padding: 2rem;
  border-radius: 10px;
  box-shadow: 0 2px 10px rgba(0,0,0,0.1);
  margin-bottom: 2rem;
}

.welcome-section h2 {
  color: #333;
  margin-bottom: 0.5rem;
}

.welcome-section p {
  color: #666;
}

.stats-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
  gap: 1rem;
  margin-bottom: 2rem;
}

.stat-card {
  background: white;
  padding: 1.5rem;
  border-radius: 10px;
  box-shadow: 0 2px 10px rgba(0,0,0,0.1);
  text-align: center;
}

.stat-number {
  font-size: 2rem;
  font-weight: bold;
  color: #667eea;
  margin-bottom: 0.5rem;
}

.stat-label {
  color: #666;
  font-size: 0.9rem;
}

.dashboard-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(300px, 1fr));
  gap: 2rem;
  margin-bottom: 2rem;
}

.dashboard-card {
  background: white;
  padding: 2rem;
  border-radius: 10px;
  box-shadow: 0 2px 10px rgba(0,0,0,0.1);
  transition: transform 0.3s ease;
}

.dashboard-card:hover {
  transform: translateY(-5px);
}

.card-header {
  display: flex;
  align-items: center;
  margin-bottom: 1rem;
}

.card-icon {
  width: 40px;
  height: 40px;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  border-radius: 8px;
  display: flex;
  align-items: center;
  justify-content: center;
  color: white;
  font-size: 1.2rem;
  margin-right: 1rem;
}

.card-title {
  font-size: 1.2rem;
  font-weight: 600;
  color: #333;
}

.card-content {
  color: #666;
  margin-bottom: 1rem;
}

.card-action {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
  border: none;
  padding: 0.75rem 1.5rem;
  border-radius: 5px;
  cursor: pointer;
  text-decoration: none;
  display: inline-block;
  transition: transform 0.2s ease;
}

.card-action:hover {
  transform: translateY(-2px);
}

.recent-activity {
  background: white;
  padding: 2rem;
  border-radius: 10px;
  box-shadow: 0 2px 10px rgba(0,0,0,0.1);
}

.recent-activity h3 {
  margin-bottom: 1rem;
  color: #333;
}

.activity-item {
  display: flex;
  align-items: center;
  padding: 1rem 0;
  border-bottom: 1px solid #eee;
}

.activity-item:last-child {
  border-bottom: none;
}

.activity-icon {
  width: 30px;
  height: 30px;
  background: #f0f0f0;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  margin-right: 1rem;
  font-size: 0.8rem;
}

.activity-content {
  flex: 1;
}

.activity-title {
  font-weight: 500;
  color: #333;
  margin-bottom: 0.25rem;
}

.activity-time {
  font-size: 0.8rem;
  color: #666;
}

@media (max-width: 768px) {
  .container {
    padding: 0 1rem;
  }

  .dashboard-grid {
    grid-template-columns: 1fr;
  }

  .stats-grid {
    grid-template-columns: repeat(2, 1fr);
  }
}
</style>
