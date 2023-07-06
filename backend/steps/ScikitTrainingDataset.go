package steps

import (
	"di/model"
	"fmt"
	"os"
)

type ScikitTrainingModel struct {
	ID         uint
	PipelineID uint
	RunID      uint
	ScikitStep Step
}

func (step ScikitTrainingModel) GetID() int {
	return int(step.ID)
}

func (step *ScikitTrainingModel) SetConfig(stepConfig model.StepDataConfig) error {
	return step.ScikitStep.SetConfig(stepConfig)
}

func (step *ScikitTrainingModel) SetPipelineID(pipelineID uint) error {
	step.PipelineID = pipelineID

	return nil
}

func (step *ScikitTrainingModel) SetRunID(runID uint) error {
	step.RunID = runID

	return nil
}

func (step *ScikitTrainingModel) GetPipelineID() uint {
	return step.PipelineID
}

func (step *ScikitTrainingModel) GetRunID() uint {
	return step.RunID
}

func (step ScikitTrainingModel) Execute(logFile *os.File) error {

	pipelinesWorkDir := os.Getenv("PIPELINES_WORK_DIR")
	currentPipelineWorkDir := pipelinesWorkDir + "/" + fmt.Sprint(step.PipelineID) + "/" + fmt.Sprint(step.RunID) + "/"
	if err := os.MkdirAll(currentPipelineWorkDir, os.ModePerm); err != nil {
		return err
	}

	return step.ScikitStep.Execute(logFile)
}
