package steps

import (
	"di/model"
	"fmt"
	"os"

	"github.com/go-git/go-git/v5"
)

type CheckoutRepo struct {
	ID         uint
	PipelineID uint
	RunID      uint
	RepoURL    string `json:"repoURL"`
}

func (step CheckoutRepo) GetID() int {
	return int(step.ID)
}

func (step *CheckoutRepo) SetData(stepConfig model.StepData) error {
	step.RepoURL = stepConfig.NameAndType.Name

	return nil
}

func (step *CheckoutRepo) SetPipelineID(pipelineID uint) error {
	step.PipelineID = pipelineID

	return nil
}

func (step *CheckoutRepo) SetRunID(runID uint) error {
	step.RunID = runID

	return nil
}

func (step *CheckoutRepo) GetPipelineID() uint {
	return step.PipelineID
}

func (step *CheckoutRepo) GetRunID() uint {
	return step.RunID
}

func (step CheckoutRepo) Execute(logFile *os.File) error {

	pipelinesWorkDir := os.Getenv("PIPELINES_WORK_DIR")
	currentPipelineWorkDir := pipelinesWorkDir + "/" + fmt.Sprint(step.PipelineID) + "/" + fmt.Sprint(step.RunID) + "/"
	if err := os.MkdirAll(currentPipelineWorkDir, os.ModePerm); err != nil {
		return err
	}

	if _, err := git.PlainClone(currentPipelineWorkDir, false, &git.CloneOptions{
		URL:      step.RepoURL,
		Progress: logFile,
	}); err != nil {
		// if err == git.ErrRepositoryAlreadyExists {
		// 	return err
		// }

		return err
	}

	return nil
}
