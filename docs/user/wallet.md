### 获取我的钱包

[GET] /v1/wallet

获取我的钱包 Map, 以 Key-Value 的 形式返回.

### 获取钱包详情

/v1/wallet/w/:currency

获取指定一个钱包的详细信息.

### 钱包转账

[POST] /v1/transfer

需要在请求头设置 `X-Pay-Password`, 指定二级密码.

需要在请求头设置 `X-Signature`, 指定数据的签名.

| 参数     | 类型     | 说明                    | 必选 |
| -------- | -------- | ----------------------- | ---- |
| currency | `string` | 钱包类型                | \*   |
| to       | `string` | 转账对象的用户纯数字 ID | \*   |
| amount   | `string` | 转账金额                | \*   |
| note     | `string` | 转账备注                |      |

!> 在发起转账前，先调用签名接口，把 JSON 格式的参数，提交到 `/v1/signature` 进行签名. 签名后赋值给 `X-Signature`

### 获取转账记录

[GET] /v1/transfer

获取我的转账记录

### 获取转账记录详情

[GET] /v1/transfer/t/:transfer_id

获取某一条转账记录的详情
