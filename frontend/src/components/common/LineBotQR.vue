<template>
  <div class="line-bot-qr-modal" v-if="show" @click="closeModal">
    <div class="modal-content" @click.stop>
      <div class="modal-header">
        <h3>加入我們的 LINE 機器人</h3>
        <button class="close-btn" @click="closeModal">×</button>
      </div>
      
      <div class="modal-body">
        <div class="qr-section">
          <div class="qr-code">
            <div ref="qrCodeRef" class="qr-container"></div>
          </div>
          <p class="qr-description">掃描 QR Code 加入 LINE 機器人</p>
        </div>
        
        <div class="alternative-methods">
          <h4>其他加入方式：</h4>
          <div class="method-item">
            <span class="method-icon">🔍</span>
            <span class="method-text">在 LINE 中搜尋：<strong>@351thdpd</strong></span>
            <button class="copy-btn" @click="copyId">複製</button>
          </div>
          <div class="method-item">
            <span class="method-icon">🔗</span>
            <span class="method-text">點擊連結直接加入</span>
            <button class="link-btn" @click="openLineApp">開啟 LINE</button>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import { ref, onMounted, nextTick } from 'vue'

export default {
  name: 'LineBotQR',
  props: {
    show: {
      type: Boolean,
      default: false
    }
  },
  emits: ['close'],
  setup(props, { emit }) {
    const qrCodeRef = ref(null)
    const lineBotId = '@351thdpd'
    // 使用多種格式嘗試
    const lineBotUrls = [
      'https://line.me/R/ti/p/351thdpd',
      'https://line.me/R/ti/p/@351thdpd',
      'https://line.me/R/ti/p/line://ti/p/351thdpd'
    ]
    const lineBotUrl = lineBotUrls[0] // 使用第一個作為預設

    const closeModal = () => {
      emit('close')
    }

    const generateQRCode = () => {
      if (!qrCodeRef.value) return

      // 清空容器
      qrCodeRef.value.innerHTML = ''

      // 使用 LINE 官方 QR Code 圖片
      const qrImg = document.createElement('img')
      qrImg.src = 'https://qr-official.line.me/gs/M_351thdpd_GW.png?oat_content=qr'
      qrImg.alt = 'LINE Bot QR Code'
      qrImg.style.width = '200px'
      qrImg.style.height = '200px'
      qrImg.style.borderRadius = '10px'
      qrImg.style.boxShadow = '0 4px 8px rgba(0,0,0,0.1)'

      qrCodeRef.value.appendChild(qrImg)
    }

    const copyId = async () => {
      try {
        await navigator.clipboard.writeText(lineBotId)
        alert('已複製 LINE 機器人 ID 到剪貼板！')
      } catch (err) {
        // 備用方案
        const textArea = document.createElement('textarea')
        textArea.value = lineBotId
        document.body.appendChild(textArea)
        textArea.select()
        document.execCommand('copy')
        document.body.removeChild(textArea)
        alert('已複製 LINE 機器人 ID 到剪貼板！')
      }
    }

    const openLineApp = () => {
      window.open(lineBotUrl, '_blank')
    }


    onMounted(() => {
      if (props.show) {
        nextTick(() => {
          generateQRCode()
        })
      }
    })

    return {
      qrCodeRef,
      lineBotId,
      closeModal,
      copyId,
      openLineApp
    }
  },
  watch: {
    show(newVal) {
      if (newVal) {
        this.$nextTick(() => {
          this.generateQRCode()
        })
      }
    }
  },
  methods: {
    generateQRCode() {
      if (!this.$refs.qrCodeRef) return

      // 清空容器
      this.$refs.qrCodeRef.innerHTML = ''

      // 使用 LINE 官方 QR Code 圖片
      const qrImg = document.createElement('img')
      qrImg.src = 'https://qr-official.line.me/gs/M_351thdpd_GW.png?oat_content=qr'
      qrImg.alt = 'LINE Bot QR Code'
      qrImg.style.width = '200px'
      qrImg.style.height = '200px'
      qrImg.style.borderRadius = '10px'
      qrImg.style.boxShadow = '0 4px 8px rgba(0,0,0,0.1)'

      this.$refs.qrCodeRef.appendChild(qrImg)
    }
  }
}
</script>

<style scoped>
.line-bot-qr-modal {
  position: fixed;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  background: rgba(0, 0, 0, 0.5);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 2000;
  animation: fadeIn 0.3s ease-out;
}

.modal-content {
  background: white;
  border-radius: 15px;
  max-width: 500px;
  width: 90%;
  max-height: 80vh;
  overflow-y: auto;
  box-shadow: 0 20px 40px rgba(0, 0, 0, 0.2);
  animation: slideInUp 0.3s ease-out;
}

.modal-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 20px;
  border-bottom: 1px solid #eee;
  background: linear-gradient(135deg, #00c300 0%, #00a000 100%);
  color: white;
  border-radius: 15px 15px 0 0;
}

.modal-header h3 {
  margin: 0;
  font-size: 1.3rem;
  font-weight: 600;
}

.close-btn {
  background: none;
  border: none;
  color: white;
  font-size: 24px;
  cursor: pointer;
  padding: 0;
  width: 30px;
  height: 30px;
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: 50%;
  transition: background 0.3s;
}

.close-btn:hover {
  background: rgba(255, 255, 255, 0.2);
}

.modal-body {
  padding: 30px;
}

.qr-section {
  text-align: center;
  margin-bottom: 30px;
}

.qr-code {
  margin-bottom: 15px;
}

.qr-container {
  display: flex;
  justify-content: center;
  align-items: center;
  min-height: 200px;
}

.qr-description {
  color: #666;
  font-size: 0.95rem;
  margin: 0;
}

.alternative-methods {
  border-top: 1px solid #eee;
  padding-top: 20px;
}

.alternative-methods h4 {
  margin: 0 0 15px 0;
  color: #333;
  font-size: 1.1rem;
}

.method-item {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 12px;
  margin: 10px 0;
  background: #f8f9fa;
  border-radius: 8px;
  border: 1px solid #e9ecef;
}

.method-icon {
  font-size: 1.2rem;
  flex-shrink: 0;
}

.method-text {
  flex: 1;
  font-size: 0.95rem;
  color: #555;
}

.copy-btn,
.link-btn {
  background: #00c300;
  color: white;
  border: none;
  padding: 6px 12px;
  border-radius: 4px;
  font-size: 0.85rem;
  cursor: pointer;
  transition: background 0.3s;
  flex-shrink: 0;
}

.copy-btn:hover,
.link-btn:hover {
  background: #00a000;
}

.link-btn {
  background: #007bff;
}

.link-btn:hover {
  background: #0056b3;
}


@keyframes fadeIn {
  from {
    opacity: 0;
  }
  to {
    opacity: 1;
  }
}

@keyframes slideInUp {
  from {
    opacity: 0;
    transform: translateY(30px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

@media (max-width: 768px) {
  .modal-content {
    width: 95%;
    margin: 20px;
  }

  .modal-body {
    padding: 20px;
  }

  .qr-container {
    min-height: 150px;
  }

  .qr-container img {
    width: 150px !important;
    height: 150px !important;
  }

  .method-item {
    flex-direction: column;
    align-items: flex-start;
    gap: 8px;
  }

  .method-text {
    margin-bottom: 5px;
  }
}
</style>
