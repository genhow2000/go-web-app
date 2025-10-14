# LINE Bot 故障排除指南

## 問題：掃描 QR Code 後顯示「無法加入好友，請確定是否正確」

### 🔍 可能的原因

1. **LINE Bot 未公開**
2. **機器人 ID 不正確**
3. **機器人未啟用加入好友功能**
4. **機器人設置不完整**

### ✅ 解決步驟

#### 1. 檢查 LINE Bot 設置

1. **前往 LINE Developers Console**

   - 網址：https://developers.line.biz/
   - 使用你的 LINE 帳號登入

2. **檢查你的 Bot Channel**

   - 確認 Channel 已創建
   - 確認 Channel 狀態為「Active」

3. **檢查 Messaging API 設置**
   - 在 Channel 設定頁面
   - 確認「Use webhook」已啟用
   - 確認機器人已公開（不是私人）

#### 2. 驗證機器人 ID

你的機器人 ID：`@351thdpd`

**測試方法**：

1. 在 LINE 中直接搜尋 `@351thdpd`
2. 如果找不到，表示機器人未公開或 ID 不正確

#### 3. 檢查機器人權限

在 LINE Developers Console 中確認：

- ✅ 機器人已公開
- ✅ 允許用戶加入好友
- ✅ Webhook 已設置
- ✅ Channel Access Token 有效

#### 4. 替代測試方法

如果 QR Code 不工作，可以嘗試：

1. **直接搜尋**：在 LINE 中搜尋 `@351thdpd`
2. **分享連結**：使用 `https://line.me/R/ti/p/351thdpd`
3. **手動添加**：提供機器人的 QR Code 圖片

### 🛠️ 修復建議

#### 如果機器人未公開：

1. 前往 LINE Developers Console
2. 選擇你的機器人 Channel
3. 在「Messaging API」設定中：
   - 確保「Use webhook」已啟用
   - 確保機器人狀態為「Active」
   - 檢查是否有任何錯誤訊息

#### 如果機器人 ID 不正確：

1. 在 LINE Developers Console 中確認正確的機器人 ID
2. 更新代碼中的機器人 ID
3. 重新生成 QR Code

#### 如果機器人設置不完整：

1. 完成所有必要的設置步驟
2. 確保 Webhook URL 已配置
3. 測試機器人回應功能

### 📱 測試步驟

1. **重新整理網頁**：http://localhost:3001
2. **點擊「立即加入」按鈕**
3. **嘗試不同的加入方式**：
   - 掃描 QR Code
   - 在 LINE 中搜尋 @351thdpd
   - 點擊「說明」按鈕查看詳細指導

### 🔧 代碼修復

如果問題持續，可能需要：

1. **更新機器人 ID**：確認 `@351thdpd` 是否正確
2. **檢查 URL 格式**：確保使用正確的 LINE Bot URL 格式
3. **添加錯誤處理**：提供更好的用戶反饋

### 📞 需要協助？

如果以上步驟都無法解決問題，請提供：

1. LINE Developers Console 的截圖
2. 機器人 Channel 的詳細設置
3. 具體的錯誤訊息
4. 測試的設備和瀏覽器資訊

---

**最後更新**：2024 年 12 月
**機器人 ID**：@351thdpd
**狀態**：需要驗證設置
