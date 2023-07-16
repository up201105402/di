package repository

import (
	"di/model"
	"errors"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type runRepositoryImpl struct {
	DB *gorm.DB
}

func NewRunRepository(gormDB *gorm.DB) model.RunRepository {
	return &runRepositoryImpl{
		DB: gormDB,
	}
}

func (repo *runRepositoryImpl) FindByID(id uint) (*model.Run, error) {

	var run = model.Run{}

	result := repo.DB.Preload("Pipeline.User").Preload(clause.Associations).First(&run, id)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, result.Error
	}

	return &run, nil
}

func (repo *runRepositoryImpl) FindByPipeline(pipelineId uint) ([]model.Run, error) {

	var runs []model.Run

	result := repo.DB.Preload("Pipeline").Preload("RunStatus").Where("pipeline_id = ?", pipelineId).Find(&runs)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, result.Error
	}

	return runs, nil
}

func (repo *runRepositoryImpl) FindRunStepStatusesByRun(runID uint) ([]model.RunStepStatus, error) {
	var runStepStatuses []model.RunStepStatus

	result := repo.DB.Preload("Run").Preload("RunStatus").Where("run_id = ?", runID).Find(&runStepStatuses)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, result.Error
	}

	return runStepStatuses, nil
}

func (repo *runRepositoryImpl) Create(run *model.Run) error {
	result := repo.DB.Create(run)

	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (repo *runRepositoryImpl) CreateRunStepStatus(runStepStatus *model.RunStepStatus) error {
	result := repo.DB.Create(runStepStatus)

	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (repo *runRepositoryImpl) Update(run *model.Run) error {
	result := repo.DB.Save(run)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return repo.Create(run)
	}

	return nil
}

func (repo *runRepositoryImpl) UpdateRunStepStatus(runStepStatus *model.RunStepStatus) error {
	result := repo.DB.Save(runStepStatus)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return repo.CreateRunStepStatus(runStepStatus)
	}

	return nil
}

func (repo *runRepositoryImpl) Delete(id uint) error {
	result := repo.DB.Delete(&model.Run{}, id)

	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (repo *runRepositoryImpl) DeleteRunStepStatus(id uint) error {
	result := repo.DB.Delete(&model.RunStepStatus{}, id)

	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (repo *runRepositoryImpl) DeleteAllRunStepStatuses(runId uint) error {
	result := repo.DB.Where("run_id = ?", runId).Delete(&model.RunStepStatus{})

	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (repo *runRepositoryImpl) GetRunStatusByID(id uint) (*model.RunStatus, error) {
	var runStatus = model.RunStatus{}

	result := repo.DB.First(&runStatus, id)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, result.Error
	}

	return &runStatus, nil
}
