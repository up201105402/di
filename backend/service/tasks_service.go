package service

import (
	"di/steps"
	"os"

	"github.com/dominikbraun/graph"
	"github.com/hibiken/asynq"
)

const (
	RunPipelineTask = "pipeline:run"
)

type taskServiceImpl struct {
	NodeTypeService NodeTypeService
	RunService      RunService
}

type RunPipelinePayload struct {
	PipelineID      uint
	RunID           uint
	GraphDefinition string
	//Graph     graph.Graph[int, steps.Step]
	StepIndex uint
}

func NewTaskService(nodeTypeService *NodeTypeService, runService *RunService) TaskService {
	return &taskServiceImpl{
		NodeTypeService: *nodeTypeService,
		RunService:      *runService,
	}
}

func (service *taskServiceImpl) SetupAsynqWorker() {

	redisHost := os.Getenv("REDIS_HOST")
	redisPort := os.Getenv("REDIS_PORT")
	redisConnection := asynq.RedisClientOpt{
		Addr: redisHost + ":" + redisPort,
	}

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

	if err := worker.Run(mux); err != nil {
		panic("Failed to config Asynq")
	}
}

func stepHash(step steps.Step) int {
	return graph.IntHash(int(step.GetID()))
}
