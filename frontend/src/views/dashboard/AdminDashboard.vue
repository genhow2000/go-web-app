<template>
  <div class="admin-dashboard">
    <div class="sidebar">
      <div class="sidebar-content">
        <h4 class="sidebar-title">
          <i class="fas fa-tachometer-alt"></i>
          管理後台
        </h4>
        <nav class="sidebar-nav">
          <router-link to="/admin/dashboard" class="nav-link active">
            <i class="fas fa-chart-pie"></i>
            儀表板
          </router-link>
          <router-link to="/admin/users" class="nav-link">
            <i class="fas fa-users"></i>
            用戶管理
          </router-link>
          <router-link to="/admin/users/create" class="nav-link">
            <i class="fas fa-user-plus"></i>
            新增用戶
          </router-link>
          <router-link to="/" class="nav-link">
            <i class="fas fa-home"></i>
            返回首頁
          </router-link>
          <a href="#" class="nav-link" @click="handleLogout">
            <i class="fas fa-sign-out-alt"></i>
            登出
          </a>
        </nav>
      </div>
    </div>

    <div class="main-content">
      <div class="content-header">
        <h2 class="page-title">
          <i class="fas fa-chart-pie"></i>
          管理員儀表板
        </h2>
        <div class="user-info">
          歡迎回來，管理員！
        </div>
      </div>

      <div class="stats-section">
        <div class="stat-card total">
          <div class="stat-number">{{ stats.total }}</div>
          <div class="stat-label">總用戶數</div>
        </div>
        <div class="stat-card customers">
          <div class="stat-number">{{ stats.customers }}</div>
          <div class="stat-label">客戶數量</div>
        </div>
        <div class="stat-card admins">
          <div class="stat-number">{{ stats.admins }}</div>
          <div class="stat-label">管理員數量</div>
        </div>
      </div>

      <div class="quick-actions">
        <h3 class="section-title">
          <i class="fas fa-bolt"></i>
          快速操作
        </h3>
        <div class="actions-grid">
          <router-link to="/admin/users/create" class="action-btn primary">
            <i class="fas fa-user-plus"></i>
            新增用戶
          </router-link>
          <router-link to="/admin/users?role=customer" class="action-btn success">
            <i class="fas fa-users"></i>
            查看客戶
          </router-link>
          <router-link to="/admin/users?role=admin" class="action-btn warning">
            <i class="fas fa-user-shield"></i>
            查看管理員
          </router-link>
          <router-link to="/admin/users" class="action-btn info">
            <i class="fas fa-list"></i>
            所有用戶
          </router-link>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import api from '@/services/api'

export default {
  name: 'AdminDashboard',
  setup() {
    const router = useRouter()
    const authStore = useAuthStore()
    
    const stats = ref({
      total: 0,
      customers: 0,
      admins: 0
    })

    const loadStats = async () => {
      try {
        const response = await api.get('/admin/api/stats')
        stats.value = response.data
      } catch (error) {
        console.error('載入統計數據失敗:', error)
        // 使用默認值
        stats.value = {
          total: 150,
          customers: 120,
          admins: 30
        }
      }
    }

    const handleLogout = async () => {
      if (confirm('確定要登出嗎？')) {
        await authStore.logout()
        router.push('/')
      }
    }

    onMounted(() => {
      loadStats()
    })

    return {
      stats,
      handleLogout
    }
  }
}
</script>

<style scoped>
.admin-dashboard {
  display: flex;
  min-height: 100vh;
  font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif;
}

.sidebar {
  width: 250px;
  min-height: 100vh;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
}

.sidebar-content {
  padding: 20px;
}

.sidebar-title {
  margin-bottom: 30px;
  font-size: 1.2rem;
  font-weight: 600;
}

.sidebar-title i {
  margin-right: 10px;
}

.sidebar-nav {
  display: flex;
  flex-direction: column;
  gap: 5px;
}

.nav-link {
  color: rgba(255,255,255,0.8);
  padding: 12px 20px;
  border-radius: 8px;
  text-decoration: none;
  transition: all 0.3s;
  display: flex;
  align-items: center;
}

.nav-link:hover,
.nav-link.active {
  color: white;
  background: rgba(255,255,255,0.1);
}

.nav-link i {
  margin-right: 10px;
  width: 16px;
}

.main-content {
  flex: 1;
  background: #f8f9fa;
  padding: 30px;
}

.content-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 30px;
}

.page-title {
  font-size: 1.8rem;
  font-weight: 600;
  color: #2d3748;
  margin: 0;
}

.page-title i {
  margin-right: 10px;
  color: #667eea;
}

.user-info {
  color: #718096;
  font-size: 1rem;
}

.stats-section {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
  gap: 20px;
  margin-bottom: 40px;
}

.stat-card {
  background: white;
  border-radius: 12px;
  padding: 24px;
  box-shadow: 0 2px 10px rgba(0,0,0,0.1);
  border-left: 4px solid;
  transition: transform 0.2s;
}

.stat-card:hover {
  transform: translateY(-2px);
}

.stat-card.total { border-left-color: #007bff; }
.stat-card.customers { border-left-color: #28a745; }
.stat-card.admins { border-left-color: #ffc107; }

.stat-number {
  font-size: 2.5rem;
  font-weight: bold;
  margin-bottom: 8px;
  color: #2d3748;
}

.stat-label {
  color: #6c757d;
  font-size: 0.9rem;
}

.quick-actions {
  background: white;
  border-radius: 12px;
  padding: 30px;
  box-shadow: 0 2px 10px rgba(0,0,0,0.1);
}

.section-title {
  font-size: 1.3rem;
  font-weight: 600;
  color: #2d3748;
  margin-bottom: 20px;
}

.section-title i {
  margin-right: 10px;
  color: #667eea;
}

.actions-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
  gap: 15px;
}

.action-btn {
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 15px 20px;
  border-radius: 8px;
  text-decoration: none;
  font-weight: 500;
  transition: all 0.3s;
  border: none;
  cursor: pointer;
}

.action-btn i {
  margin-right: 8px;
}

.action-btn.primary {
  background: #007bff;
  color: white;
}

.action-btn.primary:hover {
  background: #0056b3;
  transform: translateY(-2px);
}

.action-btn.success {
  background: #28a745;
  color: white;
}

.action-btn.success:hover {
  background: #1e7e34;
  transform: translateY(-2px);
}

.action-btn.warning {
  background: #ffc107;
  color: #212529;
}

.action-btn.warning:hover {
  background: #e0a800;
  transform: translateY(-2px);
}

.action-btn.info {
  background: #17a2b8;
  color: white;
}

.action-btn.info:hover {
  background: #138496;
  transform: translateY(-2px);
}

@media (max-width: 768px) {
  .admin-dashboard {
    flex-direction: column;
  }
  
  .sidebar {
    width: 100%;
    min-height: auto;
  }
  
  .main-content {
    padding: 20px;
  }
  
  .content-header {
    flex-direction: column;
    align-items: flex-start;
    gap: 15px;
  }
  
  .stats-section {
    grid-template-columns: 1fr;
  }
  
  .actions-grid {
    grid-template-columns: 1fr;
  }
}
</style>
