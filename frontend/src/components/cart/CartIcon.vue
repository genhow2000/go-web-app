<template>
  <div class="cart-icon" @click="goToCart">
    <div class="cart-icon-container">
      <!-- 購物車圖標 -->
      <div class="cart-symbol">
        🛒
      </div>
      
      <!-- 商品數量徽章 -->
      <div v-if="itemCount > 0" class="cart-badge">
        {{ displayCount }}
      </div>
    </div>
    
    <!-- 購物車摘要（懸停時顯示） -->
    <div v-if="showTooltip && itemCount > 0" class="cart-tooltip">
      <div class="tooltip-content">
        <div class="tooltip-header">
          <span class="tooltip-title">購物車</span>
          <span class="tooltip-count">{{ itemCount }} 件商品</span>
        </div>
        
        <div class="tooltip-items">
          <div class="tooltip-simple">
            <p>點擊查看購物車詳情</p>
          </div>
        </div>
        
        <div class="tooltip-footer">
          <div class="tooltip-total">
            總計: NT$ {{ totalPrice.toLocaleString() }}
          </div>
          <button class="tooltip-checkout-btn">查看購物車</button>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import { ref, computed, onMounted, onUnmounted } from 'vue'
import { useRouter } from 'vue-router'
import { useCartStore } from '@/stores/cart'

export default {
  name: 'CartIcon',
  props: {
    showTooltip: {
      type: Boolean,
      default: true
    },
    size: {
      type: String,
      default: 'medium', // small, medium, large
      validator: (value) => ['small', 'medium', 'large'].includes(value)
    }
  },
  setup(props) {
    const router = useRouter()
    const cartStore = useCartStore()
    
    const showTooltipState = ref(false)
    const hoverTimeout = ref(null)

    // 計算屬性
    const itemCount = computed(() => cartStore.itemCount)
    const totalPrice = computed(() => cartStore.totalPrice)

    const displayCount = computed(() => {
      if (itemCount.value > 99) {
        return '99+'
      }
      return itemCount.value.toString()
    })

    const iconSize = computed(() => {
      const sizes = {
        small: '1.5rem',
        medium: '2rem',
        large: '2.5rem'
      }
      return sizes[props.size]
    })

    // 方法
    const goToCart = () => {
      router.push('/cart')
    }

    const handleMouseEnter = () => {
      if (hoverTimeout.value) {
        clearTimeout(hoverTimeout.value)
      }
      hoverTimeout.value = setTimeout(() => {
        showTooltipState.value = true
      }, 300)
    }

    const handleMouseLeave = () => {
      if (hoverTimeout.value) {
        clearTimeout(hoverTimeout.value)
        hoverTimeout.value = null
      }
      showTooltipState.value = false
    }

    // 生命週期
    onMounted(() => {
      // 購物車數據會通過其他組件載入，這裡不需要主動載入
      // 避免在認證狀態不穩定時觸發401錯誤
    })

    onUnmounted(() => {
      if (hoverTimeout.value) {
        clearTimeout(hoverTimeout.value)
      }
    })

    return {
      itemCount,
      totalPrice,
      displayCount,
      iconSize,
      showTooltipState,
      goToCart,
      handleMouseEnter,
      handleMouseLeave
    }
  }
}
</script>

<style scoped>
.cart-icon {
  position: relative;
  cursor: pointer;
  user-select: none;
}

.cart-icon-container {
  position: relative;
  display: flex;
  align-items: center;
  justify-content: center;
  width: 50px;
  height: 50px;
  background: #f8fafc;
  border-radius: 50%;
  transition: all 0.3s ease;
  border: 2px solid transparent;
}

.cart-icon:hover .cart-icon-container {
  background: #667eea;
  border-color: #5a67d8;
  transform: scale(1.05);
}

.cart-symbol {
  font-size: v-bind(iconSize);
  transition: transform 0.3s ease;
}

.cart-icon:hover .cart-symbol {
  transform: scale(1.1);
}

