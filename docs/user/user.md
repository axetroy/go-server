### 获取用户信息

[GET] /v1/user/profile

获取用户的详细信息资料

```bash
curl -H "Authorization: Bearer 你的身份令牌" \
     http://localhost:9001/v1/user/profile
```

```json
{
  "message": "",
  "data": {
    "id": "274588402135859200",
    "username": "test1",
    "nickname": "nickname",
    "email": null,
    "phone": null,
    "status": 1,
    "gender": 0,
    "avatar": "",
    "role": [],
    "level": 0,
    "invite_code": "d9a566c5",
    "username_rename_remaining": 0,
    "pay_password": false,
    "wechat": null,
    "created_at": "2020-06-03T08:21:49.675462Z",
    "updated_at": "2020-06-06T18:34:17.420393Z"
  },
  "status": 1
}
```

### 更新用户信息

[PUT] /v1/user/profile

| 参数              | 类型     | 说明                                                                                                                                                       | 必选 |
| ----------------- | -------- | ---------------------------------------------------------------------------------------------------------------------------------------------------------- | ---- |
| username          | `string` | 用户登陆名, 某些用户可以重命名自己的登陆名。例如微信注册的帐号                                                                                             |      |
| nickname          | `string` | 用户昵称                                                                                                                                                   |      |
| gender            | `string` | 用户性别                                                                                                                                                   |      |
| avatar            | `string` | 用户头像 URL                                                                                                                                               |      |
| wechat            | `object` | 更新微信绑定的相关信息<br/> 绑定微信后，没有拿到微信的昵称/性别等信息。所以需要客户端手动调用更新信息<br/>信息由微信小程序接口接口 `wx.getUserInfo()` 获得 |      |
| wechat.nickname   | `string` | 微信昵称                                                                                                                                                   |      |
| wechat.avatar_url | `string` | 微信头像 URL                                                                                                                                               |      |
| wechat.gender     | `int`    | 性别                                                                                                                                                       |      |
| wechat.country    | `string` | 国家                                                                                                                                                       |      |
| wechat.province   | `string` | 省份                                                                                                                                                       |      |
| wechat.city       | `string` | 城市                                                                                                                                                       |      |
| wechat.language   | `string` | 语言                                                                                                                                                       |      |

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

[GET] /v1/user/invite/:invite_id

| 参数      | 类型     | 说明          | 必选 |
| --------- | -------- | ------------- | ---- |
| invite_id | `string` | 邀请数据的 ID | \*   |

### 发送邮箱验证码

[POST] /v1/user/auth/email

发送邮箱验证码至用户绑定的邮箱

### 发送手机验证码

[POST] /v1/user/auth/phone

发送短信验证码至用户绑定的手机号

### 绑定邮箱

[POST] /v1/user/bind/email

绑定用户的邮箱，在用户没有绑定有邮箱时可以调用

| 参数  | 类型     | 说明                                            | 必选 |
| ----- | -------- | ----------------------------------------------- | ---- |
| email | `string` | 要绑定的邮箱                                    | \*   |
| code  | `string` | 邮箱收到的验证，调用 `/v1/auth/code/email` 发送 | \*   |

### 解绑邮箱

[DELETE] /v1/user/unbind/email

解除绑定邮箱

| 参数 | 类型     | 说明                                            | 必选 |
| ---- | -------- | ----------------------------------------------- | ---- |
| code | `string` | 邮箱收到的验证，调用 `/v1/auth/code/email` 发送 | \*   |

### 绑定手机

[POST] /v1/user/bind/phone

绑定用户的手机，在用户没有绑定有手机时可以调用

| 参数  | 类型     | 说明                                            | 必选 |
| ----- | -------- | ----------------------------------------------- | ---- |
| phone | `string` | 要绑定的手机                                    | \*   |
| code  | `string` | 手机收到的验证，调用 `/v1/user/auth/phone` 发送 | \*   |

### 解绑手机

[DELETE] /v1/user/unbind/phone

| 参数 | 类型     | 说明                                            | 必选 |
| ---- | -------- | ----------------------------------------------- | ---- |
| code | `string` | 手机收到的验证，调用 `/v1/user/auth/phone` 发送 | \*   |

### 绑定微信

[POST] /v1/user/bind/phone

绑定微信小程序帐号，绑定后可直接用小程序登陆

| 参数 | 类型     | 说明                                          | 必选 |
| ---- | -------- | --------------------------------------------- | ---- |
| code | `string` | 微信小程序调用 `wx.login()` 之后，返回的 code | \*   |

### 解绑微信

[DELETE] /v1/user/unbind/wechat

| 参数 | 类型     | 说明                                                                                                                                                                        | 必选 |
| ---- | -------- | --------------------------------------------------------------------------------------------------------------------------------------------------------------------------- | ---- |
| code | `string` | 验证码，如果帐号已绑定手机，则为手机号收到的验证码（`/v1/user/auth/phone`），如果有为邮箱，则用邮箱收到的验证码（`/v1/user/auth/email`），否则使用 `wx.login()` 返回的 code | \*   |

### 获取扫码登录的二维码信息

[GET] /v1/user/qrcode/{:link}

参数 `link` 为 Web 端二维码的信息, 格式为 `auth://eyJzZXNzaW9uX2lkIjoiMDZiNjY2NzAtYWFjYS00ZmRkLTg1NDctMTM2YTY1N2ExNTYxIiwiZXhwaXJlZF9hdCI6IjIwMjAtMDYtMTlUMDY6MzQ6MTcuOTQ2WiJ`

返回扫码的基本信息，在哪里登录，IP 多什么，什么设备 等信息

```json
{
  "message": "",
  "data": {
    "os": "Linux",
    "browser": "Chrome",
    "version": "70.0.133123",
    "ip": "192.168.0.1"
  },
  "status": 1
}
```

### 准许扫码登录

[POST] /v1/user/qrcode/grant

准许二维码登录，这个 web 端在循环调用 `/v1/auth/qrcode/check` 时将会成功，并且返回当前帐号的信息

| 参数 | 类型     | 说明         | 必选 |
| ---- | -------- | ------------ | ---- |
| url  | `string` | 二维码的 URL | \*   |
