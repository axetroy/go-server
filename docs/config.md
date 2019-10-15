项目配置需要一个 `.env` 文件，通过环境变量的形式进行配置

| 环境变量                                       | 类型     | 说明                                                                            | 默认值          |
| ---------------------------------------------- | -------- | ------------------------------------------------------------------------------- | --------------- |
| 用户接口配置                                   | -        | -                                                                               | -               |
| USER_HTTP_PORT                                 | `int`    | 用户接口服务监听的端口                                                          | `8080`          |
| USER_HTTP_DOMAIN                               | `string` | 用户接口服务的域名                                                              | `localhost`     |
| USER_TOKEN_SECRET_KEY                          | `string` | 用户接口服务的密钥，用于签发 `token`, 该配置不可泄漏                            | `""`            |
| USER_TLS_CERT                                  | `string` | TLS 的证书文件                                                                  | `""`            |
| USER_TLS_KEY                                   | `string` | TLS 的 key 文件                                                                 | `""`            |
| 管理员接口配置                                 | -        | -                                                                               | -               |
| ADMIN_HTTP_PORT                                | `int`    | 管理员接口服务监听的端口                                                        | `8081`          |
| ADMIN_HTTP_DOMAIN                              | `string` | 管理员接口服务的域名                                                            | `localhost`     |
| ADMIN_TOKEN_SECRET_KEY                         | `string` | 管理员接口服务的密钥，用于签发 `token`, 该配置不可泄漏                          | `""`            |
| ADMIN_TLS_CERT                                 | `string` | TLS 的证书文件                                                                  | `""`            |
| ADMIN_TLS_KEY                                  | `string` | TLS 的 key 文件                                                                 | `""`            |
| ADMIN_DEFAULT_PASSWORD                         | `string` | 默认的超级管理员 admin 的密码，在第一次启动时，会向数据库添加一个超级管理员帐号 | `admin`         |
| 通用配置                                       | -        | -                                                                               | -               |
| MACHINE_ID                                     | `int`    | 机器 ID, 在集群中，每个 ID 都应该不同，用于产出不同的 ID                        | `0`             |
| GO_MOD                                         | `string` | 处于开发模式(development)/生产模式(production)                                  | `development`   |
| SIGNATURE_KEY                                  | `string` | 数据签名的密钥, 该配置不可泄漏                                                  | `signature key` |
| UPLOAD_DIR                                     | `string` | 图片上传储存的目录                                                              | `upload`        |
| UPLOAD_FILE_MAX_SIZE                           | `int`    | 文件上传的最大大小                                                              | `10485760`      |
| UPLOAD_FILE_EXTENSION                          | `string` | 允许上传的文件类型, 以为 `,` 作为分隔符                                         | `.txt,.md`      |
| UPLOAD_IMAGE_MAX_SIZE                          | `int`    | 图片上传的最大大小                                                              | `10485760`      |
| UPLOAD_IMAGE_THUMBNAIL_WIDTH                   | `int`    | 图片缩略图宽度, 单位 `px`                                                       | `100`           |
| UPLOAD_IMAGE_THUMBNAIL_HEIGHT                  | `int`    | 图片缩略图高度, 单位 `px`                                                       | `100`           |
| 数据库配置                                     | -        | -                                                                               | -               |
| DB_HOST                                        | `string` | 连接的数据库地址                                                                | `localhost`     |
| DB_PORT                                        | `int`    | 连接的数据库端口                                                                | `65432`         |
| DB_DRIVER                                      | `string` | 数据库驱动器, 即数据库类型                                                      | `postgres`      |
| DB_NAME                                        | `string` | 数据库名称                                                                      | `gotest`        |
| DB_USERNAME                                    | `string` | 连接数据库的用户名                                                              | `gotest`        |
| DB_PASSWORD                                    | `string` | 连接数据库的密码                                                                | `gotest`        |
| DB_SYNC                                        | `string` | 在应用启动时，是否同步数据库表, 可选 `on`/`off`                                 | `on`            |
| Redis 配置                                     | -        | -                                                                               | -               |
| REDIS_SERVER                                   | `string` | `redis` 服务器地址                                                              | `localhost`     |
| REDIS_PORT                                     | `string` | `redis` 服务器端口                                                              | `6379`          |
| REDIS_PASSWORD                                 | `string` | `redis` 服务器密码                                                              | `""`            |
| SMTP 服务器配置                                | -        | -                                                                               | -               |
| SMTP_SERVER                                    | `string` | SMTP 服务器                                                                     | `""`            |
| SMTP_SERVER_PORT                               | `int`    | SMTP 服务器的端口                                                               | `""`            |
| SMTP_USERNAME                                  | `string` | SMTP 服务器的用户名                                                             | `""`            |
| SMTP_PASSWORD                                  | `string` | SMTP 服务器的密码                                                               | `""`            |
| SMTP_FROM_NAME                                 | `string` | SMTP 服务器发送邮件的发送者                                                     | `""`            |
| SMTP_FROM_EMAIL                                | `string` | SMTP 服务器发送邮件的发送者的邮箱地址                                           | `""`            |
| 短信服务设置                                   | -        | -                                                                               | -               |
| TELEPHONE_PROVIDER                             | `string` | 短信服务提供商，可选 `aliyun`/`tencent`                                         | `aliyun`        |
| TELEPHONE_ALIYUN_ACCESS_KEY                    | `string` | *阿里云*的 access key                                                           | `""`            |
| TELEPHONE_ALIYUN_ACCESS_SECRET                 | `string` | *阿里云*的 access secret                                                        | `""`            |
| TELEPHONE_ALIYUN_SIGN_NAME                     | `string` | *阿里云*短信的签名名称                                                          | `""`            |
| TELEPHONE_ALIYUN_TEMPLATE_CODE_AUTH            | `string` | *阿里云*用于发送身份验证的短信模版代码                                          | `""`            |
| TELEPHONE_ALIYUN_TEMPLATE_CODE_RESET_PASSWORD  | `string` | *阿里云*用于发送重置密码的短信模版代码                                          | `""`            |
| TELEPHONE_ALIYUN_TEMPLATE_CODE_REGISTER        | `string` | *阿里云*用于发送注册帐号的短信模版代码                                          | `""`            |
| TELEPHONE_TENCENT_APP_ID                       | `string` | *腾讯云*的 AppId                                                                | `""`            |
| TELEPHONE_TENCENT_APP_KEY                      | `string` | *腾讯云*的 AppKey                                                               | `""`            |
| TELEPHONE_TENCENT_SIGN                         | `string` | *腾讯云*的 短信签名内容                                                         | `""`            |
| TELEPHONE_TENCENT_TEMPLATE_CODE_AUTH           | `string` | *腾讯云*用于发送身份验证的短信模版代码                                          | `""`            |
| TELEPHONE_TENCENT_TEMPLATE_CODE_RESET_PASSWORD | `string` | *腾讯云*用于发送重置密码的短信模版代码                                          | `""`            |
| TELEPHONE_TENCENT_TEMPLATE_CODE_REGISTER       | `string` | *腾讯云*用于发送注册帐号的短信模版代码                                          | `""`            |
| 消息队列配置                                   | -        | -                                                                               | -               |
| MSG_QUEUE_SERVER                               | `string` | 消息队列服务器地址                                                              | `localhost`     |
| MSG_QUEUE_PORT                                 | `int`    | 消息队列服务器端口                                                              | `4150`          |
| Google 认证登陆配置                            | -        | -                                                                               | -               |
| GOOGLE_AUTH2_CLIENT_ID                         | `string` | Google 登陆的 client ID                                                         | `""`            |
| GOOGLE_AUTH2_CLIENT_SECRET                     | `string` | Google 登陆的 secret                                                            | `""`            |
| 微信小程序认证登陆配置                         | -        | -                                                                               | -               |
| WECHAT_APP_ID                                  | `string` | 微信小程序的 `appid`                                                            | `""`            |
| WECHAT_SECRET                                  | `string` | 微信小程序的 `secret`                                                           | `""`            |
| oAuth 认证设置                                 | -        | -                                                                               | -               |
| OAUTH_REDIRECT_URL                             | `string` | oAuth 认证成功后跳转到的前端 URL                                                | `""`            |
| GITHUB_KEY                                     | `string` | oAuth 认证的 `Github Key`                                                       | `""`            |
| GITHUB_SECRET                                  | `string` | oAuth 认证的 `Github Secret`                                                    | `""`            |
| GITLAB_KEY                                     | `string` | oAuth 认证的 `Gitlab Key`                                                       | `""`            |
| GITLAB_SECRET                                  | `string` | oAuth 认证的 `Gitlab Secret`                                                    | `""`            |
| GOOGLE_KEY                                     | `string` | oAuth 认证的 `Google Key`                                                       | `""`            |
| GOOGLE_SECRET                                  | `string` | oAuth 认证的 `Google Secret`                                                    | `""`            |
| FACEBOOK_KEY                                   | `string` | oAuth 认证的 `Facebook Key`                                                     | `""`            |
| TWITTER_KEY                                    | `string` | oAuth 认证的 `Twitter Key`                                                      | `""`            |
| TWITTER_SECRET                                 | `string` | oAuth 认证的 `Twitter Secret`                                                   | `""`            |

