<template>
  <div class="modal" @click="handleBackdropClick">
    <div class="modal-content" @click.stop>
      <div class="modal-header">
        <h2>{{ getModalTitle() }}</h2>
        <span class="close" @click="$emit('close')">&times;</span>
      </div>
      <div class="modal-body">
        <div v-if="!data" class="loading">載入中...</div>
        <div v-else>
          <div v-for="section in data.sections" :key="section.title" class="detail-section">
            <h3>{{ section.title }}</h3>
            <div v-for="item in section.items" :key="item.label" class="detail-item">
              <span class="detail-label">{{ item.label }}</span>
              <span 
                v-if="item.status" 
                :class="['status-badge', `status-${item.status}`]"
              >
                {{ item.value }}
              </span>
              <span v-else class="detail-value">{{ item.value }}</span>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
export default {
  name: 'StatusModal',
  props: {
    type: {
      type: String,
      required: true
    },
    data: {
      type: Object,
      default: null
    }
  },
  emits: ['close'],
  setup(props) {
    const getModalTitle = () => {
      const titles = {
        'system': '系統狀態詳情',
        'database': '資料庫詳情',
        'api': 'API 服務詳情',
        'cloud': '雲端部署詳情'
      }
      return titles[props.type] || '系統詳情'
    }

    const handleBackdropClick = (event) => {
      if (event.target.classList.contains('modal')) {
        // 點擊背景關閉模態框
        // 這裡可以觸發關閉事件
      }
    }

    return {
      getModalTitle,
      handleBackdropClick
    }
  }
}
</script>

<style scoped>
.modal {
  display: block;
  position: fixed;
  z-index: 1000;
  left: 0;
  top: 0;
  width: 100%;
  height: 100%;
  background-color: rgba(0,0,0,0.5);
  backdrop-filter: blur(5px);
}

.modal-content {
  background-color: white;
  margin: 5% auto;
  padding: 0;
  border-radius: 20px;
  width: 90%;
  max-width: 800px;
  max-height: 80vh;
  overflow-y: auto;
  box-shadow: 0 20px 60px rgba(0,0,0,0.3);
  animation: modalSlideIn 0.3s ease;
}

@keyframes modalSlideIn {
  from {
    opacity: 0;
    transform: translateY(-50px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

.modal-header {
  background: linear-gradient(135deg, #667eea, #764ba2);
  color: white;
  padding: 20px 30px;
  border-radius: 20px 20px 0 0;
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.modal-header h2 {
  margin: 0;
  font-size: 1.5rem;
}

.close {
  color: white;
  font-size: 28px;
  font-weight: bold;
  cursor: pointer;
  transition: opacity 0.3s;
}

.close:hover {
  opacity: 0.7;
}

.modal-body {
  padding: 30px;
}

.detail-section {
  margin-bottom: 25px;
}

.detail-section h3 {
  color: #2d3748;
  margin-bottom: 15px;
  font-size: 1.2rem;
  border-bottom: 2px solid #e2e8f0;
  padding-bottom: 8px;
}

.detail-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 10px 0;
  border-bottom: 1px solid #f7fafc;
}

.detail-item:last-child {
  border-bottom: none;
}

.detail-label {
  font-weight: 600;
  color: #4a5568;
}

.detail-value {
  color: #2d3748;
  font-family: 'Courier New', monospace;
  background: #f7fafc;
  padding: 4px 8px;
  border-radius: 4px;
}

.status-badge {
  display: inline-block;
  padding: 4px 12px;
  border-radius: 20px;
  font-size: 0.8rem;
  font-weight: 600;
  text-transform: uppercase;
}

.status-online { background: #c6f6d5; color: #22543d; }
.status-offline { background: #fed7d7; color: #742a2a; }
.status-warning { background: #fef5e7; color: #744210; }

.loading {
  text-align: center;
  padding: 2rem;
  color: #718096;
}

@media (max-width: 768px) {
  .modal-content {
    width: 95%;
    margin: 10% auto;
  }

  .modal-header {
    padding: 15px 20px;
  }

  .modal-body {
    padding: 20px;
  }
}
</style>
