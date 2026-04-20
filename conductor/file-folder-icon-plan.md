# 区分文件与文件夹图标 UI 优化计划 (File and Folder Icon UI Optimization Plan)

> **For agentic workers:** REQUIRED SUB-SKILL: Use
superpowers:subagent-driven-development (recommended) or
superpowers:executing-plans to implement this plan task-by-task. Steps use
checkbox (`- [ ]`) syntax for tracking.

**Goal:** 在任务管理界面的“起始转存点”选择表格中，通过直观的图标区分当前项是文件还是文件夹。

**Architecture:**

1. 依赖底层返回的 `is_folder` 字段。
2. 在 `web/src/views/Tasks.vue` 中从 `lucide-vue-next` 引入 `Folder` 和 `File` 图标。
3. 修改显示“从该文件开始”的 `<el-table-column>`，利用 `<template #default="{ row }">` 判断
   `row.is_folder` 来渲染对应的图标和颜色。

**Tech Stack:** Vue 3, Element Plus, Lucide Icons

---

## Task 1: 引入图标并修改列模板

**Files:**

- Modify: `web/src/views/Tasks.vue`

- [ ] **Step 1: 引入图标组件**
  在 `<script setup>` 中，从 `lucide-vue-next` 中增加 `Folder` 和 `File` 两个图标的引入。

- [ ] **Step 2: 自定义文件名列模板**
  在 `<template>` 的分享文件表格中，找到 `prop="name"` 的 `<el-table-column>`，将其改为使用插槽的格式：

  ```vue
            <el-table-column label="从该文件开始 (含)" show-overflow-tooltip>
              <template #default="{ row }">
                <div style="display: flex; align-items: center; gap: 6px">
                  <el-icon size="16">
                    <Folder v-if="row.is_folder" color="#eab308" />
                    <File v-else color="#64748b" />
                  </el-icon>
                  <span>{{ row.name }}</span>
                </div>
              </template>
            </el-table-column>
  ```

### Task 2: 验证与测试

- [ ] **Step 1: 验证展示效果**
  重新打开或刷新前端任务创建/编辑界面，输入一个包含文件夹和文件的分享链接。观察表格中是否能正确根据类型显示金黄色的文件夹图标和灰色的文件图标。
