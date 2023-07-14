package service

import (
	"di/model"
	"di/repository"
	"encoding/json"
	"time"

	"github.com/hibiken/asynq"
	"github.com/robfig/cron/v3"
	"gorm.io/gorm"
)

type pipelineServiceImpl struct {
	PipelineRepository model.PipelineRepository
	TaskQueueClient    *asynq.Client
}

func NewPipelineService(gormDB *gorm.DB, client *asynq.Client) PipelineService {
	return &pipelineServiceImpl{
		PipelineRepository: repository.NewPipelineRepository(gormDB),
		TaskQueueClient:    client,
	}
}

func (service *pipelineServiceImpl) Get(id uint) (*model.Pipeline, error) {
	pipeline, err := service.PipelineRepository.FindByID(id)
	return pipeline, err
}

func (service *pipelineServiceImpl) GetSchedules(id uint) ([]model.PipelineSchedule, error) {
	schedules, err := service.PipelineRepository.FindScheduleByPipeline(id)
	return schedules, err
}

func (service *pipelineServiceImpl) GetByOwner(ownerId uint) ([]model.Pipeline, error) {
	pipelines, err := service.PipelineRepository.FindByOwner(ownerId)
	return pipelines, err
}

func (service *pipelineServiceImpl) Create(userId uint, name string, definition string) error {
	if err := service.PipelineRepository.Create(&model.Pipeline{UserID: userId, Name: name, Definition: definition}); err != nil {
		return err
	}

	return nil
}

func (service *pipelineServiceImpl) CreateSchedule(pipelineID uint, uniqueOcurrence time.Time, cronExpression string) error {
	if cronExpression != "" || uniqueOcurrence.Year() > 1 {

		if err := service.PipelineRepository.CreateSchedule(&model.PipelineSchedule{PipelineID: pipelineID, UniqueOcurrence: uniqueOcurrence, CronExpression: cronExpression}); err != nil {
			return err
		}

		nextExec := uniqueOcurrence

		if cronExpression != "" {
			parser := cron.NewParser(cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow)
			schedule, parseError := parser.Parse(cronExpression)

			if parseError != nil {
				return parseError
			}

			nextExec = schedule.Next(time.Now())
		}

		task, err := NewScheduledRunPipelineTask(pipelineID, cronExpression)

		if err != nil {
			return err
		}

		service.TaskQueueClient.Enqueue(task, asynq.Queue("runs"), asynq.ProcessAt(nextExec))
	}

	return nil
}

func (service *pipelineServiceImpl) Update(pipeline *model.Pipeline) error {
	err := service.PipelineRepository.Update(pipeline)

	if err != nil {
		return err
	}

	return nil
}

func (service *pipelineServiceImpl) Delete(id uint) error {
	err := service.PipelineRepository.Delete(id)

	if err != nil {
		return err
	}

	return nil
}

func (service *pipelineServiceImpl) DeletePipelineSchedule(id uint) error {
	err := service.PipelineRepository.DeletePipelineSchedule(id)

	if err != nil {
		return err
	}

	return nil
}

func NewScheduledRunPipelineTask(pipelineID uint, cronExpression string) (*asynq.Task, error) {
	payload, err := json.Marshal(ScheduledRunPipelinePayload{PipelineID: pipelineID, CronExpression: cronExpression})
	if err != nil {
		return nil, err
	}

	return asynq.NewTask(RunPipelineTask, payload, asynq.MaxRetry(0)), nil
}
