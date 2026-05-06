# 侧边栏底部 GitHub 链接 + 版本号 + 更新提示

## 背景

用户希望在侧边栏底部展示 GitHub 仓库链接和当前版本号，并在有新版本时通过红点提示。

当前问题：
- 侧边栏 (`MainLayout.vue`) 没有底部区域
- Makefile 定义了 ldflags 版本注入，但 `cmd/server/main.go` 中未声明对应变量，注入实际无效
- 项目无版本检查机制

## 设计方案

### 1. 后端：版本注入修复

**`cmd/server/main.go`** — 声明版本变量：

```go
var (
    version = "dev"
    commit  = "unknown"
    date    = "unknown"
)
```

- 开发环境默认值 `dev`，构建时通过 `-ldflags "-X main.version=..."` 覆盖
- 启动日志输出版本：`slog.Info("UCAS starting", "version", version, "commit", commit)`

### 2. 后端：版本 API 端点

**`internal/api/router.go`** — 新增路由：

```
GET /api/version
```

响应：
```json
{
  "version": "1.1.0",
  "commit": "abc1234",
  "date": "2026-05-05T12:00:00+0800"
}
```

实现：新增 `getVersion` handler 函数，直接返回 `version`/`commit`/`date` 变量。

### 3. 前端：SidebarFooter 组件

**新建 `web/src/components/SidebarFooter.vue`**

布局：
```
┌─────────────────────┐
│  ⚙ GitHub 仓库  ↗  │  ← GitHub 图标 + 文字，可点击跳转
│  v1.1.0  ● NEW      │  ← 版本号 + 红点（有更新时显示）
└─────────────────────┘
```

样式：
- 使用 `margin-top: auto` 在 sidebar flex column 中沉底
- 文字颜色与侧边栏菜单项一致（`#64748b`）
- 红点使用 CSS `::after` 伪元素，红色圆点 + "NEW" 文字
- hover 时 GitHub 链接区域有轻微背景高亮
- 暗色模式适配

更新检查逻辑：
1. `onMounted` 时并行发起两个请求：
   - `GET /api/version` → 获取当前版本
   - `GET https://api.github.com/repos/zhaocongqi/clouddrive-auto-save/releases/latest` → 获取最新 release tag
2. 比较版本号（去掉 `v` 前缀后 semver 比较）
3. 有更新 → 红点 + "NEW" 标签，整体可点击跳转 `https://github.com/zhaocongqi/clouddrive-auto-save/releases`
4. 无更新 → 仅显示版本号，不可点击
5. GitHub API 调用失败 → 静默处理，不影响页面

### 4. 前端：MainLayout 集成

**修改 `web/src/layout/MainLayout.vue`**

- 在 `</el-menu>` 后引入 `<SidebarFooter />`
- sidebar 已有 `display: flex; flex-direction: column`，组件自动沉底

## 涉及文件

| 文件 | 操作 |
|------|------|
| `cmd/server/main.go` | 修改：声明 version/commit/date 变量 |
| `internal/api/router.go` | 修改：新增 `/api/version` 路由 + handler |
| `web/src/components/SidebarFooter.vue` | 新建 |
| `web/src/layout/MainLayout.vue` | 修改：引入 SidebarFooter 组件 |

## 验证方式

1. `make dev-server` 启动后端，访问 `http://localhost:8080/api/version` 确认返回版本信息
2. `make dev-web` 启动前端，检查侧边栏底部是否显示 GitHub 链接和版本号
3. 若当前版本低于最新 release，确认红点 + "NEW" 标签显示
4. 点击 GitHub 链接确认跳转正确
5. 暗色模式下样式正常
6. `make build` 构建后版本号不再是 `dev`
