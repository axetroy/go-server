### 新增 banner

[POST] /v1/banner

| 参数         | 类型     | 说明                                                   | 必填 |
| ------------ | -------- | ------------------------------------------------------ | ---- |
| image        | `string` | 图片 URL                                               | \*   |
| href         | `string` | 图片跳转的链接                                         | \*   |
| platform     | `string` | 该 banner 图片运用在哪个平台. 分别为 `PC` 或 `APP`     | \*   |
| description  | `string` | 该 banner 的描述信息                                   |      |
| priority     | `int`    | 优先级，用于排序                                       |      |
| identifier   | `string` | APP 跳转标识符, 给 APP 跳转页面用的                    |      |
| fallback_url | `string` | 当 APP 的 identifier 无效时的备选方案，跳转的 URL 地址 |      |

### 修改 banner

[PUT] /v1/banner/:banner_id

| 参数         | 类型     | 说明                                                   | 必填 |
| ------------ | -------- | ------------------------------------------------------ | ---- |
| image        | `string` | 图片 URL                                               |      |
| href         | `string` | 图片跳转的链接                                         |      |
| platform     | `string` | 该 banner 图片运用在哪个平台. 分别为 `PC` 或 `APP`     |      |
| description  | `string` | 该 banner 的描述信息                                   |      |
| priority     | `int`    | 优先级，用于排序                                       |      |
| identifier   | `string` | APP 跳转标识符, 给 APP 跳转页面用的                    |      |
| fallback_url | `string` | 当 APP 的 identifier 无效时的备选方案，跳转的 URL 地址 |      |

### 删除 banner

[DELETE] /v1/banner/:banner_id

### 获取 banner 列表

[GET] /v1/banner

| Query 参数 | 类型     | 说明                          | 必选 |
| ---------- | -------- | ----------------------------- | ---- |
| platform   | `string` | 根据平台筛选, 可选 `pc`/`app` |      |
| active     | `bool`   | 是否可用                      |      |

### 获取 banner 详情

[GET] /v1/banner/:banner_id