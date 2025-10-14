#!/bin/bash

# 簡化版 Cloud Run URL 監控腳本
# 當 URL 變動時，自動更新配置文件並發送通知

set -e

# 配置變數
PROJECT_ID="fleet-day-383710"
SERVICE_NAME="go-app"
REGION="asia-east1"

# 檔案路徑
URL_FILE="/tmp/current_url.txt"
LOG_FILE="/tmp/url_monitor.log"

# 記錄日誌
log() {
    echo "$(date '+%Y-%m-%d %H:%M:%S') - $1" | tee -a "$LOG_FILE"
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

# 更新配置文件
update_configs() {
    local new_url="$1"
    local redirect_uri="${new_url}/auth/line/callback"
    
    log "更新配置文件..."
    
    # 更新 cloudbuild.yaml
    if [ -f "cloudbuild.yaml" ]; then
        sed -i.bak "s|https://go-app-[^/]*\.a\.run\.app|$new_url|g" cloudbuild.yaml
        log "✅ cloudbuild.yaml 已更新"
    fi
    
    # 更新 GitHub Actions
    if [ -f ".github/workflows/deploy.yml" ]; then
        sed -i.bak "s|https://go-app-[^/]*\.a\.run\.app|$new_url|g" .github/workflows/deploy.yml
        log "✅ GitHub Actions 已更新"
    fi
    
    # 更新 CORS 配置
    if [ -f "middleware/cors.go" ]; then
        sed -i.bak "s|https://go-app-[^/]*\.a\.run\.app|$new_url|g" middleware/cors.go
        log "✅ CORS 配置已更新"
    fi
    
    # 清理備份檔案
    find . -name "*.bak" -delete
    
    log "✅ 所有配置文件已更新"
    log "📝 請手動更新 LINE Console redirect_uri 為: $redirect_uri"
}

# 發送通知（可選）
send_notification() {
    local new_url="$1"
    local message="Cloud Run URL 已變動為: $new_url\n請更新 LINE Console redirect_uri 為: ${new_url}/auth/line/callback"
    
    # 如果有 webhook URL，可以發送通知
    if [ -n "$WEBHOOK_URL" ]; then
        curl -X POST "$WEBHOOK_URL" \
            -H "Content-Type: application/json" \
            -d "{\"text\": \"$message\"}" 2>/dev/null || true
    fi
    
    log "📢 通知已發送"
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
        
        # 更新配置文件
        update_configs "$current_url"
        
        # 發送通知
        send_notification "$current_url"
        
        # 儲存新的 URL
        save_current_url "$current_url"
        
        log "✅ URL 監控和更新完成"
        
        # 返回非零退出碼，表示有變更
        exit 1
    else
        log "ℹ️ URL 未變動，無需更新"
        exit 0
    fi
}

# 執行主函數
main "$@"
