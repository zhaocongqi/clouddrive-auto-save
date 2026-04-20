# 修复 139 移动云盘分享文件解析时间错误计划 (Fix 139 Share File Time Parsing)

## 1. Objective (目标)

修复在 `139 移动云盘` 中解析分享链接时，文件更新时间显示为 `0001-01-01 00:00:00` 的 Bug。

## 2. Key Files & Context (核心影响文件)

- `internal/core/cloud139/client.go`
  - 在 `ParseShare` 方法中，解析文件的 `udTime` 时，使用了错误的格式化字符串。

## 3. Root Cause Analysis (根本原因分析)

移动云盘 V6 分享接口返回的文件更新时间 `udTime` 通常是一个 14 位的字符串，格式为
`yyyyMMddHHmmss`（例如：`20260331123456`）。
在现有的代码逻辑中：

- 针对**文件夹**的解析，使用了正确的格式 `20060102150405` 进行时间转换，所以文件夹时间正常。
- 但针对**文件**的解析，代码中虽然判断了 `len(udTime) == 14`，却错误地使用带连字符的格式 `"2006-01-02
  15:04:05"`（长度19）去解析 14 位的字符串。这导致 `time.ParseInLocation` 抛出错误并返回时间零值，前端最终展现为
  `0001-01-01 00:00:00`。

## 4. Implementation Steps (实施步骤)

1. 打开 `internal/core/cloud139/client.go`。
2. 定位到 `ParseShare` 方法中解析文件 `coLst` 的部分。
3. 将错误的解析代码：

   ```go
   updateTime, _ = time.ParseInLocation("2006-01-02 15:04:05", udTime, cst)
   ```

   修改为正确的无间隔格式：

   ```go
   updateTime, _ = time.ParseInLocation("20060102150405", udTime, cst)
   ```

## 5. Verification & Testing (验证测试)

- 保存修改后，重新编译并运行服务。
- 前端刷新任务或者重新选择起始文件。
- 确认包含文件的分享链接解析出的文件时间是否恢复为正常的格式。
