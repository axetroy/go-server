### 消息列表

[GET] /v1/message

获取我的消息列表

### 消息详情

[GET] /v1/message/:message_id

获取个人消息的详情

### 标记已读

[PUT] /v1/message/:message_id/read

标记个人消息为已读

### 删除消息

[DELETE] /v1/message/:message_id

删除一条个人消息

### 全部已读

[PUT] /v1/message/read/all

全部已所有个人消息

### 已读多个消息

[PUT] /v1/message/read/batch

| 参数 | 类型       | 说明                                  | 必选 |
| ---- | ---------- | ------------------------------------- | ---- |
| ids  | `[]string` | 要设置已读的 id 数组, 长度最多 100 条 | \*   |

### 获取当前信息的状态

[GET] /v1/message/status

里面包含未读数量
