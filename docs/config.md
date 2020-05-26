项目配置需要一个 `.env` 文件，通过环境变量的形式进行配置

`.env` 文件不是必须的，也可以直接通过环境变量进行设置

### 用户端配置

> 提供用户端的接口服务

| 环境变量                                       | 类型     | 说明                                                         | 默认值       |
| ---------------------------------------------- | -------- | ------------------------------------------------------------ | ------------ |
| 通用配置                                       | -        | -                                                            | -            |
| MACHINE_ID                                     | `int`    | 机器 ID, 在集群中，每个机器 ID 都应该不同，用于产出不同的 ID | `0`          |
| GO_MOD                                         | `string` | 处于开发模式(development)/生产模式(production)               | `production` |
| USER_HTTP_PORT                                 | `int`    | 用户接口服务监听的端口                                       | `9001`       |
| USER_HTTP_DOMAIN                               | `string` | 用户接口服务的域名                                           | `localhost`  |
| USER_TOKEN_SECRET_KEY                          | `string` | 用户接口服务的密钥，用于签发 `token`, 该配置不可泄           | `""`         |
| 数据库配置                                     | -        | -                                                            | -            |
| DB_HOST                                        | `string` | 连接的数据库地址                                             | `localhost`  |
| DB_PORT                                        | `int`    | 连接的数据库端口                                             | `65432`      |
| DB_DRIVER                                      | `string` | 数据库驱动器, 即数据库类型                                   | `postgres`   |
| DB_NAME                                        | `string` | 数据库名称                                                   | `gotest`     |
| DB_USERNAME                                    | `string` | 连接数据库的用户名                                           | `gotest`     |
| DB_PASSWORD                                    | `string` | 连接数据库的密码                                             | `gotest`     |
| Redis 配置                                     | -        | -                                                            | -            |
| REDIS_SERVER                                   | `string` | `redis` 服务器地址                                           | `localhost`  |
| REDIS_PORT                                     | `string` | `redis` 服务器端口                                           | `6379`       |
| REDIS_PASSWORD                                 | `string` | `redis` 服务器密码                                           | `""`         |
| 短信服务设置                                   | -        | -                                                            | -            |
| TELEPHONE_PROVIDER                             | `string` | 短信服务提供商，可选 `aliyun`/`tencent`                      | `aliyun`     |
| TELEPHONE_ALIYUN_ACCESS_KEY                    | `string` | *阿里云*的 access key                                        | `""`         |
| TELEPHONE_ALIYUN_ACCESS_SECRET                 | `string` | *阿里云*的 access secret                                     | `""`         |
| TELEPHONE_ALIYUN_SIGN_NAME                     | `string` | *阿里云*短信的签名名称                                       | `""`         |
| TELEPHONE_ALIYUN_TEMPLATE_CODE_AUTH            | `string` | *阿里云*用于发送身份验证的短信模版代码                       | `""`         |
| TELEPHONE_ALIYUN_TEMPLATE_CODE_RESET_PASSWORD  | `string` | *阿里云*用于发送重置密码的短信模版代码                       | `""`         |
| TELEPHONE_ALIYUN_TEMPLATE_CODE_REGISTER        | `string` | *阿里云*用于发送注册帐号的短信模版代码                       | `""`         |
| TELEPHONE_TENCENT_APP_ID                       | `string` | *腾讯云*的 AppId                                             | `""`         |
| TELEPHONE_TENCENT_APP_KEY                      | `string` | *腾讯云*的 AppKey                                            | `""`         |
| TELEPHONE_TENCENT_SIGN                         | `string` | *腾讯云*的 短信签名内容                                      | `""`         |
| TELEPHONE_TENCENT_TEMPLATE_CODE_AUTH           | `string` | *腾讯云*用于发送身份验证的短信模版代码                       | `""`         |
| TELEPHONE_TENCENT_TEMPLATE_CODE_RESET_PASSWORD | `string` | *腾讯云*用于发送重置密码的短信模版代码                       | `""`         |
| TELEPHONE_TENCENT_TEMPLATE_CODE_REGISTER       | `string` | *腾讯云*用于发送注册帐号的短信模版代码                       | `""`         |
| 消息队列配置                                   | -        | -                                                            | -            |
| MSG_QUEUE_SERVER                               | `string` | 消息队列服务器地址                                           | `localhost`  |
| MSG_QUEUE_PORT                                 | `int`    | 消息队列服务器端口                                           | `4150`       |
| Google 认证登陆配置                            | -        | -                                                            | -            |
| GOOGLE_AUTH2_CLIENT_ID                         | `string` | Google 登陆的 client ID                                      | `""`         |
| GOOGLE_AUTH2_CLIENT_SECRET                     | `string` | Google 登陆的 secret                                         | `""`         |
| oAuth 认证设置                                 | -        | -                                                            | -            |
| OAUTH_REDIRECT_URL                             | `string` | oAuth 认证成功后跳转到的前端 URL                             | `""`         |
| GITHUB_KEY                                     | `string` | oAuth 认证的 `Github Key`                                    | `""`         |
| GITHUB_SECRET                                  | `string` | oAuth 认证的 `Github Secret`                                 | `""`         |
| GITLAB_KEY                                     | `string` | oAuth 认证的 `Gitlab Key`                                    | `""`         |
| GITLAB_SECRET                                  | `string` | oAuth 认证的 `Gitlab Secret`                                 | `""`         |
| GOOGLE_KEY                                     | `string` | oAuth 认证的 `Google Key`                                    | `""`         |
| GOOGLE_SECRET                                  | `string` | oAuth 认证的 `Google Secret`                                 | `""`         |
| FACEBOOK_KEY                                   | `string` | oAuth 认证的 `Facebook Key`                                  | `""`         |
| TWITTER_KEY                                    | `string` | oAuth 认证的 `Twitter Key`                                   | `""`         |
| TWITTER_SECRET                                 | `string` | oAuth 认证的 `Twitter Secret`                                | `""`         |
| 消息队列配置                                   | -        | -                                                            | -            |
| MSG_QUEUE_SERVER                               | `string` | 消息队列服务器地址                                           | `localhost`  |
| MSG_QUEUE_PORT                                 | `int`    | 消息队列服务器端口                                           | `4150`       |

