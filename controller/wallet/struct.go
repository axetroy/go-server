package wallet

type Wallet struct {
	Pure
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

type Pure struct {
	Id      string  `json:"id"`      // 用户ID
	Balance float64 `json:"balance"` // 可用余额
	Frozen  float64 `json:"frozen"`  // 冻结余额
}
