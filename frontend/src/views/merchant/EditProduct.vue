<template>
  <div class="edit-product">
    <Header />
    
    <div class="container">
      <div class="page-header">
        <h1>編輯商品</h1>
        <router-link to="/merchant/products" class="btn btn-outline">
          <i class="fas fa-arrow-left"></i>
          返回商品列表
        </router-link>
      </div>

      <div v-if="loading" class="loading">
        <div class="spinner"></div>
        <p>載入商品信息中...</p>
      </div>

      <div v-else-if="!product" class="error-state">
        <div class="error-icon">❌</div>
        <h3>商品不存在</h3>
        <p>找不到指定的商品，請檢查商品ID是否正確。</p>
        <router-link to="/merchant/products" class="btn btn-primary">
          返回商品列表
        </router-link>
      </div>

      <div v-else class="form-container">
        <form @submit.prevent="handleSubmit">
          <div class="form-row">
            <div class="form-group">
              <label for="name">商品名稱 *</label>
              <input 
                v-model="form.name"
                type="text" 
                id="name" 
                required
                placeholder="請輸入商品名稱"
              >
            </div>
            <div class="form-group">
              <label for="price">價格 *</label>
              <input 
                v-model.number="form.price"
                type="number" 
                id="price" 
                required
                min="0"
                step="0.01"
                placeholder="請輸入價格"
              >
            </div>
          </div>

          <div class="form-group">
            <label for="description">商品描述 *</label>
            <textarea 
              v-model="form.description"
              id="description" 
              required
              rows="4"
              placeholder="請輸入商品描述"
            ></textarea>
          </div>

          <div class="form-row">
            <div class="form-group">
              <label for="category">商品分類</label>
              <select v-model="form.category" id="category">
                <option value="">請選擇分類</option>
                <option value="電子產品">電子產品</option>
                <option value="服飾">服飾</option>
                <option value="家居">家居</option>
                <option value="美妝">美妝</option>
                <option value="運動">運動</option>
                <option value="食品">食品</option>
                <option value="圖書">圖書</option>
                <option value="其他">其他</option>
              </select>
            </div>
            <div class="form-group">
              <label for="stock">庫存數量</label>
              <input 
                v-model.number="form.stock"
                type="number" 
                id="stock" 
                min="0"
                placeholder="請輸入庫存數量"
              >
            </div>
          </div>

          <div class="form-row">
            <div class="form-group">
              <label for="original_price">原價</label>
              <input 
                v-model.number="form.original_price"
                type="number" 
                id="original_price" 
                min="0"
                step="0.01"
                placeholder="請輸入原價（用於顯示折扣）"
              >
            </div>
            <div class="form-group">
              <label for="sku">商品編號</label>
              <input 
                v-model="form.sku"
                type="text" 
                id="sku" 
                placeholder="請輸入商品編號"
              >
            </div>
          </div>

          <div class="form-checkboxes">
            <label class="checkbox-label">
              <input 
                v-model="form.is_featured"
                type="checkbox"
              >
              <span class="checkmark"></span>
              設為精選商品
            </label>
            <label class="checkbox-label">
              <input 
                v-model="form.is_on_sale"
                type="checkbox"
              >
              <span class="checkmark"></span>
              設為特價商品
            </label>
            <label class="checkbox-label">
              <input 
                v-model="form.is_active"
                type="checkbox"
              >
              <span class="checkmark"></span>
              上架狀態
            </label>
          </div>

          <div class="form-actions">
            <button type="button" @click="loadProduct" class="btn btn-outline">
              重置
            </button>
            <button type="submit" class="btn btn-primary" :disabled="saving">
              {{ saving ? '保存中...' : '保存修改' }}
            </button>
          </div>
        </form>
      </div>
    </div>
  </div>
</template>

<script>
import { ref, reactive, onMounted } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import Header from '@/components/common/Header.vue'
import api from '@/services/api'

export default {
  name: 'EditProduct',
  components: {
    Header
  },
  setup() {
    const router = useRouter()
    const route = useRoute()
    const loading = ref(true)
    const saving = ref(false)
    const product = ref(null)

    const form = reactive({
      name: '',
      description: '',
      price: 0,
      original_price: 0,
      category: '',
      stock: 0,
      sku: '',
      is_featured: false,
      is_on_sale: false,
      is_active: true
    })

    const loadProduct = async () => {
      loading.value = true
      try {
        const productId = route.params.id
        const response = await api.get(`/merchant/api/products/${productId}`)
        
        if (response.data) {
          product.value = response.data
          // 填充表單
          Object.assign(form, {
            name: product.value.name || '',
            description: product.value.description || '',
            price: product.value.price || 0,
            original_price: product.value.original_price || 0,
            category: product.value.category || '',
            stock: product.value.stock || 0,
            sku: product.value.sku || '',
            is_featured: product.value.is_featured || false,
            is_on_sale: product.value.is_on_sale || false,
            is_active: product.value.is_active !== false
          })
        }
      } catch (error) {
        console.error('載入商品失敗:', error)
        // 使用模擬數據
        product.value = {
          id: route.params.id,
          name: '測試商品',
          description: '這是一個測試商品',
          price: 2999,
          original_price: 3999,
          category: '電子產品',
          stock: 100,
          sku: 'TEST001',
          is_featured: true,
          is_on_sale: true,
          is_active: true
        }
        Object.assign(form, product.value)
      } finally {
        loading.value = false
      }
    }

    const handleSubmit = async () => {
      saving.value = true
      try {
        const productId = route.params.id
        const response = await api.put(`/merchant/api/products/${productId}`, form)
        
        if (response.data.success) {
          alert('商品更新成功！')
          router.push('/merchant/products')
        } else {
          alert(response.data.error || '更新失敗，請重試')
        }
      } catch (error) {
        console.error('更新商品失敗:', error)
        alert('更新失敗，請重試')
      } finally {
        saving.value = false
      }
    }

    onMounted(() => {
      loadProduct()
    })

    return {
      loading,
      saving,
      product,
      form,
      loadProduct,
      handleSubmit
    }
  }
}
</script>

