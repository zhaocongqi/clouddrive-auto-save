# 异步更新任务状态计划 (Async Task Status Plan)

## 1. 目标 (Objective)

在任务管理页面 (`Tasks.vue`) 中，通过添加定时轮询功能，实现任务执行状态的自动（异步）更新，避免用户需要手动刷新页面才能看到最新状态。

## 2. 涉及的关键文件 (Key Files & Context)

* **前端视图:** `web/src/views/Tasks.vue`

## 3. 实施步骤 (Implementation Steps)

### 3.1 修改 Vue 导入

在 `<script setup>` 顶部引入 `onUnmounted` 钩子：

```javascript
import { ref, onMounted, onUnmounted } from 'vue'
```

### 3.2 优化 fetchList 方法

修改现有的 `fetchList` 方法，增加 `silent` 参数。当 `silent` 为 `true` 时，不触发整页的 `loading`
动画，从而实现无感刷新：

```javascript
let pollTimer = null

const fetchList = async (silent = false) => {
  if (!silent) loading.value = true
  try {
    const [taskData, accountData] = await Promise.all([getTasks(), getAccounts()])
    taskList.value = taskData
    accounts.value = accountData
  } catch (err) {
    console.error(err)
  } finally {
    if (!silent) loading.value = false
  }
}
```

### 3.3 添加定时轮询逻辑

在组件挂载时 (`onMounted`) 启动定时器，每 5 秒进行一次无感刷新；在组件卸载时 (`onUnmounted`) 清除定时器以防内存泄漏：

```javascript
onMounted(() => {
  fetchList()
  pollTimer = setInterval(() => {
    fetchList(true)
  }, 5000)
})

onUnmounted(() => {
  if (pollTimer) clearInterval(pollTimer)
})
```

## 4. 验证与测试 (Verification & Testing)

1. 启动服务并在浏览器中访问任务管理页面。
2. 触发一次任务执行，无需刷新页面，观察表格中对应任务的状态是否会在短时间内（5 秒左右）自动从“运行中”更新为“成功”或“失败”。
