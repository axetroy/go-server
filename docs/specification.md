### 1.通用说明

- 接口的输入和输出，使用 UTF-8 编码
- 接口返回以 json 格式输出
- 接口在正常状态下，HTTP 返回码均为 200（不同于业务返回码 code）
- 所有接口采用 https 协议

### 2. 请求说明

- 用户处于已登录状态时，所有请求需要带入 token（通过头部 `Authorization` 字段）
- 认证方式采用 `Web JSON Token`, 根据规范，携带的 Header 应该是 `Authorization: Bearer 你的身份令牌`
- 请求数据格式为： `application/json`
- 需要使用输入交易密码时，请求头部需要加上 `X-Pay-Password` 字段指定交易密码

**通用的查询参数**:

- limit:
    - 查询多少条数据，默认 `10`
- page:
    - 查询第 n 页数据，默认 `0`
- sort:
    - 按字段排序，默认 `created_at DESC`, 按照创建时间排序
- platform:
    - 指定平台(预留字段)

### 3. 返回说明

返回数据的统一结构如下：

```json
{
  "message": "",
  "data": [],
  "status": 1,
  "meta": {
    "limit": 10,
    "page": 0,
    "total": 100,
    "num": 10,
    "sort": "",
    "platform": ""
  }
}
```

- message:
    - 接口返回的信息
    
- data:
    - 接口返回的数据
    - 这个数据可以是任意类型，`null`/`string`/`bool`/`object`...等
    
- status:
    - 状态码，非 1 状态码则为错误
    
- meta:
    - 当查询列表时，才会返回这个字段
    - limit:
        - 查询多少条数据
        - 默认 10
    - page:
        - 查询第 n 页数据，从 0 开始
        - 默认 0
    - total:
        - 存在的数据总量
    - num:
        - 当前返回的数据条数
    - sort:
        - 当前筛选的排序条件
    - platform:
        - 通过指定平台查询(预留字段)