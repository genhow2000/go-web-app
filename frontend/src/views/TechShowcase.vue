<template>
  <div class="tech-showcase">
    <a href="/" class="back-to-mall">🏪 返回商城首頁</a>
    
    <div class="container">
      <div class="header">
        <h1>🚀 阿和 Go 全端系統展示</h1>
        <p>現代化 Web 應用技術能力展示平台</p>
      </div>

      <div class="status-grid">
        <StatusCard 
          v-for="status in statusCards" 
          :key="status.type"
          :status="status"
          @click="showStatusDetails(status.type)"
        />
      </div>

      <div class="features">
        <h2>✨ 系統功能(可點擊查看MD)</h2>
        <div class="feature-grid">
          <FeatureItem 
            v-for="feature in features" 
            :key="feature.id"
            :feature="feature"
            @click="navigateToDocs(feature.docPath)"
          />
        </div>
      </div>

      <div class="tech-stack">
        <h2>🛠️ 技術棧</h2>
        <div class="tech-categories">
          <TechCategory 
            v-for="category in techCategories" 
            :key="category.title"
            :category="category"
          />
        </div>
      </div>

      <div class="cta-section">
        <h2>🎯 技術能力展示</h2>
        <p>這個系統展示了完整的全端開發能力，包括後端 API 設計、前端界面開發、資料庫設計、雲端部署等技術。</p>
        <p><strong>特別亮點：</strong>由於 Go 沒有內建的 Migration 和 Seeder 系統，我們完全自製了這些功能，展現了深度技術理解和自製能力。</p>
        <div class="cta-buttons">
          <a href="/merchant/login" class="btn btn-primary">商戶登入</a>
          <a href="/admin/login" class="btn btn-secondary">管理員登入</a>
          <a href="/admin/db/login" class="btn btn-secondary">資料庫管理</a>
          <a href="/health" class="btn btn-secondary">API 狀態</a>
        </div>
      </div>

      <div class="footer">
        <p>© 2025 Go 全端系統展示 | 技術能力展示平台</p>
      </div>
    </div>

    <!-- 狀態詳情模態框 -->
    <StatusModal 
      v-if="showModal"
      :type="selectedType"
      :data="modalData"
      @close="closeModal"
    />
  </div>
</template>

<script>
import { ref, onMounted } from 'vue'
import StatusCard from '@/components/tech/StatusCard.vue'
import StatusModal from '@/components/tech/StatusModal.vue'
import FeatureItem from '@/components/tech/FeatureItem.vue'
import TechCategory from '@/components/tech/TechCategory.vue'
import api from '@/services/api'

export default {
  name: 'TechShowcase',
  components: {
    StatusCard,
    StatusModal,
    FeatureItem,
    TechCategory
  },
  setup() {
    const showModal = ref(false)
    const selectedType = ref(null)
    const modalData = ref(null)

    const statusCards = ref([
      {
        type: 'system',
        title: '系統狀態',
        value: '運行中',
        detail: '服務器運行時間: 24小時',
        icon: 'status-online',
        isRealtime: true
      },
      {
        type: 'database',
        title: '資料庫',
        value: '已連接',
        detail: 'SQLite 資料庫 | 遷移版本: 001',
        icon: 'status-database',
        isRealtime: true
      },
      {
        type: 'api',
        title: 'API 服務',
        value: '正常',
        detail: 'RESTful API | 版本: 2.0.0',
        icon: 'status-api',
        isRealtime: true
      },
      {
        type: 'cloud',
        title: '雲端部署',
        value: '已部署',
        detail: 'Google Cloud Run | 區域: asia-east1',
        icon: 'status-cloud',
        isRealtime: true
      }
    ])

    const features = ref([
      {
        id: 1,
        icon: '👥',
        title: '用戶管理',
        description: '完整的用戶註冊、登入、權限管理系統',
        docPath: '/docs/user-management'
      },
      {
        id: 2,
        icon: '🔐',
        title: '安全認證',
        description: 'JWT Token 認證、密碼加密、會話管理',
        docPath: '/docs/auth'
      },
      {
        id: 3,
        icon: '📊',
        title: '管理後台',
        description: '功能完整的管理界面，支援 CRUD 操作',
        docPath: '/docs/admin'
      },
      {
        id: 4,
        icon: '🗄️',
        title: '資料庫管理',
        description: '內建資料庫管理工具，支援 SQL 查詢',
        docPath: '/docs/db-management'
      },
      {
        id: 5,
        icon: '📈',
        title: '系統監控',
        description: '實時系統狀態監控、日誌查看',
        docPath: '/docs/monitoring'
      },
      {
        id: 6,
        icon: '🌐',
        title: 'API 服務',
        description: 'RESTful API 設計，支援前後端分離',
        docPath: '/docs/api'
      },
      {
        id: 7,
        icon: '⚡',
        title: 'Redis 快取',
        description: '高效能記憶體快取，提升系統性能',
        docPath: '/docs/redis'
      },
      {
        id: 8,
        icon: '📄',
        title: 'MongoDB 文檔',
        description: 'NoSQL 文檔資料庫，支援複雜資料結構',
        docPath: '/docs/mongodb'
      },
      {
        id: 9,
        icon: '🔄',
        title: '自製 Migration',
        description: '完全自製的資料庫遷移系統，支援版本控制',
        docPath: '/docs/migration'
      },
      {
        id: 10,
        icon: '🌱',
        title: '自製 Seeder',
        description: '自動化測試數據生成，支援重複執行保護',
        docPath: '/docs/seeder'
      },
      {
        id: 11,
        icon: '🤖',
        title: 'AI 智能聊天',
        description: '整合 Groq 和 Gemini AI，提供智能對話服務',
        docPath: '/docs/ai-chat'
      }
    ])

    const techCategories = ref([
      {
        title: '後端技術',
        items: [
          'Go 1.21',
          'Gin Web Framework',
          'JWT 認證',
          'bcrypt 密碼加密',
          'SQLite 資料庫',
          'Redis 快取',
          'MongoDB 文檔資料庫',
          'Groq AI API',
          'Google Gemini AI'
        ]
      },
      {
        title: '前端技術',
        items: [
          'Vue.js 3',
          'Vite 構建工具',
          'Element Plus UI',
          'Pinia 狀態管理',
          'Vue Router',
          '響應式設計',
          '現代化 UI'
        ]
      },
      {
        title: '部署與運維',
        items: [
          'Docker 容器化',
          'Google Cloud Run',
          'CI/CD 自動部署',
          '自製 Migration 系統',
          '自製 Seeder 系統',
          '日誌監控系統'
        ]
      },
      {
        title: '資料庫管理',
        items: [
          'SQLite 資料庫',
          '版本控制遷移',
          '自動測試數據',
          '重複執行保護',
          '執行記錄追蹤'
        ]
      }
    ])

    const showStatusDetails = async (type) => {
      selectedType.value = type
      showModal.value = true
      
      try {
        const response = await api.get(`/api/status/${type}`)
        modalData.value = response.data
      } catch (error) {
        console.error('載入狀態詳情失敗:', error)
        modalData.value = null
      }
    }

    const closeModal = () => {
      showModal.value = false
      selectedType.value = null
      modalData.value = null
    }

    const navigateToDocs = (docPath) => {
      window.open(docPath, '_blank')
    }

    // 實時更新系統狀態
    const updateSystemStatus = async () => {
      try {
        const response = await api.get('/health')
        const data = response.data
        
        const systemCard = statusCards.value.find(card => card.type === 'system')
        if (systemCard) {
          systemCard.value = data.status === 'healthy' ? '運行中' : '異常'
        }
        
        const dbCard = statusCards.value.find(card => card.type === 'database')
        if (dbCard) {
          dbCard.value = data.database === 'connected' ? '已連接' : '斷開'
        }
      } catch (error) {
        console.error('狀態更新失敗:', error)
      }
    }

    onMounted(() => {
      updateSystemStatus()
      // 每30秒更新一次狀態
      setInterval(updateSystemStatus, 30000)
    })

    return {
      showModal,
      selectedType,
      modalData,
      statusCards,
      features,
      techCategories,
      showStatusDetails,
      closeModal,
      navigateToDocs
    }
  }
}
</script>

