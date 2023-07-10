package steps

import (
	"di/model"
	"os"
	"strconv"
)

type ScikitTestingDataset struct {
	ID          int
	PipelineID  uint
	RunID       uint
	Dataset     string
	Name        string
	IsFirstStep bool
	DataConfig  model.StepDataConfig
}

func (step ScikitTestingDataset) GetID() int {
	return int(step.ID)
}

func (step ScikitTestingDataset) GetName() string {
	return step.Name
}

func (step *ScikitTestingDataset) GetIsFirstStep() bool {
	return step.IsFirstStep
}

func (step *ScikitTestingDataset) SetData(stepDescription model.NodeDescription) error {
	step.ID, _ = strconv.Atoi(stepDescription.ID)
	step.Name = stepDescription.Data.NameAndType.Name
	step.IsFirstStep = stepDescription.Data.NameAndType.IsFirstStep
	step.Dataset = stepDescription.Data.NameAndType.Dataset
	step.DataConfig = stepDescription.Data.StepConfig
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

	switch step.Dataset {
	case "scikitBreastCancer":

	case "scikitDiabetes":

	case "scikitDigits":

	case "scikitIris":

	case "scikitLinerrud":

	case "scikitWine":

	case "scikitLoadFile":

	}

	return nil
}
