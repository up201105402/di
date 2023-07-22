package service

import (
	"context"
	"di/model"
	"di/repository"
	"di/steps"
	"di/util"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/dominikbraun/graph"
	"github.com/hibiken/asynq"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"github.com/robfig/cron/v3"
	"gopkg.in/guregu/null.v4"
	"gorm.io/gorm"
)

type runServiceImpl struct {
	RunRepository   model.RunRepository
	PipelineService PipelineService
	NodeTypeService NodeTypeService
	TaskQueueClient asynq.Client
	I18n            *i18n.Localizer
}

func NewRunService(gormDB *gorm.DB, client *asynq.Client, i18n *i18n.Localizer, pipelineService *PipelineService, stepTypeService *NodeTypeService) RunService {
	return &runServiceImpl{
		RunRepository:   repository.NewRunRepository(gormDB),
		PipelineService: *pipelineService,
		NodeTypeService: *stepTypeService,
		TaskQueueClient: *client,
		I18n:            i18n,
	}
}

func (service *runServiceImpl) Get(id uint) (*model.Run, error) {
	run, err := service.RunRepository.FindByID(id)

	if err != nil {
		errMessage := service.I18n.MustLocalize(&i18n.LocalizeConfig{
			MessageID: "run.repository.find.run.id.failed",
			TemplateData: map[string]interface{}{
				"ID":     id,
				"Reason": err.Error(),
			},
			PluralCount: 1,
		})

		return run, errors.New(errMessage)
	}

	return run, err
}

func (service *runServiceImpl) GetByPipeline(pipelineId uint) ([]model.Run, error) {
	runs, err := service.RunRepository.FindByPipeline(pipelineId)

	if err != nil {
		errMessage := service.I18n.MustLocalize(&i18n.LocalizeConfig{
			MessageID: "run.repository.find.run.pipeline.failed",
			TemplateData: map[string]interface{}{
				"ID":     pipelineId,
				"Reason": err.Error(),
			},
			PluralCount: 1,
		})

		return runs, errors.New(errMessage)
	}

	return runs, err
}

func (service *runServiceImpl) FindRunStepStatusesByRun(runID uint) ([]model.RunStepStatus, error) {
	runStepStatuses, err := service.RunRepository.FindRunStepStatusesByRun(runID)

	if err != nil {
		errMessage := service.I18n.MustLocalize(&i18n.LocalizeConfig{
			MessageID: "run.repository.find.step-status.run.failed",
			TemplateData: map[string]interface{}{
				"ID":     runID,
				"Reason": err.Error(),
			},
			PluralCount: 1,
		})

		return runStepStatuses, errors.New(errMessage)
	}

	return runStepStatuses, err
}

func (service *runServiceImpl) FindHumanFeedbackQueriesByStepID(stepID uint) ([]model.HumanFeedbackQuery, error) {
	runStepStatuses, err := service.RunRepository.FindHumanFeedbackQueriesByStepID(stepID)

	if err != nil {
		errMessage := service.I18n.MustLocalize(&i18n.LocalizeConfig{
			MessageID: "run.repository.find.human-feedback-query.step.failed",
			TemplateData: map[string]interface{}{
				"ID":     stepID,
				"Reason": err.Error(),
			},
			PluralCount: 1,
		})

		return runStepStatuses, errors.New(errMessage)
	}

	return runStepStatuses, err
}

func (service *runServiceImpl) FindHumanFeedbackRectsByHumanFeedbackQueryID(humanFeedbackQueryID uint) ([]model.HumanFeedbackRect, error) {
	feedbackRects, err := service.RunRepository.FindHumanFeedbackRectsByHumanFeedbackQueryID(humanFeedbackQueryID)

	if err != nil {
		errMessage := service.I18n.MustLocalize(&i18n.LocalizeConfig{
			MessageID: "run.repository.find.human-feedback-rects.feedback-query.failed",
			TemplateData: map[string]interface{}{
				"ID":     humanFeedbackQueryID,
				"Reason": err.Error(),
			},
			PluralCount: 1,
		})

		return feedbackRects, errors.New(errMessage)
	}

	return feedbackRects, err
}

func (service *runServiceImpl) FindHumanFeedbackQueryStatusByID(queryStatusID uint) (*model.QueryStatus, error) {
	queryStatus, err := service.RunRepository.FindHumanFeedbackQueryStatusByID(queryStatusID)

	if err != nil {
		errMessage := service.I18n.MustLocalize(&i18n.LocalizeConfig{
			MessageID: "run.repository.find.query-status.failed",
			TemplateData: map[string]interface{}{
				"ID":     queryStatusID,
				"Reason": err.Error(),
			},
			PluralCount: 1,
		})

		return nil, errors.New(errMessage)
	}

	return queryStatus, err
}

