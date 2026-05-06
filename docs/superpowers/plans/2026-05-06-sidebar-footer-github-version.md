# 侧边栏底部 GitHub 链接 + 版本号 + 更新提示 实现计划

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** 在侧边栏底部添加 GitHub 仓库链接和版本号，有新版本时红点提示

**Architecture:** 修复后端 ldflags 版本注入 + 新增 `/api/version` API 端点 + 前端 `SidebarFooter` 组件通过 GitHub API 检查最新 release 并对比版本

**Tech Stack:** Go (Gin), Vue 3, Element Plus, lucide-vue-next, axios

---

## 文件清单

| 文件 | 操作 | 职责 |
|------|------|------|
| `cmd/server/main.go` | 修改 | 声明 version/commit/date 变量，传递给 InitRouter |
| `internal/api/router.go` | 修改 | InitRouter 接收版本参数 + 新增 `/api/version` 路由和 handler |
| `internal/api/router_test.go` | 修改 | 更新 setupTestRouter 调用签名 |
| `web/src/api/version.js` | 新建 | 版本 API 请求函数 |
| `web/src/components/SidebarFooter.vue` | 新建 | 侧边栏底部组件：GitHub 链接 + 版本号 + 更新检查 |
| `web/src/layout/MainLayout.vue` | 修改 | 引入 SidebarFooter 组件 |

---

### Task 1: 修复后端版本注入 + 新增版本 API

**Files:**
- Modify: `cmd/server/main.go`
- Modify: `internal/api/router.go`
- Modify: `internal/api/router_test.go`

- [ ] **Step 1: 在 main.go 声明版本变量**

在 `cmd/server/main.go` 的 `import` 和 `func main()` 之间添加版本变量声明：

```go
package main

import (
	"log/slog"
	"os"
	"strings"

	"github.com/zcq/clouddrive-auto-save/internal/api"
	"github.com/zcq/clouddrive-auto-save/internal/core"
	"github.com/zcq/clouddrive-auto-save/internal/core/scheduler"
	"github.com/zcq/clouddrive-auto-save/internal/core/worker"
	"github.com/zcq/clouddrive-auto-save/internal/db"
	"github.com/zcq/clouddrive-auto-save/internal/utils"
)

// 版本信息，构建时通过 -ldflags 注入
var (
	version = "dev"
	commit  = "unknown"
	date    = "unknown"
)

func main() {
```

- [ ] **Step 2: 添加启动版本日志**

在 `cmd/server/main.go` 的 `utils.InitLogger(...)` 之后添加：

```go
	slog.Info("UCAS starting", "version", version, "commit", commit, "date", date)
```

- [ ] **Step 3: 修改 InitRouter 调用，传递版本信息**

将 `cmd/server/main.go` 中的：
```go
	r := api.InitRouter(wm)
```
改为：
```go
	r := api.InitRouter(wm, version, commit, date)
```

- [ ] **Step 4: 修改 router.go — InitRouter 签名 + 新增版本端点**

在 `internal/api/router.go` 中，修改 `InitRouter` 函数签名，新增版本变量和 handler：

```go
// 版本信息（由 main 包通过 InitRouter 传入）
var (
	appVersion = "dev"
	appCommit  = "unknown"
	appDate    = "unknown"
)

func InitRouter(wm *worker.Manager, version, commit, date string) *gin.Engine {
	WorkerManager = wm
	appVersion = version
	appCommit = commit
	appDate = date
	r := gin.Default()

	// 基础 API 路由组
	api := r.Group("/api")
	{
		api.GET("/version", getVersion)

		// ... existing routes unchanged ...
	}
	// ... rest unchanged ...
}
```

在 `router.go` 文件末尾（`testBarkNotification` 函数之后）添加 handler：

```go
func getVersion(c *gin.Context) {
	c.PureJSON(http.StatusOK, gin.H{
		"version": appVersion,
		"commit":  appCommit,
		"date":    appDate,
	})
}
```

- [ ] **Step 5: 更新 router_test.go 调用签名**

