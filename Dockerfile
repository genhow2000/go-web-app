# 使用官方 Go 映像作為構建階段
FROM golang:1.21-alpine AS builder

# 設置工作目錄
WORKDIR /app

# 複製 go mod 文件
COPY go.mod ./

# 下載依賴並生成 go.sum
RUN go mod download && go mod tidy

# 複製源代碼和模板
COPY . .

# 整理依賴並構建應用（使用純Go SQLite驅動）
RUN go mod tidy && CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -o main .

# 構建遷移工具
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -o migrate cmd/migrate/main.go

# 構建管理員初始化工具
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -o init-admin cmd/init-admin/main.go

# 使用輕量級的 alpine 映像作為運行階段
FROM alpine:latest

# 安裝 ca-certificates 用於 HTTPS 請求
RUN apk --no-cache add ca-certificates

# 設置工作目錄
WORKDIR /root/

# 創建資料庫目錄 (統一使用 /tmp)
RUN mkdir -p /tmp

# 從構建階段複製二進制文件和模板
COPY --from=builder /app/main .
COPY --from=builder /app/migrate .
COPY --from=builder /app/init-admin .
COPY --from=builder /app/templates ./templates
COPY --from=builder /app/migrations ./migrations

# 暴露端口
EXPOSE 8080

# 運行應用
CMD ["./main"]
