客服系统基于 WebSocket 实现

为了更好的开发，这里约定几个事实标准：

### 1. 所有的 WebSocket 帧都为 JSON

无论是客服端发出的，还是收到帧。并且都满足以下格式

客户端发出的消息，必须满足以下格式：

| 字段    | 类型     | 说明                                                             | 必填 |
| ------- | -------- | ---------------------------------------------------------------- | ---- |
| op_id   | `uuid`   | 这条消息的 `uuid`，**可选**，如果传参，则返回的 `op_id` 于它一致 |      |
| to      | `string` | 这条消息发给谁，**只有客服需要传这个字段**                       |      |
| type    | `string` | 消息类型，下面见详情                                             | \*   |
| payload | `any`    | 消息的的附带数据                                                 |      |

客户端接收到的消息满足以下格式：

| 字段    | 类型     | 说明                                                                     | 必填 |
| ------- | -------- | ------------------------------------------------------------------------ | ---- |
| op_id   | `uuid`   | 这条消息的 `uuid`，于传入的消息的 `op_id` 一致                           |      |
| from    | `string` | 这条消息来自哪里，**发送不需要这个字段，但是返回的消息可能带有这个字段** |      |
| to      | `string` | 这条消息发给谁，**只有客服需要传这个字段**                               | \*   |
| type    | `string` | 消息类型，下面见详情                                                     | \*   |
| payload | `any`    | 消息的的附带数据                                                         |      |

### 2. 连接流程

#### `用户端` 的连接流程 `ws://localhost/v1/ws/connect/user`

1 - 连接 WebSocket

2 - 发送数据体 `{"type": "auth", payload: { "token": "your token" }}` 给服务器进行身份认证

3 - 发送数据体 `{"type": "connect"}` 给服务器请求连接

4 - 收到来自服务器的推送 `{"from": "对方的 UUID", "to":"你的 UUID","type":"connect_success","payload":{"id":"对方的 id"}}` 表示和客服连接成功

5 - 发送数据 `{"type":"message_text","payload":{"text":"你好，这是一条信息!"}}` 用于发送信息给客服

6 - 客服主动与你断开连接，收到来自服务器的推送 `{"from":"对方的 UUID","to":"你的UUID","type":"disconnected"}`

#### `客服端` 的连接流程 `ws://localhost/v1/ws/connect/waiter`

1 - 连接 WebSocket

2 - 发送数据体 `{"type": "auth", payload: { "token": "your token" }}` 给服务器进行身份认证

3 - 发送数据体 `{"type": "ready"}` 给服务器表示客服已准备就绪，可以连接用户

4 - 收到来自服务器的推送 `{"from": "对方的 UUID", "to":"你的 UUID","type":"new_connection","payload":{"id":"对方的 id"}}` 表示和新的用户建立连接

5 - 发送数据 `{"to": "用户的 UUID","type":"message_text","payload":{"text":"你好，这是一条信息!"}}` 用于发送信息给用户

6 - 用户主动与你断开连接，收到来自服务器的推送 `{"from":"用户的 UUID","to":"你的 UUID","type":"disconnected"}`

### 消息类型 type

消息类型定义了 **你允许发送什么类型的消息** 和 **你会收到什么类型的消息**

#### 用户端可以发出的消息类型

| Type          | 说明                                     | 对应的 Payload                                        |
| ------------- | ---------------------------------------- | ----------------------------------------------------- |
| auth          | 身份认证                                 | `{"token": "xxxx"}`                                   |
| connect       | 请求连接一个客服                         | `null`                                                |
| disconnect    | 与客服断开连接                           | `null`                                                |
| message_text  | 发送消息文本给客服，**需要先连接到客服** | `{"text": "这是一条消息"}`                            |
| message_image | 发送消息文本给客服，**需要先连接到客服** | `{"image": "[https/](https://example.com/demo.png)"}` |
| get_history   | 获取聊天记录                             | 返回 `message_history`                                |

