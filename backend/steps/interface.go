package steps

import (
	"di/model"
	"os"
)

type Step interface {
	GetID() int
	Execute(logFile *os.File) error
	SetConfig(stepConfig model.StepDataConfig) error
	SetPipelineID(pipelineID uint) error
	SetRunID(runID uint) error
	GetPipelineID() uint
	GetRunID() uint
}

type Edge interface {
	GetID() int
	GetNextStep() *Step
	GetPreviousStep() *Step
}
