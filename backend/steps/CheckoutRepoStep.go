package steps

import (
	"fmt"
	"os"

	"github.com/go-git/go-git/v5"
)

type CheckoutRepoStep struct {
	RepoURL string `json:"repoURL"`
}

func (step *CheckoutRepoStep) Execute() {

	pipelinesWorkDir := os.Getenv("PIPELINES_WORK_DIR")
	if err := os.MkdirAll(pipelinesWorkDir, os.ModePerm); err != nil {

	}

	_, err := git.PlainClone(pipelinesWorkDir, false, &git.CloneOptions{
		URL:      step.RepoURL,
		Progress: os.Stdout,
	})

	fmt.Printf("\x1b[31;1m%s\x1b[0m\n", fmt.Sprintf("error: %s", err))
}
