[![Build Status](https://travis-ci.com/axetroy/go-server.svg?token=QMG6TLRNwECnaTsy6ssj&branch=master)](https://travis-ci.com/axetroy/go-server)
[![Coverage Status](https://coveralls.io/repos/github/axetroy/go-server/badge.svg?branch=master)](https://coveralls.io/github/axetroy/go-server?branch=master)
[![Go Report Card](https://goreportcard.com/badge/github.com/axetroy/go-server)](https://goreportcard.com/report/github.com/axetroy/go-server)
![License](https://img.shields.io/github/license/axetroy/go-server.svg)
![Repo Size](https://img.shields.io/github/repo-size/axetroy/go-server.svg)

### Golang 实现的基础服务

这是我在闲暇时间写的一些基础服务

写一些工作中常用的服务和实现，以备在以后中用到

想到哪里写哪里, 我会不断的完善它

### 包含哪些模块?

| 模块         | 说明                                                                            |
| ------------ | ------------------------------------------------------------------------------- |
| 验证模块     | 包含`注册`/`登陆`/`账号激活`/`忘记密码`/`双重身份认证`等                        |
| 用户模块     | 包含用户信息的模块, `用户资料`/`登陆密码`/`交易密码`/`用户邀请`等               |
| 授权模块     | oAuth授权登陆, 目前只支持 `Google` 账号登陆, 未来可能包括 `微信/QQ/Github`登陆 |
| 钱包模块     | 包含钱包的相关操作，`钱包转账`/`结算等`                                         |
| 财务模块     | 包含财务相关的操作，所有涉及到钱的的操作都在此模块                              |
| 横幅模块     | 对于网站相关 Banner 的操作                                                      |
| 新闻模块     | 新闻公告类的相关操作, CMS内容                                                   |
| 系统通知     | 系统通知的相关模块，主要用于管理员发送给全员的通知                              |
| 消息模块     | 主要用于管理员发送给指定的某个用户消息                                          |
| 地址模块     | 用户设置相关的地址模块，例如`收货地址`等                                        |
| 上传模块     | 包含用户`上传文件/图片`的相关操作, 包含 `hash 去重`/`图片压缩`/`生成缩略图`等   |
| 下载模块     | 包含用户`下载文件/图片`的相关操作                                               |
| 邮件模块     | 关于邮件的相关操作，例如`发送邮件`, 用于`发送验证码`之类                        |
| 静态文件模块 | 用户访问服务器的静态文件, 放置与 `/public` 目录下的文件                         |

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
> docker-compose up # 启动数据库和其他必要的依赖服务
> go run ./cmd/user/main.go # 运行用户端的接口服务
> go run ./cmd/admin/main.go # 运行管理员端的接口服务
```

## 如何构建?

```bash
make build
```

在生成的 bin 目录下查找对应平台的可执行文件运行即可

## 如何测试?

```bash
make test
```

## TODO

- [ ] RBAC 的权限控制模型
- [ ] i18n 的错误信息
- [ ] 启用消息队列
- [ ] 提供 RPC 接口
- [ ] 数据库动态分表

## License

The [MIT License](https://github.com/axetroy/go-server/blob/master/LICENSE)
