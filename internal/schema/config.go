package schema

type Config struct {
	Fields    interface{} `json:"fields"`
	CreatedAt string      `json:"created_at"`
	UpdatedAt string      `json:"updated_at"`
}

type Name struct {
	Field       string `json:"field"`       // 字段名称
	Description string `json:"description"` // 配置描述
}
