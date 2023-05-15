package service

import (
	"di/model"
	"di/repository"

	"gorm.io/gorm"
)

type runServiceImpl struct {
	RunRepository model.RunRepository
}

func NewRunService(gormDB *gorm.DB) RunService {
	return &runServiceImpl{
		RunRepository: repository.NewRunRepository(gormDB),
	}
}

func (service *runServiceImpl) Get(id uint) (*model.Run, error) {
	run, error := service.RunRepository.FindByID(id)
	return run, error
}

func (service *runServiceImpl) GetByPipeline(pipelineId uint) ([]model.Run, error) {
	pipelines, error := service.RunRepository.FindByPipeline(pipelineId)
	return pipelines, error
}

func (service *runServiceImpl) Create(pipelineId uint) error {
	// Add Initial Status
	if err := service.RunRepository.Create(&model.Run{PipelineID: pipelineId}); err != nil {
		return err
	}

	return nil
}

func (service *runServiceImpl) Execute(pipelineId uint) error {
	// TODO

	return nil
}

func (service *runServiceImpl) Update(run *model.Run) error {
	err := service.RunRepository.Update(run)

	if err != nil {
		return err
	}

	return nil
}

func (service *runServiceImpl) Delete(id uint) error {
	err := service.RunRepository.Delete(id)

	if err != nil {
		return err
	}

	return nil
}
