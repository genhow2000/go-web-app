<template>
  <div class="ai-chat-window show">
    <div class="chat-header">
      <div>
        <h3>AI 購物助手</h3>
        <div class="status">在線</div>
      </div>
      <button class="chat-close" @click="$emit('close')">×</button>
    </div>
    
    <div class="chat-messages" ref="messagesContainer">
      <div class="message assistant">
        <div class="message-content">
          您好！我是阿和商城的AI購物助手，有什麼可以幫助您的嗎？
        </div>
        <div class="message-time">{{ getCurrentTime() }}</div>
      </div>
      
      <div 
        v-for="message in messages" 
        :key="message.id"
        class="message"
        :class="message.role"
      >
        <div class="message-content">{{ message.content }}</div>
        <div class="message-time">{{ message.time }}</div>
      </div>
    </div>
    
    <div v-if="isTyping" class="typing-indicator">
      <span>AI 正在思考</span>
      <div class="typing-dots">
        <div class="typing-dot"></div>
        <div class="typing-dot"></div>
        <div class="typing-dot"></div>
      </div>
    </div>
    
    <div class="chat-input-container">
      <div class="chat-input-wrapper">
        <textarea 
          v-model="inputMessage"
          class="chat-input" 
          placeholder="輸入您的問題..."
          rows="1"
          @keypress.enter.prevent="sendMessage"
          @input="autoResize"
          ref="inputRef"
        ></textarea>
        <button 
          class="chat-send" 
          @click="sendMessage"
          :disabled="!inputMessage.trim() || isSending"
        >
          ➤
        </button>
      </div>
    </div>
  </div>
</template>

<script>
import { ref, nextTick, onMounted } from 'vue'
import api from '@/services/api'

export default {
  name: 'AIChatWindow',
  emits: ['close'],
  setup() {
    const inputMessage = ref('')
    const messages = ref([])
    const isTyping = ref(false)
    const isSending = ref(false)
    const messagesContainer = ref(null)
    const inputRef = ref(null)
    const conversationId = ref(null)

    // 獲取當前時間
    const getCurrentTime = () => {
      const now = new Date()
      return now.toLocaleTimeString('zh-TW', { 
        hour: '2-digit', 
        minute: '2-digit' 
      })
    }

    // 自動調整輸入框高度
    const autoResize = () => {
      nextTick(() => {
        if (inputRef.value) {
          inputRef.value.style.height = 'auto'
          inputRef.value.style.height = Math.min(inputRef.value.scrollHeight, 100) + 'px'
        }
      })
    }

    // 滾動到底部
    const scrollToBottom = () => {
      nextTick(() => {
        if (messagesContainer.value) {
          messagesContainer.value.scrollTop = messagesContainer.value.scrollHeight
        }
      })
    }

    // 創建對話
    const createConversation = async () => {
      try {
        const response = await api.post('/api/chat/conversations', {
          title: 'AI 購物助手對話'
        })
        return response.data.conversation_id
      } catch (error) {
        console.error('創建對話失敗:', error)
        return null
      }
    }

    // 發送消息
    const sendMessage = async () => {
      const message = inputMessage.value.trim()
      if (!message || isSending.value) return

      // 添加用戶消息
      const userMessage = {
        id: Date.now(),
        role: 'user',
        content: message,
        time: getCurrentTime()
      }
      messages.value.push(userMessage)
      inputMessage.value = ''
      autoResize()
      scrollToBottom()

      isSending.value = true
      isTyping.value = true

      try {
        // 如果沒有對話ID，先創建對話
        if (!conversationId.value) {
          conversationId.value = await createConversation()
        }

        // 發送消息到後端
        const response = await api.post('/api/chat/send', {
          conversation_id: conversationId.value,
          message: message
        })

        if (response.data.success) {
          // 添加AI回復
          if (response.data.ai_message) {
            const aiMessage = {
              id: Date.now() + 1,
              role: 'assistant',
              content: response.data.ai_message.content,
              time: getCurrentTime()
            }
            messages.value.push(aiMessage)
          } else {
            // 降級到模擬AI回復
            setTimeout(() => {
              const aiMessage = {
                id: Date.now() + 1,
                role: 'assistant',
                content: getAIResponse(message),
                time: getCurrentTime()
              }
              messages.value.push(aiMessage)
              scrollToBottom()
            }, 1000 + Math.random() * 2000)
          }
        } else {
          throw new Error(response.data.error || '發送失敗')
        }
      } catch (error) {
        console.error('發送消息失敗:', error)
        // 添加錯誤消息
        const errorMessage = {
          id: Date.now() + 1,
          role: 'assistant',
          content: '抱歉，我現在無法回應。請稍後再試。',
          time: getCurrentTime()
        }
        messages.value.push(errorMessage)
      } finally {
        isSending.value = false
        isTyping.value = false
        scrollToBottom()
      }
    }

    // 模擬AI回復
    const getAIResponse = (userMessage) => {
      const responses = [
        "我了解您的需求，讓我為您推薦一些相關商品。",
        "這是一個很好的問題！根據您的描述，我建議您查看以下分類的商品。",
        "感謝您的詢問！我可以幫您找到最適合的商品。",
        "我明白您想要什麼了，讓我為您搜索相關商品。",
        "好的，我會根據您的需求為您推薦商品。"
      ]
      
      const message = userMessage.toLowerCase()
      if (message.includes('價格') || message.includes('多少錢')) {
        return "我們有各種價格區間的商品，從經濟實惠到高端精品都有。您可以在商品頁面查看詳細價格信息。"
      } else if (message.includes('推薦') || message.includes('建議')) {
        return "根據您的需求，我推薦您查看我們的精選商品。這些商品都經過嚴格篩選，品質有保證。"
      } else if (message.includes('配送') || message.includes('運費')) {
        return "我們提供快速配送服務，24小時內發貨，3-5天送達。滿額還有免運費優惠！"
      } else if (message.includes('退換') || message.includes('售後')) {
        return "我們提供7天無理由退換貨服務，讓您買得放心。如有任何問題，我們的客服團隊隨時為您服務。"
      } else {
        return responses[Math.floor(Math.random() * responses.length)]
      }
    }

    onMounted(() => {
      scrollToBottom()
    })

    return {
      inputMessage,
      messages,
      isTyping,
      isSending,
      messagesContainer,
      inputRef,
      getCurrentTime,
      autoResize,
      sendMessage
    }
  }
}
</script>

