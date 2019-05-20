package schema

type WalletPure struct {
	Id       string `json:"id"`       // 用户ID
	Currency string `json:"currency"` // 币种
	Balance  string `json:"balance"`  // 可用余额
	Frozen   string `json:"frozen"`   // 冻结余额
}

type Wallet struct {
	WalletPure
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}
