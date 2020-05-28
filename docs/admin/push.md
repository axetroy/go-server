### 生成一条推送

[POST] /v1/push/notification

| 参数     | 类型       | 说明             | 必填 |
| -------- | ---------- | ---------------- | ---- |
| user_ids | `[]string` | 要推送的指定用户 | \*   |
| title    | `string`   | 标题             | \*   |
| content  | `string`   | 内容             | \*   |

```bash
curl -H "Authorization: Bearer 你的身份令牌" \
     -X POST \
     -d '{"user_ids": ["266972131143712768"], "title": "测试 title", "content": "测试 content"}' \
     http://localhost/v1/push/notification
```

```json
{ "message": "", "data": null, "status": 1 }
```
