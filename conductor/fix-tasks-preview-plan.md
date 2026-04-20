# 任务预览逻辑修复计划 (Task Preview Fix Plan)

**Goal:** 修复任务编辑页面“全量重命名预览”功能中原始文件名不显示以及状态判断不准确的问题。

---

### Task 1: 修正后端预览接口字段名与逻辑

**Files:**
- Modify: `internal/api/task_preview.go`

**Changes:**
1. 将返回结果中的 `old_name` 字段更名为 `original_name`，以匹配前端 `Tasks.vue` 模板及 API 文档。
2. 完善 `matched` 逻辑：
   - 如果 `Pattern`（重命名正则）为空，只要有 `Replacement` 或预览正在进行，就默认将 `matched` 设为 `true`，以确保前端显示“匹配”成功效果。
   - 如果 `Pattern` 不为空，则通过 `regexp` 进行实际匹配测试。
3. 添加 `is_filtered` 字段，默认设为 `false`，解决前端预览列表中的状态显示异常。

### Task 2: 验证修复效果

- [ ] 打开任务编辑对话框。
- [ ] 保持“重命名正则”为空，填写“替换规则”。
- [ ] 点击“全量重命名预览”，确认“原始文件名”列正常显示，且“状态”列显示为“匹配”。