func (service *runServiceImpl) Create(pipeline model.Pipeline) (model.Run, error) {
	newRun := &model.Run{PipelineID: pipeline.ID, RunStatusID: 1, Definition: pipeline.Definition}
	if err := service.RunRepository.Create(newRun); err != nil {
		errMessage := service.I18n.MustLocalize(&i18n.LocalizeConfig{
			MessageID: "run.repository.create.run.failed",
			TemplateData: map[string]interface{}{
				"Reason": err.Error(),
			},
			PluralCount: 1,
		})

		return *newRun, errors.New(errMessage)
	}

	return *newRun, nil
}

func (service *runServiceImpl) CreateRunStepStatus(runID uint, stepID int, stepName string, runStatusID uint, errorMessage string) error {
	newRunStepStatus := &model.RunStepStatus{RunID: runID, StepID: stepID, RunStatusID: runStatusID, LastRun: time.Now()}
	if err := service.RunRepository.CreateRunStepStatus(newRunStepStatus); err != nil {
		errMessage := service.I18n.MustLocalize(&i18n.LocalizeConfig{
			MessageID: "run.repository.create.step-status.failed",
			TemplateData: map[string]interface{}{
				"Reason": err.Error(),
			},
			PluralCount: 1,
		})

		return errors.New(errMessage)
	}

	return nil
}

func (service *runServiceImpl) CreateHumanFeedbackQuery(epoch uint, runID uint, stepID int, queryID uint, rects [][]uint) error {
	newHumandFeedbackQuery := &model.HumanFeedbackQuery{
		Epoch:         epoch,
		StepID:        stepID,
		QueryID:       queryID,
		RunID:         runID,
		QueryStatusID: 1, // unresolved
	}

	if err := service.RunRepository.CreateHumanFeedbackQuery(newHumandFeedbackQuery); err != nil {
		errMessage := service.I18n.MustLocalize(&i18n.LocalizeConfig{
			MessageID: "run.repository.create.human-feedback-query.failed",
			TemplateData: map[string]interface{}{
				"Reason": err.Error(),
			},
			PluralCount: 1,
		})

		return errors.New(errMessage)
	}

	for _, rect := range rects {
		humanFeedbackRect := &model.HumanFeedbackRect{
			X1:                   rect[0],
			Y1:                   rect[1],
			X2:                   rect[2],
			Y2:                   rect[3],
			HumanFeedbackQueryID: newHumandFeedbackQuery.ID,
			Selected:             false,
		}

		if err := service.RunRepository.CreateHumanFeedbackRect(humanFeedbackRect); err != nil {
			errMessage := service.I18n.MustLocalize(&i18n.LocalizeConfig{
				MessageID: "run.repository.create.human-feedback-rect.failed",
				TemplateData: map[string]interface{}{
					"Reason": err.Error(),
				},
				PluralCount: 1,
			})

			return errors.New(errMessage)
		}
	}

	return nil
}

func (service *runServiceImpl) Execute(runID uint) error {

	run, err := service.RunRepository.FindByID(runID)

	if err != nil {
		errMessage := service.I18n.MustLocalize(&i18n.LocalizeConfig{
			MessageID: "run.repository.find.run.id.failed",
			TemplateData: map[string]interface{}{
				"ID":     runID,
				"Reason": err.Error(),
			},
			PluralCount: 1,
		})

		return errors.New(errMessage)
	}

	if err := service.UpdateRunStatus(runID, 2, 0, ""); err != nil {
		return err
	}

	err = service.RunRepository.DeleteAllRunStepStatuses(runID)

	if err != nil {
		errMessage := service.I18n.MustLocalize(&i18n.LocalizeConfig{
			MessageID: "run.repository.delete.step-status.all.failed",
			TemplateData: map[string]interface{}{
				"ID":     runID,
				"Reason": err.Error(),
			},
			PluralCount: 1,
		})

		return errors.New(errMessage)
	}

	pipeline, err := service.PipelineService.Get(run.PipelineID)

	if err != nil {
		return err
	}

	runPipelineTask, err := service.NewRunPipelineTask(pipeline.ID, runID, pipeline.Definition)

	if err != nil {
		return err
	}

	if _, err := service.TaskQueueClient.Enqueue(runPipelineTask, asynq.Queue("runs"), asynq.Timeout(0)); err != nil {
		errMessage := service.I18n.MustLocalize(&i18n.LocalizeConfig{
			MessageID: "tasks.client.enqueue.failed",
			TemplateData: map[string]interface{}{
				"Queue":  "runs",
				"Reason": err.Error(),
			},
			PluralCount: 1,
		})

		return errors.New(errMessage)
	}

	return nil
}

