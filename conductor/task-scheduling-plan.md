# Task Scheduling Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Implement a global task scheduling switch with local task overrides ("Global Default + Local Override" pattern) using `robfig/cron/v3`.

**Architecture:** 
- A new global settings model (`Setting`) to persist `global_schedule_enabled` and `global_schedule_cron`.
- Task model updated with `ScheduleMode` (`global`, `custom`, `off`). 
- Scheduler manages a global entry (triggering all "global" tasks) and individual entries for "custom" tasks.
- Frontend includes a global setting panel and a radio-button selection for each task's schedule mode.

**Tech Stack:** Go, Gin, GORM, `robfig/cron/v3`, Vue 3, Element Plus.

---

### Task 1: Update Database Models

**Files:**
- Modify: `internal/db/db.go`

- [x] **Step 1: Add ScheduleMode to Task**
Add `ScheduleMode` field to `Task` struct. Keep the existing `Cron` field for custom expressions.
- [x] **Step 2: Ensure Setting struct exists**
Check if `Setting` struct exists in `db.go`, if not, add it. And add `&Setting{}` to `AutoMigrate`.

### Task 2: Core Scheduler Updates

**Files:**
- Modify: `internal/core/scheduler/scheduler.go`

- [x] **Step 1: Update Scheduler struct**
Update Scheduler struct to have CustomEntryIDs and GlobalEntryID.
- [x] **Step 2: Implement Global and Custom Update logic**
Implement `UpdateGlobalSchedule(cronExpr string, enabled bool) error`.
Implement `UpdateTask(taskID uint, mode string, customCron string) error`.
Update `RemoveTask` to clear custom entries.

### Task 3: API Layer Integration

**Files:**
- Modify: `internal/api/router.go`
- Modify: `cmd/server/main.go`

- [x] **Step 1: Add System Settings Endpoints**
Add `GET /api/settings/schedule` and `POST /api/settings/schedule`.
- [x] **Step 2: Update Task Endpoints**
Update `createTask` and `updateTask` to handle `schedule_mode`. Call `scheduler.Global.UpdateTask` after save.
- [x] **Step 3: Lifecycle Management**
In `cmd/server/main.go`, load global settings from DB and initialize the global scheduler. Loop through tasks and initialize custom schedules.

### Task 4: Frontend UI Updates

**Files:**
- Modify: `web/src/views/Tasks.vue`
- Modify: `web/src/api/task.js` (add settings api)

- [x] **Step 1: Add Global Schedule Settings**
Add a top bar or drawer for Global Scheduling Settings (Switch & Cron input). Call the new settings API to sync.
- [x] **Step 2: Update Task Dialog**
Change the task edit form to use a radio group for `ScheduleMode`:
- 跟随全局 (Global)
- 自定义调度 (Custom) -> Shows custom cron input
- 关闭自动调度 (Off)
- [x] **Step 3: Update Table View**
Display schedule status in the table based on `schedule_mode` and `cron`.

### Task 5: Build and Verification

- [x] **Step 1:** Run Go tests.
- [x] **Step 2:** Build backend and run. Ensure database migration succeeds.
- [x] **Step 3:** Run frontend, test global toggle and individual overrides.
