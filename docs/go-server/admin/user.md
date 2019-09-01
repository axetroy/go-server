### 会员列表

[GET] /v1/user

获取所有的会员列表

### 会员详情

[GET] /v1/user/u/:user_id

获取指定会员详情

### 修改会员资料

[PUT] /v1/user/u/:user_id

| 参数     | 类型     | 说明         | 必填 |
| -------- | -------- | ------------ | ---- |
| nickname | `string` | 用户昵称     |      |
| gender   | `string` | 用户性别     |      |
| avatar   | `string` | 用户头像 URL |      |

### 修改会员登陆密码

[PUT] /v1/user/u/:user_id/password

| 参数         | 类型     | 说明   | 必填 |
| ------------ | -------- | ------ | ---- |
| new_password | `string` | 新密码 | \*   |
