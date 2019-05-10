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
- [ ] Banner 轮播图

## TODO

- [ ] i18n 的错误信息
- [ ] 分离出管理员接口
- [ ] 启用消息队列
- [ ] 提供 RPC 接口
- [ ] 数据库动态分表

## 用户接口

### 验证类

<details><summary>用户注册 <code>[POST] /v1/auth/signup</code></summary>
<p>

请求参数

| 参数        | 说明                                                                      | 必选 |
| ----------- | ------------------------------------------------------------------------- | ---- |
| username    | 通过用户名来注册, username, email, phone 三选一                           |      |
| email       | 通过邮箱来注册, username, email, phone 三选一                             |      |
| phone       | 通过手机来注册, username, email, phone 三选一, 目前手机注册无法发送验证码 |      |
| password    | 账号密码                                                                  | *    |
| invite_code | 邀请码                                                                    |      |

</p>

</details>

<details><summary>用户登陆 <code>[POST] /v1/auth/signin</code></summary>
<p>

| 参数     | 说明                                     | 必选 |
| -------- | ---------------------------------------- | ---- |
| account  | 用户账号, username/email/phone中的一个   | *    |
| password | 账号密码                                 | *    |
| code     | TODO: 手机验证码, 手机可以通过验证码登陆 |      |

</p>

</details>

<details><summary>账号激活 <code>[POST] /v1/auth/activation</code></summary>
<p>

| 参数 | 说明                                        | 必选 |
| ---- | ------------------------------------------- | ---- |
| code | 激活码，激活码来自服务器发到的邮箱/手机短信 | *    |

</p>

</details>

<details><summary>忘记密码 <code>[POST] /v1/auth/password/reset</code></summary>
<p>

| 参数         | 说明                                        | 必选 |
| ------------ | ------------------------------------------- | ---- |
| code         | 重置码，重置码来自服务器发到的邮箱/手机短信 | *    |
| new_password | 新的密码                                    |      | * |

</p>

</details>

### 用户类

<details><summary>获取用户信息<code>[GET] /v1/user/profile</code></summary>
<p>

获取用户的详细信息资料

</p>

</details>

<details><summary>更新用户信息<code>[PUT] /v1/user/profile</code></summary>
<p>

| 参数     | 说明         | 必选 |
| -------- | ------------ | ---- |
| nickname | 用户昵称     |      |
| gender   | 用户性别     |      |
| avatar   | 用户头像 URL |      |

</p>

</details>

<details><summary>修改登陆密码<code>[PUT] /v1/user/password</code></summary>
<p>

| 参数          | 说明   | 必选 |
| ------------- | ------ | ---- |
| old_passworld | 旧密码 | *    |
| new_password  | 新密码 | *    |

</p>

</details>

<details><summary>设置二级密码<code>[POST] /v1/user/password2</code></summary>
<p>

| 参数             | 说明         | 必选 |
| ---------------- | ------------ | ---- |
| password         | 二级密码     | *    |
| password_confirm | 二级密码确认 | *    |

</p>

</details>

<details><summary>修改二级密码<code>[PUT] /v1/user/password2</code></summary>
<p>

| 参数         | 说明       | 必选 |
| ------------ | ---------- | ---- |
| old_password | 旧二级密码 | *    |
| new_password | 新二级密码 | *    |

</p>

</details>

<details><summary>发送重置二级密码的邮件/短信<code>[POST] /v1/user/password2/reset</code></summary>
<p>

如果用户有手机，则发送手机验证码，如果有邮箱，则发送邮件

</p>

</details>

<details><summary>重置二级密码<code>[PUT] /v1/user/password2/reset</code></summary>
<p>

| 参数         | 说明             | 必选 |
| ------------ | ---------------- | ---- |
| code         | 二级密码的重置码 | *    |
| new_password | 新二级密码       | *    |

</p>

</details>

<details><summary>我的邀请列表<code>[GET] /v1/user/invite/list</code></summary>
<p>

获取我的邀请列表

</p>

</details>

<details><summary>获取单条邀请信息<code>[GET] /v1/user/invite/detail/:invite_id</code></summary>
<p>

| 参数      | 说明         | 必选 |
| --------- | ------------ | ---- |
| invite_id | 邀请数据的ID | *    |

</p>

</details>

<details><summary>上传头像<code>[POST] /v1/user/avatar</code></summary>
<p>

头像上传为 Form 表单

| 参数 | 说明                                  | 必选 |
| ---- | ------------------------------------- | ---- |
| file | 要上传的头像图片，仅支持 jpg/jpeg/png | *    |

</p>

