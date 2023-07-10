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
	PipelineID  uint `json:"pipelineId"`
	Pipeline    Pipeline
	StatusID    uint
	Status      RunStatus
	Description string
	LastRun     time.Time
}

type CreateRunReq struct {
	Execute bool `json:"execute"`
}

type ExecuteRunReq struct {
	ID uint `json:"id"`
}
