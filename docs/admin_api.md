## 管理员接口

### 用户类

<details><summary>管理员登陆<code>[POST] /v1/login</code></summary>

<p>

| 参数     | 说明       | 必选 |
| -------- | ---------- | ---- |
| username | 管理员账号 | \*   |
| password | 账号密码   | \*   |

</p>

</details>

<details><summary>获取会员列表<code>[GET] /v1/user</code></summary>

<p>

获取所有的会员列表

</p>

</details>

<details><summary>获取指定会员详情<code>[GET] /v1/user/u/:user_id</code></summary>

<p>

获取指定会员详情

</p>

</details>

<details><summary>修改会员资料<code>[PUT] /v1/user/u/:user_id</code></summary>

<p>

| 参数     | 说明         | 必选 |
| -------- | ------------ | ---- |
| nickname | 用户昵称     |      |
| gender   | 用户性别     |      |
| avatar   | 用户头像 URL |      |

</p>

</details>

<details><summary>修改会员密码<code>[PUT] /v1/user/u/:user_id/password</code></summary>

<p>

| 参数         | 说明   | 必选 |
| ------------ | ------ | ---- |
| new_password | 新密码 | \*   |

</p>

</details>

### 管理员类

<details><summary>创建管理员<code>[POST] /v1/admin</code></summary>

仅限于超级管理员

<p>

| 参数     | 说明                       | 必选 |
| -------- | -------------------------- | ---- |
| account  | 管理员账号                 | \*   |
| password | 账号密码                   | \*   |
| name     | 管理员名称，注册后不可修改 | \*   |

</p>

</details>

<details><summary>修改管理员信息<code>[PUT] /v1/admin/a/:admin_id</code></summary>

仅限于超级管理员

<p>

| 参数   | 说明                                                         | 必选 |
| ------ | ------------------------------------------------------------ | ---- |
| name   | 管理员名字                                                   |      |
| status | 管理员状态, 可选 `-1`(未激活)/`0`(默认状态)/`-100`(已被禁用) |      |

</p>

</details>

<details><summary>管理员列表<code>[GET] /v1/admin</code></summary>

获取管理员列表

</details>

<details><summary>获取管理员自己的信息<code>[GET] /v1/admin/profile</code></summary>

<p>

获取管理员的个人信息

</p>

</details>

<details><summary>获取指定管理员信息<code>[GET] /v1/admin/a/:admin_id</code></summary>

<p>

获取指定管理员信息

</p>

</details>

### RBAC 鉴权

</details>

<details><summary>获取用户角色列表<code>[GET] /v1/role</code></summary>

<p>

获取当前的用户角色列表

</p>

</details>

<details><summary>获取用户角色详情<code>[GET] /v1/role/r/:name</code></summary>

<p>

获取用户角色详情

</p>

</details>

<details><summary>创建用户角色<code>[POST] /v1/role</code></summary>

<p>

创建一个用户角色

| 参数        | 说明                          | 必选 |
| ----------- | ----------------------------- | ---- |
| name        | 角色名称 `string`, 角色名唯一 | \*   |
| description | 角色描述 `string`             | \*   |
| accession   | 角色所拥有的权限 `[]string`   | \*   |
| note        | 角色备注 `string`             |      |

</p>

</details>

<details><summary>更新用户角色<code>[PUT] /v1/role/r/:name</code></summary>

<p>

更新一个用户角色, `内置角色` 无法更新

| 参数        | 说明                        | 必选 |
| ----------- | --------------------------- | ---- |
| description | 角色描述 `string`           |      |
| accession   | 角色所拥有的权限 `[]string` |      |
| note        | 角色备注 `string`           |      |

</p>

</details>

<details><summary>删除用户角色<code>[DELETE] /v1/role/r/:name</code></summary>

<p>

删除用户角色, `内置角色` 无法删除

> 如果有任何一个用户属于这个角色，则不允许删除

</p>

</details>

<details><summary>更改用户角色<code>[PUT] /v1/role/u/:user_id</code></summary>

<p>

更改用户的角色, 一个用户可以赋予多种角色

| 参数  | 说明                      | 必选 |
| ----- | ------------------------- | ---- |
| roles | 角色名称, 类型 `[]string` | \*   |

</p>

</p>

</details>

<details><summary>获取权限列表<code>[GET] /v1/role/accession</code></summary>

<p>

获取所有权限

</p>

</details>

### 新闻资讯类

<details><summary>添加新闻资讯<code>[POST] /v1/news</code></summary>

<p>

| 参数    | 说明                                                          | 必选 |
| ------- | ------------------------------------------------------------- | ---- |
| title   | 资讯标题                                                      | \*   |
| content | 资讯内容                                                      | \*   |
| type    | 资讯的类型, 取值 `news`(新闻资讯) or `announcement`(官方公告) | \*   |
| tags    | 资讯标签，字符串数组                                          |      |

</p>

</details>

