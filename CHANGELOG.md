# 更新日志 (Changelog)

## [Unreleased]

（暂无）

## [1.6.1] - 2026-06-09

### 🔧 维护与部署 (Chore)

- **Docker 镜像标签统一**：镜像标签从 `1.x.x` 改为 `v1.x.x` 格式，与 git tag 保持一致。
- **CI 依赖版本统一**：三个 Docker 发布工作流的 `docker/metadata-action` 统一升级到 v6。
- **版本注入修复**：Docker 构建时注入二进制的 VERSION 去掉 `v` 前缀，与 GoReleaser 行为对齐，修复前端显示 `vv1.5.0` 双前缀问题。
- **README 版本同步**：Docker 示例镜像版本从 `v1.1.1` 更新为 `v1.6.0`。

## [1.6.0] - 2026-06-08

### ✨ 核心特性 (Features)

- **API 静态 Token 认证**：新增 `UCAS_API_KEY` 环境变量，支持 `X-API-Key` 请求头和 `?token=` 查询参数认证，为空则跳过（向后兼容）。
- **凭据 AES-GCM 加密存储**：新增 `UCAS_SECRET_KEY` 环境变量（64 字符 hex），Cookie/AuthToken 使用 AES-256-GCM 加密后存入数据库，启动时自动迁移明文凭据。
- **SSE Pinia Store 统一管理**：新建 `web/src/stores/sse.js`，Dashboard/Tasks/Search 三视图共享单一 EventSource 连接，支持引用计数和指数退避重连。
- **插件配置更新 API**：`PUT /api/plugins/:name/config` 从 TODO stub 改为完整实现，配置持久化到 Setting 表。
- **任务重试机制**：失败任务自动分类（致命/可恢复），可恢复错误按指数退避重试（30s→60s→120s→...→最大 3600s），支持配置最大重试次数和忽略后缀去重。

### 🔧 重构 (Refactoring)

- **Task 输入 DTO 隔离**：新增 `taskInputDTO` 白名单结构体，`createTask`/`updateTask` 绑定 DTO 而非完整 `db.Task`，防止客户端篡改运行时状态字段。补全 `filter`/`run_days`/`start_date` 三个缺失字段。
- **Account 凭据脱敏**：所有返回账号信息的端点统一通过 `toAccountDTO()` 输出，Cookie/AuthToken 不再出现在 API 响应中。
- **前端组件拆分**：Tasks.vue（1993→438 行）拆分为 TaskTable/TaskForm/utils 子组件；Settings.vue（1112→145 行）拆分为 Schedule/Notify/Plugin/Search 子组件。
- **API 响应格式统一**：全部端点统一为 `c.PureJSON`，消除 `c.JSON` 的 HTML 转义不一致问题。
- **前端 API 标准化**：3 处原生 `fetch` 替换为 axios 统一调用，request.js 拦截器注入 API Key。
- **预定义魔法匹配**：新增 `$TV`、`$BLACK_WORD`、`$SHOW_MAGIC`、`$TV_MAGIC` 四种预定义正则规则，任务中直接用 `$名称` 引用。

### 🛡️ 安全 (Security)

- **Docker 非 root 运行**：Dockerfile 添加 `ucas` 用户（UID 1000），容器不再以 root 身份运行。
- **Mock 文件构建隔离**：4 个 mock 文件添加 `//go:build e2e` 标签，不编入生产二进制。
- **CI latest 标签修复**：`docker-publish-main.yml` 的 `latest` 改为 `main`，避免与 tag 发布冲突。

### 🐛 修复 (Bug Fixes)

