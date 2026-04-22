# 初始交互体验与系统启动优化计划 (Initial UX & Startup Fix Plan)

**Goal:**

1. 优化账号管理和任务管理页面在“无数据”时的展示效果，添加带有插画的引导式 `<el-empty>` 占位和创建按钮。
2. 修复后端首次启动时，因全局调度配置（`global_schedule_enabled` 和 `global_schedule_cron`）未初始化而导致
   GORM 抛出 `record not found` 错误的问题。

---

## Task 1: 优化 Accounts.vue 的空状态展示

**Files:**

- Modify: `web/src/views/Accounts.vue`

**Changes:**

1. 在表格视图 (`table-card`) 和卡片视图 (`card-view-container`) 且数据为空时，隐藏原有的简单提示。
2. 引入 Element Plus 的 `<el-empty>` 组件，增加友好的引导语（“您还没有绑定任何云盘账号”），并使用插槽 `#default`
   提供一个主按钮（“立即绑定账号”），点击后触发 `openAddDialog`。

## Task 2: 优化 Tasks.vue 的空状态展示

**Files:**

- Modify: `web/src/views/Tasks.vue`

**Changes:**

1. 在任务表格数据为空时，隐藏原有的空表格。
2. 引入 `<el-empty>` 组件，增加引导语（“当前没有任何转存任务”），并使用插槽 `#default`
   提供一个主按钮（“创建新任务”），点击后触发 `openAddDialog`。

## Task 3: 修复后端启动时的 Record Not Found 错误

**Files:**

- Modify: `cmd/server/main.go`

**Changes:**

1. 将获取全局调度设置时使用的 `First()` 替换为 `Find()`（或显式处理 `gorm.ErrRecordNotFound`）。使用
   `Find()` 是最简单的避免 `record not found` 报错的方式，当查不到数据时，返回的设置对象属性为空字符串，后续
   `UpdateGlobalSchedule` 也能正常处理。

```go
 // 修改前
 db.DB.Where("key = ?", "global_schedule_enabled").First(&enabledSetting)
 db.DB.Where("key = ?", "global_schedule_cron").First(&cronSetting)
 
 // 修改后
 db.DB.Where("key = ?", "global_schedule_enabled").Find(&enabledSetting)
 db.DB.Where("key = ?", "global_schedule_cron").Find(&cronSetting)
```