<style scoped>
.tech-showcase {
  font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, 'Helvetica Neue', Arial, sans-serif;
  line-height: 1.6;
  color: #333;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  min-height: 100vh;
}

.back-to-mall {
  position: fixed;
  top: 20px;
  left: 20px;
  background: rgba(255, 255, 255, 0.9);
  color: #667eea;
  padding: 10px 20px;
  border-radius: 25px;
  text-decoration: none;
  font-weight: bold;
  transition: all 0.3s ease;
  z-index: 100;
}

.back-to-mall:hover {
  background: white;
  transform: translateY(-2px);
  box-shadow: 0 5px 15px rgba(0,0,0,0.2);
}

.container {
  max-width: 1200px;
  margin: 0 auto;
  padding: 20px;
}

.header {
  text-align: center;
  color: white;
  margin-bottom: 40px;
}

.header h1 {
  font-size: 3rem;
  margin-bottom: 10px;
  text-shadow: 2px 2px 4px rgba(0,0,0,0.3);
}

.header p {
  font-size: 1.2rem;
  opacity: 0.9;
}

.status-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(300px, 1fr));
  gap: 20px;
  margin-bottom: 40px;
}

.features {
  background: white;
  border-radius: 15px;
  padding: 30px;
  margin-bottom: 40px;
  box-shadow: 0 10px 30px rgba(0,0,0,0.1);
}

.features h2 {
  color: #2d3748;
  margin-bottom: 20px;
  text-align: center;
}

.feature-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
  gap: 20px;
}

.tech-stack {
  background: white;
  border-radius: 15px;
  padding: 30px;
  margin-bottom: 40px;
  box-shadow: 0 10px 30px rgba(0,0,0,0.1);
}

.tech-stack h2 {
  color: #2d3748;
  margin-bottom: 20px;
  text-align: center;
}

.tech-categories {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(250px, 1fr));
  gap: 20px;
}

.cta-section {
  text-align: center;
  color: white;
}

.cta-buttons {
  display: flex;
  gap: 20px;
  justify-content: center;
  flex-wrap: wrap;
  margin-top: 20px;
}

.btn {
  display: inline-block;
  padding: 12px 30px;
  border-radius: 25px;
  text-decoration: none;
  font-weight: bold;
  transition: all 0.3s ease;
  border: 2px solid transparent;
}

.btn-primary {
  background: white;
  color: #667eea;
  border-color: white;
}

.btn-primary:hover {
  background: transparent;
  color: white;
  border-color: white;
}

.btn-secondary {
  background: transparent;
  color: white;
  border-color: white;
}

.btn-secondary:hover {
  background: white;
  color: #667eea;
}

.footer {
  text-align: center;
  color: white;
  margin-top: 40px;
  opacity: 0.8;
}

@media (max-width: 768px) {
  .header h1 {
    font-size: 2rem;
  }
  
  .cta-buttons {
    flex-direction: column;
    align-items: center;
  }
}
</style>