- **Broadcaster 竞态修复**：`Shutdown()` 与 `run()` 发送循环的竞态条件修复，通过 `defer closeAllClients()` 避免 send to closed channel panic。
- **Worker context 及时释放**：4 处 `defer cancel()` 改为对应代码块结束后立即调用。
- **Notify Manager 锁优化**：`Send`/`Test` 方法在 RLock 内只做快照，释放锁后再执行网络调用。
- **quark json.Unmarshal 错误处理**：`SaveLink` 中 `json.Unmarshal` 添加错误检查。
- **正则预编译**：cloud139 和 pansou 中 14 处函数内 `regexp.MustCompile` 提取为包级变量。
- **IsEncrypted 误判修复**：从简单 `:` 检查改为验证 `base64(12B nonce):base64(ciphertext)` 格式。
- **SSE addHistoryLogs 覆盖修复**：从替换改为追加，避免丢失实时日志。
- **ScheduleSettings 键污染修复**：`fetchScheduleSettings` 只合并白名单内的 key，避免引入 bark_* 等无关 key。
- **Makefile test 去重**：合并为单次运行，消除重复测试。
- **drivers map 并发保护**：添加 `sync.RWMutex`。

### 🎨 界面优化 (UI/UX)

- **Dashboard 日志时间戳**：日志从纯字符串改为 `{text, time}` 对象格式，显示实际接收时间而非渲染时间。
- **前端死代码清理**：删除未使用的 echarts 图表组件。
- **vite 路径别名**：添加 `@` 路径别名支持。

### 🔧 重构 (Refactoring)

- **API 响应格式统一**：全部端点统一为扁平格式（直接返回业务数据），消除了信封格式 `{code, data}` 与扁平格式共存的不一致问题。Settings.vue 中 4 处绕过 request.js 的 raw fetch 调用改为统一使用 axios 实例。
- **预定义魔法匹配**：新增 `$TV`、`$BLACK_WORD`、`$SHOW_MAGIC`、`$TV_MAGIC` 四种预定义正则规则，任务中直接用 `$名称` 引用。
- **链接验证 API**：新增 `GET /api/search/validate` 端点，自动识别夸克/移动云盘平台并验证分享链接有效性。
- **运行星期配置**：任务新增 `run_days` 字段，支持按星期过滤运行日（1=周一, 7=周日），比手写 Cron 更直观。
- **忽略后缀去重**：任务新增 `ignore_extension` 开关，01.mp4 和 01.mkv 视为同一文件避免重复转存。
- **Tasks 页面状态筛选**：新增全部/等待中/运行中/成功/失败筛选 Tab + 搜索框。
- **全局快捷键**：Ctrl+S 保存任务、Ctrl+R 运行所有任务、未保存修改关闭拦截。
- **CloudSaver & PanSou 搜索集成**：
  - **CloudSaver 搜索源**：实现 JWT Token 认证机制，支持自动登录续期，搜索结果自动清洗（提取夸克链接、解析标题描述、时间格式化）。
  - **PanSou 搜索源**：实现免认证搜索，支持按网盘类型过滤和结果合并去重，`note` 字段自动解析为标题和描述。
  - **搜索配置管理**：新增 `/api/search/config` API（GET/PUT），支持从 Setting 表加载/保存搜索源配置，运行时热更新。
  - **统一搜索客户端**：重构 `Client` 支持配置化源创建、结果去重（按 shareurl）和排序（按时间降序）。
  - **前端适配**：Search.vue 新增标签和频道显示，Settings.vue 新增搜索源配置 Tab。
  - **API 文档**：新增 `docs/cloudsaver_api.md` 和 `docs/pansou_api.md`，详细记录两个搜索源的接口规范和调用流程。

### 🎨 界面优化与重构 (UI/UX Refactoring)

