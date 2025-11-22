package wallet

import (
	"github.com/jmonteiro/picpay-like/core/types"
)

type WalletService struct {
	store types.WalletStore
}

func NewWalletService(store types.WalletStore) *WalletService {
	return &WalletService{
		store: store,
	}
}

func (s *WalletService) CreateWallet(wallet types.Wallet) error {

	return s.store.CreateWallet(wallet)
}

func (s *WalletService) GetWalletByID(id int) (*types.Wallet, error) {
	return s.store.GetWalletByID(id)
}

func (s *WalletService) GetWalletsByUserID(userID int) ([]*types.Wallet, error) {
	return s.store.GetWalletsByUserID(userID)
}
