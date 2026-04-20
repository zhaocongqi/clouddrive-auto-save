# 文档同步更新计划 (Docs Sync Plan)

**Goal:** 更新 `README.md`、`docs/api/tasks.md` 和 `docs/cloud_drive_apis.md`，使其准确反映近期对转存核心逻辑（时间排序、起始点截断、正则过滤、基于新名去重）以及魔法变量的使用说明。

---

### Task 1: 更新 Tasks API 文档 (docs/api/tasks.md)

**Files:**
- Modify: `docs/api/tasks.md`

**Changes:**
1. 重写 `8. 自动重命名与去重逻辑说明 (Workflow)` 章节，详细描述：
   - **时间降序排列**: 系统强制按更新时间从新到旧排序。
   - **双重过滤机制 (AND 关系)**: 明确说明系统仅转存 **“在起始文件之后（含）更新”** 且 **“符合正则匹配规则”** 的文件。
   - **预测式智能去重**: 描述基于重命名后新名字的预检逻辑。
2. 增加 **魔法变量 (Magic Variables)** 章节，详细列出 `{TASKNAME}`, `{YEAR}`, `{DATE}`, `{CHINESE}`, `{EXT}` 的提取规则和作用。

### Task 2: 更新主 README 文档 (README.md)

**Files:**
- Modify: `README.md`

**Changes:**
1. 新增 **“魔法变量与重命名示例”** 章节，提供直观的示例（如：`[{DATE}] {TASKNAME}.{EXT}`）。
2. 在核心特性或使用说明中，强调起始点选择与正则匹配之间的 **“与 (AND)”** 关系。

### Task 3: 更新 Cloud Drive APIs 文档 (docs/cloud_drive_apis.md)

**Files:**
- Modify: `docs/cloud_drive_apis.md`

**Changes:**
1. 完善夸克网盘 `2.4 分享与转存接口 (PC端)` 中的 `执行保存 (save)` 和 `查询异步任务 (task)` 说明。
2. 明确标注转存是异步行为，返回 `task_id`，且系统已实现自动轮询机制。

### Task 4: 提交更改

- 本地提交，不执行推送。
