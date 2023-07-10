package model

import (
	"time"

	"gorm.io/gorm"
)

type RunStatus struct {
	gorm.Model
	Name        string
	IsFinal     bool
	Description string
}

type Run struct {
	gorm.Model
	PipelineID   uint `json:"pipelineId"`
	Pipeline     Pipeline
	RunStatusID  uint
	RunStatus    RunStatus
	ErrorMessage string
	Definition   string
	LastRun      time.Time
}

type RunStepStatus struct {
	gorm.Model
	StepID       int
	RunID        uint
	Run          Run
	RunStatusID  uint
	RunStatus    RunStatus
	ErrorMessage string
	LastRun      time.Time
}

type CreateRunReq struct {
	Execute bool `json:"execute"`
}

type ExecuteRunReq struct {
	ID uint `json:"id"`
}
