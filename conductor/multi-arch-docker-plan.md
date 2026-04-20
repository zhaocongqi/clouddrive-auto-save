# 异构镜像支持与构建流水线优化计划 (Multi-Arch Docker Pipeline Plan)

**Goal:** 优化 GitHub Actions 工作流，利用 Docker 官方的 `buildx` 与 `qemu`
功能支持构建并发布多架构（异构）镜像（如 `linux/amd64` 和 `linux/arm64`），全面提升容器镜像的兼容性。同时优化 Dockerfile
以支持将外部版本号注入二进制文件。

---

## Task 1: 优化 Dockerfile 接收版本变量

**Files:**

- Modify: `Dockerfile`

**Changes:**

1. 在编译阶段引入 `ARG VERSION=latest`。
2. 将该参数通过 `-ldflags="-X main.version=${VERSION}"` 注入编译步骤，确保 Docker
   容器内运行的程序也能打印正确的版本号。

### Task 2: 改造主分支流水线 (Main Branch)

**Files:**

- Modify: `.github/workflows/docker-publish-main.yml`

**Changes:**

1. 引入 `docker/setup-qemu-action@v3` 以提供跨架构仿真环境。
2. 引入 `docker/setup-buildx-action@v3` 创建并引导多平台构建器。
3. 摒弃 `Makefile` 在 CI 中的调用，直接使用业界标准的 `docker/build-push-action@v5`：
   - 启用 `platforms: linux/amd64,linux/arm64`
   - 传递 `build-args: VERSION=latest`
   - 自动执行推送。

### Task 3: 改造标签发布流水线 (Tags)

**Files:**

- Modify: `.github/workflows/docker-publish-tags.yml`

**Changes:**

1. 同样引入 QEMU 和 Buildx 运行环境。
2. 使用 `docker/build-push-action@v5` 进行多平台构建。
3. 传递 `build-args: VERSION=${{ env.TAG_NAME }}`。
4. **双标签推送**：打 tag 发布时，不仅推送 `${TAG_NAME}` 标签（如 `v1.0.0`），还会同步更新 `latest`
   标签，这是开源项目的最佳实践。

### Task 4: 提交与推送

- [ ] 提交修改到主分支。
- [ ] 自动执行推送以验证新流水线能否正确编译并上传异构镜像。