.cart-badge {
  position: absolute;
  top: -5px;
  right: -5px;
  min-width: 20px;
  height: 20px;
  background: #e53e3e;
  color: white;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 0.75rem;
  font-weight: 600;
  border: 2px solid white;
  animation: pulse 2s infinite;
}

@keyframes pulse {
  0% {
    transform: scale(1);
  }
  50% {
    transform: scale(1.1);
  }
  100% {
    transform: scale(1);
  }
}

.cart-tooltip {
  position: absolute;
  top: 100%;
  right: 0;
  margin-top: 10px;
  z-index: 1000;
  opacity: v-bind(showTooltipState ? 1 : 0);
  visibility: v-bind(showTooltipState ? 'visible' : 'hidden');
  transition: all 0.3s ease;
  transform: translateY(v-bind(showTooltipState ? '0' : '-10px'));
}

.tooltip-content {
  background: white;
  border-radius: 12px;
  box-shadow: 0 10px 25px rgba(0,0,0,0.15);
  border: 1px solid #e2e8f0;
  min-width: 300px;
  max-width: 400px;
  overflow: hidden;
}

.tooltip-header {
  padding: 16px 20px;
  background: #667eea;
  color: white;
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.tooltip-title {
  font-weight: 600;
  font-size: 1.1rem;
}

.tooltip-count {
  font-size: 0.9rem;
  opacity: 0.9;
}

.tooltip-items {
  padding: 16px 20px;
  max-height: 200px;
  overflow-y: auto;
}

.tooltip-item {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 8px 0;
  border-bottom: 1px solid #f1f5f9;
}

.tooltip-item:last-child {
  border-bottom: none;
}

.tooltip-item-image {
  width: 40px;
  height: 40px;
  object-fit: cover;
  border-radius: 6px;
  background: #f7fafc;
}

.tooltip-item-placeholder {
  width: 40px;
  height: 40px;
  display: flex;
  align-items: center;
  justify-content: center;
  background: #f7fafc;
  border-radius: 6px;
  font-size: 1.2rem;
}

.tooltip-item-info {
  flex: 1;
  min-width: 0;
}

.tooltip-item-name {
  font-size: 0.9rem;
  font-weight: 500;
  color: #2d3748;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.tooltip-item-quantity {
  font-size: 0.8rem;
  color: #718096;
}

.tooltip-more {
  text-align: center;
  padding: 8px 0;
  color: #718096;
  font-size: 0.85rem;
  font-style: italic;
}

.tooltip-footer {
  padding: 16px 20px;
  background: #f8fafc;
  border-top: 1px solid #e2e8f0;
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.tooltip-total {
  font-weight: 600;
  color: #2d3748;
  font-size: 1rem;
}

.tooltip-checkout-btn {
  background: #667eea;
  color: white;
  border: none;
  padding: 8px 16px;
  border-radius: 6px;
  font-size: 0.9rem;
  font-weight: 500;
  cursor: pointer;
  transition: background 0.2s;
}

.tooltip-checkout-btn:hover {
  background: #5a67d8;
}

/* 響應式設計 */
@media (max-width: 768px) {
  .cart-icon-container {
    width: 45px;
    height: 45px;
  }
  
  .cart-symbol {
    font-size: 1.5rem;
  }
  
  .tooltip-content {
    min-width: 280px;
    max-width: 320px;
  }
  
  .tooltip-header {
    padding: 12px 16px;
  }
  
  .tooltip-items {
    padding: 12px 16px;
  }
  
  .tooltip-footer {
    padding: 12px 16px;
    flex-direction: column;
    gap: 10px;
  }
  
  .tooltip-checkout-btn {
    width: 100%;
  }
}

/* 動畫效果 */
.cart-icon {
  @media (prefers-reduced-motion: no-preference) {
    .cart-icon-container {
      transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
    }
    
    .cart-symbol {
      transition: transform 0.3s cubic-bezier(0.4, 0, 0.2, 1);
    }
    
    .cart-badge {
      transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
    }
  }
}
</style>
