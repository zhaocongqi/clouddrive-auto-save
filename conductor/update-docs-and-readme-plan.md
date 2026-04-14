# 更新文档说明计划 (Update Documentation Plan)

## 1. 目标 (Objective)
根据近期实施的几个重要修复（账号删除安全校验、139云盘家庭容量计算修复、网盘接口响应日志 SSE 广播），更新项目的核心文档 `README.md` 和 API 接口文档 `@docs/api/**`，确保文档与最新系统行为一致。

## 2. 涉及文件 (Key Files & Context)
- `README.md`: 项目主页说明。
- `docs/api/accounts.md`: 账号管理的 API 文档。
- `docs/api/dashboard.md`: 仪表盘统计与日志的 API 文档。
- `docs/cloud_drive_apis.md`: 云盘接口文档（可能需要补充底层日志拦截相关）。

## 3. 实施细节 (Implementation Steps)

### 3.1. 更新 `docs/api/accounts.md`
- **补充删除账号接口 (`DELETE /accounts/:id`)**：
  - 明确列出该接口。
  - **重点说明安全拦截机制**：如果账号下存在关联任务，删除操作将被拒绝，返回 `409 Conflict` (或 `400`) 并提示“该账号有关联的 X 个任务，请先删除任务”。
- **完善容量字段说明 (`Account` 对象)**：
  - 在 `capacity_used` 和 `capacity_total` 字段的说明中，增加关于 139 云盘家庭空间的批注：总容量包含家庭空间，而已用容量仅计算个人空间的使用情况。

### 3.2. 更新 `docs/api/dashboard.md`
- **新增实时日志接口 (`GET /dashboard/logs`)**：
  - 增加该接口的文档，说明其返回格式为 `text/event-stream` (Server-Sent Events, SSE)。
  - 说明其作用是向前端实时广播网盘底层的请求/响应日志。
- **新增获取近期日志接口 (`GET /dashboard/logs/recent`)**：
  - 增加该接口的文档，说明其用于获取系统缓存的近期（历史）日志数组。

### 3.3. 更新 `README.md`
- 在 **✨ 核心特性 (Features)** 部分追加最新的优化点：
  - **深度排错与监控**：通过 Server-Sent Events (SSE) 将网盘底层 API 请求响应日志实时广播至 Dashboard 进行展示。
  - **数据安全与精准统计**：增加了账号删除时的任务关联强校验（防止误删），并修正了 139 云盘家庭空间的容量精确统计模型。

### 3.4. 检查并更新 `docs/cloud_drive_apis.md` (可选)
- 检查该文件是否包含底层请求拦截或日志输出的说明，如有必要，简述底层客户端（139/Quark）现已通过 `GlobalBroadcaster` 全局广播 API 响应。

## 4. 验证与测试 (Verification & Testing)
1. 检查 `README.md` 渲染是否正常，新增特性描述是否清晰。
2. 检查 `docs/api/accounts.md` 中删除接口及容量批注是否完整、无 Markdown 语法错误。
3. 检查 `docs/api/dashboard.md` 中新增的 `/dashboard/logs` 和 `/dashboard/logs/recent` 接口说明是否符合文档规范。
4. 确保所更新的文档内容完全匹配之前相关的修复方案和实际代码实现。