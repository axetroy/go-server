package transfer

type Log struct {
	Id        int64   `json:"id"`
	Currency  string  `json:"currency"`
	From      int64   `json:"from"`
	To        int64   `json:"to"`
	Amount    float64 `json:"amount"`
	Status    int     `json:"status"`
	Note      *string `json:"note"`
	CreatedAt string  `json:"created_at"`
	UpdatedAt string  `json:"updated_at"`
}
