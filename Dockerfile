# --- 第一阶段：构建前端 ---
FROM node:20-alpine AS web-builder
WORKDIR /app/web

# 设置 npm 镜像源并安装依赖 (利用 Docker 缓存)
RUN npm config set registry https://registry.npmmirror.com
COPY web/package*.json ./
RUN npm install

# 拷贝源码并构建
COPY web/ .
RUN npm run build

# --- 第二阶段：构建后端 ---
FROM golang:1.25-alpine AS server-builder
# 替换 Alpine 镜像源加速
RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.tuna.tsinghua.edu.cn/g' /etc/apk/repositories

WORKDIR /app

# 设置 Go 代理并下载依赖 (利用 Docker 缓存)
ENV GOPROXY=https://goproxy.cn,direct
COPY go.mod go.sum ./
RUN go mod download

# 拷贝源码并编译
COPY . .
COPY --from=web-builder /app/web/dist ./internal/api/dist
RUN go build -tags embed -ldflags="-s -w" -o ucas cmd/server/main.go

# --- 第三阶段：最终镜像 ---
FROM alpine:latest
# 替换 Alpine 镜像源加速 apk add
RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.tuna.tsinghua.edu.cn/g' /etc/apk/repositories && \
    apk add --no-cache ca-certificates tzdata

WORKDIR /app
COPY --from=server-builder /app/ucas .
COPY --from=web-builder /app/web/dist ./web/dist

# 配置环境变量默认值
ENV LOG_LEVEL=INFO
ENV DB_PATH=/app/data/data.db
ENV TZ=Asia/Shanghai

# 创建数据目录
RUN mkdir -p /app/data

EXPOSE 8080
CMD ["./ucas"]
