# API 接口日志补全计划 (API Logging Enhancement Plan)

> **For agentic workers:** REQUIRED SUB-SKILL: Use
superpowers:subagent-driven-development or superpowers:executing-plans to
implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for
tracking.

**Goal:** 为后端的其他 API 接口补齐日志输出能力，以便在 Dashboard 的控制台中能够实时观察所有 API
的调用过程、参数以及执行结果，方便开发调试和线上定位问题。

**Architecture:**

- 基于当前的 `log.Printf` 机制（已经被劫持并广播到 SSE）。
- 修改 `internal/api/router.go`、`internal/api/folder.go` 以及
  `internal/api/task_preview.go` 中的核心处理函数。
- 在接口调用的入口、发生错误处、以及成功返回处增加结构化的日志信息（例如：`[API] 正在创建任务: ...`, `[API] 错误: ...`,
  `[API] 完成: ...`）。

**Tech Stack:** Go, Gin

---

## Task 1: 补充 Folder 模块接口日志

**Files:**

- Modify: `internal/api/folder.go`

- [ ] **Step 1: 增加 log 包的导入**
  在文件顶部的 import 中增加 `"log"`。

- [ ] **Step 2: 增强 getAccountFolders 日志**
  - 在函数入口打印获取目录树的请求参数。
  - 在遇到错误时打印详细错误日志。
  - 成功获取目录后，打印获取到的文件夹数量。

- [ ] **Step 3: 增强 createAccountFolder 日志**
  - 在入口处打印准备创建的文件夹信息。
  - 错误时打印失败详情。
  - 创建成功后打印结果信息。

## Task 2: 补充 Task 模块与通用模块接口日志

**Files:**

- Modify: `internal/api/router.go`
- Modify: `internal/api/task_preview.go`

- [ ] **Step 1: 增强 router.go 中的日志**
  为以下核心操作增加生命周期日志：
  - `checkAccount`: 打印校验账号的起始和结果。
  - `createAccount`, `updateAccount`, `deleteAccount`: 记录账号数据变更操作。
  - `createTask`, `updateTask`, `deleteTask`, `runTask`: 记录任务变更与启动操作（在 `runTask`
    已有的成功处增加 `log.Printf`）。

- [ ] **Step 2: 增强 previewTask 的日志**
  在 `internal/api/task_preview.go` 中的 `previewTask` 函数：
  - 开始时打印正在预览的任务和分享链接。
  - 成功计算后打印预览了多少个文件。

## Task 3: 验证与测试

- [ ] **Step 1: 重启后端服务**
  应用修改后重启 server 进程。
- [ ] **Step 2: 界面操作触发**
  在前端执行“校验账号”、“加载目录”、“修改任务”、“运行任务”等操作，确认控制台和 Dashboard 右侧日志窗口均能实时出现清晰的 `[API]`
前缀日志。
