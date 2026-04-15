# 修复夸克网盘在根目录新建文件夹的报错问题

## 1. Objective (目标)
解决用户反馈的在任务管理中点击“选择目录”并在**根目录**新建文件夹时，夸克网盘会报“无法确定当前选中目录的 ID，请重新展开该节点”的错误（而移动云盘不报错）的问题。

## 2. Key Files & Context (关键文件与上下文)
- `web/src/views/Tasks.vue`
  - 触发点：`handleInlineCreateFolder` 方法。
  - 错误原因：`pathIdMap.value[currentPath]`（即路径到 ID 的映射表）在某些特殊边界情况下（例如切换任务编辑或初始化时遗漏），对于根目录 `/` 的映射值变为 `undefined`。这导致了 `if (currentID === undefined)` 的拦截逻辑被意外触发。
  - 为什么 139 云盘不容易触发？可能与任务回显状态或特定于操作习惯相关。无论如何，对于“根目录”这个特殊节点，其父 ID 在后端处理中具有确定的回退逻辑（Quark 会自动将 `""` 转为 `"0"`，139 自动处理为空或 `/`），因此前端可以进行绝对安全的防御性降级。

## 3. Implementation Steps (实施步骤)
修改 `web/src/views/Tasks.vue` 中的 `handleInlineCreateFolder` 函数，增加对根目录 `/` 的防御性回退逻辑：

```javascript
const handleInlineCreateFolder = async () => {
  if (!newFolderName.value.trim()) {
    return ElMessage.warning('请输入文件夹名称')
  }
  
  const currentPath = selectedTreePath.value || '/'
  let currentID = pathIdMap.value[currentPath]

  // 【核心修复】：防御性回退
  // 如果当前是在根目录（'/'），但由于切换等原因映射表丢失了根目录的 ID 缓存，
  // 我们强制将其指定为空字符串。
  // 后端接口（如夸克驱动）对于 parent_id 为空字符串的请求，会自动安全地解析为根节点（"0"）。
  if (currentID === undefined && currentPath === '/') {
    currentID = ''
  }

  if (currentID === undefined) {
    return ElMessage.warning('无法确定当前选中目录的 ID，请重新展开该节点')
  }

  creatingFolder.value = true
  // ... 后续正常发送 API 请求
}
```

## 4. Verification & Testing (验证与测试)
1. 应用更改并重新构建前端项目（`npm run build`）。
2. 在浏览器中打开并选择一个夸克网盘账号。
3. 点击“选择保存目录”弹出树形结构框。
4. 保持未选中任何子节点（或主动使得路径为 `/`），在底部输入新文件夹名并点击“新建”。
5. 验证是否直接创建成功，不再抛出“无法确定当前选中目录的 ID”的异常。
