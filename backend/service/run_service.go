package service

import (
	"context"
	"di/model"
	"di/repository"
	"di/steps"
	"di/util"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/dominikbraun/graph"
	"github.com/hibiken/asynq"
	"gorm.io/gorm"
)

type runServiceImpl struct {
	RunRepository    model.RunRepository
	PipelineService  PipelineService
	NodeTypeService  NodeTypeService
	TasksQueueClient asynq.Client
}

func NewRunService(gormDB *gorm.DB, client *asynq.Client, pipelineService *PipelineService, stepTypeService *NodeTypeService) RunService {
	return &runServiceImpl{
		RunRepository:    repository.NewRunRepository(gormDB),
		PipelineService:  *pipelineService,
		NodeTypeService:  *stepTypeService,
		TasksQueueClient: *client,
	}
}

func (service *runServiceImpl) Get(id uint) (*model.Run, error) {
	run, error := service.RunRepository.FindByID(id)
	return run, error
}

func (service *runServiceImpl) GetByPipeline(pipelineId uint) ([]model.Run, error) {
	pipelines, error := service.RunRepository.FindByPipeline(pipelineId)
	return pipelines, error
}

func (service *runServiceImpl) Create(pipelineId uint) error {
	// Add Initial Status
	newRun := &model.Run{PipelineID: pipelineId, StatusID: 1}
	if err := service.RunRepository.Create(newRun); err != nil {
		return err
	}

	return nil
}

func (service *runServiceImpl) Execute(runID uint) error {
	// demarshal stringified pipeline definition json

	run, err := service.RunRepository.FindByID(runID)

	if err != nil {
		log.Printf("Could not retrieve run with id %v\n", runID)
		return err
	}

	pipeline, err := service.PipelineService.Get(run.PipelineID)

	if err != nil {
		log.Printf("Could not retrieve pipeline with id %v\n", run.PipelineID)
		return err
	}

	runPipelineTask, err := service.NewRunPipelineTask(pipeline.ID, runID, pipeline.Definition, 0)

	if err != nil {
		return err
	}

	if _, err := service.TasksQueueClient.Enqueue(
		runPipelineTask,
		asynq.Queue("runs"),
	); err != nil {
		return err
	}

	return nil
}

func (service *runServiceImpl) Update(run *model.Run) error {
	err := service.RunRepository.Update(run)

	if err != nil {
		return err
	}

	return nil
}

func (service *runServiceImpl) Delete(id uint) error {
	err := service.RunRepository.Delete(id)

	if err != nil {
		return err
	}

	return nil
}

func (service *runServiceImpl) NewRunPipelineTask(pipelineID uint, runID uint, graph string, stepIndex uint) (*asynq.Task, error) {
	payload, err := json.Marshal(RunPipelinePayload{PipelineID: pipelineID, RunID: runID, GraphDefinition: graph, StepIndex: stepIndex})
	if err != nil {
		return nil, err
	}
	return asynq.NewTask(RunPipelineTask, payload), nil
}

func (service *runServiceImpl) HandleRunPipelineTask(ctx context.Context, t *asynq.Task) error {
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
		step, _ := service.NodeTypeService.NewStepInstance(runPipelinePayload.PipelineID, runPipelinePayload.RunID, stepDescription.Data.Type, stepDescription)

		if step != nil {
			return *step
		}

		return nil
	})

	edgs := util.Map(stepDescriptions, func(stepDescription model.NodeDescription) steps.Edge {
		edge, _ := service.NodeTypeService.NewEdgeInstance(runPipelinePayload.PipelineID, runPipelinePayload.RunID, stepDescription.Type, stepDescription)

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

	firstStepID := 0

	for _, step := range stps {
		if step.GetIsFirstStep() {
			firstStepID = step.GetID()
		}
		pipelineGraph.AddVertex(step)
	}

	for _, edge := range edgs {
		pipelineGraph.AddEdge(edge.GetSourceID(), edge.GetTargetID())
	}

	pipelinesWorkDir := os.Getenv("PIPELINES_WORK_DIR")
	currentPipelineWorkDir := pipelinesWorkDir + "/" + fmt.Sprint(runPipelinePayload.PipelineID) + "/" + fmt.Sprint(runPipelinePayload.RunID) + "/"

	if err := os.RemoveAll(currentPipelineWorkDir); err != nil {
		return asynq.SkipRetry
	}

	if err := os.MkdirAll(currentPipelineWorkDir, os.ModePerm); err != nil {
		return asynq.SkipRetry
	}

	logFileName := "run.log"
	logFile, err := os.Create(currentPipelineWorkDir + logFileName)

	if err != nil {
		log.Print(err.Error())
		return asynq.SkipRetry
	}

	graph.BFS(pipelineGraph, firstStepID, func(id int) bool {
		step, _ := pipelineGraph.Vertex(id)
		log.Printf("BFS  %v\n", step)

		if err := step.Execute(logFile); err != nil {
			service.Get(step.GetRunID())
			// service.RunService.Update()
			logFile.Close()
		}

		return false
	})

	return asynq.SkipRetry
}
