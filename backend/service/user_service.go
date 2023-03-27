package service

import (
	"di/model"
	"di/repository"
	"di/util"
	"log"

	"gorm.io/gorm"
)

type userService struct {
	UserRepository model.UserRepository
}

func NewUserService(gormDB *gorm.DB) model.UserService {
	return &userService{
		UserRepository: repository.NewUserRepository(gormDB),
	}
}

func (service *userService) Get(id uint) (*model.User, error) {
	user, error := service.UserRepository.FindByID(id)
	return user, error
}

func (service *userService) Signup(username string, password string) error {
	pw, err := util.HashPassword(password)

	if err != nil {
		log.Printf("Unable to signup user for email: %v\n", username)
		return err
	}

	if err := service.UserRepository.Create(&model.User{Username: username, Password: pw}); err != nil {
		return err
	}

	return nil
}

func (service *userService) Signin(user *model.User) error {
	uFetched, err := service.UserRepository.FindByUsername(user.Username)

	if err != nil {
		return err
	}

	match, err := util.ComparePasswords(uFetched.Password, user.Password)

	if err != nil {
		return err
	}

	if !match {
		return err
	}

	*user = *uFetched
	return nil
}

func (service *userService) UpdateDetails(user *model.User) error {
	// Update user in UserRepository
	err := service.UserRepository.Update(user)

	if err != nil {
		return err
	}

	return nil
}
