# Go 專案 Makefile

.PHONY: help build run clean migrate seed test

# 預設目標
help:
	@echo "可用的命令:"
	@echo "  build     - 建構應用程式"
	@echo "  run       - 運行應用程式"
	@echo "  clean     - 清理建構檔案"
	@echo "  migrate   - 執行資料庫遷移"
	@echo "  seed      - 執行資料庫 seeder"
	@echo "  test      - 執行測試"
	@echo "  docker-build - 建構 Docker 映像"
	@echo "  docker-run   - 運行 Docker 容器"

# 建構應用程式
build:
	go build -o bin/main .

# 運行應用程式
run:
	go run .

# 清理建構檔案
clean:
	rm -rf bin/

# 執行資料庫遷移
migrate:
	go run cmd/migrate/main.go

# 執行資料庫 seeder
seed:
	go run cmd/seed/main.go run

# 清除測試數據
seed-clear:
	go run cmd/seed/main.go clear

# 只執行用戶 seeder
seed-user:
	go run cmd/seed/main.go user

# 執行測試
test:
	go test ./...

# Docker 相關命令
docker-build:
	docker-compose build

docker-run:
	docker-compose up -d

docker-logs:
	docker-compose logs -f

docker-stop:
	docker-compose down

# 在 Docker 容器中執行命令
docker-migrate:
	docker-compose exec go-app ./migrate

docker-seed:
	docker-compose exec go-app ./seed run

docker-seed-clear:
	docker-compose exec go-app ./seed clear

docker-seed-user:
	docker-compose exec go-app ./seed user

# 進入 Docker 容器
docker-shell:
	docker-compose exec go-app /bin/sh
