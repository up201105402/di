package steps

import (
	"di/model"
	"os"
)

type ScikitTrainingModel struct {
	ID         uint
	PipelineID uint
	RunID      uint
	DataConfig model.StepDataConfig
}

func (step ScikitTrainingModel) GetID() int {
	return int(step.ID)
}

func (step *ScikitTrainingModel) SetConfig(stepConfig model.StepDataConfig) error {
	step.DataConfig = stepConfig
	return nil
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
