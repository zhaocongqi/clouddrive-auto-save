# 保存路径 UI 优化实施计划 (Save Path UI Optimization Plan)

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** 重构任务管理弹窗中的“保存路径”选择交互，将其从受限的下拉树形组件升级为“带前缀的自由输入框 + 独立大屏树形弹窗 + 底部内嵌新建操作”的复合结构。

**Architecture:** 
1. 将原有的 `el-tree-select` 替换为 `el-input`。使用插槽 (slots) 添加前缀（显示账号平台如 `[Quark]`）和后缀按钮（“选择目录”）。
2. 新增一个专门用于选择目录的 `el-dialog`，内部包含完整的 `el-tree` 组件，以提供更大的可视区域。
3. 在新弹窗的 `footer` 区域内嵌一个用于新建文件夹的小型表单（输入框+按钮），允许用户在当前选中的树节点下快速创建并立即刷新节点。

**Tech Stack:** Vue 3, Element Plus

---

### Task 1: 改造表单主界面的输入框

**Files:**
- Modify: `web/src/views/Tasks.vue`

- [ ] **Step 1: 替换现有的 el-tree-select**
  在 `<template>` 中，找到“保存路径”对应的 `el-form-item`。将其内部的 `path-input-group` 替换为带有前后缀插槽的 `el-input`。
  ```vue
        <el-form-item label="保存路径" required>
          <div class="path-input-group">
            <el-input 
              v-model="form.save_path" 
              placeholder="可手动输入或点击右侧选择" 
              class="save-path-input"
            >
              <template #prepend v-if="selectedAccountPlatform">
                [{{ selectedAccountPlatform }}]
              </template>
              <template #append>
                <el-button @click="openFolderDialog">选择目录</el-button>
              </template>
            </el-input>
          </div>
        </el-form-item>
  ```

- [ ] **Step 2: 添加账号平台前缀的计算属性**
  在 `<script setup>` 中，引入 `computed` 并创建一个计算属性，根据当前选中的 `form.account_id` 匹配出平台的中文名称。
  ```javascript
  import { ref, onMounted, onUnmounted, computed } from 'vue'
  // ...
  const selectedAccountPlatform = computed(() => {
    if (!form.value.account_id) return ''
    const account = accounts.value.find(acc => acc.id === form.value.account_id)
    if (account) {
      return account.platform === '139' ? '移动云盘' : 'Quark'
    }
    return ''
  })
  ```

- [ ] **Step 3: 调整样式**
  在 `<style scoped>` 中，为新的输入框添加样式，确保前缀部分具有足够的辨识度。
  ```css
  .save-path-input {
    width: 100%;
  }
  .save-path-input :deep(.el-input-group__prepend) {
    background-color: var(--el-fill-color-light);
    color: var(--el-text-color-regular);
    font-weight: bold;
  }
  ```

### Task 2: 实现独立大屏目录选择弹窗

**Files:**
- Modify: `web/src/views/Tasks.vue`

- [ ] **Step 1: 增加目录选择弹窗的模板**
  在 `<template>` 底部（如原有预览对话框上方）增加一个新的 `el-dialog`。
  ```vue
      <!-- 目录选择独立弹窗 -->
      <el-dialog 
        v-model="folderDialogVisible" 
        title="选择保存目录" 
        width="600px" 
        class="folder-dialog"
        append-to-body
        destroy-on-close
      >
        <div class="folder-tree-container">
          <el-tree
            ref="folderTreeRef"
            lazy
            :load="loadFolders"
            :props="{ label: 'label', isLeaf: 'isLeaf' }"
            node-key="path"
            highlight-current
            @current-change="handleTreeCurrentChange"
            :expand-on-click-node="false"
          >
            <template #default="{ node, data }">
              <span class="custom-tree-node">
                <span>{{ node.label }}</span>
              </span>
            </template>
          </el-tree>
        </div>
        
        <template #footer>
          <div class="folder-dialog-footer">
            <div class="create-folder-action">
              <el-input 
                v-model="newFolderName" 
                placeholder="在此目录下新建文件夹" 
                size="small" 
                @keyup.enter="handleInlineCreateFolder"
              >
                <template #append>
                  <el-button @click="handleInlineCreateFolder" :loading="creatingFolder">新建</el-button>
                </template>
              </el-input>
            </div>
            <div class="dialog-actions">
              <el-button @click="folderDialogVisible = false">取消</el-button>
              <el-button type="primary" @click="confirmFolderSelection">确认选择</el-button>
            </div>
          </div>
        </template>
      </el-dialog>
  ```

