package password

type Level string

const (
	LevelLow    Level = "low"
	LevelMiddle Level = "middle"
	LevelHigh   Level = "high"
)

type Password struct {
	Text  string
	Level Level // 安全级别
}

func New(text string) *Password {

	// 纯数字/纯字母 -> 低级别
	// 数字 + 字母 -> 中级别
	// 数字 + 字母 + 符号 -> 高级别

	return &Password{
		Text:  text,
		Level: LevelLow,
	}
}

func (c *Password) isValid() bool {
	if len(c.Text) < 6 {
		return false
	}

	return true
}
