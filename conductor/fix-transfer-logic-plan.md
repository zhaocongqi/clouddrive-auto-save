# 任务转存逻辑与起始文件功能修复计划 (Fix Task Transfer & StartFile Plan)

**Goal:** 修复任务执行时的核心逻辑，确保“从起始文件向前转存”、“正则匹配过滤”以及“基于新名称的去重”功能恢复正常。

---

## Task 1: 重构 Worker 执行核心逻辑

**Files:**

- Modify: `internal/core/worker/worker.go`

**Changes:**

1. **获取分享列表并排序**:
   - 调用 `driver.ParseShare` 获取原始列表。
   - **新增排序**: 必须显式按 `UpdateTime` 从新到旧（降序）排序，确保“起始文件”判断逻辑的稳定性。
2. **应用起始文件过滤 (StartFile Logic)**:
   - 如果 `task.StartFileID` 不为空：
     - 找到该 ID 在已排序列表中的位置。
     - **截断列表**: 仅保留该文件本身及其之后更新（即索引更靠前）的所有文件。
3. **集成重命名与智能去重**:
   - 遍历过滤后的列表。
   - 预计算 `NewName`。
   - 校验正则 `Matched` 状态，若设置了正则且不匹配，则跳过（实现“只转存匹配文件”）。
   - 拿 `NewName` 去目标目录 Map 中比对，若已存在则跳过。
4. **提交转存**: 最终将 `filteredIDs` 提交。

## Task 2: 同步修复 API 解析接口

**Files:**

- Modify: `internal/api/task_preview.go`

**Changes:**

1. **`parseShareLinkInfo` 增强**:
   - 接收 `Pattern`, `Replacement`, `Name` 参数。
   - 获取目标目录列表，进行重命名预览和同名预检。
   - 确保返回的 `is_existed` 和 `new_name` 准确无误。

## Task 3: 验证修复效果

- [ ] **验证起始点**: 选择分享中的某一个文件作为起始点，确认之前的旧文件不会被转存。
- [ ] **验证正则过滤**: 设置一个过滤规则，确认不符合规则的文件被标记为“跳过”或“不匹配”。
- [ ] **验证去重**: 在网盘中手动创建一个与重命名后同名的文件，确认解析时显示“已在网盘”。

## Task 4: 提交代码

- 本地提交，不执行推送。
