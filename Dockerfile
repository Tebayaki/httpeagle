# 阶段1: 构建阶段
FROM golang:1.21-alpine AS builder

# 设置工作目录
WORKDIR /app

# 复制go.mod（如果使用依赖管理）
COPY go.mod ./

# 下载依赖（如果有）
RUN go mod download || true

# 复制源码
COPY main.go .

# 编译为静态二进制文件
# CGO_ENABLED=0: 禁用CGO，生成纯静态二进制
# GOOS=linux: 目标操作系统
# GOARCH=amd64: 目标架构（可选arm64）
# -ldflags="-s -w": 去除调试信息，减小体积
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o httpeagle main.go

# 安装UPX并压缩二进制
RUN apk add --no-cache upx && \
    upx --best --lzma /app/httpeagle

# 阶段2: 运行阶段（最小镜像）
FROM alpine:latest

# 安装CA证书（用于HTTPS）
RUN apk --no-cache add ca-certificates

# 创建非root用户运行
RUN addgroup -g 1000 -S appgroup && \
    adduser -u 1000 -S appuser -G appgroup

# 创建工作目录
RUN mkdir -p /images /certs

# 从构建阶段复制二进制文件
COPY --from=builder /app/httpeagle /usr/local/bin/

# 修改文件所有权
RUN chown -R appuser:appgroup /images /certs

# 切换为非root用户
USER appuser

# 暴露端口
EXPOSE 41596

# 健康检查
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:41596/health || exit 1

# 启动服务
CMD ["httpeagle"]