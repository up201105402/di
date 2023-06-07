package service

import (
	"di/model"
	"di/repository"
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
	StepTypeService  StepTypeService
	TasksQueueClient asynq.Client
}

func NewRunService(gormDB *gorm.DB, client *asynq.Client, pipelineService *PipelineService, stepTypeService *StepTypeService) RunService {
	return &runServiceImpl{
		RunRepository:    repository.NewRunRepository(gormDB),
		PipelineService:  *pipelineService,
		StepTypeService:  *stepTypeService,
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
	newRun := &model.Run{PipelineID: pipelineId}
	if err := service.RunRepository.Create(newRun); err != nil {
		return err
	}

	runTask, err := tasks.NewRunPipelineTask(newRun.ID, 0)

	if err != nil {
		return err
	}

	if _, err := service.TasksQueueClient.Enqueue(
		runTask,
		asynq.Queue("runs"),
	); err != nil {
		return err
	}

	return nil
}

func (service *runServiceImpl) Execute(runID uint) error {
	// demarshal stringified pipeline definition json

	run, err := service.RunRepository.FindByID(runID)
	pipeline, err := service.PipelineService.Get(run.PipelineID)

	if err != nil {
		log.Printf("Could not retrieve run with id %v\n", runID)
		return err
	}

	// s := "[{\"type\":\"checkoutRepo\",\"dimensions\":{\"width\":135,\"height\":52},\"handleBounds\":{\"target\":[{\"id\":\"0_output\",\"position\":\"right\",\"x\":132.15625,\"y\":23,\"width\":6,\"height\":6}]},\"computedPosition\":{\"x\":-92.47021757726134,\"y\":-27.568760519928844,\"z\":1000},\"selected\":true,\"dragging\":false,\"resizing\":false,\"initialized\":true,\"data\":{\"nameAndType\":{\"nodeName\":\"Step 1\",\"nodeType\":\"checkoutRepo\"},\"stepConfig\":{\"repoURL\":\"http://dafdfds.com/dskfajdsf.git\"},\"isFirstStep\":true},\"events\":{},\"id\":\"0\",\"label\":\"Step 1\",\"position\":{\"x\":-92.47021757726134,\"y\":-27.568760519928844},\"class\":\"light\"}]"
	// var val []Step // <---- This must be an array to match input
	// if err := json.Unmarshal([]byte(s), &val); err != nil {
	// 	panic(err)
	// }
	// fmt.Println(val[0].ID)

	var stepDescriptions []model.StepDescription

	if err := json.Unmarshal([]byte(pipeline.Definition), &stepDescriptions); err != nil {
		log.Printf("Unable to unmarshal pipeline %v definition\n", pipeline.ID)
		return err
	}

	steps := util.Map(stepDescriptions, func(stepDescription model.StepDescription) model.Step {
		step, _ := service.StepTypeService.NewStepInstance(pipeline.ID, stepDescription.Type, stepDescription.Data.StepConfig)
		return *step
	})

	steps = util.Filter(steps, func(step model.Step) bool {
		return step != nil
	})

	pipelineGraph := graph.New(stepHash, graph.Directed())

	for _, step := range steps {

		if step.IsVertex() {
			pipelineGraph.AddVertex(step)
		}

		if step.IsEdge() {
			// TODO
			pipelineGraph.AddEdge()
		}
	}

	pipelineGraph.AddEdge()

	runPipelineTask, err := tasks.NewRunPipelineTask(pipelineGraph, 0)

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

func stepHash(step model.Step) int {
	return graph.IntHash(int(step.GetID()))
}
