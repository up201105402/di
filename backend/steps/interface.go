package steps

import (
	"di/model"
	"os"
)

type Step interface {
	GetID() int
	Execute(logFile *os.File) error
	SetData(stepDescription model.NodeDescription) error
	SetPipelineID(pipelineID uint) error
	SetRunID(runID uint) error
	GetPipelineID() uint
	GetRunID() uint
	GetIsFirstStep() bool
}

type Edge interface {
	SetData(stepDescription model.NodeDescription)
	GetSourceID() int
	GetTargetID() int
}
