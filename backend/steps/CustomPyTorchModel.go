package steps

import (
	"bytes"
	"di/model"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"github.com/nicksnyder/go-i18n/v2/i18n"
)

type CustomPyTorchModel struct {
	ID           int
	PipelineID   uint
	RunID        uint
	IsFirstStep  bool
	Name         string
	ScriptType   string
	InlineScript string
	Filename     string
}

func (step CustomPyTorchModel) GetID() int {
	return int(step.ID)
}

func (step CustomPyTorchModel) GetName() string {
	return step.Name
}

func (step *CustomPyTorchModel) GetIsFirstStep() bool {
	return step.IsFirstStep
}

func (step *CustomPyTorchModel) SetData(stepDescription model.NodeDescription) error {
	step.ID, _ = strconv.Atoi(stepDescription.Data.ID)
	step.Name = stepDescription.Data.NameAndType.Name
	step.IsFirstStep = stepDescription.Data.NameAndType.IsFirstStep
	step.ScriptType = stepDescription.Data.NameAndType.ScriptType
	step.InlineScript = stepDescription.Data.StepConfig.InlineScript.String
	step.Filename = stepDescription.Data.StepConfig.Filename.String

	return nil
}

func (step *CustomPyTorchModel) SetPipelineID(pipelineID uint) error {
	step.PipelineID = pipelineID

	return nil
}

func (step *CustomPyTorchModel) SetRunID(runID uint) error {
	step.RunID = runID

	return nil
}

func (step *CustomPyTorchModel) GetPipelineID() uint {
	return step.PipelineID
}

func (step *CustomPyTorchModel) GetRunID() uint {
	return step.RunID
}

func (step *CustomPyTorchModel) GetIsStaggered() bool {
	return false
}

func (step CustomPyTorchModel) Execute(logFile *os.File, feedbackRects [][]model.HumanFeedbackRect, I18n *i18n.Localizer) ([]model.HumanFeedbackQueryPayload, error) {

	runLogger := log.New(logFile, "", log.Ldate|log.Ltime|log.Lmicroseconds|log.Llongfile)

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
	customModelsDir := currentPipelineWorkDir + "custom_models/"

	if err := os.MkdirAll(customModelsDir, os.ModePerm); err != nil {
		errMessage := I18n.MustLocalize(&i18n.LocalizeConfig{
			MessageID: "os.cmd.mkdir.dir.failed",
			TemplateData: map[string]interface{}{
				"Path":   customModelsDir,
				"Reason": err.Error(),
			},
			PluralCount: 1,
		})

		return nil, errors.New(errMessage)
	}

	filename := "custom_model.py"

	if step.ScriptType == "inline" {
		customModelFile, err := os.Create(customModelsDir + filename)

		if err != nil {
			errMessage := fmt.Sprintf("Error creating script file from inline script: %v", err.Error())
			return nil, errors.New(errMessage)
		}

		defer customModelFile.Close()

		processedString := strings.ReplaceAll(step.InlineScript, "\u00a0", " ")
		customModelFile.WriteString(processedString)

	} else if step.ScriptType == "file" {
		fileUploadDir, exists := os.LookupEnv("FILE_UPLOAD_DIR")

		if !exists {
			errMessage := I18n.MustLocalize(&i18n.LocalizeConfig{
				MessageID: "env.variable.find.failed",
				TemplateData: map[string]interface{}{
					"Name": "FILE_UPLOAD_DIR",
				},
				PluralCount: 1,
			})

			log.Printf(errMessage)
			return nil, errors.New(errMessage)
		}

		sourceFile, err := os.Open(fileUploadDir + "pipelines/" + fmt.Sprint(step.PipelineID) + "/" + step.Filename)

		if err != nil {
			errMessage := fmt.Sprintf("Error creating script file from file script: %v", err.Error())
			return nil, errors.New(errMessage)
		}

		defer sourceFile.Close()

		customModelFile, err := os.Create(customModelsDir + filename)

		if err != nil {
			errMessage := fmt.Sprintf("Error creating script file from file script: %v", err.Error())
			return nil, errors.New(errMessage)
		}

		defer customModelFile.Close()

		_, err = io.Copy(customModelFile, sourceFile)

		if err != nil {
			errMessage := fmt.Sprintf("Error creating script file from file script: %v", err.Error())
			return nil, errors.New(errMessage)
		}
	} else {
		errMessage := "Invalid script type for shell script step!"
		return nil, errors.New(errMessage)
	}

	cmd := exec.Command("python3", filename)
	cmd.Dir = currentPipelineWorkDir
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	runLogger.Println(stderr.String())
	runLogger.Println(stdout.String())

	return nil, err
}
