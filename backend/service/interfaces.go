package service

import (
	"context"
	"di/model"
	"di/steps"

	"github.com/hibiken/asynq"
)

type Services struct {
	UserService     UserService
	PipelineService PipelineService
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
	NewPairFromUser(ctx context.Context, u *model.User, prevTokenID uint) (*model.TokenPair, error)
	Signout(ctx context.Context, uid uint) error
	ValidateIDToken(tokenString string) (*model.User, error)
	ValidateRefreshToken(refreshTokenString string) (*model.RefreshToken, error)
}

type PipelineService interface {
	Get(id uint) (*model.Pipeline, error)
	GetByOwner(ownerId uint) ([]model.Pipeline, error)
	Create(userId uint, name string, definition string) error
	Update(pipeline *model.Pipeline) error
	Delete(id uint) error
}

type RunService interface {
	Get(id uint) (*model.Run, error)
	GetByPipeline(pipelineId uint) ([]model.Run, error)
	Create(pipelineId uint) error
	Execute(runID uint) error
	Update(run *model.Run) error
	Delete(id uint) error
	NewRunPipelineTask(pipelineID uint, runID uint, graph string, stepIndex uint) (*asynq.Task, error)
	HandleRunPipelineTask(ctx context.Context, t *asynq.Task) error
}

type NodeTypeService interface {
	NewStepInstance(pipelineID uint, runID uint, nodeType string, nodeDescription model.NodeDescription) (*steps.Step, error)
	NewEdgeInstance(pipelineID uint, runID uint, edgeType string, nodeDescription model.NodeDescription) (*steps.Edge, error)
}

type TaskService interface {
	SetupAsynqWorker()
}
