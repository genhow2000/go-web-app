<template>
  <div class="cart-summary">
    <div class="summary-header">
      <h3>購物車摘要</h3>
    </div>

    <div class="summary-content">
      <!-- 商品統計 -->
      <div class="summary-stats">
        <div class="stat-item">
          <span class="stat-label">商品數量</span>
          <span class="stat-value">{{ itemCount }} 件</span>
        </div>
        <div class="stat-item">
          <span class="stat-label">商品種類</span>
          <span class="stat-value">{{ items.length }} 種</span>
        </div>
      </div>

      <!-- 價格明細 -->
      <div class="price-breakdown">
        <div class="price-item">
          <span class="price-label">商品總價</span>
          <span class="price-value">NT$ {{ totalPrice.toLocaleString() }}</span>
        </div>
        
        <!-- 運費 -->
        <div class="price-item">
          <span class="price-label">運費</span>
          <span class="price-value">
            <span v-if="shippingFee === 0" class="free-shipping">免費</span>
            <span v-else>NT$ {{ shippingFee.toLocaleString() }}</span>
          </span>
        </div>

        <!-- 折扣 -->
        <div class="price-item discount" v-if="discount > 0">
          <span class="price-label">折扣</span>
          <span class="price-value discount-value">-NT$ {{ discount.toLocaleString() }}</span>
        </div>

        <!-- 分隔線 -->
        <div class="price-divider"></div>

        <!-- 總計 -->
        <div class="price-item total">
          <span class="price-label">總計</span>
          <span class="price-value total-value">NT$ {{ finalTotal.toLocaleString() }}</span>
        </div>
      </div>

      <!-- 優惠券輸入 -->
      <div class="coupon-section">
        <div class="coupon-input-group">
          <input 
            v-model="couponCode" 
            type="text" 
            placeholder="輸入優惠券代碼"
            class="coupon-input"
            @keyup.enter="applyCoupon"
          >
          <button 
            @click="applyCoupon" 
            :disabled="!couponCode.trim() || applyingCoupon"
            class="coupon-btn"
          >
            {{ applyingCoupon ? '使用中...' : '使用' }}
          </button>
        </div>
        
        <!-- 已使用的優惠券 -->
        <div v-if="appliedCoupon" class="applied-coupon">
          <span class="coupon-name">{{ appliedCoupon.name }}</span>
          <button @click="removeCoupon" class="remove-coupon-btn">×</button>
        </div>
      </div>

      <!-- 操作按鈕 -->
      <div class="summary-actions">
        <button 
          @click="handleCheckout" 
          :disabled="!canCheckout"
          class="checkout-btn"
        >
          立即結算
        </button>
        
        <button 
          @click="handleClearCart" 
          class="clear-btn"
        >
          清空購物車
        </button>
      </div>

      <!-- 安全提示 -->
      <div class="security-note">
        <div class="security-icon">🔒</div>
        <p>您的購物車信息已加密保護</p>
      </div>
    </div>
  </div>
</template>

<script>
import { ref, computed } from 'vue'

export default {
  name: 'CartSummary',
  props: {
    items: {
      type: Array,
      default: () => []
    },
    totalPrice: {
      type: Number,
      default: 0
    },
    itemCount: {
      type: Number,
      default: 0
    }
  },
  emits: ['clear-cart', 'checkout'],
  setup(props, { emit }) {
    const couponCode = ref('')
    const appliedCoupon = ref(null)
    const applyingCoupon = ref(false)

    // 運費計算
    const shippingFee = computed(() => {
      // 免費運費門檻
      const freeShippingThreshold = 1000
      if (props.totalPrice >= freeShippingThreshold) {
        return 0
      }
      // 固定運費
      return 100
    })

    // 折扣計算
    const discount = computed(() => {
      if (appliedCoupon.value) {
        return appliedCoupon.value.discount
      }
      return 0
    })

    // 最終總計
    const finalTotal = computed(() => {
      return Math.max(0, props.totalPrice + shippingFee.value - discount.value)
    })

    // 是否可以結算
    const canCheckout = computed(() => {
      return props.itemCount > 0 && props.totalPrice > 0
    })

    // 應用優惠券
    const applyCoupon = async () => {
      if (!couponCode.value.trim()) return

      applyingCoupon.value = true
      try {
        // 這裡應該調用API驗證優惠券
        // 暫時模擬API調用
        await new Promise(resolve => setTimeout(resolve, 1000))
        
        // 模擬優惠券數據
        const mockCoupons = {
          'WELCOME10': { name: '新用戶優惠券', discount: 100 },
          'SAVE50': { name: '節省50元', discount: 50 },
          'FREESHIP': { name: '免運費', discount: shippingFee.value }
        }

        const coupon = mockCoupons[couponCode.value.toUpperCase()]
        if (coupon) {
          appliedCoupon.value = coupon
          couponCode.value = ''
        } else {
          alert('優惠券代碼無效')
        }
      } catch (err) {
        alert('優惠券驗證失敗')
      } finally {
        applyingCoupon.value = false
      }
    }

    // 移除優惠券
    const removeCoupon = () => {
      appliedCoupon.value = null
    }

    // 結算
    const handleCheckout = () => {
      emit('checkout')
    }

    // 清空購物車
    const handleClearCart = () => {
      emit('clear-cart')
    }

    return {
      couponCode,
      appliedCoupon,
      applyingCoupon,
      shippingFee,
      discount,
      finalTotal,
      canCheckout,
      applyCoupon,
      removeCoupon,
      handleCheckout,
      handleClearCart
    }
  }
}
</script>

