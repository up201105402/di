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

type trainerServiceImpl struct {
	TrainerRepository repository.TrainerRepository
	TaskQueueClient   *asynq.Client
	I18n              *i18n.Localizer
}

func NewTrainerService(gormDB *gorm.DB, client *asynq.Client, i18n *i18n.Localizer) TrainerService {
	return &trainerServiceImpl{
		TrainerRepository: repository.NewTrainerRepository(gormDB),
		TaskQueueClient:   client,
		I18n:              i18n,
	}
}

func (service *trainerServiceImpl) Get(id uint) (*model.Trainer, error) {
	trainer, err := service.TrainerRepository.FindByID(id)

	if err != nil {
		errMessage := service.I18n.MustLocalize(&i18n.LocalizeConfig{
			MessageID: "trainer.repository.find.trainer.id.failed",
			TemplateData: map[string]interface{}{
				"ID":     id,
				"Reason": err.Error(),
			},
			PluralCount: 1,
		})

		return trainer, errors.New(errMessage)
	}

	return trainer, err
}

func (service *trainerServiceImpl) GetByOwner(ownerId uint) ([]model.Trainer, error) {
	trainers, err := service.TrainerRepository.FindByOwner(ownerId)

	if err != nil {
		errMessage := service.I18n.MustLocalize(&i18n.LocalizeConfig{
			MessageID: "trainer.repository.find.trainer.owner.failed",
			TemplateData: map[string]interface{}{
				"OwnerID": ownerId,
				"Reason":  err.Error(),
			},
			PluralCount: 1,
		})

		return trainers, errors.New(errMessage)
	}

	return trainers, err
}

func (service *trainerServiceImpl) Create(userId uint, name string) error {
	if err := service.TrainerRepository.Create(&model.Trainer{UserID: userId, Name: name}); err != nil {
		errMessage := service.I18n.MustLocalize(&i18n.LocalizeConfig{
			MessageID: "trainer.repository.create.trainer.failed",
			TemplateData: map[string]interface{}{
				"Reason": err.Error(),
			},
			PluralCount: 1,
		})

		return errors.New(errMessage)
	}

	return nil
}

func (service *trainerServiceImpl) Update(trainer *model.Trainer) error {
	err := service.TrainerRepository.Update(trainer)

	if err != nil {
		errMessage := service.I18n.MustLocalize(&i18n.LocalizeConfig{
			MessageID: "trainer.repository.update.trainer.failed",
			TemplateData: map[string]interface{}{
				"ID":     trainer.ID,
				"Reason": err.Error(),
			},
			PluralCount: 1,
		})

		return errors.New(errMessage)
	}

	return nil
}

func (service *trainerServiceImpl) Delete(id uint) error {

	trainer, err := service.Get(id)

	if err != nil {
		return err
	}

	err = service.TrainerRepository.Delete(id)

	if err != nil {
		errMessage := service.I18n.MustLocalize(&i18n.LocalizeConfig{
			MessageID: "trainer.repository.delete.trainer.failed",
			TemplateData: map[string]interface{}{
				"ID":     id,
				"Reason": err.Error(),
			},
			PluralCount: 1,
		})

		return errors.New(errMessage)
	}

	if err := os.Remove(trainer.Path); err != nil {
		errMessage := service.I18n.MustLocalize(&i18n.LocalizeConfig{
			MessageID: "os.cmd.mkdir.dir.failed",
			TemplateData: map[string]interface{}{
				"Path":   trainer.Path,
				"Reason": err.Error(),
			},
			PluralCount: 1,
		})

		log.Println(errMessage)
		return errors.New(errMessage)
	}

	return nil
}
