package steps

import (
	"di/model"
	"os"
	"strconv"
)

type ScikitTrainingDataset struct {
	ID          int
	PipelineID  uint
	RunID       uint
	Dataset     string
	IsFirstStep bool
	DataConfig  model.StepDataConfig
}

func (step ScikitTrainingDataset) GetID() int {
	return int(step.ID)
}

func (step *ScikitTrainingDataset) GetIsFirstStep() bool {
	return step.IsFirstStep
}

func (step *ScikitTrainingDataset) SetData(stepDescription model.NodeDescription) error {
	step.ID, _ = strconv.Atoi(stepDescription.ID)
	step.IsFirstStep = stepDescription.Data.NameAndType.IsFirstStep
	step.Dataset = stepDescription.Data.NameAndType.Dataset
	step.DataConfig = stepDescription.Data.StepConfig
	return nil
}

func (step *ScikitTrainingDataset) SetPipelineID(pipelineID uint) error {
	step.PipelineID = pipelineID

	return nil
}

func (step *ScikitTrainingDataset) SetRunID(runID uint) error {
	step.RunID = runID

	return nil
}

func (step *ScikitTrainingDataset) GetPipelineID() uint {
	return step.PipelineID
}

func (step *ScikitTrainingDataset) GetRunID() uint {
	return step.RunID
}

func (step ScikitTrainingDataset) Execute(logFile *os.File) error {

	logFile.WriteString("Executing...")

	// cmd := exec.Command("python3",
	// 	"/usr/src/di/backend/scikit/python/datasets/load_dataset_from_csv.py",
	// 	"-f1", "backend/scikit/python/datasets/data.csv",
	// 	"-f2", "backend/scikit/python/datasets/target.csv",
	// 	"-d", "backend/scikit/python/datasets/filtered_data.csv")

	// if err := cmd.Run(); err != nil {
	// 	return err
	// }

	return nil
}
