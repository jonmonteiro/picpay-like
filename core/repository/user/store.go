package repository

import (
	"database/sql"
	"fmt"

	"github.com/jmonteiro/picpay-like/core/domain/user"
)

type Store struct {
	db *sql.DB
}

func NewStore(db *sql.DB) domain.UserStore {
	return &Store{db: db}
}

func (s *Store) CreateUser(user domain.User) error {
	_, err := s.db.Exec("INSERT INTO users (email, password) VALUES ($1, $2)", user.Email, user.Password)
	if err != nil {
		return err
	}

	return nil
}

func (s *Store) GetUserByID(id int) (*domain.User, error) {
	row := s.db.QueryRow(`SELECT id, email, password FROM users WHERE id=$1`, id)

	var u domain.User
	err := row.Scan(&u.ID, &u.Email, &u.Password)
	if err != nil {
		return nil, err
	}

	return &u, nil
}

func (s *Store) GetUserByEmail(email string) (*domain.User, error) {
	rows, err := s.db.Query("SELECT * FROM users WHERE email = $1", email)
	if err != nil {
		return nil, err
	}

	u := new(domain.User)
	for rows.Next() {
		u, err = scanRowsIntoUser(rows)
		if err != nil {
			return nil, err
		}
	}

	if u.ID == 0 {
		return nil, fmt.Errorf("user not found")
	}

	return u, nil
}

func (s *Store) GetUsers() ([]*domain.User, error) {
	rows, err := s.db.Query(`SELECT id, email, password FROM users`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []*domain.User
	for rows.Next() {
		var u domain.User
		if err := rows.Scan(&u.ID, &u.Email, &u.Password); err != nil {
			return nil, err
		}
		users = append(users, &u)
	}

	return users, rows.Err()
}

func (s *Store) UpdateUser(u domain.RegisterUserPayload, id int) error {
	_, err := s.db.Exec(`UPDATE users SET email=$1, password=$2 WHERE id=$3`, u.Email, u.Password, id)
	return err
}

func (s *Store) DeleteUser(id int) error {
	_, err := s.db.Exec(`DELETE FROM users WHERE id=$1`, id)
	return err
}

func scanRowsIntoUser(rows *sql.Rows) (*domain.User, error) {
	user := new(domain.User)

	err := rows.Scan(
		&user.ID,
		&user.Email,
		&user.Password,
	)
	if err != nil {
		return nil, err
	}

	return user, nil
}