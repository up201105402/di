package steps

import "di/model"

type Step interface {
	GetID() int
	Execute() error
	SetConfig(stepConfig model.StepDataConfig) error
	SetPipelineID(pipelineID uint) error
}

type Edge interface {
	GetID() int
	GetNextStep() *Step
	GetPreviousStep() *Step
}
