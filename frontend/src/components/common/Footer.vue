<template>
  <footer class="footer">
    <div class="container">
      <div class="footer-content">
        <div class="footer-section">
          <h3>購物指南</h3>
          <ul>
            <li><a href="#">購物流程</a></li>
            <li><a href="#">配送說明</a></li>
            <li><a href="#">退換貨政策</a></li>
            <li><a href="#">常見問題</a></li>
          </ul>
        </div>
        <div class="footer-section">
          <h3>客戶服務</h3>
          <ul>
            <li><a href="#">聯絡我們</a></li>
            <li><a href="#">在線客服</a></li>
            <li><a href="#">意見反饋</a></li>
            <li><a href="#">售後服務</a></li>
          </ul>
        </div>
        <div class="footer-section">
          <h3>關於我們</h3>
          <ul>
            <li><a href="#">公司簡介</a></li>
            <li><a href="#">招聘信息</a></li>
            <li><a href="#">合作夥伴</a></li>
            <li><a href="#">隱私政策</a></li>
          </ul>
        </div>
        <div class="footer-section">
          <h3>關注我們</h3>
          <ul>
            <li><a href="#">Facebook</a></li>
            <li><a href="#">Instagram</a></li>
            <li><a href="#">Twitter</a></li>
            <li><a href="#">YouTube</a></li>
          </ul>
        </div>
      </div>
      <div class="footer-bottom">
        <p>&copy; 2025 阿和商城. 版權所有 | 技術展示平台</p>
      </div>
    </div>
  </footer>
</template>

<script>
export default {
  name: 'Footer'
}
</script>

<style scoped>
.footer {
  background: #2d3748;
  color: white;
  padding: 3rem 0 1rem;
}

.container {
  max-width: 1200px;
  margin: 0 auto;
  padding: 0 20px;
}

.footer-content {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(250px, 1fr));
  gap: 2rem;
  margin-bottom: 2rem;
}

.footer-section h3 {
  margin-bottom: 1rem;
  color: #e2e8f0;
}

.footer-section ul {
  list-style: none;
}

.footer-section ul li {
  margin-bottom: 0.5rem;
}

.footer-section ul li a {
  color: #a0aec0;
  text-decoration: none;
  transition: color 0.3s;
}

.footer-section ul li a:hover {
  color: white;
}

.footer-bottom {
  text-align: center;
  padding-top: 2rem;
  border-top: 1px solid #4a5568;
  color: #a0aec0;
}
</style>
