package service

import (
	"di/steps"
	"os"

	"github.com/dominikbraun/graph"
	"github.com/hibiken/asynq"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"gopkg.in/guregu/null.v4"
)

const (
	RunPipelineTask          = "pipeline:run"
	ScheduledRunPipelineTask = "pipeline:scheduled_run"
)

type taskServiceImpl struct {
	I18n            *i18n.Localizer
	NodeTypeService StepService
	RunService      RunService
}

type RunPipelinePayload struct {
	PipelineID      uint
	RunID           uint
	GraphDefinition string
	StepID          null.Int
}

type ScheduledRunPipelinePayload struct {
	PipelineID         uint
	PipelineScheduleID uint
}

func NewTaskService(i18n *i18n.Localizer, nodeTypeService *StepService, runService *RunService) TaskService {
	return &taskServiceImpl{
		I18n:            i18n,
		NodeTypeService: *nodeTypeService,
		RunService:      *runService,
	}
}

func (service *taskServiceImpl) SetupAsynqWorker() {

	redisHost, exists := os.LookupEnv("REDIS_HOST")

	if !exists {
		panic("REDIS_HOST is not defined!")
	}

	redisPort, exists := os.LookupEnv("REDIS_PORT")

	if !exists {
		panic("REDIS_PORT is not defined!")
	}

	redisConnection := asynq.RedisClientOpt{Addr: redisHost + ":" + redisPort}

	worker := asynq.NewServer(redisConnection, asynq.Config{
		Concurrency: 10,
		Queues: map[string]int{
			"runs": 1,
		},
	})

	mux := asynq.NewServeMux()

	mux.HandleFunc(
		RunPipelineTask,
		service.RunService.HandleRunPipelineTask,
	)

	mux.HandleFunc(
		ScheduledRunPipelineTask,
		service.RunService.HandleScheduledRunPipelineTask,
	)

	if err := worker.Run(mux); err != nil {
		panic("Failed to config Asynq")
	}
}

func stepHash(step steps.Step) int {
	return graph.IntHash(int(step.GetID()))
}
