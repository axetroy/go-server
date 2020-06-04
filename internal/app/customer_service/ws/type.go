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
)

// 用户收到的类型
const (
	TypeResponseUserAuthSuccess    TypeResponseUser = "auth_success"    // 初始化，告诉用户当前的链接 ID
	TypeResponseUserNotConnect     TypeResponseUser = "not_connect"     // 尚未连接
	TypeResponseUserConnectSuccess TypeResponseUser = "connect_success" // 连接成功，现在可以开始对话
	TypeResponseUserDisconnected   TypeResponseUser = "disconnected"    // 客服与用户断开连接
	TypeResponseUserConnectQueue   TypeResponseUser = "connect_queue"   // 正在排队，请等待
	TypeResponseUserMessageText    TypeResponseUser = "message_text"    // 用户收到文本消息
	TypeResponseUserError          TypeResponseUser = "error"           // 用户收到一个错误
)

// 客服发出的消息类型
const (
	TypeRequestWaiterAuth        TypeRequestWaiter = "auth"         // 身份认证
	TypeRequestWaiterReady       TypeRequestWaiter = "ready"        // 客服已准备就绪，可以开始接收客人
	TypeRequestWaiterMessageText TypeRequestWaiter = "message_text" // 客服发出文本消息
	TypeRequestWaiterDisconnect  TypeRequestWaiter = "disconnect"   // 请求断开连接
)

// 客服收到的消息
const (
	TypeResponseWaiterMessageText   TypeResponseWaiter = "message_text"   // 客服收到文本消息
	TypeResponseWaiterNewConnection TypeResponseWaiter = "new_connection" // 有新连接
	TypeResponseWaiterDisconnected  TypeResponseWaiter = "disconnected"   // 有新连接断开
	TypeResponseWaiterError         TypeResponseWaiter = "error"          // 有新连接断开
)
