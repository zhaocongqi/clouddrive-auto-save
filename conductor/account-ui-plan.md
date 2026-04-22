# 账号管理 UI 重新设计实施计划 (Account Management UI Redesign Plan)

> **For agentic workers:** REQUIRED SUB-SKILL: Use
superpowers:subagent-driven-development (recommended) or
superpowers:executing-plans to implement this plan task-by-task. Steps use
checkbox (`- [ ]`) syntax for tracking.

**Goal:** 优化账号管理页面的 UI。在保持整体表格形态的基础上进行深度美化，并增加渐变进度条及剩余空间提示。保留模态弹窗交互但优化字段布局和说明。

**Architecture:**

1. **优化版表格 (Enhanced Table)**：对 `el-table` 的列宽、内边距进行调整，可能加入平台图标（Icon）以代替简单的 Tag。
2. **渐变水平进度条 + 剩余提示**：修改“存储空间”列的自定义模板，将原始容量信息计算出剩余空间，并使用自定义 CSS 样式的渐变色覆盖
   `el-progress`。
3. **模态弹窗 (Dialog)**：优化添加/编辑账号的对话框排版，将帮助信息更好地融入 UI，调整多行文本框的大小和样式。

**Tech Stack:** Vue 3, Element Plus, CSS

---

## Task 1: 优化表格展示结构与样式

**Files:**

- Modify: `web/src/views/Accounts.vue`

- [ ] **Step 1: 引入图标并美化平台列**
  - 从 `lucide-vue-next` 中引入 `HardDrive`、`Cloud` 等图标用于区分不同云盘，或者直接用好看的 Tag 样式。
  - 在 `Accounts.vue` 的 `<script setup>` 中补充相关 Icon。
  - 在 `<template>` 中，修改“平台”列的展示。

- [ ] **Step 2: 重新设计容量展示列**
  - 在“存储空间”列，计算剩余空间。
  - 添加显示文字：例如 `已用 10GB / 总共 100GB (剩 90GB)`。
  - 为 `el-progress` 添加特定的 class（如 `.gradient-progress`），稍后在 CSS 中实现渐变效果。

## Task 2: 实现进度条渐变及全局样式美化

**Files:**

- Modify: `web/src/views/Accounts.vue`

- [ ] **Step 1: 编写容量文字和渐变进度条样式**
  - 在 `<style scoped>` 中增加进度条渐变效果。可以通过穿透 `::v-deep(.el-progress-bar__inner)`
    来修改背景色为 `linear-gradient(...)`。
  - 增加文字排版的 Flex 布局，使“已用/总共”和“剩余”分别靠两边对齐。

- [ ] **Step 2: 优化表格和页面整体边距**
  - 调整 `.table-card` 的阴影和圆角，使其更加现代化。
  - 为不同列的数据（如时间、状态）增加合适的颜色。

## Task 3: 优化编辑/添加弹窗的表单布局

**Files:**

- Modify: `web/src/views/Accounts.vue`

- [ ] **Step 1: 重排表单项**
  - 将弹窗中的说明文字（help 属性）改为 `el-alert` 或插槽形式的表单提示，更加显眼。
  - 对 139 特有的输入字段（Authorization 和 Cookie）的说明文案进行清晰化处理，使用分割线或标签页。

- [ ] **Step 2: 确保逻辑一致性**
  - 检查弹窗关闭和保存的交互，确保 loading 状态正确显示。

## Task 4: 测试与验证

- [ ] **Step 1: 验证表格展示**
  启动前端环境，观察账号列表：
  1. 平台列是否清晰美观。
  2. 进度条是否呈现渐变色，并且剩余空间显示正确。
  3. 各列宽度是否协调。

- [ ] **Step 2: 验证表单交互**
  点击添加/编辑按钮，验证弹窗布局是否比之前更整洁、易读。
