# 修复分享解析中的重命名预览与同名预检计划 (Fix Share Parse Logic Plan)

**Goal:** 完善 `/api/tasks/parse_share` 接口，使其支持根据任务正则实时预览新文件名，并自动预检目标目录以标记已存在的文件。

---

## Task 1: 升级后端 `parseShareLinkInfo` 接口

**Files:**

- Modify: `internal/api/task_preview.go`

**Changes:**

1. **扩展请求结构体**: 增加 `Pattern`, `Replacement`, `Name` 字段。
2. **获取目标目录快照**:
   - 调用 `driver.PrepareTargetPath` 获取目标目录的真实 ID。
   - 调用 `driver.ListFiles` 获取该目录下已有的文件列表，并存入 Map 以供快速比对。
3. **集成重命名与预检循环**:
   - 遍历分享文件列表。
   - 调用 `renamer.Processor` 计算每个文件的 `NewName`。
   - 判定 `Matched` 状态（处理 Pattern 为空的情况）。
   - **智能去重**: 拿计算出的 `NewName` 去目标目录快照中比对，如果存在，则标记 `IsExisted = true`。

## Task 2: 验证修复效果

- [ ] 打开任务编辑对话框，设置正则和保存路径。
- [ ] 点击“选择文件”打开解析弹窗。
- [ ] **验证 1**: 确认文件列表中的“新名称”列显示了正则处理后的效果。
- [ ] **验证 2**: 如果目标目录已存在同名文件，确认该行显示为“已在网盘”且呈现灰色禁用状态。

## Task 3: 自动提交

- 提交代码更改。
- 不自动推送。
