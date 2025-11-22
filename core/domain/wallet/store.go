package wallet

import (
	"database/sql"
	"github.com/jmonteiro/picpay-like/core/types"
)

type Store struct {
	db *sql.DB
}

func NewStore(db *sql.DB) types.WalletStore {
	return &Store{db: db}
}

func (s *Store) CreateWallet(wallet types.Wallet) error {
	_, err := s.db.Exec("INSERT INTO wallets (user_id, balance) VALUES ($1, $2)", wallet.UserID, wallet.Balance)
	if err != nil {
		return err
	}

	return nil
}

func (s *Store) GetWalletByID(id int) (*types.Wallet, error) {
	row := s.db.QueryRow(`SELECT id, user_id, balance FROM wallets WHERE id=$1`, id)

	var w types.Wallet
	err := row.Scan(&w.ID, &w.UserID, &w.Balance)
	if err != nil {
		return nil, err
	}

	return &w, nil
}

func (s *Store) GetWalletByUserID(userID int) (*types.Wallet, error) {
	row := s.db.QueryRow(`SELECT id, user_id, balance FROM wallets WHERE user_id=$1`, userID)

	var w types.Wallet
	err := row.Scan(&w.ID, &w.UserID, &w.Balance)
	if err != nil {
		return nil, err
	}

	return &w, nil
}


