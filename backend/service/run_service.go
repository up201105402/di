package service

import (
	"di/model"
	"di/repository"
	"encoding/json"
	"log"

	"github.com/hibiken/asynq"
	"gorm.io/gorm"
)

type runServiceImpl struct {
	RunRepository    model.RunRepository
	PipelineService  PipelineService
	TasksQueueClient asynq.Client
}

func NewRunService(gormDB *gorm.DB, client *asynq.Client, pipelineService *PipelineService) RunService {

	return &runServiceImpl{
		RunRepository:    repository.NewRunRepository(gormDB),
		PipelineService:  *pipelineService,
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
	if err := service.RunRepository.Create(&model.Run{PipelineID: pipelineId}); err != nil {
		return err
	}

	return nil
}

func (service *runServiceImpl) Execute(pipelineId uint) error {
	// demarshal stringified pipeline definition json

	pipeline, err := service.PipelineService.Get(pipelineId)

	if err != nil {
		log.Printf("Could not retrieve pipeline with id %v\n", pipelineId)
		return err
	}

	// s := "[{\"type\":\"checkoutRepo\",\"dimensions\":{\"width\":135,\"height\":52},\"handleBounds\":{\"target\":[{\"id\":\"0_output\",\"position\":\"right\",\"x\":132.15625,\"y\":23,\"width\":6,\"height\":6}]},\"computedPosition\":{\"x\":-92.47021757726134,\"y\":-27.568760519928844,\"z\":1000},\"selected\":true,\"dragging\":false,\"resizing\":false,\"initialized\":true,\"data\":{\"nameAndType\":{\"nodeName\":\"Step 1\",\"nodeType\":\"checkoutRepo\"},\"stepConfig\":{\"repoURL\":\"http://dafdfds.com/dskfajdsf.git\"},\"isFirstStep\":true},\"events\":{},\"id\":\"0\",\"label\":\"Step 1\",\"position\":{\"x\":-92.47021757726134,\"y\":-27.568760519928844},\"class\":\"light\"}]"
	// var val []Step // <---- This must be an array to match input
	// if err := json.Unmarshal([]byte(s), &val); err != nil {
	// 	panic(err)
	// }
	// fmt.Println(val[0].ID)

	var steps []model.StepDescription

	if err := json.Unmarshal([]byte(pipeline.Definition), &steps); err != nil {
		log.Printf("Unable to unmarshal pipeline %v definition\n", pipelineId)
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