将 `internal/api/router_test.go` 中的：
```go
	r := InitRouter(wm)
```
改为：
```go
	r := InitRouter(wm, "test", "test", "test")
```

- [ ] **Step 6: 验证编译和测试**

Run: `cd /home/zcq/Github/clouddrive-auto-save && go vet ./... && go test ./internal/api/... -v -count=1`
Expected: 编译通过，测试全部 PASS

- [ ] **Step 7: 手动验证版本端点**

Run: `cd /home/zcq/Github/clouddrive-auto-save && go run cmd/server/main.go &` (后台启动)
然后: `curl http://localhost:8080/api/version`
Expected: `{"commit":"unknown","date":"unknown","version":"dev"}`
清理: `kill %1`

- [ ] **Step 8: Commit**

```bash
git add cmd/server/main.go internal/api/router.go internal/api/router_test.go
git commit -m "feat(api): 修复版本注入并新增 /api/version 端点

- 在 main.go 声明 version/commit/date 变量（ldflags 注入目标）
- 启动日志输出版本信息
- InitRouter 接收版本参数
- 新增 GET /api/version 端点返回当前版本信息"
```

---

### Task 2: 前端版本 API 请求函数

**Files:**
- Create: `web/src/api/version.js`

- [ ] **Step 1: 创建版本 API 请求模块**

创建 `web/src/api/version.js`：

```js
import request from './request'

export function getVersion() {
  return request.get('/version')
}
```

- [ ] **Step 2: Commit**

```bash
git add web/src/api/version.js
git commit -m "feat(web): 添加版本 API 请求函数"
```

---

### Task 3: 创建 SidebarFooter 组件

**Files:**
- Create: `web/src/components/SidebarFooter.vue`

- [ ] **Step 1: 创建 SidebarFooter 组件**

创建 `web/src/components/SidebarFooter.vue`：

```vue
<template>
  <div class="sidebar-footer">
    <a
      class="github-link"
      href="https://github.com/zhaocongqi/clouddrive-auto-save"
      target="_blank"
      rel="noopener noreferrer"
    >
      <Github :size="16" />
      <span>GitHub 仓库</span>
      <ExternalLink :size="12" />
    </a>
    <div
      class="version-info"
      :class="{ 'has-update': hasUpdate }"
      @click="hasUpdate && openReleases()"
    >
      <span class="version-text">v{{ currentVersion }}</span>
      <span v-if="hasUpdate" class="update-badge">NEW</span>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { Github, ExternalLink } from 'lucide-vue-next'
import { getVersion } from '../api/version'

const currentVersion = ref('...')
const hasUpdate = ref(false)

const GITHUB_RELEASES_URL = 'https://github.com/zhaocongqi/clouddrive-auto-save/releases'
const GITHUB_API_URL = 'https://api.github.com/repos/zhaocongqi/clouddrive-auto-save/releases/latest'

function parseVersion(v) {
  const cleaned = v.replace(/^v/, '')
  const parts = cleaned.split('.').map(Number)
  return parts.length === 3 && parts.every(n => !isNaN(n)) ? parts : null
}

function compareVersions(current, latest) {
  const c = parseVersion(current)
  const l = parseVersion(latest)
  if (!c || !l) return false
  for (let i = 0; i < 3; i++) {
    if (l[i] > c[i]) return true
    if (l[i] < c[i]) return false
  }
  return false
}

function openReleases() {
  window.open(GITHUB_RELEASES_URL, '_blank', 'noopener,noreferrer')
}

onMounted(async () => {
  // 获取当前版本
  try {
    const res = await getVersion()
    currentVersion.value = res.version || 'dev'
  } catch {
    currentVersion.value = 'dev'
  }

  // 如果是 dev 版本，跳过更新检查
  if (currentVersion.value === 'dev') return

  // 检查 GitHub 最新 release
  try {
    const resp = await fetch(GITHUB_API_URL)
    if (!resp.ok) return
    const data = await resp.json()
    const latestTag = data.tag_name || ''
    hasUpdate.value = compareVersions(currentVersion.value, latestTag)
  } catch {
    // 静默失败
  }
})
</script>

<style scoped>
.sidebar-footer {
  margin-top: auto;
  padding: 12px 16px;
  border-top: 1px solid rgba(0, 0, 0, 0.05);
}

html.dark .sidebar-footer {
  border-top: 1px solid rgba(255, 255, 255, 0.05);
}

.github-link {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 8px 12px;
  border-radius: 8px;
  color: #64748b;
  text-decoration: none;
  font-size: 13px;
  transition: background 0.2s, color 0.2s;
}

.github-link:hover {
  background: #f1f5f9;
  color: #334155;
}

html.dark .github-link {
  color: #94a3b8;
}

html.dark .github-link:hover {
  background: rgba(255, 255, 255, 0.05);
  color: #e2e8f0;
}

.version-info {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 4px 12px;
  font-size: 12px;
  color: #94a3b8;
}

.version-info.has-update {
  cursor: pointer;
  color: #64748b;
  border-radius: 8px;
  padding: 6px 12px;
  margin: 2px 0;
  transition: background 0.2s;
}

.version-info.has-update:hover {
  background: #f1f5f9;
}

html.dark .version-info.has-update:hover {
  background: rgba(255, 255, 255, 0.05);
}

.update-badge {
  display: inline-flex;
  align-items: center;
  padding: 1px 6px;
  border-radius: 10px;
  background: #ef4444;
  color: white;
  font-size: 10px;
  font-weight: 600;
  letter-spacing: 0.5px;
  animation: pulse 2s ease-in-out infinite;
}

@keyframes pulse {
  0%, 100% { opacity: 1; }
  50% { opacity: 0.7; }
}
</style>
```

