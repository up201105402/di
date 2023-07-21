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

type RunRepository interface {
	FindByID(runID uint) (*Run, error)
	FindByPipeline(pipelineID uint) ([]Run, error)
	FindRunStepStatusesByRun(runID uint) ([]RunStepStatus, error)
	FindHumanFeedbackQueriesByStepID(stepID uint) ([]HumanFeedbackQuery, error)
	FindHumanFeedbackRectsByHumanFeedbackQueryID(humanFeedbackQueryID uint) ([]HumanFeedbackRect, error)
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
	DeleteAllRunStepStatuses(runID uint) error
	GetRunStatusByID(runID uint) (*RunStatus, error)
}
