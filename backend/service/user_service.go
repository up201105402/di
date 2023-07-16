package service

import (
	"di/model"
	"di/repository"
	"di/util"
	"errors"

	"github.com/nicksnyder/go-i18n/v2/i18n"
	"gorm.io/gorm"
)

type userServiceImpl struct {
	I18n           *i18n.Localizer
	UserRepository repository.UserRepository
}

func NewUserService(gormDB *gorm.DB, i18n *i18n.Localizer) UserService {
	return &userServiceImpl{
		I18n:           i18n,
		UserRepository: repository.NewUserRepository(gormDB),
	}
}

func (service *userServiceImpl) Get(id uint) (*model.User, error) {
	user, err := service.UserRepository.FindByID(id)

	if err != nil {
		errMessage, _ := service.I18n.Localize(&i18n.LocalizeConfig{
			MessageID: "user.repository.find.user.id.failed",
			TemplateData: map[string]interface{}{
				"ID":     id,
				"Reason": err.Error(),
			},
		})

		return user, errors.New(errMessage)
	}

	return user, err
}

func (service *userServiceImpl) GetByUsername(username string) (*model.User, error) {
	user, err := service.UserRepository.FindByUsername(username)

	if err != nil {
		errMessage, _ := service.I18n.Localize(&i18n.LocalizeConfig{
			MessageID: "user.repository.find.user.username.failed",
			TemplateData: map[string]interface{}{
				"Username": username,
				"Reason":   err.Error(),
			},
		})

		return user, errors.New(errMessage)
	}

	return user, err
}

func (service *userServiceImpl) Signup(username string, password string) error {
	pw, err := util.HashPassword(password)

	if err != nil {
		errMessage, _ := service.I18n.Localize(&i18n.LocalizeConfig{
			MessageID: "user.repository.create.user.failed",
			TemplateData: map[string]interface{}{
				"Reason": err.Error(),
			},
		})

		return errors.New(errMessage)
	}

	if err := service.UserRepository.Create(&model.User{Username: username, Password: pw}); err != nil {
		errMessage, _ := service.I18n.Localize(&i18n.LocalizeConfig{
			MessageID: "user.repository.create.user.failed",
			TemplateData: map[string]interface{}{
				"Reason": err.Error(),
			},
		})

		return errors.New(errMessage)
	}

	return nil
}

func (service *userServiceImpl) Signin(user *model.User) error {
	uFetched, err := service.UserRepository.FindByUsername(user.Username)

	if err != nil {
		errMessage, _ := service.I18n.Localize(&i18n.LocalizeConfig{
			MessageID: "user.repository.find.user.username.failed",
			TemplateData: map[string]interface{}{
				"Username": user.Username,
				"Reason":   err.Error(),
			},
		})

		return errors.New(errMessage)
	}

	match, err := util.ComparePasswords(uFetched.Password, user.Password)

	if err != nil {
		errMessage, _ := service.I18n.Localize(&i18n.LocalizeConfig{
			MessageID: "user.service.user.signin.failed",
			TemplateData: map[string]interface{}{
				"Reason": err.Error(),
			},
		})

		return errors.New(errMessage)
	}

	if !match {
		errMessage, _ := service.I18n.Localize(&i18n.LocalizeConfig{
			MessageID: "user.service.user.signin.password.failed",
		})

		return errors.New(errMessage)
	}

	*user = *uFetched
	return nil
}

func (service *userServiceImpl) UpdateDetails(user *model.User) error {
	err := service.UserRepository.Update(user)

	if err != nil {
		errMessage, _ := service.I18n.Localize(&i18n.LocalizeConfig{
			MessageID: "user.repository.update.user.failed",
			TemplateData: map[string]interface{}{
				"Reason": err.Error(),
			},
		})

		return errors.New(errMessage)
	}

	return nil
}
