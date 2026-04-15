# 重构目录树加载逻辑与修复根目录新建文件夹问题

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** 通过引入虚拟“根目录”节点（方案 B）来彻底重构 `web/src/views/Tasks.vue` 中的树形组件加载逻辑，移除旧有的硬编码和补丁判断，从而优雅地解决根目录新建文件夹的报错及视图不同步问题。

**Architecture:** 
- **虚拟根节点**：在 Element Plus `el-tree` 首次加载 (`node.level === 0`) 时，不直接请求后端 API，而是手动解析并返回一个统一的 `[{ label: '根目录', path: '/', id: '', isLeaf: false }]` 节点。
- **层级顺延**：所有真实的 API 数据拉取动作顺延到 `node.level >= 1` 时发生。
- **逻辑精简**：因为任何被选中的节点（哪怕是最顶层的根目录）在树中都有了物理实体，所以获取选中 ID、在父节点下 `append` 新文件夹、同步高亮状态等操作都变成了常规逻辑，不再需要写繁琐的 `if (currentPath === '/')` 补丁。

**Tech Stack:** Vue 3, Element Plus

---

### Task 1: 引入虚拟根节点并修改加载逻辑

**Files:**
- Modify: `web/src/views/Tasks.vue`

- [ ] **Step 1: 给 `el-tree` 增加默认展开属性**
  修改 `<el-tree>` 标签，添加 `:default-expanded-keys="['/']"`，使用户在打开弹窗时，根目录节点能自动展开并加载第一层数据。

- [ ] **Step 2: 重写 `loadFolders` 方法**
  拦截 `node.level === 0`，返回固定的根目录对象。真实的后端请求从 `node.level >= 1` 开始。
  ```javascript
  // 加载网盘目录树
  const loadFolders = async (node, resolve) => {
    if (!form.value.account_id) {
      return resolve([])
    }

    // 0层：直接返回虚拟根节点
    if (node.level === 0) {
      pathIdMap.value['/'] = '' // 确保映射表里有根节点的合法记录
      return resolve([{ label: '根目录', path: '/', id: '', isLeaf: false }])
    }

    // 1层及以上：使用父节点（虚拟根或真实目录）的属性请求后端
    const parentID = node.data.id
    const parentPath = node.data.path

    if (node.level === 1) {
      loadingFolders.value = true
    }
    try {
      const folders = await getFolders(form.value.account_id, parentID, parentPath)
      folders.forEach(f => {
        pathIdMap.value[f.path] = f.id
      })
      resolve(folders)
    } catch (err) {
      console.error('加载目录失败:', err)
      resolve([])
    } finally {
      if (node.level === 1) {
        loadingFolders.value = false
      }
    }
  }
  ```

### Task 2: 精简新建文件夹逻辑 (移除旧补丁)

**Files:**
- Modify: `web/src/views/Tasks.vue`

- [ ] **Step 1: 整理 `handleInlineCreateFolder` 逻辑**
  删除之前添加的 `pathIdMap` 强行回退和各种 `undefined` 检查，由于虚拟根节点的存在，`selectedTreePath` 和 `selectedTreeId` 会始终保持可靠的配对。
  ```javascript
  // 内嵌新建文件夹处理
  const handleInlineCreateFolder = async () => {
    if (!newFolderName.value.trim()) {
      return ElMessage.warning('请输入文件夹名称')
    }

    const currentPath = selectedTreePath.value || '/'
    // 从选中状态取 ID，若为空则默认取虚拟根节点的空字符串 ''
    const currentID = selectedTreeId.value || '' 

    creatingFolder.value = true
    try {
      const res = await createFolder(form.value.account_id, currentID, currentPath, newFolderName.value.trim())
      ElMessage.success('文件夹创建成功')
      pathIdMap.value[res.path] = res.id
      
      if (folderTreeRef.value) {
        // 由于存在虚拟根节点，getNode 一定能找到有效的父节点（包括 '/'）
        const currentNode = folderTreeRef.value.getNode(currentPath)
        if (currentNode) {
          folderTreeRef.value.append(res, currentNode)
          currentNode.expanded = true
        }
        selectedTreePath.value = res.path
        selectedTreeId.value = res.id
        folderTreeRef.value.setCurrentKey(res.path)
      }
      newFolderName.value = ''
    } catch (err) {
      console.error('创建文件夹失败:', err)
    } finally {
      creatingFolder.value = false
    }
  }
  ```

### Task 3: 验证与测试

- [ ] **Step 1: 编译测试**
  在 `web/` 目录下运行编译命令确保无报错。
- [ ] **Step 2: 界面交互测试**
  1. 验证“选择保存目录”弹窗弹出时，“根目录”是否自动显示并展开。
  2. 选中“根目录”并新建文件夹，确认前端界面是否立即同步出现了新建立的文件夹，且高亮状态切换到新文件夹。
  3. 验证旧有的“无法确定当前选中目录”的报错是否完全消除。
