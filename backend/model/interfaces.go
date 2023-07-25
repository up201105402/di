package model

import (
	"context"
	"time"
)

type TokenRepository interface {
	SetRefreshToken(ctx context.Context, userID uint, tokenID uint, expiresIn time.Duration) error
	DeleteRefreshToken(ctx context.Context, userID uint, prevTokenID uint) error
	DeleteUserRefreshTokens(ctx context.Context, userID uint) error
}

type PipelineRepository interface {
	FindByID(pipelineID uint) (*Pipeline, error)
	FindPipelineScheduleByID(pipelineScheduleID uint) (*PipelineSchedule, error)
	FindPipelineScheduleByPipeline(pipelineID uint) ([]PipelineSchedule, error)
	GetAllPipelineSchedules() ([]PipelineSchedule, error)
	FindByOwner(ownerID uint) ([]Pipeline, error)
	Create(pipeline *Pipeline) error
	CreatePipelineSchedule(pipelineSchedule *PipelineSchedule) error
	Update(pipeline *Pipeline) error
	Delete(pipelineID uint) error
	DeletePipelineSchedule(pipelineID uint) error
}

type DatasetRepository interface {
	FindByID(datasetID uint) (*Dataset, error)
	FindScriptsByDatasetID(datasetID uint) ([]DatasetScript, error)
	FindScriptByID(scriptID uint) (*DatasetScript, error)
	FindByOwner(ownerID uint) ([]Dataset, error)
	Create(dataset *Dataset) error
	CreateDatasetScript(datasetScript *DatasetScript) error
	Update(dataset *Dataset) error
	Delete(datasetID uint) error
	DeleteDatasetScript(datasetScriptId uint) error
}

type TrainerRepository interface {
	FindByID(trainerID uint) (*Trainer, error)
	FindByOwner(ownerID uint) ([]Trainer, error)
	Create(trainer *Trainer) error
	Update(trainer *Trainer) error
	Delete(trainerID uint) error
}

type TesterRepository interface {
	FindByID(testerID uint) (*Tester, error)
	FindByOwner(ownerID uint) ([]Tester, error)
	Create(tester *Tester) error
	Update(tester *Tester) error
	Delete(testerID uint) error
}

type TrainedRepository interface {
	FindByID(trainedID uint) (*Trained, error)
	FindByOwner(ownerID uint) ([]Trained, error)
	Create(trained *Trained) error
	Update(trained *Trained) error
	Delete(trainedID uint) error
}

type RunRepository interface {
	FindByID(runID uint) (*Run, error)
	FindByPipeline(pipelineID uint) ([]Run, error)
	FindRunStepStatusesByRun(runID uint) ([]RunStepStatus, error)
	FindHumanFeedbackQueriesByStepID(runID uint, stepID uint) ([]HumanFeedbackQuery, error)
	FindHumanFeedbackQueriesByRunID(runID uint) ([]HumanFeedbackQuery, error)
	FindHumanFeedbackQueryByID(queryID uint) (*HumanFeedbackQuery, error)
	FindHumanFeedbackRectsByHumanFeedbackQueryID(humanFeedbackQueryID uint) ([]HumanFeedbackRect, error)
	FindHumanFeedbackQueryStatusByID(queryStatusID uint) (*QueryStatus, error)
	Create(run *Run) error
	CreateRunStepStatus(runStepStatus *RunStepStatus) error
	CreateHumanFeedbackQuery(humanFeedbackQuery *HumanFeedbackQuery) error
	CreateHumanFeedbackRect(humanFeedbackRect *HumanFeedbackRect) error
	Update(run *Run) error
	UpdateRunStepStatus(runStepStatus *RunStepStatus) error
	UpdateHumanFeedbackQuery(query *HumanFeedbackQuery) error
	UpdateHumanFeedbackRect(rect *HumanFeedbackRect) error
	Delete(runID uint) error
	DeleteRunStepStatus(runID uint) error
	DeleteAllHumanFeedbackQueriesByRunID(runID uint) error
	DeleteAllRunStepStatuses(runID uint) error
	GetRunStatusByID(runID uint) (*RunStatus, error)
}
