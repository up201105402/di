package steps

import (
	"di/model"
	"di/util"
	"di/util/errors"
	"os"
)

var ScikitUnsupervisedModelTypes = []string{
	"leastSquares",
	"ridgeRegression",
	"ridgeRegressionCV",
	"ridgeClassifier",
	"ridgeClassifierCV",
	"lasso",
	"lassoCV",
	"lassoLars",
	"lassoLarsCV",
	"lassoLarsIC",
	"multiTaskLasso",
	"multiTaskLassoCV",
	"elasticNet",
	"elasticNetCV",
	"multiTaskElasticNet",
	"multiTaskElasticNetCV",
	"lars",
	"larsCV",
	"omp",
	"ompCV",
	"bayesianRidge",
	"bayesianARD",
	"logisticRegression",
	"logisticRegressionCV",
	"tweedieRegressor",
	"poissonRegressor",
	"gammaRegressor",
	"sgdClassifier",
	"sgdRegressor",
	"perceptron",
	"passiveAgressiveClassifier",
	"passiveAgressiveRegressor",
	"huberRegression",
	"ransacRegression",
	"theilSenRegression",
	"quantileRegression",
}

type ScikitUnsupervisedModel struct {
	ID         uint
	PipelineID uint
	RunID      uint
	Model      string
	DataConfig model.StepDataConfig
}

func (step ScikitUnsupervisedModel) GetID() int {
	return int(step.ID)
}

func (step *ScikitUnsupervisedModel) SetModel(model string) error {
	if util.StringArrayContains(ScikitUnsupervisedModelTypes, model) {
		step.Model = model
		return nil
	}

	return errors.NewNotFound("model", model)
}

func (step *ScikitUnsupervisedModel) SetConfig(stepConfig model.StepDataConfig) error {
	step.DataConfig = stepConfig
	return nil
}

func (step *ScikitUnsupervisedModel) SetPipelineID(pipelineID uint) error {
	step.PipelineID = pipelineID
	return nil
}

func (step *ScikitUnsupervisedModel) SetRunID(runID uint) error {
	step.RunID = runID

	return nil
}

func (step *ScikitUnsupervisedModel) GetPipelineID() uint {
	return step.PipelineID
}

func (step *ScikitUnsupervisedModel) GetRunID() uint {
	return step.RunID
}

func (step ScikitUnsupervisedModel) Execute(logFile *os.File) error {

	return nil
}
