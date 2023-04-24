package service

import (
	"di/model"
	"di/repository"

	"gorm.io/gorm"
)

type pipelineServiceImpl struct {
	PipelineRepository model.PipelineRepository
}

func NewPipelineService(gormDB *gorm.DB) model.PipelineService {
	return &pipelineServiceImpl{
		PipelineRepository: repository.NewPipelineRepository(gormDB),
	}
}

func (service *pipelineServiceImpl) Get(id uint) (*model.Pipeline, error) {
	pipeline, error := service.PipelineRepository.FindByID(id)
	return pipeline, error
}

func (service *pipelineServiceImpl) GetByOwner(ownerId uint) ([]model.Pipeline, error) {
	pipelines, error := service.PipelineRepository.FindByOwner(ownerId)
	return pipelines, error
}

func (service *pipelineServiceImpl) Create(userId uint, name string, definition string) error {
	if err := service.PipelineRepository.Create(&model.Pipeline{UserID: userId, Name: name, Definition: definition}); err != nil {
		return err
	}

	return nil
}

func (service *pipelineServiceImpl) Update(pipeline *model.Pipeline) error {
	err := service.PipelineRepository.Update(pipeline)

	if err != nil {
		return err
	}

	return nil
}

func (service *pipelineServiceImpl) Delete(id uint) error {
	err := service.PipelineRepository.Delete(id)

	if err != nil {
		return err
	}

	return nil
}
