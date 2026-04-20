# 修复任务编辑后账号信息消失的问题 (Fix Task Account Info Disappearance Plan)

**Goal:** 修复在创建或更新任务后，任务列表中的账号昵称（Account Nickname）会暂时消失直到刷新页面的问题。

---

### Task 1: 优化后端任务创建与更新接口

**Files:**
- Modify: `internal/api/router.go`

**Changes:**
1. **`createTask` 函数**: 在 `db.DB.Create(&task)` 之后，增加 `db.DB.Preload("Account").First(&task, task.ID)`。
2. **`updateTask` 函数**: 在 `db.DB.Model(&task).Updates(updateData)` 之后，增加 `db.DB.Preload("Account").First(&task, task.ID)`。

**Reasoning:**
目前的后端接口在创建或更新任务后，直接返回了 `task` 对象，但并没有预加载关联的 `Account` 对象。前端在收到响应后，通过 `Object.assign` 更新本地列表，导致原本存在的 `account` 字段被覆盖为空（或不含昵称的空结构），从而引起 UI 上的账号信息消失。通过在返回前强制预加载，可以确保前端始终持有完整的关联数据。

---

### Task 2: 验证修复效果

- [ ] 创建一个新任务。
- [ ] 确认任务创建成功后，列表中的账号昵称正常显示。
- [ ] 编辑该任务并保存。
- [ ] 确认保存后，任务名下方的账号昵称依然存在，没有消失。

---

### Task 3: 自动提交

- 提交代码更改。
- 不自动推送。
