#!/bin/bash

# 獲取 Git 資訊
GIT_COMMIT=$(git rev-parse HEAD)
GIT_BRANCH=$(git rev-parse --abbrev-ref HEAD)
BUILD_TIME=$(date -u +"%Y-%m-%dT%H:%M:%SZ")

echo "構建資訊："
echo "Git Commit: $GIT_COMMIT"
echo "Git Branch: $GIT_BRANCH"
echo "Build Time: $BUILD_TIME"

# 構建 Docker 映像並傳入 Git 資訊
docker-compose build --build-arg GIT_COMMIT="$GIT_COMMIT" --build-arg GIT_BRANCH="$GIT_BRANCH" --build-arg BUILD_TIME="$BUILD_TIME"

# 啟動容器
docker-compose up -d

echo "構建完成！版本號應該會顯示為: 2.0.0-${GIT_COMMIT:0:7}"
