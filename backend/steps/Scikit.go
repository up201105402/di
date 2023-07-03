package steps

import (
	"di/model"
	"di/steps/scikit"
	"fmt"
	"os"
	"reflect"
	"strings"
)

type Scikit struct {
	ID         uint
	PipelineID uint
	RunID      uint
	ScikitStep Step
}

func (step Scikit) GetID() int {
	return int(step.ID)
}

func (step *Scikit) SetConfig(stepConfig model.StepDataConfig) error {
	return step.ScikitStep.SetConfig(stepConfig)
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

	return step.ScikitStep.Execute(logFile)
}

func initScikitModelTypeRegistry() map[string]reflect.Type {
	var scikitModelTypeRegistry = make(map[string]reflect.Type)

	stepTypes := []interface{}{
		scikit.LeastSquares{},
		scikit.RidgeRegAndClassification{},
		scikit.Lasso{},
	}

	for _, v := range stepTypes {
		splitString := strings.SplitAfter(fmt.Sprintf("%T", v), ".")
		scikitModelTypeRegistry[splitString[len(splitString)-1]] = reflect.TypeOf(v)
	}

	return scikitModelTypeRegistry
}
