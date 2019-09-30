### 用户注册

[POST] /v1/auth/signup

| 参数        | 类型     | 说明                                                                      | 必选 |
| ----------- | -------- | ------------------------------------------------------------------------- | ---- |
| username    | `string` | 通过用户名来注册, username, email, phone 三选一                           |      |
| email       | `string` | 通过邮箱来注册, username, email, phone 三选一                             |      |
| phone       | `string` | 通过手机来注册, username, email, phone 三选一, 目前手机注册无法发送验证码 |      |
| password    | `string` | 账号密码                                                                  | \*   |
| invite_code | `string` | 邀请码                                                                    |      |

### 用户登陆

[POST] /v1/auth/signin

| 参数     | 类型     | 说明                                     | 必选 |
| -------- | -------- | ---------------------------------------- | ---- |
| account  | `string` | 用户账号, username/email/phone 中的一个  | \*   |
| password | `string` | 账号密码                                 | \*   |
| code     | `string` | TODO: 手机验证码, 手机可以通过验证码登陆 |      |

### 微信小程序登陆

[POST] /v1/auth/signin/wechat

| 参数 | 类型     | 说明                              | 必选 |
| ---- | -------- | --------------------------------- | ---- |
| code | `string` | 微信小程序授权成功后返回的 `code` | \*   |

在小程序端，用户授权之后，把 `code` 发送给后端，后端根据 `openid` 锁定唯一帐号，并且返回用户身份资料

如果该微信用户还没有注册平台帐号, 则会带 `Response Header` 中附带 `X-Wechat-OpenID`

需要微信端再调用一次接口(`/v1/auth/signin/wechat_complete`)，获取用户的身份资料，例如`手机号码`/`邮箱`, 再次提交

### 微信小程序帐号信息补全

[PUT] /v1/auth/signin/wechat_complete

| 参数     | 类型     | 说明                              | 必选 |
| -------- | -------- | --------------------------------- | ---- |
| code     | `string` | 微信小程序授权成功后返回的 `code` | \*   |
| username | `string` | 帐号名                            | \*   |
| phone    | `string` | 微信用户的手机号                  | \*   |

### 账号激活

[POST] /v1/auth/activation

| 参数 | 类型     | 说明                                        | 必选 |
| ---- | -------- | ------------------------------------------- | ---- |
| code | `string` | 激活码，激活码来自服务器发到的邮箱/手机短信 | \*   |

### 忘记密码

[POST] /v1/auth/password/reset

| 参数         | 类型     | 说明                                        | 必选 |
| ------------ | -------- | ------------------------------------------- | ---- |
| code         | `string` | 重置码，重置码来自服务器发到的邮箱/手机短信 | \*   |
| new_password | `string` | 新的密码                                    |      | \* |