func (service *runServiceImpl) Resume(runID uint) error {
	run, err := service.RunRepository.FindByID(runID)

	if err != nil {
		log.Printf(err.Error())
		return err
	}

	if run.RunStatusID != 5 {
		errMessage := service.I18n.MustLocalize(&i18n.LocalizeConfig{
			MessageID: "run.service.resume.status.error",
			TemplateData: map[string]interface{}{
				"ID": run.ID,
			},
			PluralCount: 1,
		})
		log.Printf(errMessage)
		return errors.New(errMessage)
	}

	runStepStatuses, getError := service.FindRunStepStatusesByRun(run.ID)

	if getError != nil {
		log.Printf(getError.Error())
		return err
	}

	runStepStatuses = util.Filter(runStepStatuses, func(runStateStatus model.RunStepStatus) bool {
		return runStateStatus.RunStatusID == 5
	})

	if len(runStepStatuses) == 0 {
		errMessage := service.I18n.MustLocalize(&i18n.LocalizeConfig{
			MessageID: "run.service.resume.status.error",
			TemplateData: map[string]interface{}{
				"ID": run.ID,
			},
			PluralCount: 1,
		})
		log.Printf(errMessage)
		return errors.New(errMessage)
	}

	humanFeedbackQueries, err := service.FindHumanFeedbackQueriesByStepID(uint(runStepStatuses[0].StepID))

	if err != nil {
		errMessage := service.I18n.MustLocalize(&i18n.LocalizeConfig{
			MessageID: "run.service.resume.status.error",
			TemplateData: map[string]interface{}{
				"ID": run.ID,
			},
			PluralCount: 1,
		})
		log.Printf(errMessage)
		return errors.New(errMessage)
	}

	if len(humanFeedbackQueries) == 0 {
		errMessage := service.I18n.MustLocalize(&i18n.LocalizeConfig{
			MessageID: "run.service.resume.status.error",
			TemplateData: map[string]interface{}{
				"ID": run.ID,
			},
			PluralCount: 1,
		})
		log.Printf(errMessage)
		return errors.New(errMessage)
	}

	runPipelineTask, err := service.NewResumeRunPipelineTask(run.Pipeline.ID, runID, run.Definition, runStepStatuses[0].StepID)

	if err != nil {
		return err
	}

	if _, err := service.TaskQueueClient.Enqueue(runPipelineTask, asynq.Queue("runs"), asynq.Timeout(0)); err != nil {
		errMessage := service.I18n.MustLocalize(&i18n.LocalizeConfig{
			MessageID: "tasks.client.enqueue.failed",
			TemplateData: map[string]interface{}{
				"Queue":  "runs",
				"Reason": err.Error(),
			},
			PluralCount: 1,
		})

		return errors.New(errMessage)
	}

	return nil
}

func (service *runServiceImpl) Update(run *model.Run) error {
	err := service.RunRepository.Update(run)

	if err != nil {
		errMessage := service.I18n.MustLocalize(&i18n.LocalizeConfig{
			MessageID: "run.repository.update.run.failed",
			TemplateData: map[string]interface{}{
				"ID":     run.ID,
				"Reason": err.Error(),
			},
			PluralCount: 1,
		})

		return errors.New(errMessage)
	}

	return nil
}

func (service *runServiceImpl) UpdateRunStepStatus(runStepStatus *model.RunStepStatus) error {
	err := service.RunRepository.UpdateRunStepStatus(runStepStatus)

	if err != nil {
		errMessage := service.I18n.MustLocalize(&i18n.LocalizeConfig{
			MessageID: "run.repository.update.step-status.failed",
			TemplateData: map[string]interface{}{
				"ID":     runStepStatus.ID,
				"Reason": err.Error(),
			},
			PluralCount: 1,
		})

		return errors.New(errMessage)
	}

	return nil
}

func (service *runServiceImpl) UpdateHumanFeedbackQuery(query *model.HumanFeedbackQuery) error {

	err := service.RunRepository.UpdateHumanFeedbackQuery(query)

	if err != nil {
		errMessage := service.I18n.MustLocalize(&i18n.LocalizeConfig{
			MessageID: "run.repository.update.human-feedback-query.failed",
			TemplateData: map[string]interface{}{
				"ID":     query.ID,
				"Reason": err.Error(),
			},
			PluralCount: 1,
		})

		return errors.New(errMessage)
	}

	return nil
}

func (service *runServiceImpl) UpdateHumanFeedbackRects(rects []model.HumanFeedbackRect) error {

	for _, rect := range rects {
		err := service.RunRepository.UpdateHumanFeedbackRect(&rect)

		if err != nil {
			errMessage := service.I18n.MustLocalize(&i18n.LocalizeConfig{
				MessageID: "run.repository.update.human-feedback-rects.failed",
				TemplateData: map[string]interface{}{
					"ID":     rect.HumanFeedbackQueryID,
					"Reason": err.Error(),
				},
				PluralCount: 1,
			})

			return errors.New(errMessage)
		}
	}

	return nil
}