<style scoped>
.edit-product {
  min-height: 100vh;
  background: #f8f9fa;
}

.container {
  max-width: 800px;
  margin: 0 auto;
  padding: 20px;
}

.page-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 30px;
  background: white;
  padding: 20px;
  border-radius: 10px;
  box-shadow: 0 2px 10px rgba(0,0,0,0.1);
}

.page-header h1 {
  margin: 0;
  color: #2d3748;
}

.loading {
  text-align: center;
  padding: 60px 20px;
  color: #718096;
}

.spinner {
  display: inline-block;
  width: 40px;
  height: 40px;
  border: 4px solid #f3f3f3;
  border-top: 4px solid #667eea;
  border-radius: 50%;
  animation: spin 1s linear infinite;
  margin-bottom: 20px;
}

@keyframes spin {
  0% { transform: rotate(0deg); }
  100% { transform: rotate(360deg); }
}

.error-state {
  text-align: center;
  padding: 60px 20px;
  color: #718096;
}

.error-icon {
  font-size: 4rem;
  margin-bottom: 20px;
}

.error-state h3 {
  margin-bottom: 10px;
  color: #4a5568;
}

.error-state p {
  margin-bottom: 30px;
}

.form-container {
  background: white;
  padding: 30px;
  border-radius: 10px;
  box-shadow: 0 2px 10px rgba(0,0,0,0.1);
}

.form-row {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 20px;
  margin-bottom: 20px;
}

.form-group {
  margin-bottom: 20px;
}

.form-group label {
  display: block;
  margin-bottom: 8px;
  font-weight: 500;
  color: #4a5568;
}

.form-group input,
.form-group textarea,
.form-group select {
  width: 100%;
  padding: 12px;
  border: 1px solid #e2e8f0;
  border-radius: 6px;
  font-size: 14px;
  transition: border-color 0.3s;
}

.form-group input:focus,
.form-group textarea:focus,
.form-group select:focus {
  outline: none;
  border-color: #667eea;
}

.form-group textarea {
  resize: vertical;
  min-height: 100px;
}

.form-checkboxes {
  display: flex;
  gap: 30px;
  margin-bottom: 30px;
  flex-wrap: wrap;
}

.checkbox-label {
  display: flex;
  align-items: center;
  cursor: pointer;
  font-weight: 500;
  color: #4a5568;
}

.checkbox-label input[type="checkbox"] {
  display: none;
}

.checkmark {
  width: 20px;
  height: 20px;
  border: 2px solid #e2e8f0;
  border-radius: 4px;
  margin-right: 8px;
  position: relative;
  transition: all 0.3s;
}

.checkbox-label input[type="checkbox"]:checked + .checkmark {
  background: #667eea;
  border-color: #667eea;
}

.checkbox-label input[type="checkbox"]:checked + .checkmark::after {
  content: '✓';
  position: absolute;
  top: 50%;
  left: 50%;
  transform: translate(-50%, -50%);
  color: white;
  font-size: 12px;
  font-weight: bold;
}

.form-actions {
  display: flex;
  gap: 15px;
  justify-content: flex-end;
  padding-top: 20px;
  border-top: 1px solid #e2e8f0;
}

.btn {
  padding: 12px 24px;
  border-radius: 6px;
  text-decoration: none;
  font-size: 14px;
  font-weight: 500;
  border: none;
  cursor: pointer;
  transition: all 0.3s;
  display: inline-flex;
  align-items: center;
  gap: 8px;
}

.btn-primary {
  background: #667eea;
  color: white;
}

.btn-primary:hover:not(:disabled) {
  background: #5a6fd8;
}

.btn-primary:disabled {
  background: #a0aec0;
  cursor: not-allowed;
}

.btn-outline {
  background: transparent;
  color: #667eea;
  border: 1px solid #667eea;
}

.btn-outline:hover {
  background: #667eea;
  color: white;
}

@media (max-width: 768px) {
  .page-header {
    flex-direction: column;
    gap: 15px;
    align-items: flex-start;
  }
  
  .form-row {
    grid-template-columns: 1fr;
  }
  
  .form-checkboxes {
    flex-direction: column;
    gap: 15px;
  }
  
  .form-actions {
    flex-direction: column;
  }
}
</style>
