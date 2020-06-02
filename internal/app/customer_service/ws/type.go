package ws

type TypeForWriteMessage string         // 所有的类型总称
type TypeFromUser TypeForWriteMessage   // 用户发出的消息类型
type TypeToUser TypeForWriteMessage     // 输出给用户的消息类型
type TypeFromWaiter TypeForWriteMessage // 客服发出的消息类型
type TypeToWaiter TypeForWriteMessage   // 输出给客服的消息类型

// 用户可以发出的消息类型
const (
	TypeFromUserConnect     TypeFromUser = "connect"      // 请求连接一个客服
	TypeFromUserDisconnect  TypeFromUser = "disconnect"   // 请求和客服断开连接
	TypeFromUserMessageText TypeFromUser = "message_text" // 发送文本消息
)

// 用户收到的类型
const (
	TypeToUserInitialize     TypeToUser = "initialize"      // 初始化，告诉用户当前的链接 ID
	TypeToUserNotConnect     TypeToUser = "not_connect"     // 尚未连接
	TypeToUserConnectSuccess TypeToUser = "connect_success" // 连接成功，现在可以开始对话
	TypeToUserDisconnected   TypeToUser = "disconnected"    // 与客服断开连接
	TypeToUserConnectQueue   TypeToUser = "connect_queue"   // 正在排队，请等待
	TypeToUserMessageText    TypeToUser = "message_text"    // 用户收到文本消息
	TypeToUserError          TypeToUser = "error"           // 用户收到一个错误
)

// 客服发出的消息类型
const (
	TypeFromWaiterReady       TypeFromWaiter = "ready"        // 客服已准备就绪，可以开始接收客人
	TypeFromWaiterMessageText TypeFromWaiter = "message_text" // 客服发出文本消息
)

// 客服收到的消息
const (
	TypeToWaiterTypeInitializeToUser TypeToWaiter = "initialize"     // 初始化，告诉客服当前的链接 ID
	TypeToWaiterMessageText          TypeToWaiter = "message_text"   // 客服收到文本消息
	TypeToWaiterNewConnection        TypeToWaiter = "new_connection" // 有新连接
	TypeToWaiterDisconnected         TypeToWaiter = "disconnected"   // 有新连接断开
	TypeToWaiterError                TypeToWaiter = "error"          // 有新连接断开
)
