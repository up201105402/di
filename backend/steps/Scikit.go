package steps

import (
	"di/model"
	"fmt"
	"os"
	"reflect"
	"strings"

	"github.com/go-git/go-git/v5"
)

type Scikit struct {
	ID         uint
	PipelineID uint
	RunID      uint
	RepoURL    string `json:"repoURL"`
}

func (step Scikit) GetID() int {
	return int(step.ID)
}

func (step *Scikit) SetConfig(stepConfig model.StepDataConfig) error {
	step.RepoURL = stepConfig.RepoURL

	return nil
}

func (step *Scikit) SetPipelineID(pipelineID uint) error {
	step.PipelineID = pipelineID

	return nil
}

func (step *Scikit) SetRunID(runID uint) error {
	step.RunID = runID

	return nil
}

func (step *Scikit) GetPipelineID() uint {
	return step.PipelineID
}

func (step *Scikit) GetRunID() uint {
	return step.RunID
}

func (step Scikit) Execute(logFile *os.File) error {

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

func initScikitModelTypeRegistry() map[string]reflect.Type {
	var scikitModelTypeRegistry = make(map[string]reflect.Type)

	stepTypes := []interface{}{steps.CheckoutRepo{}}

	for _, v := range stepTypes {
		splitString := strings.SplitAfter(fmt.Sprintf("%T", v), ".")
		scikitModelTypeRegistry[splitString[len(splitString)-1]] = reflect.TypeOf(v)
	}

	return scikitModelTypeRegistry
}
