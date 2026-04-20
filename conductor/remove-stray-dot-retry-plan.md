# 再次尝试移除状态列点号计划 (Retry Remove Stray Dot Plan)

## 1. Objective (目标)

彻底根除任务列表“状态”列中 “LINK ERROR” 右侧顽固存在的异常点号。

## 2. Key Files & Context (核心影响文件)

- `web/src/views/Tasks.vue`
  - 状态列模板区域。

## 3. Root Cause Analysis (根本原因分析)

- 尽管进行了多次清理，但点号可能隐藏在标签闭合处的换行或不可见空白字符中。
- 之前的 fuzzy match 修正可能未能完全覆盖所有字符。

## 4. Implementation Steps (实施步骤)

1. **重写整个 el-table-column**:
   - 采用最简化的标签闭合方式。
   - 移除所有可能产生歧义的空格。
   - 修正 `div` 标签的闭合缩进。

## 5. Verification & Testing (验证与测试)

- 强制刷新浏览器缓存。
- 确认点号消失。
