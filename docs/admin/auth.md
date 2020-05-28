### 管理员登陆

[POST] /v1/login

| 参数     | 类型     | 说明       | 必填 |
| -------- | -------- | ---------- | ---- |
| username | `string` | 管理员账号 | \*   |
| password | `string` | 账号密码   | \*   |

```bash
curl -X POST \
     -d '{"username": "admin", "password": "admin"}' \
     http://localhost/v1/login
```

```json
{
  "message": "",
  "data": {
    "id": "266237936893165568",
    "username": "admin",
    "name": "admin",
    "accession": [],
    "is_super": true,
    "status": 0,
    "created_at": "2020-05-11T07:20:03.617075Z",
    "updated_at": "2020-05-11T07:20:03.617075Z",
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1aWQiOiJNalkyTWpNM09UTTJPRGt6TVRZMU5UWTQiLCJhdWQiOiIyNjYyMzc5MzY4OTMxNjU1NjgiLCJleHAiOjE1OTA2NTcwMTgsImp0aSI6IjI2NjIzNzkzNjg5MzE2NTU2OCIsImlhdCI6MTU5MDYzNTQxOCwiaXNzIjoiYWRtaW4iLCJuYmYiOjE1OTA2MzU0MTh9.azzP0TDfFuY4ybDDj-6sS0Cfj9uN3MxV4lNf7gCNVzc"
  },
  "status": 1
}
```
