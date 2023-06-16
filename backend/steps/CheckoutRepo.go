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
	RepoURL    string `json:"repoURL"`
}

func (step CheckoutRepo) GetID() int {
	return int(step.ID)
}

func (step *CheckoutRepo) SetConfig(stepConfig model.StepDataConfig) error {
	step.RepoURL = stepConfig.RepoURL

	return nil
}

func (step *CheckoutRepo) SetPipelineID(pipelineID uint) error {
	step.PipelineID = pipelineID

	return nil
}

func (step CheckoutRepo) Execute() error {

	pipelinesWorkDir := os.Getenv("PIPELINES_WORK_DIR")
	currentPipelineWorkDir := pipelinesWorkDir + "/" + fmt.Sprint(step.PipelineID)
	if err := os.MkdirAll(currentPipelineWorkDir, os.ModePerm); err != nil {
		return err
	}

	if _, err := git.PlainClone(pipelinesWorkDir, false, &git.CloneOptions{
		URL:      step.RepoURL,
		Progress: os.Stdout,
	}); err != nil {
		return err
	}

	return nil
}
