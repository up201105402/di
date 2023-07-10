package steps

import (
	"di/model"
	"os"
)

type ScikitTrainingDataset struct {
	ID         uint
	PipelineID uint
	RunID      uint
	Dataset    string
	DataConfig model.StepDataConfig
}

func (step ScikitTrainingDataset) GetID() int {
	return int(step.ID)
}

func (step *ScikitTrainingDataset) SetData(stepData model.StepData) error {
	step.Dataset = stepData.NameAndType.Dataset
	step.DataConfig = stepData.StepConfig
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
