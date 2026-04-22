# --- 第一阶段：构建前端 (固定使用构建机原生架构以提速) ---
FROM --platform=$BUILDPLATFORM node:20-alpine AS web-builder
WORKDIR /app/web

# 安装依赖 (利用 Docker 缓存)
COPY web/package*.json ./
RUN npm install

# 拷贝源码并构建
COPY web/ .
RUN npm run build

# --- 第二阶段：构建后端 (固定使用原生架构，通过环境变量交叉编译) ---
FROM --platform=$BUILDPLATFORM golang:1.25-alpine AS server-builder

WORKDIR /app

# 下载依赖 (利用 Docker 缓存)
COPY go.mod go.sum ./
RUN go mod download

# 接收 Buildx 自动注入的跨平台参数
ARG TARGETOS
ARG TARGETARCH

# 拷贝源码并编译
COPY . .
COPY --from=web-builder /app/web/dist ./internal/api/dist
ARG VERSION=latest
# 设置 CGO_ENABLED=0 并使用目标架构参数进行极速交叉编译
RUN CGO_ENABLED=0 GOOS=${TARGETOS} GOARCH=${TARGETARCH} go build -tags embed -ldflags="-s -w -X main.version=${VERSION}" -o ucas cmd/server/main.go

# --- 第三阶段：最终镜像 (适配目标架构) ---
FROM alpine:latest
# 安装运行环境依赖
RUN apk add --no-cache ca-certificates tzdata

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
