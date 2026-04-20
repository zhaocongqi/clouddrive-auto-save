# 系统设置 (Settings)

## 1. 获取全局调度设置

获取系统当前的全局定时触发规则及开关状态。

- **URL**: `/settings/schedule`
- **Method**: `GET`
- **Response**:

```json
{
  "enabled": true,
  "cron": "0 0 2 * * *"
}
```

---

## 2. 更新全局调度设置

修改全局调度规则。修改后系统会立即更新后台调度引擎。

- **URL**: `/settings/schedule`
- **Method**: `POST`
- **Payload**:
| 字段 | 类型 | 说明 |
| :--- | :--- | :--- |
| `enabled` | bool | 全局总开关 |
| `cron` | string | 全局 Cron 表达式 (例如 `0 0 2 * * *` 代表每天凌晨 2 点) |

---

## 3. 调度优先级逻辑 (Logic Priority)

- **跟随全局 (Global)**: 任务执行受全局开关 `enabled` 与全局 `cron` 的共同约束。
- **自定义 (Custom)**: 任务拥有独立的 Cron 表达式，不受全局配置影响。
- **手动 (Off)**: 系统不自动执行，仅在用户点击 UI 上的“运行”按钮时触发。
