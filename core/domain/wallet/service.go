package wallet

import (
	"fmt"

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
	_, err := s.store.GetWalletByUserID(wallet.UserID)
	if err == nil {
		return fmt.Errorf("wallet for user with ID %d already exists", wallet.UserID)
	}
	return s.store.CreateWallet(wallet)
}

func (s *WalletService) GetWalletByID(id int) (*types.Wallet, error) {
	return s.store.GetWalletByID(id)
}

func (s *WalletService) GetWalletByUserID(userID int) (*types.Wallet, error) {
	return s.store.GetWalletByUserID(userID)
}