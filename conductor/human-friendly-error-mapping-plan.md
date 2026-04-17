# 夸克/139 错误信息人性化清洗计划 (User-Friendly Error Mapping Plan)

## 1. Objective (目标)
彻底移除后端透传的底层技术错误（如 `Quark API HTTP error: 404, body: ...`），将其替换为通俗易懂的中文语言提示。

## 2. Key Files & Context (核心影响文件)
- `internal/core/quark/client.go`: 优化 `doRequest` 中的错误处理逻辑。
- `internal/core/cloud139/client.go`: 同步优化 139 的错误回传。

## 3. Implementation Steps (实施步骤)

### Step 3.1: 建立夸克错误码映射表
在 `internal/core/quark/client.go` 的 `doRequest` 中，针对非 200 响应执行以下逻辑：
1. 解析 JSON Body。
2. 提取 `code` 字段。
3. 根据 `code` 返回人性化文本：
   - `41010`: "分享文件包含违规内容，已被系统屏蔽。"
   - `24000`: "提取码不正确，请重新输入。"
   - `24001`: "该分享已失效，可能已被取消或删除。"
   - 其他 404: "分享链接已失效或链接格式不正确。"
   - 其他 400: "请求参数错误或非法访问。"

### Step 3.2: 移除技术前缀
无论 139 还是夸克，返回给前端的错误信息不再包含 "API error"、"HTTP error" 等词汇。
- **之前**: `Quark API 错误: 文件涉及违规内容`
- **现在**: `[Fatal] 该文件涉及违规内容，无法执行转存`

## 4. Verification & Testing (验证与测试)
- 测试违规链接：确认提示语变为“该分享文件涉及违规内容，已被系统屏蔽。”
- 测试错误提取码：确认提示语变为“提取码不正确，请重新输入。”
- 确认不再出现任何 JSON 字符串。
