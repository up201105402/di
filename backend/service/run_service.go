package service

import (
	"di/model"
	"di/repository"
	"log"

	"github.com/hibiken/asynq"
	"gorm.io/gorm"
)

type runServiceImpl struct {
	RunRepository    model.RunRepository
	PipelineService  PipelineService
	NodeTypeService  NodeTypeService
	TasksQueueClient asynq.Client
	TaskService      TaskService
}

func NewRunService(gormDB *gorm.DB, client *asynq.Client, pipelineService *PipelineService, stepTypeService *NodeTypeService, taskService *TaskService) RunService {
	return &runServiceImpl{
		RunRepository:    repository.NewRunRepository(gormDB),
		PipelineService:  *pipelineService,
		NodeTypeService:  *stepTypeService,
		TasksQueueClient: *client,
		TaskService:      *taskService,
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

	runPipelineTask, err := service.TaskService.NewRunPipelineTask(pipeline.ID, pipeline.Definition, 0)

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
