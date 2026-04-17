# 重命名感知的起始转存点选择优化计划 (Rename-Aware Checkpoint Selection Plan)

## 1. Objective (目标)
优化“选择起始转存文件”的功能，使其能够感知当前的重命名规则。系统将基于重命名后的文件名执行同名检查，并在界面上同时展示原始名与预览名，实现“所见即所得”的选择体验，确保同名提示的准确性。

## 2. Key Files & Context (核心影响文件)
- **核心模型**: `internal/core/drive.go`
  - 在 `FileInfo` 中新增 `NewName` 字段。
- **API 控制器**: `internal/api/task_preview.go`
  - 增强 `parseShareLinkInfo`，集成重命名处理器。
- **前端视图**: `web/src/views/Tasks.vue`
  - 传递重命名参数，并重构表格以展示双重名称。

## 3. Implementation Steps (实施步骤)

### Step 3.1: 后端模型与逻辑重构
1. **修改 `internal/core/drive.go`**:
   - 在 `FileInfo` 中添加：`NewName string `json:"new_name"``。
2. **升级 `internal/api/task_preview.go`**:
   - 接收参数增加：`Pattern`, `Replacement`, `Name` (任务名)。
   - 在获取分享列表后，若存在重命名规则，调用 `renamer` 引擎为每个文件生成 `NewName`。
   - **同名检查升级**：在对比目标目录时，优先使用 `NewName` 进行比对，确保“已同步”标签反映的是最终入库后的状态。

### Step 3.2: 前端交互升级
1. **修改 `web/src/views/Tasks.vue`**:
   - `openStartFileDialog` 调用时，将表单中的 `pattern`, `replacement`, `name` 一并传给后端。
2. **重构选择表格列**:
   - 调整“文件名”列，主标题显示 `NewName`（如果存在），副标题显示 `Original Name`（灰色小字）。
   - 确保“状态”列基于重命名后的结果显示。

## 4. Verification & Testing (验证与测试)
1. **编译测试 (关键)**: 运行 `go build` 确保核心逻辑、驱动、API 层无任何语法错误。
2. **逻辑验证**: 设置一个重命名规则（如加后缀），确认弹窗中是否显示了新名字，且已存在于网盘的同名文件（基于新名）被正确置灰。
3. **鲁棒性验证**: 确认不设置任何正则规则时，系统行为退化为原始的文件名比对，保持兼容性。
