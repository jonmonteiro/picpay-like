package types

type WalletStore interface {
	CreateWallet(wallet Wallet) error
	GetWalletByID(id int) (*Wallet, error)
	GetWalletsByUserID(userID int) ([]*Wallet, error)
}

type Wallet struct {
	ID         int     `json:"id"`
	UserID     int     `json:"user_id"`
	CardNumber string  `json:"card_number"`
	Balance    float64 `json:"balance"`
}
