package service

import (
	"di/model"
	"di/repository"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/hibiken/asynq"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"github.com/robfig/cron/v3"
	"gorm.io/gorm"
)

type pipelineServiceImpl struct {
	PipelineRepository model.PipelineRepository
	TaskQueueClient    *asynq.Client
	I18n               *i18n.Localizer
}

func NewPipelineService(gormDB *gorm.DB, client *asynq.Client, i18n *i18n.Localizer) PipelineService {
	return &pipelineServiceImpl{
		PipelineRepository: repository.NewPipelineRepository(gormDB),
		TaskQueueClient:    client,
		I18n:               i18n,
	}
}

func (service *pipelineServiceImpl) SyncAsyncTasks() {
	redisHost, exists := os.LookupEnv("REDIS_HOST")

	if !exists {
		panic("REDIS_HOST is not defined!")
	}

	redisPort, exists := os.LookupEnv("REDIS_PORT")

	if !exists {
		panic("REDIS_PORT is not defined!")
	}

	inspector := asynq.NewInspector(asynq.RedisClientOpt{Addr: redisHost + ":" + redisPort})
	scheduledTasks, _ := inspector.ListScheduledTasks("runs")
	pipelineSchedules, err := service.GetAllPipelineSchedules()

	if err != nil {
		return
	}

	for i := 0; i < len(pipelineSchedules); i++ {
		isQueued := false

		for j := 0; j < len(scheduledTasks); j++ {
			if scheduledTasks[j].Type == ScheduledRunPipelineTask {
				var scheduledRunPipelinePayload ScheduledRunPipelinePayload
				if err := json.Unmarshal(scheduledTasks[j].Payload, &scheduledRunPipelinePayload); err != nil {
					errStr := fmt.Errorf("json.Unmarshal failed: %v: %w", err, asynq.SkipRetry)
					log.Println(errStr)
					return
				}

				if scheduledRunPipelinePayload.PipelineScheduleID == pipelineSchedules[i].ID {
					isQueued = true
				}
			}
		}

		if !isQueued {

			nextExec := pipelineSchedules[i].UniqueOcurrence

			if pipelineSchedules[i].CronExpression != "" {
				parser := cron.NewParser(cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow)
				schedule, parseError := parser.Parse(pipelineSchedules[i].CronExpression)

				if parseError != nil {
					return
				}

				nextExec = schedule.Next(time.Now())
			}

			if nextExec.Compare(time.Now()) > 0 {
				err := service.enqueueTask(pipelineSchedules[i].UniqueOcurrence, pipelineSchedules[i].CronExpression, pipelineSchedules[i].PipelineID, pipelineSchedules[i].ID)
				if err != nil {
					return
				}
			}
		}
	}
}

func (service *pipelineServiceImpl) Get(id uint) (*model.Pipeline, error) {
	pipeline, err := service.PipelineRepository.FindByID(id)

	if err != nil {
		errMessage := service.I18n.MustLocalize(&i18n.LocalizeConfig{
			MessageID: "pipeline.repository.find.pipeline.id.failed",
			TemplateData: map[string]interface{}{
				"ID":     id,
				"Reason": err.Error(),
			},
			PluralCount: 1,
		})

		return pipeline, errors.New(errMessage)
	}

	return pipeline, err
}

func (service *pipelineServiceImpl) GetPipelineSchedule(id uint) (*model.PipelineSchedule, error) {
	pipelineSchedule, err := service.PipelineRepository.FindPipelineScheduleByID(id)

	if err != nil {
		errMessage := service.I18n.MustLocalize(&i18n.LocalizeConfig{
			MessageID: "pipeline.repository.find.schedule.id.failed",
			TemplateData: map[string]interface{}{
				"ID":     id,
				"Reason": err.Error(),
			},
			PluralCount: 1,
		})

		return pipelineSchedule, errors.New(errMessage)
	}

	return pipelineSchedule, err
}

func (service *pipelineServiceImpl) GetPipelineSchedules(id uint) ([]model.PipelineSchedule, error) {
	schedules, err := service.PipelineRepository.FindPipelineScheduleByPipeline(id)

	if err != nil {
		errMessage := service.I18n.MustLocalize(&i18n.LocalizeConfig{
			MessageID: "pipeline.repository.find.schedule.pipelineID.failed",
			TemplateData: map[string]interface{}{
				"ID":     id,
				"Reason": err.Error(),
			},
			PluralCount: 1,
		})

		return schedules, errors.New(errMessage)
	}

	return schedules, err
}

func (service *pipelineServiceImpl) GetAllPipelineSchedules() ([]model.PipelineSchedule, error) {
	schedules, err := service.PipelineRepository.GetAllPipelineSchedules()

	if err != nil {
		errMessage := service.I18n.MustLocalize(&i18n.LocalizeConfig{
			MessageID: "pipeline.repository.find.schedule.all.failed",
			TemplateData: map[string]interface{}{
				"Reason": err.Error(),
			},
			PluralCount: 1,
		})

		return schedules, errors.New(errMessage)
	}

	return schedules, err
}

