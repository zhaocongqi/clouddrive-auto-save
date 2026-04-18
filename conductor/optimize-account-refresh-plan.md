# 优化账号校验后的 UI 刷新计划

**目标：** 避免用户点击单个账号的“校验”时刷新整页数据（`fetchList()`），改为仅原地更新表格中该行的数据，提升交互体验。

**背景：**
后端 `checkAccount` API 在成功时已返回更新后的账号对象，失败时也会在错误负载中带回最新的账号快照（包含更新后的 `last_check` 和 `status`）。

**修改内容 (`web/src/views/Accounts.vue`)：**
1. 更新 `handleCheck(row)` 函数，捕获 API 的返回值。
2. 成功时：直接使用 `Object.assign(row, updatedAccount)` 原地同步数据。
3. 失败时：从错误拦截器的响应中提取 `account` 对象，同样执行原地同步。
4. 移除 `finally` 代码块中的 `fetchList()` 调用。

这将使 UI 响应瞬间完成，并减少全量刷新的网络开销。
