<template>
  <div class="add-to-cart">
    <!-- 數量選擇器 -->
    <div class="quantity-selector" v-if="showQuantitySelector">
      <button 
        @click="decreaseQuantity" 
        :disabled="quantity <= 1 || loading"
        class="quantity-btn quantity-decrease"
      >
        −
      </button>
      <input 
        v-model.number="quantity" 
        @change="validateQuantity"
        type="number" 
        min="1" 
        :max="maxQuantity"
        :disabled="loading"
        class="quantity-input"
      >
      <button 
        @click="increaseQuantity" 
        :disabled="quantity >= maxQuantity || loading"
        class="quantity-btn quantity-increase"
      >
        +
      </button>
    </div>

    <!-- 加入購物車按鈕 -->
    <button 
      @click="handleAddToCart" 
      :disabled="!canAddToCart || loading"
      :class="buttonClass"
    >
      <span v-if="loading" class="loading-spinner"></span>
      <span v-else class="button-icon">🛒</span>
      <span class="button-text">
        {{ buttonText }}
      </span>
    </button>

    <!-- 狀態提示 -->
    <div v-if="message" :class="messageClass">
      {{ message }}
    </div>
  </div>
</template>

<script>
import { ref, computed, watch } from 'vue'
import { useCartStore } from '@/stores/cart'
import { useAuthStore } from '@/stores/auth'

export default {
  name: 'AddToCartButton',
  props: {
    product: {
      type: Object,
      required: true
    },
    showQuantitySelector: {
      type: Boolean,
      default: true
    },
    variant: {
      type: String,
      default: 'primary', // primary, secondary, small
      validator: (value) => ['primary', 'secondary', 'small'].includes(value)
    },
    disabled: {
      type: Boolean,
      default: false
    }
  },
  emits: ['added-to-cart', 'error'],
  setup(props, { emit }) {
    const cartStore = useCartStore()
    const authStore = useAuthStore()

    const quantity = ref(1)
    const loading = ref(false)
    const message = ref('')
    const messageType = ref('')

    // 計算屬性
    const maxQuantity = computed(() => {
      return props.product.stock || 0
    })


    const canAddToCart = computed(() => {
      return !props.disabled && 
             props.product.is_active && 
             props.product.stock > 0 &&
             quantity.value > 0
    })

    const buttonClass = computed(() => {
      const baseClass = 'add-to-cart-btn'
      const variantClass = `btn-${props.variant}`
      const loadingClass = loading.value ? 'btn-loading' : ''
      const disabledClass = !canAddToCart.value ? 'btn-disabled' : ''
      
      return [baseClass, variantClass, loadingClass, disabledClass].join(' ')
    })

    const buttonText = computed(() => {
      if (loading.value) {
        return '處理中...'
      }
      if (!props.product.is_active) {
        return '商品已下架'
      }
      if (props.product.stock <= 0) {
        return '缺貨中'
      }
      return '加入購物車'
    })

    const messageClass = computed(() => {
      return ['message', `message-${messageType.value}`].join(' ')
    })

    // 方法
    const decreaseQuantity = () => {
      if (quantity.value > 1) {
        quantity.value--
      }
    }

    const increaseQuantity = () => {
      if (quantity.value < maxQuantity.value) {
        quantity.value++
      }
    }

    const validateQuantity = () => {
      if (quantity.value < 1) {
        quantity.value = 1
      } else if (quantity.value > maxQuantity.value) {
        quantity.value = maxQuantity.value
      }
    }

    const showMessage = (text, type = 'info') => {
      message.value = text
      messageType.value = type
      setTimeout(() => {
        message.value = ''
        messageType.value = ''
      }, 3000)
    }

    const handleAddToCart = async () => {
      if (!canAddToCart.value || loading.value) return

      // 檢查是否已登入
      if (!authStore.isAuthenticated) {
        showMessage('請先登入', 'warning')
        // 可以觸發登入流程
        return
      }

      loading.value = true
      message.value = ''

      try {
        await cartStore.addToCart(props.product.id, quantity.value)
        
        showMessage(`已將 ${quantity.value} 件「${props.product.name}」加入購物車！`, 'success')
        emit('added-to-cart', {
          product: props.product,
          quantity: quantity.value
        })
      } catch (err) {
        showMessage(err.message || '加入購物車失敗', 'error')
        emit('error', err)
      } finally {
        loading.value = false
      }
    }

    // 監聽商品變化，重置數量
    watch(() => props.product.id, () => {
      quantity.value = 1
      message.value = ''
    })

    return {
      quantity,
      loading,
      message,
      maxQuantity,
      canAddToCart,
      buttonClass,
      buttonText,
      messageClass,
      decreaseQuantity,
      increaseQuantity,
      validateQuantity,
      handleAddToCart
    }
  }
}
</script>

