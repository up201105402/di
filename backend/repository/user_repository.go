package repository

import (
	"di/model"
	"errors"
	"log"

	"gorm.io/gorm"
)

type userRepository struct {
	DB *gorm.DB
}

// NewUserRepository is a User Repository factory
func NewUserRepository(gormDB *gorm.DB) model.UserRepository {
	return &userRepository{
		DB: gormDB,
	}
}

// Creates a new user
func (repo *userRepository) Create(user *model.User) error {
	result := repo.DB.Create(user)

	if result.Error != nil {
		log.Printf("Failed to create user. Reason: %v\n", result.Error)
		return result.Error
	}

	return nil
}

// Updates a user's properties
func (repo *userRepository) Update(user *model.User) error {
	result := repo.DB.Save(user)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		repo.Create(user)
	}

	return nil
}

// FindByID fetches a user by id
func (repo *userRepository) FindByID(id uint) (*model.User, error) {

	var user = model.User{}

	result := repo.DB.First(&user, id)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		log.Printf("Failed to get user with id: %v. Reason: %v\n", id, result.Error)
		return nil, result.Error
	}

	return &user, nil
}

// FindByUsername fetches a user by username
func (repo *userRepository) FindByUsername(username string) (*model.User, error) {
	var user = model.User{Username: username}

	result := repo.DB.First(&user, "username = ?", username)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		log.Printf("Failed to get user with username: %v. Reason: %v\n", username, result.Error)
		return nil, result.Error
	}

	return &user, nil
}
