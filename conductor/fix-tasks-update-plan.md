# 修复任务管理界面响应式更新计划

**目标：** 解决 `Tasks.vue` 表格中任务状态无法通过 SSE 事件自动更新，导致用户必须手动刷新页面的问题。

**背景：**
`Tasks.vue` 组件监听 SSE 流并尝试使用 `Object.assign` 将事件负载合并到本地 `taskList.value`
数组中。然而，如果后端返回的任务对象没有预加载 `Account` 关联，`Object.assign` 会用一个零值对象覆盖现有的
`row.account`，这可能破坏 UI 显示，或无法正确触发 Vue 对 `status` 和 `message` 单元格的深层响应式监听。

**修改内容 (`web/src/views/Tasks.vue`)：**

1. 放弃使用粗暴的 `Object.assign`，改为针对性地更新执行期间会变化的字段：
   - `status`
   - `message`
   - `last_run`
   - `percent`
   - `stage`
2. 这种精确赋值保证了 Vue 3 的 Proxy 能拦截到属性变动并重新渲染特定单元格，同时不会破坏 `account` 等嵌套关联数据。
3. 在任务达到 `Success` 或 `Failed` 状态时增加一次 `fetchList(true)` 的静默兜底刷新，确保任务完成时全行数据的一致性。
