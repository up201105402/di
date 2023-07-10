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
	"time"

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

func (service *runServiceImpl) CreateRunStepStatus(runID uint, stepID int, runStatusID uint, errorMessage string) error {
	newRunStepStatus := &model.RunStepStatus{RunID: runID, StepID: stepID, StatusID: runStatusID, LastRun: time.Now()}
	if err := service.RunRepository.CreateRunStepStatus(newRunStepStatus); err != nil {
		return err
	}

	return nil
}

func (service *runServiceImpl) Execute(runID uint) error {
	// demarshal stringified pipeline definition json

	run, err := service.RunRepository.FindByID(runID)

	if err != nil {
		log.Printf("Could not retrieve run with id %v", runID)
		return err
	}

	service.UpdateRunStatus(runID, 2, "")

	err = service.RunRepository.DeleteAllRunStepStatuses(runID)

	if err != nil {
		log.Printf("Could not delete run step statuses for run with id %v", runID)
		return err
	}

	pipeline, err := service.PipelineService.Get(run.PipelineID)

	if err != nil {
		log.Printf("Could not retrieve pipeline with id %v", run.PipelineID)
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

func (service *runServiceImpl) UpdateRunStepStatus(runStepStatus *model.RunStepStatus) error {
	err := service.RunRepository.UpdateRunStepStatus(runStepStatus)

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

func (service *runServiceImpl) DeleteRunStepStatus(id uint) error {
	err := service.RunRepository.DeleteRunStepStatus(id)

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
	return asynq.NewTask(RunPipelineTask, payload, asynq.MaxRetry(0)), nil
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
		log.Printf("Unable to unmarshal pipeline %v definition", runPipelinePayload.PipelineID)
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
		pipelineGraph.AddEdge(edge.GetTargetID(), edge.GetSourceID())
	}

	pipelinesWorkDir := os.Getenv("PIPELINES_WORK_DIR")
	runLogskDir := os.Getenv("RUN_LOGS_DIR")

	if pipelinesWorkDir == "" {
		service.UpdateRunStatus(runPipelinePayload.RunID, 3, "PIPELINES_WORK_DIR is not defined!")
		return asynq.SkipRetry
	}

	currentPipelineWorkDir := pipelinesWorkDir + "/" + fmt.Sprint(runPipelinePayload.PipelineID) + "/" + fmt.Sprint(runPipelinePayload.RunID) + "/"
	currentRunLogDir := runLogskDir + "/" + fmt.Sprint(runPipelinePayload.PipelineID) + "/" + fmt.Sprint(runPipelinePayload.RunID) + "/"

	if err := os.RemoveAll(currentPipelineWorkDir); err != nil {
		service.UpdateRunStatus(runPipelinePayload.RunID, 3, "Error removing files from "+currentPipelineWorkDir)
		return asynq.SkipRetry
	}

	if err := os.MkdirAll(currentPipelineWorkDir, os.ModePerm); err != nil {
		service.UpdateRunStatus(runPipelinePayload.RunID, 3, "Error creating directory "+currentPipelineWorkDir)
		return asynq.SkipRetry
	}

	if err := os.RemoveAll(currentRunLogDir); err != nil {
		service.UpdateRunStatus(runPipelinePayload.RunID, 3, "Error removing files from "+currentPipelineWorkDir)
		return asynq.SkipRetry
	}

	if err := os.MkdirAll(currentRunLogDir, os.ModePerm); err != nil {
		service.UpdateRunStatus(runPipelinePayload.RunID, 3, "Error creating directory "+currentPipelineWorkDir)
		return asynq.SkipRetry
	}

	logFileName := "run.log"
	logFile, err := os.Create(currentRunLogDir + logFileName)

	if err != nil {
		log.Print(err.Error())
		service.UpdateRunStatus(runPipelinePayload.RunID, 3, "Error creating log file  "+currentRunLogDir+logFileName)
		return asynq.SkipRetry
	}

	runLogger := log.New(logFile, "", log.Ldate|log.Ltime|log.Lmicroseconds|log.Llongfile)

	msg := fmt.Sprintf("Executing run %d", runPipelinePayload.RunID)
	log.Println(msg)
	runLogger.Println(msg)

	hasError := false

	graph.BFS(pipelineGraph, firstStepID, func(id int) bool {
		step, _ := pipelineGraph.Vertex(id)

		msg = fmt.Sprintf("Executing step %s (%d) ...", step.GetName(), step.GetID())
		log.Println(msg)
		runLogger.Println(msg)

		service.RunRepository.CreateRunStepStatus(&model.RunStepStatus{RunID: runPipelinePayload.RunID, StepID: id, StatusID: 2, LastRun: time.Now()})

		if err != nil {
			msg = fmt.Sprintf("Error updating step %d status: %s", step.GetID(), err.Error())
			log.Println(msg)
			runLogger.Println(msg)

			service.RunRepository.UpdateRunStepStatus(&model.RunStepStatus{StepID: id, RunID: runPipelinePayload.RunID, StatusID: 3, ErrorMessage: err.Error()})

			hasError = true

			return true
		}

		if err := step.Execute(logFile); err != nil {

			msg = fmt.Sprintf("Error executing step %d: %s", step.GetID(), err.Error())
			runLogger.Println(msg)
			log.Println(msg)

			service.RunRepository.UpdateRunStepStatus(&model.RunStepStatus{StepID: id, RunID: runPipelinePayload.RunID, StatusID: 3, ErrorMessage: err.Error()})

			hasError = true

			return true
		} else {
			runLogger.Println(fmt.Sprintf("Step %s (%d) executed successfully!", step.GetName(), step.GetID()))
		}

		service.RunRepository.UpdateRunStepStatus(&model.RunStepStatus{StepID: id, RunID: runPipelinePayload.RunID, StatusID: 4})
		return false
	})

	if hasError {
		service.UpdateRunStatus(runPipelinePayload.RunID, 3, msg)
		msg := fmt.Sprintf("Execution of run %d failed!", runPipelinePayload.RunID)
		log.Println(msg)
		runLogger.Println(msg)
	} else {
		service.UpdateRunStatus(runPipelinePayload.RunID, 4, msg)
		msg := fmt.Sprintf("Execution of run %d successful!", runPipelinePayload.RunID)
		log.Println(msg)
		runLogger.Println(msg)
	}

	logFile.Close()

	return asynq.SkipRetry
}

func (service *runServiceImpl) UpdateRunStatus(runID uint, statusID uint, errorMessage string) {
	run, _ := service.RunRepository.FindByID(runID)
	run.StatusID = statusID
	run.ErrorMessage = errorMessage
	run.LastRun = time.Now()
	service.RunRepository.Update(run)
}