<style scoped>
.ai-chat-window {
  position: fixed;
  bottom: 90px;
  left: 20px;
  width: 350px;
  height: 500px;
  background: white;
  border-radius: 15px;
  box-shadow: 0 10px 30px rgba(0,0,0,0.2);
  z-index: 1001;
  display: flex;
  flex-direction: column;
  overflow: hidden;
  border: 1px solid #e2e8f0;
  animation: slideUp 0.3s ease-out;
}

@keyframes slideUp {
  from {
    opacity: 0;
    transform: translateY(20px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

.chat-header {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
  padding: 15px 20px;
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.chat-header h3 {
  margin: 0;
  font-size: 16px;
  font-weight: 600;
}

.chat-header .status {
  font-size: 12px;
  opacity: 0.9;
}

.chat-close {
  background: none;
  border: none;
  color: white;
  font-size: 20px;
  cursor: pointer;
  padding: 0;
  width: 24px;
  height: 24px;
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: 50%;
  transition: background 0.3s;
}

.chat-close:hover {
  background: rgba(255,255,255,0.2);
}

.chat-messages {
  flex: 1;
  padding: 20px;
  overflow-y: auto;
  background: #f8f9fa;
}

.message {
  margin-bottom: 15px;
  display: flex;
  align-items: flex-start;
}

.message.user {
  justify-content: flex-end;
}

.message.assistant {
  justify-content: flex-start;
}

.message-content {
  max-width: 80%;
  padding: 10px 15px;
  border-radius: 18px;
  font-size: 14px;
  line-height: 1.4;
}

.message.user .message-content {
  background: #667eea;
  color: white;
  border-bottom-right-radius: 5px;
}

.message.assistant .message-content {
  background: white;
  color: #333;
  border: 1px solid #e2e8f0;
  border-bottom-left-radius: 5px;
}

.message-time {
  font-size: 11px;
  color: #999;
  margin-top: 5px;
  text-align: right;
}

.message.assistant .message-time {
  text-align: left;
}

.chat-input-container {
  padding: 15px 20px;
  background: white;
  border-top: 1px solid #e2e8f0;
}

.chat-input-wrapper {
  display: flex;
  gap: 10px;
  align-items: center;
}

.chat-input {
  flex: 1;
  padding: 10px 15px;
  border: 1px solid #e2e8f0;
  border-radius: 25px;
  outline: none;
  font-size: 14px;
  resize: none;
  max-height: 100px;
}

.chat-input:focus {
  border-color: #667eea;
}

.chat-send {
  background: #667eea;
  color: white;
  border: none;
  border-radius: 50%;
  width: 40px;
  height: 40px;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  transition: background 0.3s;
}

.chat-send:hover:not(:disabled) {
  background: #5a6fd8;
}

.chat-send:disabled {
  background: #ccc;
  cursor: not-allowed;
}

.typing-indicator {
  display: flex;
  align-items: center;
  gap: 5px;
  color: #666;
  font-size: 12px;
  font-style: italic;
  padding: 0 20px 10px;
}

.typing-dots {
  display: flex;
  gap: 2px;
}

.typing-dot {
  width: 4px;
  height: 4px;
  background: #666;
  border-radius: 50%;
  animation: typing 1.4s infinite;
}

.typing-dot:nth-child(2) {
  animation-delay: 0.2s;
}

.typing-dot:nth-child(3) {
  animation-delay: 0.4s;
}

@keyframes typing {
  0%, 60%, 100% {
    transform: translateY(0);
  }
  30% {
    transform: translateY(-10px);
  }
}

@media (max-width: 768px) {
  .ai-chat-window {
    width: calc(100vw - 40px);
    left: 20px;
    right: 20px;
    height: 400px;
  }
}
</style>
