# 存储空间展示优化设计方案

## 1. 目标 (Goal)

优化账号管理页面中的存储空间进度条视觉效果，并修复在已用容量超过总容量时出现的 "NaN undefined" 显示错误。

## 2. 视觉改进：状态感知进度条 (方案 B)

进度条将不再使用固定的蓝紫渐变，而是根据当前已用百分比动态切换色彩，增强视觉语义。

- **正常状态 (< 70%)**: 使用绿色渐变 `linear-gradient(90deg, #10b981 0%, #34d399 100%)`。
- **预警状态 (70% - 90%)**: 使用黄色渐变 `linear-gradient(90deg, #f59e0b 0%, #fbbf24 100%)`。
- **危险/超额状态 (>= 90%)**: 使用红色渐变 `linear-gradient(90deg, #ef4444 0%, #f87171 100%)`。

## 3. 逻辑修复：容量文本展示

修复 `Accounts.vue` 中容量剩余量的计算与展示逻辑。

- **计算公式**: `remaining = capacity_total - capacity_used`。
- **展示策略**:
  - **未超额 (remaining >= 0)**: 显示文本 `剩 {{ formatBytes(remaining) }}`，颜色为绿色。
  - **已超额 (remaining < 0)**: 显示文本 `已超额 {{ formatBytes(Math.abs(remaining)) }}`，颜色为红色。
- **进度条百分比**: 传递给 `el-progress` 的 `percentage` 值需通过 `Math.min(100, percentage)` 进行截断，确保最大为 100%。

## 4. 涉及文件

- `web/src/views/Accounts.vue`: 修改模板中的文本显示逻辑、CSS 样式以及 `calcPercentage` / `getCapacityStatus` 辅助函数。