- [ ] **Step 2: 验证前端编译**

Run: `cd /home/zcq/Github/clouddrive-auto-save/web && npm run build`
Expected: 构建成功，无错误

- [ ] **Step 3: Commit**

```bash
git add web/src/components/SidebarFooter.vue
git commit -m "feat(web): 新增 SidebarFooter 组件

- GitHub 仓库链接（可点击跳转）
- 当前版本号显示
- 通过 GitHub API 检查最新 release
- 有新版本时显示红点 + NEW 标签，点击跳转 Releases 页
- 暗色模式适配"
```

---

### Task 4: 集成 SidebarFooter 到 MainLayout

**Files:**
- Modify: `web/src/layout/MainLayout.vue`

- [ ] **Step 1: 在 MainLayout 中引入并使用 SidebarFooter**

在 `<script setup>` 中添加 import（在现有 import 之后）：

```js
import SidebarFooter from '../components/SidebarFooter.vue'
```

在模板中，将 `</el-menu>` 之后的 `</el-aside>` 之前插入组件：

```vue
      </el-menu>
      <SidebarFooter />
    </el-aside>
```

- [ ] **Step 2: 验证前端开发服务器**

Run: `cd /home/zcq/Github/clouddrive-auto-save/web && npm run dev &` (后台启动)
在浏览器访问 `http://localhost:5173`，检查：
- 侧边栏底部显示 GitHub 链接和版本号
- 点击 GitHub 链接可跳转
- 暗色模式下样式正常

- [ ] **Step 3: Commit**

```bash
git add web/src/layout/MainLayout.vue
git commit -m "feat(web): 在侧边栏底部集成 SidebarFooter 组件"
```

---

### Task 5: 端到端验证

- [ ] **Step 1: 完整构建验证**

Run: `cd /home/zcq/Github/clouddrive-auto-save && make build`
Expected: 构建成功，`bin/ucas` 生成

- [ ] **Step 2: 验证构建版本号注入**

Run: `cd /home/zcq/Github/clouddrive-auto-save && VERSION=1.1.0 make build-server`
然后: `curl http://localhost:8080/api/version` (启动后)
Expected: `{"version":"1.1.0","commit":"<hash>","date":"<date>"}`

- [ ] **Step 3: 运行完整测试**

Run: `cd /home/zcq/Github/clouddrive-auto-save && make check`
Expected: lint + vet + test 全部通过

- [ ] **Step 4: 最终 Commit（如有遗漏）**

```bash
git status
# 确认没有未提交的变更
```
