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

# 整理依賴並構建應用
RUN go mod tidy && CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

# 使用輕量級的 alpine 映像作為運行階段
FROM alpine:latest

# 安裝 ca-certificates 用於 HTTPS 請求
RUN apk --no-cache add ca-certificates

# 設置工作目錄
WORKDIR /root/

# 從構建階段複製二進制文件和模板
COPY --from=builder /app/main .
COPY --from=builder /app/templates ./templates

# 暴露端口
EXPOSE 8080

# 運行應用
CMD ["./main"]