### 管理员端配置

> 提供管理员端的接口服务

| 环境变量               | 类型     | 说明                                                         | 默认值       |
| ---------------------- | -------- | ------------------------------------------------------------ | ------------ |
| 通用配置               | -        | -                                                            | -            |
| MACHINE_ID             | `int`    | 机器 ID, 在集群中，每个机器 ID 都应该不同，用于产出不同的 ID | `0`          |
| GO_MOD                 | `string` | 处于开发模式(development)/生产模式(production)               | `production` |
| ADMIN_HTTP_PORT        | `int`    | 管理员接口服务监听的端口                                     | `9002`       |
| ADMIN_HTTP_DOMAIN      | `string` | 管理员接口服务的域名                                         | `localhost`  |
| ADMIN_TOKEN_SECRET_KEY | `string` | 管理员接口服务的密钥，用于签发 `token`, 该配置不可泄         | `""`         |
| ADMIN_DEFAULT_PASSWORD | `string` | 第一次启动时，默认的管理员密码                               | `"admin"`    |
| 数据库配置             | -        | -                                                            | -            |
| DB_HOST                | `string` | 连接的数据库地址                                             | `localhost`  |
| DB_PORT                | `int`    | 连接的数据库端口                                             | `65432`      |
| DB_DRIVER              | `string` | 数据库驱动器, 即数据库类型                                   | `postgres`   |
| DB_NAME                | `string` | 数据库名称                                                   | `gotest`     |
| DB_USERNAME            | `string` | 连接数据库的用户名                                           | `gotest`     |
| DB_PASSWORD            | `string` | 连接数据库的密码                                             | `gotest`     |
| Redis 配置             | -        | -                                                            | -            |
| REDIS_SERVER           | `string` | `redis` 服务器地址                                           | `localhost`  |
| REDIS_PORT             | `string` | `redis` 服务器端口                                           | `6379`       |
| REDIS_PASSWORD         | `string` | `redis` 服务器密码                                           | `""`         |
| 消息队列配置           | -        | -                                                            | -            |
| MSG_QUEUE_SERVER       | `string` | 消息队列服务器地址                                           | `localhost`  |
| MSG_QUEUE_PORT         | `int`    | 消息队列服务器端口                                           | `4150`       |

### 资源服务器

> 提供静态资源接口服务

| 环境变量                      | 类型     | 说明                                           | 默认值               |
| ----------------------------- | -------- | ---------------------------------------------- | -------------------- |
| 通用配置                      | -        | -                                              | -                    |
| GO_MOD                        | `string` | 处于开发模式(development)/生产模式(production) | `production`         |
| RESOURCE_HTTP_PORT            | `int`    | 资源接口服务监听的端口                         | `9003`               |
| RESOURCE_HTTP_DOMAIN          | `string` | 资源接口服务的域名                             | `localhost`          |
| UPLOAD_DIR                    | `string` | 图片上传储存的目录                             | `upload`             |
| UPLOAD_FILE_MAX_SIZE          | `int`    | 文件上传的最大大小                             | `1024*1024*10` = 10M |
| UPLOAD_FILE_EXTENSION         | `string` | 允许上传的文件类型, 以为 `,` 作为分隔符        | `.txt,.md`           |
| UPLOAD_IMAGE_MAX_SIZE         | `int`    | 图片上传的最大大小                             | `1024*1024*10` = 10M |
| UPLOAD_IMAGE_THUMBNAIL_WIDTH  | `int`    | 图片缩略图宽度, 单位 `px`                      | `100`                |
| UPLOAD_IMAGE_THUMBNAIL_HEIGHT | `int`    | 图片缩略图高度, 单位 `px`                      | `100`                |

### 消息队列服务器

> 消费队列里面的消息

| 环境变量                | 类型     | 说明                                           | 默认值       |
| ----------------------- | -------- | ---------------------------------------------- | ------------ |
| 通用配置                | -        | -                                              | -            |
| GO_MOD                  | `string` | 处于开发模式(development)/生产模式(production) | `production` |
| MSG_QUEUE_SERVER        | `string` | 消息队列服务器地址                             | `localhost`  |
| MSG_QUEUE_PORT          | `int`    | 消息队列服务器端口                             | `4150`       |
| 数据库配置              | -        | -                                              | -            |
| DB_HOST                 | `string` | 连接的数据库地址                               | `localhost`  |
| DB_PORT                 | `int`    | 连接的数据库端口                               | `65432`      |
| DB_DRIVER               | `string` | 数据库驱动器, 即数据库类型                     | `postgres`   |
| DB_NAME                 | `string` | 数据库名称                                     | `gotest`     |
| DB_USERNAME             | `string` | 连接数据库的用户名                             | `gotest`     |
| DB_PASSWORD             | `string` | 连接数据库的密码                               | `gotest`     |
| 推送服务器              | -        | -                                              | -            |
| ONE_SIGNAL_APP_ID       | `string` | 推送服务器 one signal 的 APP ID                | ``           |
| ONE_SIGNAL_REST_API_KEY | `string` | 推送服务器 one signal 的 REST API KEY          | ``           |
