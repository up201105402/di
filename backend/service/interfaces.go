package service

import (
	"context"
	"di/model"
	"di/steps"
	"time"

	"github.com/hibiken/asynq"
)

type Services struct {
	UserService     UserService
	PipelineService PipelineService
	DatasetService  DatasetService
	TrainerService  TrainerService
	TesterService   TesterService
	TrainedService  TrainedModelService
	RunService      RunService
	TokenService    TokenService
}

type UserService interface {
	Get(id uint) (*model.User, error)
	GetByUsername(username string) (*model.User, error)
	Signup(username string, password string) error
	Signin(user *model.User) error
	UpdateDetails(user *model.User) error
}

type TokenService interface {
	NewFirstPairFromUser(ctx context.Context, u *model.User) (*model.TokenPair, error)
	NewPairFromUser(ctx context.Context, u *model.User, refreshToken model.RefreshToken) (*model.TokenPair, error)
	Signout(ctx context.Context, uid uint) error
	ValidateIDToken(tokenString string) (*model.User, error)
	ValidateRefreshToken(refreshTokenString string) (*model.RefreshToken, error)
}

type PipelineService interface {
	SyncAsyncTasks()
	Get(id uint) (*model.Pipeline, error)
	GetPipelineSchedule(id uint) (*model.PipelineSchedule, error)
	GetPipelineSchedules(id uint) ([]model.PipelineSchedule, error)
	GetByOwner(ownerId uint) ([]model.Pipeline, error)
	Create(userId uint, name string, definition string) error
	CreatePipelineSchedule(pipelineID uint, uniqueOcurrence time.Time, cronExpression string) error
	Update(pipeline *model.Pipeline) error
	Delete(id uint) error
	DeletePipelineSchedule(id uint) error
}

type DatasetService interface {
	Get(id uint) (*model.Dataset, error)
	GetDatasetScripts(id uint) ([]model.DatasetScript, error)
	GetDatasetScript(scriptID uint) (*model.DatasetScript, error)
	GetByOwner(ownerId uint) ([]model.Dataset, error)
	Create(userId uint, name string) error
	CreateDatasetScript(datasetID uint, scriptName string, filePath string) error
	Update(dataset *model.Dataset) error
	Delete(id uint) error
	DeleteDatasetScript(datasetScriptId uint) error
}

type TrainerService interface {
	Get(id uint) (*model.Trainer, error)
	GetByOwner(ownerId uint) ([]model.Trainer, error)
	Create(userId uint, name string) error
	Update(dataset *model.Trainer) error
	Delete(id uint) error
}

type TesterService interface {
	Get(id uint) (*model.Tester, error)
	GetByOwner(ownerId uint) ([]model.Tester, error)
	Create(userId uint, name string) error
	Update(tester *model.Tester) error
	Delete(id uint) error
}

type TrainedModelService interface {
	Get(id uint) (*model.Trained, error)
	GetByOwner(ownerId uint) ([]model.Trained, error)
	Create(userId uint, name string) (*model.Trained, error)
	Update(trained *model.Trained) error
	Delete(id uint) error
}

type RunService interface {
	Get(id uint) (*model.Run, error)
	GetByPipeline(pipelineId uint) ([]model.Run, error)
	FindRunStepStatusesByRun(runID uint) ([]model.RunStepStatus, error)
	FindHumanFeedbackQueriesByRunID(runID uint) ([]model.HumanFeedbackQuery, error)
	FindHumanFeedbackQueryByID(queryID uint) (*model.HumanFeedbackQuery, error)
	FindHumanFeedbackQueriesByStepID(runID uint, runStepStatusID uint) ([]model.HumanFeedbackQuery, error)
	FindHumanFeedbackRectsByHumanFeedbackQueryID(humanFeedbackQueryID uint) ([]model.HumanFeedbackRect, error)
	FindHumanFeedbackQueryStatusByID(queryStatusID uint) (*model.QueryStatus, error)
	Create(pipeline model.Pipeline) (model.Run, error)
	CreateRunStepStatus(runID uint, stepID int, stepName string, runStatusID uint, errorMessage string) error
	CreateHumanFeedbackQuery(epoch uint, runID uint, stepID int, queryID uint, rects [][]uint) error
	Execute(runID uint) error
	Resume(runID uint) error
	Update(run *model.Run) error
	UpdateRunStepStatus(run *model.RunStepStatus) error
	UpdateHumanFeedbackQuery(query *model.HumanFeedbackQuery) error
	UpdateHumanFeedbackRects(rects []model.HumanFeedbackRect) error
	Delete(id uint) error
	DeleteRunStepStatus(id uint) error
	NewRunPipelineTask(pipelineID uint, runID uint, graph string) (*asynq.Task, error)
	NewResumeRunPipelineTask(pipelineID uint, runID uint, graph string, stepID int) (*asynq.Task, error)
	HandleRunPipelineTask(ctx context.Context, t *asynq.Task) error
	HandleScheduledRunPipelineTask(ctx context.Context, t *asynq.Task) error
	UpdateRunStatus(runID uint, statusID uint, stepWaitingFeedback int, errorMessage string) error
}

type RunStepStatusService interface {
	Get(id uint) (*model.Run, error)
	GetByRun(runID uint) ([]model.Run, error)
	Create(runID uint, runStatus model.RunStatus, errorMessage string) error
	Update(run *model.RunStepStatus) error
	Delete(id uint) error
	UpdateRunStatus()
}

type StepService interface {
	NewStepInstance(pipelineID uint, runID uint, nodeType string, nodeDescription model.NodeDescription) (*steps.Step, error)
	NewEdgeInstance(pipelineID uint, runID uint, edgeType string, nodeDescription model.NodeDescription) (*steps.Edge, error)
}

type TaskService interface {
	SetupAsynqWorker()
}
