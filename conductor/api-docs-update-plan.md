# API 文档深度校验与更新计划 (API Docs Update Plan)

**Goal:** 基于全量代码审查结果，修正 `@docs/api/` 目录下文档与当前实际代码实现之间的不匹配、遗漏和过时内容。

---

## Task 1: 补充任务管理接口文档 (tasks.md)

**Files:**

- Modify: `docs/api/tasks.md`

**Changes:**

1. **完善 `updateTask` 逻辑**: 在“更新任务配置”小节中，明确声明当 `share_url` 或 `extract_code` 变动时，系统将自动重置状态并清空 `[Fatal]` 错误。
2. **新增 `dismiss` 接口**: 补充 `POST /tasks/:id/dismiss` (忽略/消除错误状态) 接口的完整说明，因为该接口在代码中已实现并投入使用，但文档中缺失。

## Task 2: 修正账号目录接口文档 (accounts.md)

**Files:**

- Modify: `docs/api/accounts.md`

**Changes:**

1. **修正 Payload**: 在“新建网盘文件夹” (`POST /accounts/:id/folders`) 小节中，移除 `parent_id` 字段（代码实现中并未接收此参数，仅依靠 `parent_path` 和 `name` 即可实现兼容）。
2. **完善返回结构**: 在“获取网盘目录树” (`GET /accounts/:id/folders`) 的 `FolderItem` 结构中，补充说明目前返回的对象同时包含 `name` 和 `label` 字段（两者值相同，为兼容不同组件而设）。

## Task 3: 格式校验与自动提交

- 执行 `make lint-md` 修复排版。
- 将更改提交至本地仓库。
