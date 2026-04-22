# 修复目录树不显示文件夹名称的问题 (Fix Folder Tree Display Plan)

**Goal:** 修复编辑任务时，目录选择弹窗中的目录树不显示文件夹名称的问题。原因是后端接口返回的字段名为 `name`，而前端 `el-tree` 组件配置期望的是 `label` 字段。

---

## Task 1: 修正后端目录接口返回结构

**Files:**

- Modify: `internal/api/folder.go`

**Changes:**

1. 修改 `getAccountFolders` 函数中的结果组装逻辑。
2. 将返回的对象从 `core.FileInfo` 转换为包含 `label` 和 `isLeaf` 字段的 Map 结构。
3. 确保 `label` 字段被赋予文件夹的名称，`isLeaf` 统一设为 `false`（支持继续展开）。

## Task 2: 验证修复效果

- [ ] 打开任务编辑对话框。
- [ ] 点击“选择目录”。
- [ ] 确认目录树中正常显示“根目录”及下属各级文件夹名称。
