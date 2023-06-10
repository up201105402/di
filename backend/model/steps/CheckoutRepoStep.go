package steps

import (
	"fmt"
	"os"

	"github.com/go-git/go-git/v5"
)

type CheckoutRepoStep struct {
	ID         uint
	PipelineID uint
	RepoURL    string `json:"repoURL"`
}

func (step *CheckoutRepoStep) GetID() uint {
	return step.ID
}

func (step *CheckoutRepoStep) Execute() error {

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
