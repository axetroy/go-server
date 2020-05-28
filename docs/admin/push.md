推送通知都会附带数据 `Data` 结构如下

```go
type NotificationBody struct {
	Event   NotificationClickEvent `json:"event"`   // 事件名
	Payload interface{}            `json:"payload"` // 数据体
}
```

例如 `{"event": "none", "payload": "任何可能的数据"}`

其中包含以下事件

| Event                   | 说明                             | Payload                                            |
| ----------------------- | -------------------------------- | -------------------------------------------------- |
| none                    | 无意义的事件，APP 应该什么都不做 | `none`                                             |
| login_abnormal          | 用户登录异常的推送               | `none`                                             |
| new_system_notification | 管理员发送一个新的系统通知的推送 | `{"id": "xxxx", "title": "xxx", "content": "xxx"}` |

### 生成一条推送

[POST] /v1/push/notification

| 参数     | 类型       | 说明                              | 必填 |
| -------- | ---------- | --------------------------------- | ---- |
| user_ids | `[]string` | 要推送的指定用户                  | \*   |
| title    | `string`   | 标题                              | \*   |
| content  | `string`   | 内容                              | \*   |
| payload  | `object`   | 通知附带的数据，用于 App 点开通知 |      |

```bash
curl -H "Authorization: Bearer 你的身份令牌" \
     -X POST \
     -d '{"user_ids": ["266972131143712768"], "title": "测试 title", "content": "测试 content"}' \
     http://localhost/v1/push/notification
```

```json
{ "message": "", "data": null, "status": 1 }
```
