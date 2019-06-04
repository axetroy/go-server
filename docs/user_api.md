## 用户接口

### 验证类

<details><summary>用户注册 <code>[POST] /v1/auth/signup</code></summary>
<p>

请求参数

| 参数        | 类型     | 说明                                                                      | 必选 |
| ----------- | -------- | ------------------------------------------------------------------------- | ---- |
| username    | `string` | 通过用户名来注册, username, email, phone 三选一                           |      |
| email       | `string` | 通过邮箱来注册, username, email, phone 三选一                             |      |
| phone       | `string` | 通过手机来注册, username, email, phone 三选一, 目前手机注册无法发送验证码 |      |
| password    | `string` | 账号密码                                                                  | \*   |
| invite_code | `string` | 邀请码                                                                    |      |

</p>

</details>

<details><summary>用户登陆 <code>[POST] /v1/auth/signin</code></summary>
<p>

| 参数     | 类型     | 说明                                     | 必选 |
| -------- | -------- | ---------------------------------------- | ---- |
| account  | `string` | 用户账号, username/email/phone 中的一个  | \*   |
| password | `string` | 账号密码                                 | \*   |
| code     | `string` | TODO: 手机验证码, 手机可以通过验证码登陆 |      |

</p>

</details>

<details><summary>账号激活 <code>[POST] /v1/auth/activation</code></summary>
<p>

| 参数 | 类型     | 说明                                        | 必选 |
| ---- | -------- | ------------------------------------------- | ---- |
| code | `string` | 激活码，激活码来自服务器发到的邮箱/手机短信 | \*   |

</p>

</details>

<details><summary>忘记密码 <code>[POST] /v1/auth/password/reset</code></summary>
<p>

| 参数         | 类型     | 说明                                        | 必选 |
| ------------ | -------- | ------------------------------------------- | ---- |
| code         | `string` | 重置码，重置码来自服务器发到的邮箱/手机短信 | \*   |
| new_password | `string` | 新的密码                                    |      | \* |

</p>

</details>

### oAuth2

<details><summary>Google 登陆 <code>[GET] /v1/oauth2/google</code></summary>
<p>

前端跳转到这个 URL 进行 Google 授权登陆

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

| 参数          | 类型     | 说明   | 必选 |
| ------------- | -------- | ------ | ---- |
| old_passworld | `string` | 旧密码 | \*   |
| new_password  | `string` | 新密码 | \*   |

</p>

</details>

<details><summary>设置二级密码<code>[POST] /v1/user/password2</code></summary>
<p>

| 参数             | 类型     | 说明         | 必选 |
| ---------------- | -------- | ------------ | ---- |
| password         | `string` | 二级密码     | \*   |
| password_confirm | `string` | 二级密码确认 | \*   |

</p>

</details>

<details><summary>修改二级密码<code>[PUT] /v1/user/password2</code></summary>
<p>

| 参数         | 类型     | 说明       | 必选 |
| ------------ | -------- | ---------- | ---- |
| old_password | `string` | 旧二级密码 | \*   |
| new_password | `string` | 新二级密码 | \*   |

</p>

</details>

<details><summary>发送重置二级密码的邮件/短信<code>[POST] /v1/user/password2/reset</code></summary>
<p>

如果用户有手机，则发送手机验证码，如果有邮箱，则发送邮件

</p>

</details>

<details><summary>重置二级密码<code>[PUT] /v1/user/password2/reset</code></summary>
<p>

| 参数         | 类型     | 说明             | 必选 |
| ------------ | -------- | ---------------- | ---- |
| code         | `string` | 二级密码的重置码 | \*   |
| new_password | `string` | 新二级密码       | \*   |

</p>

</details>

<details><summary>我的邀请列表<code>[GET] /v1/user/invite</code></summary>
<p>

获取我的邀请列表

</p>

</details>

<details><summary>获取单条邀请信息<code>[GET] /v1/user/invite/i/:invite_id</code></summary>
<p>

| 参数      | 类型     | 说明          | 必选 |
| --------- | -------- | ------------- | ---- |
| invite_id | `string` | 邀请数据的 ID | \*   |

</p>

</details>

<details><summary>上传头像<code>[POST] /v1/user/avatar</code></summary>
<p>

头像上传为 Form 表单

