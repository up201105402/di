package service

import (
	"di/model"
	"di/repository"
	"errors"

	"github.com/hibiken/asynq"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"gorm.io/gorm"
)

type datasetServiceImpl struct {
	DatasetRepository model.DatasetRepository
	TaskQueueClient   *asynq.Client
	I18n              *i18n.Localizer
}

func NewDatasetService(gormDB *gorm.DB, client *asynq.Client, i18n *i18n.Localizer) DatasetService {
	return &datasetServiceImpl{
		DatasetRepository: repository.NewDatasetRepository(gormDB),
		TaskQueueClient:   client,
		I18n:              i18n,
	}
}

func (service *datasetServiceImpl) Get(id uint) (*model.Dataset, error) {
	dataset, err := service.DatasetRepository.FindByID(id)

	if err != nil {
		errMessage := service.I18n.MustLocalize(&i18n.LocalizeConfig{
			MessageID: "dataset.repository.find.dataset.id.failed",
			TemplateData: map[string]interface{}{
				"ID":     id,
				"Reason": err.Error(),
			},
			PluralCount: 1,
		})

		return dataset, errors.New(errMessage)
	}

	return dataset, err
}

func (service *datasetServiceImpl) GetDatasetScripts(id uint) ([]model.DatasetScript, error) {
	datasetScripts, err := service.DatasetRepository.FindScriptsByDatasetID(id)

	if err != nil {
		errMessage := service.I18n.MustLocalize(&i18n.LocalizeConfig{
			MessageID: "dataset.repository.find.scripts.id.failed",
			TemplateData: map[string]interface{}{
				"ID":     id,
				"Reason": err.Error(),
			},
			PluralCount: 1,
		})

		return datasetScripts, errors.New(errMessage)
	}

	return datasetScripts, err
}

func (service *datasetServiceImpl) GetDatasetScript(scriptID uint) (*model.DatasetScript, error) {
	datasetScripts, err := service.DatasetRepository.FindScriptByID(scriptID)

	if err != nil {
		errMessage := service.I18n.MustLocalize(&i18n.LocalizeConfig{
			MessageID: "dataset.repository.find.scripts.id.failed",
			TemplateData: map[string]interface{}{
				"ID":     scriptID,
				"Reason": err.Error(),
			},
			PluralCount: 1,
		})

		return datasetScripts, errors.New(errMessage)
	}

	return datasetScripts, err
}

func (service *datasetServiceImpl) CreateDatasetScript(datasetID uint, scriptName string, filePath string) error {
	if err := service.DatasetRepository.CreateDatasetScript(&model.DatasetScript{Name: scriptName, Path: filePath, DatasetID: datasetID}); err != nil {
		errMessage := service.I18n.MustLocalize(&i18n.LocalizeConfig{
			MessageID: "dataset.repository.create.dataset.failed",
			TemplateData: map[string]interface{}{
				"Reason": err.Error(),
			},
			PluralCount: 1,
		})

		return errors.New(errMessage)
	}

	return nil
}

func (service *datasetServiceImpl) GetByOwner(ownerId uint) ([]model.Dataset, error) {
	datasets, err := service.DatasetRepository.FindByOwner(ownerId)

	if err != nil {
		errMessage := service.I18n.MustLocalize(&i18n.LocalizeConfig{
			MessageID: "dataset.repository.find.dataset.owner.failed",
			TemplateData: map[string]interface{}{
				"OwnerID": ownerId,
				"Reason":  err.Error(),
			},
			PluralCount: 1,
		})

		return datasets, errors.New(errMessage)
	}

	return datasets, err
}

func (service *datasetServiceImpl) Create(userId uint, name string, entryPoint string) error {
	if err := service.DatasetRepository.Create(&model.Dataset{UserID: userId, Name: name, EntryPoint: entryPoint}); err != nil {
		errMessage := service.I18n.MustLocalize(&i18n.LocalizeConfig{
			MessageID: "dataset.repository.create.dataset.failed",
			TemplateData: map[string]interface{}{
				"Reason": err.Error(),
			},
			PluralCount: 1,
		})

		return errors.New(errMessage)
	}

	return nil
}

func (service *datasetServiceImpl) Update(dataset *model.Dataset) error {
	err := service.DatasetRepository.Update(dataset)

	if err != nil {
		errMessage := service.I18n.MustLocalize(&i18n.LocalizeConfig{
			MessageID: "dataset.repository.update.dataset.failed",
			TemplateData: map[string]interface{}{
				"ID":     dataset.ID,
				"Reason": err.Error(),
			},
			PluralCount: 1,
		})

		return errors.New(errMessage)
	}

	return nil
}

func (service *datasetServiceImpl) Delete(id uint) error {
	err := service.DatasetRepository.Delete(id)

	if err != nil {
		errMessage := service.I18n.MustLocalize(&i18n.LocalizeConfig{
			MessageID: "dataset.repository.delete.dataset.failed",
			TemplateData: map[string]interface{}{
				"ID":     id,
				"Reason": err.Error(),
			},
			PluralCount: 1,
		})

		return errors.New(errMessage)
	}

	return nil
}

func (service *datasetServiceImpl) DeleteDatasetScript(datasetScriptId uint) error {
	err := service.DatasetRepository.DeleteDatasetScript(datasetScriptId)

	if err != nil {
		errMessage := service.I18n.MustLocalize(&i18n.LocalizeConfig{
			MessageID: "dataset.repository.delete.dataset.failed",
			TemplateData: map[string]interface{}{
				"ID":     datasetScriptId,
				"Reason": err.Error(),
			},
			PluralCount: 1,
		})

		return errors.New(errMessage)
	}

	return nil
}