</details>

### 收货地址

<details><summary>添加收货地址<code>[POST] /v1/user/address/create</code></summary>
<p>

| 参数          | 说明                       | 必选 |
| ------------- | -------------------------- | ---- |
| name          | 收件人                     | *    |
| phone         | 收件人手机号               | *    |
| province_code | 省份代码，6位数            | *    |
| city_code     | 城市代码，6位数            | *    |
| area_code     | 县城代码，6位数            | *    |
| address       | 详细地址，具体的街道门牌号 | *    |
| is_default    | 是否设置为默认地址         | *    |

</p>

</details>

<details><summary>更新收货地址<code>[PUT] /v1/user/address/update/:address_id</code></summary>
<p>

| 参数          | 说明                       | 必选 |
| ------------- | -------------------------- | ---- |
| name          | 收件人                     |      |
| phone         | 收件人手机号               |      |
| province_code | 省份代码，6位数            |      |
| city_code     | 城市代码，6位数            |      |
| area_code     | 县城代码，6位数            |      |
| address       | 详细地址，具体的街道门牌号 |      |
| is_default    | 是否设置为默认地址         |      |

</p>

</details>

<details><summary>删除收货地址<code>[DELETE] /v1/user/address/delete/:address_id</code></summary>
<p>

删除收货地址

</p>

</details>

<details><summary>收货地址列表<code>[GET] /v1/user/address/list</code></summary>
<p>

获取我的收货地址列表

</p>

</details>

<details><summary>获取默认收货地址<code>[GET] /v1/user/address/default</code></summary>
<p>

获取我的默认收货地址

</p>

</details>

<details><summary>获取全国地区码列表<code>[GET] /v1/area</code></summary>
<p>

获取全国地区码列表

</p>

</details>

### 钱包类


<details><summary>获取我的钱包<code>[GET] /v1/wallet/map</code></summary>
<p>

获取我的钱包 Map.

</p>

</details>

<details><summary>获取单个钱包信息<code>[GET] /v1/wallet/currency/:currency</code></summary>
<p>

获取指定一个钱包的详细信息.

</p>

</details>

<details><summary>钱包转账<code>[POST] /v1/transfer</code></summary>
<p>

需要在请求头设置 `X-Pay-Password`, 指定二级密码.

| 参数     | 说明                   | 必选 |
| -------- | ---------------------- | ---- |
| currency | 钱包类型               | *    |
| to       | 转账对象的用户纯数字ID | *    |
| amount   | 转账金额               | *    |
| note     | 转账备注               |      |

</p>

</details>

<details><summary>获取转账记录<code>[GET] /v1/transfer/history</code></summary>
<p>

获取我的转账记录

</p>

</details>

<details><summary>获取转账记录详情<code>[GET] /v1/transfer/detail/:transfer_id</code></summary>
<p>

获取某一条转账记录的详情

</p>

</details>

### 财务类

<details><summary>财务日志<code>[GET] /v1/finance/history</code></summary>
<p>

获取财务日志

</p>

</details>

### 系统通知类

<details><summary>系统通知列表<code>[GET] /v1/notification/list</code></summary>
<p>

获取系统通知列表

</p>

</details>

<details><summary>系统通知详情<code>[GET] /v1/notification/detail/:notification_id</code></summary>
<p>

获取某个系统通知详情

</p>

</details>

<details><summary>标记系统通知已读<code>[GET] /v1/notification/read/:notification_id</code></summary>
<p>

标记系统通知为已读

</p>

</details>

### 个人消息类

<details><summary>个人消息列表<code>[GET] /v1/message/list</code></summary>
<p>

获取个人列表

</p>

</details>

<details><summary>个人消息详情<code>[GET] /v1/message/detail/:message_id</code></summary>
<p>

获取某个系统通知详情

</p>

</details>

<details><summary>标记个人消息已读<code>[GET] /v1/message/read/:message_id</code></summary>
<p>

标记个人消息为已读

</p>

</details>

### 新闻资讯类

<details><summary>资讯列表<code>[GET] /v1/news/list</code></summary>
<p>

获取资讯列表

</p>

</details>

<details><summary>资讯详情<code>[GET] /v1/news/detail/:news_id</code></summary>
<p>

获取某个资讯详情

</p>

</details>

### 邮件服务

> 要使用邮件服务，需要在 `.env` 文件中配置 SMTP 服务

<details><summary>发送账号激活邮件<code>[POST] /v1//email/send/activation</code></summary>
<p>

发送账号激活邮件

| 参数 | 说明             | 必选 |
| ---- | ---------------- | ---- |
| to   | 要激活的账号邮箱 | *    |

