#!/bin/bash

# 優化暫存變更腳本
echo "🔧 優化暫存變更..."

# 1. 移除重複的腳本文件
echo "📁 移除重複的腳本文件..."
rm -f start-local.sh
rm -f test-vue-docker.sh

# 2. 恢復有用的文檔文件
echo "📄 恢復有用的文檔文件..."
git restore UNIT_TEST_PLAN.md
git restore USER_MANAGEMENT.md
git restore VUE_MIGRATION_TEST_RESULTS.md

# 3. 恢復有用的腳本文件
echo "🔧 恢復有用的腳本文件..."
git restore setup-env.sh
git restore verify_docker.sh
git restore verify_system.sh

echo "✅ 優化完成！"
echo ""
echo "📋 優化內容："
echo "  - 移除了重複的腳本文件"
echo "  - 恢復了有用的文檔文件"
echo "  - 恢復了有用的腳本文件"
echo ""
echo "💡 現在可以提交變更："
echo "  git add ."
echo "  git commit -m 'feat: 改進開發環境和用戶認證功能'"
echo "  git push"
