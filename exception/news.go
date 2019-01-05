package exception

var (
	NewsInvalidType = New("错误的文章类型")
	NewsNotExist    = New("文章不存在")
)
