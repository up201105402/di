package steps

import (
	"di/model"
	"fmt"
	"os"
	"strconv"

	"github.com/go-git/go-git/v5"
)

type CheckoutRepo struct {
	ID          int
	PipelineID  uint
	RunID       uint
	IsFirstStep bool
	RepoURL     string `json:"repoURL"`
}

func (step CheckoutRepo) GetID() int {
	return int(step.ID)
}

func (step *CheckoutRepo) GetIsFirstStep() bool {
	return step.IsFirstStep
}

func (step *CheckoutRepo) SetData(stepDescription model.NodeDescription) error {
	step.ID, _ = strconv.Atoi(stepDescription.ID)
	step.IsFirstStep = stepDescription.Data.NameAndType.IsFirstStep
	step.RepoURL = stepDescription.Data.StepConfig.RepoURL

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
