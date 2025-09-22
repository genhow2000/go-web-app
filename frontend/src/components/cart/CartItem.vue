<template>
  <div class="cart-item">
    <!-- 商品圖片 -->
    <div class="item-image">
      <img 
        v-if="item.product?.image_url" 
        :src="item.product.image_url" 
        :alt="item.product.name"
        @error="handleImageError"
      >
      <div v-else class="image-placeholder">
        <div class="placeholder-icon">📦</div>
      </div>
    </div>

    <!-- 商品信息 -->
    <div class="item-info">
      <h3 class="item-name">{{ item.product?.name || '商品名稱' }}</h3>
      <p class="item-description">{{ item.product?.description || '商品描述' }}</p>
      
      <!-- 商品標籤 -->
      <div class="item-badges" v-if="item.product">
        <span v-if="item.product.is_featured" class="badge badge-featured">精選</span>
        <span v-if="item.product.is_on_sale" class="badge badge-sale">特價</span>
      </div>

      <!-- 價格信息 -->
      <div class="item-price">
        <span class="current-price">NT$ {{ item.price.toLocaleString() }}</span>
        <span v-if="item.product?.original_price && item.product.original_price > item.price" 
              class="original-price">
          NT$ {{ item.product.original_price.toLocaleString() }}
        </span>
      </div>

      <!-- 庫存狀態 -->
      <div class="stock-status" v-if="item.product">
        <span v-if="item.product.stock <= 0" class="out-of-stock">缺貨</span>
        <span v-else-if="item.product.stock < 10" class="low-stock">庫存不足</span>
        <span v-else class="in-stock">有庫存</span>
      </div>
    </div>

    <!-- 數量控制 -->
    <div class="item-controls">
      <div class="quantity-control">
        <button 
          @click="decreaseQuantity" 
          :disabled="quantity <= 1"
          class="quantity-btn quantity-decrease"
        >
          −
        </button>
        <input 
          v-model.number="quantity" 
          @change="updateQuantity"
          @blur="validateQuantity"
          type="number" 
          min="1" 
          :max="maxQuantity"
          class="quantity-input"
        >
        <button 
          @click="increaseQuantity" 
          :disabled="quantity >= maxQuantity"
          class="quantity-btn quantity-increase"
        >
          +
        </button>
      </div>

      <!-- 小計 -->
      <div class="item-subtotal">
        NT$ {{ subtotal.toLocaleString() }}
      </div>

      <!-- 移除按鈕 -->
      <button 
        @click="removeItem" 
        class="remove-btn"
        title="移除商品"
      >
        🗑️
      </button>
    </div>
  </div>
</template>

<script>
import { ref, computed, watch } from 'vue'

export default {
  name: 'CartItem',
  props: {
    item: {
      type: Object,
      required: true
    }
  },
  emits: ['update-quantity', 'remove-item'],
  setup(props, { emit }) {
    const quantity = ref(props.item.quantity)
    const isUpdating = ref(false)

    // 計算最大數量（庫存限制）
    const maxQuantity = computed(() => {
      return props.item.product?.stock || 999
    })

    // 計算小計
    const subtotal = computed(() => {
      return props.item.price * quantity.value
    })

    // 減少數量
    const decreaseQuantity = () => {
      if (quantity.value > 1) {
        quantity.value--
        updateQuantity()
      }
    }

    // 增加數量
    const increaseQuantity = () => {
      if (quantity.value < maxQuantity.value) {
        quantity.value++
        updateQuantity()
      }
    }

    // 更新數量
    const updateQuantity = async () => {
      if (isUpdating.value) return
      
      isUpdating.value = true
      try {
        await emit('update-quantity', props.item.product_id, quantity.value)
      } catch (err) {
        // 如果更新失敗，恢復原數量
        quantity.value = props.item.quantity
      } finally {
        isUpdating.value = false
      }
    }

    // 驗證數量
    const validateQuantity = () => {
      if (quantity.value < 1) {
        quantity.value = 1
      } else if (quantity.value > maxQuantity.value) {
        quantity.value = maxQuantity.value
      }
      
      if (quantity.value !== props.item.quantity) {
        updateQuantity()
      }
    }

    // 移除商品
    const removeItem = () => {
      if (confirm('確定要移除這個商品嗎？')) {
        emit('remove-item', props.item.product_id)
      }
    }

    // 圖片錯誤處理
    const handleImageError = (event) => {
      event.target.style.display = 'none'
      const placeholder = event.target.nextElementSibling
      if (placeholder) {
        placeholder.style.display = 'flex'
      }
    }

    // 監聽 props 變化
    watch(() => props.item.quantity, (newQuantity) => {
      quantity.value = newQuantity
    })

    return {
      quantity,
      maxQuantity,
      subtotal,
      decreaseQuantity,
      increaseQuantity,
      updateQuantity,
      validateQuantity,
      removeItem,
      handleImageError
    }
  }
}
</script>

