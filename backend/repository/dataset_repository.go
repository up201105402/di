package repository

import (
	"di/model"
	"errors"
	"log"

	"gorm.io/gorm"
)

type datasetRepositoryImpl struct {
	DB *gorm.DB
}

// NewPipelineRepository is a Pipeline Repository factory
func NewDatasetRepository(gormDB *gorm.DB) DatasetRepository {
	return &datasetRepositoryImpl{
		DB: gormDB,
	}
}

func (repo *datasetRepositoryImpl) FindByID(id uint) (*model.Dataset, error) {
	var dataset = model.Dataset{}

	result := repo.DB.Preload("User").First(&dataset, id)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, result.Error
	}

	return &dataset, nil
}

func (repo *datasetRepositoryImpl) FindByOwner(ownerId uint) ([]model.Dataset, error) {

	var datasets []model.Dataset

	result := repo.DB.Preload("User").Where("user_id = ?", ownerId).Order("id desc").Find(&datasets)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, result.Error
	}

	return datasets, nil
}

func (repo *datasetRepositoryImpl) FindScriptsByDatasetID(datasetId uint) ([]model.DatasetScript, error) {

	var datasetsScripts []model.DatasetScript

	result := repo.DB.Where("dataset_id = ?", datasetId).Order("id desc").Find(&datasetsScripts)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, result.Error
	}

	return datasetsScripts, nil
}

func (repo *datasetRepositoryImpl) FindScriptByID(scriptID uint) (*model.DatasetScript, error) {

	var datasetsScript *model.DatasetScript

	result := repo.DB.Where("id = ?", scriptID).Find(&datasetsScript)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, result.Error
	}

	return datasetsScript, nil
}

func (repo *datasetRepositoryImpl) Create(dataset *model.Dataset) error {
	result := repo.DB.Create(dataset)

	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (repo *datasetRepositoryImpl) CreateDatasetScript(datasetScript *model.DatasetScript) error {
	result := repo.DB.Create(datasetScript)

	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (repo *datasetRepositoryImpl) Update(dataset *model.Dataset) error {
	result := repo.DB.Save(dataset)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		repo.Create(dataset)
	}

	return nil
}

func (repo *datasetRepositoryImpl) Delete(id uint) error {
	result := repo.DB.Delete(&model.Dataset{}, id)

	if result.Error != nil {
		log.Printf("Failed to delete pipeline. Reason: %v\n", result.Error)
		return result.Error
	}

	return nil
}

func (repo *datasetRepositoryImpl) DeleteDatasetScript(datasetScriptId uint) error {
	result := repo.DB.Delete(&model.DatasetScript{}, datasetScriptId)

	if result.Error != nil {
		return result.Error
	}

	return nil
}