func (service *pipelineServiceImpl) GetByOwner(ownerId uint) ([]model.Pipeline, error) {
	pipelines, err := service.PipelineRepository.FindByOwner(ownerId)

	if err != nil {
		errMessage := service.I18n.MustLocalize(&i18n.LocalizeConfig{
			MessageID: "pipeline.repository.find.pipeline.owner.failed",
			TemplateData: map[string]interface{}{
				"OwnerID": ownerId,
				"Reason":  err.Error(),
			},
			PluralCount: 1,
		})

		return pipelines, errors.New(errMessage)
	}

	return pipelines, err
}

func (service *pipelineServiceImpl) Create(userId uint, name string, definition string) error {
	if err := service.PipelineRepository.Create(&model.Pipeline{UserID: userId, Name: name, Definition: definition}); err != nil {
		errMessage := service.I18n.MustLocalize(&i18n.LocalizeConfig{
			MessageID: "pipeline.repository.create.pipeline.failed",
			TemplateData: map[string]interface{}{
				"Reason": err.Error(),
			},
			PluralCount: 1,
		})

		return errors.New(errMessage)
	}

	return nil
}

func (service *pipelineServiceImpl) CreatePipelineSchedule(pipelineID uint, uniqueOcurrence time.Time, cronExpression string) error {
	if cronExpression != "" || uniqueOcurrence.Year() > 1 {

		pipelineSchedule := &model.PipelineSchedule{PipelineID: pipelineID, UniqueOcurrence: uniqueOcurrence, CronExpression: cronExpression}
		if err := service.PipelineRepository.CreatePipelineSchedule(pipelineSchedule); err != nil {
			errMessage := service.I18n.MustLocalize(&i18n.LocalizeConfig{
				MessageID: "pipeline.repository.create.pipeline.failed",
				TemplateData: map[string]interface{}{
					"Reason": err.Error(),
				},
				PluralCount: 1,
			})

			return errors.New(errMessage)
		}

		err := service.enqueueTask(uniqueOcurrence, cronExpression, pipelineID, pipelineSchedule.ID)
		if err != nil {
			return err
		}
	}

	return nil
}

func (service *pipelineServiceImpl) enqueueTask(uniqueOcurrence time.Time, cronExpression string, pipelineID uint, pipelineScheduleID uint) error {
	nextExec := uniqueOcurrence

	if cronExpression != "" {
		parser := cron.NewParser(cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow)
		schedule, parseError := parser.Parse(cronExpression)

		if parseError != nil {
			return parseError
		}

		nextExec = schedule.Next(time.Now())
	}

	task, err := NewScheduledRunPipelineTask(pipelineID, pipelineScheduleID)

	if err != nil {
		return err
	}

	service.TaskQueueClient.Enqueue(task, asynq.Queue("runs"), asynq.Timeout(0), asynq.ProcessAt(nextExec))
	return nil
}

func (service *pipelineServiceImpl) Update(pipeline *model.Pipeline) error {
	err := service.PipelineRepository.Update(pipeline)

	if err != nil {
		errMessage := service.I18n.MustLocalize(&i18n.LocalizeConfig{
			MessageID: "pipeline.repository.create.pipeline.failed",
			TemplateData: map[string]interface{}{
				"ID":     pipeline.ID,
				"Reason": err.Error(),
			},
			PluralCount: 1,
		})

		return errors.New(errMessage)
	}

	return nil
}

func (service *pipelineServiceImpl) Delete(id uint) error {
	err := service.PipelineRepository.Delete(id)

	if err != nil {
		errMessage := service.I18n.MustLocalize(&i18n.LocalizeConfig{
			MessageID: "pipeline.repository.delete.pipeline.failed",
			TemplateData: map[string]interface{}{
				"ID":     id,
				"Reason": err.Error(),
			},
			PluralCount: 1,
		})

		return errors.New(errMessage)
	}

	return nil
}

func (service *pipelineServiceImpl) DeletePipelineSchedule(id uint) error {
	err := service.PipelineRepository.DeletePipelineSchedule(id)

	if err != nil {
		errMessage := service.I18n.MustLocalize(&i18n.LocalizeConfig{
			MessageID: "pipeline.repository.delete.schedule.failed",
			TemplateData: map[string]interface{}{
				"ID":     id,
				"Reason": err.Error(),
			},
			PluralCount: 1,
		})

		return errors.New(errMessage)
	}

	return nil
}

func NewScheduledRunPipelineTask(pipelineID uint, pipelineScheduleID uint) (*asynq.Task, error) {
	payload, err := json.Marshal(ScheduledRunPipelinePayload{PipelineID: pipelineID, PipelineScheduleID: pipelineScheduleID})
	if err != nil {
		return nil, err
	}

	return asynq.NewTask(ScheduledRunPipelineTask, payload, asynq.MaxRetry(0)), nil
}
