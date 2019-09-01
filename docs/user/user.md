### 获取用户信息

[GET] /v1/user/profile

获取用户的详细信息资料

### 更新用户信息

[PUT] /v1/user/profile

| 参数     | 类型     | 说明         | 必选 |
| -------- | -------- | ------------ | ---- |
| nickname | `string` | 用户昵称     |      |
| gender   | `string` | 用户性别     |      |
| avatar   | `string` | 用户头像 URL |      |

### 修改登陆密码

[PUT] /v1/user/password

| 参数          | 类型     | 说明   | 必选 |
| ------------- | -------- | ------ | ---- |
| old_passworld | `string` | 旧密码 | \*   |
| new_password  | `string` | 新密码 | \*   |

### 设置二级密码

[POST] /v1/user/password2

| 参数             | 类型     | 说明         | 必选 |
| ---------------- | -------- | ------------ | ---- |
| password         | `string` | 二级密码     | \*   |
| password_confirm | `string` | 二级密码确认 | \*   |

### 修改二级密码

[PUT] /v1/user/password2

| 参数         | 类型     | 说明       | 必选 |
| ------------ | -------- | ---------- | ---- |
| old_password | `string` | 旧二级密码 | \*   |
| new_password | `string` | 新二级密码 | \*   |

### 发送重置二级密码的邮件/短信

[POST] /v1/user/password2/reset

如果用户有手机，则发送手机验证码，如果有邮箱，则发送邮件

### 重置二级密码

[PUT] /v1/user/password2/reset

| 参数         | 类型     | 说明             | 必选 |
| ------------ | -------- | ---------------- | ---- |
| code         | `string` | 二级密码的重置码 | \*   |
| new_password | `string` | 新二级密码       | \*   |

### 邀请列表

[GET] /v1/user/invite

获取我的邀请列表

### 邀请详情

[GET] /v1/user/invite/i/:invite_id

| 参数      | 类型     | 说明          | 必选 |
| --------- | -------- | ------------- | ---- |
| invite_id | `string` | 邀请数据的 ID | \*   |

### 上传头像

[POST] /v1/user/avatar

头像上传为 Form 表单

| 参数 | 类型   | 说明                                  | 必选 |
| ---- | ------ | ------------------------------------- | ---- |
| file | `file` | 要上传的头像图片，仅支持 jpg/jpeg/png | \*   |
