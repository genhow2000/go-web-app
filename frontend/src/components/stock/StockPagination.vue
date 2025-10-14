<template>
  <div class="pagination">
    <!-- 每頁筆數選擇 -->
    <div class="per-page-selector">
      <label>每頁顯示：</label>
      <select v-model="localPerPage" @change="onPerPageChange" class="per-page-select">
        <option value="10">10</option>
        <option value="20">20</option>
        <option value="50">50</option>
        <option value="100">100</option>
      </select>
      <span class="per-page-label">筆</span>
    </div>

    <!-- 分頁資訊 -->
    <div class="pagination-info">
      <span>第 {{ currentPage }} 頁，共 {{ totalPages }} 頁</span>
      <span class="total-count">（總共 {{ totalCount }} 筆）</span>
    </div>

    <!-- 分頁按鈕 -->
    <div class="pagination-buttons">
      <!-- 第一頁 -->
      <button 
        @click="goToPage(1)" 
        :disabled="currentPage === 1"
        class="page-btn first-page"
        :class="{ disabled: currentPage === 1 }"
      >
        首頁
      </button>

      <!-- 上一頁 -->
      <button 
        @click="goToPage(currentPage - 1)" 
        :disabled="currentPage === 1"
        class="page-btn prev-page"
        :class="{ disabled: currentPage === 1 }"
      >
        <i class="arrow">‹</i>
      </button>

      <!-- 頁碼按鈕 -->
      <div class="page-numbers">
        <template v-for="page in visiblePages" :key="page">
          <!-- 省略號 -->
          <span v-if="page === '...'" class="ellipsis">...</span>
          <!-- 頁碼按鈕 -->
          <button 
            v-else
            @click="goToPage(page)" 
            class="page-btn page-number"
            :class="{ active: page === currentPage }"
          >
            {{ page }}
          </button>
        </template>
      </div>

      <!-- 下一頁 -->
      <button 
        @click="goToPage(currentPage + 1)" 
        :disabled="currentPage === totalPages"
        class="page-btn next-page"
        :class="{ disabled: currentPage === totalPages }"
      >
        <i class="arrow">›</i>
      </button>

      <!-- 最後一頁 -->
      <button 
        @click="goToPage(totalPages)" 
        :disabled="currentPage === totalPages"
        class="page-btn last-page"
        :class="{ disabled: currentPage === totalPages }"
      >
        末頁
      </button>
    </div>

    <!-- 跳轉到指定頁面 -->
    <div class="page-jumper">
      <label>跳轉到：</label>
      <input 
        v-model="jumpToPage" 
        @keyup.enter="jumpToPageHandler"
        type="number" 
        :min="1" 
        :max="totalPages"
        class="jump-input"
        placeholder="頁碼"
      >
      <button @click="jumpToPageHandler" class="jump-btn">跳轉</button>
    </div>
  </div>
</template>

<script>
import { ref, computed, watch } from 'vue'

export default {
  name: 'StockPagination',
  props: {
    currentPage: {
      type: Number,
      required: true
    },
    totalPages: {
      type: Number,
      required: true
    },
    totalCount: {
      type: Number,
      required: true
    },
    perPage: {
      type: Number,
      default: 20
    }
  },
  emits: ['page-change', 'per-page-change'],
  setup(props, { emit }) {
    const localPerPage = ref(props.perPage)
    const jumpToPage = ref('')

    // 計算可見的頁碼範圍
    const visiblePages = computed(() => {
      const pages = []
      const current = props.currentPage
      const total = props.totalPages
      
      if (total <= 7) {
        // 總頁數少於等於7頁，顯示所有頁碼
        for (let i = 1; i <= total; i++) {
          pages.push(i)
        }
      } else {
        // 總頁數大於7頁，使用省略號
        if (current <= 4) {
          // 當前頁在前4頁
          for (let i = 1; i <= 5; i++) {
            pages.push(i)
          }
          pages.push('...')
          pages.push(total)
        } else if (current >= total - 3) {
          // 當前頁在後4頁
          pages.push(1)
          pages.push('...')
          for (let i = total - 4; i <= total; i++) {
            pages.push(i)
          }
        } else {
          // 當前頁在中間
          pages.push(1)
          pages.push('...')
          for (let i = current - 1; i <= current + 1; i++) {
            pages.push(i)
          }
          pages.push('...')
          pages.push(total)
        }
      }
      
      return pages
    })

    // 跳轉到指定頁面
    const goToPage = (page) => {
      if (page >= 1 && page <= props.totalPages && page !== props.currentPage) {
        emit('page-change', page)
      }
    }

    // 每頁筆數變更
    const onPerPageChange = () => {
      emit('per-page-change', localPerPage.value)
    }

    // 跳轉處理
    const jumpToPageHandler = () => {
      const page = parseInt(jumpToPage.value)
      if (page >= 1 && page <= props.totalPages) {
        goToPage(page)
        jumpToPage.value = ''
      }
    }

    // 監聽 perPage 變化
    watch(() => props.perPage, (newPerPage) => {
      localPerPage.value = newPerPage
    })

    return {
      localPerPage,
      jumpToPage,
      visiblePages,
      goToPage,
      onPerPageChange,
      jumpToPageHandler
    }
  }
}
</script>