<style scoped>
.add-to-cart {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.quantity-selector {
  display: flex;
  align-items: center;
  border: 1px solid #e2e8f0;
  border-radius: 6px;
  overflow: hidden;
  width: fit-content;
}

.quantity-btn {
  width: 36px;
  height: 36px;
  border: none;
  background: #f7fafc;
  color: #4a5568;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 1.1rem;
  font-weight: 600;
  transition: background 0.2s;
}

.quantity-btn:hover:not(:disabled) {
  background: #edf2f7;
}

.quantity-btn:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.quantity-input {
  width: 60px;
  height: 36px;
  border: none;
  text-align: center;
  font-size: 1rem;
  font-weight: 500;
  color: #2d3748;
  background: white;
}

.quantity-input:focus {
  outline: none;
  background: #f7fafc;
}

.add-to-cart-btn {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 8px;
  padding: 12px 20px;
  border: none;
  border-radius: 6px;
  font-size: 1rem;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.2s;
  position: relative;
  overflow: hidden;
}

/* 按鈕變體 */
.btn-primary {
  background: #667eea;
  color: white;
}

.btn-primary:hover:not(:disabled) {
  background: #5a67d8;
  transform: translateY(-1px);
  box-shadow: 0 4px 12px rgba(102, 126, 234, 0.3);
}

.btn-secondary {
  background: #f7fafc;
  color: #4a5568;
  border: 1px solid #e2e8f0;
}

.btn-secondary:hover:not(:disabled) {
  background: #edf2f7;
  border-color: #cbd5e0;
}

.btn-small {
  padding: 8px 16px;
  font-size: 0.9rem;
}

/* 按鈕狀態 */
.btn-in-cart {
  background: #38a169;
  color: white;
}

.btn-in-cart:hover:not(:disabled) {
  background: #2f855a;
}

.btn-loading {
  cursor: not-allowed;
}

.btn-disabled {
  opacity: 0.5;
  cursor: not-allowed;
  transform: none !important;
  box-shadow: none !important;
}

.loading-spinner {
  width: 16px;
  height: 16px;
  border: 2px solid transparent;
  border-top: 2px solid currentColor;
  border-radius: 50%;
  animation: spin 1s linear infinite;
}

@keyframes spin {
  0% { transform: rotate(0deg); }
  100% { transform: rotate(360deg); }
}

.button-icon {
  font-size: 1.1rem;
}

.button-text {
  font-weight: 600;
}

.message {
  padding: 8px 12px;
  border-radius: 4px;
  font-size: 0.9rem;
  font-weight: 500;
  text-align: center;
}

.message-success {
  background: #f0fdf4;
  color: #059669;
  border: 1px solid #bbf7d0;
}

.message-error {
  background: #fef2f2;
  color: #dc2626;
  border: 1px solid #fecaca;
}

.message-warning {
  background: #fffbeb;
  color: #d97706;
  border: 1px solid #fed7aa;
}

.message-info {
  background: #eff6ff;
  color: #2563eb;
  border: 1px solid #bfdbfe;
}

/* 響應式設計 */
@media (max-width: 480px) {
  .add-to-cart-btn {
    padding: 10px 16px;
    font-size: 0.9rem;
  }
  
  .quantity-btn {
    width: 32px;
    height: 32px;
  }
  
  .quantity-input {
    width: 50px;
    height: 32px;
  }
}
</style>
