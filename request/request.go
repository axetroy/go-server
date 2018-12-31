package request

var (
	DefaultLimit = 10
	DefaultPage  = 0
	MaxLimit     = 100
)

type Query struct {
	Limit    int     `json:"limit"`
	Page     int     `json:"page"`
	Sort     string  `json:"sort"`
	Platform *string `json:"platform"`
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
