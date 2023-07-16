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

	result := repo.DB.Preload("User").First(&pipeline, id)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, result.Error
	}

	return &pipeline, nil
}

func (repo *pipelineRepositoryImpl) FindPipelineScheduleByID(pipelineScheduleID uint) (*model.PipelineSchedule, error) {
	var pipelineSchedule = model.PipelineSchedule{}

	result := repo.DB.First(&pipelineSchedule, pipelineScheduleID)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, result.Error
	}

	return &pipelineSchedule, nil
}

func (repo *pipelineRepositoryImpl) GetAllPipelineSchedules() ([]model.PipelineSchedule, error) {
	var pipelineSchedules []model.PipelineSchedule

	result := repo.DB.Find(&pipelineSchedules)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, result.Error
	}

	return pipelineSchedules, nil
}

func (repo *pipelineRepositoryImpl) FindByOwner(ownerId uint) ([]model.Pipeline, error) {

	var pipelines []model.Pipeline

	result := repo.DB.Preload("User").Where("user_id = ?", ownerId).Order("id").Find(&pipelines)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, result.Error
	}

	return pipelines, nil
}

func (repo *pipelineRepositoryImpl) FindPipelineScheduleByPipeline(pipelineID uint) ([]model.PipelineSchedule, error) {
	var pipelineSchedules []model.PipelineSchedule

	result := repo.DB.Where("pipeline_id = ?", pipelineID).Order("id").Find(&pipelineSchedules)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, result.Error
	}

	return pipelineSchedules, nil
}

func (repo *pipelineRepositoryImpl) Create(pipeline *model.Pipeline) error {
	result := repo.DB.Create(pipeline)

	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (repo *pipelineRepositoryImpl) CreatePipelineSchedule(pipelineSchedule *model.PipelineSchedule) error {
	result := repo.DB.Create(pipelineSchedule)

	if result.Error != nil {
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

func (repo *pipelineRepositoryImpl) DeletePipelineSchedule(id uint) error {
	result := repo.DB.Delete(&model.PipelineSchedule{}, id)

	if result.Error != nil {
		log.Printf("Failed to delete pipeline schedule. Reason: %v\n", result.Error)
		return result.Error
	}

	return nil
}
