package service

import (
	"context"
	"di/model"
	"di/steps"
	"di/util"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/dominikbraun/graph"
	"github.com/hibiken/asynq"
)

const (
	RunPipelineTask = "pipeline:run"
)

type taskServiceImpl struct {
	NodeTypeService NodeTypeService
}

type RunPipelinePayload struct {
	PipelineID      uint
	GraphDefinition string
	//Graph     graph.Graph[int, steps.Step]
	StepIndex uint
}

func NewTaskService(nodeTypeService *NodeTypeService) TaskService {
	return &taskServiceImpl{
		NodeTypeService: *nodeTypeService,
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
		service.HandleRunPipelineTask,
	)

	if err := worker.Run(mux); err != nil {
		panic("Failed to config Asynq")
	}
}

func (service *taskServiceImpl) NewRunPipelineTask(pipelineID uint, graph string, stepIndex uint) (*asynq.Task, error) {
	payload, err := json.Marshal(RunPipelinePayload{GraphDefinition: graph, StepIndex: stepIndex})
	if err != nil {
		return nil, err
	}
	return asynq.NewTask(RunPipelineTask, payload), nil
}

func (service *taskServiceImpl) HandleRunPipelineTask(ctx context.Context, t *asynq.Task) error {
	var runPipelinePayload RunPipelinePayload
	if err := json.Unmarshal(t.Payload(), &runPipelinePayload); err != nil {
		log.Println("json.Unmarshal failed: %v: %w", err, asynq.SkipRetry)
		errStr := fmt.Errorf("json.Unmarshal failed: %v: %w", err, asynq.SkipRetry)
		return errStr
	}

	var stepDescriptions []model.NodeDescription

	if err := json.Unmarshal([]byte(runPipelinePayload.GraphDefinition), &stepDescriptions); err != nil {
		log.Printf("Unable to unmarshal pipeline %v definition\n", runPipelinePayload.PipelineID)
		return err
	}

	pipelineGraph := graph.New(stepHash, graph.Directed(), graph.Acyclic())

	stps := util.Map(stepDescriptions, func(stepDescription model.NodeDescription) steps.Step {
		step, _ := service.NodeTypeService.NewStepInstance(runPipelinePayload.PipelineID, stepDescription.Type, stepDescription.Data.StepConfig)

		if step != nil {
			return *step
		}

		return nil
	})

	edgs := util.Map(stepDescriptions, func(stepDescription model.NodeDescription) steps.Edge {
		edge, _ := service.NodeTypeService.NewEdgeInstance(runPipelinePayload.PipelineID, stepDescription.Type, stepDescription.Data.StepConfig)

		if edge != nil {
			return *edge
		}

		return nil
	})

	stps = util.Filter(stps, func(step steps.Step) bool {
		return step != nil
	})

	edgs = util.Filter(edgs, func(edge steps.Edge) bool {
		return edge != nil
	})

	for _, step := range stps {
		pipelineGraph.AddVertex(step)
	}

	for _, edge := range edgs {
		previousStepID := (*edge.GetPreviousStep()).GetID()
		nextStepID := (*edge.GetNextStep()).GetID()
		pipelineGraph.AddEdge(previousStepID, nextStepID)
	}

	graph.BFS(pipelineGraph, 0, func(id int) bool {
		log.Printf("BFS %v\n", id)
		return false
	})

	return asynq.SkipRetry
}

func stepHash(step steps.Step) int {
	return graph.IntHash(int(step.GetID()))
}
