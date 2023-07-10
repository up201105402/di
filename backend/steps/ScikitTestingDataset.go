package steps

import (
	"di/model"
	"os"
)

type ScikitTestingDataset struct {
	ID         uint
	PipelineID uint
	RunID      uint
	Dataset    string
	DataConfig model.StepDataConfig
}

func (step ScikitTestingDataset) GetID() int {
	return int(step.ID)
}

func (step *ScikitTestingDataset) SetData(stepData model.StepData) error {
	step.Dataset = stepData.NameAndType.Dataset
	step.DataConfig = stepData.StepConfig
	return nil
}

func (step *ScikitTestingDataset) SetPipelineID(pipelineID uint) error {
	step.PipelineID = pipelineID

	return nil
}

func (step *ScikitTestingDataset) SetRunID(runID uint) error {
	step.RunID = runID

	return nil
}

func (step *ScikitTestingDataset) GetPipelineID() uint {
	return step.PipelineID
}

func (step *ScikitTestingDataset) GetRunID() uint {
	return step.RunID
}

func (step ScikitTestingDataset) Execute(logFile *os.File) error {

	logFile.WriteString("Executing...")

	return nil
}
