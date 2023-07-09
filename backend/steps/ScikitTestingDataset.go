package steps

import (
	"di/model"
	"os"
)

type ScikitTestingModel struct {
	ID         uint
	PipelineID uint
	RunID      uint
	DataConfig model.StepDataConfig
}

func (step ScikitTestingModel) GetID() int {
	return int(step.ID)
}

func (step *ScikitTestingModel) SetConfig(stepConfig model.StepDataConfig) error {
	step.DataConfig = stepConfig

	return nil
}

func (step *ScikitTestingModel) SetPipelineID(pipelineID uint) error {
	step.PipelineID = pipelineID

	return nil
}

func (step *ScikitTestingModel) SetRunID(runID uint) error {
	step.RunID = runID

	return nil
}

func (step *ScikitTestingModel) GetPipelineID() uint {
	return step.PipelineID
}

func (step *ScikitTestingModel) GetRunID() uint {
	return step.RunID
}

func (step ScikitTestingModel) Execute(logFile *os.File) error {

	return nil
}
