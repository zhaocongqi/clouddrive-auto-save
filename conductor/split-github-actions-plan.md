# 拆分 GitHub Actions 流水线计划 (Split GitHub Actions Plan)

**Goal:** 将单一的 Docker 发布流水线拆分为两个独立的文件，以解决 GitHub Actions 原生不支持在 `tags` 推送时应用
`paths-ignore` 过滤规则的问题。实现在更新文档时不触发主分支的镜像构建。

---

## Task 1: 删除旧的单一流水线

**Files:**

- Delete: `.github/workflows/docker-publish.yml`

### Task 2: 创建主分支发布流水线

**Files:**

- Create: `.github/workflows/docker-publish-main.yml`

**Changes:**

- 触发条件: 仅 `push` 到 `main` 分支。
- 过滤条件: 使用 `paths-ignore` 忽略 `docs/**`, `conductor/**`, `README.md`,
  `CHANGELOG.md`, `LICENSE`, `.gitignore`。
- 构建逻辑: 直接构建并推送 `latest` 标签。

### Task 3: 创建 Tags 发布流水线

**Files:**

- Create: `.github/workflows/docker-publish-tags.yml`

**Changes:**

- 触发条件: 仅 `push` 到匹配 `v*` 的 `tags`。
- 构建逻辑: 提取标签名称（如 `v1.0.0`），构建并推送该版本的镜像。

### Task 4: 提交与推送

- [ ] 提交所有流水线修改。
- [ ] 推送到远程仓库，以验证过滤功能生效。
