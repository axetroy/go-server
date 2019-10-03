项目配置需要一个 `.env` 文件，通过环境变量的形式进行配置

| 环境变量                      | 类型     | 说明                                                     | 默认值          |
| ----------------------------- | -------- | -------------------------------------------------------- | --------------- |
| 用户接口配置                  | -        | -                                                        | -               |
| USER_HTTP_PORT                | `int`    | 用户接口服务监听的端口                                   | `8080`          |
| USER_HTTP_DOMAIN              | `string` | 用户接口服务的域名                                       | `localhost`     |
| USER_TOKEN_SECRET_KEY         | `string` | 用户接口服务的密钥，用于签发 `token`, 该配置不可泄漏     | `""`            |
| USER_TLS_CERT                 | `string` | TLS 的证书文件                                           | `""`            |
| USER_TLS_KEY                  | `string` | TLS 的 key 文件                                          | `""`            |
| 管理员接口配置                | -        | -                                                        | -               |
| ADMIN_HTTP_PORT               | `int`    | 管理员接口服务监听的端口                                 | `8081`          |
| ADMIN_HTTP_DOMAIN             | `string` | 管理员接口服务的域名                                     | `localhost`     |
| ADMIN_TOKEN_SECRET_KEY        | `string` | 管理员接口服务的密钥，用于签发 `token`, 该配置不可泄漏   | `""`            |
| ADMIN_TLS_CERT                | `string` | TLS 的证书文件                                           | `""`            |
| ADMIN_TLS_KEY                 | `string` | TLS 的 key 文件                                          | `""`            |
| 通用配置                      | -        | -                                                        | -               |
| MACHINE_ID                    | `int`    | 机器 ID, 在集群中，每个 ID 都应该不同，用于产出不同的 ID | `0`             |
| GO_MOD                        | `string` | 处于开发模式(development)/生产模式(production)           | `development`   |
| SIGNATURE_KEY                 | `string` | 数据签名的密钥, 该配置不可泄漏                           | `signature key` |
| UPLOAD_DIR                    | `string` | 图片上传储存的目录                                       | `upload`        |
| UPLOAD_FILE_MAX_SIZE          | `int`    | 文件上传的最大大小                                       | `10485760`      |
| UPLOAD_FILE_EXTENSION         | `string` | 允许上传的文件类型, 以为 `,` 作为分隔符                  | `.txt,.md`      |
| UPLOAD_IMAGE_MAX_SIZE         | `int`    | 图片上传的最大大小                                       | `10485760`      |
| UPLOAD_IMAGE_THUMBNAIL_WIDTH  | `int`    | 图片缩略图宽度, 单位 `px`                                | `100`           |
| UPLOAD_IMAGE_THUMBNAIL_HEIGHT | `int`    | 图片缩略图高度, 单位 `px`                                | `100`           |
| 数据库配置                    | -        | -                                                        | -               |
| DB_HOST                       | `string` | 连接的数据库地址                                         | `localhost`     |
| DB_PORT                       | `int`    | 连接的数据库端口                                         | `65432`         |
| DB_DRIVER                     | `string` | 数据库驱动器, 即数据库类型                               | `postgres`      |
| DB_NAME                       | `string` | 数据库名称                                               | `gotest`        |
| DB_USERNAME                   | `string` | 连接数据库的用户名                                       | `gotest`        |
| DB_PASSWORD                   | `string` | 连接数据库的密码                                         | `gotest`        |
| DB_SYNC                       | `string` | 在应用启动时，是否同步数据库表, 可选 `on`/`off`          | `on`            |
| Redis 配置                    | -        | -                                                        | -               |
| REDIS_SERVER                  | `string` | `redis` 服务器地址                                       | `localhost`     |
| REDIS_PORT                    | `string` | `redis` 服务器端口                                       | `6379`          |
| REDIS_PASSWORD                | `string` | `redis` 服务器密码                                       | `""`            |
| SMTP 服务器配置               | -        | -                                                        | -               |
| SMTP_SERVER                   | `string` | SMTP 服务器                                              | `""`            |
| SMTP_SERVER_PORT              | `int`    | SMTP 服务器的端口                                        | `""`            |
| SMTP_USERNAME                 | `string` | SMTP 服务器的用户名                                      | `""`            |
| SMTP_PASSWORD                 | `string` | SMTP 服务器的密码                                        | `""`            |
| SMTP_FROM_NAME                | `string` | SMTP 服务器发送邮件的发送者                              | `""`            |
| SMTP_FROM_EMAIL               | `string` | SMTP 服务器发送邮件的发送者的邮箱地址                    | `""`            |
| 消息队列配置                  | -        | -                                                        | -               |
| MSG_QUEUE_SERVER              | `string` | 消息队列服务器地址                                       | `localhost`     |
| MSG_QUEUE_PORT                | `int`    | 消息队列服务器端口                                       | `4150`          |
| Google 认证登陆配置           | -        | -                                                        | -               |
| GOOGLE_AUTH2_CLIENT_ID        | `string` | Google 登陆的 client ID                                  | `""`            |
| GOOGLE_AUTH2_CLIENT_SECRET    | `string` | Google 登陆的 secret                                     | `""`            |
| 微信小程序认证登陆配置        | -        | -                                                        | -               |
| WECHAT_APP_ID                 | `string` | 微信小程序的 `appid`                                     | `""`            |
| WECHAT_SECRET                 | `string` | 微信小程序的 `secret`                                    | `""`            |

