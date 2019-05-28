[![Build Status](https://travis-ci.com/axetroy/go-server.svg?token=QMG6TLRNwECnaTsy6ssj&branch=master)](https://travis-ci.com/axetroy/go-server)
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

- [x] Web 框架 [Gin](https://github.com/gin-gonic/gin)
- [x] 数据库 Postgres
- [x] 缓存 Redis
- [x] 身份认证机制 [JWT](http://jwt.io)/[oAuth2](https://oauth.net/2/)
- [x] 数据库操作 [GORM](https://github.com/jinzhu/gorm)
- [x] 消息队列 [nsq](https://github.com/nsqio/nsq)
- [ ] RBAC 的权限模型
- [x] Docker 一键启动应用

### 包含哪些模块?

| 模块         | 说明                                                                                       |
| ------------ | ------------------------------------------------------------------------------------------ |
| 验证模块     | 包含`注册`/`登陆`/`账号激活`/`忘记密码`/`双重身份认证`等                                   |
| 用户模块     | 包含用户信息的模块, `用户资料`/`登陆密码`/`交易密码`/`用户邀请`等                          |
| 授权模块     | oAuth授权登陆, 目前只支持 `Google` 账号登陆, 未来可能包括 `微信/QQ/Github`登陆             |
| 钱包模块     | 包含钱包的相关操作，`钱包转账`/`结算`等                                                    |
| 财务模块     | 所有涉及到钱的的操作都会被记录在此, 例如`转账记录`/`消费记录`等                            |
| 横幅模块     | 对于网站相关 `Banner` 的操作，可根据不同的平台设置不同的横幅，例如 PC 端大屏的与 APP不相同 |
| 新闻模块     | 新闻公告类的相关操作, CMS内容                                                              |
| 系统通知     | 系统通知的相关模块，主要用于`管理员发送给全员的通知`                                       |
| 消息模块     | 用户的个人消息模块, 主要用于`管理员发送给某个用户的通知`                                   |
| 地址模块     | 用户设置相关的地址模块，例如`收货地址`等                                                   |
| 上传模块     | 包含用户`上传文件/图片`的相关操作, 包含 `hash 去重`/`图片压缩`/`生成缩略图`等              |
| 下载模块     | 包含用户`下载文件/图片`的相关操作                                                          |
| 邮件模块     | 关于邮件的相关操作，例如`发送邮件`, 用于`发送验证码`之类                                   |
| 静态文件模块 | 用户访问服务器的静态文件, 放置与 `/public` 目录下的文件                                    |
| 反馈模块     | TODO: 用户反馈模块，用户`建议反馈`/`BUG反馈`等                                            |

## 提供了哪些接口？

- [用户端](docs/user_api.md)
- [管理员端](docs/admin_api.md)

## 如何运行?

首选确保你安装有:

- [Golang](https://golang.org/)
- [Docker](https://www.docker.com/)
- [Docker Compose](https://docs.docker.com/compose/)

再根据以下命令运行

```bash
> go get -v github.com/axetroy/go-server # 拉取项目
> cd $GOPATH/github/axetroy/go-server # 切换到项目目录
> docker-compose -f docker-compose.mq.yml up # 启动消息队列
> docker-compose up # 启动数据库和HTTP服务
> go run ./cmd/message_queue/main.go # 启动消息队列
> go run ./cmd/user/main.go # 运行用户端的接口服务
> go run ./cmd/admin/main.go # 运行管理员端的接口服务
```

## 如何构建?

```bash
make build
```

在生成的 bin 目录下查找对应平台的可执行文件运行即可

```bash
> cd $GOPATH/github/axetroy/go-server
> tree ./bin
```

文件说明:

```
./bin
├── admin_linux_x64                         # 管理员端的启动文件
├── admin_linux_x86
├── admin_osx_64
├── admin_osx_x86
├── admin_win_x64.exe
├── admin_win_x86.exe
├── message_queue_linux_x64                 # 消息队列启动文件
├── message_queue_linux_x86
├── message_queue_osx_64
├── message_queue_osx_x86
├── message_queue_win_x64.exe
├── message_queue_win_x86.exe
├── user_linux_x64                          # 用户端的启动文件
├── user_linux_x86
├── user_osx_64
├── user_osx_x86
├── user_win_x64.exe
└── user_win_x86.exe
```

## 如何测试?

```bash
make test
```

## TODO

- [ ] RBAC 的权限控制模型
- [ ] i18n 的错误信息
- [ ] 提供 RPC 接口
- [ ] 数据库动态分表

## License

The [MIT License](https://github.com/axetroy/go-server/blob/master/LICENSE)
