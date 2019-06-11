## 管理员接口

### 用户类

<details><summary>管理员登陆<code>[POST] /v1/login</code></summary>

<p>

| 参数     | 类型     | 说明       | 必填 |
| -------- | -------- | ---------- | ---- |
| username | `string` | 管理员账号 | \*   |
| password | `string` | 账号密码   | \*   |

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

| 参数     | 类型     | 说明         | 必填 |
| -------- | -------- | ------------ | ---- |
| nickname | `string` | 用户昵称     |      |
| gender   | `string` | 用户性别     |      |
| avatar   | `string` | 用户头像 URL |      |

</p>

</details>

<details><summary>修改会员密码<code>[PUT] /v1/user/u/:user_id/password</code></summary>

<p>

| 参数         | 类型     | 说明   | 必填 |
| ------------ | -------- | ------ | ---- |
| new_password | `string` | 新密码 | \*   |

</p>

</details>

### 管理员类

<details><summary>创建管理员<code>[POST] /v1/admin</code></summary>

仅限于超级管理员

<p>

| 参数     | 类型     | 说明                       | 必填 |
| -------- | -------- | -------------------------- | ---- |
| account  | `string` | 管理员账号                 | \*   |
| password | `string` | 账号密码                   | \*   |
| name     | `string` | 管理员名称，注册后不可修改 | \*   |

</p>

</details>

<details><summary>修改管理员信息<code>[PUT] /v1/admin/a/:admin_id</code></summary>

仅限于超级管理员

<p>

| 参数   | 类型     | 说明                                                         | 必填 |
| ------ | -------- | ------------------------------------------------------------ | ---- |
| name   | `string` | 管理员名字                                                   |      |
| status | `int`    | 管理员状态, 可选 `-1`(未激活)/`0`(默认状态)/`-100`(已被禁用) |      |

</p>

</details>

<details><summary>管理员列表<code>[GET] /v1/admin</code></summary>

获取管理员列表

</details>

<details><summary>获取管理员自己的信息<code>[GET] /v1/profile</code></summary>

<p>

获取管理员的个人信息

</p>

</details>

<details><summary>获取指定管理员信息<code>[GET] /v1/admin/a/:admin_id</code></summary>

<p>

获取指定管理员信息

</p>

</details>

<details><summary>删除管理员<code>[DELETE] /v1/admin/a/:admin_id</code></summary>

<p>

删除删除管理员

</p>

</details>

<details><summary>获取管理员的所有权限列表<code>[GET] /v1/admin/accession</code></summary>

<p>

获取管理员的所有权限列表

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

| 参数        | 类型       | 说明                 | 必填 |
| ----------- | ---------- | -------------------- | ---- |
| name        | `string`   | 角色名称, 角色名唯一 | \*   |
| description | `string`   | 角色描述             | \*   |
| accession   | `[]string` | 角色所拥有的权限列表 | \*   |
| note        | `string`   | 角色备注             |      |

</p>

</details>

<details><summary>更新用户角色<code>[PUT] /v1/role/r/:name</code></summary>

<p>

更新一个用户角色, `内置角色` 无法更新

| 参数        | 类型       | 说明             | 必填 |
| ----------- | ---------- | ---------------- | ---- |
| description | `string`   | 角色描述         |      |
| accession   | `[]string` | 角色所拥有的权限 |      |
| note        | `string`   | 角色备注         |      |

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

| 参数  | 类型       | 说明                                            | 必填 |
| ----- | ---------- | ----------------------------------------------- | ---- |
| roles | `[]string` |  要更改成的角色, 当前角色会覆盖掉用户原有的角色 | \*   |

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

| 参数    | 类型       | 说明                                                         | 必填 |
| ------- | ---------- | ------------------------------------------------------------ | ---- |
| title   | `string`   | 资讯标题                                                     | \*   |
| content | `string`   | 资讯内容                                                     | \*   |
| type    | `string`   | 资讯的类型,取值 `news`(新闻资讯) or `announcement`(官方公告) | \*   |
| tags    | `[]string` | 资讯标签，字符串数组                                         |      |

</p>

</details>

<details><summary>更新新闻资讯<code>[PUT] /v1/news/n/:news_id</code></summary>

<p>

| 参数    | 类型       | 说明                                                          | 必填 |
| ------- | ---------- | ------------------------------------------------------------- | ---- |
| title   | `string`   | 资讯标题                                                      |      |
| content | `string`   | 资讯内容                                                      |      |
| type    | `string`   | 资讯的类型, 取值 `news`(新闻资讯) or `announcement`(官方公告) |      |
| tags    | `[]string` | 资讯标签，字符串数组                                          |      |

</p>

</details>

<details><summary>获取单个资讯信息<code>[GET] /v1/news/n/:news_id</code></summary>

<p>

获取单个资讯信息

</p>

</details>

<details><summary>获取资讯列表<code>[GET] /v1/news</code></summary>

<p>

获取资讯列表

| 参数   | 类型     | 说明       | 必填 |
| ------ | -------- | ---------- | ---- |
| type   | `string` | 资讯的类型 |      |
| status | `string` | 资讯的状态 |      |

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

| 参数    | 类型     | 说明     | 必填 |
| ------- | -------- | -------- | ---- |
| title   | `string` | 通知标题 | \*   |
| content | `string` | 通知内容 | \*   |
| note    | `string` | 备注     |      |

</p>

</details>

<details><summary>修改系统通知<code>[PUT] /v1/notification/n/:notification_id</code></summary>

<p>

| 参数    | 类型     | 说明     | 必填 |
| ------- | -------- | -------- | ---- |
| title   | `string` | 通知标题 |      |
| content | `string` | 通知内容 |      |
| note    | `string` | 备注     |      |

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

