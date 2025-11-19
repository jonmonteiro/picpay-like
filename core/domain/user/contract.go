package domain

type UserStore interface {
	CreateUser(User) error
	GetUserByID(id int) (*User, error)
	GetUserByEmail(email string) (*User, error)
    GetUsers() ([]*User, error)
	UpdateUser(user RegisterUserPayload, id int) error	
	DeleteUser(id int) error
}
