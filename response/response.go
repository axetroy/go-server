package response

type Meta struct {
	Limit    int     `json:"limit"`
	Page     int     `json:"page"`
	Total    int64   `json:"total"`
	Num      int     `json:"num"`
	Sort     string  `json:"sort"`
	Platform *string `json:"platform"`
}

const (
	StatusSuccess = 1
	StatusFail    = 0

	DefaultLimit = 10
	DefaultPage  = 1

	MaxLimit = 100
)

type Response struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
	Status  int         `json:"status"`
}

type List struct {
	Response
	Meta *Meta `json:"meta"`
}