</p>

</details>

<details><summary>发送登陆密码重置邮件<code>[POST] /v1//email/send/reset_password</code></summary>
<p>

发送账号激活邮件

| 参数 | 说明             | 必选 |
| ---- | ---------------- | ---- |
| to   | 要激活的账号邮箱 | *    |

</p>

</details>

### 文件下载

<details><summary>上传文件<code>[POST] /v1//upload/file</code></summary>
<p>

Form 表单文件上传, 目前仅支持单个文件上传

| 参数 | 说明         | 必选 |
| ---- | ------------ | ---- |
| file | 要上传的文件 | *    |

</p>

</details>

<details><summary>上传图片<code>[POST] /v1//upload/image</code></summary>
<p>

Form 表单图片上传, 目前仅支持单张图片上传

| 参数 | 说明         | 必选 |
| ---- | ------------ | ---- |
| file | 要上传的图片 | *    |

</p>

</details>

<details><summary>下载文件<code>[GET] /v1/download/file/:filename</code></summary>
<p>

下载文件, `filename` 为上传时返回的字段

</p>

</details>

<details><summary>下载图片<code>[GET] /v1/download/image/:filename</code></summary>
<p>

下载图片, `filename` 为上传时返回的字段

</p>

</details>

<details><summary>下载缩略图<code>[GET] /v1/download/thumbnail/:filename</code></summary>
<p>

下载图片, `filename` 为上传时返回的字段

</p>

</details>

### 静态文件服务

<details><summary>静态文件<code>[GET] /v1/public/:filename</code></summary>

<p>

在 `public` 目录下的静态文件服务

</p>

</details>

## 管理员接口

### 用户类

<details><summary>登陆<code>[POST] /v1/admin/login</code></summary>

<p>

| 参数     | 说明       | 必选 |
| -------- | ---------- | ---- |
| username | 管理员账号 | *    |
| password | 账号密码   | *    |

</p>

</details>

<details><summary>获取管理员信息<code>[GET] /v1/admin/profile</code></summary>


<p>

获取管理员的个人信息

</p>

</details>


### 管理员类

<details><summary>创建管理员<code>[POST] /v1/admin/create</code></summary>

仅限于超级管理员

<p>

| 参数     | 说明                       | 必选 |
| -------- | -------------------------- | ---- |
| account  | 管理员账号                 | *    |
| password | 账号密码                   | *    |
| name     | 管理员名称，注册后不可修改 | *    |

</p>

</details>

### 新闻资讯类

<details><summary>添加新闻资讯<code>[POST] /v1/admin/news/create</code></summary>

<p>

| 参数    | 说明                                                          | 必选 |
| ------- | ------------------------------------------------------------- | ---- |
| title   | 资讯标题                                                      | *    |
| content | 资讯内容                                                      | *    |
| type    | 资讯的类型, 取值 `news`(新闻资讯) or `announcement`(官方公告) | *    |
| tags    | 资讯标签，字符串数组                                          |      |

</p>

</details>

<details><summary>更新新闻资讯<code>[PUT] /v1/admin/news/update/:news_id</code></summary>

<p>

| 参数    | 说明                                                          | 必选 |
| ------- | ------------------------------------------------------------- | ---- |
| title   | 资讯标题                                                      |      |
| content | 资讯内容                                                      |      |
| type    | 资讯的类型, 取值 `news`(新闻资讯) or `announcement`(官方公告) |      |
| tags    | 资讯标签，字符串数组                                          |      |

</p>

</details>

### 系统通知类

<details><summary>新增系统通知<code>[POST] /v1/admin/notification/create</code></summary>

<p>

| 参数    | 说明     | 必选 |
| ------- | -------- | ---- |
| title   | 通知标题 | *    |
| content | 通知内容 | *    |
| note    | 备注     |      |

</p>

</details>

<details><summary>修改系统通知<code>[PUT] /v1/admin/notification/update/:notification_id</code></summary>

<p>

| 参数    | 说明     | 必选 |
| ------- | -------- | ---- |
| title   | 通知标题 |      |
| content | 通知内容 |      |
| note    | 备注     |      |

</p>

</details>

### 个人消息类

<details><summary>新增个人消息<code>[POST] /v1/admin/message/create</code></summary>

<p>

| 参数    | 说明     | 必选 |
| ------- | -------- | ---- |
| uid     | 用户ID   |      |
| title   | 通知标题 |      |
| content | 通知内容 |      |

</p>

</details>

## License

The [MIT License](https://github.com/axetroy/go-server/blob/master/LICENSE)
