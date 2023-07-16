package steps

import (
	"bytes"
	"di/model"
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strconv"
	"time"
)

type PythonScript struct {
	ID           int
	PipelineID   uint
	RunID        uint
	IsFirstStep  bool
	Name         string
	ScriptType   string
	InlineScript string
	Filename     string
}

func (step PythonScript) GetID() int {
	return int(step.ID)
}

func (step PythonScript) GetName() string {
	return step.Name
}

func (step *PythonScript) GetIsFirstStep() bool {
	return step.IsFirstStep
}

func (step *PythonScript) SetData(stepDescription model.NodeDescription) error {
	step.ID, _ = strconv.Atoi(stepDescription.ID)
	step.Name = stepDescription.Data.NameAndType.Name
	step.IsFirstStep = stepDescription.Data.NameAndType.IsFirstStep
	step.ScriptType = stepDescription.Data.NameAndType.ScriptType
	step.InlineScript = stepDescription.Data.StepConfig.InlineScript.String
	step.Filename = stepDescription.Data.StepConfig.Filename.String

	return nil
}

func (step *PythonScript) SetPipelineID(pipelineID uint) error {
	step.PipelineID = pipelineID

	return nil
}

func (step *PythonScript) SetRunID(runID uint) error {
	step.RunID = runID

	return nil
}

func (step *PythonScript) GetPipelineID() uint {
	return step.PipelineID
}

func (step *PythonScript) GetRunID() uint {
	return step.RunID
}

func (step PythonScript) Execute(logFile *os.File) error {

	runLogger := log.New(logFile, "", log.Ldate|log.Ltime|log.Lmicroseconds|log.Llongfile)

	pipelinesWorkDir, exists := os.LookupEnv("PIPELINES_WORK_DIR")

	if !exists {
		errMessage := "Pipelines work directory is not defined!"
		return errors.New(errMessage)
	}

	currentPipelineWorkDir := pipelinesWorkDir + "/" + fmt.Sprint(step.PipelineID) + "/" + fmt.Sprint(step.RunID) + "/"

	filename := step.Filename

	if step.ScriptType == "inline" {
		filename = pipelinesWorkDir + "file_" + fmt.Sprintf("%d", time.Now().Unix())
		scriptFile, err := os.Create(filename)
		defer scriptFile.Close()

		if err != nil {
			errMessage := fmt.Sprintf("Error creating script file from inline script: %v", err.Error())
			return errors.New(errMessage)
		}

		scriptFile.WriteString(step.InlineScript)

	} else if step.ScriptType != "file" {
		errMessage := "Invalid script type for shell script step!"
		return errors.New(errMessage)
	}

	cmd := exec.Command("python3", filename)
	cmd.Dir = currentPipelineWorkDir
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	runLogger.Println(stderr.String())
	runLogger.Println(stdout.String())

	return err
}