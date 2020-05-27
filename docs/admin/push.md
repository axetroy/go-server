### 生成一条推送

[POST] /v1/push/notification

| 参数     | 类型       | 说明             | 必填 |
| -------- | ---------- | ---------------- | ---- |
| user_ids | `[]string` | 要推送的指定用户 | \*   |
| title    | `string`   | 标题             | \*   |
| content  | `string`   | 内容             | \*   |
