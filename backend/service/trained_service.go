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

type trainedServiceImpl struct {
	TrainedRepository repository.TrainedModelRepository
	TaskQueueClient   *asynq.Client
	I18n              *i18n.Localizer
}

func NewTrainedService(gormDB *gorm.DB, client *asynq.Client, i18n *i18n.Localizer) TrainedModelService {
	return &trainedServiceImpl{
		TrainedRepository: repository.NewTrainedRepository(gormDB),
		TaskQueueClient:   client,
		I18n:              i18n,
	}
}

func (service *trainedServiceImpl) Get(id uint) (*model.Trained, error) {
	trained, err := service.TrainedRepository.FindByID(id)

	if err != nil {
		errMessage := service.I18n.MustLocalize(&i18n.LocalizeConfig{
			MessageID: "trained.repository.find.trained.id.failed",
			TemplateData: map[string]interface{}{
				"ID":     id,
				"Reason": err.Error(),
			},
			PluralCount: 1,
		})

		return trained, errors.New(errMessage)
	}

	return trained, err
}

func (service *trainedServiceImpl) GetByOwner(ownerId uint) ([]model.Trained, error) {
	traineds, err := service.TrainedRepository.FindByOwner(ownerId)

	if err != nil {
		errMessage := service.I18n.MustLocalize(&i18n.LocalizeConfig{
			MessageID: "trained.repository.find.trained.owner.failed",
			TemplateData: map[string]interface{}{
				"OwnerID": ownerId,
				"Reason":  err.Error(),
			},
			PluralCount: 1,
		})

		return traineds, errors.New(errMessage)
	}

	return traineds, err
}

func (service *trainedServiceImpl) Create(userId uint, name string) (*model.Trained, error) {
	trained := &model.Trained{UserID: userId, Name: name}
	if err := service.TrainedRepository.Create(trained); err != nil {
		errMessage := service.I18n.MustLocalize(&i18n.LocalizeConfig{
			MessageID: "trained.repository.create.trained.failed",
			TemplateData: map[string]interface{}{
				"Reason": err.Error(),
			},
			PluralCount: 1,
		})

		return nil, errors.New(errMessage)
	}

	return trained, nil
}

func (service *trainedServiceImpl) Update(trained *model.Trained) error {
	err := service.TrainedRepository.Update(trained)

	if err != nil {
		errMessage := service.I18n.MustLocalize(&i18n.LocalizeConfig{
			MessageID: "trained.repository.update.trained.failed",
			TemplateData: map[string]interface{}{
				"ID":     trained.ID,
				"Reason": err.Error(),
			},
			PluralCount: 1,
		})

		return errors.New(errMessage)
	}

	return nil
}

func (service *trainedServiceImpl) Delete(id uint) error {

	trained, err := service.Get(id)

	if err != nil {
		return err
	}

	err = service.TrainedRepository.Delete(id)

	if err != nil {
		errMessage := service.I18n.MustLocalize(&i18n.LocalizeConfig{
			MessageID: "trained.repository.delete.trained.failed",
			TemplateData: map[string]interface{}{
				"ID":     id,
				"Reason": err.Error(),
			},
			PluralCount: 1,
		})

		return errors.New(errMessage)
	}

	if err := os.Remove(trained.Path); err != nil {
		errMessage := service.I18n.MustLocalize(&i18n.LocalizeConfig{
			MessageID: "os.cmd.mkdir.dir.failed",
			TemplateData: map[string]interface{}{
				"Path":   trained.Path,
				"Reason": err.Error(),
			},
			PluralCount: 1,
		})

		log.Println(errMessage)
		return errors.New(errMessage)
	}

	return nil
}