| 参数 | 类型   | 说明                                  | 必选 |
| ---- | ------ | ------------------------------------- | ---- |
| file | `file` | 要上传的头像图片，仅支持 jpg/jpeg/png | \*   |

</p>

</details>

### 收货地址

<details><summary>添加收货地址<code>[POST] /v1/user/address</code></summary>
<p>

| 参数          | 类型     | 说明                       | 必选 |
| ------------- | -------- | -------------------------- | ---- |
| name          | `string` | 收件人                     | \*   |
| phone         | `string` | 收件人手机号               | \*   |
| province_code | `string` | 省份代码，6 位数           | \*   |
| city_code     | `string` | 城市代码，6 位数           | \*   |
| area_code     | `string` | 县城代码，6 位数           | \*   |
| address       | `string` | 详细地址，具体的街道门牌号 | \*   |
| is_default    | `bool`   | 是否设置为默认地址         | \*   |

</p>

</details>

<details><summary>更新收货地址<code>[PUT] /v1/user/address/a/:address_id</code></summary>
<p>

| 参数          | 类型     | 说明                       | 必选 |
| ------------- | -------- | -------------------------- | ---- |
| name          | `string` | 收件人                     |      |
| phone         | `string` | 收件人手机号               |      |
| province_code | `string` | 省份代码，6 位数           |      |
| city_code     | `string` | 城市代码，6 位数           |      |
| area_code     | `string` | 县城代码，6 位数           |      |
| address       | `string` | 详细地址，具体的街道门牌号 |      |
| is_default    | `bool`   | 是否设置为默认地址         |      |

</p>

</details>

<details><summary>删除收货地址<code>[DELETE] /v1/user/address/a/:address_id</code></summary>
<p>

删除收货地址

</p>

</details>

<details><summary>收货地址列表<code>[GET] /v1/user/address</code></summary>
<p>

获取我的收货地址列表

</p>

</details>

<details><summary>获取默认收货地址<code>[GET] /v1/user/address/default</code></summary>
<p>

获取我的默认收货地址

</p>

</details>

<details><summary>获取某一个地址<code>[GET] /v1/user/address/a/:address_id</code></summary>
<p>

获取某一个地址的详细信息

</p>

</details>

<details><summary>获取全国地区码列表<code>[GET] /v1/area</code></summary>
<p>

获取全国地区码列表

</p>

</details>

### 钱包类

<details><summary>获取我的钱包<code>[GET] /v1/wallet</code></summary>
<p>

获取我的钱包 Map.

</p>

</details>

<details><summary>获取单个钱包信息<code>[GET] /v1/wallet/w/:currency</code></summary>
<p>

获取指定一个钱包的详细信息.

</p>

</details>

<details><summary>钱包转账<code>[POST] /v1/transfer</code></summary>
<p>

需要在请求头设置 `X-Pay-Password`, 指定二级密码.

| 参数     | 类型     | 说明                    | 必选 |
| -------- | -------- | ----------------------- | ---- |
| currency | `string` | 钱包类型                | \*   |
| to       | `string` | 转账对象的用户纯数字 ID | \*   |
| amount   | `string` | 转账金额                | \*   |
| note     | `string` | 转账备注                |      |

</p>

</details>

<details><summary>获取转账记录<code>[GET] /v1/transfer</code></summary>
<p>

获取我的转账记录

</p>

</details>

<details><summary>获取转账记录详情<code>[GET] /v1/transfer/t/:transfer_id</code></summary>
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

<details><summary>系统通知列表<code>[GET] /v1/notification</code></summary>
<p>

获取系统通知列表

</p>

</details>

<details><summary>系统通知详情<code>[GET] /v1/notification/n/:notification_id</code></summary>
<p>

获取某个系统通知详情

</p>

</details>

<details><summary>标记系统通知已读<code>[PUT] /v1/notification/n/:notification_id/read</code></summary>
<p>

标记系统通知为已读

</p>

</details>

### 个人消息类

<details><summary>个人消息列表<code>[GET] /v1/message</code></summary>
<p>

获取我的消息列表

</p>

</details>

<details><summary>个人消息详情<code>[GET] /v1/message/m/:message_id</code></summary>
<p>

获取某个系统通知详情

</p>

</details>

<details><summary>标记个人消息已读<code>[PUT] /v1/message/m/:message_id/read</code></summary>
<p>

标记个人消息为已读

</p>

</details>

