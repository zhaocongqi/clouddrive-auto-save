# 修复监控面板失败任务持久化显示计划

**目标：** 确保当用户在“任务执行监控”面板点击“X”忽略一个失败任务时，该任务会保持隐藏状态，而不会在下一次 5 秒自动轮询时重新出现。

**修改内容：**

### 1. 后端 API (`internal/api/router.go`)
- **新增接口**：`POST /api/tasks/:id/dismiss`。
  - 逻辑：在数据库中将该任务的 `stage` 字段更新为 `"Dismissed"`。
- **优化查询**：在 `getDashboardStats` 接口中，更新 `runningTasksList` 的查询语句。
  - 过滤条件：排除掉状态为 `status = "failed"` 且阶段为 `stage = "Dismissed"` 的任务。
  - 效果：确保只有“新发生的”或“未处理”的失败任务会被展示。

### 2. 前端 API (`web/src/api/task.js`)
- 增加 `dismissTask(taskId)` 函数，用于调用新增的后端接口。

### 3. 前端 UI (`web/src/views/Dashboard.vue`)
- 更新 `dismissTask` 处理函数：
  - 改为异步函数 `async`。
  - 在从本地 `runningTasks` 数组移除卡片之前，先调用 `dismissTask` 接口。
  - 将“忽略”状态持久化到数据库。

**重新出现逻辑：**
- 无需修改 `worker.go`。当任务被再次手动触发或定时触发时，Worker 会自动调用 `updateProgress` 将 `stage` 重置为 `"Started"`，这会自动清除 `"Dismissed"` 标记，使得任务在运行或再次失败时能重新出现在监控面板上。
