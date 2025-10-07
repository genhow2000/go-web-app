# LINE Bot 機器人設置指南

## 概述

本指南將幫助你設置 LINE Bot 機器人，與現有的 LINE OAuth 登入系統整合。LINE Bot 可以通過 Webhook 接收用戶訊息並自動回覆。

## 目錄

1. [LINE Bot 開發者帳號設置](#1-line-bot-開發者帳號設置)
2. [創建 LINE Bot Channel](#2-創建-line-bot-channel)
3. [配置 Webhook URL](#3-配置-webhook-url)
4. [後端實現](#4-後端實現)
5. [前端整合](#5-前端整合)
6. [測試與部署](#6-測試與部署)
7. [進階功能](#7-進階功能)

## 1. LINE Bot 開發者帳號設置

### 1.1 註冊 LINE Developers 帳號

1. 前往 [LINE Developers Console](https://developers.line.biz/)
2. 使用你的 LINE 帳號登入
3. 完成開發者帳號驗證

### 1.2 創建 Provider

1. 在 LINE Developers Console 中點擊 "Create"
2. 選擇 "Provider"
3. 填寫 Provider 資訊：
   - Provider name: 你的應用名稱
   - Provider description: 簡短描述

## 2. 創建 LINE Bot Channel

### 2.1 創建 Messaging API Channel

1. 在 Provider 頁面點擊 "Create"
2. 選擇 "Messaging API"
3. 填寫 Channel 資訊：
   - Channel name: 機器人顯示名稱
   - Channel description: 機器人描述
   - Category: 選擇適當的分類
   - Subcategory: 選擇子分類

### 2.2 獲取 Channel 資訊

創建完成後，記錄以下重要資訊：

```
Channel ID: 你的 Channel ID
Channel Secret: 你的 Channel Secret
Channel Access Token: 你的 Channel Access Token
```

### 2.3 配置 Bot 設定

1. 在 Channel 設定頁面：
   - 啟用 "Use webhook"
   - 設定 "Webhook URL" (稍後配置)
   - 選擇 "Allow bot to join group chats" (如需要)
   - 設定 "Auto-reply messages" (可選)

## 3. 配置 Webhook URL

### 3.1 本地開發環境

對於本地開發，你需要使用 ngrok 或其他隧道服務：

```bash
# 安裝 ngrok
brew install ngrok  # macOS
# 或下載從 https://ngrok.com/

# 啟動隧道
ngrok http 8080
```

記錄 ngrok 提供的 HTTPS URL，例如：`https://abc123.ngrok.io`

### 3.2 生產環境

對於生產環境，使用你的實際域名：

```
https://yourdomain.com/webhook/line
```

## 4. 後端實現

### 4.1 添加 LINE Bot 配置

在 `config/config.go` 中添加 LINE Bot 配置：

```go
// LineBotConfig LINE Bot 配置
type LineBotConfig struct {
    ChannelID     string `json:"channel_id"`
    ChannelSecret string `json:"channel_secret"`
    ChannelToken  string `json:"channel_token"`
    WebhookURL    string `json:"webhook_url"`
}

// OAuthConfig OAuth配置
type OAuthConfig struct {
    LINE    LineOAuthConfig `json:"line"`
    LineBot LineBotConfig   `json:"line_bot"`  // 新增
}
```

### 4.2 創建 LINE Bot 服務

創建 `services/line_bot_service.go`：

```go
package services

import (
    "bytes"
    "crypto/hmac"
    "crypto/sha256"
    "encoding/base64"
    "encoding/json"
    "fmt"
    "io"
    "net/http"
    "go-simple-app/config"
    "go-simple-app/logger"
    "github.com/sirupsen/logrus"
)

type LineBotService struct {
    config *config.LineBotConfig
    client *http.Client
}

type LineWebhookEvent struct {
    Type       string                 `json:"type"`
    ReplyToken string                 `json:"replyToken"`
    Source     LineSource             `json:"source"`
    Timestamp  int64                  `json:"timestamp"`
    Message    *LineMessage           `json:"message"`
    Postback   *LinePostback          `json:"postback"`
}

type LineSource struct {
    Type    string `json:"type"`
    UserID  string `json:"userId"`
    GroupID string `json:"groupId"`
    RoomID  string `json:"roomId"`
}

type LineMessage struct {
    ID   string `json:"id"`
    Type string `json:"type"`
    Text string `json:"text"`
}

type LinePostback struct {
    Data string `json:"data"`
}

type LineReplyMessage struct {
    ReplyToken string           `json:"replyToken"`
    Messages   []LineTextMessage `json:"messages"`
}

type LineTextMessage struct {
    Type string `json:"type"`
    Text string `json:"text"`
}

func NewLineBotService(config *config.LineBotConfig) *LineBotService {
    return &LineBotService{
        config: config,
        client: &http.Client{},
    }
}

// 驗證 Webhook 簽名
func (s *LineBotService) VerifySignature(body []byte, signature string) bool {
    hash := hmac.New(sha256.New, []byte(s.config.ChannelSecret))
    hash.Write(body)
    expectedSignature := base64.StdEncoding.EncodeToString(hash.Sum(nil))
    return hmac.Equal([]byte(signature), []byte(expectedSignature))
}

// 處理 Webhook 事件
func (s *LineBotService) HandleWebhook(body []byte) error {
    var webhookData struct {
        Events []LineWebhookEvent `json:"events"`
    }

    if err := json.Unmarshal(body, &webhookData); err != nil {
        return fmt.Errorf("failed to parse webhook data: %v", err)
    }

    for _, event := range webhookData.Events {
        if err := s.processEvent(event); err != nil {
            logger.Error("處理 LINE 事件失敗", err, logrus.Fields{
                "event_type": event.Type,
                "user_id": event.Source.UserID,
            })
        }
    }

    return nil
}

// 處理單個事件
func (s *LineBotService) processEvent(event LineWebhookEvent) error {
    switch event.Type {
    case "message":
        if event.Message != nil && event.Message.Type == "text" {
            return s.handleTextMessage(event)
        }
    case "postback":
        if event.Postback != nil {
            return s.handlePostback(event)
        }
    }
    return nil
}

// 處理文字訊息
func (s *LineBotService) handleTextMessage(event LineWebhookEvent) error {
    userMessage := event.Message.Text
    replyToken := event.ReplyToken
    userID := event.Source.UserID

    logger.Info("收到 LINE 訊息", logrus.Fields{
        "user_id": userID,
        "message": userMessage,
    })

    // 這裡可以整合你的 AI 服務
    responseText := s.generateResponse(userMessage, userID)

    return s.replyMessage(replyToken, responseText)
}

// 處理 Postback 事件
func (s *LineBotService) handlePostback(event LineWebhookEvent) error {
    data := event.Postback.Data
    replyToken := event.ReplyToken
    userID := event.Source.UserID

    logger.Info("收到 LINE Postback", logrus.Fields{
        "user_id": userID,
        "data": data,
    })

    // 處理 Postback 數據
    responseText := s.handlePostbackData(data, userID)

    return s.replyMessage(replyToken, responseText)
}

// 生成回應
func (s *LineBotService) generateResponse(userMessage, userID string) string {
    // 這裡可以整合你的 AI 聊天服務
    // 例如調用現有的 chat_service

    // 簡單的回應範例
    switch userMessage {
    case "你好", "hi", "hello":
        return "你好！我是你的購物助手，有什麼可以幫助你的嗎？"
    case "商品", "產品", "products":
        return "你可以在我們的網站上瀏覽商品：https://yourdomain.com"
    case "購物車", "cart":
        return "你可以在這裡查看購物車：https://yourdomain.com/cart"
    default:
        return "感謝你的訊息！如需更多幫助，請訪問我們的網站：https://yourdomain.com"
    }
}

// 處理 Postback 數據
func (s *LineBotService) handlePostbackData(data, userID string) string {
    // 根據 Postback 數據處理不同的動作
    switch data {
    case "view_products":
        return "這是我們的商品目錄：https://yourdomain.com/products"
    case "view_cart":
        return "這是你的購物車：https://yourdomain.com/cart"
    case "contact_support":
        return "如需客服協助，請聯繫我們：support@yourdomain.com"
    default:
        return "收到你的選擇，正在處理中..."
    }
}

// 回覆訊息
func (s *LineBotService) replyMessage(replyToken, text string) error {
    replyData := LineReplyMessage{
        ReplyToken: replyToken,
        Messages: []LineTextMessage{
            {
                Type: "text",
                Text: text,
            },
        },
    }

    jsonData, err := json.Marshal(replyData)
    if err != nil {
        return fmt.Errorf("failed to marshal reply data: %v", err)
    }

    req, err := http.NewRequest("POST", "https://api.line.me/v2/bot/message/reply", bytes.NewBuffer(jsonData))
    if err != nil {
        return fmt.Errorf("failed to create request: %v", err)
    }

    req.Header.Set("Content-Type", "application/json")
    req.Header.Set("Authorization", "Bearer "+s.config.ChannelToken)

    resp, err := s.client.Do(req)
    if err != nil {
        return fmt.Errorf("failed to send reply: %v", err)
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        body, _ := io.ReadAll(resp.Body)
        return fmt.Errorf("LINE API error: %d, %s", resp.StatusCode, string(body))
    }

    logger.Info("成功回覆 LINE 訊息", logrus.Fields{
        "reply_token": replyToken,
        "text": text,
    })

    return nil
}

// 發送推播訊息
func (s *LineBotService) SendPushMessage(userID, text string) error {
    pushData := struct {
        To       string           `json:"to"`
        Messages []LineTextMessage `json:"messages"`
    }{
        To: userID,
        Messages: []LineTextMessage{
            {
                Type: "text",
                Text: text,
            },
        },
    }

    jsonData, err := json.Marshal(pushData)
    if err != nil {
        return fmt.Errorf("failed to marshal push data: %v", err)
    }

    req, err := http.NewRequest("POST", "https://api.line.me/v2/bot/message/push", bytes.NewBuffer(jsonData))
    if err != nil {
        return fmt.Errorf("failed to create request: %v", err)
    }

    req.Header.Set("Content-Type", "application/json")
    req.Header.Set("Authorization", "Bearer "+s.config.ChannelToken)

    resp, err := s.client.Do(req)
    if err != nil {
        return fmt.Errorf("failed to send push message: %v", err)
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        body, _ := io.ReadAll(resp.Body)
        return fmt.Errorf("LINE API error: %d, %s", resp.StatusCode, string(body))
    }

    return nil
}
```

### 4.3 創建 LINE Bot 控制器

創建 `controllers/line_bot_controller.go`：

```go
package controllers

import (
    "io"
    "net/http"
    "go-simple-app/services"
    "go-simple-app/logger"
    "github.com/gin-gonic/gin"
    "github.com/sirupsen/logrus"
)

type LineBotController struct {
    lineBotService *services.LineBotService
}

func NewLineBotController(lineBotService *services.LineBotService) *LineBotController {
    return &LineBotController{
        lineBotService: lineBotService,
    }
}

// Webhook 處理器
func (c *LineBotController) Webhook(ctx *gin.Context) {
    // 讀取請求體
    body, err := io.ReadAll(ctx.Request.Body)
    if err != nil {
        logger.Error("讀取 LINE Webhook 請求體失敗", err, nil)
        ctx.JSON(http.StatusBadRequest, gin.H{"error": "Failed to read request body"})
        return
    }

    // 驗證簽名
    signature := ctx.GetHeader("X-Line-Signature")
    if !c.lineBotService.VerifySignature(body, signature) {
        logger.Error("LINE Webhook 簽名驗證失敗", nil, logrus.Fields{
            "signature": signature,
        })
        ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid signature"})
        return
    }

    // 處理 Webhook 事件
    if err := c.lineBotService.HandleWebhook(body); err != nil {
        logger.Error("處理 LINE Webhook 失敗", err, nil)
        ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to handle webhook"})
        return
    }

    ctx.JSON(http.StatusOK, gin.H{"status": "success"})
}

// 發送推播訊息
func (c *LineBotController) SendPushMessage(ctx *gin.Context) {
    var request struct {
        UserID string `json:"user_id" binding:"required"`
        Text   string `json:"text" binding:"required"`
    }

    if err := ctx.ShouldBindJSON(&request); err != nil {
        ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    if err := c.lineBotService.SendPushMessage(request.UserID, request.Text); err != nil {
        logger.Error("發送 LINE 推播訊息失敗", err, logrus.Fields{
            "user_id": request.UserID,
            "text": request.Text,
        })
        ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send push message"})
        return
    }

    ctx.JSON(http.StatusOK, gin.H{"status": "success"})
}
```

### 4.4 更新配置

更新 `config/config.go` 中的配置加載：

```go
OAuth: OAuthConfig{
    LINE: LineOAuthConfig{
        ClientID:     getEnv("LINE_CLIENT_ID", ""),
        ClientSecret: getEnv("LINE_CLIENT_SECRET", ""),
        RedirectURL:  getEnv("LINE_REDIRECT_URL", "http://localhost:8080/auth/line/callback"),
        Scopes:       []string{"profile", "openid"},
    },
    LineBot: LineBotConfig{
        ChannelID:     getEnv("LINE_BOT_CHANNEL_ID", ""),
        ChannelSecret: getEnv("LINE_BOT_CHANNEL_SECRET", ""),
        ChannelToken:  getEnv("LINE_BOT_CHANNEL_TOKEN", ""),
        WebhookURL:    getEnv("LINE_BOT_WEBHOOK_URL", "https://yourdomain.com/webhook/line"),
    },
},
```

### 4.5 添加路由

在 `routes/routes.go` 中添加 LINE Bot 路由：

```go
// LINE Bot 路由
lineBotController := controllers.NewLineBotController(services.NewLineBotService(&config.OAuth.LineBot))
api.POST("/webhook/line", lineBotController.Webhook)
api.POST("/line/push", lineBotController.SendPushMessage)
```

## 5. 前端整合

### 5.1 添加 LINE Bot 相關 API

在 `frontend/src/services/api.js` 中添加：

```javascript
// LINE Bot 相關 API
export const lineBotAPI = {
  // 發送推播訊息
  sendPushMessage: async (userId, text) => {
    const response = await fetch("/api/line/push", {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
        Authorization: `Bearer ${localStorage.getItem("auth_token")}`,
      },
      body: JSON.stringify({ user_id: userId, text }),
    });
    return response.json();
  },
};
```

### 5.2 在管理後台添加 LINE Bot 管理

創建 `frontend/src/views/admin/LineBotManagement.vue`：

```vue
<template>
  <div class="line-bot-management">
    <h2>LINE Bot 管理</h2>

    <div class="bot-info">
      <h3>機器人資訊</h3>
      <p>Channel ID: {{ botInfo.channelId }}</p>
      <p>Webhook URL: {{ botInfo.webhookUrl }}</p>
      <p>狀態: {{ botInfo.status }}</p>
    </div>

    <div class="send-message">
      <h3>發送推播訊息</h3>
      <form @submit.prevent="sendMessage">
        <div class="form-group">
          <label>用戶 ID:</label>
          <input v-model="messageForm.userId" type="text" required />
        </div>
        <div class="form-group">
          <label>訊息內容:</label>
          <textarea v-model="messageForm.text" required></textarea>
        </div>
        <button type="submit" :disabled="sending">
          {{ sending ? "發送中..." : "發送訊息" }}
        </button>
      </form>
    </div>

    <div class="message-history">
      <h3>訊息歷史</h3>
      <div
        v-for="message in messageHistory"
        :key="message.id"
        class="message-item"
      >
        <p><strong>用戶:</strong> {{ message.userId }}</p>
        <p><strong>內容:</strong> {{ message.text }}</p>
        <p><strong>時間:</strong> {{ message.timestamp }}</p>
      </div>
    </div>
  </div>
</template>

<script>
import { lineBotAPI } from "@/services/api.js";

export default {
  name: "LineBotManagement",
  data() {
    return {
      botInfo: {
        channelId: "",
        webhookUrl: "",
        status: "active",
      },
      messageForm: {
        userId: "",
        text: "",
      },
      messageHistory: [],
      sending: false,
    };
  },
  methods: {
    async sendMessage() {
      this.sending = true;
      try {
        await lineBotAPI.sendPushMessage(
          this.messageForm.userId,
          this.messageForm.text
        );
        this.messageHistory.unshift({
          id: Date.now(),
          userId: this.messageForm.userId,
          text: this.messageForm.text,
          timestamp: new Date().toLocaleString(),
        });
        this.messageForm = { userId: "", text: "" };
        alert("訊息發送成功！");
      } catch (error) {
        console.error("發送訊息失敗:", error);
        alert("發送訊息失敗，請重試");
      } finally {
        this.sending = false;
      }
    },
  },
};
</script>

<style scoped>
.line-bot-management {
  max-width: 800px;
  margin: 0 auto;
  padding: 20px;
}

.bot-info,
.send-message,
.message-history {
  background: #f5f5f5;
  padding: 20px;
  margin: 20px 0;
  border-radius: 8px;
}

.form-group {
  margin: 15px 0;
}

.form-group label {
  display: block;
  margin-bottom: 5px;
  font-weight: bold;
}

.form-group input,
.form-group textarea {
  width: 100%;
  padding: 10px;
  border: 1px solid #ddd;
  border-radius: 4px;
}

.form-group textarea {
  height: 100px;
  resize: vertical;
}

button {
  background: #007bff;
  color: white;
  padding: 10px 20px;
  border: none;
  border-radius: 4px;
  cursor: pointer;
}

button:disabled {
  background: #ccc;
  cursor: not-allowed;
}

.message-item {
  background: white;
  padding: 15px;
  margin: 10px 0;
  border-radius: 4px;
  border-left: 4px solid #007bff;
}
</style>
```

## 6. 測試與部署

### 6.1 環境變數設置

在 `.env` 文件中添加：

```bash
# LINE Bot 配置
LINE_BOT_CHANNEL_ID=your_channel_id
LINE_BOT_CHANNEL_SECRET=your_channel_secret
LINE_BOT_CHANNEL_TOKEN=your_channel_token
LINE_BOT_WEBHOOK_URL=https://yourdomain.com/webhook/line
```

### 6.2 本地測試

1. 啟動應用：

```bash
go run main.go
```

2. 使用 ngrok 創建隧道：

```bash
ngrok http 8080
```

3. 在 LINE Developers Console 中設置 Webhook URL：

```
https://your-ngrok-url.ngrok.io/webhook/line
```

4. 測試機器人：
   - 掃描 QR Code 添加機器人為好友
   - 發送訊息測試

### 6.3 生產部署

1. 更新 LINE Developers Console 中的 Webhook URL
2. 確保 HTTPS 證書有效
3. 部署到你的生產環境

## 7. 進階功能

### 7.1 整合 AI 聊天服務

修改 `line_bot_service.go` 中的 `generateResponse` 方法：

```go
func (s *LineBotService) generateResponse(userMessage, userID string) string {
    // 整合現有的 AI 聊天服務
    aiService := services.NewAIManager(s.aiConfig)

    response, err := aiService.GenerateResponse(userMessage, userID)
    if err != nil {
        logger.Error("AI 回應生成失敗", err, logrus.Fields{
            "user_id": userID,
            "message": userMessage,
        })
        return "抱歉，我暫時無法回應你的訊息。"
    }

    return response
}
```

### 7.2 添加富媒體訊息

支援圖片、按鈕、輪播等富媒體訊息：

```go
type LineTemplateMessage struct {
    Type     string                 `json:"type"`
    AltText  string                 `json:"altText"`
    Template LineTemplate           `json:"template"`
}

type LineTemplate struct {
    Type     string                 `json:"type"`
    Text     string                 `json:"text"`
    Actions  []LineAction           `json:"actions"`
}

type LineAction struct {
    Type  string `json:"type"`
    Label string `json:"label"`
    Data  string `json:"data"`
    URI   string `json:"uri"`
}
```

### 7.3 用戶狀態管理

追蹤用戶對話狀態：

```go
type UserSession struct {
    UserID    string    `json:"user_id"`
    State     string    `json:"state"`
    Context   map[string]interface{} `json:"context"`
    UpdatedAt time.Time `json:"updated_at"`
}
```

### 7.4 數據分析

記錄機器人使用統計：

```go
type BotAnalytics struct {
    Date        string `json:"date"`
    UserCount   int    `json:"user_count"`
    MessageCount int   `json:"message_count"`
    ResponseTime int   `json:"response_time"`
}
```

## 8. 安全注意事項

1. **驗證 Webhook 簽名**：確保所有請求都經過簽名驗證
2. **限制 API 存取**：使用適當的認證和授權
3. **資料保護**：遵循 GDPR 和相關隱私法規
4. **錯誤處理**：避免在錯誤訊息中洩露敏感資訊
5. **速率限制**：實施適當的速率限制防止濫用

## 9. 監控與維護

1. **日誌記錄**：記錄所有重要事件和錯誤
2. **性能監控**：監控回應時間和成功率
3. **用戶反饋**：收集用戶反饋並持續改進
4. **定期更新**：保持 LINE Bot API 和相關依賴的更新

## 10. 常見問題

### Q: Webhook 驗證失敗

A: 檢查 Channel Secret 是否正確，確保簽名計算正確

### Q: 無法發送推播訊息

A: 確認 Channel Access Token 有效，檢查用戶是否已加好友

### Q: 機器人無回應

A: 檢查 Webhook URL 是否可訪問，查看伺服器日誌

### Q: 訊息格式錯誤

A: 確保 JSON 格式正確，檢查必填欄位

---

這個指南提供了完整的 LINE Bot 設置流程，你可以根據實際需求進行調整和擴展。記住要測試所有功能並確保安全性。
