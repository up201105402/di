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

	result := repo.DB.Preload("Pipeline").Preload("RunStatus").Where("pipeline_id = ?", pipelineId).Order("id desc").Find(&runs)

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

func (repo *runRepositoryImpl) FindHumanFeedbackQueriesByStepID(runID uint, stepID uint) ([]model.HumanFeedbackQuery, error) {
	var humanFeedbackQueries []model.HumanFeedbackQuery

	result := repo.DB.Preload("QueryStatus").Where("run_id = ? and step_id = ?", runID, stepID).Find(&humanFeedbackQueries)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, result.Error
	}

	return humanFeedbackQueries, nil
}

func (repo *runRepositoryImpl) FindHumanFeedbackRectsByHumanFeedbackQueryID(humanFeedbackQueryID uint) ([]model.HumanFeedbackRect, error) {
	var humanFeedbackRects []model.HumanFeedbackRect

	result := repo.DB.Preload("HumanFeedbackQuery").Where("human_feedback_query_id = ?", humanFeedbackQueryID).Find(&humanFeedbackRects)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, result.Error
	}

	return humanFeedbackRects, nil
}

func (repo *runRepositoryImpl) FindHumanFeedbackQueryStatusByID(queryStatusID uint) (*model.QueryStatus, error) {
	var queryStatus model.QueryStatus

	result := repo.DB.Where("id = ?", queryStatusID).Find(&queryStatus)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, result.Error
	}

	return &queryStatus, nil
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

func (repo *runRepositoryImpl) CreateHumanFeedbackQuery(humanFeedbackQuery *model.HumanFeedbackQuery) error {
	result := repo.DB.Create(humanFeedbackQuery)

	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (repo *runRepositoryImpl) CreateHumanFeedbackRect(humanFeedbackRect *model.HumanFeedbackRect) error {
	result := repo.DB.Create(humanFeedbackRect)

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

func (repo *runRepositoryImpl) UpdateHumanFeedbackQuery(query *model.HumanFeedbackQuery) error {
	result := repo.DB.Save(query)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return repo.CreateHumanFeedbackQuery(query)
	}

	return nil
}

func (repo *runRepositoryImpl) UpdateHumanFeedbackRect(rect *model.HumanFeedbackRect) error {
	result := repo.DB.Save(rect)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return repo.CreateHumanFeedbackRect(rect)
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
