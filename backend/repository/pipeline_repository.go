package repository

import (
	"di/model"
	"errors"
	"log"

	"gorm.io/gorm"
)

type pipelineRepositoryImpl struct {
	DB *gorm.DB
}

// NewPipelineRepository is a Pipeline Repository factory
func NewPipelineRepository(gormDB *gorm.DB) model.PipelineRepository {
	return &pipelineRepositoryImpl{
		DB: gormDB,
	}
}

func (repo *pipelineRepositoryImpl) FindByID(id uint) (*model.Pipeline, error) {

	var pipeline = model.Pipeline{}

	result := repo.DB.First(&pipeline, id)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		log.Printf("Failed to get pipeline with id: %v. Reason: %v\n", id, result.Error)
		return nil, result.Error
	}

	return &pipeline, nil
}

func (repo *pipelineRepositoryImpl) FindByOwner(ownerId uint) ([]model.Pipeline, error) {

	var pipelines []model.Pipeline

	result := repo.DB.Preload("User").Where("user_id = ?", ownerId).Order("id").Find(&pipelines)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		log.Printf("Failed to get pipelines with username: %v. Reason: %v\n", ownerId, result.Error)
		return nil, result.Error
	}

	return pipelines, nil
}

func (repo *pipelineRepositoryImpl) Create(pipeline *model.Pipeline) error {
	result := repo.DB.Create(pipeline)

	if result.Error != nil {
		log.Printf("Failed to create pipeline. Reason: %v\n", result.Error)
		return result.Error
	}

	return nil
}

func (repo *pipelineRepositoryImpl) Update(pipeline *model.Pipeline) error {
	result := repo.DB.Save(pipeline)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		repo.Create(pipeline)
	}

	return nil
}

func (repo *pipelineRepositoryImpl) Delete(id uint) error {
	result := repo.DB.Delete(&model.Pipeline{}, id)

	if result.Error != nil {
		log.Printf("Failed to delete pipeline. Reason: %v\n", result.Error)
		return result.Error
	}

	return nil
}