<details><summary>更新新闻资讯<code>[PUT] /v1/news/n/:news_id</code></summary>

<p>

| 参数    | 说明                                                          | 必选 |
| ------- | ------------------------------------------------------------- | ---- |
| title   | 资讯标题                                                      |      |
| content | 资讯内容                                                      |      |
| type    | 资讯的类型, 取值 `news`(新闻资讯) or `announcement`(官方公告) |      |
| tags    | 资讯标签，字符串数组                                          |      |

</p>

</details>

<details><summary>获取单个资讯信息<code>[GET] /v1/news/n/:news_id</code></summary>

<p>

获取单个资讯信息

</p>

</details>

<details><summary>获取资讯列表<code>[GET] /v1/news</code></summary>

<p>

获取单个资讯信息

</p>

</details>

<details><summary>删除资讯<code>[DELETE] /v1/news/n/:news_id</code></summary>

<p>

删除单个资讯

</p>

</details>

### 系统通知类

<details><summary>新增系统通知<code>[POST] /v1/notification</code></summary>

<p>

| 参数    | 说明     | 必选 |
| ------- | -------- | ---- |
| title   | 通知标题 | \*   |
| content | 通知内容 | \*   |
| note    | 备注     |      |

</p>

</details>

<details><summary>修改系统通知<code>[PUT] /v1/notification/n/:notification_id</code></summary>

<p>

| 参数    | 说明     | 必选 |
| ------- | -------- | ---- |
| title   | 通知标题 |      |
| content | 通知内容 |      |
| note    | 备注     |      |

</p>

</details>

<details><summary>删除系统通知<code>[DELETE] /v1/notification/n/:notification_id</code></summary>

<p>

管理员删除系统通知

</p>

</details>

<details><summary>获取系统通知列表<code>[GET] /v1/notification</code></summary>

<p>

管理员获取系统通知列表

</p>

</details>

<details><summary>获取系统通知详情<code>[GET] /v1/notification/n/:notification_id</code></summary>

<p>

管理员获取系统通知详情

</p>

</details>

### 个人消息类

<details><summary>新增个人消息<code>[POST] /v1/message</code></summary>

<p>

| 参数    | 说明     | 必选 |
| ------- | -------- | ---- |
| uid     | 用户 ID  |      |
| title   | 通知标题 |      |
| content | 通知内容 |      |

</p>

</details>

<details><summary>删除个人消息<code>[DELETE] /v1/message/m/:message_id</code></summary>

<p>

删除个人消息

</p>

</details>

<details><summary>更改个人消息<code>[PUT] /v1/message/m/:message_id</code></summary>

<p>

| 参数    | 说明     | 必选 |
| ------- | -------- | ---- |
| title   | 消息标题 |      |
| content | 消息内容 |      |

</p>

</details>

<details><summary>消息列表<code>[GET] /v1/message</code></summary>
<p>

获取我的消息列表

</p>

</details>

<details><summary>消息详情<code>[GET] /v1/message/m/:message_id</code></summary>
<p>

获取某个系统通知详情

</p>

</details>

### Banner 轮播图

<details><summary>新增 banner<code>[POST] /v1/banner</code></summary>

<p>

| 参数         | 说明                                                   | 必选 |
| ------------ | ------------------------------------------------------ | ---- |
| image        | 图片 URL                                               | \*   |
| href         | 图片跳转的链接                                         | \*   |
| platform     | 该 banner 图片运用在哪个平台. 分别为 `PC` 或 `APP`     | \*   |
| description  | 该 banner 的描述信息                                   |      |
| priority     | 优先级，用于排序                                       |      |
| identifier   | APP 跳转标识符, 给 APP 跳转页面用的                    |      |
| fallback_url | 当 APP 的 identifier 无效时的备选方案，跳转的 URL 地址 |      |

</p>

</details>

<details><summary>修改 banner<code>[PUT] /v1/banner/b/:banner_id</code></summary>

<p>

| 参数         | 说明                                                   | 必选 |
| ------------ | ------------------------------------------------------ | ---- |
| image        | 图片 URL                                               |      |
| href         | 图片跳转的链接                                         |      |
| platform     | 该 banner 图片运用在哪个平台. 分别为 `PC` 或 `APP`     |      |
| description  | 该 banner 的描述信息                                   |      |
| priority     | 优先级，用于排序                                       |      |
| identifier   | APP 跳转标识符, 给 APP 跳转页面用的                    |      |
| fallback_url | 当 APP 的 identifier 无效时的备选方案，跳转的 URL 地址 |      |

</p>

</details>

<details><summary>删除 banner<code>[DELETE] /v1/banner/b/:banner_id</code></summary>

<p>

删除一条 banner

</p>

</details>

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

### 系统信息

<details><summary>获取当前服务器的信息<code>[GET] /v1/system</code></summary>

<p>

获取当前服务器的信息, 包括内存/CPU/磁盘等

</p>

</details>
