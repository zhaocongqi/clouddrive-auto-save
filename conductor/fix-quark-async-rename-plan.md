# 修复预览逻辑与夸克异步转存问题计划 (Fix Preview & Quark Async Issue Plan)

**Goal:**

1. 修复预览逻辑中“不匹配正则的文件也显示重命名后名字”的问题。
2. 修复夸克网盘转存后立即重命名失败的问题（原因是夸克转存是异步任务，需等待完成后才能重命名）。

---

## Task 1: 优化重命名预览展示逻辑

**Files:**

- Modify: `internal/api/task_preview.go`

**Changes:**

1. 修改 `previewTask` 函数：
   - 只有在 `matched` 为 `true` 时，才将计算出的 `newName` 赋值给返回结果。
   - 如果不匹配正则，`new_name` 保持为原始文件名，确保预览效果与实际执行逻辑一致。
2. 修改 `parseShareLinkInfo` 函数：
   - 同样增加正则匹配判定。
   - 仅对匹配正则的文件应用预览名，不匹配的文件 `NewName` 保持为原始 `Name`。

## Task 2: 修复夸克异步转存导致的重命名失效

**Files:**

- Modify: `internal/core/quark/client.go`

**Changes:**

1. **解析 TaskID**: 修改 `SaveLink` 函数，解析转存 API 返回的 `task_id`。
2. **实现任务等待逻辑**:
   - 新增 `waitTask(taskID string)` 方法。
   - 轮询 `GET /1/clouddrive/task?task_id=...` 接口。
   - 当 `status` 变为 `2` (成功) 时返回；如果任务失败或超时，则返回错误。
3. **阻塞转存**: 在 `SaveLink` 返回前，调用 `waitTask` 确保文件已真实进入目标目录。这样 Worker 在下一步执行 `ListFiles` 时就能看到这些新文件，从而成功重命名。

## Task 3: 验证修复效果

- [ ] **验证预览**: 设置正则规则，确认预览列表中只有匹配的文件显示了新名称，不匹配的显示原名。
- [ ] **验证夸克重命名**: 执行一个包含重命名规则的夸克转存任务，确认转存完成后，文件能被成功重命名，而不是保留原始名称。
