package steps

import (
	"di/model"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/nicksnyder/go-i18n/v2/i18n"
	"gopkg.in/guregu/null.v4"
)

type Tester struct {
	ID              int
	PipelineID      uint
	RunID           uint
	IsFirstStep     bool
	Name            string
	TesterName      string
	TesterID        uint
	Filepath        string
	Epochs          null.Int
	IsStaggered     bool
	CustomArguments null.String
}

func (step Tester) GetID() int {
	return int(step.ID)
}

func (step Tester) GetName() string {
	return step.Name
}

func (step *Tester) GetIsFirstStep() bool {
	return step.IsFirstStep
}

func (step *Tester) SetData(stepDescription model.NodeDescription) error {
	step.ID, _ = strconv.Atoi(stepDescription.Data.ID)
	step.Name = stepDescription.Data.NameAndType.Name
	step.IsFirstStep = stepDescription.Data.NameAndType.IsFirstStep
	step.TesterName = stepDescription.Data.StepConfig.TesterName
	step.TesterID = stepDescription.Data.StepConfig.TesterID
	step.Filepath = stepDescription.Data.StepConfig.TesterPath
	step.Epochs = stepDescription.Data.StepConfig.Epochs
	step.Filepath = stepDescription.Data.StepConfig.TesterPath
	step.IsStaggered = stepDescription.Data.StepConfig.IsStaggered
	step.CustomArguments = stepDescription.Data.StepConfig.CustomArguments

	return nil
}

func (step *Tester) SetPipelineID(pipelineID uint) error {
	step.PipelineID = pipelineID

	return nil
}

func (step *Tester) SetRunID(runID uint) error {
	step.RunID = runID

	return nil
}

func (step *Tester) GetPipelineID() uint {
	return step.PipelineID
}

func (step *Tester) GetRunID() uint {
	return step.RunID
}

func (step *Tester) GetIsStaggered() bool {
	return true
}

func (step Tester) Execute(logFile *os.File, feedbackRects [][]model.HumanFeedbackRect, I18n *i18n.Localizer) ([]model.HumanFeedbackQueryPayload, error) {

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

	err := step.copyOriginalTesterFile(currentPipelineWorkDir, I18n)

	if err != nil {
		return nil, err
	}

	var args []string

	args = append(args, filepath.Base(step.Filepath))

	if step.CustomArguments.Valid {
		args = append(args, step.CustomArguments.String)
	}

	cmd := exec.Command("python3", args...)
	cmd.Dir = currentPipelineWorkDir
	cmd.Stdout = logFile
	cmd.Stderr = logFile

	cmdErr := cmd.Run()

	if cmdErr != nil {
		return nil, cmdErr
	}

	return nil, nil
}

func (step Tester) copyOriginalTesterFile(currentPipelineWorkDir string, I18n *i18n.Localizer) error {
	fileUploadDir, exists := os.LookupEnv("FILE_UPLOAD_DIR")

	if !exists {
		errMessage := I18n.MustLocalize(&i18n.LocalizeConfig{
			MessageID: "env.variable.find.failed",
			TemplateData: map[string]interface{}{
				"Name": "FILE_UPLOAD_DIR",
			},
			PluralCount: 1,
		})

		return errors.New(errMessage)
	}

	relativeFilePath := strings.Split(step.Filepath, "/files")[1]
	path := filepath.Join(fileUploadDir, relativeFilePath)

	sourceFile, err := os.Open(path)

	if err != nil {
		errMessage := fmt.Sprintf("Error creating script file from file script: %v", err.Error())
		return errors.New(errMessage)
	}

	defer sourceFile.Close()

	destinationFilePath := filepath.Join(currentPipelineWorkDir, filepath.Base(step.Filepath))
	testerFileDestination, err := os.Create(destinationFilePath)

	if err != nil {
		errMessage := fmt.Sprintf("Error creating script file from file script: %v", err.Error())
		return errors.New(errMessage)
	}

	defer testerFileDestination.Close()

	_, err = io.Copy(testerFileDestination, sourceFile)

	if err != nil {
		errMessage := fmt.Sprintf("Error creating script file from file script: %v", err.Error())
		return errors.New(errMessage)
	}

	return nil
}
