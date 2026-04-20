# GitHub Actions 流水线深度强化计划 (GitHub Workflow Hardening Plan)

**Goal:** 通过引入依赖缓存、安全审计、并发控制以及标准化的元数据管理，全面提升 GitHub Actions 的执行效率、安全性和配置质量。

---

### Task 1: 强化 CI 基础检查 (ci.yml)

**Files:**
- Modify: `.github/workflows/ci.yml`

**Changes:**
1. **并发控制**: 增加 `concurrency` 配置，同一分支的新推送将自动取消正在运行的旧任务。
2. **依赖缓存**: 
   - 为 `actions/setup-go` 开启 `cache: true`。
   - 为 `actions/setup-node` 开启 `cache: 'npm'`。
3. **安全审计**: 在测试前增加 `golang.org/x/vuln/cmd/govulncheck@latest` 扫描步骤，检测依赖漏洞。

### Task 2: 规范化 Docker 发布流水线 (Main & Tags)

**Files:**
- Modify: `.github/workflows/docker-publish-main.yml`
- Modify: `.github/workflows/docker-publish-tags.yml`

**Changes:**
1. **元数据管理**: 引入 `docker/metadata-action@v5`。
   - 移除手动计算 `TAG_NAME` 的 Bash 脚本。
   - 自动根据分支（latest）或标签（v*）生成规范的 Docker 标签和 Labels。
2. **并发控制**: 增加 `concurrency` 配置。

### Task 3: 优化 Release 发布流水线 (release.yml)

**Files:**
- Modify: `.github/workflows/release.yml`

**Changes:**
1. **并发控制**: 增加 `concurrency` 配置。
2. **构建优化**: 确保 Node 缓存路径与 `web/` 目录匹配，减少前端编译时间。

### Task 4: 提交与验证

- [ ] 提交所有流水线优化代码。
- [ ] 观察 GitHub Actions 运行耗时，验证缓存是否生效。
