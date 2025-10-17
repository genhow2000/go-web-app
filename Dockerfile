# 使用官方 Go 映像作為構建階段
FROM golang:1.21-alpine AS builder

# 安裝必要的編譯工具和 Git
RUN apk add --no-cache gcc musl-dev git

# 設置工作目錄
WORKDIR /app

# 複製 go mod 文件
COPY go.mod ./

# 下載依賴並生成 go.sum
RUN go mod download && go mod tidy

# 複製源代碼和模板
COPY . .

# 獲取 Git 資訊並設置構建參數
RUN git config --global --add safe.directory /app
ARG GIT_COMMIT=""
ARG GIT_BRANCH=""
ARG BUILD_TIME=""

# 整理依賴並構建應用（使用純Go SQLite驅動）
RUN go mod tidy && CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags "-X main.GitCommit=${GIT_COMMIT} -X main.GitBranch=${GIT_BRANCH} -X main.BuildTime=${BUILD_TIME}" -o main .

# 構建 migration 工具
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o migrate cmd/migrate/main.go

# 構建 seeder 工具
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o seed cmd/seed/main.go

# Vue.js 前端構建階段
FROM node:18-alpine AS frontend-builder

WORKDIR /app/frontend

# 複製前端依賴文件
COPY frontend/package*.json ./

# 安裝前端依賴
RUN npm install

# 複製前端源代碼
COPY frontend/ .

# 構建前端應用
RUN npm run build

# 使用輕量級的 alpine 映像作為運行階段
FROM alpine:latest

# 安裝 ca-certificates 用於 HTTPS 請求
RUN apk --no-cache add ca-certificates

# 設置工作目錄
WORKDIR /root/

# 創建資料庫目錄
RUN mkdir -p /tmp && mkdir -p /app/data

# 從構建階段複製二進制文件和模板
COPY --from=builder /app/main .
COPY --from=builder /app/migrate .
COPY --from=builder /app/seed .
COPY --from=builder /app/templates ./templates
COPY --from=builder /app/md ./md
COPY --from=builder /app/migrations ./migrations

# 從前端構建階段複製構建文件
COPY --from=frontend-builder /app/frontend/dist ./static/dist

# 暴露端口
EXPOSE 8080

# 運行應用
CMD ["./main"]
