# CI 流水线优化计划 (CI Workflow Optimization Plan)

**Goal:** 优化 `.github/workflows/ci.yml`，使其只在后端 Go 代码及相关配置文件发生变动时才触发 `make check`，从而避免在仅修改文档或前端代码时浪费 CI 资源。

---

### Task 1: 优化 CI 触发条件

**Files:**
- Modify: `.github/workflows/ci.yml`

**Changes:**
1. 为 `push` 和 `pull_request` 添加 `paths` 过滤规则。
2. 仅当以下文件或目录发生变动时触发流水线：
   - `**.go` (Go 源码)
   - `go.mod` (依赖声明)
   - `go.sum` (依赖校验)
   - `Makefile` (构建脚本)
   - `.github/workflows/ci.yml` (流水线配置本身)

### Task 2: 自动提交

- [ ] 本地提交修改，不自动推送。