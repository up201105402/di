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
	FindScheduleByID(pipelineID uint) ([]PipelineSchedule, error)
	FindByOwner(ownerID uint) ([]Pipeline, error)
	Create(pipeline *Pipeline) error
	CreateSchedule(pipeline *PipelineSchedule) error
	Update(pipeline *Pipeline) error
	Delete(pipelineID uint) error
}

type RunRepository interface {
	FindByID(runID uint) (*Run, error)
	FindByPipeline(pipelineID uint) ([]Run, error)
	FindRunStepStatusesByRun(runID uint) ([]RunStepStatus, error)
	Create(run *Run) error
	CreateRunStepStatus(runStepStatus *RunStepStatus) error
	Update(run *Run) error
	UpdateRunStepStatus(runStepStatus *RunStepStatus) error
	Delete(runID uint) error
	DeleteRunStepStatus(runID uint) error
	DeleteAllRunStepStatuses(runID uint) error
	GetRunStatusByID(runID uint) (*RunStatus, error)
}
