# 账号管理 (Accounts)

## 1. 获取账号列表
获取系统中所有已绑定的云盘账号。

- **URL**: `/accounts`
- **Method**: `GET`
- **Response**: `Array<Account>`

### Account 对象结构
| 字段 | 类型 | 说明 |
| :--- | :--- | :--- |
| `id` | uint | 账号唯一标识 |
| `platform` | string | `139` 或 `quark` |
| `nickname` | string | 云盘真实昵称 |
| `vip_name` | string | 会员等级（如：白银会员、SVIP） |
| `capacity_used`| int64 | 已用容量 (Bytes) |
| `capacity_total`| int64 | 总容量 (Bytes) |
| `status` | int | 状态 (1:正常, 0:失效) |
| `last_check` | string | 最后一次校验的时间戳 |

---

## 2. 添加新账号
绑定一个新的移动云盘或夸克网盘账号。

- **URL**: `/accounts`
- **Method**: `POST`
- **Payload**:
| 字段 | 类型 | 必填 | 说明 |
| :--- | :--- | :--- | :--- |
| `platform` | string | 是 | `139` 或 `quark` |
| `account_name`| string | 是 | 备注名或手机号 |
| `cookie` | string | 否 | 浏览器全量 Cookie |
| `auth_token` | string | 否 | 仅 139 支持，抓包获取的 Basic 串 |

---

## 3. 更新账号信息
修改已存在的账号配置。

- **URL**: `/accounts/:id`
- **Method**: `PUT`
- **Payload**: 与添加账号一致（仅需传递待修改字段）。

---

## 4. 账号有效性校验
手动触发后端模拟登录，校验凭证是否有效并更新昵称、容量及会员信息。

- **URL**: `/accounts/:id/check`
- **Method**: `POST`
- **Response**: 返回更新后的账号对象。

---

## 5. 删除账号
彻底移除该账号信息。

- **URL**: `/accounts/:id`
- **Method**: `DELETE`