例如以下配置

```env
##################### 用户端专有配置 #####################
USER_HTTP_PORT=9000 # 用户端的 HTTP 监听端口. 默认 8080
USER_HTTP_DOMAIN=http://localhost:9000 # 用户端的 API 域名
USER_TOKEN_SECRET_KEY=user # 用户端的 JWT token 密钥
USER_TLS_CERT="" # TLS 的证书文件
USER_TLS_KEY="" # TLS 的 key 文件

##################### 管理员专有配置 #####################
ADMIN_HTTP_PORT=9091 # 管理员端的 HTTP 监听端口. 默认 8081
ADMIN_HTTP_DOMAIN=http://localhost:9091 # 用户端的 API 域名
ADMIN_TOKEN_SECRET_KEY=admin # 管理员端的 JWT token 密钥
ADMIN_TLS_CERT="" # TLS 的证书文件
ADMIN_TLS_KEY="" # TLS 的 key 文件
ADMIN_DEFAULT_PASSWORD="admin" # 默认的超级管理员 admin 的密码，在第一次启动时，会向数据库添加一个超级管理员帐号。默认值: admin

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
DB_SYNC=off # 在应用启动时，是否同步数据库表, 可选 on/off, 默认 off

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

# 短信服务设置
TELEPHONE_PROVIDER="aliyun" # 选用哪一家的短信服务，可选 `aliyun`

# 阿里云短信
TELEPHONE_ALIYUN_ACCESS_KEY="${TELEPHONE_ALIYUN_ACCESS_KEY}" # 阿里云的 access key
TELEPHONE_ALIYUN_ACCESS_SECRET="${TELEPHONE_ALIYUN_ACCESS_SECRET}" # 阿里云的 access secret
TELEPHONE_ALIYUN_SIGN_NAME="${TELEPHONE_ALIYUN_SIGN_NAME}" # 阿里云短信的签名名称
TELEPHONE_ALIYUN_TEMPLATE_CODE_AUTH="${TELEPHONE_ALIYUN_TEMPLATE_CODE_AUTH}" # 用于发送身份验证的短信模版代码
TELEPHONE_ALIYUN_TEMPLATE_CODE_RESET_PASSWORD="${TELEPHONE_ALIYUN_TEMPLATE_CODE_RESET_PASSWORD}" # 用于发送重置密码的短信模版代码
TELEPHONE_ALIYUN_TEMPLATE_CODE_REGISTER="${TELEPHONE_ALIYUN_TEMPLATE_CODE_REGISTER}" # 用于发送注册帐号的短信模版代码

# 腾讯云短信
TELEPHONE_TENCENT_APP_ID="${TELEPHONE_TENCENT_APP_ID}" # sdkappid请填写您在 短信控制台 添加应用后生成的实际 SDK AppID
TELEPHONE_TENCENT_APP_KEY="${TELEPHONE_TENCENT_APP_KEY}" # sdkappid 对应的 appkey，需要业务方高度保密
TELEPHONE_TENCENT_SIGN="${TELEPHONE_TENCENT_SIGN}" # 短信签名内容，使用 UTF-8 编码，必须填写已审核通过的签名。签名信息可登录 短信控制台 查看
TELEPHONE_TENCENT_TEMPLATE_CODE_AUTH="${TELEPHONE_TENCENT_TEMPLATE_CODE_AUTH}" # 用于发送身份验证的短信模版代码
TELEPHONE_TENCENT_TEMPLATE_CODE_RESET_PASSWORD="${TELEPHONE_TENCENT_TEMPLATE_CODE_RESET_PASSWORD}" # 用于发送重置密码的短信模版代码
TELEPHONE_TENCENT_TEMPLATE_CODE_REGISTER="${TELEPHONE_TENCENT_TEMPLATE_CODE_REGISTER}" # 用于发送注册帐号的短信模版代码

# 消息队列配置
MSG_QUEUE_SERVER = 127.0.0.1 # 消息队列服务器地址. 默认 127.0.0.1
MSG_QUEUE_PORT = 4150 # 消息队列服务器端口. 默认 4150

# OAuth2 认证服务
OAUTH_REDIRECT_URL="${OAUTH_REDIRECT_URL}" # 认证成功后，跳转到前端的 URL 地址, 携带 code 给前端拿到用户相关的 token
GITHUB_KEY="${GITHUB_KEY}"
GITHUB_SECRET="${GITHUB_SECRET}"
GITLAB_KEY="${GITLAB_KEY}"
GITLAB_SECRET="${SECRET}"
GOOGLE_KEY="${GOOGLE_KEY}"
GOOGLE_SECRET="${GOOGLE_SECRET}"
FACEBOOK_KEY="${FACEBOOK_KEY}"
FACEBOOK_SECRET="${SECRET}"
TWITTER_KEY="${TWITTER_KEY}"
TWITTER_SECRET="${TWITTER_SECRET}"

# 微信小程序认证
WECHAT_APP_ID = "${WECHAT_APP_ID}"
WECHAT_SECRET = "${WECHAT_SECRET}"
```
