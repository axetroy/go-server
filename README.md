[![Build Status](https://travis-ci.com/axetroy/go-server.svg?token=QMG6TLRNwECnaTsy6ssj&branch=master)](https://travis-ci.com/axetroy/go-server)
[![Coverage Status](https://coveralls.io/repos/github/axetroy/go-server/badge.svg?branch=master)](https://coveralls.io/github/axetroy/go-server?branch=master)
[![Go Report Card](https://goreportcard.com/badge/github.com/axetroy/go-server)](https://goreportcard.com/report/github.com/axetroy/go-server)
![License](https://img.shields.io/github/license/axetroy/go-server.svg)
![Repo Size](https://img.shields.io/github/repo-size/axetroy/go-server.svg)

### Golang 实现的基础服务

这是我在闲暇时间写的一些基础服务

写一些工作中常用的服务和实现，以备在以后中用到

想到哪里写哪里, 我会不断的完善它

### 包含哪些服务

- [x] 验证类
  - [x] 注册
  - [x] 登陆
  - [x] 账号激活
  - [x] 忘记密码
  - [x] 双重身份验证
  - [ ] 接入短信验证码服务商
  - [ ] 图片验证码

- [ ] 用户类
  - [x] 登出
  - [x] 获取用户资料
  - [x] 更改用户资料
  - [x] 修改登陆密码
  - [x] 忘记登陆密码
  - [x] 设置交易密码
  - [x] 修改交易密码
  - [x] 忘记交易密码
  - [x] 获取用户已邀请的用户列表
  - 用户头像
    - [x] 上传用户头像
    - [ ] 第三方头像
  - oAuth2 第三方登陆
    - [ ] 微信
    - [ ] QQ
    - [x] Google
    - [ ] Github
  - [x] 收货地址服务

- [x] 钱包类
  - [x] 用户钱包
  - [x] 钱包转账

- [ ] 财务流水
  - [ ] 财务日志

- [x] 新闻公告
- [x] 系统通知
- [x] 个人通知

- [x] 上传类
  - [x] 文件上传
    - [x] 获取上传的文件
    - [x] 下载上传的文件
    - [x] 限制文件大小/类型
  - [x] 图片上传
    - [x] 生成缩略图
    - [x] 下载图片
    - [x] 限制图片大小/类型
- [x] 邮件服务

- [x] 静态文件服务
- [ ] 帮助中心
- [x] Banner 轮播图
- [ ] 反馈系统(优化建议/BUG反馈等)

## 提供了哪些接口？

- [用户端](docs/user_api.md)
- [管理员端](docs/admin_api.md)

## 如何运行

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

## TODO

- [ ] RBAC 的权限控制模型
- [ ] i18n 的错误信息
- [ ] 启用消息队列
- [ ] 提供 RPC 接口
- [ ] 数据库动态分表

## License

The [MIT License](https://github.com/axetroy/go-server/blob/master/LICENSE)
