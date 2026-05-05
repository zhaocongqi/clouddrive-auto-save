# Bark 批量通知集成设计

## 概述

将 Bark 推送与任务管理模块深度集成，支持单任务和批量执行两种通知模式。单任务完成后立即推送；批量执行（全部运行）时等待所有任务完成，汇总为一条通知。

## 背景与动机

当前 Bark 推送仅在 `finishTask` 中统一触发，不区分单任务和批量执行场景。当用户点击"全部运行"时，会收到 N 条独立通知，造成信息轰炸。需要引入批量追踪机制，在批量执行完成后发送一条汇总通知。

## 核心需求

- **单任务执行**：完成后立即推送（保持现有行为）
- **批量执行**：等所有任务完成后，发送一条汇总推送
- **汇总内容**：统计摘要 + 每个任务详情（名称、状态、耗时、文件数）+ 转存文件列表
- **批量中的失败**：不立即通知，等全部完成再汇总
- **通知粒度**：仅全局开关（保持现有，不按任务独立控制）
- **消息格式**：单任务保持现有格式

## 方案选型

| 方案 | 描述 | 优点 | 缺点 |
|------|------|------|------|
| A: Worker 层 BatchTracker | Worker 包中新增追踪器，finishTask 区分单/批量模式 | 改动集中，逻辑自包含 | 进程重启丢失批次状态 |
| B: Notify 层聚合器 | notify 包中新增聚合器，worker 完成后上报 | worker 改动最小 | notify 承担状态管理职责 |
| C: 数据库批次追踪 | 新增 batch_executions 表持久化 | 可追溯历史 | 需要 DB 迁移，过度设计 |

**选定方案 A**：改动范围可控，逻辑清晰，进程重启丢失批次状态在实际使用中影响很小。

## 架构设计

### 数据结构

```go
// worker/batch_tracker.go

// BatchResult 收集单个任务的执行结果
type BatchResult struct {
    TaskName string
    Status   string    // "success" | "failed"
    Message  string
    Files    []string
    Duration time.Duration
}

// BatchTracker 追踪批量执行状态
type BatchTracker struct {
    mu      sync.Mutex
    batches map[string]*batchInfo
}

type batchInfo struct {
    total   int
    count   int            // 已完成计数
    results []BatchResult
    startAt time.Time
}
```

### BatchTracker 方法

```go
// NewBatchTracker 创建追踪器实例
func NewBatchTracker() *BatchTracker

// RegisterBatch 注册一个新批次
func (t *BatchTracker) RegisterBatch(batchID string, total int)

// ReportTask 上报单个任务完成结果
// 当 count == total 时，触发汇总推送并清理批次数据
func (t *BatchTracker) ReportTask(batchID string, result BatchResult)
```

`BatchTracker` 在 `NewManager` 中初始化：`tracker: NewBatchTracker()`。

### Worker 层改动

**Job 结构体扩展**：

```go
type Job struct {
    Task    *db.Task
    BatchID string  // 为空表示单任务执行
}
```

**Manager 新增字段**：

```go
type Manager struct {
    // ... 现有字段
    tracker *BatchTracker
}
```

**finishTask 签名变更**：

```go
// 改前
func (m *Manager) finishTask(task *db.Task, status, message string, files []string, startTime time.Time)

// 改后
func (m *Manager) finishTask(job Job, status, message string, files []string, startTime time.Time)
```

**finishTask 内部逻辑变更**：

```go
func (m *Manager) finishTask(job Job, status, message string, files []string, startTime time.Time) {
    task := job.Task
    duration := time.Since(startTime)

    // ... 现有 DB 更新和 SSE 广播逻辑不变 ...

    if job.BatchID != "" {
        // 批量模式：上报结果给 tracker，不单独推送
        m.tracker.ReportTask(job.BatchID, BatchResult{
            TaskName: task.Name, Status: status,
            Message: message, Files: files, Duration: duration,
        })
    } else {
        // 单任务模式：保持现有行为，立即推送
        notify.SendTaskNotification(task.Name, status, message, files, duration)
    }
}
```

### Notify 层新增

在 `notify` 包中新增 `SendBatchNotification`：

```go
// notify/batch.go

// BatchResult 批量任务结果（避免循环依赖，notify 包自定义）
type BatchResult struct {
    TaskName string
    Status   string
    Message  string
    Files    []string
    Duration time.Duration
}

// SendBatchNotification 发送批量执行汇总通知
func SendBatchNotification(results []BatchResult, totalDuration time.Duration) {
    // 1. 统计成功/失败数
    // 2. 构造标题
    // 3. 构造正文（任务详情 + 文件列表）
    // 4. 根据是否有失败选择级别和铃声
    // 5. 异步发送
}
```

**类型转换**：`BatchTracker.ReportTask` 在全部完成时，将 `[]worker.BatchResult` 转换为 `[]notify.BatchResult`（字段完全一致，逐元素拷贝），然后调用 `notify.SendBatchNotification`。

### API 层改动

`runAllTasks` handler 变更：

```go
func runAllTasks(c *gin.Context) {
    // ... 现有查询任务逻辑不变 ...

    // 新增：生成批次 ID 并注册
    batchID := fmt.Sprintf("batch_%d", time.Now().UnixMilli())
    wm.RegisterBatch(batchID, len(tasks))

    for _, task := range tasks {
        // ... 现有状态更新逻辑不变 ...
        wm.Submit(worker.Job{Task: &task, BatchID: batchID})
    }

    c.JSON(200, gin.H{"message": fmt.Sprintf("已触发 %d 个任务", len(tasks))})
}
```

## 消息格式

### 单任务（保持现有）

- 成功标题：`"✅ 转存任务完成: {taskName}"`
- 失败标题：`"❌ 转存任务失败: {taskName}"`
- 正文：任务结果消息 + 执行耗时 + 文件列表（最多 10 个）

### 批量汇总

**标题**：
- 全部成功：`"📊 批量转存完成: 全部 5 个任务成功"`
- 部分失败：`"📊 批量转存完成: 3成功 / 2失败"`

**正文示例**：
```
总耗时: 3分25秒

✅ 任务A - 转存成功 (新增 5 个文件) - 耗时 45s
✅ 任务B - 转存成功 (新增 3 个文件) - 耗时 32s
❌ 任务C - 解析分享失败: 链接过期 - 耗时 2s
✅ 任务D - 转存成功 (新增 8 个文件) - 耗时 1m12s
❌ 任务E - 转存失败: 空间不足 - 耗时 5s

转存文件列表:
- movie1.mp4
- movie2.mp4
- episode01.mkv
... 等共 16 个文件
```

**推送级别**：有失败任务使用 `bark_failure_level` 配置，全部成功使用 `bark_success_level`。铃声同理。

## 文件变更清单

| 文件 | 变更类型 | 说明 |
|------|----------|------|
| `internal/core/worker/batch_tracker.go` | 新增 | BatchTracker 实现 |
| `internal/core/worker/worker.go` | 修改 | Job 扩展、Manager 新增 tracker、finishTask 逻辑变更 |
| `internal/core/notify/batch.go` | 新增 | SendBatchNotification 实现 |
| `internal/api/router.go` | 修改 | runAllTasks 生成 batchID 并注册 |

## 约束与边界

- 进程重启会丢失进行中的批次状态（可接受，批量执行通常几分钟内完成）
- 不支持按任务独立控制通知开关
- 不改变现有的 SSE 广播机制
- 不改变现有的 Bark 配置项结构