<style scoped>
.cart-item {
  display: flex;
  align-items: center;
  gap: 20px;
  padding: 20px;
  background: white;
  border-radius: 10px;
  box-shadow: 0 2px 8px rgba(0,0,0,0.1);
  transition: box-shadow 0.3s;
}

.cart-item:hover {
  box-shadow: 0 4px 12px rgba(0,0,0,0.15);
}

.item-image {
  flex-shrink: 0;
  width: 120px;
  height: 120px;
  border-radius: 8px;
  overflow: hidden;
  background: #f7fafc;
}

.item-image img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.image-placeholder {
  width: 100%;
  height: 100%;
  display: flex;
  align-items: center;
  justify-content: center;
  background: #f7fafc;
}

.placeholder-icon {
  font-size: 2rem;
  color: #a0aec0;
}

.item-info {
  flex: 1;
  min-width: 0;
}

.item-name {
  font-size: 1.2rem;
  font-weight: 600;
  color: #2d3748;
  margin-bottom: 8px;
  line-height: 1.4;
}

.item-description {
  color: #718096;
  font-size: 0.9rem;
  margin-bottom: 12px;
  line-height: 1.4;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
}

.item-badges {
  margin-bottom: 12px;
}

.badge {
  display: inline-block;
  padding: 4px 8px;
  border-radius: 4px;
  font-size: 0.75rem;
  font-weight: 500;
  margin-right: 8px;
}

.badge-featured {
  background: #fef5e7;
  color: #d69e2e;
}

.badge-sale {
  background: #fed7d7;
  color: #e53e3e;
}

.item-price {
  margin-bottom: 8px;
}

.current-price {
  font-size: 1.1rem;
  font-weight: 600;
  color: #2d3748;
}

.original-price {
  font-size: 0.9rem;
  color: #a0aec0;
  text-decoration: line-through;
  margin-left: 8px;
}

.stock-status {
  font-size: 0.8rem;
  font-weight: 500;
}

.in-stock {
  color: #38a169;
}

.low-stock {
  color: #d69e2e;
}

.out-of-stock {
  color: #e53e3e;
}

.item-controls {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 15px;
  flex-shrink: 0;
}

.quantity-control {
  display: flex;
  align-items: center;
  border: 1px solid #e2e8f0;
  border-radius: 6px;
  overflow: hidden;
}

.quantity-btn {
  width: 32px;
  height: 32px;
  border: none;
  background: #f7fafc;
  color: #4a5568;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 1.2rem;
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
  height: 32px;
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

.item-subtotal {
  font-size: 1.1rem;
  font-weight: 600;
  color: #2d3748;
  text-align: center;
}

.remove-btn {
  width: 36px;
  height: 36px;
  border: none;
  background: #fed7d7;
  color: #e53e3e;
  border-radius: 50%;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 1rem;
  transition: all 0.2s;
}

.remove-btn:hover {
  background: #feb2b2;
  transform: scale(1.1);
}

/* 響應式設計 */
@media (max-width: 768px) {
  .cart-item {
    flex-direction: column;
    text-align: center;
    gap: 15px;
  }
  
  .item-image {
    width: 100px;
    height: 100px;
  }
  
  .item-controls {
    flex-direction: row;
    justify-content: space-between;
    width: 100%;
  }
  
  .quantity-control {
    order: 1;
  }
  
  .item-subtotal {
    order: 2;
  }
  
  .remove-btn {
    order: 3;
  }
}

@media (max-width: 480px) {
  .cart-item {
    padding: 15px;
  }
  
  .item-name {
    font-size: 1.1rem;
  }
  
  .item-description {
    font-size: 0.85rem;
  }
}
</style>
