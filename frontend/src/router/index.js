import { createRouter, createWebHistory } from 'vue-router'
import { useAuthStore } from '@/stores/auth'

// 頁面組件
import HomePage from '@/views/HomePage.vue'
import TechShowcase from '@/views/TechShowcase.vue'
import CustomerLogin from '@/views/auth/CustomerLogin.vue'
import CustomerRegister from '@/views/auth/CustomerRegister.vue'
import MerchantLogin from '@/views/auth/MerchantLogin.vue'
import AdminLogin from '@/views/auth/AdminLogin.vue'
import MerchantRegister from '@/views/auth/MerchantRegister.vue'
import CustomerDashboard from '@/views/dashboard/CustomerDashboard.vue'
import MerchantDashboard from '@/views/dashboard/MerchantDashboard.vue'
import AdminDashboard from '@/views/dashboard/AdminDashboard.vue'
import MerchantProducts from '@/views/merchant/MerchantProducts.vue'
import CreateProduct from '@/views/merchant/CreateProduct.vue'
import EditProduct from '@/views/merchant/EditProduct.vue'
import CategoryPage from '@/views/CategoryPage.vue'
import ProductDetail from '@/views/ProductDetail.vue'
import CartPage from '@/views/CartPage.vue'

const routes = [
  {
    path: '/',
    name: 'Home',
    component: HomePage
  },
  {
    path: '/tech-showcase',
    name: 'TechShowcase',
    component: TechShowcase
  },
  {
    path: '/category/:category',
    name: 'CategoryPage',
    component: CategoryPage
  },
  {
    path: '/product/:id',
    name: 'ProductDetail',
    component: ProductDetail
  },
  {
    path: '/cart',
    name: 'CartPage',
    component: CartPage,
    meta: { requiresAuth: true, role: 'customer' }
  },
  // 登入頁面
  {
    path: '/customer/login',
    name: 'CustomerLogin',
    component: CustomerLogin
  },
  {
    path: '/merchant/login',
    name: 'MerchantLogin',
    component: MerchantLogin
  },
  {
    path: '/admin/login',
    name: 'AdminLogin',
    component: AdminLogin
  },
  // 註冊頁面
  {
    path: '/register',
    name: 'CustomerRegister',
    component: CustomerRegister
  },
  {
    path: '/merchant/register',
    name: 'MerchantRegister',
    component: MerchantRegister
  },
  // 受保護的儀表板頁面
  {
    path: '/customer/dashboard',
    name: 'CustomerDashboard',
    component: CustomerDashboard,
    meta: { requiresAuth: true, role: 'customer' }
  },
  {
    path: '/merchant/dashboard',
    name: 'MerchantDashboard',
    component: MerchantDashboard,
    meta: { requiresAuth: true, role: 'merchant' }
  },
  {
    path: '/admin/dashboard',
    name: 'AdminDashboard',
    component: AdminDashboard,
    meta: { requiresAuth: true, role: 'admin' }
  },
  // 商戶商品管理
  {
    path: '/merchant/products',
    name: 'MerchantProducts',
    component: MerchantProducts,
    meta: { requiresAuth: true, role: 'merchant' }
  },
  {
    path: '/merchant/products/create',
    name: 'CreateProduct',
    component: CreateProduct,
    meta: { requiresAuth: true, role: 'merchant' }
  },
  {
    path: '/merchant/products/:id/edit',
    name: 'EditProduct',
    component: EditProduct,
    meta: { requiresAuth: true, role: 'merchant' }
  }
]

const router = createRouter({
  history: createWebHistory(),
  routes
})

// 路由守衛
router.beforeEach(async (to, from, next) => {
  const authStore = useAuthStore()
  
  // 確保認證狀態已初始化
  if (!authStore.user && !authStore.token) {
    try {
      await authStore.initAuth()
    } catch (error) {
      console.error('認證初始化失敗', error)
    }
  }
  
  if (to.meta.requiresAuth) {
    if (!authStore.isAuthenticated) {
      // 根據角色重定向到對應的登入頁面
      const role = to.meta.role
      if (role === 'customer') {
        next('/customer/login')
      } else if (role === 'merchant') {
        next('/merchant/login')
      } else if (role === 'admin') {
        next('/admin/login')
      } else {
        next('/')
      }
      return
    }
    
    // 檢查角色權限
    if (to.meta.role && authStore.user?.role !== to.meta.role) {
      next('/')
      return
    }
  }
  
  next()
})

export default router
