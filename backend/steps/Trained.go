package steps

import (
	"di/model"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"

	"github.com/nicksnyder/go-i18n/v2/i18n"
)

type Trained struct {
	ID          int
	PipelineID  uint
	RunID       uint
	IsFirstStep bool
	Name        string
	TrainedName string
	TrainedID   uint
	Filepath    string
}

func (step Trained) GetID() int {
	return int(step.ID)
}

func (step Trained) GetName() string {
	return step.Name
}

func (step *Trained) GetIsFirstStep() bool {
	return step.IsFirstStep
}

func (step *Trained) SetData(stepDescription model.NodeDescription) error {
	step.ID, _ = strconv.Atoi(stepDescription.ID)
	step.Name = stepDescription.Data.NameAndType.Name
	step.IsFirstStep = stepDescription.Data.NameAndType.IsFirstStep
	step.TrainedName = stepDescription.Data.StepConfig.TrainedName
	step.TrainedID = stepDescription.Data.StepConfig.TrainedID
	step.Filepath = stepDescription.Data.StepConfig.TrainedPath

	return nil
}

func (step *Trained) SetPipelineID(pipelineID uint) error {
	step.PipelineID = pipelineID

	return nil
}

func (step *Trained) SetRunID(runID uint) error {
	step.RunID = runID

	return nil
}

func (step *Trained) GetPipelineID() uint {
	return step.PipelineID
}

func (step *Trained) GetRunID() uint {
	return step.RunID
}

func (step *Trained) GetIsStaggered() bool {
	return false
}

func (step Trained) Execute(logFile *os.File, feedbackRects [][]model.HumanFeedbackRect, I18n *i18n.Localizer) ([]model.HumanFeedbackQueryPayload, error) {

	pipelinesWorkDir, exists := os.LookupEnv("PIPELINES_WORK_DIR")

	if !exists {
		errMessage := I18n.MustLocalize(&i18n.LocalizeConfig{
			MessageID: "env.variable.find.failed",
			TemplateData: map[string]interface{}{
				"Name": "PIPELINES_WORK_DIR",
			},
			PluralCount: 1,
		})

		return nil, errors.New(errMessage)
	}

	currentPipelineWorkDir := pipelinesWorkDir + "/" + fmt.Sprint(step.PipelineID) + "/" + fmt.Sprint(step.RunID) + "/"
	sourceFile, err := os.Open(step.Filepath)

	if err != nil {
		errMessage := fmt.Sprintf("Error creating script file from file script: %v", err.Error())
		return nil, errors.New(errMessage)
	}

	defer sourceFile.Close()

	destinationFilePath := filepath.Join(currentPipelineWorkDir, "trained", filepath.Base(step.Filepath))
	trainedFileDestination, err := os.Create(destinationFilePath)

	if err != nil {
		errMessage := fmt.Sprintf("Error creating script file from file script: %v", err.Error())
		return nil, errors.New(errMessage)
	}

	defer trainedFileDestination.Close()

	_, err = io.Copy(trainedFileDestination, sourceFile)

	if err != nil {
		errMessage := fmt.Sprintf("Error creating script file from file script: %v", err.Error())
		return nil, errors.New(errMessage)
	}

	return nil, err
}
