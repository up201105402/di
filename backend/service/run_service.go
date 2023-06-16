package service

import (
	"di/model"
	"di/repository"
	"di/steps"
	"di/tasks"
	"di/util"
	"encoding/json"
	"log"

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

	var stepDescriptions []model.NodeDescription

	if err := json.Unmarshal([]byte(pipeline.Definition), &stepDescriptions); err != nil {
		log.Printf("Unable to unmarshal pipeline %v definition\n", pipeline.ID)
		return err
	}

	pipelineGraph := graph.New(stepHash, graph.Directed())

	// var stps []steps.Step
	// var edges []steps.Edge

	// for index, stepDescription := range stepDescriptions {
	// 	step, _ := service.NodeTypeService.NewStepInstance(pipeline.ID, stepDescription.Type, stepDescription.Data.StepConfig)
	// 	edge, _ := service.NodeTypeService.NewEdgeInstance(pipeline.ID, stepDescription.Type, stepDescription.Data.StepConfig)

	// 	if step != nil {
	// 		stps = append(stps, *step)
	// 	}

	// 	if edge != nil {
	// 		edges = append(edges, *edge)
	// 	}
	// }

	stps := util.Map(stepDescriptions, func(stepDescription model.NodeDescription) steps.Step {
		step, _ := service.NodeTypeService.NewStepInstance(pipeline.ID, stepDescription.Type, stepDescription.Data.StepConfig)
		return *step
	})

	edgs := util.Map(stepDescriptions, func(stepDescription model.NodeDescription) steps.Edge {
		edge, _ := service.NodeTypeService.NewEdgeInstance(pipeline.ID, stepDescription.Type, stepDescription.Data.StepConfig)
		return *edge
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

	runPipelineTask, err := tasks.NewRunPipelineTask(pipelineGraph, 0)

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

func stepHash(step steps.Step) int {
	return graph.IntHash(int(step.GetID()))
}
