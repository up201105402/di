package service

import (
	"di/model"
	"di/repository"
	"di/util"
	"log"

	"gorm.io/gorm"
)

type userServiceImpl struct {
	UserRepository repository.UserRepository
}

func NewUserService(gormDB *gorm.DB) UserService {
	return &userServiceImpl{
		UserRepository: repository.NewUserRepository(gormDB),
	}
}

func (service *userServiceImpl) Get(id uint) (*model.User, error) {
	user, error := service.UserRepository.FindByID(id)
	return user, error
}

func (service *userServiceImpl) GetByUsername(username string) (*model.User, error) {
	user, error := service.UserRepository.FindByUsername(username)
	return user, error
}

func (service *userServiceImpl) Signup(username string, password string) error {
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

func (service *userServiceImpl) Signin(user *model.User) error {
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

func (service *userServiceImpl) UpdateDetails(user *model.User) error {
	// Update user in UserRepository
	err := service.UserRepository.Update(user)

	if err != nil {
		return err
	}

	return nil
}