- **科技暗黑毛玻璃风格升级**：全局部署了 Glassmorphism 样式和霓虹微光设计语言，对 Element Plus 全局组件进行了发光微光样式重塑。
- **控制台 (Dashboard) 重构**：改为左右分栏的极客主控面板，整合了实时活跃队列与历史时间线，保证了实时刷新下历史重试操作能被正常触发。
- **任务管理 (Tasks) 交互优化**：将原本繁重的全屏创建/编辑 Dialog 重构为右侧滑入的 `<el-drawer>` 抽屉表单；实装了基于正则的智能粘贴解析器，支持从复合文案中一键提取网盘分享链接和提取码。
- **账号管理 (Accounts) 空间可视化**：改版为精美的卡片网格布局，并将原本的数据表格容量展现形式升级为直观的可视化 `el-progress type="circle"` 圆环空间配额图。
- **页面路由精简**：废弃了独立的通知和插件配置页面，导航栏精简为五大核心模块，并将调度配置、推送通道和扩展插件集中整合至新的 Tab 式 [Settings.vue](file:///home/zcq/Github/clouddrive-auto-save/web/src/views/Settings.vue) 页面中。

### 🛡️ 稳定性 (Stability)

- **E2E 精准定位器重构**：修复了在嵌套 `el-tabs` 复杂 DOM 树下由于重名“保存/测试”按钮以及同级 `el-switch` 导致的定位歧义超时问题。为各通道表单添加了专属 class（如 `.bark-form`, `.wechat-form` 等）并优化了 Playwright 选择器，使 74 个 E2E 集成测试用例 100% 顺利通过。
- **E2E 测试大幅增强**：先前新增了 26 个测试（48 → 74），覆盖仪表盘 FAB 导航、账号编辑、OpenList 配置、Bark 高级设置、侧边栏导航、暗色模式、提取码、表单验证等场景。
- **修复 Dashboard 测试 flaky**：Dashboard 的 SSE 长连接会接收真实后端事件触发 `fetchStats()` 与 `page.route` mock 竞态，通过同时 mock SSE 端点 (`/api/dashboard/logs`) 阻断干扰。


## [1.3.0] - 2026-05-18

### ✨ 核心特性 (Features)

- **分享链接子目录浏览**：
  - **子目录导航**：选择文件弹窗支持点击文件夹进入子目录，面包屑导航返回上级。
  - **139 平台目录记忆**：新增 `share_parent_id` 字段，移动云盘选择子目录后再次打开自动定位。
  - **浏览分享内容**：新增"浏览分享内容并选择目录"按钮，可将子目录设为新的分享链接起始点。
- **任务执行增强**：
  - **parse_share API** 新增 `parent_id` 参数，支持浏览子目录内容。
  - **起始文件选择**：选择子目录后，起始转存文件弹窗自动以该目录为根。

### 🎨 界面优化 (UI)

- **按钮布局重构**：分享链接的"浏览内容"和"打开链接"按钮从 input append 插槽改为独立 flex 行布局，解决对齐和重叠问题。
- **模式差异交互**：浏览分享内容模式显示"进入"按钮；选择起始文件模式在根目录显示 radio 选择列。

### ⚡ 性能优化 (Performance)

- **E2E 测试提速 80%**：引入 storageState setup project 预置账号，移除 6 处硬编码 `waitForTimeout(5000)`，workers 从 1 提升至 4。29 个测试从 ~80s 降至 ~16s。

### 🛡️ 稳定性 (Stability)

- **E2E 测试覆盖**：新增 3 个分享链接子目录浏览 E2E 测试用例。
- **单元测试**：新增 `share_parent_id` 持久化的 API 单元测试。
- **Mock 增强**：夸克和 139 的 HTTP Mock 支持子目录路由。

## [1.1.1] - 2026-05-06

### ✨ 核心特性 (Features)

- **全局视觉体验升级**：
  - **图标重构**：将全站图标重构为“层叠云朵”设计，使用原生 SVG 替换静态资源。不仅更新了系统 Favicon，还在侧边栏左上角集成了带有平滑呼吸灯动画的新 `CloudLogo` 组件。
- **系统状态与引导增强**：
  - **版本与源码导航**：新增侧边栏底部组件 (`SidebarFooter`)，支持显示当前系统版本、自动检测 GitHub 新版本发布，并提供快捷源码库跳转。
  - **API 版本端点**：后端引入 `/api/version` 路由并修复了构建版本信息的注入机制，方便前端获取运行时数据。

### 📚 文档完善 (Documentation)

- **API 手册拆分**：随着网盘底层驱动的深入支持，将原本的整合版 API 手册拆分为独立且更加详尽的《移动云盘 (139) API 手册》与《夸克网盘 (Quark) API 手册》，极大提升了二次开发可读性。

## [1.1.0] - 2026-05-05

### ✨ 核心特性 (Features)

- **Bark 通知深度进化**：
  - **高级配置支持**：新增推送级别（关键/即时/间歇/被动）、自定义铃声、图标及自动归档设置。
  - **智能批量汇总**：引入 `BatchTracker` 追踪器与 `SendBatchNotification` 逻辑，在“执行所有任务”时将多次通知自动压缩为单条汇总消息，彻底告警“消息轰炸”。
  - **数据竞争修复**：深度重构通知模块，消除并发场景下的数据竞争风险。
- **任务引擎效能提升**：
  - **耗时统计分析**：所有转存任务均集成耗时统计功能，并随通知同步展示。
  - **原生批量模式支持**：Worker 层支持 Batch 模式，通过动态生成的 `BatchID` 实现任务流的精确追踪。
- **调度系统交互重构**：
  - **双模式 Cron 配置**：设置页面支持“简易/高级”双模式切换，满足不同用户的定时需求，并增强了表达式校验与持久化能力。
  - **UI 布局优化**：精细化调整设置页面头部布局，提升视觉一致性。

### 🛡️ 稳定性与安全 (Stability)

- **网盘驱动健壮性**：
  - **夸克异常拦截**：增强对失效或空分享链接的识别能力，防止误判为转存成功。
- **系统安全性**：
  - **ID 碰撞防御**：在 `BatchID` 生成算法中引入随机后缀，彻底消除高并发下的 ID 碰撞风险。

### 🔧 维护与部署 (Chore)

- **开发协作增强**：新增 `CLAUDE.md` 开发协作指南，为 AI 助手提供标准化的项目上下文与命令规范。
- **文档体系完善**：补充 Bark API 详细说明文档，并在 UI 中集成 Bark 配置教程链接。
- **质量管控升级**：统一全局 Go 代码格式 (`gofmt`)，优化 Markdown 校验规则（自动忽略 `node_modules`）。

## [1.0.0] - 2026-04-29

### ✨ 核心特性 (Features)

- **多平台深度驱动**：整合了移动云盘 (139) 与夸克网盘 (Quark) 的转存、管理及重命名能力。
- **智能化整理与重命名**：内置多种魔法变量（`{YEAR}`, `{DATE}`, `{OLDNAME}` 等）与正则重命名模板，支持转存后自动更名。新增 `{OLDNAME}` 魔法变量。
- **交互体验升级**：任务管理界面账号选择增加平台分组；增加分享链接快捷外跳功能，并支持自动提取、拼接并复制提取码至剪贴板，极大简化未安装网盘脚本插件时的用户操作。
- **高并发调度引擎**：基于 Go Worker Pool 实现任务并发处理，支持秒级 Cron 定时规则。
- **现代化仪表盘**：全响应式设计，集成实时进度监控、全站状态同步及历史日志检索。
- **结构化分级日志**：引入 `log/slog`，支持 `LOG_LEVEL` 动态过滤，并将日志实时同步至前端仪表盘。
- **一键批量控制**：支持一键运行所有健康任务，自动跳过 [Fatal] 任务与运行中任务。

### 🛡️ 稳定性与安全 (Stability)

- **网盘异常处理强化**：深度优化网盘接口报错逻辑。修复 139 移动云盘错误码 float64 解析导致的异常；针对夸克网盘缺失提取码 (`41008`)、提取码错误 (`41007`, `24000`, `41009`) 等细分场景提供精确的中文 `[Fatal]` 致命错误拦截与 UI 警示，防止任务陷入死循环。
- **全链路测试体系**：引入 Playwright 自动化测试框架，构建并完善了包含各种边缘场景动态 HTTP Mock 在内的 E2E 测试，大幅提升发布质量。
- **自动状态自愈**：服务端启动时自动检测并重置异常中断的任务，防止死锁。
- **依赖注入重构**：核心模块解耦，显著提升代码可测试性。
- **敏感数据脱敏**：控制台与日志流对 `Cookie`、`Token` 等敏感信息进行了自动屏蔽/等级控制。

### 🔧 维护与部署 (Chore)

- **Docker 支持**：提供多阶段构建的轻量化 Docker 镜像方案。
- **自动化流水线**：集成 GitHub Actions CI 与 GoReleaser 自动发布流程。
- **Makefile 增强**：内置 `check`, `test-html`, `lint`, `e2e-test` 等质量管控工具。
