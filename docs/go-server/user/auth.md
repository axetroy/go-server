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