| 参数    | 类型     | 说明     | 必填 |
| ------- | -------- | -------- | ---- |
| uid     | `string` | 用户 ID  | \*   |
| title   | `string` | 通知标题 | \*   |
| content | `string` | 通知内容 | \*   |

</p>

</details>

<details><summary>删除个人消息<code>[DELETE] /v1/message/m/:message_id</code></summary>

<p>

删除个人消息

</p>

</details>

<details><summary>更改个人消息<code>[PUT] /v1/message/m/:message_id</code></summary>

<p>

| 参数    | 类型     | 说明     | 必填 |
| ------- | -------- | -------- | ---- |
| title   | `string` | 消息标题 |      |
| content | `string` | 消息内容 |      |

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

| 参数         | 类型     | 说明                                                   | 必填 |
| ------------ | -------- | ------------------------------------------------------ | ---- |
| image        | `string` | 图片 URL                                               | \*   |
| href         | `string` | 图片跳转的链接                                         | \*   |
| platform     | `string` | 该 banner 图片运用在哪个平台. 分别为 `PC` 或 `APP`     | \*   |
| description  | `string` | 该 banner 的描述信息                                   |      |
| priority     | `int`    | 优先级，用于排序                                       |      |
| identifier   | `string` | APP 跳转标识符, 给 APP 跳转页面用的                    |      |
| fallback_url | `string` | 当 APP 的 identifier 无效时的备选方案，跳转的 URL 地址 |      |

</p>

</details>

<details><summary>修改 banner<code>[PUT] /v1/banner/b/:banner_id</code></summary>

<p>

| 参数         | 类型     | 说明                                                   | 必填 |
| ------------ | -------- | ------------------------------------------------------ | ---- |
| image        | `string` | 图片 URL                                               |      |
| href         | `string` | 图片跳转的链接                                         |      |
| platform     | `string` | 该 banner 图片运用在哪个平台. 分别为 `PC` 或 `APP`     |      |
| description  | `string` | 该 banner 的描述信息                                   |      |
| priority     | `int`    | 优先级，用于排序                                       |      |
| identifier   | `string` | APP 跳转标识符, 给 APP 跳转页面用的                    |      |
| fallback_url | `string` | 当 APP 的 identifier 无效时的备选方案，跳转的 URL 地址 |      |

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

| Query 参数 | 类型     | 说明                          | 必选 |
| ---------- | -------- | ----------------------------- | ---- |
| platform   | `string` | 根据平台筛选, 可选 `pc`/`app` |      |
| active     | `bool`   | 是否可用                      |      |

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

### 用户反馈

<details><summary>获取反馈列表<code>[GET] /v1/report</code></summary>

<p>

获取反馈列表

| 参数   | 类型     | 说明                                   | 必选 |
| ------ | -------- | -------------------------------------- | ---- |
| uid    | `string` | 根据某个用户的`uid`筛选                |      |
| type   | `string` | 根据`类型`筛选                         |      |
| status | `int`    | 根据`状态`筛选, `0` 未解决, `1` 已解决 |      |

</p>

</details>

<details><summary>获取反馈详情<code>[GET] /v1/report/r/:report_id</code></summary>

<p>

获取一条反馈的详情

</p>

</details>

<details><summary>更新反馈信息<code>[PUT] /v1/report/r/:report_id</code></summary>

<p>

更新反馈信息, 主要标记是否已解决

| 参数   | 类型   | 说明                                   | 必选 |
| ------ | ------ | -------------------------------------- | ---- |
| status | `int`  | 根据`状态`筛选, `0` 未解决, `1` 已解决 |      |
| lock   | `bool` | 是否锁定该反馈, 锁定之后用户无法再更新 |      |

</p>

</details>

### 菜单模块

<details><summary>创建页面菜单<code>[POST] /v1/menu/m/:menu_id</code></summary>

<p>

创建页面菜单

| 参数      | 类型       | 说明                                        | 必选 |
| --------- | ---------- | ------------------------------------------- | ---- |
| name      | `string`   | 菜单名                                      | \*   |
| url       | `string`   | 菜单对应的页面 url                          |      |
| icon      | `string`   | 菜单图标                                    |      |
| accession | `[]string` | 菜单对应的权限                              |      |
| sort      | `int`      | 菜单排序, 值越大，菜单越靠前                |      |
| parent_id | `string`   | 该菜单的父级菜单, 如果 不填写，则为顶级菜单 |      |

</p>

</details>

<details><summary>获取页面菜单列表<code>[GET] /v1/menu</code></summary>

<p>

获取页面菜单列表

</p>

</details>

<details><summary>获取页面菜单详情<code>[GET] /v1/menu/m/:menu_id</code></summary>

<p>

获取一条页面菜单的详情

</p>

</details>

<details><summary>更新页面菜单<code>[PUT] /v1/menu/m/:menu_id</code></summary>

<p>

更新页面菜单

| 参数      | 类型       | 说明                         | 必选 |
| --------- | ---------- | ---------------------------- | ---- |
| name      | `string`   | 菜单名                       |      |
| url       | `string`   | 菜单对应的页面 url           |      |
| icon      | `string`   | 菜单图标                     |      |
| accession | `[]string` | 菜单对应的权限               |      |
| sort      | `int`      | 菜单排序, 值越大，菜单越靠前 |      |
| parent_id | `string`   | 该菜单的父级菜单             |      |

</p>

</details>

<details><summary>删除页面菜单<code>[DELETE] /v1/menu/m/:menu_id</code></summary>

<p>

删除页面菜单

</p>

</details>
