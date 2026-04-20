# 自动发布 GitHub Releases 产物计划 (GoReleaser Pipeline Plan)

**Goal:** 创建一个专用的 GitHub Action 流水线，当推送 Tag (如 `v1.0.0`) 时，自动调用 GoReleaser
工具，将交叉编译出的各个平台二进制文件（Windows, Linux, macOS）打包并作为附件发布到 GitHub Releases 页面。

---

## Task 1: 创建 Release 工作流

**Files:**

- Create: `.github/workflows/release.yml`

**Changes:**

1. **触发条件**: 监听 `v*` 格式的 Tag 推送。
2. **权限配置**: 授予 `contents: write` 权限，以允许 Action 在代码库中创建 Release 页面并上传附件。
3. **环境准备**:
   - `actions/checkout@v4`，需配置 `fetch-depth: 0` 以确保 GoReleaser 能读取完整的 Git
     历史来生成准确的变更日志 (Changelog)。
   - `actions/setup-go@v5` 配置 Go 1.25 环境。
   - `actions/setup-node@v4` 配置 Node.js 20 环境（**关键步骤**：因为后端启用了 `-tags=embed`，必须先编译
     Vue 前端）。
4. **编译前端**: 运行 `make build-web` 生成 `web/dist` 静态资源。
5. **执行发布**: 使用官方 `goreleaser/goreleaser-action@v5` 运行 `release --clean` 命令，自动读取
   `.goreleaser.yaml` 配置，完成交叉编译与 GitHub Releases 发布。

### Task 2: 提交与推送

- [ ] 提交新增的工作流文件。
- [ ] 推送至远程仓库。