<details><summary>删除个人消息<code>[DELETE] /v1/message/m/:message_id</code></summary>
<p>

删除一条个人消息

</p>

</details>

### 新闻资讯类

<details><summary>资讯列表<code>[GET] /v1/news</code></summary>
<p>

获取资讯列表

</p>

</details>

<details><summary>资讯详情<code>[GET] /v1/news/n/:news_id</code></summary>
<p>

获取某个资讯详情

</p>

</details>

### 邮件服务

> 要使用邮件服务，需要在 `.env` 文件中配置 SMTP 服务

<details><summary>发送账号激活邮件<code>[POST] /v1/email/send/activation</code></summary>
<p>

发送账号激活邮件

| 参数 | 类型     | 说明             | 必选 |
| ---- | -------- | ---------------- | ---- |
| to   | `string` | 要激活的账号邮箱 | \*   |

</p>

</details>

<details><summary>发送登陆密码重置邮件<code>[POST] /v1/email/send/password/reset</code></summary>
<p>

发送账号激活邮件

| 参数 | 类型     | 说明             | 必选 |
| ---- | -------- | ---------------- | ---- |
| to   | `string` | 要激活的账号邮箱 | \*   |

</p>

</details>

### 上传类

<details><summary>上传文件<code>[POST] /v1/upload/file</code></summary>
<p>

Form 表单文件上传, 目前仅支持单个文件上传

| 参数 | 类型     | 说明         | 必选 |
| ---- | -------- | ------------ | ---- |
| file | `string` | 要上传的文件 | \*   |

</p>

</details>

<details><summary>上传图片<code>[POST] /v1/upload/image</code></summary>
<p>

Form 表单图片上传, 目前仅支持单张图片上传

| 参数 | 类型     | 说明         | 必选 |
| ---- | -------- | ------------ | ---- |
| file | `string` | 要上传的图片 | \*   |

</p>

</details>

### 下载类

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

下载缩略图, `filename` 为上传时返回的字段

</p>

</details>

### 资源类

<details><summary>获取上传文件的纯文本<code>[GET] /v1/resource/file/:filename</code></summary>
<p>

获取上传文件的纯文本, `filename` 为上传时返回的字段

</p>

</details>

<details><summary>获取上传的图片<code>[GET] /v1/resource/image/:filename</code></summary>
<p>

获取上传的图片, `filename` 为上传时返回的字段

</p>

</details>

<details><summary>获取上传的缩略图<code>[GET] /v1/resource/thumbnail/:filename</code></summary>
<p>

获取上传的缩略图, `filename` 为上传时返回的字段

</p>

</details>

### 静态文件服务

<details><summary>静态文件<code>[GET] /v1/public/:filename</code></summary>

<p>

在 `public` 目录下的静态文件服务

</p>

</details>

### Banner 轮播图

<details><summary>获取 banner 列表<code>[GET] /v1/banner</code></summary>

<p>

获取 banner 列表

</p>

</details>

<details><summary>获取 banner 详情<code>[GET] /v1/banner/b/:banner_id</code></summary>

<p>

获取一条 banner 的详情

</p>

</details>

### 用户反馈

<details><summary>获取我的反馈列表<code>[GET] /v1/report</code></summary>

<p>

获取我的反馈列表

| 参数   | 类型     | 说明                                   | 必选 |
| ------ | -------- | -------------------------------------- | ---- |
| type   | `string` | 根据`类型`筛选                         |      |
| status | `int`    | 根据`状态`筛选, `0` 未解决, `1` 已解决 |      |

</p>

</details>

<details><summary>提交反馈<code>[POST] /v1/report</code></summary>

<p>

提交一条用户反馈

| 参数        | 类型       | 说明                                                                 | 必选 |
| ----------- | ---------- | -------------------------------------------------------------------- | ---- |
| title       | `string`   | 标题                                                                 | \*   |
| content     | `string`   | 内容                                                                 | \*   |
| type        | `string`   | 反馈的类型, 目前支持 `bug`/`feature`/`suggestion`/`other`            | \*   |
| screenshots | `[]string` | 反馈附带的截图，一个数组的图片路径，该路径是用图片上传接口得到的路径 |      |

</p>

</details>

<details><summary>获取反馈详情<code>[GET] /v1/report/r/:report_id</code></summary>

<p>

获取一条反馈的详情

</p>

</details>
