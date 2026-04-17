# 移除状态列多余点号计划 (Remove Stray Dot in Status Column Plan)

## 1. Objective (目标)
移除在任务列表“状态”列中，当显示 “LINK ERROR” 时其右侧出现的异常点号（stray dot），恢复界面的纯净度。

## 2. Key Files & Context (核心影响文件)
- `web/src/views/Tasks.vue`
  - 第 31-54 行附近的 `el-table-column` (状态)。

## 3. Root Cause Analysis (根本原因分析)
- 在之前的代码合并或替换过程中，可能在 `</el-tooltip>` 或 `</el-tag>` 闭合标签之后意外留下了一个点号 `.`。
- 虽然在普通的 `read_file` 输出中难以察觉，但在浏览器渲染时会被识别为文本内容并显示在标签右侧。

## 4. Implementation Steps (实施步骤)
1. **执行精准替换**:
   - 将“状态”列的整个模板代码块替换为经过校准、不含任何多余字符的版本。
   - 特别检查 `</el-tooltip>` 和 `</template>` 之间的空白区域。

## 5. Verification & Testing (验证与测试)
1. 刷新前端页面。
2. 触发或查看一个带有致命错误（LINK ERROR）的任务，确认标签右侧的点号已彻底消失。
