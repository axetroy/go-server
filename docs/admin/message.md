### 新增个人消息

[POST] /v1/message

| 参数    | 类型     | 说明     | 必填 |
| ------- | -------- | -------- | ---- |
| uid     | `string` | 用户 ID  | \*   |
| title   | `string` | 通知标题 | \*   |
| content | `string` | 通知内容 | \*   |

### 删除个人消息

[DELETE] /v1/message/:message_id

### 更改个人消息

[PUT] /v1/message/:message_id

| 参数    | 类型     | 说明     | 必填 |
| ------- | -------- | -------- | ---- |
| title   | `string` | 消息标题 |      |
| content | `string` | 消息内容 |      |

### 消息列表

[GET] /v1/message

### 消息详情

[GET] /v1/message/:message_id