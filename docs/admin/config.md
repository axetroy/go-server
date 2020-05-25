### 配置名称

当前可用的配置名称字段

| 配置名称   | 描述                                  |
| ---------- | ------------------------------------- |
| smtp       | SMTP 邮件服务的配置，用于发送邮件服务 |
| phone      | 手机相关的配置，用于发送手机短信      |
| wechat_app | 微信小程序相关的配置                  |

### 获取配置名称列表

[GET] /v1/config/name

获取所有的配置的 `名称/描述` 列表

### 创建配置

[POST] /v1/config/:config_name

仅限于超级管理员，且只能创建一次

根据不同的 `config_name` 需要传入不同的配置

#### smtp

| 参数       | 类型     | 说明                 | 必填 |
| ---------- | -------- | -------------------- | ---- |
| server     | `string` | SMTP 服务器地址      | \*   |
| port       | `int`    | SMTP 服务器端口      | \*   |
| username   | `string` | 登录 SMTP 的用户名   | \*   |
| password   | `string` | 登录 SMTP 的密码     | \*   |
| from_name  | `string` | 邮件发送者的名字     | \*   |
| from_email | `string` | 邮件发送者的邮箱地址 | \*   |

#### wechat_app

| 参数   | 类型     | 说明                | 必填 |
| ------ | -------- | ------------------- | ---- |
| app_id | `string` | 微信小程序的 APP ID | \*   |
| secret | `string` | 微信小程序的 secret | \*   |

### 修改配置

[PUT] /v1/config/:config_name

仅限于超级管理员

参数与创建接口一致，并且需要传入完整的字段

### 获取指定的配置

[GET] /v1/config/:config_name

获取指定的配置

### 配置列表

[GET] /v1/config

获取配置列表