<style scoped>
.cart-summary {
  background: white;
  border-radius: 10px;
  box-shadow: 0 2px 10px rgba(0,0,0,0.1);
  overflow: hidden;
}

.summary-header {
  background: #667eea;
  color: white;
  padding: 20px;
  text-align: center;
}

.summary-header h3 {
  margin: 0;
  font-size: 1.3rem;
  font-weight: 600;
}

.summary-content {
  padding: 25px;
}

.summary-stats {
  margin-bottom: 25px;
}

.stat-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 8px 0;
  border-bottom: 1px solid #f1f5f9;
}

.stat-item:last-child {
  border-bottom: none;
}

.stat-label {
  color: #64748b;
  font-size: 0.9rem;
}

.stat-value {
  color: #1e293b;
  font-weight: 500;
}

.price-breakdown {
  margin-bottom: 25px;
}

.price-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 10px 0;
}

.price-label {
  color: #64748b;
  font-size: 0.95rem;
}

.price-value {
  color: #1e293b;
  font-weight: 500;
}

.free-shipping {
  color: #059669;
  font-weight: 600;
}

.price-item.discount .price-label {
  color: #059669;
}

.discount-value {
  color: #059669;
  font-weight: 600;
}

.price-divider {
  height: 1px;
  background: #e2e8f0;
  margin: 15px 0;
}

.price-item.total {
  padding: 15px 0;
  border-top: 2px solid #e2e8f0;
}

.price-item.total .price-label {
  font-size: 1.1rem;
  font-weight: 600;
  color: #1e293b;
}

.total-value {
  font-size: 1.3rem;
  font-weight: 700;
  color: #667eea;
}

.coupon-section {
  margin-bottom: 25px;
}

.coupon-input-group {
  display: flex;
  gap: 10px;
  margin-bottom: 15px;
}

.coupon-input {
  flex: 1;
  padding: 10px 12px;
  border: 1px solid #d1d5db;
  border-radius: 6px;
  font-size: 0.9rem;
  transition: border-color 0.2s;
}

.coupon-input:focus {
  outline: none;
  border-color: #667eea;
  box-shadow: 0 0 0 3px rgba(102, 126, 234, 0.1);
}

.coupon-btn {
  padding: 10px 16px;
  background: #667eea;
  color: white;
  border: none;
  border-radius: 6px;
  font-size: 0.9rem;
  font-weight: 500;
  cursor: pointer;
  transition: background 0.2s;
  white-space: nowrap;
}

.coupon-btn:hover:not(:disabled) {
  background: #5a67d8;
}

.coupon-btn:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.applied-coupon {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 10px 12px;
  background: #f0fdf4;
  border: 1px solid #bbf7d0;
  border-radius: 6px;
}

.coupon-name {
  color: #059669;
  font-size: 0.9rem;
  font-weight: 500;
}

.remove-coupon-btn {
  background: none;
  border: none;
  color: #059669;
  font-size: 1.2rem;
  cursor: pointer;
  padding: 0;
  width: 20px;
  height: 20px;
  display: flex;
  align-items: center;
  justify-content: center;
}

.remove-coupon-btn:hover {
  background: #dcfce7;
  border-radius: 50%;
}

.summary-actions {
  display: flex;
  flex-direction: column;
  gap: 12px;
  margin-bottom: 20px;
}

.checkout-btn {
  width: 100%;
  padding: 15px;
  background: #667eea;
  color: white;
  border: none;
  border-radius: 8px;
  font-size: 1.1rem;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.2s;
}

.checkout-btn:hover:not(:disabled) {
  background: #5a67d8;
  transform: translateY(-1px);
  box-shadow: 0 4px 12px rgba(102, 126, 234, 0.3);
}

.checkout-btn:disabled {
  opacity: 0.5;
  cursor: not-allowed;
  transform: none;
  box-shadow: none;
}

.clear-btn {
  width: 100%;
  padding: 12px;
  background: #f8fafc;
  color: #64748b;
  border: 1px solid #e2e8f0;
  border-radius: 6px;
  font-size: 0.95rem;
  cursor: pointer;
  transition: all 0.2s;
}

.clear-btn:hover {
  background: #f1f5f9;
  color: #475569;
}

.security-note {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 12px;
  background: #f8fafc;
  border-radius: 6px;
  text-align: center;
}

.security-icon {
  font-size: 1rem;
}

.security-note p {
  margin: 0;
  color: #64748b;
  font-size: 0.85rem;
}

/* 響應式設計 */
@media (max-width: 768px) {
  .summary-content {
    padding: 20px;
  }
  
  .coupon-input-group {
    flex-direction: column;
  }
  
  .coupon-btn {
    width: 100%;
  }
}
</style>
