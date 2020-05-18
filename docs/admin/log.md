### 获取用户的登陆日志列表

[GET] /v1/log/login

获取用户的登陆日志列表, 筛选条件如下

| 参数    | 类型     | 说明         | 必选 |
| ------- | -------- | ------------ | ---- |
| uid     | `string` | 用户 ID      |      |
| type    | `int`    | 登陆类型     |      |
| command | `int`    | 当前状态     |      |
| ip      | `string` | 根据 IP 筛选 |      |

### 获取用户的登陆日志详情

[GET] /v1/log/login/:log_id