例如以下配置

```env
##################### 用户端专有配置 #####################
USER_HTTP_PORT=9090 # 用户端的 HTTP 监听端口. 默认 8080
USER_HTTP_DOMAIN=http://127.0.0.1:8080 # 用户端的 API 域名
USER_TOKEN_SECRET_KEY=user # 用户端的 JWT token 密钥

##################### 管理员专有配置 #####################
ADMIN_HTTP_PORT=9091 # 管理员端的 HTTP 监听端口. 默认 8081
ADMIN_HTTP_DOMAIN=http://127.0.0.1:8081 # 用户端的 API 域名
ADMIN_TOKEN_SECRET_KEY=admin # 管理员端的 JWT token 密钥


######################## 公共配置 ########################
# 通用
MACHINE_ID="0" # 机器 ID, 在集群中，每个ID都应该不同，用于产出不同的 ID
GO_MOD="production" # 处于开发模式(development)/生产模式(production), 默认 development
SIGNATURE_KEY="signature key" # 数据签名的密钥, 该配置不可泄漏
UPLOAD_DIR=upload # 图片上传储存的目录
UPLOAD_FILE_MAX_SIZE=10485760 # 文件上传的最大大小，这里是 1024 * 1024 * 10 = 10M
UPLOAD_FILE_EXTENSION=".txt,.md" # 允许上传的文件类型
UPLOAD_IMAGE_MAX_SIZE=10485760 # 图片上传的最大大小，这里是 1024 * 1024 * 10 = 10M
UPLOAD_IMAGE_THUMBNAIL_WIDTH=100 # 图片缩略图宽度
UPLOAD_IMAGE_THUMBNAIL_HEIGHT=100 # 图片的缩略图高度

# 主数据库设置
DB_HOST="${DB_HOST}" # 默认 localhost
DB_PORT="${DB_PORT}" # 默认 "65432", postgres 官方端口 54321
DB_DRIVER="${DB_DRIVER}" # 默认 "postgres"
DB_NAME="${DB_NAME}" # 默认 "gotest"
DB_USERNAME="${DB_USERNAME}" # 默认 "gotest"
DB_PASSWORD="${DB_PASSWORD}" # 默认 "gotest"
DB_SYNC="on" # 在应用启动时，是否同步数据库表, 可选 on/off, 默认 on

# Redis 缓存服务器配置
REDIS_SERVER=localhost #  Redis 服务器地址
REDIS_PORT=6379 # Redis 端口
REDIS_PASSWORD=password # 连接服务器密码

# SMTP 服务器配置，用于发送邮件
SMTP_SERVER = smtp.qq.com # 邮件服务器
SMTP_SERVER_PORT = 465 # 邮件服务器端口
SMTP_USERNAME = 450409405 # 邮件服务器用户名
SMTP_PASSWORD = "${SMTP_PASSWORD}" # 邮件服务器密码
SMTP_FROM_NAME = Axetroy # 邮件发送者名
SMTP_FROM_EMAIL = 450409405@qq.com # 邮件发送地址

# 消息队列配置
MSG_QUEUE_SERVER = 127.0.0.1 # 消息队列服务器地址. 默认 127.0.0.1
MSG_QUEUE_PORT = 4150 # 消息队列服务器端口. 默认 4150

# OAuth2 认证服务
GOOGLE_AUTH2_CLIENT_ID = "${GOOGLE_AUTH2_CLIENT_ID}" # Google oAuth2 的 client ID
GOOGLE_AUTH2_CLIENT_SECRET = "${GOOGLE_AUTH2_CLIENT_SECRET}" # Google oAuth2 的 client secret

# 微信小程序认证
WECHAT_APP_ID = "${WECHAT_APP_ID}"
WECHAT_SECRET = "${WECHAT_SECRET}"
```
