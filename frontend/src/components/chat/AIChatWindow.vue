<template>
  <div class="ai-chat-window show">
    <div class="chat-header">
      <div>
        <h3>{{ getChatTitle() }}</h3>
        <div class="status">在線</div>
        <div v-if="stockContext" class="stock-context">
          正在分析：{{ stockContext.name }} ({{ stockContext.code }})
        </div>
      </div>
      <button class="chat-close" @click="$emit('close')">×</button>
    </div>
    
    <div class="chat-messages" ref="messagesContainer">
      <div class="message assistant">
        <div class="message-content">
          {{ getWelcomeMessage() }}
        </div>
        <div class="message-time">{{ getCurrentTime() }}</div>
      </div>
      
      <div 
        v-for="message in messages" 
        :key="message.id"
        class="message"
        :class="{ 
          'user': message.role === 'user',
          'assistant': message.role === 'assistant',
          'system': message.role === 'system',
          'error': message.isError,
          'warning': message.isWarning
        }"
      >
        <div class="message-content">
          <span v-if="message.role === 'system'" class="error-icon">
            <span v-if="message.errorType === 'api_error'">🔌</span>
            <span v-else-if="message.errorType === 'quota_exceeded'">⚠️</span>
            <span v-else-if="message.errorType === 'rate_limited'">⏰</span>
            <span v-else>⚠️</span>
          </span>
          {{ message.content }}
        </div>
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
    
    <!-- 股票相關預設問題 -->
    <div v-if="stockContext && messages.length === 0" class="quick-questions">
      <h4>快速提問：</h4>
      <div class="question-buttons">
        <button @click="askQuestion('這支股票值得買嗎？')" class="question-btn">
          這支股票值得買嗎？
        </button>
        <button @click="askQuestion('分析這支股票的技術指標')" class="question-btn">
          分析技術指標
        </button>
        <button @click="askQuestion('這支股票的投資風險如何？')" class="question-btn">
          投資風險分析
        </button>
        <button @click="askQuestion('這支股票的基本面如何？')" class="question-btn">
          基本面分析
        </button>
      </div>
    </div>

    <div class="chat-input-container">
      <div class="chat-input-wrapper">
        <textarea 
          v-model="inputMessage"
          class="chat-input" 
          :placeholder="getInputPlaceholder()"
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
  props: {
    stockContext: {
      type: Object,
      default: null
    }
  },
  emits: ['close'],
  setup(props) {
    const inputMessage = ref('')
    const messages = ref([])
    const isTyping = ref(false)
    const isSending = ref(false)
    const messagesContainer = ref(null)
    const inputRef = ref(null)
    const conversationId = ref(null)

    // 獲取聊天標題
    const getChatTitle = () => {
      return props.stockContext ? 'AI 股票助手' : 'AI 購物助手'
    }

    // 獲取歡迎訊息
    const getWelcomeMessage = () => {
      console.log('getWelcomeMessage - 股票上下文:', props.stockContext)
      if (props.stockContext) {
        return `您好！我是AI股票助手，正在為您分析 ${props.stockContext.name} (${props.stockContext.code})。我可以幫您分析這支股票的投資價值、技術指標、風險評估等。有什麼問題可以問我！`
      }
      return '您好！我是阿和商城的AI購物助手，有什麼可以幫助您的嗎？'
    }

    // 獲取輸入框提示文字
    const getInputPlaceholder = () => {
      if (props.stockContext) {
        return `詢問關於 ${props.stockContext.name} 的問題...`
      }
      return '輸入您的問題...'
    }

    // 快速提問
    const askQuestion = (question) => {
      inputMessage.value = question
      sendMessage()
    }

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

        // 準備發送的消息
        const messageData = {
          conversation_id: conversationId.value,
          message: message
        }
        
        // 如果有股票上下文，只傳遞簡短的關鍵資訊
        if (props.stockContext) {
          console.log('股票上下文:', props.stockContext)
          messageData.stock_context = {
            code: props.stockContext.code,
            name: props.stockContext.name,
            current_price: props.stockContext.price?.price || 0,
            change: props.stockContext.price?.change || 0,
            market: props.stockContext.market
          }
          console.log('發送的股票上下文:', messageData.stock_context)
        } else {
          console.log('沒有股票上下文')
        }

        // 發送消息到後端
        const response = await api.post('/api/chat/send', messageData)

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
            
            // 檢查是否有API錯誤或使用上限提醒
            if (response.data.api_error) {
              const errorMessage = {
                id: Date.now() + 2,
                role: 'system',
                content: response.data.api_error,
                time: getCurrentTime(),
                isError: true,
                errorType: response.data.error_type || 'api_error'
              }
              messages.value.push(errorMessage)
            }
            
            // 檢查使用統計和警告
            if (response.data.usage_stats) {
              const dailyCount = response.data.usage_stats.daily_requests || 0
              const dailyLimit = response.data.usage_stats.daily_limit || 0
              
              if (dailyCount >= dailyLimit) {
                const limitMessage = {
                  id: Date.now() + 3,
                  role: 'system',
                  content: '您今日的使用次數已達上限，請登入會員以提高使用限制',
                  time: getCurrentTime(),
                  isWarning: true
                }
                messages.value.push(limitMessage)
              } else if (response.data.warning) {
                const warningMessage = {
                  id: Date.now() + 3,
                  role: 'system',
                  content: response.data.warning,
                  time: getCurrentTime(),
                  isWarning: true
                }
                messages.value.push(warningMessage)
              }
            }
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
      getChatTitle,
      getWelcomeMessage,
      getInputPlaceholder,
      askQuestion,
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
  top: 50%;
  left: 50%;
  transform: translate(-50%, -50%);
  width: 600px;
  max-height: 80vh;
  min-height: 500px;
  background: white;
  border-radius: 15px;
  box-shadow: 0 20px 60px rgba(0,0,0,0.3);
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

.stock-context {
  font-size: 12px;
  color: #667eea;
  font-weight: 600;
  margin-top: 4px;
  padding: 4px 8px;
  background: rgba(102, 126, 234, 0.1);
  border-radius: 12px;
  display: inline-block;
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
  max-height: none;
  min-height: 0;
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
  word-wrap: break-word;
  white-space: pre-wrap;
  overflow-wrap: break-word;
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

.message.system {
  justify-content: center;
  margin: 10px 20%;
}

.message.system .message-content {
  background: #f8f9fa;
  color: #6c757d;
  border: 1px solid #e9ecef;
  border-radius: 8px;
  font-size: 13px;
  text-align: center;
  max-width: 100%;
}

.message.error .message-content {
  background: #f8d7da;
  color: #721c24;
  border-color: #f5c6cb;
}

.message.warning .message-content {
  background: #fff3cd;
  color: #856404;
  border-color: #ffeaa7;
}

.error-icon {
  margin-right: 8px;
  font-size: 16px;
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

/* 快速問題區域 */
.quick-questions {
  padding: 15px 20px;
  background: #f8f9fa;
  border-top: 1px solid #e2e8f0;
}

.quick-questions h4 {
  margin: 0 0 10px 0;
  font-size: 14px;
  color: #4a5568;
  font-weight: 600;
}

.question-buttons {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
}

.question-btn {
  padding: 8px 12px;
  background: white;
  border: 1px solid #e2e8f0;
  border-radius: 20px;
  color: #4a5568;
  font-size: 12px;
  cursor: pointer;
  transition: all 0.2s ease;
  white-space: nowrap;
}

.question-btn:hover {
  background: #667eea;
  color: white;
  border-color: #667eea;
  transform: translateY(-1px);
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

/* 響應式設計 */
@media (max-width: 768px) {
  .ai-chat-window {
    width: calc(100vw - 40px);
    height: 80vh;
    max-height: 600px;
    min-height: 400px;
    top: 50%;
    left: 50%;
    transform: translate(-50%, -50%);
    margin: 0 20px;
    position: fixed;
    z-index: 1001;
  }
  
  .question-buttons {
    flex-direction: column;
  }
  
  .question-btn {
    text-align: center;
    white-space: normal;
  }
}

/* 小屏幕手機優化 */
@media (max-width: 480px) {
  .ai-chat-window {
    width: calc(100vw - 20px);
    height: 85vh;
    max-height: 500px;
    margin: 0 10px;
  }
  
  .chat-header h3 {
    font-size: 14px;
  }
  
  .chat-messages {
    padding: 15px;
  }
  
  .message-content {
    max-width: 85%;
    font-size: 13px;
  }
}

@keyframes typing {
  0%, 60%, 100% {
    transform: translateY(0);
  }
  30% {
    transform: translateY(-10px);
  }
}
</style>
