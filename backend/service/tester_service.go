package service

import (
	"di/model"
	"di/repository"
	"errors"
	"log"
	"os"

	"github.com/hibiken/asynq"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"gorm.io/gorm"
)

type testerServiceImpl struct {
	TesterRepository model.TesterRepository
	TaskQueueClient  *asynq.Client
	I18n             *i18n.Localizer
}

func NewTesterService(gormDB *gorm.DB, client *asynq.Client, i18n *i18n.Localizer) TesterService {
	return &testerServiceImpl{
		TesterRepository: repository.NewTesterRepository(gormDB),
		TaskQueueClient:  client,
		I18n:             i18n,
	}
}

func (service *testerServiceImpl) Get(id uint) (*model.Tester, error) {
	tester, err := service.TesterRepository.FindByID(id)

	if err != nil {
		errMessage := service.I18n.MustLocalize(&i18n.LocalizeConfig{
			MessageID: "tester.repository.find.tester.id.failed",
			TemplateData: map[string]interface{}{
				"ID":     id,
				"Reason": err.Error(),
			},
			PluralCount: 1,
		})

		return tester, errors.New(errMessage)
	}

	return tester, err
}

func (service *testerServiceImpl) GetByOwner(ownerId uint) ([]model.Tester, error) {
	testers, err := service.TesterRepository.FindByOwner(ownerId)

	if err != nil {
		errMessage := service.I18n.MustLocalize(&i18n.LocalizeConfig{
			MessageID: "tester.repository.find.tester.owner.failed",
			TemplateData: map[string]interface{}{
				"OwnerID": ownerId,
				"Reason":  err.Error(),
			},
			PluralCount: 1,
		})

		return testers, errors.New(errMessage)
	}

	return testers, err
}

func (service *testerServiceImpl) Create(userId uint, name string) error {
	if err := service.TesterRepository.Create(&model.Tester{UserID: userId, Name: name}); err != nil {
		errMessage := service.I18n.MustLocalize(&i18n.LocalizeConfig{
			MessageID: "tester.repository.create.tester.failed",
			TemplateData: map[string]interface{}{
				"Reason": err.Error(),
			},
			PluralCount: 1,
		})

		return errors.New(errMessage)
	}

	return nil
}

func (service *testerServiceImpl) Update(tester *model.Tester) error {
	err := service.TesterRepository.Update(tester)

	if err != nil {
		errMessage := service.I18n.MustLocalize(&i18n.LocalizeConfig{
			MessageID: "tester.repository.update.tester.failed",
			TemplateData: map[string]interface{}{
				"ID":     tester.ID,
				"Reason": err.Error(),
			},
			PluralCount: 1,
		})

		return errors.New(errMessage)
	}

	return nil
}

func (service *testerServiceImpl) Delete(id uint) error {

	tester, err := service.Get(id)

	if err != nil {
		return err
	}

	err = service.TesterRepository.Delete(id)

	if err != nil {
		errMessage := service.I18n.MustLocalize(&i18n.LocalizeConfig{
			MessageID: "tester.repository.delete.tester.failed",
			TemplateData: map[string]interface{}{
				"ID":     id,
				"Reason": err.Error(),
			},
			PluralCount: 1,
		})

		return errors.New(errMessage)
	}

	if err := os.Remove(tester.Path); err != nil {
		errMessage := service.I18n.MustLocalize(&i18n.LocalizeConfig{
			MessageID: "os.cmd.mkdir.dir.failed",
			TemplateData: map[string]interface{}{
				"Path":   tester.Path,
				"Reason": err.Error(),
			},
			PluralCount: 1,
		})

		log.Println(errMessage)
		return errors.New(errMessage)
	}

	return nil
}