<style scoped>
.pagination {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 20px;
  padding: 20px;
  background: white;
  border-radius: 8px;
  box-shadow: 0 2px 8px rgba(0,0,0,0.1);
  flex-wrap: wrap;
}

.per-page-selector {
  display: flex;
  align-items: center;
  gap: 8px;
  white-space: nowrap;
}

.per-page-selector label {
  color: #4a5568;
  font-size: 14px;
  font-weight: 500;
}

.per-page-select {
  padding: 6px 10px;
  border: 2px solid #e2e8f0;
  border-radius: 6px;
  font-size: 14px;
  background: white;
  cursor: pointer;
  transition: border-color 0.3s ease;
}

.per-page-select:focus {
  outline: none;
  border-color: #4299e1;
}

.per-page-label {
  color: #718096;
  font-size: 14px;
}

.pagination-info {
  display: flex;
  align-items: center;
  gap: 5px;
  color: #4a5568;
  font-size: 14px;
  white-space: nowrap;
}

.total-count {
  color: #718096;
}

.pagination-buttons {
  display: flex;
  align-items: center;
  gap: 4px;
}

.page-btn {
  display: flex;
  align-items: center;
  justify-content: center;
  min-width: 36px;
  height: 36px;
  padding: 0 12px;
  border: 2px solid #e2e8f0;
  background: white;
  color: #4a5568;
  border-radius: 6px;
  cursor: pointer;
  font-size: 14px;
  font-weight: 500;
  transition: all 0.3s ease;
  user-select: none;
}

.page-btn:hover:not(.disabled) {
  border-color: #4299e1;
  color: #4299e1;
  background: #f7fafc;
}

.page-btn.active {
  background: #4299e1;
  border-color: #4299e1;
  color: white;
}

.page-btn.disabled {
  opacity: 0.5;
  cursor: not-allowed;
  background: #f7fafc;
}

.page-btn.disabled:hover {
  border-color: #e2e8f0;
  color: #4a5568;
  background: #f7fafc;
}

.page-numbers {
  display: flex;
  align-items: center;
  gap: 4px;
}

.ellipsis {
  display: flex;
  align-items: center;
  justify-content: center;
  min-width: 36px;
  height: 36px;
  color: #718096;
  font-size: 14px;
  font-weight: 500;
}

.arrow {
  font-size: 16px;
  font-weight: bold;
}

.page-jumper {
  display: flex;
  align-items: center;
  gap: 8px;
  white-space: nowrap;
}

.page-jumper label {
  color: #4a5568;
  font-size: 14px;
  font-weight: 500;
}

.jump-input {
  width: 60px;
  padding: 6px 8px;
  border: 2px solid #e2e8f0;
  border-radius: 6px;
  font-size: 14px;
  text-align: center;
  transition: border-color 0.3s ease;
}

.jump-input:focus {
  outline: none;
  border-color: #4299e1;
}

.jump-btn {
  padding: 6px 12px;
  background: #4299e1;
  color: white;
  border: none;
  border-radius: 6px;
  cursor: pointer;
  font-size: 14px;
  font-weight: 500;
  transition: background-color 0.3s ease;
}

.jump-btn:hover {
  background: #3182ce;
}

/* 響應式設計 */
@media (max-width: 768px) {
  .pagination {
    flex-direction: column;
    gap: 15px;
    align-items: stretch;
  }

  .per-page-selector,
  .pagination-info,
  .page-jumper {
    justify-content: center;
  }

  .pagination-buttons {
    justify-content: center;
    flex-wrap: wrap;
  }

  .page-btn {
    min-width: 32px;
    height: 32px;
    padding: 0 8px;
    font-size: 13px;
  }

  .first-page,
  .last-page {
    display: none;
  }
}

@media (max-width: 480px) {
  .pagination {
    padding: 15px;
  }

  .page-numbers {
    gap: 2px;
  }

  .page-btn {
    min-width: 28px;
    height: 28px;
    padding: 0 6px;
    font-size: 12px;
  }

  .ellipsis {
    min-width: 28px;
    height: 28px;
    font-size: 12px;
  }
}
</style>
