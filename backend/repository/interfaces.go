package repository

import "di/model"

type UserRepository interface {
	FindByID(id uint) (*model.User, error)
	FindByUsername(username string) (*model.User, error)
	Create(u *model.User) error
	Update(u *model.User) error
}
