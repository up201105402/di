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

func (service *runServiceImpl) CreateRunStepStatus(runID uint, stepID int, runStatusID uint, errorMessage string) error {
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

	if err := service.UpdateRunStatus(runID, 2, ""); err != nil {
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

	runPipelineTask, err := service.NewRunPipelineTask(pipeline.ID, runID, pipeline.Definition, 0)

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

func (service *runServiceImpl) NewRunPipelineTask(pipelineID uint, runID uint, graph string, stepIndex uint) (*asynq.Task, error) {
	payload, err := json.Marshal(RunPipelinePayload{PipelineID: pipelineID, RunID: runID, GraphDefinition: graph, StepIndex: stepIndex})
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

	return executeRunPipelineTask(runPipelinePayload, service)
}

func executeRunPipelineTask(runPipelinePayload RunPipelinePayload, service *runServiceImpl) error {
	var stepDescriptions []model.NodeDescription

	if err := json.Unmarshal([]byte(runPipelinePayload.GraphDefinition), &stepDescriptions); err != nil {
		log.Printf("Unable to unmarshal pipeline %v definition", runPipelinePayload.PipelineID)
		return asynq.SkipRetry
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
		pipelineGraph.AddEdge(edge.GetTargetID(), edge.GetSourceID())
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

		if err := service.UpdateRunStatus(runPipelinePayload.RunID, 3, errMessage); err != nil {
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

		if err := service.UpdateRunStatus(runPipelinePayload.RunID, 3, errMessage); err != nil {
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

		if err := service.UpdateRunStatus(runPipelinePayload.RunID, 3, errMessage); err != nil {
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

		if err := service.UpdateRunStatus(runPipelinePayload.RunID, 3, errMessage); err != nil {
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

		if err := service.UpdateRunStatus(runPipelinePayload.RunID, 3, errMessage); err != nil {
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

		if err := service.UpdateRunStatus(runPipelinePayload.RunID, 3, errMessage); err != nil {
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

		service.UpdateRunStatus(runPipelinePayload.RunID, 3, errMessage)
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
		service.UpdateRunStatus(runPipelinePayload.RunID, 3, errMessage)
		return asynq.SkipRetry
	}

	runLogger := log.New(logFile, "", log.Ldate|log.Ltime|log.Lmicroseconds|log.Llongfile)

	msg := fmt.Sprintf("Executing run %d", runPipelinePayload.RunID)
	log.Println(msg)
	runLogger.Println(msg)

	hasError := false
	var stepErr error

	graph.BFS(pipelineGraph, firstStepID, func(id int) bool {
		step, _ := pipelineGraph.Vertex(id)

		msg = fmt.Sprintf("Executing step %s (%d) ...", step.GetName(), step.GetID())
		log.Println(msg)
		runLogger.Println(msg)

		runStepStatus := &model.RunStepStatus{RunID: runPipelinePayload.RunID, StepID: id, RunStatusID: 2, LastRun: time.Now()}
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

		if err := step.Execute(logFile, service.I18n); err != nil {

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

		if err := service.updateStepRunStatus(runStepStatus, 4, ""); err != nil {
			runLogger.Println(err.Error())
			log.Println(err.Error())
			stepErr = err
			hasError = true
			return true
		}

		return false
	})

	if hasError {
		if err := service.UpdateRunStatus(runPipelinePayload.RunID, 3, msg); err != nil {
			log.Println(err.Error())
			runLogger.Println(err.Error())
		}

		errMessage := service.I18n.MustLocalize(&i18n.LocalizeConfig{
			MessageID: "run.service.execute.run.failed",
			TemplateData: map[string]interface{}{
				"ID":     runPipelinePayload.RunID,
				"Reason": stepErr.Error(),
			},
			PluralCount: 1,
		})

		log.Println(errMessage)
		runLogger.Println(errMessage)
	} else {
		if err := service.UpdateRunStatus(runPipelinePayload.RunID, 4, ""); err != nil {
			log.Println(err.Error())
			runLogger.Println(err.Error())
		}

		errMessage := service.I18n.MustLocalize(&i18n.LocalizeConfig{
			MessageID: "run.service.execute.run.success",
			TemplateData: map[string]interface{}{
				"ID": runPipelinePayload.RunID,
			},
			PluralCount: 1,
		})

		log.Println(errMessage)
		runLogger.Println(errMessage)
	}

	logFile.Close()

	return asynq.SkipRetry
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
		StepIndex:       0,
	}

	return executeRunPipelineTask(*runPipelinePayload, service)
}

func (service *runServiceImpl) UpdateRunStatus(runID uint, statusID uint, errorMessage string) error {
	run, _ := service.RunRepository.FindByID(runID)
	runStatus, _ := service.RunRepository.GetRunStatusByID(statusID)
	run.RunStatusID = statusID
	run.RunStatus = *runStatus
	run.ErrorMessage = errorMessage
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
