#!/bin/bash

# Cloud Run URL 監控和 LINE Console 自動更新腳本
# 當 Cloud Run URL 變動時，自動更新 LINE Console 的 redirect_uri

set -e

# 配置變數
PROJECT_ID="fleet-day-383710"
SERVICE_NAME="go-app"
REGION="asia-east1"
LINE_CHANNEL_ID="2008159551"
LINE_CHANNEL_SECRET="2cca495d6b53e8b2a2d684ee87113f01"

# 檔案路徑
URL_FILE="/tmp/current_url.txt"
LOG_FILE="/tmp/url_monitor.log"

# 記錄日誌
log() {
    echo "$(date '+%Y-%m-%d %H:%M:%S') - $1" >> "$LOG_FILE"
}

# 獲取當前 Cloud Run URL
get_current_url() {
    gcloud run services describe "$SERVICE_NAME" \
        --region="$REGION" \
        --format="value(status.url)" 2>/dev/null || echo ""
}

# 獲取上次記錄的 URL
get_last_url() {
    if [ -f "$URL_FILE" ]; then
        cat "$URL_FILE"
    else
        echo ""
    fi
}

# 儲存當前 URL
save_current_url() {
    echo "$1" > "$URL_FILE"
}

# 更新 LINE Console redirect_uri
update_line_redirect_uri() {
    local new_url="$1"
    local redirect_uri="${new_url}/auth/line/callback"
    
    log "嘗試更新 LINE Console redirect_uri 為: $redirect_uri"
    
    # 使用 LINE Management API 更新 redirect_uri
    # 注意：這需要 LINE Channel Access Token，需要先獲取
    local channel_access_token=$(get_line_channel_access_token)
    
    if [ -n "$channel_access_token" ]; then
        local response=$(curl -s -X PUT "https://api.line.me/v2/oauth/redirectUri" \
            -H "Authorization: Bearer $channel_access_token" \
            -H "Content-Type: application/json" \
            -d "{\"redirectUri\": \"$redirect_uri\"}" 2>/dev/null)
        
        if echo "$response" | grep -q "success\|200"; then
            log "✅ LINE Console redirect_uri 更新成功: $redirect_uri"
            return 0
        else
            log "❌ LINE Console redirect_uri 更新失敗: $response"
            return 1
        fi
    else
        log "❌ 無法獲取 LINE Channel Access Token"
        return 1
    fi
}

# 獲取 LINE Channel Access Token
get_line_channel_access_token() {
    # 使用 Channel ID 和 Channel Secret 獲取 Access Token
    local response=$(curl -s -X POST "https://api.line.me/v2/oauth/accessToken" \
        -H "Content-Type: application/x-www-form-urlencoded" \
        -d "grant_type=client_credentials&client_id=$LINE_CHANNEL_ID&client_secret=$LINE_CHANNEL_SECRET" 2>/dev/null)
    
    echo "$response" | grep -o '"access_token":"[^"]*"' | cut -d'"' -f4
}

# 更新 GitHub Actions 環境變數
update_github_secrets() {
    local new_url="$1"
    
    log "更新 GitHub Secrets..."
    
    # 使用 GitHub CLI 更新 secrets
    if command -v gh &> /dev/null; then
        gh secret set SERVICE_URL --body "$new_url" 2>/dev/null || log "⚠️ 無法更新 GitHub Secrets"
    else
        log "⚠️ GitHub CLI 未安裝，請手動更新 GitHub Secrets"
    fi
}

# 更新 Cloud Build 配置
update_cloudbuild_config() {
    local new_url="$1"
    local cloudbuild_file="cloudbuild.yaml"
    
    log "更新 cloudbuild.yaml 配置..."
    
    # 備份原檔案
    cp "$cloudbuild_file" "${cloudbuild_file}.backup"
    
    # 更新 URL
    sed -i.bak "s|https://go-app-[^/]*\.a\.run\.app|$new_url|g" "$cloudbuild_file"
    
    # 檢查是否有變更
    if ! diff -q "$cloudbuild_file" "${cloudbuild_file}.backup" > /dev/null; then
        log "✅ cloudbuild.yaml 已更新"
        # 提交變更
        git add "$cloudbuild_file"
        git commit -m "自動更新 Cloud Run URL: $new_url" || true
        git push || log "⚠️ 無法推送變更到 Git"
    else
        log "ℹ️ cloudbuild.yaml 無需更新"
    fi
    
    # 清理備份檔案
    rm -f "${cloudbuild_file}.backup" "${cloudbuild_file}.bak"
}

# 主函數
main() {
    log "開始監控 Cloud Run URL..."
    
    # 獲取當前 URL
    local current_url=$(get_current_url)
    local last_url=$(get_last_url)
    
    if [ -z "$current_url" ]; then
        log "❌ 無法獲取 Cloud Run URL"
        exit 1
    fi
    
    log "當前 URL: $current_url"
    log "上次 URL: $last_url"
    
    # 檢查 URL 是否變動
    if [ "$current_url" != "$last_url" ]; then
        log "🔄 檢測到 URL 變動，開始更新配置..."
        
        # 更新 LINE Console
        if update_line_redirect_uri "$current_url"; then
            log "✅ LINE Console 更新成功"
        else
            log "❌ LINE Console 更新失敗，請手動檢查"
        fi
        
        # 更新 GitHub Secrets
        update_github_secrets "$current_url"
        
        # 更新 Cloud Build 配置
        update_cloudbuild_config "$current_url"
        
        # 儲存新的 URL
        save_current_url "$current_url"
        
        log "✅ URL 監控和更新完成"
    else
        log "ℹ️ URL 未變動，無需更新"
    fi
}

# 執行主函數
main "$@"
