# 网盘接口响应日志广播计划 (Cloud Drive Response Logging Plan)

## 1. 目标 (Objective)
在移动云盘 (139) 和夸克网盘 (Quark) 的底层请求客户端中，全局拦截所有的 API 响应。将网盘服务器返回的数据通过 `utils.GlobalBroadcaster.Broadcast` 广播至实时日志流中，并在 Dashboard 统一展示。同时，移除旧版代码中散落的零星调试日志，保证日志规范一致。

## 2. 涉及文件 (Key Files & Context)
- `internal/core/cloud139/client.go`
- `internal/core/quark/client.go`
- 依赖项：`github.com/zcq/clouddrive-auto-save/internal/utils`

## 3. 实施步骤 (Implementation Steps)

### 3.1 修改移动云盘客户端 (`cloud139/client.go`)
- **引入依赖**：确保在 `import` 中添加 `"github.com/zcq/clouddrive-auto-save/internal/utils"`。
- **全局拦截**：在 `doRequest` 方法中，读取完 `resp.Body` 之后，新增日志广播代码。解析 `apiURL` 以提取干净的接口路径，然后构造标准的日志字符串：
  ```go
  u, _ := url.Parse(apiURL)
  apiPath := apiURL
  if u != nil {
      apiPath = u.Path
  }
  msg := fmt.Sprintf("[139 Debug] 接口 %s 响应: %s", apiPath, string(respBody))
  log.Printf(msg)
  utils.GlobalBroadcaster.Broadcast(msg)
  ```
- **清理冗余日志**：在 `GetInfo` 等方法中找到并移除 `log.Printf("[139 Debug] getUser 原始响应: %s", string(resp))` 等重复的打印代码。

### 3.2 修改夸克客户端 (`quark/client.go`)
- **引入依赖**：确保在 `import` 中添加 `"github.com/zcq/clouddrive-auto-save/internal/utils"`。
- **全局拦截**：在 `doRequest` 方法中，读取完 `resp.Body` 之后，新增相似的日志广播代码。解析 `fullURL` 的 `Path`：
  ```go
  u, _ := url.Parse(fullURL)
  apiPath := fullURL
  if u != nil {
      apiPath = u.Path
  }
  msg := fmt.Sprintf("[Quark Debug] 接口 %s 响应: %s", apiPath, string(respBody))
  log.Printf(msg)
  utils.GlobalBroadcaster.Broadcast(msg)
  ```
- **清理冗余日志**：在 `ParseShare` 等方法中找到并移除 `log.Printf("[Quark Debug] ParseShare 详情原始响应: %s", string(resp))` 等重复的打印代码。

## 4. 验证与测试 (Verification & Testing)
1. 编译并运行服务端。
2. 在 UI 中手动触发“账号有效性校验”或使用 Postman 发起对后端的请求。
3. 观察后端控制台的终端输出是否正常打印格式化的 `[139 Debug]` 或 `[Quark Debug]` 日志。
4. 观察 Dashboard 的实时日志流，确认网盘的原始响应均能通过 SSE 正确传达到前端展示。