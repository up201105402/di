package model

type UserService interface {
	Get(id uint) (*User, error)
	Signup(username string, password string) error
	Signin(user *User) error
	UpdateDetails(user *User) error
}

type UserRepository interface {
	FindByID(id uint) (*User, error)
	FindByUsername(username string) (*User, error)
	Create(u *User) error
	Update(u *User) error
}