- [ ] **Step 2: 声明相关的响应式变量**
  在 `<script setup>` 中添加弹窗所需的控制变量。
  ```javascript
  // 独立目录弹窗相关
  const folderDialogVisible = ref(false)
  const folderTreeRef = ref(null)
  const selectedTreePath = ref('')
  const newFolderName = ref('')
  const creatingFolder = ref(false)
  ```

- [ ] **Step 3: 实现弹窗的打开与确认逻辑**
  编写打开弹窗时初始化状态的逻辑，以及点击确认时将选中的路径回填到主表单的逻辑。
  ```javascript
  const openFolderDialog = () => {
    if (!form.value.account_id) {
      return ElMessage.warning('请先选择执行账号')
    }
    // 尝试将当前输入框的路径作为默认选中（如果存在于树中）
    selectedTreePath.value = form.value.save_path || '/'
    newFolderName.value = ''
    folderDialogVisible.value = true
  }

  const handleTreeCurrentChange = (data) => {
    selectedTreePath.value = data.path
  }

  const confirmFolderSelection = () => {
    if (selectedTreePath.value) {
      form.value.save_path = selectedTreePath.value
    }
    folderDialogVisible.value = false
  }
  ```

### Task 3: 改造新建文件夹的交互逻辑

**Files:**
- Modify: `web/src/views/Tasks.vue`

- [ ] **Step 1: 实现内嵌的新建逻辑**
  使用新的 `handleInlineCreateFolder` 替换掉原先基于 `ElMessageBox.prompt` 的 `handleCreateFolder`。新逻辑将直接在选中的树节点下发送请求，并在成功后动态更新树组件。
  ```javascript
  const handleInlineCreateFolder = async () => {
    if (!newFolderName.value.trim()) {
      return ElMessage.warning('请输入文件夹名称')
    }
    
    // 如果没有选中树节点，默认在根目录创建
    const currentPath = selectedTreePath.value || '/'
    const currentID = pathIdMap.value[currentPath]

    if (currentID === undefined) {
      return ElMessage.warning('无法确定当前选中目录的 ID，请重新展开该节点')
    }

    creatingFolder.value = true
    try {
      const res = await createFolder(form.value.account_id, currentID, currentPath, newFolderName.value.trim())
      ElMessage.success('文件夹创建成功')
      
      // 更新映射表
      pathIdMap.value[res.path] = res.id
      
      // 动态向 ElTree 追加新节点
      if (folderTreeRef.value) {
        const currentNode = folderTreeRef.value.getNode(currentPath)
        if (currentNode) {
          folderTreeRef.value.append(res, currentNode)
          currentNode.expanded = true // 确保父节点展开以显示新节点
        }
      }
      
      // 自动选中新建的文件夹
      selectedTreePath.value = res.path
      if (folderTreeRef.value) {
         folderTreeRef.value.setCurrentKey(res.path)
      }
      
      // 清空输入框
      newFolderName.value = ''
    } catch (err) {
      console.error('创建文件夹失败:', err)
    } finally {
      creatingFolder.value = false
    }
  }
  ```

- [ ] **Step 2: 补充弹窗布局样式**
  在 `<style scoped>` 中添加弹窗内部的布局样式，特别是底部 Footer 的 Flex 布局。
  ```css
  .folder-tree-container {
    height: 400px;
    overflow-y: auto;
    border: 1px solid var(--el-border-color-light);
    border-radius: 4px;
    padding: 8px;
  }

  .folder-dialog-footer {
    display: flex;
    justify-content: space-between;
    align-items: center;
    width: 100%;
  }

  .create-folder-action {
    flex: 1;
    margin-right: 20px;
    max-width: 300px;
  }

  .dialog-actions {
    flex-shrink: 0;
  }
  ```

### Task 4: 验证与清理

**Files:**
- Modify: `web/src/views/Tasks.vue`

- [ ] **Step 1: 清理旧代码**
  确保彻底移除了旧的 `handleCreateFolder` 函数以及不再使用的导入。

- [ ] **Step 2: 测试全流程**
  1.  打开任务编辑弹窗。
  2.  观察“保存路径”输入框前是否有 `[Quark]` 或 `[移动云盘]` 标签。
  3.  手动在输入框内粘贴或修改一个复杂路径，测试是否能正常保存任务。
  4.  点击“选择目录”按钮，测试大屏弹窗是否能正常懒加载目录树。
  5.  选中一个子目录，在弹窗底部输入名字并点击“新建”，测试是否能原地创建并在树中动态显示。
  6.  点击“确认选择”，测试路径是否能正确回填到输入框。