![ci](https://github.com/axetroy/go-server/workflows/ci/badge.svg)
[![Coverage Status](https://coveralls.io/repos/github/axetroy/go-server/badge.svg?branch=master)](https://coveralls.io/github/axetroy/go-server?branch=master)
[![Go Report Card](https://goreportcard.com/badge/github.com/axetroy/go-server)](https://goreportcard.com/report/github.com/axetroy/go-server)
![License](https://img.shields.io/github/license/axetroy/go-server.svg)
![Repo Size](https://img.shields.io/github/repo-size/axetroy/go-server.svg)

### Golang 实现的基础服务

这是我在闲暇时间写的一些基础服务, 基本上大多数后端服务都需要用到的。

它用来帮助我快速开始一个项目，而不是重头开始写，浪费大量时间和精力。

想到哪里写哪里, 我会不断的完善它。

> 为什么不写成微服务形式，模块分离? 项目体量没有达到。

### 技术栈

- [x] Web 框架 [iris](https://github.com/kataras/iris)
- [x] 数据库 Postgres
- [x] 缓存 Redis
- [x] 身份认证机制 Session/[JWT](http://jwt.io)/[oAuth2](https://oauth.net/2/)
- [x] 数据库操作 [GORM](https://github.com/jinzhu/gorm)
- [x] 消息队列 [nsq](https://github.com/nsqio/nsq)
- [x] RBAC 的鉴权模型
- [x] Docker 一键启动应用

### 包含哪些模块?

| 模块         | 说明                                                                                        |
| ------------ | ------------------------------------------------------------------------------------------- |
| 验证模块     | 包含`注册`/`登陆`/`账号激活`/`忘记密码`/`双重身份认证`等                                    |
| 用户模块     | 包含用户信息的模块, `用户资料`/`登陆密码`/`扫码登录`/`交易密码`/`用户邀请`等                |
| 授权模块     | oAuth 授权登陆, 目前支持 `微信小程序`/`Google`/`Github`/`Gitlab`/`Twitter`/`Facebook`       |
| 钱包模块     | 包含钱包的相关操作，`钱包转账`/`结算`等                                                     |
| 财务模块     | 所有涉及到钱的的操作都会被记录在此, 例如`转账记录`/`消费记录`等                             |
| 横幅模块     | 对于网站相关 `Banner` 的操作，可根据不同的平台设置不同的横幅，例如 PC 端大屏的与 APP 不相同 |
| 新闻模块     | 新闻公告类的相关操作, CMS 内容                                                              |
| 系统通知     | 系统通知的相关模块，主要用于`管理员发送给全员的通知`                                        |
| 消息模块     | 用户的个人消息模块, 主要用于`管理员发送给某个用户的通知`                                    |
| 推送模块     | APP 的推送通知模块，接入第三方 `onesignal`                                                  |
| 地区模块     | 包含全国各地区的行政码划分等                                                                |
| 地址模块     | 用户设置相关的地址模块，例如`收货地址`等                                                    |
| 上传模块     | 包含用户`上传文件/图片`的相关操作, 包含 `hash 去重`/`图片压缩`/`图片缩放`/`图片裁剪`等      |
| 下载模块     | 包含用户`下载文件/图片`的相关操作                                                           |
| 邮件模块     | 关于邮件的相关操作，例如`发送邮件`, 用于`发送验证码`之类                                    |
| 短信模块     | 用于发送短信验证码，接入第三方服务`阿里云`/`腾讯云`                                         |
| 静态文件模块 | 用户访问服务器的静态文件, 放置与 `/public` 目录下的文件                                     |
| 反馈模块     | 用户反馈模块，用户`建议反馈`/`BUG反馈`等                                                    |
| 页面菜单模块 | 定义`后台页面菜单`/`页面权限`等                                                             |
| 日志模块     | `系统日志`/`登陆日志`/`操作日志`/`异常日志`等                                               |
| 帮助中心     | 可嵌套的帮助中心模块                                                                        |
| 配置中心     | 由后台接口控制的配置中心，配置例如 SMTP/短信/小程序 等信息                                  |
| 定时任务模块 | 进行一系列的定时任务，例如切割表，迁移冷数据等                                              |
| 客服系统     | 基于 Websocket 的客服系统，实时消息推送，一对多服务关系                                     |

## 如何使用?

首先搭建项目需要的`依赖数据库/服务`, 这里推荐使用 `Docker`。

在 [docker](docker) 目录中提供了 `docker-compose` 配置文件，方便一键搭建。

然后获取[构建好的可执行文件](https://github.com/axetroy/go-server/releases), 找到对应的平台，并且下载。或者自行构建。

你需要使用 5 个文件

1. message_queue_server

> 启用消息队列消费服务器，用于消费在队列里面的事物。

2. resource_server

> 静态文件服务器。用于静态文件的上传/下载/缩略图等

3. user_server

> 监听用户相关的接口服务

4. admin_server

> 监听管理员相关的接口服务

5。scheduled_server

> 定时任务

6. customer_service

> 客服系统

然后复制 [.env](.env) 到可执行文件目录下，运行可执行文件即可。例如 `./user_server start`

快速下载可执行文件

```bash
curl -fsSL https://raw.githubusercontent.com/axetroy/go-server/master/install.sh | bash -s v0.5.2
```

## 如何进行本地开发?

首选确保你安装有:

- [Golang](https://golang.org/) >= 1.14.x
- [Docker](https://www.docker.com/)
- [Docker Compose](https://docs.docker.com/compose/)

再根据以下命令运行

```bash
# 克隆项目
$ go get -v github.com/axetroy/go-server # 拉取项目

# 启用项目依赖(数据库，消息队列等)
$ cd $GOPATH/github/axetroy/go-server # 切换到项目目录
$ cd docker
$ ./start.sh

# 启动接口服务
$ cd $GOPATH/github/axetroy/go-server # 切换到项目目录
$ go run ./cmd/admin/main.go migrate # 同步数据库表
$ go run ./cmd/message_queue/main.go # 启动消息队列
$ go run ./cmd/user/main.go # 运行用户端的接口服务
$ go run ./cmd/admin/main.go # 运行管理员端的接口服务
$ go run ./cmd/resource/main.go # 运行资源类的接口服务
$ go run ./cmd/scheduled/main.go # 运行定时任务
$ go run ./cmd/customer_service/main.go # 运行客服系统
```

可以通过 [.env](.env) 文件进行配置

## 如何构建?

```bash
$ make build
```

在生成的 bin 目录下查找对应平台的可执行文件运行即可

## 如何测试?

```bash
make test
```

## 鸣谢

特别感谢 [JetBrains](https://www.jetbrains.com/?from=go-server) 为本开源项目提供免费的 [IntelliJ GoLand](https://www.jetbrains.com/go/?from=go-server) 授权

<p>
 <a href="https://www.jetbrains.com/?from=ferry">
   <img height="200" src="https://www.fdevops.com/wp-content/uploads/2020/09/1599213857-jetbrains-variant-4.png">
 </a>
</p>

## License

The [MIT License](LICENSE)