func (service *runServiceImpl) Delete(id uint) error {
	err := service.RunRepository.Delete(id)

	if err != nil {
		errMessage := service.I18n.MustLocalize(&i18n.LocalizeConfig{
			MessageID: "run.repository.delete.run.failed",
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

func (service *runServiceImpl) DeleteRunStepStatus(id uint) error {
	err := service.RunRepository.DeleteRunStepStatus(id)

	if err != nil {
		errMessage := service.I18n.MustLocalize(&i18n.LocalizeConfig{
			MessageID: "run.repository.delete.step-status.id.failed",
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

func (service *runServiceImpl) NewRunPipelineTask(pipelineID uint, runID uint, graph string) (*asynq.Task, error) {
	payload, err := json.Marshal(RunPipelinePayload{PipelineID: pipelineID, RunID: runID, GraphDefinition: graph})
	if err != nil {
		return nil, err
	}
	return asynq.NewTask(RunPipelineTask, payload, asynq.MaxRetry(0)), nil
}

func (service *runServiceImpl) NewResumeRunPipelineTask(pipelineID uint, runID uint, graph string, stepID int) (*asynq.Task, error) {
	payload, err := json.Marshal(RunPipelinePayload{PipelineID: pipelineID, RunID: runID, GraphDefinition: graph, StepID: null.NewInt(int64(stepID), true)})
	if err != nil {
		return nil, err
	}
	return asynq.NewTask(RunPipelineTask, payload, asynq.MaxRetry(0)), nil
}

func (service *runServiceImpl) HandleRunPipelineTask(ctx context.Context, t *asynq.Task) error {
	var runPipelinePayload RunPipelinePayload
	if err := json.Unmarshal(t.Payload(), &runPipelinePayload); err != nil {
		errStr := fmt.Errorf("json.Unmarshal failed: %v", err.Error())
		log.Println(errStr)
		return asynq.SkipRetry
	}

	if runPipelinePayload.StepID.Valid {
		return service.resumeRunPipelineTask(runPipelinePayload)
	}

	return executeRunPipelineTask(runPipelinePayload, service)
}

func executeRunPipelineTask(runPipelinePayload RunPipelinePayload, service *runServiceImpl) error {

	pipelineGraph, err := service.createPipelineGraph(runPipelinePayload)

	if err != nil {
		log.Printf(err.Error())

		if err := service.UpdateRunStatus(runPipelinePayload.RunID, 3, 0, err.Error()); err != nil {
			log.Println(err.Error())
		}

		return asynq.SkipRetry
	}

	pipelinesWorkDir, exists := os.LookupEnv("PIPELINES_WORK_DIR")

	if !exists {
		errMessage := service.I18n.MustLocalize(&i18n.LocalizeConfig{
			MessageID: "env.variable.find.failed",
			TemplateData: map[string]interface{}{
				"Name": "PIPELINES_WORK_DIR",
			},
			PluralCount: 1,
		})

		log.Println(errMessage)

		if err := service.UpdateRunStatus(runPipelinePayload.RunID, 3, 0, errMessage); err != nil {
			log.Println(err.Error())
		}

		return asynq.SkipRetry
	}

	runLogsDir, exists := os.LookupEnv("RUN_LOGS_DIR")

	if !exists {
		errMessage := service.I18n.MustLocalize(&i18n.LocalizeConfig{
			MessageID: "env.variable.find.failed",
			TemplateData: map[string]interface{}{
				"Name": "RUN_LOGS_DIR",
			},
			PluralCount: 1,
		})

		log.Println(errMessage)

		if err := service.UpdateRunStatus(runPipelinePayload.RunID, 3, 0, errMessage); err != nil {
			log.Println(err.Error())
		}

		return asynq.SkipRetry
	}

	currentPipelineWorkDir := pipelinesWorkDir + "/" + fmt.Sprint(runPipelinePayload.PipelineID) + "/" + fmt.Sprint(runPipelinePayload.RunID) + "/"
	currentRunLogDir := runLogsDir + "/pipelines/" + fmt.Sprint(runPipelinePayload.PipelineID) + "/" + fmt.Sprint(runPipelinePayload.RunID) + "/"

	if err := os.RemoveAll(currentPipelineWorkDir); err != nil {
		errMessage := service.I18n.MustLocalize(&i18n.LocalizeConfig{
			MessageID: "env.variable.find.failed",
			TemplateData: map[string]interface{}{
				"Path":   currentPipelineWorkDir,
				"Reason": err.Error(),
			},
			PluralCount: 2,
		})

		log.Println(errMessage)

		if err := service.UpdateRunStatus(runPipelinePayload.RunID, 3, 0, errMessage); err != nil {
			log.Println(err.Error())
		}

		return asynq.SkipRetry
	}

	if err := os.MkdirAll(currentPipelineWorkDir, os.ModePerm); err != nil {
		errMessage := service.I18n.MustLocalize(&i18n.LocalizeConfig{
			MessageID: "os.cmd.mkdir.dir.failed",
			TemplateData: map[string]interface{}{
				"Path":   currentPipelineWorkDir,
				"Reason": err.Error(),
			},
			PluralCount: 1,
		})

		log.Println(errMessage)

		if err := service.UpdateRunStatus(runPipelinePayload.RunID, 3, 0, errMessage); err != nil {
			log.Println(err.Error())
		}

		return asynq.SkipRetry
	}

	if err := os.RemoveAll(currentRunLogDir); err != nil {
		errMessage := service.I18n.MustLocalize(&i18n.LocalizeConfig{
			MessageID: "env.variable.find.failed",
			TemplateData: map[string]interface{}{
				"Path":   currentRunLogDir,
				"Reason": err.Error(),
			},
			PluralCount: 1,
		})

		log.Println(errMessage)

		if err := service.UpdateRunStatus(runPipelinePayload.RunID, 3, 0, errMessage); err != nil {
			log.Println(err.Error())
		}

		return asynq.SkipRetry
	}

	if err := os.MkdirAll(currentRunLogDir, os.ModePerm); err != nil {
		errMessage := service.I18n.MustLocalize(&i18n.LocalizeConfig{
			MessageID: "os.cmd.mkdir.dir.failed",
			TemplateData: map[string]interface{}{
				"Path":   currentRunLogDir,
				"Reason": err.Error(),
			},
			PluralCount: 1,
		})

		log.Println(errMessage)

		if err := service.UpdateRunStatus(runPipelinePayload.RunID, 3, 0, errMessage); err != nil {
			log.Println(err.Error())
		}

		return asynq.SkipRetry
	}

	logFileName, exists := os.LookupEnv("RUN_LOG_FILE_NAME")

	if !exists {
		errMessage := service.I18n.MustLocalize(&i18n.LocalizeConfig{
			MessageID: "env.variable.find.failed",
			TemplateData: map[string]interface{}{
				"Name": "RUN_LOG_FILE_NAME",
			},
			PluralCount: 1,
		})

		service.UpdateRunStatus(runPipelinePayload.RunID, 3, 0, errMessage)
		log.Println(errors.New(errMessage))
		return asynq.SkipRetry
	}

	logFile, err := os.Create(currentRunLogDir + logFileName)

	if err != nil {
		errMessage := service.I18n.MustLocalize(&i18n.LocalizeConfig{
			MessageID: "os.cmd.create.file.failed",
			TemplateData: map[string]interface{}{
				"Name":   currentRunLogDir + logFileName,
				"Reason": err.Error(),
			},
			PluralCount: 1,
		})

		log.Print(errMessage)
		service.UpdateRunStatus(runPipelinePayload.RunID, 3, 0, errMessage)
		return asynq.SkipRetry
	}

	return service.traverseAndExecuteSteps(runPipelinePayload.RunID, pipelineGraph, 0, logFile)
}

func (service *runServiceImpl) resumeRunPipelineTask(runPipelinePayload RunPipelinePayload) error {

	pipelineGraph, err := service.createPipelineGraph(runPipelinePayload)

	if err != nil {
		log.Printf(err.Error())

		if err := service.UpdateRunStatus(runPipelinePayload.RunID, 3, 0, err.Error()); err != nil {
			log.Println(err.Error())
		}

		return asynq.SkipRetry
	}

	pipelinesWorkDir, exists := os.LookupEnv("PIPELINES_WORK_DIR")

	if !exists {
		errMessage := service.I18n.MustLocalize(&i18n.LocalizeConfig{
			MessageID: "env.variable.find.failed",
			TemplateData: map[string]interface{}{
				"Name": "PIPELINES_WORK_DIR",
			},
			PluralCount: 1,
		})

		log.Println(errMessage)

		if err := service.UpdateRunStatus(runPipelinePayload.RunID, 3, 0, errMessage); err != nil {
			log.Println(err.Error())
		}

		return asynq.SkipRetry
	}

	runLogsDir, exists := os.LookupEnv("RUN_LOGS_DIR")

	if !exists {
		errMessage := service.I18n.MustLocalize(&i18n.LocalizeConfig{
			MessageID: "env.variable.find.failed",
			TemplateData: map[string]interface{}{
				"Name": "RUN_LOGS_DIR",
			},
			PluralCount: 1,
		})

		log.Println(errMessage)

		if err := service.UpdateRunStatus(runPipelinePayload.RunID, 3, 0, errMessage); err != nil {
			log.Println(err.Error())
		}

		return asynq.SkipRetry
	}

	currentPipelineWorkDir := pipelinesWorkDir + "/" + fmt.Sprint(runPipelinePayload.PipelineID) + "/" + fmt.Sprint(runPipelinePayload.RunID) + "/"
	currentRunLogDir := runLogsDir + "/pipelines/" + fmt.Sprint(runPipelinePayload.PipelineID) + "/" + fmt.Sprint(runPipelinePayload.RunID) + "/"
	logFileName, exists := os.LookupEnv("RUN_LOG_FILE_NAME")

	if !exists {
		errMessage := service.I18n.MustLocalize(&i18n.LocalizeConfig{
			MessageID: "env.variable.find.failed",
			TemplateData: map[string]interface{}{
				"Name": "RUN_LOG_FILE_NAME",
			},
			PluralCount: 1,
		})

		service.UpdateRunStatus(runPipelinePayload.RunID, 3, 0, errMessage)
		log.Println(errors.New(errMessage))
		return asynq.SkipRetry
	}

	logFile, err := os.Open(currentRunLogDir + logFileName)

	if err != nil {
		errMessage := service.I18n.MustLocalize(&i18n.LocalizeConfig{
			MessageID: "os.cmd.read.file.failed",
			TemplateData: map[string]interface{}{
				"Name":   currentRunLogDir + logFileName,
				"Reason": err.Error(),
			},
			PluralCount: 1,
		})

		log.Print(errMessage)
		service.UpdateRunStatus(runPipelinePayload.RunID, 3, 0, errMessage)
		return asynq.SkipRetry
	}

	_, err = os.Stat(currentPipelineWorkDir)

	if err != nil {
		errMessage := service.I18n.MustLocalize(&i18n.LocalizeConfig{
			MessageID: "os.cmd.read.dir.failed",
			TemplateData: map[string]interface{}{
				"Name":   currentPipelineWorkDir,
				"Reason": err.Error(),
			},
			PluralCount: 1,
		})

		log.Print(errMessage)
		service.UpdateRunStatus(runPipelinePayload.RunID, 3, 0, errMessage)
		return asynq.SkipRetry
	}

	return service.traverseAndExecuteSteps(runPipelinePayload.RunID, pipelineGraph, int(runPipelinePayload.StepID.Int64), logFile)
}

func (service *runServiceImpl) traverseAndExecuteSteps(runID uint, pipelineGraph graph.Graph[int, steps.Step], startAtStepID int, logFile *os.File) error {

	runLogger := log.New(logFile, "", log.Ldate|log.Ltime|log.Lmicroseconds|log.Llongfile)

	msg := fmt.Sprintf("Executing run %d", runID)
	log.Println(msg)
	runLogger.Println(msg)

	hasError := false
	var stepErr error
	hasFeedback := false
	stepWaitingFeedback := 0

	graph.BFS(pipelineGraph, int(startAtStepID), func(id int) bool {
		step, _ := pipelineGraph.Vertex(id)

		msg := fmt.Sprintf("Executing step %s (%d) ...", step.GetName(), step.GetID())
		log.Println(msg)
		runLogger.Println(msg)

		var feebackRects [][]model.HumanFeedbackRect

		run, _ := service.Get(uint(runID))

		var runStepStatus *model.RunStepStatus

		if step.GetIsStaggered() && run.RunStatusID == 5 {
			runStepStatuses, getError := service.FindRunStepStatusesByRun(run.ID)

			if getError != nil {
				log.Printf(getError.Error())
				hasError = true
				stepErr = getError
				return true
			}

			runStepStatuses = util.Filter(runStepStatuses, func(runStateStatus model.RunStepStatus) bool {
				return runStateStatus.RunStatusID == 5 && runStateStatus.StepID == startAtStepID
			})

			if len(runStepStatuses) == 0 {
				errMessage := service.I18n.MustLocalize(&i18n.LocalizeConfig{
					MessageID: "run.service.resume.steps.status.error",
					TemplateData: map[string]interface{}{
						"ID": run.ID,
					},
					PluralCount: 1,
				})
				log.Printf(errMessage)
				hasError = true
				stepErr = errors.New(errMessage)
				return true
			}

			runStepStatus = &runStepStatuses[0]

			humanFeedbackQueries, err := service.FindHumanFeedbackQueriesByStepID(uint(runStepStatus.StepID))

			if err != nil {
				errMessage := service.I18n.MustLocalize(&i18n.LocalizeConfig{
					MessageID: "run.handler.feedback.find.fail",
					TemplateData: map[string]interface{}{
						"ID":     runStepStatuses[0].StepID,
						"Reason": err.Error(),
					},
					PluralCount: 1,
				})
				log.Printf(errMessage)
				hasError = true
				stepErr = errors.New(errMessage)

				if err := service.updateStepRunStatus(runStepStatus, 3, err.Error()); err != nil {
					log.Println(err.Error())
					runLogger.Println(err.Error())
				}

				return true
			}

			for _, humanFeedbackQuery := range humanFeedbackQueries {
				rects, err := service.FindHumanFeedbackRectsByHumanFeedbackQueryID(humanFeedbackQuery.ID)

				if err != nil {
					log.Printf(err.Error())
					hasError = true
					stepErr = errors.New(err.Error())
					return true
				}

				feebackRects = append(feebackRects, rects)
			}
		} else {
			runStepStatus := &model.RunStepStatus{RunID: runID, StepID: id, Name: step.GetName(), RunStatusID: 2, LastRun: time.Now()}
			err := service.RunRepository.CreateRunStepStatus(runStepStatus)

			if err != nil {
				log.Println(err.Error())
				runLogger.Println(err.Error())

				if err := service.updateStepRunStatus(runStepStatus, 3, err.Error()); err != nil {
					log.Println(err.Error())
					runLogger.Println(err.Error())
				}

				hasError = true
				return true
			}
		}

		if feedbackPayload, err := step.Execute(logFile, feebackRects, service.I18n); err != nil {

			stepErr = err
			hasError = true

			errMessage := service.I18n.MustLocalize(&i18n.LocalizeConfig{
				MessageID: "run.service.execute.step.failed",
				TemplateData: map[string]interface{}{
					"ID":     step.GetID(),
					"Reason": err.Error(),
				},
				PluralCount: 1,
			})

			runLogger.Println(errMessage)
			log.Println(errMessage)

			if err := service.updateStepRunStatus(runStepStatus, 3, err.Error()); err != nil {
				runLogger.Println(errMessage)
				log.Println(errMessage)
			}

			return true
		} else {
			for _, feedback := range feedbackPayload {
				err = service.CreateHumanFeedbackQuery(feedback.Epoch, feedback.RunID, feedback.StepID, feedback.QueryID, feedback.Rects)
				hasFeedback = true
			}

			if err != nil {
				errMessage := service.I18n.MustLocalize(&i18n.LocalizeConfig{
					MessageID: "run.service.execute.step.failed",
					TemplateData: map[string]interface{}{
						"ID":     step.GetID(),
						"Reason": err.Error(),
					},
					PluralCount: 1,
				})

				runLogger.Println(errMessage)
				log.Println(errMessage)

				if err := service.updateStepRunStatus(runStepStatus, 3, err.Error()); err != nil {
					runLogger.Println(errMessage)
					log.Println(errMessage)
				}
			} else {
				errMessage := service.I18n.MustLocalize(&i18n.LocalizeConfig{
					MessageID: "run.service.execute.step.success",
					TemplateData: map[string]interface{}{
						"ID": step.GetID(),
					},
					PluralCount: 1,
				})

				runLogger.Println(errMessage)
				log.Println(errMessage)
			}
		}

		var err error

		if hasFeedback {
			err = service.updateStepRunStatus(runStepStatus, 5, "") // waiting feedback
		} else {
			err = service.updateStepRunStatus(runStepStatus, 4, "") // success
		}

		if err != nil {
			runLogger.Println(err.Error())
			log.Println(err.Error())
			stepErr = err
			hasError = true
			return true
		}

		if hasFeedback {
			stepWaitingFeedback = step.GetID()
			return true
		}

		return false
	})

	if hasError {
		if err := service.UpdateRunStatus(runID, 3, 0, msg); err != nil {
			log.Println(err.Error())
			runLogger.Println(err.Error())
		}

		errMessage := service.I18n.MustLocalize(&i18n.LocalizeConfig{
			MessageID: "run.service.execute.run.failed",
			TemplateData: map[string]interface{}{
				"ID":     runID,
				"Reason": stepErr.Error(),
			},
			PluralCount: 1,
		})

		log.Println(errMessage)
		runLogger.Println(errMessage)
	} else {
		runStatusID := 4
		if hasFeedback {
			runStatusID = 5
		}
		if err := service.UpdateRunStatus(runID, uint(runStatusID), stepWaitingFeedback, ""); err != nil {
			log.Println(err.Error())
			runLogger.Println(err.Error())
		}

		errMessage := service.I18n.MustLocalize(&i18n.LocalizeConfig{
			MessageID: "run.service.execute.run.success",
			TemplateData: map[string]interface{}{
				"ID": runID,
			},
			PluralCount: 1,
		})

		log.Println(errMessage)
		runLogger.Println(errMessage)
	}

	logFile.Close()

	return asynq.SkipRetry
}

func (service *runServiceImpl) createPipelineGraph(runPipelinePayload RunPipelinePayload) (graph.Graph[int, steps.Step], error) {
	var stepDescriptions []model.NodeDescription

	if err := json.Unmarshal([]byte(runPipelinePayload.GraphDefinition), &stepDescriptions); err != nil {
		errMessage := service.I18n.MustLocalize(&i18n.LocalizeConfig{
			MessageID: "run.service.execute.demarshal.error",
			TemplateData: map[string]interface{}{
				"ID":     runPipelinePayload.RunID,
				"Reason": err.Error(),
			},
			PluralCount: 1,
		})
		log.Printf(errMessage)
		return nil, errors.New(errMessage)
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

	for _, step := range stps {
		pipelineGraph.AddVertex(step)
	}

	for _, edge := range edgs {
		pipelineGraph.AddEdge(edge.GetTargetID(), edge.GetSourceID())
	}

	return pipelineGraph, nil
}

func (service *runServiceImpl) HandleScheduledRunPipelineTask(ctx context.Context, t *asynq.Task) error {
	var scheduledRunPipelinePayload ScheduledRunPipelinePayload
	if err := json.Unmarshal(t.Payload(), &scheduledRunPipelinePayload); err != nil {
		errStr := fmt.Errorf("json.Unmarshal failed: %v: %w", err, asynq.SkipRetry)
		log.Println(errStr)
		return errStr
	}

	pipelineSchedule, err := service.PipelineService.GetPipelineSchedule(scheduledRunPipelinePayload.PipelineScheduleID)

	if err != nil {
		log.Println(err.Error())
		return asynq.SkipRetry
	}

	pipeline, err := service.PipelineService.Get(scheduledRunPipelinePayload.PipelineID)

	if err != nil {
		log.Println(err.Error())
		return asynq.SkipRetry
	}

	run, err := service.Create(*pipeline)

	if err != nil {
		log.Println(err.Error())
		return asynq.SkipRetry
	}

	if pipelineSchedule.CronExpression != "" {
		parser := cron.NewParser(cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow)
		schedule, parseError := parser.Parse(pipelineSchedule.CronExpression)

		if parseError != nil {
			return parseError
		}

		nextExec := schedule.Next(time.Now())

		task, err := NewScheduledRunPipelineTask(pipelineSchedule.PipelineID, pipelineSchedule.ID)

		if err != nil {
			return err
		}

		if _, err = service.TaskQueueClient.Enqueue(task, asynq.Queue("runs"), asynq.Timeout(0), asynq.ProcessAt(nextExec)); err != nil {
			errMessage := service.I18n.MustLocalize(&i18n.LocalizeConfig{
				MessageID: "tasks.client.enqueue.failed",
				TemplateData: map[string]interface{}{
					"Queue":  "runs",
					"Reason": err.Error(),
				},
				PluralCount: 1,
			})

			return errors.New(errMessage)
		}
	}

	runPipelinePayload := &RunPipelinePayload{
		PipelineID:      pipeline.ID,
		RunID:           run.ID,
		GraphDefinition: pipeline.Definition,
	}

	return executeRunPipelineTask(*runPipelinePayload, service)
}

func (service *runServiceImpl) UpdateRunStatus(runID uint, statusID uint, stepWaitingFeedback int, errorMessage string) error {
	run, _ := service.RunRepository.FindByID(runID)
	runStatus, _ := service.RunRepository.GetRunStatusByID(statusID)
	run.RunStatusID = statusID
	run.RunStatus = *runStatus
	run.ErrorMessage = errorMessage
	run.StepWaitingFeedback = stepWaitingFeedback
	run.LastRun = time.Now()

	if err := service.RunRepository.Update(run); err != nil {
		errMessage := service.I18n.MustLocalize(&i18n.LocalizeConfig{
			MessageID: "run.repository.update.run.failed",
			TemplateData: map[string]interface{}{
				"ID":     runID,
				"Reason": err.Error(),
			},
			PluralCount: 1,
		})

		return errors.New(errMessage)
	}

	return nil
}

func (service *runServiceImpl) updateStepRunStatus(runStepStatus *model.RunStepStatus, statusID uint, errorMessage string) error {
	runStatus, _ := service.RunRepository.GetRunStatusByID(statusID)
	runStepStatus.RunStatusID = statusID
	runStepStatus.RunStatus = *runStatus
	runStepStatus.ErrorMessage = errorMessage
	runStepStatus.LastRun = time.Now()

	if err := service.RunRepository.UpdateRunStepStatus(runStepStatus); err != nil {
		errMessage := service.I18n.MustLocalize(&i18n.LocalizeConfig{
			MessageID: "run.repository.update.step-status.failed",
			TemplateData: map[string]interface{}{
				"ID":     runStepStatus.ID,
				"Reason": err.Error(),
			},
			PluralCount: 1,
		})

		return errors.New(errMessage)
	}

	return nil
}
