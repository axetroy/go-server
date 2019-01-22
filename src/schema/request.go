package schema

var (
	DefaultLimit = 10
	DefaultPage  = 0
	MaxLimit     = 100
)

type Query struct {
	Limit    int     `json:"limit" form:"limit"`
	Page     int     `json:"page" form:"page"`
	Sort     string  `json:"sort" form:"sort"`
	Platform *string `json:"platform" form:"platform"`
}

func (q *Query) Normalize() *Query {

	if q.Limit <= 0 {
		q.Limit = DefaultLimit // 默认查询10条
	} else if q.Limit > MaxLimit {
		q.Limit = MaxLimit
	}

	if q.Page <= 0 {
		q.Page = DefaultPage
	}

	return q
}
