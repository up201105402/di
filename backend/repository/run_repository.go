package repository

import (
	"di/model"
	"errors"
	"log"

	"gorm.io/gorm"
)

type runRepositoryImpl struct {
	DB *gorm.DB
}

// NewPipelineRepository is a Pipeline Repository factory
func NewRunRepository(gormDB *gorm.DB) model.RunRepository {
	return &runRepositoryImpl{
		DB: gormDB,
	}
}

func (repo *runRepositoryImpl) FindByID(id uint) (*model.Run, error) {

	var run = model.Run{}

	result := repo.DB.First(&run, id)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		log.Printf("Failed to get pipeline with id: %v. Reason: %v\n", id, result.Error)
		return nil, result.Error
	}

	return &run, nil
}

func (repo *runRepositoryImpl) FindByPipeline(pipelineId uint) ([]model.Run, error) {

	var runs []model.Run

	result := repo.DB.Preload("Pipeline").Where("pipeline_id = ?", pipelineId).Find(&runs)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		log.Printf("Failed to get runs for pipeline %v. Reason: %v\n", pipelineId, result.Error)
		return nil, result.Error
	}

	return runs, nil
}

func (repo *runRepositoryImpl) Create(run *model.Run) error {
	result := repo.DB.Create(run)

	if result.Error != nil {
		log.Printf("Failed to create run. Reason: %v\n", result.Error)
		return result.Error
	}

	return nil
}

func (repo *runRepositoryImpl) Update(run *model.Run) error {
	result := repo.DB.Save(run)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		repo.Create(run)
	}

	return nil
}

func (repo *runRepositoryImpl) Delete(id uint) error {
	result := repo.DB.Delete(&model.Run{}, id)

	if result.Error != nil {
		log.Printf("Failed to delete run. Reason: %v\n", result.Error)
		return result.Error
	}

	return nil
}
