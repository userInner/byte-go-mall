

# 使用官方的 Golang 镜像作为构建环境
FROM golang:1.23.5-alpine AS builder

# 设置工作目录
WORKDIR /app

ENV GOPROXY=https://goproxy.cn,direct
# 复制 go.mod 和 go.sum 文件
COPY go.mod go.sum ./

# 下载依赖
RUN go mod download

# 复制项目文件（注意：这里复制的是整个项目）
COPY . .

# 切换到 app/user 目录
WORKDIR /app/app/user

# 构建项目
RUN CGO_ENABLED=0 GOOS=linux go build -o user-service .

# 使用 Alpine 作为运行时环境
FROM alpine:3.19

# 设置工作目录
WORKDIR /app

# 从构建阶段复制可执行文件
COPY --from=builder /app/app/user/user-service .

# 运行服务
CMD ["./user-service"]