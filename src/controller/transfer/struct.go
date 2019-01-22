package transfer

type Log struct {
	Id        string  `json:"id"`
	From      string  `json:"from"`
	To        string  `json:"to"`
	Currency  string  `json:"currency"`
	Amount    float64 `json:"amount"`
	Status    int     `json:"status"`
	Note      *string `json:"note"`
	CreatedAt string  `json:"created_at"`
	UpdatedAt string  `json:"updated_at"`
}