#### 用户端会收到的消息类型

| Type                  | 说明                               | 对应的 Payload                                                                  |
| --------------------- | ---------------------------------- | ------------------------------------------------------------------------------- |
| auth_success          | 身份认证成功                       | `{"id":"274588402135859200","username":"test1","nickname":"test1","avatar":""}` |
| not_connect           | 尚未与客服连接                     | `...`                                                                           |
| connect_success       | 连接客服成功                       | `{"uuid": "客服的 UUID"}`                                                       |
| disconnected          | 连接已断开                         | `null`                                                                          |
| connect_queue         | 请求连接客服，但是正忙,正在排队    | `{"location": 100}`                                                             |
| message_text          | 收到来自客服的文本消息             | `{"text": "这是一条消息"}`                                                      |
| message_text_success  | `message_text` 的成功回执          | `{"text": "这是一条消息"}`                                                      |
| message_image         | 收到来自客服的图片消息             | `{"image": "https://example.com/demo.png"}`                                     |
| message_image_success | `message_image` 的成功回执         | `{"image": "https://example.com/demo.png"}`                                     |
| message_history       | 系统推送过来的聊天记录             | `[...]`                                                                         |
| idle                  | 连接空闲，则在接下来的时间断开连接 | `{"message": "xxxx" }`                                                          |
| error                 | 操作错误                           | `{"message": "这是错误信息"}`                                                   |

#### 客服端可以发出的消息类型

| Type                | 说明                                          | 对应的 Payload              |
| ------------------- | --------------------------------------------- | --------------------------- |
| auth                | 身份认证                                      | `{"token": "xxxx"}`         |
| ready               | 客服已就绪，可以连接客户                      | `null`                      |
| unready             | 客服暂停接口，不会再接收新的用户分配          | `null`                      |
| disconnect          | 与指定的用户断开连接                          | `{"uuid": "xxx"}`           |
| message_text        | 发送消息文本给用户, **需要指定 to 字段**      | `{"text": "这是一条消息"}`  |
| message_image       | 发送消息文本给用户, **需要指定 to 字段**      | `{"image": "这是一条消息"}` |
| get_history         | 获取聊天记录, payload 需要指定 `user_id` 字段 | 返回 `message_history`      |
| get_history_session | 获取会话记录                                  | `null`                      |

#### 客服端会收到的消息类型

| Type                  | 说明                                 | 对应的 Payload                                                                  |
| --------------------- | ------------------------------------ | ------------------------------------------------------------------------------- |
| auth_success          | 身份认证成功                         | `{"id":"274588402135859200","username":"test1","nickname":"test1","avatar":""}` |
| new_connection        | 有新的用户与客服连接                 | `{"uuid": "用户的 UUID"}`                                                       |
| disconnected          | 用户主动与客服断开                   | `null`                                                                          |
| unready               | 客服暂停接口，不会再接收新的用户分配 | `null`                                                                          |
| kickout               | 被踢下线                             | `null`                                                                          |
| message_text          | 收到来自用户的文本消息               | `{"text": "这是一条消息"}`                                                      |
| message_text_success  | `message_text` 的成功回执            | `{"text": "这是一条消息"}`                                                      |
| message_image         | 收到来自用户的图片消息               | `{"image": "https://example.com/demo.png"}`                                     |
| message_image_success | `message_image` 的成功回执           | `{"image": "https://example.com/demo.png"}`                                     |
| message_history       | 系统推送过来的聊天记录               | `[...]`                                                                         |
| session_history       | 系统推送过来的客服会话记录           | `[...]`                                                                         |
| error                 | 操作错误                             | `{"message": "这是错误信息"}`                                                   |

对应的 type 源码: [https://github.com/axetroy/go-server/blob/master/internal/app/customer_service/ws/type.go]
对应的 payload 源码: [https://github.com/axetroy/go-server/blob/master/internal/app/customer_service/ws/type_payload.go]
