// Copyright 2019-2020 Axetroy. All rights reserved. MIT license.
package ws

type typeForWriteMessage string             // 所有的类型总称
type TypeRequestUser typeForWriteMessage    // 用户发出的消息类型
type TypeResponseUser typeForWriteMessage   // 输出给用户的消息类型
type TypeRequestWaiter typeForWriteMessage  // 客服发出的消息类型
type TypeResponseWaiter typeForWriteMessage // 输出给客服的消息类型

// 用户可以发出的消息类型
const (
	TypeRequestUserAuth        TypeRequestUser = "auth"         // 认证帐号
	TypeRequestUserConnect     TypeRequestUser = "connect"      // 请求连接一个客服
	TypeRequestUserDisconnect  TypeRequestUser = "disconnect"   // 请求和客服断开连接
	TypeRequestUserMessageText TypeRequestUser = "message_text" // 发送文本消息
	TypeRequestUserGetHistory  TypeRequestUser = "get_history"  // 请求获取用户聊天记录，应该返回 `message_history`
)

// 用户收到的类型
const (
	TypeResponseUserAuthSuccess        TypeResponseUser = "auth_success"         // 初始化，告诉用户当前的链接 ID
	TypeResponseUserNotConnect         TypeResponseUser = "not_connect"          // 尚未连接
	TypeResponseUserConnectSuccess     TypeResponseUser = "connect_success"      // 连接成功，现在可以开始对话
	TypeResponseUserMessageHistory     TypeResponseUser = "message_history"      // 用户的聊天记录
	TypeResponseUserDisconnected       TypeResponseUser = "disconnected"         // 客服与用户断开连接
	TypeResponseUserConnectQueue       TypeResponseUser = "connect_queue"        // 正在排队，请等待
	TypeResponseUserMessageText        TypeResponseUser = "message_text"         // 用户收到文本消息
	TypeResponseUserMessageTextSuccess TypeResponseUser = "message_text_success" // message_text 的回执
	TypeResponseUserIdle               TypeResponseUser = "idle"                 // 当前连接出于空闲状态
	TypeResponseUserError              TypeResponseUser = "error"                // 用户收到一个错误
)

// 客服发出的消息类型
const (
	TypeRequestWaiterAuth         TypeRequestWaiter = "auth"         // 身份认证
	TypeResponseWaiterAuthSuccess TypeRequestWaiter = "auth_success" // 身份认证
	TypeRequestWaiterReady        TypeRequestWaiter = "ready"        // 客服已准备就绪，可以开始接收客人
	TypeRequestWaiterUnReady      TypeRequestWaiter = "unready"      // 客服进入未就绪状态，意味着客服暂停接客
	TypeRequestWaiterMessageText  TypeRequestWaiter = "message_text" // 客服发出文本消息
	TypeRequestWaiterDisconnect   TypeRequestWaiter = "disconnect"   // 请求断开连接
	TypeRequestWaiterGetHistory   TypeRequestWaiter = "get_history"  // 请求获取用户聊天记录，应该返回 `message_history`, 需要指定 payload
)

// 客服收到的消息
const (
	TypeResponseWaiterMessageText        TypeResponseWaiter = "message_text"         // 客服收到文本消息
	TypeResponseWaiterMessageTextSuccess TypeResponseWaiter = "message_text_success" // message_text 成功的回执
	TypeResponseWaiterNewConnection      TypeResponseWaiter = "new_connection"       // 有新连接
	TypeResponseWaiterDisconnected       TypeResponseWaiter = "disconnected"         // 有新连接断开
	TypeResponseWaiterKickOut            TypeResponseWaiter = "kickout"              // 被踢下线
	TypeResponseWaiterMessageHistory     TypeResponseWaiter = "message_history"      // 用户的聊天记录
	TypeResponseWaiterError              TypeResponseWaiter = "error"                // 有新连接断开
)
