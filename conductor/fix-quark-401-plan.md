# 修复夸克分享解析错误双重提示计划 (Fix Quark Share Parsing Error Tips Plan)

## 1. Objective (目标)
修复在编辑任务选择起始文件时，如果夸克分享链接违规或过期，前端会同时弹出两根错误提示（一个是冗长丑陋的后端 JSON 报错，另一个是固定的“解析失败”提示）的体验问题。

## 2. Key Files & Context (核心影响文件)
- `internal/core/quark/client.go`: 底层请求封装 `doRequest`，目前在遇到非 200 响应时直接将原始的 JSON 字符串返回。
- `web/src/views/Tasks.vue`: 前端的 `openStartFileDialog` 异常处理模块，目前会在捕获错误后手动弹出一个 `ElMessage.error`。

## 3. Implementation Steps (实施步骤)

### Step 3.1: 后端驱动层错误信息清洗
1. 打开 `internal/core/quark/client.go`。
2. 定位到 `doRequest` 方法的 HTTP 状态码非 200 处理逻辑处：
   ```go
   if resp.StatusCode != http.StatusOK {
        ...
   }
   ```
3. 在此逻辑内，增加 JSON 解析步骤。尝试从响应 Body 中提取出真实的 `message` 字段。
4. 如果提取成功，且 HTTP 状态码为 400 或 404（通常是链接失效、违规），将其格式化为 `[Fatal] 夸克接口报错: {message}` 的简短文本并返回。否则返回原始状态码错误。

### Step 3.2: 前端 UI 冗余提示清理
1. 打开 `web/src/views/Tasks.vue`。
2. 定位到 `openStartFileDialog` 中的 `catch (err)` 代码块。
3. 删除硬编码的 `ElMessage.error('解析分享链接失败，请检查链接或账号状态')`，因为前端 `request.js` 拦截器已经会根据后端返回的精简错误信息进行弹窗提示，保留双重弹窗显得多余。

## 4. Verification & Testing (验证与测试)
- 测试违规/失效的夸克链接：打开任务编辑点击“选择文件”。
- 预期结果：前端只弹出一个红色的报错框，内容类似 `[Fatal] 夸克接口报错: 文件涉及违规内容`，不再显示原始的长串 JSON，也不再有第二个重复的警告。