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
	"time"

	"github.com/nicksnyder/go-i18n/v2/i18n"
)

type ShellScript struct {
	ID           int
	PipelineID   uint
	RunID        uint
	IsFirstStep  bool
	Name         string
	ScriptType   string
	InlineScript string
	Filename     string
}

func (step ShellScript) GetID() int {
	return int(step.ID)
}

func (step ShellScript) GetName() string {
	return step.Name
}

func (step *ShellScript) GetIsFirstStep() bool {
	return step.IsFirstStep
}

func (step *ShellScript) SetData(stepDescription model.NodeDescription) error {
	step.ID, _ = strconv.Atoi(stepDescription.Data.ID)
	step.Name = stepDescription.Data.NameAndType.Name
	step.IsFirstStep = stepDescription.Data.NameAndType.IsFirstStep
	step.ScriptType = stepDescription.Data.NameAndType.ScriptType
	step.InlineScript = stepDescription.Data.StepConfig.InlineScript.String
	step.Filename = stepDescription.Data.StepConfig.Filename.String

	return nil
}

func (step *ShellScript) SetPipelineID(pipelineID uint) error {
	step.PipelineID = pipelineID

	return nil
}

func (step *ShellScript) SetRunID(runID uint) error {
	step.RunID = runID

	return nil
}

func (step *ShellScript) GetPipelineID() uint {
	return step.PipelineID
}

func (step *ShellScript) GetRunID() uint {
	return step.RunID
}

func (step *ShellScript) GetIsStaggered() bool {
	return false
}

func (step ShellScript) Execute(logFile *os.File, feedbackRects [][]model.HumanFeedbackRect, I18n *i18n.Localizer) ([]model.HumanFeedbackQueryPayload, error) {

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

	filename := step.Filename

	if step.ScriptType == "inline" {
		filename = pipelinesWorkDir + "file_" + fmt.Sprintf("%d", time.Now().Unix())
		scriptFile, err := os.Create(filename)
		defer scriptFile.Close()

		if err != nil {
			errMessage := fmt.Sprintf("Error creating script file from inline script: %v", err.Error())
			return nil, errors.New(errMessage)
		}

		processedString := strings.ReplaceAll(step.InlineScript, "\u00a0", " ")
		scriptFile.WriteString(processedString)

	} else if step.ScriptType != "file" {
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

		shellScriptFile, err := os.Create(currentPipelineWorkDir + filename)

		if err != nil {
			errMessage := fmt.Sprintf("Error creating script file from file script: %v", err.Error())
			return nil, errors.New(errMessage)
		}

		defer shellScriptFile.Close()

		_, err = io.Copy(shellScriptFile, sourceFile)

		if err != nil {
			errMessage := fmt.Sprintf("Error creating script file from file script: %v", err.Error())
			return nil, errors.New(errMessage)
		}
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
