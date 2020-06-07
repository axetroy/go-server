### 系统通知列表

[GET] /v1/notification

获取系统通知列表

### 系统通知详情

[GET] /v1/notification/:notification_id

获取某个系统通知详情

### 标记系统通知已读

[PUT] /v1/notification/:notification_id/read

标记系统通知为已读

### 批量标记已读

[PUT] /v1/notification/read/batch

批量标记已读

| 参数 | 类型       | 说明             | 必选 |
| ---- | ---------- | ---------------- | ---- |
| ids  | `[]string` | 要已读的 ID 数组 | \*   |
