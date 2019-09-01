### 新增系统通知

[POST] /v1/notification

| 参数    | 类型     | 说明     | 必填 |
| ------- | -------- | -------- | ---- |
| title   | `string` | 通知标题 | \*   |
| content | `string` | 通知内容 | \*   |
| note    | `string` | 备注     |      |

### 修改系统通知

[PUT] /v1/notification/n/:notification_id

| 参数    | 类型     | 说明     | 必填 |
| ------- | -------- | -------- | ---- |
| title   | `string` | 通知标题 |      |
| content | `string` | 通知内容 |      |
| note    | `string` | 备注     |      |

### 删除系统通知

[DELETE] /v1/notification/n/:notification_id

### 获取系统通知列表

[GET] /v1/notification

### 获取系统通知详情

[GET] /v1/notification/n/:notification_id