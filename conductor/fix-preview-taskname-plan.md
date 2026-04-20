# 修复预览时任务名硬编码问题的计划 (Fix Preview TaskName Plan)

**Goal:** 解决在“重命名预览”功能中，`{TASKNAME}` 变量始终显示为 "Preview" 的问题，使其能够实时反映用户当前输入的任务名称。

---

### Task 1: 修正后端预览接口以接收任务名

**Files:**
- Modify: `internal/api/task_preview.go`

**Changes:**
1. 修改 `previewTask` 的请求结构体，增加 `Name` 字段（JSON 标签为 `name`）。
2. 在循环处理文件预览时，将原本硬编码传给 `processor.Process` 的 `"Preview"` 字符串改为使用 `req.Name`。
3. 如果 `req.Name` 为空，则默认使用 `"Task"` 字符占位。

### Task 2: 验证修复效果

- [ ] 在创建任务弹窗中输入任务名称（例如：`功夫熊猫4`）。
- [ ] 在替换规则中使用 `{TASKNAME}` 变量（例如：`{TASKNAME}.{EXT}`）。
- [ ] 点击全量预览，确认生成的预览名中包含 `功夫熊猫4` 而不是 `Preview`。

### Task 3: 提交更改
- 提交后端代码。
- 不自动推送。
