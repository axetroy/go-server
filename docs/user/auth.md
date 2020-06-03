### 用户注册

#### 用户名注册

[POST] /v1/auth/signup

| 参数        | 类型     | 说明     | 必选 |
| ----------- | -------- | -------- | ---- |
| username    | `string` | 用户名   | \*   |
| password    | `string` | 账号密码 | \*   |
| invite_code | `string` | 邀请码   |      |

```bash
curl -X POST \
     -d '{"username": "test1", "password": "test1"}' \
     http://localhost:9001/v1/auth/signup
```

```json
{
  "message": "",
  "data": {
    "id": "274588402135859200",
    "username": "test1",
    "nickname": "test1",
    "email": null,
    "phone": null,
    "status": 1,
    "gender": 0,
    "avatar": "",
    "role": ["user"],
    "level": 0,
    "invite_code": "d9a566c5",
    "username_rename_remaining": 0,
    "pay_password": false,
    "wechat": null,
    "created_at": "2020-06-03T16:21:49.675462+08:00",
    "updated_at": "2020-06-03T16:21:49.675462+08:00"
  },
  "status": 1
}
```

#### 邮箱注册

[POST] /v1/auth/signup/email

| 参数        | 类型     | 说明                                            | 必选 |
| ----------- | -------- | ----------------------------------------------- | ---- |
| email       | `string` | 邮箱地址                                        | \*   |
| password    | `string` | 账号密码                                        | \*   |
| code        | `string` | 邮箱验证码, 通过 `/v1/email/send/register` 发送 | \*   |
| invite_code | `string` | 邀请码                                          |      |

#### 手机注册

[POST] /v1/auth/signup/phone

| 参数        | 类型     | 说明                                                | 必选 |
| ----------- | -------- | --------------------------------------------------- | ---- |
| phone       | `string` | 手机号                                              | \*   |
| code        | `string` | 手机号收到的验证码, 通过 `/v1/auth/code/phone` 发送 | \*   |
| invite_code | `string` | 邀请码                                              |      |

### 用户登陆

[POST] /v1/auth/signin

| 参数     | 类型     | 说明                                    | 必选 |
| -------- | -------- | --------------------------------------- | ---- |
| account  | `string` | 用户账号, username/email/phone 中的一个 | \*   |
| password | `string` | 账号密码                                | \*   |

```bash
curl -X POST \
     -d '{"account": "test1", "password": "test1"}' \
     http://localhost:9001/v1/auth/signin
```

```json
{
  "message": "",
  "data": {
    "id": "266972131143712768",
    "username": "test1",
    "nickname": "this",
    "email": null,
    "phone": null,
    "status": 1,
    "gender": 0,
    "avatar": "http://example.com/v1/resource/image/26ce518102f9907c2ba9b94927bcfa3e.jpg",
    "role": ["user"],
    "level": 0,
    "invite_code": "935cd3fb",
    "username_rename_remaining": 0,
    "pay_password": false,
    "wechat": null,
    "created_at": "2020-05-13T07:57:29.167257Z",
    "updated_at": "2020-05-27T06:03:56.339296Z",
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1aWQiOiJNalkyT1RjeU1UTXhNVFF6TnpFeU56WTQiLCJhdWQiOiIyNjY5NzIxMzExNDM3MTI3NjgiLCJleHAiOjE1OTA2NjA3MzUsImp0aSI6IjI2Njk3MjEzMTE0MzcxMjc2OCIsImlhdCI6MTU5MDYzOTEzNSwiaXNzIjoidXNlciIsIm5iZiI6MTU5MDYzOTEzNX0.cV8Q6gARJEJnVyMlzKUhPN6HqeNYq2e9_cTxO3rDZq8"
  },
  "status": 1
}
```

### 手机号登陆

[POST] /v1/auth/signin/phone

使用`手机号`+`短信验证码`的形式登陆帐号

| 参数  | 类型     | 说明                                                  | 必选 |
| ----- | -------- | ----------------------------------------------------- | ---- |
| phone | `string` | 手机号                                                | \*   |
| code  | `string` | 手机收到的短信验证码，通过 `/v1/auth/code/phone` 发送 | \*   |

### 邮箱登陆

[POST] /v1/auth/signin/email

使用`邮箱地址`+`邮箱证码`的形式登陆帐号

| 参数  | 类型     | 说明                                              | 必选 |
| ----- | -------- | ------------------------------------------------- | ---- |
| email | `string` | 邮箱地址                                          | \*   |
| code  | `string` | 邮箱收到的验证码，通过 `/v1/auth/code/email` 发送 | \*   |

### 微信小程序登陆

[POST] /v1/auth/signin/wechat

| 参数 | 类型     | 说明                              | 必选 |
| ---- | -------- | --------------------------------- | ---- |
| code | `string` | 微信小程序授权成功后返回的 `code` | \*   |

如果该微信用户没有绑定过帐号，则默认创建一个. 密码随机，建议创建帐号后修改密码。

### oAuth 认证登陆

[POST] /v1/auth/signin/oauth

| 参数 | 类型     | 说明                                                     | 必选 |
| ---- | -------- | -------------------------------------------------------- | ---- |
| code | `string` | oAuth 认证接口(`/v1/oauth/:provider`)成功后返回的 `code` | \*   |

如果该没有绑定过帐号，则默认创建一个. 密码随机，建议创建帐号后修改密码。

### 忘记密码

[POST] /v1/auth/password/reset

| 参数         | 类型     | 说明                                        | 必选 |
| ------------ | -------- | ------------------------------------------- | ---- |
| code         | `string` | 重置码，重置码来自服务器发到的邮箱/手机短信 | \*   |
| new_password | `string` | 新的密码                                    |      | \* |

### 发送邮箱验证码

[POST] /v1/auth/code/email

用户验证该邮箱是这个用户所有

| 参数  | 类型     | 说明     | 必选 |
| ----- | -------- | -------- | ---- |
| email | `string` | 邮箱地址 | \*   |

### 发送短信验证码

[POST] /v1/auth/code/phone

用户验证该手机号是这个用户所有

| 参数  | 类型     | 说明   | 必选 |
| ----- | -------- | ------ | ---- |
| phone | `string` | 手机号 | \*   |
