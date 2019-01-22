package exception

var (
	// upload
	NotSupportType = New("不支持该文件类型")
	OutOfSize      = New("超出文件大小限制")
)
