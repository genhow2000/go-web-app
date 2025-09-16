<template>
  <div class="product-card" @click="$emit('view', product.id)">
    <div class="product-image-placeholder">
      <div class="product-icon">📦</div>
    </div>
    <div class="product-info">
      <h3 class="product-name">{{ product.name }}</h3>
      <p class="product-description">{{ product.description }}</p>
      <div class="product-badges">
        <span v-if="product.is_featured" class="badge badge-featured">精選</span>
        <span v-if="product.is_on_sale" class="badge badge-sale">特價</span>
      </div>
      <div class="product-price">
        <span class="current-price">NT$ {{ product.price.toLocaleString() }}</span>
        <span v-if="product.original_price" class="original-price">
          NT$ {{ product.original_price.toLocaleString() }}
        </span>
      </div>
      <div class="product-actions">
        <button 
          class="btn-add-cart" 
          @click.stop="$emit('add-to-cart', product.id)"
        >
          加入購物車
        </button>
        <button 
          class="btn-favorite" 
          @click.stop="$emit('toggle-favorite', product.id)"
        >
          ❤️
        </button>
      </div>
    </div>
  </div>
</template>

<script>
export default {
  name: 'ProductCard',
  props: {
    product: {
      type: Object,
      required: true
    }
  },
  emits: ['view', 'add-to-cart', 'toggle-favorite']
}
</script>

<style scoped>
.product-card {
  background: white;
  border-radius: 15px;
  overflow: hidden;
  box-shadow: 0 5px 15px rgba(0,0,0,0.1);
  transition: all 0.3s;
  cursor: pointer;
}

.product-card:hover {
  transform: translateY(-5px);
  box-shadow: 0 15px 30px rgba(0,0,0,0.15);
}

.product-image-placeholder {
  width: 100%;
  height: 200px;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: 8px 8px 0 0;
}

.product-icon {
  font-size: 48px;
  color: white;
  text-shadow: 0 2px 4px rgba(0,0,0,0.3);
}

.product-info {
  padding: 1.5rem;
}

.product-name {
  font-size: 1.1rem;
  font-weight: 600;
  color: #2d3748;
  margin-bottom: 0.5rem;
}

.product-description {
  color: #718096;
  font-size: 0.9rem;
  margin-bottom: 1rem;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
}

.product-badges {
  display: flex;
  gap: 0.5rem;
  margin-bottom: 1rem;
}

.badge {
  padding: 0.25rem 0.75rem;
  border-radius: 20px;
  font-size: 0.8rem;
  font-weight: 600;
}

.badge-featured {
  background: #fef5e7;
  color: #744210;
}

.badge-sale {
  background: #fed7d7;
  color: #742a2a;
}

.product-price {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  margin-bottom: 1rem;
}

.current-price {
  font-size: 1.3rem;
  font-weight: bold;
  color: #e53e3e;
}

.original-price {
  font-size: 1rem;
  color: #a0aec0;
  text-decoration: line-through;
}

.product-actions {
  display: flex;
  gap: 0.5rem;
}

.btn-add-cart {
  flex: 1;
  background: #667eea;
  color: white;
  border: none;
  padding: 0.75rem;
  border-radius: 8px;
  font-weight: 600;
  cursor: pointer;
  transition: background 0.3s;
}

.btn-add-cart:hover {
  background: #5a6fd8;
}

.btn-favorite {
  background: #f7fafc;
  color: #4a5568;
  border: 1px solid #e2e8f0;
  padding: 0.75rem;
  border-radius: 8px;
  cursor: pointer;
  transition: all 0.3s;
}

.btn-favorite:hover {
  background: #edf2f7;
}
</style>
