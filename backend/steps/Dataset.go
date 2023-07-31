package steps

import (
	"di/model"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/nicksnyder/go-i18n/v2/i18n"
)

type Dataset struct {
	ID          int
	PipelineID  uint
	RunID       uint
	IsFirstStep bool
	Name        string
	DatasetName string
	DatasetID   uint
	Filepath    string
}

func (step Dataset) GetID() int {
	return int(step.ID)
}

func (step Dataset) GetName() string {
	return step.Name
}

func (step *Dataset) GetIsFirstStep() bool {
	return step.IsFirstStep
}

func (step *Dataset) SetData(stepDescription model.NodeDescription) error {
	step.ID, _ = strconv.Atoi(stepDescription.Data.ID)
	step.Name = stepDescription.Data.NameAndType.Name
	step.IsFirstStep = stepDescription.Data.NameAndType.IsFirstStep
	step.DatasetName = stepDescription.Data.StepConfig.DatasetName
	step.DatasetID = stepDescription.Data.StepConfig.DatasetID
	step.Filepath = stepDescription.Data.StepConfig.DatasetPath

	return nil
}

func (step *Dataset) SetPipelineID(pipelineID uint) error {
	step.PipelineID = pipelineID

	return nil
}

func (step *Dataset) SetRunID(runID uint) error {
	step.RunID = runID

	return nil
}

func (step *Dataset) GetPipelineID() uint {
	return step.PipelineID
}

func (step *Dataset) GetRunID() uint {
	return step.RunID
}

func (step *Dataset) GetIsStaggered() bool {
	return false
}

func (step Dataset) Execute(logFile *os.File, feedbackRects [][]model.HumanFeedbackRect, I18n *i18n.Localizer) ([]model.HumanFeedbackQueryPayload, error) {

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

	fileUploadDir, exists := os.LookupEnv("FILE_UPLOAD_DIR")

	if !exists {
		errMessage := I18n.MustLocalize(&i18n.LocalizeConfig{
			MessageID: "env.variable.find.failed",
			TemplateData: map[string]interface{}{
				"Name": "FILE_UPLOAD_DIR",
			},
			PluralCount: 1,
		})

		return nil, errors.New(errMessage)
	}

	relativeFilePath := strings.Split(step.Filepath, "/files")[1]
	path := filepath.Join(fileUploadDir, relativeFilePath)

	sourceFile, err := os.Open(path)

	if err != nil {
		errMessage := fmt.Sprintf("Error creating script file from file script: %v", err.Error())
		return nil, errors.New(errMessage)
	}

	defer sourceFile.Close()

	datasetsDir := filepath.Join(currentPipelineWorkDir, "datasets")

	if err := os.MkdirAll(datasetsDir, os.ModePerm); err != nil {
		errMessage := I18n.MustLocalize(&i18n.LocalizeConfig{
			MessageID: "os.cmd.mkdir.dir.failed",
			TemplateData: map[string]interface{}{
				"Path":   datasetsDir,
				"Reason": err.Error(),
			},
			PluralCount: 1,
		})

		return nil, errors.New(errMessage)
	}

	destinationFilePath := filepath.Join(datasetsDir, filepath.Base(step.Filepath))
	datasetFileDestination, err := os.Create(destinationFilePath)

	if err != nil {
		errMessage := fmt.Sprintf("Error creating script file from file script: %v", err.Error())
		return nil, errors.New(errMessage)
	}

	defer datasetFileDestination.Close()

	_, err = io.Copy(datasetFileDestination, sourceFile)

	if err != nil {
		errMessage := fmt.Sprintf("Error creating script file from file script: %v", err.Error())
		return nil, errors.New(errMessage)
	}

	return nil, err
}
